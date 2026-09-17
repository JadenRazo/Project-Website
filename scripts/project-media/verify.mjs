import assert from "node:assert/strict";
import { mkdir, mkdtemp, writeFile } from "node:fs/promises";
import { join } from "node:path";
import { tmpdir } from "node:os";

const engines = await import(process.env.PLAYWRIGHT_MODULE || "playwright");
const { default: AxeBuilder } = await import(
  process.env.AXE_MODULE || "@axe-core/playwright"
);
const engine = process.env.ENGINE || "chromium";
const base = process.env.PREVIEW_URL || "http://127.0.0.1:4190";
const out =
  process.env.EVIDENCE_DIR ||
  (await mkdtemp(join(tmpdir(), "project-media-qa-")));
const mediaPattern = /\.(mp4|webm)(?:\?|$)/;
await mkdir(out, { recursive: true, mode: 0o700 });
console.log("Browser evidence:", out);
const results = [];
const browser = await engines[engine].launch({
  headless: true,
  args:
    engine === "chromium" ? ["--no-sandbox", "--disable-dev-shm-usage"] : [],
});
const isolate = (page) =>
  page.route("**/api/**", (route) =>
    route.fulfill({ status: 200, contentType: "application/json", body: "{}" }),
  );
const playing = (page) =>
  page.waitForFunction(
    () => {
      const v = document.querySelector("video");
      return v && !v.paused && v.currentTime > 0.2;
    },
    {},
    { timeout: 15000 },
  );
const hideDocument = (page, hidden) =>
  page.evaluate((hidden) => {
    // Deterministically exercise the visibility lifecycle in a headless browser.
    Object.defineProperty(document, "hidden", {
      configurable: true,
      value: hidden,
    });
    document.dispatchEvent(new Event("visibilitychange"));
  }, hidden);

