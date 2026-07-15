# Migrations — removed and renamed components

This file maps every removed or renamed PopUI component to its replacement.
It is written for both humans and AI coding agents: **if you upgraded
`github.com/invopop/popui.go` and now get a compile error like
`undefined: popui.SomeComponent` or `undefined: props.SomeComponent`, search
this file for that identifier** — every removed name appears here verbatim,
with the replacement and a before/after snippet.

Entries are grouped by the release that removed them, newest first.

## July 2026 consolidation (PR #53)

| Removed | Replacement |
|---|---|
| `popui.CardFile`, `props.CardFile` | `popui.FileDownload`, `props.FileDownload` |
| `popui.CardFileInfo`, `props.CardFileInfo` | `popui.FileDownloadInfo`, `props.FileDownloadInfo` |
| `popui.CardFieldset`, `props.CardFieldset` | `popui.FieldsetCard`, `props.FieldsetCard` |
| `popui.FlashMessage`, `props.FlashMessage` | `popui.Toast`, `props.Toast` |
| `popui.ButtonCopy`, `props.ButtonCopy` | `popui.Button` with `props.Button.Copy` |
| `popui.ButtonCopyLink`, `props.ButtonCopyLink` | compose two `popui.Button`s (copy + link) |
| `popui.FileBadge`, `props.FileBadge` | `popui.Avatar` with `Color` + initials |
| `popui.TagStatusIcon`, `props.TagStatusIcon` | renamed `popui.StatusBadge`, `props.StatusBadge` |
| `popui.DropdownSelect`, `props.DropdownSelect`, `props.DropdownSelectOption` | none — see entry; option struct is `props.FilterOption` |
| `popui.SidePanelList` family | `popui.DescriptionList` + `popui.DescriptionListItem` |
| `popui.SidebarCollapsibleSection` | `popui.SidebarSection` (non-collapsing) |
| `popui.Page`, `PageContainer`, `PageHeader`, `PageContent`, `PageSections`, `PageSection`, `PageActions`, `PageTitle` | `popui.App`, `Header`, `Main`, `Section`, `Title` |
| `popui.PopupLayout`, `props.PopupLayout` | `popui.App` |
| `popui.PopupConfigContainer`, `popui.PopupConfigFooter` | `popui.Main`, `popui.Footer` |
| `props.WizardHeader`, `props.WizardContent`, `props.WizardFooter` | none (never had components) |
| `scripts` package (`scripts.ButtonCopy`) | `popui.EmbeddedJS()` (popui.js covers copy buttons) |

### CardFile / CardFileInfo → FileDownload / FileDownloadInfo

Same structure, new names in the File family, plus `Borderless`/`Hover`
variants and a `Value` secondary line on the info slot.

```templ
// before
@popui.CardFile(props.CardFile{}) {
	@popui.CardFileInfo(props.CardFileInfo{Label: "invoice.xml"}) { ... }
}
// after
@popui.FileDownload() {
	@popui.FileDownloadInfo(props.FileDownloadInfo{Label: "invoice.xml"}) { ... }
}
```

### CardFieldset → FieldsetCard

Rename only — same props (`Title`, `Description`), now in the Fieldset family.

### FlashMessage → Toast

`Toast` is the notification system: dark floating panel with `Type`
(`success`, `error`, `info`), optional `Description`, `ActionLabel`/
`ActionHref`, `Position` and `Duration`. Render it with an `ID`, then show it
with `popui.showToast('<id>')` or a `data-toast-trigger="<id>"` attribute —
unlike FlashMessage it does not auto-show on render.

```templ
// before
@popui.FlashMessage(props.FlashMessage{Type: "success", Message: "Saved"})
// after
@popui.Toast(props.Toast{ID: "saved-toast", Type: props.ToastTypeSuccess, Message: "Saved"})
// + trigger from JS: popui.showToast('saved-toast')
```

The server-side `flash` helper package (`popui.go/flash`) is unchanged and
works with Toast.

### ButtonCopy → Button with Copy

Copy-to-clipboard is a `Button` capability now.

