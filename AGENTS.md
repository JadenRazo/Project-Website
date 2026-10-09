# Engineering portfolio

Present Jaden's engineering work through understandable, inspectable evidence.
Keep claims tied to repositories, tests, incidents or dated measurements and
their limits. Preserve the current presentation and project films for focused
repairs; older code and routes are not permission for a redesign.

## Select the owning source

- `frontend/`: React/Vite UI; its package manifest owns frontend commands.
- `backend/cmd/` and `backend/internal/`: Go services and shared behavior; inspect
  the affected entry point and configuration before running it.
- `scripts/ci/` and `.github/workflows/ci.yml`: existing build-debt ratchet and
  review checks. Never expand the broken baseline to conceal a regression;
  remove repaired packages from it.
- `.github/workflows/deploy-content.yml`: separate `live`-branch frontend release.
  It qualifies the README's broad statement that checked-in workflows do not deploy.

## Validation and side effects

For frontend changes, use `npm run type-check`, `npm run lint` and
`npm run build` from `frontend/`, with focused rendered-page checks where affected.
For backend changes, use the relevant packages, `go build -o /dev/null ./cmd/...`
and `go vet ./cmd/...` from `backend/`; CI lists the scoped race tests and runs
`scripts/ci/build-ratchet.sh` from root. `go build ./...` has recorded debt and
must not be represented as passing. Read the workflow for the current test scope.
Prose-only edits need source/link checks, not product builds.

Keep auth, contact-message and visitor data private. Use isolated data for
runtime checks. Root start/fresh/kill scripts are process-management operations;
inspect them before use and preserve other sessions' services.

`main` review CI is distinct from deployment: qualifying pushes to `live` and
manual runs on `live` can publish to S3 and invalidate CloudFront. Preserve
conditional homepage replacement, retained rollback bytes, assets-before-HTML
ordering, old assets and runtime configuration. Deployment must stay within its
frontend authority, without changing APIs, DNS or databases. A green build does
not prove publication or live availability.

Maintain `AGENTS.md`, preserve concurrent edits and do not load legacy instruction
files. Docs explain purpose and action; PRs explain the concrete behavior change,
checks actually run and material limits. Use concise outcome-based commits
(`type: change`) and distinguish local validation from deployment.
