import assert from 'node:assert/strict'
import { mkdir, writeFile } from 'node:fs/promises'

const engines = await import(process.env.PLAYWRIGHT_MODULE || 'playwright')
const { default: AxeBuilder } = await import(
  process.env.AXE_MODULE || '@axe-core/playwright'
)
const engine = process.env.ENGINE || 'chromium'
const base = process.env.PREVIEW_URL || 'http://127.0.0.1:4190'
const out = process.env.EVIDENCE_DIR || '/tmp/jadenrazo-video-evidence/qa'
const mediaPattern = /\.(mp4|webm)(?:\?|$)/
await mkdir(out, { recursive: true })
const results = []
const browser = await engines[engine].launch({
  headless: true,
  args:
    engine === 'chromium' ? ['--no-sandbox', '--disable-dev-shm-usage'] : [],
})
const isolate = (page) =>
  page.route('**/api/**', (route) =>
    route.fulfill({ status: 200, contentType: 'application/json', body: '{}' }),
  )
const playing = (page) =>
  page.waitForFunction(
    () => {
      const v = document.querySelector('video')
      return v && !v.paused && v.currentTime > 0.2
    },
    {},
    { timeout: 15000 },
  )
const hideDocument = (page, hidden) =>
  page.evaluate((hidden) => {
    // Deterministically exercise the visibility lifecycle in a headless browser.
    Object.defineProperty(document, 'hidden', {
      configurable: true,
      value: hidden,
    })
    document.dispatchEvent(new Event('visibilitychange'))
  }, hidden)

