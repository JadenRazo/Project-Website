# Jaden Razo portfolio

The source for [jadenrazo.dev](https://jadenrazo.dev): an evidence-led portfolio for AWS cloud, DevOps, platform, and SRE roles, backed by a React frontend and a collection of Go services.

The homepage points engineering claims to inspectable repositories, tests, incident records, failure exercises, runbooks, or explicit limitations. It no longer uses lines-of-code totals, satisfaction counters, or a forced intro animation as a substitute for evidence.

## Repository map

| Path | Purpose |
| --- | --- |
| `frontend/` | Vite, React 18, TypeScript, Tailwind CSS, and styled-components application |
| `backend/cmd/` | Go entry points for the API, admin tools, messaging, URL shortening, and workers |
| `backend/internal/` | Shared auth, projects, status, messaging, visitor, and infrastructure packages |
| `scripts/ci/` | Build-debt ratchet used by CI |
| `docs/` | Architecture and supporting documentation |
| `deployments/`, `backend/deployments/` | Example process, proxy, and container configuration |

## Verification

GitHub Actions is review-only and receives no deployment credentials. It currently verifies:

- every Go binary under `backend/cmd/...` builds and passes `go vet`
- the existing auth and performance test packages pass with the race detector
- packages outside the known broken baseline cannot regress, and newly repaired packages must be removed from that baseline
- the frontend installs from its lockfile, typechecks, lints, and produces a Vite build
- CodeQL scans JavaScript/TypeScript and the deployable Go commands with the extended security query suite

Run the primary checks locally:

```bash
cd frontend
npm ci
npm run type-check
npm run lint
npm run build

cd ../backend
go mod verify
go build -o /dev/null ./cmd/...
go test -race ./internal/common/auth/... ./internal/performance/...
```

## Known limitations

- `go build ./...` does not yet pass. Historical packages under `backend/internal/messaging/` contain duplicate declarations and references to a removed type. The deployable `cmd` binaries build; CI records the remaining package debt as a ratchet rather than hiding it behind a permanently red job.
- The repository contains several generations of frontend and backend code. The evidence-led homepage is the current presentation surface; older standalone routes remain for compatibility.
- Local end-to-end use needs environment-specific PostgreSQL, Redis, proxy, and secret configuration. No production credentials are committed.
- The checked-in CI does not deploy the public site. A green workflow proves the review checks above, not production availability.

## Development

```bash
cd frontend
npm ci
npm run dev
```

API requests default to the same origin. Set `VITE_API_URL` for a build-time endpoint or provide `window._env_.REACT_APP_API_URL` through `public/env-config.js` at runtime. Backend commands have different environment needs; inspect the relevant entry point and example configuration before running one.

## Security and privacy

The application handles contact messages, authentication data, and visitor telemetry. Do not commit runtime secrets, database contents, or user data. See [SECURITY.md](SECURITY.md) for private reporting guidance.

Workflow actions are pinned to immutable commits. Dependabot groups maintenance updates monthly and ignores automatic major-version migrations so upgrades remain reviewable.

## License

MIT for original project code. Third-party dependencies and bundled assets retain their own licenses.
