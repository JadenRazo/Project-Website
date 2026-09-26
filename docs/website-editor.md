# Website editor integration

`raizhost/content-map.json` exposes public website copy and controls;
`content.json` stores values, and `website.json` registers native routes and
components. `site.json` describes the existing live-branch publishing workflow.
The shared kit in `raizhost/editor` is exported from raizhost-app, with hashes.

The home, About, Contact, Projects, Portfolio, Blog and Status routes support
page metadata/navigation/visibility plus portable text/photo/button sections.
The homepage's native sections can be reordered or removed. Interior route
bodies remain native components with editable surrounding sections. Existing
route addresses remain stable, and application/API namespaces are reserved.
New pages render through the React router. Design settings feed both CSS tokens
and the existing styled-component theme, preserving original defaults when blank.

Public JSX copy uses `SiteText` or existing content hooks; evidence/navigation
lists and the portrait use mapped controls. React owns rendering, stable section
keys and interaction state. The editor dispatches `raizhost:document-update`
for complete drafts and `raizhost:content-update` for individual changes. The
existing preview MemoryRouter uses the requested route, including pages that
have not been published. The app's sandbox blocks production form submissions.

Backend data is still edited through its existing authorized application
interfaces: blog article records, certifications, project records, status/API
responses, messaging and account data are not copied into repository JSON.
Native interactive modules, film assets and service endpoints remain intact.

Verify `npm run test:editor-compat --prefix frontend` and
`npm run build --prefix frontend`, then run the app's authenticated composition
journey and inspect the real native runtime at desktop/375/412/430px. Test a
published portable section after reopening, deleting/hiding native sections,
design reset, added pages, navigation and film/form interactions.

The existing `deploy-content.yml` workflow owns frontend publication on `live`.
The app publishes JSON/uploads; it does not publish backend code. Deploy a
compatible app/site before connecting this source. Old app versions reject
these new field types; rollbacks require coordinated preservation of source
configuration and saved owner content.
