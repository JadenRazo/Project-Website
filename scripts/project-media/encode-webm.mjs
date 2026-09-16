import { readFile, rename, writeFile } from 'node:fs/promises'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { spawn } from 'node:child_process'
import { once } from 'node:events'
import { createHash } from 'node:crypto'

const root = resolve(dirname(fileURLToPath(import.meta.url)), '../..')
const manifestPath = resolve(root, 'frontend/src/data/project-film-assets.json')
const assets = JSON.parse(await readFile(manifestPath, 'utf8'))
const only = process.argv.find((arg) => arg.startsWith('--only='))?.slice(7)
for (const [id, asset] of Object.entries(assets)) {
  if (only && id !== only) continue
  for (const portrait of [false, true]) {
    const input = resolve(
      root,
      'frontend/public' + asset[portrait ? 'mobileSrc' : 'src'],
    )
    const base = `${id}-${portrait ? 'portrait' : 'wide'}`
    const target = resolve(root, `frontend/public/media/projects/${base}.webm`)
    const encoder = spawn(
      'ffmpeg',
      [
        '-hide_banner',
        '-loglevel',
        'error',
        '-i',
        input,
        '-an',
        '-c:v',
        'libvpx-vp9',
        '-b:v',
        '0',
        '-crf',
        '32',
        '-row-mt',
        '1',
        '-threads',
        '2',
        '-deadline',
        'good',
        '-cpu-used',
        '4',
        '-y',
        target,
      ],
      { stdio: ['ignore', 'ignore', 'pipe'] },
    )
    let error = ''
    encoder.stderr.on('data', (data) => {
      error += data
    })
    const [code] = await once(encoder, 'close')
    if (code !== 0) throw new Error(error)
    const buffer = await readFile(target)
    const hash = createHash('sha256').update(buffer).digest('hex').slice(0, 10)
    const hashed = `${base}-${hash}.webm`
    await rename(target, resolve(dirname(target), hashed))
    asset[portrait ? 'mobileWebmSrc' : 'webmSrc'] = `/media/projects/${hashed}`
    console.log(base, 'VP9', buffer.length, 'bytes')
  }
  await writeFile(manifestPath, JSON.stringify(assets, null, 2) + '\n')
}