try {
  for (const width of [375, 412, 430, 768, 1440]) {
    const context = await browser.newContext({
      viewport: { width, height: 900 },
      reducedMotion: "reduce",
    });
    const page = await context.newPage(),
      media = [],
      errors = [];
    await isolate(page);
    page.on("request", (r) => {
      if (mediaPattern.test(r.url())) media.push(r.url());
    });
    page.on("pageerror", (e) => errors.push(e.message));
    await page.goto(base, { waitUntil: "networkidle" });
    const cards = page.locator(".project-card");
    const first = cards.first();
    assert.equal(await cards.count(), 6, "Six existing project cards");
    assert.equal(
      await page.getByRole("tab").count(),
      0,
      "No separate tabbed video section",
    );
    assert.equal(
      await page.locator("video").count(),
      0,
      "No mounted player until intent",
    );
    assert.equal(
      await page
        .locator("#contact")
        .getByRole("link", { name: "LinkedIn", exact: true })
        .getAttribute("href"),
      "https://www.linkedin.com/in/JadenRazo",
    );
    assert.equal(
      await page.evaluate(
        () => document.documentElement.scrollWidth <= innerWidth,
      ),
      true,
      `No overflow at ${width}`,
    );
    // Follow the actual homepage CTA instead of jumping straight to the player.
    await page
      .getByRole("link", { name: "Watch the projects", exact: true })
      .click();
    await page.waitForTimeout(200);
    const playBox = await first.locator(".film-play").boundingBox();
    assert.ok(
      playBox.y >= 70 && playBox.y + playBox.height <= 900,
      "Play control is visible after the homepage CTA",
    );
    await page.screenshot({ path: `${out}/${engine}-${width}-showcase.png` });
    await first.screenshot({ path: `${out}/${engine}-${width}-card.png` });
    for (let i = 0; i < 6; i++) await cards.nth(i).scrollIntoViewIfNeeded();
    assert.equal(
      media.length,
      0,
      "Browsing all six cards must not download video",
    );
    const axe = await new AxeBuilder({ page })
      .include("#projects")
      .withTags(["wcag2a", "wcag2aa", "wcag21aa"])
      .analyze();
    assert.deepEqual(
      axe.violations.map((v) => ({
        id: v.id,
        impact: v.impact,
        nodes: v.nodes.map((n) => n.target),
      })),
      [],
      `Axe at ${width}`,
    );
    await first.locator(".film-frame").scrollIntoViewIfNeeded();
    await first.locator(".film-play").focus();
    await page.keyboard.press("Enter");
    await playing(page);
    const video = page.locator("video");
    assert.equal(await video.count(), 1);
    assert.equal(await video.evaluate((v) => Math.round(v.duration)), 24);
    assert.equal(
      await video.evaluate((v) => document.activeElement === v),
      true,
      "Focus follows the disappearing play control",
    );
    const source = await video.getAttribute("src");
    assert.ok(
      source.includes(width < 640 ? "portrait" : "wide"),
      "Viewport-specific composition",
    );
    assert.equal(
      new Set(media).size,
      1,
      "Only one composition/format downloaded",
    );
    await first.getByRole("button", { name: /Play chapter 3/ }).click();
    await page.waitForFunction(
      () => document.querySelector("video")?.currentTime >= 16,
    );
    await page.evaluate(() => window.scrollTo({ top: 0, behavior: "instant" }));
    await page.waitForFunction(() => document.querySelector("video")?.paused);
    await first.locator(".film-frame").scrollIntoViewIfNeeded();
    await page.waitForTimeout(150);
    assert.equal(
      await video.evaluate((v) => v.paused),
      true,
      "Returning to a card does not resume motion",
    );
    if (width === 375) {
      const box = await video.boundingBox();
      await page.mouse.move(box.x + 30, box.y + box.height - 22);
      await page.mouse.click(
        box.x + 30,
        box.y + box.height - (engine === "webkit" ? 22 : 48),
      );
      await playing(page);
      await page.waitForFunction(
        () => !document.querySelector(".film-loading"),
      );
    }
    const detached = await video.elementHandle();
    const second = cards.nth(1);
    await second.locator(".film-frame").scrollIntoViewIfNeeded();
    await second.locator(".film-play").click();
    await playing(page);
    assert.equal(
      await page.locator("video").count(),
      1,
      "Only one active video across all six cards",
    );
    assert.deepEqual(
      await detached.evaluate((v) => ({
        connected: v.isConnected,
        paused: v.paused,
        src: v.getAttribute("src"),
      })),
      { connected: false, paused: true, src: null },
      "Starting another card releases the previous player",
    );
    assert.equal(
      await first.locator(".film-play").count(),
      1,
      "Previous card returns to its poster",
    );
    await second.locator(".film-transcript summary").focus();
    await page.keyboard.press("Enter");
    assert.equal(
      await second.locator(".film-transcript").getAttribute("open"),
      "",
    );
    assert.deepEqual(errors, [], "No page exceptions");
    results.push({
      width,
      axeViolations: 0,
      noInitialVideoRequests: true,
      source,
      chapterSeek: true,
      offscreenPause: true,
      keyboard: true,
      focus: true,
      overflow: false,
      heroLinkShowsPlay: true,
      videoInsideEveryCard: true,
      oneActivePlayer: true,
      selectionCleanup: true,
    });
    console.log(engine, width, "passed");
    await context.close();
  }

  // Decode every chapter in all six actual card slots, including both formats.
  for (const width of [412, 1440]) {
    const page = await browser.newPage({ viewport: { width, height: 1000 } });
    await isolate(page);
    await page.goto(base, { waitUntil: "networkidle" });
    for (let i = 0; i < 6; i++) {
      const card = page.locator(".project-card").nth(i);
      await card.locator(".film-frame").scrollIntoViewIfNeeded();
      await card.locator(".film-play").click();
      await playing(page);
      assert.equal(await page.locator("video").count(), 1);
      for (const t of [3, 11, 19]) {
        await page.locator("video").evaluate(async (v, t) => {
          v.pause();
          await new Promise((resolve) => {
            v.addEventListener("seeked", resolve, { once: true });
            v.currentTime = t;
          });
        }, t);
        await card
          .locator(".film-frame")
          .screenshot({
            path: `${out}/${engine}-${width}-film-${i + 1}-${t}.png`,
          });
      }
    }
    results.push({ width, allSixFilmsDecode: true, allChaptersSeek: true });
    await page.close();
  }

  const page = await browser.newPage({
    viewport: { width: 412, height: 1000 },
  });
  await isolate(page);
  await page.goto(base, { waitUntil: "networkidle" });
  const first = page.locator(".project-card").first();
  await first.locator(".film-frame").scrollIntoViewIfNeeded();
  await first.locator(".film-play").click();
  await playing(page);
  await page.setViewportSize({ width: 800, height: 1000 });
  const ratio = await first
    .locator(".film-frame")
    .evaluate(
      (e) =>
        e.clientWidth /
        (e.clientHeight - parseFloat(getComputedStyle(e).paddingBottom)),
    );
  assert.ok(
    Math.abs(ratio - 0.8) < 0.005,
    "Loaded composition retains its aspect ratio",
  );
  await first.locator(".film-frame").scrollIntoViewIfNeeded();
  await first.getByRole("button", { name: /Play chapter 1/ }).click();
  await playing(page);
  await hideDocument(page, true);
  assert.equal(
    await page.locator("video").evaluate((v) => v.paused),
    true,
    "Hidden document pauses",
  );
  await hideDocument(page, false);
  await page.waitForTimeout(200);
  assert.equal(
    await page.locator("video").evaluate((v) => v.paused),
    true,
    "Visible document does not auto-resume",
  );
  await first.getByRole("button", { name: /Play chapter 1/ }).click();
  await playing(page);
  const detached = await page.locator("video").elementHandle();
  await page.route(mediaPattern, (route) => route.abort("failed"));
  const second = page.locator(".project-card").nth(1);
  await second.locator(".film-frame").scrollIntoViewIfNeeded();
  await second.locator(".film-play").click();
  await second
    .getByText("The video couldn’t load.")
    .waitFor({ timeout: 15000 });
  assert.deepEqual(
    await detached.evaluate((v) => ({
      connected: v.isConnected,
      paused: v.paused,
      src: v.getAttribute("src"),
    })),
    { connected: false, paused: true, src: null },
    "A failed new film still cleans up the previous one",
  );
  await page.unroute(mediaPattern);
  await second.getByRole("button", { name: "Try again" }).focus();
  await page.keyboard.press("Enter");
  await playing(page);
  assert.equal(
    await page.locator("video").evaluate((v) => document.activeElement === v),
    true,
  );
  await page.locator("video").evaluate((v) => {
    v.currentTime = 23.8;
  });
  await second
    .getByRole("button", { name: "Watch again" })
    .waitFor({ timeout: 5000 });
  await second.getByRole("button", { name: "Watch again" }).focus();
  await page.keyboard.press("Enter");
  await playing(page);
  assert.equal(
    await page.locator("video").evaluate((v) => document.activeElement === v),
    true,
  );
  assert.ok(
    await page.locator("video").evaluate((v) => v.currentTime < 2),
    "Replay begins at zero",
  );

  let release;
  const gate = new Promise((resolve) => {
    release = resolve;
  });
  await page.route(mediaPattern, async (route) => {
    await gate;
    await route.continue();
  });
  const third = page.locator(".project-card").nth(2);
  await third.locator(".film-frame").scrollIntoViewIfNeeded();
  await third.locator(".film-play").click();
  await third.getByText("Loading video…").waitFor();
  await hideDocument(page, true);
  release();
  await page.waitForTimeout(750);
  assert.equal(
    await page.locator("video").evaluate((v) => v.paused),
    true,
    "Hidden pending playback stays paused",
  );
  await hideDocument(page, false);
  await page.unroute(mediaPattern);
  await third.getByRole("button", { name: /Play chapter 1/ }).click();
  await playing(page);
  results.push({
    rotation: true,
    hiddenPause: true,
    noAutomaticResume: true,
    selectionCleanup: true,
    errorRecovery: true,
    replay: true,
    delayedResponseCancellation: true,
  });
  await page.close();
  await writeFile(
    `${out}/results-${engine}.json`,
    JSON.stringify(results, null, 2) + "\n",
  );
  console.log(JSON.stringify(results, null, 2));
} finally {
  await browser.close();
}
