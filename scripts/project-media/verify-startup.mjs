import assert from "node:assert/strict";

const engines = await import(process.env.PLAYWRIGHT_MODULE || "playwright");
const engine = process.env.ENGINE || "chromium";
const base = process.env.PREVIEW_URL || "http://127.0.0.1:4190";
const browser = await engines[engine].launch({
  headless: true,
  args: engine === "chromium" ? ["--no-sandbox", "--disable-dev-shm-usage"] : [],
});
const playing = (page) => page.waitForFunction(() => {
  const video = document.querySelector("video");
  return video && !video.paused && !video.error && video.currentTime > 0.2;
}, null, { timeout: 15000 });
const open = async (width, rejectStartup = false) => {
  const page = await browser.newPage({ viewport: { width, height: 1000 } });
  await page.route("**/api/**", (route) => route.fulfill({
    status: 200, contentType: "application/json", body: "{}",
  }));
  await page.addInitScript((rejectStartup) => {
    window.filmPlayCalls = [];
    const nativePlay = HTMLMediaElement.prototype.play;
    let rejectNext = rejectStartup;
    HTMLMediaElement.prototype.play = function () {
      window.filmPlayCalls.push({
        visible: getComputedStyle(this).visibility === "visible",
        controls: this.controls,
        userGesture: navigator.userActivation.isActive,
      });
      if (rejectNext) {
        rejectNext = false;
        return Promise.reject(new DOMException("Playback is not ready yet", "InvalidStateError"));
      }
      return nativePlay.call(this);
    };
  }, rejectStartup);
  await page.goto(base, { waitUntil: "networkidle" });
  return page;
};

try {
  for (const width of [390, 1440]) {
    const page = await open(width);
    const cards = page.locator(".project-card");
    for (let index = 0; index < 6; index++) {
      await cards.nth(index).locator(".film-play").click();
      await playing(page);
    }
    const calls = await page.evaluate(() => window.filmPlayCalls);
    assert.equal(calls.length, 6, "Each first click starts one playback request");
    for (const call of calls) {
      assert.equal(call.visible, true, "Player is visible before requesting playback");
      assert.equal(call.controls, true, "Native controls are available before playback");
      assert.equal(call.userGesture, true, "Playback starts within the click gesture");
    }
    await page.close();
    console.log(engine, width, "all six first clicks start visible, interactive players");
  }

  const page = await open(390, true);
  const requests = [];
  page.on("request", request => {
    if (/\.(mp4|webm)(?:\?|$)/.test(request.url())) requests.push(request.url());
  });
  const card = page.locator(".project-card").first();
  await card.locator(".film-play").click();
  await page.waitForTimeout(500);
  assert.equal(await card.locator(".film-recovery").count(), 0, "A startup rejection is not a file-loading error");
  assert.equal(await card.locator(".film-loading").count(), 0);
  assert.equal(await page.locator("video").evaluate(video => video.controls), true);
  assert.equal(new Set(requests).size, 1, "A valid source is retained after startup rejection");
  const source = await page.locator("video").getAttribute("src");
  // Resume within the same page instead of sending the visitor to a new tab.
  await card.getByRole("button", { name: /Play chapter 1/ }).click();
  await playing(page);
  assert.equal(await page.locator("video").getAttribute("src"), source);
  assert.equal(new Set(requests).size, 1);
  await page.close();
  console.log(engine, "startup rejection keeps working inline controls and the original source");
} finally {
  await browser.close();
}
