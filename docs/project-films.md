# Homepage project films

The homepage keeps the six projects on the deployed September 16, 2026 site:
RaizHost, CloudCostMCP, TicketHacker, llm-lint, SRE Reference App, and SRE Landing
Zone. Each has a 24-second, silent film with three eight-second chapters. The
project stage appears immediately after the introduction.

## Viewing contract

- One mounted player. Selecting a project replaces and releases the old player.
- No video source, video request, autoplay, or media preload before explicit play.
- A desktop composition (1600 × 1000) and a separately laid-out phone composition
  (896 × 1120). The initial play gesture chooses the composition; it keeps that
  geometry through rotation. Selecting a new film adapts to the current viewport.
- Native controls for seeking, playback speed, captions, and fullscreen where the
  browser supports them. Keyboard project selection, chapter shortcuts, an inline
  transcript, and explicit retry/open-video recovery remain available.
- Playback pauses when the stage leaves the viewport or the document is hidden.
  Returning does not resume it. Reduced motion also disables UI transitions.
- H.264 MP4 with faststart is preferred when supported; VP9 WebM supplies a
  fallback. Only one format/composition is downloaded. Both are silent, 30 fps.
- Every distributed asset has a content hash in its URL, including posters and
  WebVTT descriptions. Long-lived immutable caching is safe for these paths.

## What the films show

| Project | Source and boundary |
| --- | --- |
| RaizHost | Public site screenshots captured September 16, 2026; source-backed editor/publish diagrams. Architecture source `raizhost-architecture@ccdfe89d95f3935cd8051f54769d43b8aeb79b01`. No customer editor session or production publish is shown. |
| CloudCostMCP | Actual CLI output at `96bab3231ee1f1e877e12a2d9d872741abd397a0`, using a sample Terraform `aws_instance.web`, `t3.small`, `us-east-1`, 730 hours/month. $24.18 is the captured estimate, including $9 synthetic egress, not a current quote or actual bill. Pricing metadata is in the raw evidence. |
| TicketHacker | Actual dashboard at `2c2e0481e4a10e68bf1df555dc1e29a6a143cf8b`, with local intercepted sample API responses: inbox, high-priority filter, resolved ticket. The film visibly labels sample data; no backend or live messaging integration is demonstrated. |
| llm-lint | Actual version 0.4.1 scan, rule, and fix-preview output from a disposable repository containing a tracked `.cursorrules`. The proposed repair is not applied. Configurable policy does not override required attribution. |
| SRE Reference App | Source `791ac0458a75c65060c170c846fbb12c48172570`, `docs/chaos-experiments.md`. Animated explanation of a historical controlled task-stop exercise. The recorded 78 seconds is one observation; time is compressed. |
| SRE Landing Zone | Source `9a5ffcbcb2d6ced9b9d2f78710c9d30ba03966d4`. Animation of the recorded May 2026 five-account lab, pilot-light standby, and tag-scoped idle stop. It does not assert current AWS resource state. |

Raw capture inputs and CLI output live in `scripts/project-media/captures/` and
`evidence/`. CLI screens are readable excerpts of those outputs. Architecture
sequences explain the cited design; they are not recordings of live cloud drills.
Update the capture date, evidence, transcript, and film together when facts change.

## Regeneration and verification

Requires Node 22+, ffmpeg with libx264/libvpx/libwebp, and Chromium. Install the
frontend and media-tool dependencies with `npm ci` in their respective directories,
then run `npx playwright install chromium` from `scripts/project-media`.

```sh
# From scripts/project-media; checked-in captures are enough to rebuild.
npm run render
# From frontend, after the asset manifest is final:
npm run build
npm run preview -- --host 127.0.0.1 --port 4190
# In another terminal, from scripts/project-media:
npm run verify
```

`render.mjs --only=raizhost` and `encode-webm.mjs --only=raizhost` support focused
updates; finish with `node fingerprint.mjs`. Streaming PNG frames directly into
ffmpeg avoids large intermediate frame folders. The renderer uses repeatable
compositions; each asset hash reflects the actual encoded bytes.

`capture-products.mjs` captures the public RaizHost site plus an isolated local
TicketHacker dashboard (`TICKET_URL`, default `http://127.0.0.1:4191`). Supply
`TICKET_REVISION` with the exact dashboard commit. It refuses non-local TicketHacker
URLs and intercepts that dashboard's API requests with fictional fixtures.

`verify.mjs` checks 375/412/430/768/1440 px, project selection, zero initial video
requests, all 36 chapter/composition combinations, accessibility, keyboard focus,
rotation, visibility pause, interrupted loading, retry, replay, and cleanup.
`ENGINE=webkit` or `ENGINE=firefox` runs another installed browser. Headless
visibility transitions are simulated explicitly. Real-device Safari acceptance
still requires a physical iPhone. `PLAYWRIGHT_MODULE`, `AXE_MODULE`, `PREVIEW_URL`,
and `EVIDENCE_DIR` allow existing QA installations and a built local preview.

## Release boundary

The production site was independently checked on September 16, 2026: S3 behind
CloudFront, with June 22 HTML. GitHub main had newer site changes that were not
live. A build of this branch includes those maintained changes as well as this
showcase. Review the complete preview before publishing; deployment is a distinct
authorized step. The broken LinkedIn shortener is replaced by the owner's direct
`https://www.linkedin.com/in/JadenRazo` in all five active/fallback source locations.

## Separate product follow-ups

During capture, CloudCost's `what-if` result for changing the fixture from t3.small
to t3.micro unexpectedly increased the compute estimate ($15.18 → $56.94). The raw
output is retained for investigation; this behavior is not presented in the film.
TicketHacker's detail route hit a Tiptap extensions error in the local dependency
installation. The film uses the actual working inbox/filter/status flow. These
are separate product issues, not silently claimed as repaired here.

## Player acceptance notes

Inline playback reserves a 76px native-control area below the composition, so a
paused phone video leaves its explanatory text readable. The player checks both
`play` and `playing` events: WebKit can resume muted media on re-entry without
emitting another `playing` event. The pause stays enforced through the re-entry
paint; native controls then work even when their events remain inside the
browser's shadow tree. Explicit play/chapter actions can resume immediately.
Advancing frames clear a stale loading label after a WebKit seek.

The reproducible CLI capture is `capture-cli.mjs`; set `CLOUDCOST_CLI` to the local
built CLI file, `CLOUDCOST_REVISION` to its verified commit, and `LLM_LINT_BIN` to the
local lint executable. It creates a disposable fixture, records both stdout and
stderr, and never applies the proposed repair. Review changed evidence before
regenerating films; raw values are deliberately not updated on a schedule.