try {
  for (const width of [375, 412, 430, 768, 1440]) {
    const context = await browser.newContext({
      viewport: { width, height: 1000 },
      reducedMotion: 'reduce',
    })
    const page = await context.newPage(),
      media = [],
      errors = []
    await isolate(page)
    page.on('request', (r) => {
      if (mediaPattern.test(r.url())) media.push(r.url())
    })
    page.on('pageerror', (e) => errors.push(e.message))
    await page.goto(base, { waitUntil: 'networkidle' })
    await page.locator('#projects').scrollIntoViewIfNeeded()
    await page.waitForTimeout(200)
    assert.equal(media.length, 0, 'No video request before intent')
    assert.equal(await page.getByRole('tab').count(), 6)
    assert.equal(
      await page
        .locator('#contact')
        .getByRole('link', { name: 'LinkedIn', exact: true })
        .getAttribute('href'),
      'https://www.linkedin.com/in/JadenRazo',
    )
    assert.equal(
      await page.evaluate(
        () => document.documentElement.scrollWidth <= innerWidth,
      ),
      true,
      `No overflow at ${width}`,
    )
    for (let i = 0; i < 6; i++) await page.getByRole('tab').nth(i).click()
    assert.equal(media.length, 0, 'Switching projects must not fetch video')
    await page.getByRole('tab').first().click()
    await page.locator('.film-frame').scrollIntoViewIfNeeded()
    await page.screenshot({ path: `${out}/${engine}-${width}-showcase.png` })
    await page
      .locator('#projects')
      .screenshot({ path: `${out}/${engine}-${width}-section.png` })
    const axe = await new AxeBuilder({ page })
      .include('#projects')
      .withTags(['wcag2a', 'wcag2aa', 'wcag21aa'])
      .analyze()
    assert.deepEqual(
      axe.violations.map((v) => ({
        id: v.id,
        impact: v.impact,
        nodes: v.nodes.map((n) => n.target),
      })),
      [],
      `Axe at ${width}`,
    )
    await page
      .getByRole('button', { name: /Play RaizHost walkthrough/ })
      .focus()
    await page.keyboard.press('Enter')
    await playing(page)
    const video = page.locator('video')
    assert.equal(await video.evaluate((v) => Math.round(v.duration)), 24)
    assert.equal(
      await video.evaluate((v) => document.activeElement === v),
      true,
      'Focus follows disappearing play button',
    )
    const source = await video.getAttribute('src')
    assert.ok(
      source.includes(width < 640 ? 'portrait' : 'wide'),
      'Viewport-specific composition',
    )
    assert.equal(
      new Set(media).size,
      1,
      'Only one composition/format downloaded',
    )
    await page.getByRole('button', { name: /Play chapter 3/ }).click()
    await page.waitForFunction(
      () => document.querySelector('video')?.currentTime >= 16,
    )
    await page.evaluate(() => window.scrollTo({ top: 0, behavior: 'instant' }))
    await page.waitForFunction(() => document.querySelector('video')?.paused)
    await page.locator('.film-frame').scrollIntoViewIfNeeded()
    await page.waitForTimeout(150)
    assert.equal(
      await video.evaluate((v) => v.paused),
      true,
      'Returning to the stage does not resume motion',
    )
    // One native-control pointer smoke check per engine. Native button order
    // changes with width; our keyboard selectors/chapters cover every layout.
    if (width === 375) {
      const box = await video.boundingBox()
      await page.mouse.move(box.x + 30, box.y + box.height - 22)
      await page.mouse.click(box.x + 30, box.y + box.height - (engine === 'webkit' ? 22 : 48))
      await playing(page)
      await page.waitForFunction(() => !document.querySelector('.film-loading'))
    }
    await page.getByRole('tab').first().focus()
    await page.keyboard.press('End')
    assert.equal(
      await page.getByRole('tab').last().getAttribute('aria-selected'),
      'true',
    )
    await page.keyboard.press('Home')
    assert.equal(
      await page.getByRole('tab').first().getAttribute('aria-selected'),
      'true',
    )
    await page.keyboard.press('ArrowRight')
    assert.equal(
      await page.getByRole('tab').nth(1).getAttribute('aria-selected'),
      'true',
    )
    await page.locator('.film-transcript summary').click()
    assert.equal(
      await page.locator('.film-transcript').getAttribute('open'),
      '',
    )
    assert.deepEqual(errors, [], 'No page exceptions')
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
    })
    console.log(engine, width, 'passed')
    await context.close()
  }

  // Decode every chapter of every composition, rather than only the first film.
  for (const width of [412, 1440]) {
    const page = await browser.newPage({ viewport: { width, height: 1000 } })
    await isolate(page)
    await page.goto(base, { waitUntil: 'networkidle' })
    for (let i = 0; i < 6; i++) {
      await page.getByRole('tab').nth(i).click()
      await page.locator('.film-frame').scrollIntoViewIfNeeded()
      await page.locator('.film-play').click()
      await playing(page)
      for (const t of [3, 11, 19]) {
        await page.locator('video').evaluate(async (v, t) => {
          v.pause()
          await new Promise((resolve) => {
            v.addEventListener('seeked', resolve, { once: true })
            v.currentTime = t
          })
        }, t)
        await page.locator('.film-frame').screenshot({
          path: `${out}/${engine}-${width}-film-${i + 1}-${t}.png`,
        })
      }
    }
    results.push({ width, allSixFilmsDecode: true, allChaptersSeek: true })
    await page.close()
  }

  const page = await browser.newPage({ viewport: { width: 412, height: 1000 } })
  await isolate(page)
  await page.goto(base, { waitUntil: 'networkidle' })
  await page.locator('.film-frame').scrollIntoViewIfNeeded()
  await page.locator('.film-play').click()
  await playing(page)
  // A rotation must not put a portrait video inside a landscape frame.
  await page.setViewportSize({ width: 800, height: 1000 })
  const ratio = await page
    .locator('.film-frame')
    .evaluate(
      (e) =>
        e.clientWidth /
        (e.clientHeight - parseFloat(getComputedStyle(e).paddingBottom)),
    )
  assert.ok(
    Math.abs(ratio - 0.8) < 0.005,
    'Loaded composition retains its aspect ratio',
  )
  await page.locator('.film-frame').scrollIntoViewIfNeeded()
  await page.getByRole('button', { name: /Play chapter 1/ }).click()
  await playing(page)
  await hideDocument(page, true)
  assert.equal(
    await page.locator('video').evaluate((v) => v.paused),
    true,
    'Hidden document pauses',
  )
  await hideDocument(page, false)
  await page.waitForTimeout(200)
  assert.equal(
    await page.locator('video').evaluate((v) => v.paused),
    true,
    'Visible document does not auto-resume',
  )
  await page.getByRole('button', { name: /Play chapter 1/ }).click()
  await playing(page)
  const detached = await page.locator('video').elementHandle()
  await page.getByRole('tab').nth(1).click()
  assert.deepEqual(
    await detached.evaluate((v) => ({
      connected: v.isConnected,
      paused: v.paused,
      src: v.getAttribute('src'),
    })),
    { connected: false, paused: true, src: null },
    'Switching cleans up the previous player',
  )

  // Media failure has a useful recovery path, including keyboard focus.
  await page.route(mediaPattern, (route) => route.abort('failed'))
  await page.locator('.film-frame').scrollIntoViewIfNeeded()
  await page.locator('.film-play').click()
  await page.getByText('The video couldn’t load.').waitFor({ timeout: 15000 })
  await page.unroute(mediaPattern)
  await page.getByRole('button', { name: 'Try again' }).focus()
  await page.keyboard.press('Enter')
  await playing(page)
  assert.equal(
    await page.locator('video').evaluate((v) => document.activeElement === v),
    true,
  )
  await page.locator('video').evaluate((v) => {
    v.currentTime = 23.8
  })
  await page
    .getByRole('button', { name: 'Watch again' })
    .waitFor({ timeout: 5000 })
  await page.getByRole('button', { name: 'Watch again' }).focus()
  await page.keyboard.press('Enter')
  await playing(page)
  assert.equal(
    await page.locator('video').evaluate((v) => document.activeElement === v),
    true,
  )
  assert.ok(
    await page.locator('video').evaluate((v) => v.currentTime < 2),
    'Replay begins at zero',
  )

  // A delayed response must not restart playback after the page is hidden.
  await page.getByRole('tab').nth(2).click()
  let release
  const gate = new Promise((resolve) => {
    release = resolve
  })
  await page.route(mediaPattern, async (route) => {
    await gate
    await route.continue()
  })
  await page.locator('.film-frame').scrollIntoViewIfNeeded()
  await page.locator('.film-play').click()
  await page.getByText('Loading video…').waitFor()
  await hideDocument(page, true)
  release()
  await page.waitForTimeout(750)
  assert.equal(
    await page.locator('video').evaluate((v) => v.paused),
    true,
    'Hidden pending playback stays paused',
  )
  await hideDocument(page, false)
  await page.unroute(mediaPattern)
  await page.getByRole('button', { name: /Play chapter 1/ }).click()
  await playing(page)
  results.push({
    rotation: true,
    hiddenPause: true,
    noAutomaticResume: true,
    selectionCleanup: true,
    errorRecovery: true,
    replay: true,
    delayedResponseCancellation: true,
  })
  await page.close()
  await writeFile(
    `${out}/results-${engine}.json`,
    JSON.stringify(results, null, 2) + '\n',
  )
  console.log(JSON.stringify(results, null, 2))
} finally {
  await browser.close()
}
