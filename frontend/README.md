# Portfolio frontend

The Vite, React 18, and TypeScript application for [jadenrazo.dev](https://jadenrazo.dev).

The homepage is the current presentation surface. It positions Jaden Razo for cloud, DevOps, platform, and SRE roles and links claims to public repositories, tests, incident records, runbooks, and explicit limitations. Supporting routes retain the project index, blog, status view, contact form, and administrative tools.

## Local verification

Use Node.js 20 or newer. CI currently runs Node.js 22.

```bash
npm ci
npm run type-check
npm run lint
npm run build
```

Run the development server with:

```bash
npm run dev
```

Vite serves the site at `http://localhost:3000` by default and proxies `/api` to `localhost:8080` and `/ws` to `localhost:8082`.

## Configuration

API requests use same-origin URLs unless one of these is set:

- `VITE_API_URL` at development or build time
- `window._env_.REACT_APP_API_URL` through `public/env-config.js` at runtime

The latter remains for compatibility with the current deployment packaging.

## Design and accessibility

- System font stacks avoid a third-party runtime dependency.
- Navigation uses native scrolling and honors reduced-motion preferences.
- The homepage has a visible-on-focus skip link and a 48px interaction target floor.
- Responsive QA covers 375px, 412px, 430px, tablet, and desktop layouts.
- The social preview at `public/images/og-image.png` is captured from the real homepage rather than stock artwork.

## Build boundary

`npm run build` writes static assets to `build/`. The repository CI verifies the build but does not deploy it or receive production credentials. Backend-powered routes require the corresponding Go services and environment-specific data stores.

See the repository-level [README](../README.md), [CONTRIBUTING](../CONTRIBUTING.md), and [SECURITY](../SECURITY.md) documents for broader project guidance.
