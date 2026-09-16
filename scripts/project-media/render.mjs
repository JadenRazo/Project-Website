import { readFile, writeFile, mkdir } from 'node:fs/promises'
import { resolve, dirname } from 'node:path'
import { fileURLToPath, pathToFileURL } from 'node:url'
import { spawn, spawnSync } from 'node:child_process'
import { once } from 'node:events'
import { createHash } from 'node:crypto'
import { build } from '../../frontend/node_modules/esbuild/lib/main.js'

const here = dirname(fileURLToPath(import.meta.url)),
  root = resolve(here, '../..')
const out = resolve(root, 'frontend/public/media/projects')
const { chromium } = await import(process.env.PLAYWRIGHT_MODULE || 'playwright')
const compiled = await build({
  entryPoints: [resolve(root, 'frontend/src/data/projectFilms.ts')],
  bundle: true,
  write: false,
  format: 'esm',
  platform: 'node',
})
const { projectFilms } = await import(
  `data:text/javascript;base64,${Buffer.from(compiled.outputFiles[0].text).toString('base64')}`
)
const manifest = JSON.parse(
  await readFile(
    resolve(root, 'frontend/src/data/project-film-assets.json'),
    'utf8',
  ),
)
const only = process.argv.find((a) => a.startsWith('--only='))?.slice(7)
const postersOnly = process.argv.includes('--posters-only')
await mkdir(out, { recursive: true })
const images = {}
for (const name of [
  'raizhost-wide',
  'raizhost-portrait',
  'ticket-inbox',
  'ticket-context',
  'ticket-resolved',
])
  images[name] =
    `data:image/png;base64,${(await readFile(resolve(here, `captures/${name}.png`))).toString('base64')}`
const evidence = {
  cost: JSON.parse(
    await readFile(resolve(here, 'evidence/cloudcost-output.json'), 'utf8'),
  ),
}
const browser = await chromium.launch({
  headless: true,
  args: ['--no-sandbox', '--disable-dev-shm-usage'],
})
try {
  const page = await browser.newPage()
  await page.addScriptTag({ path: resolve(here, 'film-canvas.js') })
  await page.evaluate(
    async ({ images, evidence }) => {
      window.paint = await window.createFilmRenderer(images, evidence)
    },
    { images, evidence },
  )
  for (const film of projectFilms.filter((f) => !only || f.id === only)) {
    const asset = {}
    for (const portrait of [false, true]) {
      const aspect = portrait ? 'portrait' : 'wide',
        base = `${film.id}-${aspect}`
      const poster = Buffer.from(
        await page.evaluate(
          ({ film, portrait }) => window.paint(film, 6, portrait, true),
          { film, portrait },
        ),
        'base64',
      )
      const p = spawnSync(
        'ffmpeg',
        [
          '-hide_banner',
          '-loglevel',
          'error',
          '-f',
          'image2pipe',
          '-i',
          'pipe:0',
          '-frames:v',
          '1',
          '-c:v',
          'libwebp',
          '-quality',
          '85',
          '-y',
          `${out}/${base}.webp`,
        ],
        { input: poster },
      )
      if (p.status !== 0) throw new Error(p.stderr.toString())
      asset[portrait ? 'mobilePoster' : 'poster'] =
        `/media/projects/${base}.webp`
      if (!postersOnly) {
        const target = `${out}/${base}.mp4`
        const encoder = spawn(
          'ffmpeg',
          [
            '-hide_banner',
            '-loglevel',
            'error',
            '-f',
            'image2pipe',
            '-framerate',
            '30',
            '-vcodec',
            'png',
            '-i',
            'pipe:0',
            '-an',
            '-c:v',
            'libx264',
            '-preset',
            'fast',
            '-crf',
            '23',
            '-threads',
            '2',
            '-pix_fmt',
            'yuv420p',
            '-movflags',
            '+faststart',
            '-y',
            target,
          ],
          { stdio: ['pipe', 'ignore', 'pipe'] },
        )
        const finished = once(encoder, 'close')
        let err = ''
        encoder.stderr.on('data', (d) => (err += d))
        encoder.stdin.on('error', () => {})
        let lastKey = '',
          lastBuffer
        for (let frame = 0; frame < 720; frame++) {
          const t = frame / 30,
            local = t % 8,
            scene = Math.floor(t / 8)
          const fixed =
            film.id === 'cloudcostmcp' ||
            film.id === 'llm-lint' ||
            film.id === 'tickethacker' ||
            (film.id === 'raizhost' && scene === 0)
          // Terminal holds are intentional reading time. Packet flows keep moving.
          const sampled = fixed && local > 2 ? scene * 8 + 6 : t
          const key = String(sampled)
          if (key !== lastKey) {
            lastBuffer = Buffer.from(
              await page.evaluate(
                ({ film, t, portrait }) => window.paint(film, t, portrait),
                { film, t: sampled, portrait },
              ),
              'base64',
            )
            lastKey = key
          }
          try {
            if (!encoder.stdin.write(lastBuffer))
              await once(encoder.stdin, 'drain')
          } catch (error) {
            throw new Error(err || error.message)
          }
        }
        encoder.stdin.end()
        const [code] = await finished
        if (code !== 0) throw new Error(err)
        const hash = createHash('sha256')
          .update(await readFile(target))
          .digest('hex')
          .slice(0, 10)
        const hashed = `${base}-${hash}.mp4`
        const { rename } = await import('node:fs/promises')
        await rename(target, `${out}/${hashed}`)
        asset[portrait ? 'mobileSrc' : 'src'] = `/media/projects/${hashed}`
        console.log(
          film.id,
          aspect,
          (await readFile(`${out}/${hashed}`)).length,
          'bytes',
        )
      }
    }
    const stamp = (n) => `00:00:${String(n).padStart(2, '0')}.000`
    const vtt =
      'WEBVTT\n\n' +
      film.chapters
        .map(
          (ch, i) =>
            `${stamp(ch.at)} --> ${stamp((i + 1) * 8)}\n${ch.title}. ${ch.description}\n`,
        )
        .join('\n')
    await writeFile(`${out}/${film.id}.vtt`, vtt)
    manifest[film.id] = {
      ...manifest[film.id],
      ...asset,
      captions: `/media/projects/${film.id}.vtt`,
    }
    await writeFile(
      resolve(root, 'frontend/src/data/project-film-assets.json'),
      JSON.stringify(manifest, null, 2) + '\n',
    )
  }
} finally {
  await browser.close()
}
