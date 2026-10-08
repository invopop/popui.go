# Applications admin prototype → access handoff

This directory is a clickable, backend-free prototype of the Access admin
**Applications** area (app list, configuration, actions, availability and
NATS credentials), built with PopUI. It exists so that the UI can be designed
here while the backend evolves separately in `invopop/access`, and then merged.

`popui build` also writes every GET page of the prototype under
`public/examples/apps/…/index.html` (`static.go` → `StaticPages`), so the
published docs site and its Netlify deploy previews can be browsed and
reviewed; HTMX navigation works there because each file is a full document
and htmx extracts the body. Saving, issuing, reordering and deleting need
the dev server and 404 on the static site.

Run it with `air` (or `go run ./cmd/popui serve`) and open
<http://localhost:3000/examples/apps>. Append `?sudo=0` to any URL to see the
non-sudo rendering.

Origin: [invopop/access#271](https://github.com/invopop/access/pull/271)
("redesign application & action admin UI"). The templates here started as a
copy of that PR's `apps.templ`, `app_actions.templ` and
`app_credentials.templ`, with the domain types swapped for local stand-ins.

## What is real and what is fake

| Concern | Prototype | access |
|---|---|---|
| Templates (`*.templ`) | **Real.** These are the deliverable. | Port verbatim into `internal/interfaces/web/components/admin/`. |
| Alpine controllers (`layoutHead` in `layout.templ`) | Real, inlined in the page head. | Already live in `assets/scripts/admin.js` (`actionId`, `appTagList`, `appAllowOrgIDList`, `appChannelList`). Keep the JS there; do not inline. |
| `navTabs` / `navTab` (`tabs.templ`) | Copied from access. | Already exists; no change. |
| `Layout` (`layout.templ`) | Look-alike shell; only "Applications" navigates. | Use the real `admin.Layout`. |
| Models in `data.go` (`Org`, `Application`, `Channel`, `Action`, `ActionResult`, `Credential`, collections) | Stand-ins with the **same field names** as `models.Application`, `gateway.Action`, `models.NATSCredential`. | Replace with the real types; see type map below. |
| Vocabularies in `data.go` (categories, scopes, schemas, verbs, renewal periods, addons, currencies) | Hard-coded copies. | Come from `models.ApplicationCategories()`, `scope.Options()` + `Key.Description()`, the `admin` package vars, `tax.AllAddonDefs()`, `currency.Definitions()`. |
| Forms in `forms.go` (`ApplicationForm`, `ActionForm`, `ActionResultForm`, converters) | Copied from `admin/apps.go` and `admin/app_actions.go`, with `cbc.Key`/`scope.List` flattened to `string`/`[]string`. | Keep the access versions. |
| Handlers in `apps.go` | Fake. GET renders fixtures; POST/PUT re-render what was posted with a flash. Only the availability toggle mutates the in-memory fixtures. | Use `adminAppsController` as in the PR. |
| Sudo | `isSudo(ctx)` reads a context flag set from `?sudo=`. | `scope.FromContext(ctx).Contains(scope.Sudo)` — same call site, same name. |
| Flash | `popui.go/flash` — identical. | Identical. |

## Listing tab (new, no backend yet)

The **Listing** tab holds the marketplace-facing content of an app, modelled
on the Stripe App Marketplace listing form. `Listing` in `listing_data.go`
is the shape the backend needs on `models.Application` (or a sibling doc):

| Field | Limit | Notes |
|---|---|---|
| `About` | 1000 chars | long description |
| `Publisher` | 80 chars | "Built by"; defaults to the org name |
| `Features[0..2]` | title 80, description 300, image ≥1600px PNG/JPG ≤5MB | ordered; three slots always render and the `listingFeatures` Alpine controller (in `layoutHead`, move to `admin.js`) shows/hides/clears them for add, remove and the empty state; posted as parallel `feature_titles[]` / `feature_descriptions[]` / `feature_image_urls[]` plus `feature_order` from the deck; slots without a title are dropped |
| `Countries` | ISO alpha-2 list | empty means worldwide; use gobl `l10n` for the options |
| `WebsiteURL`, `PrivacyURL`, `TermsURL`, `FAQURL` | URLs | `PrivacyURL` required before an app goes public |
| `SupportEmail`, `SupportURL` | | at least one channel |

The App tab's Name now carries `maxlength="35"`, matching the listing rule;
enforce it in `Validate()` alongside the 120-char Description.

Routes: `GET/PUT /apps/:id/listing`. The PUT re-renders what was posted (no
persistence) like the app form; feature images need the same multipart
handling as icon/logo.

## Documents tab (new, no backend yet)

The **Documents** tab manages the example GOBL documents shown on the app's
docs page ("Documents" tab of `spain.mdx`). It does not exist in access; the
prototype defines the shape the backend needs:

- `Document{ID, Title, Description, Schema, Channel, Content}` — `Content`
  is minimal GOBL JSON.
- Display is derived, not stored, exactly as for Workflows and Actions: a
  titled section per **channel** (app channel order, then the default
  channel) and inside it one reorderable deck per **schema**
  (`DocumentGroup`). There are no user-defined categories. Groups appear
  with their first document and disappear with their last.
- Routes under `/apps/:id/documents` (see `apps.go` → `Register`): list,
  create/edit/delete, and `groups/:channel/:schema/reorder` taking the deck's
  comma-separated `order` input and answering 204 (`hx-swap="none"`).
- These handlers **mutate the in-memory fixtures** (`documents.go`).

Backend work: a `Documents` store keyed by app, the routes above on
`adminAppsController`, and a publisher that renders the groups into the docs
page or an API the docs can pull from.

## Workflows tab (new, no backend yet)

The **Workflows** tab manages the example workflows shown on the app's docs
page ("Workflows" tab of `spain.mdx`, one `snippets/workflows/es/*.mdx` per
workflow). Same status as Documents: not in access, shaped by the prototype.

- `Workflow{ID, Name, Description, Channel, Schema, Steps, Rescue}`;
  `WorkflowStep{Name, Provider, Config}` matches the Console's workflow JSON,
  so `Workflow.Content()` is what readers import.
- Display is derived, not stored: a titled section per channel (app channel
  order, then the default channel) and inside it one deck per schema
  (`WorkflowGroup`). Order is kept per group; groups appear when the first
  workflow with that channel+schema is added and disappear when the last is
  removed.
- Routes under `/apps/:id/workflows` (see `apps.go` → `Register`), with
  `groups/:channel/:schema/reorder` taking the deck's `order` input and
  answering 204 (`hx-swap="none"`). The default channel is spelled
  `default` and the schema slash becomes a dash (`bill-invoice`).
- The editor keeps the steps as JSON in a monospace textarea; the form's
  name, description and schema win over the JSON's on save.
- These handlers mutate the in-memory fixtures (`workflows.go`).

Backend work: a `Workflows` store per app and the routes above, plus the
same docs publisher as Documents.

## Integration tab

`GET/PUT /apps/:id/integration` replaces the Credentials tab. It holds, top
to bottom: How to connect (read-only), then one form with Delivery, OAuth
and Enrollment (header Save, `form="app-integration-form"`), then the NATS
credentials table with Issue/Revoke (NATS apps only; HTTP apps get a
Rotate secret button in How to connect instead). The credential routes
moved under `/integration/credentials`. **Creating an app lands on this
tab** so the one-time OAuth client secret and HTTP signing secret are shown
once; `create` must therefore render `AppIntegration` with the posted form.
The App tab keeps Name, ID, Description, Icon, Logo, Categories, Channels,
URLs (config, launch, info), Tags and Metadata — the PR #271 shape minus
OAuth, Scopes and API URL (Scopes sit with OAuth on Integration).

## Third-party developer additions (blockers 2 and 3 of the gaps doc)

New fields and read-only content added for outside developers; none exist in
access today.

| Where | Field / content | Backend |
|---|---|---|
| Integration → Delivery | `Transport` (`nats` default / `http`), `APIURL` moved here and required for HTTP, `SigningSecret` shown once + `SigningSecretHash` | New columns on `models.Application`; HTTP delivery in the gateway: `POST` API URL, `Invopop-Signature: t=<unix>,v1=<HMAC-SHA256(secret, t+"."+body)>`, 5-min replay window, 30-s timeout, `200` result / `202` async. `POST …/integration/rotate-secret` mints a new secret and shows it once (prototype does not store it). Likewise `POST …/integration/reset-client-secret` replaces a lost OAuth client secret (new hash, shown once) — a "Reset secret" button sits where the secret would be. The OAuth explainer is an `Accordion` ("How the app gets a token for a workspace"), not a tooltip. |
| Integration → OAuth (moved from the App form, **with the Scopes multiselect**) | `oauthFlowInfo`: client-credentials example (`POST /access/v1/oauth/token` with `enrollment_id`) | Documentation only; confirm the real endpoint and parameter names. |
| Integration → Enrollment | `SettingsSchema` (JSON Schema for enrollment `data`), `AllowedOrigins[]` (browser origins for CORS; `appOriginList` Alpine controller → `admin.js`) | Store both; validate `PATCH /access/v1/enrollment` data against the schema; answer CORS preflights for the listed origins (plus the Config URL origin). |
| Integration → How to connect | NATS server `tls://connect.ngs.global`, subject `gw.<app>.task`, SDK pointer, retry semantics; HTTP apps see the webhook details instead; the one-time modal repeats server + subject | Documentation only, `natsURL` constant. |
| Action → Processing | Descriptions on "Exclude silo entry" / "Optional silo entry", retry policy note | Documentation only. |
| Action → Configuration | `ConfigSchema` (JSON Schema for the step config) | New field on `gateway.Action`; Console renders the step form from it. |
| Action → Contract | `taskRequestExample` / `taskResultExample` read-only JSON | Documentation only; align with the protobuf messages. |
| Action → Results | Status select offers NA/OK/KO/SKIP/ERR/QUEUED (`actionResultStatuses`) | Confirm numeric values against `gateway.ActionResultStatus`; the admin today only accepts 0–2. |

## Type map

| Prototype (`examples/apps`) | access |
|---|---|
| `*Org` | `*models.Org` |
| `*Application`, `*ApplicationCollection` | `*models.Application`, `*models.ApplicationCollection` |
| `Channel{Key string}` | `models.Channel{Key cbc.Key}` — templates use `ch.Key.String()` |
| `Application.Categories []string` | `[]cbc.Key` — `opt.Key.In(form.Categories...)` |
| `Application.Scopes []string` | `scope.List` — `form.Scopes.Has(scp)` |
| `Application.Visibility string` | `cbc.Key` — compare against `models.AppVisibility*` and call `.String()` in the Alpine `x-data` |
| `Application.DisabledAt *time.Time` | `*at.Timestamp` |
| `*Action`, `*ActionCollection` | `*gateway.Action`, `*gateway.ActionCollection` (`github.com/invopop/gateway/protocol`) |
| `ActionResult.Status int32` | `gateway.ActionResultStatus` |
| `*Credential`, `*CredentialCollection` | `*models.NATSCredential`, `*models.NATSCredentialCollection` |
| `Definition{Key, Name, Desc string}` | `*cbc.Definition` — `.Key.String()`, `.Name.String()`, `.Desc.String()` |
| `Scope{Key, Description}` | `scope.Key` + `scope.Key.Description()` |
| `ApplicationModelToForm`, `newApplicationForm` | `admin.ApplicationModelToForm`, `admin.ApplicationModelToForm(models.NewApplication(""))` |
| `actionFormFromAction`, `ActionForm.ToAction` | `admin.actionFormFromGateway`, `ActionForm.ToGatewayAction` |
| `bindResults(url.Values)` | not needed — access decodes `results[i].*` with `go-playground/form` |

## Template signatures

The access templates take a `*models.FeatureCollection` that the prototype
omits because nothing renders it. Restore the parameter when porting:

| Prototype | access (PR #271) |
|---|---|
| `AppsIndex(org, apps, showDeactivated)` | same |
| `AppsEdit(org, form, err)` | `AppsEdit(org, form, fc, err)` |
| `AppsEditAvailability(org, form, err)` | same |
| `AppActions(org, app, ac, err)` | same |
| `AppActionEdit(org, app, act, err)` | `AppActionEdit(org, app, act, fc, err)` |
| `AppIntegration(org, app, form, enabled, revocationEnabled, creds, issued, err)` | replaces `AppCredentials`; adds the app form for delivery / OAuth / enrollment |

## URL helpers

`baseURL(org, "apps", id)` and `orgURL(…)` keep the access call shape but
drop the org segment: here they yield `/examples/apps/<id>`, in access
`/admin/<org>/apps/<id>`. Routes in `apps.go` → `Register` map 1:1 onto
`adminAppsController.serve`.

`navAttrs(url)` is a small helper that returns the
`hx-get` / `hx-target` / `hx-push-url` triple used on every tab, row,
breadcrumb and header button. access spells the three attributes out inline;
either copy the helper into `admin.go` (recommended, it removes a lot of
noise) or expand it back.

## Deliberate differences from PR #271

These are UI decisions made in the prototype that the port should keep:

- **Fixed sidebar.** `Layout` passes `props.App{FillViewport: true}` so the
  sidebar stays pinned and only `Main` scrolls. Carry this over to
  `admin.Layout` in access (it is a one-line change to its `popui.App` call).
- **Description** is a single-line `Input` with `maxlength="120"` spanning
  both columns directly under Name and App ID, instead of a five-row
  `Textarea` below the grid. Backend: enforce the 120-character limit on
  `Application.Description` in `Validate()`.
- **Field tooltips.** All carry `Delay: tooltipDelay` (300ms) so they do not flash on mouse-over. Every labelled field on the app form and the sudo
  availability block carries a `props.Tooltip` (the "(i)" next to the label)
  explaining its purpose. Keep the copy when porting; it is the only in-app
  documentation these fields have.
- **Icon and logo are uploads.** `Icon URL` / `Logo URL` inputs became
  `popui.FileUpload` fields (`name="icon"`, `name="logo"`, `accept="image/*"`)
  with the current URL carried in hidden `icon_url` / `logo_url` inputs so an
  untouched form keeps it. **Backend work:** the create/update handlers must
  accept multipart (`hx-encoding="multipart/form-data"` on the form), store
  the file (assets bucket) and write the resulting URL back onto
  `Application.IconURL` / `LogoURL`. The prototype ignores the files.
- **Tab order** is App, Listing, Actions, Workflows, Documents, Integration, Availability — Listing, Workflows and Documents are new; Credentials is renamed **Integration** and absorbs the OAuth section from the App form plus the new Delivery and Enrollment sections; Availability moves to the end.
- **Tab counters.** Actions, Workflows, Documents and Credentials tabs show their count
  ("Actions (6)") when non-zero. The prototype derives them in `tabCounts`
  from the fixtures; in access pass a `tabCount` struct into `appTabs` from
  the controller, which already loads those collections.
- **Reorder forms carry `hx-disinherit="*"`.** Each deck sits in a
  `<form hx-post=… hx-trigger="popui-card-deck-reorder" hx-swap="none">`;
  htmx inherits those onto the cards' own `hx-get`, which then push a URL
  but swap nothing. Keep the disinherit when porting.
- **Bottom padding.** Every `popui.Main` takes `Class: mainClass` (`pb-16`,
  64px) so scrolled pages end with space under the last section. Worth
  making a popui default on `Main` rather than repeating it in access.
- **Page title is the app name** on every tabbed page (App, Listing,
  Actions, Workflows, Documents, Credentials, Availability), with no
  description under it; the tab strip sits directly beneath. "Create app"
  is the title for a new app. The per-page headings from PR #271 ("App
  configuration", "App actions", "NATS credentials"…) are gone.
- **Action form order**: an "Info" card (Name, Description, an **Icon upload** replacing the Icon URL input — same multipart handling and hidden `icon_url` as the app icon — and the
  read-only ID for existing actions) comes first, then the Action ID
  composer for new actions, then a "Pricing" card (Pops cost, Renewal
  period, Seat, Unroll), then URLs, Schemas, Processing, Results, Silo
  options, Metadata and Sudo as in PR #271.
- **Header Save buttons** use `props.Button{Form: appFormID}` instead of
  `Attributes: {"form": …}` — same HTML, typed prop.
- **Tabs are shared.** `appTabs(org, form, selected)` now drives all four
  pages (`""`, `"actions"`, `"availability"`, `"credentials"`); the PR had
  separate `appActionTabs` / `appCredentialsTabs` copies. The actions and
  credentials pages call `appTabs(org, ApplicationModelToForm(app), …)`.
- **"Issue credentials"** moved from the article body into the page header
  as the primary action, matching "Create app" / "New action".
- **Credential result modal** uses `Textarea{Monospace: true}`, a
  `ButtonGroup` with a Download button (`props.Button{Href, Download}`)
  instead of a raw `<a download>`, and the public key is a copy-to-clipboard
  button (`props.Button{Copy: key, CopyPrefixLength: 10, CopySuffixLength: 6}`).
- **Revoke** is a small (`sm`) danger button and drops the inline htmx
  spinner markup; `hx-disabled-elt` already disables it while in flight.
- **Action results**: the remove control is an icon-size transparent button
  with an aria-label; "Add result" carries the Add icon.
- **Action sudo block** gets a `Separator` + `TitleGroup("Sudo")` like every
  other section instead of a `Fieldset{Legend: "Sudo"}`.
- **Actions, Workflows and Documents share one layout**: channel sections, schema decks, `schemaLabel(schema, kind)` headings ("Invoice actions", "Invoice workflows", "Invoice documents").
- **Channels editor** is a bare `CardDeck` (no card `Fieldset` around it):
  the `CardDeckHead` carries a single "Key / Label" column label and the
  small "+ Add row" button on the right, each row is a `Card` with the two
  inputs and an icon-small transparent remove button. The deck is always
  rendered, so an app with no channels shows just the head.
- **Actions are card decks, not a table.** Like Workflows: a titled section
  per channel (app channel order, then "Default channel") and inside it one
  reorderable `CardDeck` per schema set (`ActionCollection.Groups`), headed by
  the schema label with the schemas in mono. Cards show name over mono ID,
  a grey "Disabled" tag, one billing tag ("Seat · monthly" or "Releases
  seat"), the pops cost as a `popsTag` (the Console's PopsTag: `icons.Pops` tile + mono count, dash when free), and a kebab with Edit plus, for sudo, Enable /
  Disable (`POST …/actions/:aid/enable|disable`, with an `hx-confirm` on
  disable). Reorder posts `order` to `…/actions/groups/:channel/:schema/
  reorder` (204, `hx-swap="none"`). **Backend:** `gateway.Action` has no
  order field; either add one or keep a per-app ordered ID list, and the
  gateway's `ActionCollection` must honour it.
- Column headers in the credentials table are sentence case
  ("Public key", "Issued by").

## Porting checklist for the access PR

1. Bump popui.go in `access/go.mod` to a tag that includes this prototype's
   PopUI features (`props.Button.Form`, `props.Textarea.Monospace`,
   `props.Article.Title/Description`, `Select{Multiple, Searchable}` with
   `SelectOption.Description`). Regenerate `styles/admin.css`
   (`go generate ./internal/interfaces/web/assets`) so the new utility classes
   in these templates are compiled.
2. Copy `apps.templ`, `app_actions.templ`, `app_credentials.templ` over the
   PR #271 versions; re-apply the type map above (`.String()`, `.In`, `.Has`),
   re-add the `fc *models.FeatureCollection` parameter, and swap
   `Layout`/`appsLayout` for the real ones.
3. Add `navAttrs` to `admin/admin.go`, or expand it inline.
4. Keep the Alpine controllers in `admin.js`; delete nothing there. Do not
   copy `layoutHead` from the prototype.
5. Run `templ generate` with the version pinned in access's `go.mod`, then
   `go test ./internal/interfaces/web/components/admin/` — the PR's
   `apps_test.go` / `app_actions_test.go` assertions (two index columns,
   header Save `form=` wiring, multiselects, sudo-only sections) all still
   hold for these templates; `examples/apps/apps_test.go` here checks the
   same landmarks through the prototype router.
6. Diff the rendered HTML of each page against the prototype at the same
   route to catch anything the type swap changed.

## Keeping the two tracks in sync

- Front-end changes happen here first; each change to a `.templ` in this
  directory should keep the access field names so the port stays mechanical.
- Backend changes that add or rename a form field should be mirrored in
  `forms.go` + `data.go` (one line each) so the prototype keeps rendering the
  real shape.
- If a new PopUI component is needed, add it to the library proper (repo
  root + docs), not to this directory.
