import { readFile, rename, writeFile } from 'node:fs/promises'
import { dirname, resolve, extname, basename } from 'node:path'
import { fileURLToPath } from 'node:url'
import { createHash } from 'node:crypto'

const root = resolve(dirname(fileURLToPath(import.meta.url)), '../..')
const manifestPath = resolve(root, 'frontend/src/data/project-film-assets.json')
const assets = JSON.parse(await readFile(manifestPath, 'utf8'))
for (const asset of Object.values(assets)) {
  for (const key of ['poster', 'mobilePoster', 'captions']) {
    const input = resolve(root, 'frontend/public' + asset[key])
    const buffer = await readFile(input)
    const hash = createHash('sha256').update(buffer).digest('hex').slice(0, 10)
    const stem = basename(input, extname(input)).replace(/-[a-f0-9]{10}$/, '')
    const filename = `${stem}-${hash}${extname(input)}`
    const target = resolve(dirname(input), filename)
    if (input !== target) await rename(input, target)
    asset[key] = `/media/projects/${filename}`
  }
}
await writeFile(manifestPath, JSON.stringify(assets, null, 2) + '\n')
