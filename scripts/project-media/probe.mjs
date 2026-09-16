import assert from 'node:assert/strict'
import { readFile, writeFile } from 'node:fs/promises'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { spawnSync } from 'node:child_process'
import { createHash } from 'node:crypto'

const root = resolve(dirname(fileURLToPath(import.meta.url)), '../..')
const manifest = JSON.parse(
  await readFile(
    resolve(root, 'frontend/src/data/project-film-assets.json'),
    'utf8',
  ),
)
const results = []
assert.equal(Object.keys(manifest).length, 6)
for (const [id, asset] of Object.entries(manifest)) {
  for (const [key, path] of Object.entries(asset)) {
    const input = resolve(root, 'frontend/public' + path)
    const buffer = await readFile(input)
    const hash = createHash('sha256').update(buffer).digest('hex').slice(0, 10)
    assert.ok(path.includes(`-${hash}.`), `Content hash matches: ${path}`)
    if (key.includes('Src') || key === 'src' || key === 'webmSrc') {
      const p = spawnSync(
        'ffprobe',
        ['-v', 'error', '-show_streams', '-show_format', '-of', 'json', input],
        { encoding: 'utf8' },
      )
      assert.equal(p.status, 0, p.stderr)
      const info = JSON.parse(p.stdout),
        stream = info.streams[0]
      const portrait = key.startsWith('mobile')
      assert.equal(info.streams.length, 1, 'Silent video has one video stream')
      assert.equal(stream.codec_type, 'video')
      assert.equal(stream.codec_name, path.endsWith('.mp4') ? 'h264' : 'vp9')
      assert.equal(stream.pix_fmt, 'yuv420p')
      assert.equal(stream.width, portrait ? 896 : 1600)
      assert.equal(stream.height, portrait ? 1120 : 1000)
      assert.equal(stream.avg_frame_rate, '30/1')
      assert.ok(Math.abs(Number(info.format.duration) - 24) < 0.1)
      assert.ok(buffer.length < 1000000, `${path} exceeds the 1 MB film budget`)
      if (path.endsWith('.mp4'))
        assert.ok(
          buffer.indexOf('moov') < buffer.indexOf('mdat'),
          'MP4 supports faststart',
        )
      results.push({
        id,
        key,
        path,
        bytes: buffer.length,
        codec: stream.codec_name,
        width: stream.width,
        height: stream.height,
        duration: info.format.duration,
        fps: 30,
      })
    } else if (key === 'captions') {
      const vtt = buffer.toString('utf8')
      assert.ok(vtt.startsWith('WEBVTT'))
      const cues = vtt.split('\n').filter(line => {
        const times = line.split(' --> ')
        return times.length === 2 && times.every(time => /^\d{2}:\d{2}:\d{2}\.\d{3}$/.test(time))
      })
      assert.equal(cues.length, 3)
    }
  }
}
assert.equal(results.length, 24, 'Six films × two compositions × two codecs')
const report = {
  films: results,
  totalVideoBytes: results.reduce((sum, film) => sum + film.bytes, 0),
  largestVideoBytes: Math.max(...results.map((film) => film.bytes)),
}
if (process.env.PROBE_REPORT)
  await writeFile(
    process.env.PROBE_REPORT,
    JSON.stringify(report, null, 2) + '\n',
  )
console.log(JSON.stringify(report, null, 2))
