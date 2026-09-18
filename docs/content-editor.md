# Portfolio content editing

This integration starts from the released project-card film revision
`5ad54207574e8099c583fce3f9eed937600d2c1f`. It preserves the current cloud/DevOps
design, all six films, media hashes, transcripts, playback controls and links.
It does not restore the older Build/Create/Innovate opening.

`raizhost/content-map.json` declares owner controls and `raizhost/content.json`
is the source of editable copy. The React adapter accepts only declared paths,
types, bounds and select options; updates change state without replacing
framework-owned structures. Repeated film names and categories share one value.

Owners can edit the homepage headline, background, operating principles, contact
copy and form labels, walkthrough copy, and About/Contact/Projects introductions.
They can show or hide project categories and technology tags.
Film activation, keyboard controls, playback, captions and evidence remain
intact. Media replacement and recorded statements require a coordinated update;
the editor does not expose arbitrary CSS/JavaScript or rewrite API-backed records.

The working preview uses memory routing and storage fallbacks inside the app's
opaque iframe. The public site retains its normal router and persistent storage.
Try interactions exercises the films and local form controls; submissions are
disabled. Restart preview reapplies the saved draft.

## Proposed activation

After production approval, create `live` at the exact reviewed release revision.
Keep only approved production code there so an owner wording change cannot also
release pending work from `main`. Provision the narrowly scoped OIDC role
`raizcloud-site-portfolio-editor-deploy` from the operations proposal, restricted
to this repository's `live` branch and the existing portfolio bucket/distribution.
Then activate `deploy-content.yml`, verify its exact first deployment, and attach
`JadenRazo/Project-Website` / `live` to the portfolio tenant in the editor.

Content directory: `raizhost`. Uploads: `frontend/public/uploads`, URL prefix
`/uploads`. Drafts and working previews stay in the app; there is initially no
separate shared preview branch. Backend, APIs, DNS and messaging are unchanged.

Before release, compare the latest production revision and preserve any approved
intervening changes. The earlier `portfolio-live` experiment targets a superseded
design and must not be deployed.

Checks: `npm run test:editor-compat --prefix frontend`, `npm run build --prefix
frontend`, the app's cross-repository preview checks, and the existing film QA.
