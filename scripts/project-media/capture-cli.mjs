import { mkdir, mkdtemp, writeFile } from 'node:fs/promises'
import { spawnSync } from 'node:child_process'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { tmpdir } from 'node:os'

// Point at verified local builds. This script never installs or executes a
// project from an unreviewed remote URL and never applies llm-lint repairs.
const { CLOUDCOST_CLI, CLOUDCOST_REVISION, LLM_LINT_BIN } = process.env
if (!CLOUDCOST_CLI || !CLOUDCOST_REVISION || !LLM_LINT_BIN) {
  throw new Error(
    'Set CLOUDCOST_CLI, CLOUDCOST_REVISION, and LLM_LINT_BIN to verified local builds',
  )
}
const here = dirname(fileURLToPath(import.meta.url))
const scratch = await mkdtemp(resolve(tmpdir(), 'portfolio-cli-capture-'))
const lintDir = resolve(scratch, 'lint-fixture')
await mkdir(lintDir)
const run = (bin, args, cwd = scratch) => {
  const result = spawnSync(bin, args, {
    cwd,
    encoding: 'utf8',
    timeout: 120000,
    env: {
      ...process.env,
      NO_COLOR: '1',
      CLOUDCOST_CACHE_PATH: resolve(scratch, 'pricing-cache.db'),
    },
  })
  if (result.error) throw result.error
  return {
    command: [bin, ...args].join(' '),
    status: result.status,
    stdout: result.stdout,
    stderr: result.stderr,
  }
}
await writeFile(
  resolve(scratch, 'main.tf'),
  'resource "aws_instance" "web" {\n  ami = "ami-example"\n  instance_type = "t3.small"\n}\n',
)
const cost = {
  analyze: run('node', [CLOUDCOST_CLI, 'analyze', 'main.tf', '--json']),
  estimate: run('node', [
    CLOUDCOST_CLI,
    'estimate',
    'main.tf',
    '--provider',
    'aws',
    '--region',
    'us-east-1',
    '--json',
  ]),
}
for (const result of Object.values(cost))
  if (result.status !== 0) throw new Error(result.stderr)
await writeFile(
  resolve(lintDir, '.cursorrules'),
  'Sample local editor configuration.\n',
)
await writeFile(
  resolve(lintDir, 'README.md'),
  '# Example repository\nDisposable capture fixture.\n',
)
for (const args of [
  ['init', '-q'],
  ['add', '.cursorrules', 'README.md'],
]) {
  const result = run('git', args, lintDir)
  if (result.status !== 0) throw new Error(result.stderr)
}
const lint = {
  scan: run(LLM_LINT_BIN, ['scan', lintDir], lintDir),
  rule: run(LLM_LINT_BIN, ['rules', 'show', 'LLM006'], lintDir),
  preview: run(
    LLM_LINT_BIN,
    ['scan', lintDir, '--fix-preview', '--fix-git-history', 'none'],
    lintDir,
  ),
}
if (
  lint.scan.status !== 1 ||
  !lint.scan.stdout.includes('LLM006') ||
  lint.rule.status !== 0 ||
  lint.preview.status !== 1 ||
  !`${lint.preview.stdout}\n${lint.preview.stderr}`.includes('would fix')
) {
  throw new Error(
    'Unexpected lint output; inspect the fixture before replacing capture evidence',
  )
}
await mkdir(resolve(here, 'evidence'), { recursive: true })
await writeFile(
  resolve(here, 'evidence/cloudcost-output.json'),
  JSON.stringify(cost, null, 2) + '\n',
)
await writeFile(
  resolve(here, 'evidence/lint-output.json'),
  JSON.stringify(lint, null, 2) + '\n',
)
await writeFile(
  resolve(here, 'evidence/cli-capture.json'),
  JSON.stringify(
    {
      capturedAt: new Date().toISOString(),
      cloudcostRevision: CLOUDCOST_REVISION,
      lintVersion: run(LLM_LINT_BIN, ['--version']).stdout.trim(),
      scratch,
    },
    null,
    2,
  ) + '\n',
)
console.log(
  'Captured real CLI evidence. Review every changed value and caption before rendering.',
  scratch,
)
