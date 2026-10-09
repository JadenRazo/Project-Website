import assert from "node:assert/strict";

const engines = await import(process.env.PLAYWRIGHT_MODULE || "playwright");
const engine = process.env.ENGINE || "chromium";
const base = process.env.PREVIEW_URL || "http://127.0.0.1:4190";
const mediaPattern = /\.(mp4|webm)(?:\?|$)/;
const browser = await engines[engine].launch({
  headless: true,
  args: engine === "chromium" ? ["--no-sandbox", "--disable-dev-shm-usage"] : [],
});
const playing = (page, minimum = 0.2) => page.waitForFunction(
  (minimum) => {
    const video = document.querySelector("video");
    return video && !video.paused && video.currentTime > minimum;
  }, minimum, { timeout: 15000 },
);
const open = async (width) => {
  const page = await browser.newPage({ viewport: { width, height: 1000 } });
  await page.route("**/api/**", (route) => route.fulfill({
    status: 200, contentType: "application/json", body: "{}",
  }));
  // A positive capability hint does not guarantee that the decoder works.
  await page.addInitScript(() => {
    const nativeCanPlayType = HTMLMediaElement.prototype.canPlayType;
    HTMLMediaElement.prototype.canPlayType = function (type) {
      return type.includes("video/mp4")
        ? "probably"
        : nativeCanPlayType.call(this, type);
    };
  });
  await page.goto(base, { waitUntil: "networkidle" });
  return page;
};

try {
  for (const width of [412, 1440]) {
    const page = await open(width);
    const requests = [];
    page.on("request", (request) => {
      if (mediaPattern.test(request.url())) requests.push(request.url());
    });
    await page.route(/\.mp4(?:\?|$)/, (route) => route.abort("failed"));
    const cards = page.locator(".project-card");
    for (let i = 0; i < 6; i++) {
      const card = cards.nth(i);
      await card.locator(".film-frame").scrollIntoViewIfNeeded();
      // A chapter selection must survive replacing a failed source.
      if (i === 0) await card.getByRole("button", { name: /Play chapter 3/ }).click();
      else await card.locator(".film-play").click();
      await playing(page, i === 0 ? 16 : 0.2);
      assert.equal(await page.locator("video").count(), 1);
      const source = await page.locator("video").getAttribute("src");
      assert.ok(source.endsWith(".webm"), "Failed MP4 switches to WebM");
      assert.ok(source.includes(width < 640 ? "portrait" : "wide"));
      assert.equal(await card.locator(".film-recovery").count(), 0);
    }
    assert.equal(new Set(requests).size, 12, "One primary and one fallback per film");
    await page.close();
    console.log(engine, width, "all six films recover; chapter seek preserved");
  }

  const page = await open(412);
  const requests = [];
  page.on("request", (request) => {
    if (mediaPattern.test(request.url())) requests.push(request.url());
  });
  await page.route(mediaPattern, (route) => route.abort("failed"));
  const card = page.locator(".project-card").first();
  await card.locator(".film-frame").scrollIntoViewIfNeeded();
  await card.locator(".film-play").click();
  await card.getByText("The video couldn’t load.").waitFor();
  assert.equal(new Set(requests).size, 2, "Exhausted formats do not loop");
  assert.ok((await card.locator(".film-recovery a").getAttribute("href")).endsWith(".webm"));
  await page.unroute(mediaPattern);
  await card.getByRole("button", { name: "Try again" }).click();
  await playing(page);
  await page.close();
  console.log(engine, "both formats failing still offers a working retry");

  if (engine === "chromium") {
    const page = await open(1440);
    const card = page.locator(".project-card").first();
    await card.locator(".film-frame").scrollIntoViewIfNeeded();
    await card.locator(".film-play").click();
    await playing(page);
    console.log(engine, "playback survives an inaccurate MP4 capability hint");
    await page.close();
  }
} finally {
  await browser.close();
}