```templ
// before
@popui.ButtonCopy(props.ButtonCopy{Value: id, PrefixLength: 8, SuffixLength: 4})
// after
@popui.Button(props.Button{Copy: id, CopyPrefixLength: 8, CopySuffixLength: 4})
```

The `Variant: "ghost"` of ButtonCopy maps to `props.ButtonVariantTransparent`.
Copy buttons default to `type="button"`, as before.

### ButtonCopyLink → two Buttons

```templ
// before
@popui.ButtonCopyLink(props.ButtonCopyLink{Value: id, URL: url})
// after
<span class="inline-flex items-center gap-1">
	@popui.Button(props.Button{Copy: id})
	@popui.Button(props.Button{Variant: props.ButtonVariantTransparent, Size: props.ButtonSizeIcon, Href: url, Target: "_blank", Rel: "noopener noreferrer"}) {
		@icons.ExternalLink()
	}
</span>
```

### FileBadge → Avatar

`Avatar` renders up to three characters and takes a `Color` foreground token
name (`document-xml`, `document-pdf`, `document-png`, `accent`, `success`,
`warning`, `critical`, `info`).

```templ
// before
@popui.FileBadge(props.FileBadge{Type: "xml"})
// after
@popui.Avatar(props.Avatar{Initial: "XML", Color: "document-xml", Size: props.AvatarSizeLarge, Class: "rounded-md font-mono"})
```

### TagStatusIcon → StatusBadge

Pure rename: component, props type, and the `data-tag-status-icon-*`
attributes (now `data-status-badge-*`). Same `Status` vocabulary
(`success`, `failed`, `warning`, `running`) and `Label` prop.

### DropdownSelect — removed without direct replacement

Built for filter rows and superseded by `popui.Filter`'s inline option lists;
never adopted elsewhere. If you used it as a filter, use `popui.Filter` with
`props.FilterInput{Type: props.FilterTypeSelect}` (or `FilterTypeMultiple`).
`props.DropdownSelectOption` lives on as `props.FilterOption` (same fields:
`Value`, `Label`, `Color`). The `dropdownSelect` Alpine controller was removed
from popui.js. For a plain form select, use `popui.Select`.

### SidePanelList family → DescriptionList

`SidePanelList`, `SidePanelItem`, `SidePanelRow`, `SidePanelHeading` and
`SidePanelEmpty` are covered by `DescriptionList` detail rows — see the
"Detail View" example on the Description List docs page.

```templ
// before
@popui.SidePanelList() {
	@popui.SidePanelItem(props.SidePanelItem{Label: "Role", Value: "Seller"})
}
// after
@popui.DescriptionList() {
	@popui.DescriptionListItem(props.DescriptionListItem{Inline: true, Label: "Role", Value: "Seller"})
}
```

Copyable rows (`SidePanelRow{Copyable: true}`) map to
`DescriptionListItem.Actions` with a `Copy` action. Headings and empty states
were plain styled text — an `<h3 class="py-3 text-base font-medium
text-foreground">` and a `<p class="text-base
text-foreground-default-secondary">` respectively.

### SidebarCollapsibleSection — removed

Use `popui.SidebarSection` (non-collapsing). If collapsing is essential, wrap
items in a native `<details>` element.

### Page family / PopupLayout / PopupConfig — removed legacy layouts

All were deprecated aliases of the App layout system:

- `Page`, `PopupLayout`, `PopupConfigContainer` → `popui.App` (page shell)
- `PageHeader`, `PageActions` → `popui.Header` (breadcrumbs + action slot)
- `PageContent`, `PageContainer`, `PageSections` → `popui.Main`
- `PageSection` → `popui.Section` (same `Title`/`Description` props)
- `PageTitle` → `popui.Title`
- `PopupConfigFooter` → `popui.Footer`

### scripts package — removed

`scripts.ButtonCopy()` (a standalone `<script>` for copy buttons) is
redundant: `popui.EmbeddedJS()` already wires copy buttons via popui.js.

## Docs-site anchors (same release)

Sidebar groups changed: `#guides-*` pages are now `#foundations-*`
(`#guides-icons` → `#foundations-icons`, `#guides-tokens` →
`#foundations-tokens`) and `#layout-*` pages are now `#components-*`
(`#layout-card` → `#components-card`, etc.).
