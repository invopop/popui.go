# Invopop Design System — PopUI

This document captures the Invopop design system as implemented in **PopUI** (`github.com/invopop/popui.go`). It is intended as a reference for generating accurate, consistent UIs that match the Invopop platform.

---

## ⚠️ Prime Directive: Always Use PopUI Components

**Never implement UI primitives from scratch.** PopUI already provides a complete component library. When generating or editing any Invopop UI — regardless of context, framework, or how the request is phrased — your first instinct must be to reach for the existing PopUI component, not to build something new.

This rule applies even when:
- The request seems simple ("just a button")
- A custom implementation looks equivalent
- You are writing plain HTML/CSS rather than Go Templ
- The request doesn't explicitly mention PopUI

### The Full Component Inventory

Before writing any UI element, check whether PopUI already covers it:

| Need | PopUI Component |
|------|----------------|
| Any clickable action | `Button` (variants: default, primary, secondary, danger, transparent) |
| Icon-only action | `Button` with `Size: "icon"` |
| Group of buttons | `ButtonGroup` |
| Copy-to-clipboard | `ButtonCopy` |
| Page/app shell | `App` |
| Top bar with actions | `Header` |
| Bottom bar with actions | `Footer` |
| Primary content area | `Main` > `Article` > `Block` |
| Side navigation | `Sidebar` > `SidebarSection` > `SidebarItem` |
| Right panel / detail pane | `Aside` |
| Content card | `Card` > `CardHeader` + `CardContent` |
| Segmented content card | `SegmentedCard` |
| Metric display | `CardDashboard` > `CardDashboardItem` |
| Progress bar | `CardProgressBar` |
| File attachment display | `CardFile` > `CardFileInfo` |
| Data table | `Table` |
| Table pagination | `TablePagination` |
| Text input | `Input` |
| Multi-line text | `Textarea` |
| Dropdown / select | `Select` |
| Checkbox | `Checkbox` |
| Toggle switch | `Checkbox` with `Variant: "switch"` |
| Radio buttons | `Radio` |
| File input | `InputFile` |
| Avatar upload | `FileUpload` |
| Form container | `Form` |
| Grouped form fields | `Fieldset` |
| Titled settings section (tinted body) | `FieldsetCard` |
| Grouped checkboxes/radios | `OptionGroup` |
| Text field label | `Label` |
| User avatar | `Avatar` |
| Breadcrumb navigation | `Breadcrumbs` > `Breadcrumb` |
| Expandable section | `Accordion` > `AccordionTrigger` + `AccordionContent` |
| Dropdown/context actions | `ContextMenu` > `ContextMenuItem` |
| Floating content panel | `Popover` |
| Key-value data display | `DescriptionList` > `DescriptionListItem` |
| Tab navigation | `Tabs` |
| Horizontal divider | `Separator` |
| Workflow/invoice status | `TagStatus` |
| Toast / success message | `FlashMessage` |
| In-app notification | `Notification` |
| Empty / error state | `PageState` |
| Country flag | `Flag` |
| Range slider | `Slider` |
| Image | `Image` |
| Inline icon | Any icon from `github.com/invopop/icons` |

### When It's OK to Build Custom

Only build a custom component when **all** of the following are true:

1. You have checked the inventory above and no existing component covers the need
2. The user has explicitly asked for something new or described a component that doesn't exist in PopUI
3. The custom element is composed using PopUI tokens (semantic color classes, spacing, typography) — not raw hex values or arbitrary Tailwind colors

Even then: wrap native HTML elements with PopUI token classes, don't reach for custom CSS or inline styles.

### Recognizing Disguised Requests

Some requests sound custom but map directly to PopUI components:

| Phrasing | Correct component |
|----------|-------------------|
| "a save button" | `Button` with `Variant: "primary"` |
| "a delete/remove button" | `Button` with `Variant: "danger"` |
| "a subtle/ghost action" | `Button` with `Variant: "transparent"` |
| "a settings gear icon button" | `Button` with `Variant: "transparent"`, `Size: "icon"`, `@icons.Settings()` |
| "a toggle / on-off switch" | `Checkbox` with `Variant: "switch"` |
| "a status badge / pill" | `TagStatus` |
| "a modal or drawer" (if side panel) | `Aside` |
| "a dropdown menu / three-dot menu" | `ContextMenu` |
| "a key-value list" | `DescriptionList` |
| "a settings section / form card" | `Fieldset` with `Variant: "card"` (bordered) or `FieldsetCard` (titled, tinted body) |
| "an empty state / zero state" | `PageState` |
| "a success toast" | `FlashMessage` |
| "a collapsible / expandable section" | `Accordion` |
| "a copy button for an ID or hash" | `ButtonCopy` |

---

## 1. Visual Theme & Atmosphere

PopUI is a **light-mode-primary** design system built for productivity applications — invoicing dashboards, workflow UIs, and configuration panels. It reads as clean, professional, and information-dense without feeling corporate or sterile. The visual language is rooted in functional clarity: borders define space, whitespace is purposeful, and color is used semantically rather than decoratively.

The palette is a neutral off-white/off-black range with semantic color layers for status, document types, and accent states. A workspace accent color is configurable via `--workspace-accent-color` CSS variable, but the system defaults to a neutral blue. Corners are **moderately rounded** — not pill-shaped, not sharp — using a tight 4–12px scale. Shadows are minimal and button-specific; depth is primarily achieved through borders.

**Key Characteristics:**
- Light-mode native with a white/near-white canvas
- Inter for all UI text; CommitMono for code and mono labels
- Semantic token system: colors are named by role (`foreground`, `background`, `border`, `icon`), not by value
- Configurable workspace accent color (`--workspace-accent-color`)
- Moderate border radius scale (4px–12px); no pill shapes except in specific badge contexts
- Status-driven color vocabulary: draft, processing, registered, sent, received, paid, completed, error, rejected, void, invalid
- Document-type colors: PDF, XML, PNG each have their own foreground/background/border/icon tokens
- Minimal shadow system — only buttons use shadows; everything else relies on borders

---

## 2. Color System

PopUI uses a **semantic token system** built on Tailwind CSS. Colors are never referenced by hex value in templates — always by token name. This keeps the system themeable and consistent.

### Token Structure

Tokens follow the pattern `{role}-{modifier?}`:
- `foreground` / `foreground-default-secondary` / `foreground-default-tertiary`
- `background` / `background-default-secondary` / `background-default-tertiary`
- `border` / `border-default-secondary` / `border-default-tertiary`
- `icon` / `icon-default-secondary` / `icon-default-bold`

### Foreground Tokens

| Token | Description |
|-------|-------------|
| `foreground` | Default text color (primary) |
| `foreground-default-secondary` | Secondary text |
| `foreground-default-tertiary` | Tertiary / muted text |
| `foreground-inverse` | Text on dark backgrounds |
| `foreground-inverse-secondary` | Secondary text on dark |
| `foreground-accent` | Accent-colored text |
| `foreground-selected` | Selected state text |
| `foreground-success` | Success-state text |
| `foreground-warning` | Warning-state text |
| `foreground-critical` | Critical/error text |
| `foreground-info` | Info text |
| `foreground-attention` | Attention text |

### Background Tokens

| Token | Description |
|-------|-------------|
| `background` | Default page/surface background |
| `background-default-secondary` | Secondary surface (slightly off-white) |
| `background-default-tertiary` | Tertiary surface |
| `background-default-tertiary-hover` | Tertiary hover state |
| `background-default-bold` | Bold/strong background |
| `background-default-negative` | Negative/inverse background |
| `background-accent` | Accent-colored background |
| `background-accent-hover` | Accent hover state |
| `background-selected` | Selected element background |
| `background-selected-hover` | Selected hover |
| `background-success` | Success background |
| `background-success-bold` | Bold success background |
| `background-warning` | Warning background |
| `background-warning-bold` | Bold warning background |
| `background-critical` | Critical/error background |
| `background-critical-bold` | Bold critical background |
| `background-info` | Info background |
| `background-info-bold` | Bold info background |
| `background-attention` | Attention background |

### Border Tokens

| Token | Description |
|-------|-------------|
| `border` | Default border |
| `border-default-secondary` | Secondary border |
| `border-default-secondary-hover` | Secondary border on hover |
| `border-default-tertiary` | Tertiary / subtle border |
| `border-inverse` | Inverse border (on dark) |
| `border-inverse-secondary` | Secondary inverse border |
| `border-selected` | Selected state border |
| `border-selected-bold` | Bold selected border |
| `border-warning` | Warning border |
| `border-warning-bold` | Bold warning border |
| `border-critical-bold` | Critical/error border |

### Icon Tokens

| Token | Description |
|-------|-------------|
| `icon` | Default icon color |
| `icon-default-secondary` | Secondary icon |
| `icon-default-bold` | Bold icon |
| `icon-inverse` | Inverse icon (on dark) |
| `icon-inverse-secondary` | Secondary inverse icon |
| `icon-accent` | Accent icon |
| `icon-selected` | Selected icon |
| `icon-success` | Success icon |
| `icon-warning` | Warning icon |
| `icon-critical` | Critical icon |
| `icon-info` | Info icon |
| `icon-attention` | Attention icon |

### Document Type Tokens

Each document type (PDF, XML, PNG/Image) has its own complete color set:

| Token Pattern | Types |
|---------------|-------|
| `foreground-document-{type}` | pdf, xml, png |
| `background-document-{type}` | pdf, xml, png |
| `border-document-{type}` | pdf, xml, png |
| `icon-document-{type}` | pdf, xml, png |

> **Usage tip:** Use `border border-dashed border-border-document-png bg-background-document-png` for file drop zones and placeholder content areas — this is the canonical pattern seen throughout the codebase.

### Status Colors

Workflow and invoice statuses each have their own token. Use these for `TagStatus` and state indicators:

| Token | Status |
|-------|--------|
| `status-draft` | Draft |
| `status-processing` | Processing |
| `status-registered` | Registered |
| `status-sent` | Sent |
| `status-received` | Received |
| `status-paid` | Paid |
| `status-completed` | Completed |
| `status-error` | Error |
| `status-rejected` | Rejected |
| `status-void` | Void |
| `status-invalid` | Invalid |
| `status-empty` | Empty |

---

## 3. Typography

### Font Families

| Token | Family | Use |
|-------|--------|-----|
| `font-sans` | Inter, sans-serif | All UI text |
| `font-mono` | CommitMono, monospace | Code, technical labels, IDs |

### Font Weights

| Class | Weight | Use |
|-------|--------|-----|
| `font-normal` | 400 | Body, descriptions, labels |
| `font-medium` | 500 | Navigation, emphasized text |
| `font-semibold` | 600 | Headings, titles, strong labels |

No bold (700) in the system. Three weights only.

### Text Size Scale

| Class | Size / Line Height | Use |
|-------|--------------------|-----|
| `text-2xl` | 24px / 32px | Page headings, section titles |
| `text-xl` | 20px / 28px | Card headings, sub-sections |
| `text-lg` | 16px / 24px | Standard body, input labels |
| `text-md` | 15px / 22px | Navigation, secondary body |
| `text-base` | 14px / 20px | Default text, captions, descriptions |
| `text-sm` | 12px / 16px | Small labels, metadata, breadcrumbs |
| `text-xs` | 10px / 12px | Micro labels, badges, footnotes |

> **Note:** Invopop's `text-base` is **14px**, not the Tailwind default of 16px. This is a deliberate design decision — the system runs one size smaller than standard Tailwind defaults.

### Letter Spacing

Letter spacing is tied to size — tighter at larger sizes, slightly looser at smaller:

| Class | Spacing | Paired with |
|-------|---------|-------------|
| `tracking-tightest` | -0.29px | `text-2xl` |
| `tracking-tighter` | -0.24px | `text-xl` |
| `tracking-tight` | -0.16px | `text-lg` |
| `tracking-normal` | -0.07px | `text-base` |
| `tracking-wide` | -0.036px | `text-sm` |

### Typography Principles

- Use `text-foreground` for primary content, `text-foreground-default-secondary` for supporting text, `text-foreground-default-tertiary` for muted/placeholder text
- Apply negative letter-spacing at large sizes (`text-2xl`, `text-xl`) — this is deliberate
- `font-mono` / CommitMono is reserved for code, IDs, hash values, and technical labels. Never use it for regular content

---

## 4. Spacing & Layout

### Border Radius Scale

| Class | Value | Use |
|-------|-------|-----|
| `rounded-sm` / `rounded` | 4px | Small elements: badges, chips |
| `rounded-md` | 6px | Buttons, inputs |
| `rounded-lg` | 8px | Cards, panels |
| `rounded-xl` | 10px | Larger containers |
| `rounded-2xl` | 12px | Feature sections, dashed placeholder areas |

No pill shapes (9999px) for standard UI elements. Pill/full-round is avoided in the product UI.

### Layout Components

The app layout is structured as a composition of semantic regions:

```
App
├── Sidebar (optional, collapsible, 240px fixed)
├── Header (sticky, breadcrumb + actions)
├── Main (scrollable content area)
│   ├── Article (content column, optional FullWidth)
│   │   └── Block (flex grouping within article)
│   └── Table + TablePagination (full-width data)
├── Aside (optional, 400px fixed right panel)
└── Footer (sticky, actions right-aligned)
```

**Key layout rules:**
- `Main` with `Center: true` applies a max-width centered layout
- `Article` without `FullWidth` is a constrained content column
- `Article` with `FullWidth: true` spans full width but keeps horizontal padding
- `Aside` is fixed at 400px width, scrollable, on the right
- `Sidebar` is fixed at 240px, collapses on mobile, supports dark variant
- `Footer` content is always right-aligned (`justify-end`)

### Spacing Scale

Base unit: 4px (Tailwind default).

Common spacings in use:
- Component gaps: `gap-2` (8px), `gap-4` (16px)
- Card padding: `p-4` (16px) to `p-6` (24px)
- Article padding: `p-4` to `p-6`
- Section spacing: `space-y-4` (16px) between form groups
- Button group spacing: `gap-2` between buttons

---

## 5. Shadows

Shadow system is minimal — only buttons use shadows. Cards and containers rely entirely on borders.

| Token | Use |
|-------|-----|
| `shadow` | Default shadow (card-level ambient) |
| `shadow-button-default` | Default/secondary button |
| `shadow-button-primary` | Primary button |
| `shadow-none` | Explicit no-shadow override |

---

## 6. Components

### Buttons

**Variants:** `default` | `primary` | `secondary` | `danger` | `transparent`

**Sizes:** `sm` | `md` (default) | `lg` | `icon`

| Variant | Visual |
|---------|--------|
| `default` | Neutral — muted background, standard border |
| `primary` | Strong — uses accent/bold background, high contrast |
| `secondary` | Outlined — border visible, lower visual weight than primary |
| `danger` | Destructive — critical color background |
| `transparent` | Ghost — no background, no border; for icon buttons and inline actions |

- Buttons are grouped in `ButtonGroup` containers with `gap-2` spacing
- `ButtonGroup` supports `Align` prop: `left` (default), `center`, `right`
- Icon-only buttons use `Size: "icon"` with a single icon child
- Buttons can render as `<a>` by setting `Href`

**Usage patterns:**
- Primary actions in Header slots: `Variant: "primary"` with icon + label
- Destructive actions: `Variant: "danger"` with Delete icon
- Toolbar/contextual: `Variant: "transparent"` with `Size: "icon"`
- Footer actions: right-aligned `ButtonGroup` with primary + secondary

### Cards

Cards are the primary content container. They combine a `Card` > `CardHeader` > `CardContent` structure.

```
Card
├── CardHeader (avatar + title + subtitle)
├── CardContent (flexible content area)
├── CardProgressBar (metric with progress)
├── CardDashboard > CardDashboardItem (metrics grid)
└── CardFile > CardFileInfo (file attachments)
```

- Cards can render as links via `Href` prop
- Cards support `Disabled` state
- `CardHeader` takes an `Avatar` (or custom component) as a child slot
- Card borders use `border` token; no shadows on cards

### Inputs & Forms

**Input sizes:** `sm` | `md` (default) | `lg`

**Input types:** text, email, password, date, and all standard HTML types

Form structure:
```
Form
└── Fieldset (optional card variant) or FieldsetCard (titled, tinted body)
    ├── Input (with Label, Placeholder, optional hint)
    ├── Textarea
    ├── Select
    ├── Checkbox (standard or switch variant)
    ├── Radio
    ├── InputFile / FileUpload
    └── ButtonGroup (submit/cancel)
```

- `Fieldset` with `Variant: "card"` renders as a bordered card — use for grouped settings
- `FieldsetCard` renders a title + description heading above a tinted card body — use for titled settings sections and grouped form blocks
- `Fieldset` with custom `Class: "bg-background-default-secondary"` for highlighted/sudo sections
- `OptionGroup` wraps multiple checkboxes or radios with a shared label
- `Label` can be separate from `Input` for custom hint/link patterns
- `Checkbox` supports `Variant: "switch"` for toggle switches

### Tables

```
Table
└── TablePagination (current/total pages, rows-per-page, item count, custom actions)
```

- Tables are full-width within `Main` (no `Article` wrapper)
- `TablePagination` supports rows-per-page options, navigation links, and custom action slots
- Items label is customizable ("invoices", "orders", etc.)

### Navigation

**Sidebar:**
```
Sidebar (variant: "" | "dark")
├── SidebarHeader (logo + collapse button)
├── SidebarContent
│   └── SidebarSection (title + items)
│       └── SidebarItem (href, selected state)
└── SidebarFooter (user/secondary actions)
```

- `SidebarItem` with `Selected: true` applies active highlight
- Sidebar supports dark variant (`Variant: "dark"`) with inverted colors
- Collapse button uses `data-sidebar-hide` attribute (Alpine.js)
- `SidebarHeader` expects: logo/icon + workspace name + collapse button

**Breadcrumbs:**
```
Breadcrumbs
└── Breadcrumb (label, href, custom children, HTMX attributes)
```

- Last breadcrumb has no `Href` (current page)
- Breadcrumbs support full `Attributes` pass-through for HTMX `hx-get`, `hx-target`, etc.

### Tags & Status

`TagStatus` component for displaying workflow/invoice state. Uses `status-{state}` color tokens.

**States:** draft, processing, registered, sent, received, paid, completed, error, rejected, void, invalid, empty

### Other Components

| Component | Description |
|-----------|-------------|
| `Avatar` | Circular user avatar; `Size: "lg"` or default (small). Accepts `Initial` text or `Image` child |
| `Accordion` | Native `details/summary` expandable sections |
| `ContextMenu` | Trigger button + dropdown menu; supports `RightAlign` |
| `DescriptionList` | `dl`-based term/value pairs for data display |
| `Separator` | Dashed horizontal divider |
| `FlashMessage` | Toast-style success/error feedback |
| `ButtonCopy` | Copy-to-clipboard with truncation (prefix/suffix lengths) |
| `Popover` | Floating content panel |
| `Notification` | In-app notification display |
| `PageState` | Empty/error state with illustration, title, description, CTA |
| `Slider` | Range input |
| `Tabs` | Tab navigation component |
| `Flag` | Country flag via ISO 3166-1 alpha-2 code (requires `flag-icons` CSS) |

---

## 7. Icons

Icons come from `github.com/invopop/icons`. They render as inline SVG components and should always be used alongside text labels in buttons (except icon-only `Size: "icon"` buttons).

**Full icon set includes (non-exhaustive):**
Add, Alert, Archive, ArrowDown/Left/Right/Up, Batch, Bell, Billing, Bot, Box, Brackets, Calendar, CheckBadge, CheckCircle, ChevronDown/Left/Right/Up, Close, Code, Command, Dashboard, Delete, Download, Edit, Envelope, Exchange, ExternalLink, Failure, Filter, Flag, Folder, Hashtag, History, Info, Invoice, Invopop, Key, List, Loader, Lock, Logout, Menu, Note, Notification, Ok, Options, Order, Payment, Pin, Preview, Published, Pulse, Queue, Receipt, Reload, Replace, Reset, Rocket, Running, Save, Search, Send, Settings, Sign, Skip, Slider, SortAscending, SortDescending, SquareCheck, Stack, Stamp, Status, Stop, Success, Support, Tables, Tag, Team, Tick, Upload, User, Van, View, Warning, WarningBold, Workflow, Workspace

**Usage:**
```go
@icons.Add()
@icons.Delete()
@icons.Settings()
```

---

## 8. App Layout Templates

These are the **canonical, mandatory layout templates** for all Invopop apps. Do not invent a new layout structure — always start from the appropriate template below and adapt the content inside it.

### How to choose a template

| Situation | Template to use |
|-----------|----------------|
| Default app, list view, dashboard, settings page | **Basic App** |
| Creating, adding, editing, wizard, multi-step flow | **App Popup** |
| Any page whose primary content is a data table | **App with Table and Pagination** |

---

### Template 1: Basic App _(default for most pages)_

Use this as the starting point for any standard page: dashboards, list views, detail views, settings, and any content page that isn't a creation/edit flow or a full table.

```go
package examples

import (
	"github.com/invopop/icons"
	"github.com/invopop/popui.go"
	"github.com/invopop/popui.go/props"
)

templ AppExample() {
	@popui.App(props.App{Title: "Application"}) {
		@popui.Header(props.Header{
			Breadcrumbs: []props.Breadcrumb{{Label: "Example"}},
		}) {
			@popui.Button(props.Button{Variant: "primary"}) {
				@icons.Add()
				<span>New Item</span>
			}
		}
		@popui.Main(props.Main{Center: true}) {
			@popui.Article() {
				@popui.Block(props.Block{
					Class: "items-center justify-center border border-dashed border-border-document-png bg-background-document-png rounded-2xl text-foreground-default-tertiary",
				}) {
					This is where the page contents goes.
				}
			}
		}
		@popui.Footer() {
			<span class="text-sm text-foreground-default-secondary">&copy; 2025 Invopop Inc.</span>
		}
	}
}
```

**Key rules for this template:**
- `Header` always receives `Breadcrumbs` (at minimum one item for the current page)
- Primary action (e.g. "New X") goes as a child of `Header`, right-aligned automatically
- `Main` uses `Center: true` to constrain content width
- `Footer` contains copyright or secondary navigation — not form actions
- Replace the `Block` placeholder with actual page content

---

### Template 2: App Popup _(creation, edit, wizard, and add flows)_

Use this for any flow where the user is creating or editing something: new invoice, edit settings, onboarding wizard, multi-step form. The close button in the header lets the user cancel and return.

```go
package examples

import (
	"github.com/invopop/icons"
	"github.com/invopop/popui.go"
	"github.com/invopop/popui.go/props"
)

templ AppPopupExample() {
	@popui.App(props.App{Title: "Application Popup"}) {
		@popui.Header(props.Header{
			Title: appPopupTitle(),
		}) {
			@popui.Button(props.Button{Variant: props.ButtonVariantSecondary}) {
				Add New
			}
		}
		@popui.Main() {
			@popui.Article(props.Article{FullWidth: true}) {
				<div class="w-full flex flex-1 items-center justify-center min-h-10 p-4 gap-2 border border-dashed border-border-document-png bg-background-document-png rounded-2xl text-foreground-default-tertiary">
					This is where the page contents goes.
				</div>
			}
		}
		@popui.Footer() {
			@popui.Button(props.Button{Variant: props.ButtonVariantPrimary}) {
				Next
			}
		}
	}
}

templ appPopupTitle() {
	@popui.Button(props.Button{
		Variant: props.ButtonVariantTransparent,
		Class:   "p-2",
	}) {
		@icons.Close()
	}
	<span>
		Popup Example
	</span>
}
```

**Key rules for this template:**
- `Header` uses the `Title` prop (not `Breadcrumbs`) — the title slot contains a close/back button + page name
- The close button is always `Variant: "transparent"` with `@icons.Close()` and `Class: "p-2"`
- `Main` does **not** use `Center: true` — content is full-width inside `Article{FullWidth: true}`
- `Footer` contains the primary submit/next action (right-aligned automatically) — not copyright
- The footer primary button label reflects the action: "Save", "Next", "Create", "Submit", etc.
- Secondary actions (cancel, back) go in Footer as `Variant: "secondary"` before the primary

---

### Template 3: App with Table and Pagination _(list pages with tabular data)_

Use this whenever the primary content of a page is a data table. Do not wrap `Table` inside `Article` — it sits directly inside `Main`.

```go
package examples

import (
	"github.com/invopop/icons"
	"github.com/invopop/popui.go"
	"github.com/invopop/popui.go/props"
)

templ AppWithTableExample() {
	@popui.App(props.App{Title: "Application With Table"}) {
		@popui.Header(props.Header{
			Breadcrumbs: []props.Breadcrumb{{Label: "Example"}},
		}) {
			@popui.Button(props.Button{Variant: "primary"}) {
				@icons.Add()
				<span>New Item</span>
			}
		}
		@popui.Main() {
			@popui.Table() {
				<thead>
					<tr>
						<th>Invoice ID</th>
						<th>Customer</th>
						<th>Amount</th>
						<th>Status</th>
					</tr>
				</thead>
				<tbody>
					<tr>
						<td>INV-001</td>
						<td>Acme Corp</td>
						<td>$1,234.56</td>
						<td>Paid</td>
					</tr>
				</tbody>
			}
			@popui.TablePagination(props.TablePagination{
				CurrentPage:        1,
				TotalPages:         260,
				TotalItems:         25991,
				RowsPerPage:        100,
				RowsPerPageOptions: []int{10, 25, 50, 100, 250},
				ShowRowsPerPage:    true,
				ItemsLabel:         "invoices",
				FirstPageURL:       "?page=1",
				PrevPageURL:        "?page=0",
				NextPageURL:        "?page=2",
				LastPageURL:        "?page=260",
			}) {
				@popui.ButtonGroup() {
					@popui.Button(props.Button{
						Type:    "button",
						Size:    "sm",
						Variant: "secondary",
					}) {
						Create Invoice
					}
				}
			}
		}
	}
}
```

**Key rules for this template:**
- `Table` is a direct child of `Main` — never wrap it in `Article` or `Block`
- `TablePagination` immediately follows `Table`, also a direct child of `Main`
- Always populate `ItemsLabel` with the entity name (e.g. "invoices", "workflows", "suppliers")
- Always include `RowsPerPageOptions` and `ShowRowsPerPage: true`
- The `TablePagination` child slot accepts a `ButtonGroup` for contextual actions (e.g. bulk actions, export)
- `Main` does **not** use `Center: true` — tables are full-width
- `Header` uses `Breadcrumbs`, same as Template 1

---

## 9. Do's and Don'ts

### Do
- Always use semantic tokens (`text-foreground`, `bg-background`, `border-border`) — never raw hex or Tailwind color utilities
- Use `border border-dashed border-border-document-png bg-background-document-png rounded-2xl text-foreground-default-tertiary` for placeholder/empty content zones — this is the canonical pattern
- Use `gap-2` between buttons in `ButtonGroup`, `space-y-4` between form fields
- Apply `tracking-tightest` at `text-2xl`, `tracking-tighter` at `text-xl` — negative tracking at large sizes is intentional
- Use `font-mono` (CommitMono) for IDs, hashes, codes, and technical values
- Use `TagStatus` with `status-{state}` tokens for any workflow/invoice state display
- Wrap form fields in `Fieldset` with `Variant: "card"` or `FieldsetCard` for grouped settings sections
- Use `Variant: "transparent"` + `Size: "icon"` for toolbar/contextual icon buttons
- Place primary actions in the Header slot (top right); place confirmation actions in Footer (right-aligned)

### Don't
- **Don't implement buttons, inputs, cards, tables, forms, navigation, or any other UI primitive from scratch** — always use the PopUI component. This applies even for "simple" or "one-off" cases.
- **Don't write raw `<button>`, `<input>`, `<select>`, `<form>`, `<table>`, `<nav>` HTML elements** if a PopUI component exists for it — and one almost always does. Check the component inventory first.
- **Don't recreate a component because the usage looks slightly different** — use the `Class` prop to extend styling, not a full reimplementation.
- Don't use Tailwind color utilities like `text-gray-500`, `bg-blue-100` — use semantic tokens instead
- Don't use font-bold (700) — the system uses 400, 500, 600 only
- Don't use `rounded-full` / pill shapes for product UI elements
- Don't use CommitMono for UI labels, navigation, or headings — Inter only
- Don't add box-shadows to cards or containers — only buttons use shadows
- Don't add color sections or gradient backgrounds — the product UI is white/neutral throughout
- Don't place Footer content left-aligned — Footer is always `justify-end`
- Don't use inline hex colors or opacity variants — the token system handles all states

---

## 10. Agent Prompt Guide

### Quick Token Reference
```
Primary text:        text-foreground
Secondary text:      text-foreground-default-secondary
Muted text:         text-foreground-default-tertiary
Page background:    bg-background
Secondary surface:  bg-background-default-secondary
Default border:     border border-border
Accent:             (workspace accent color via CSS var)
Success:            text-foreground-success / bg-background-success
Warning:            text-foreground-warning / bg-background-warning
Critical:           text-foreground-critical / bg-background-critical
```

### Common Class Patterns

**Placeholder / empty zone:**
```
border border-dashed border-border-document-png bg-background-document-png rounded-2xl text-foreground-default-tertiary
```

**Muted label:**
```
text-sm text-foreground-default-secondary
```

**Code / ID display:**
```
font-mono text-sm text-foreground
```

**Page heading:**
```
text-2xl font-semibold tracking-tightest text-foreground
```

**Section subtitle:**
```
text-base text-foreground-default-secondary
```

**Standard card:**
```
rounded-lg border border-border bg-background p-4
```

**Highlighted/sudo fieldset:**
```
bg-background-default-secondary (passed as Class to Fieldset with card variant)
```

### Example Component Compositions

**Primary action button in header:**
```go
@popui.Button(props.Button{Variant: "primary"}) {
    @icons.Add()
    <span>New Invoice</span>
}
```

**Status tag:**
```go
@popui.TagStatus(props.TagStatus{Status: "paid"})
```

**Empty state:**
```go
@popui.PageState(props.PageState{
    Illustration: popui.EmptyStateIcon(),
    Title: "No invoices found",
    Description: "Create your first invoice to get started.",
}) {
    @popui.Button(props.Button{Variant: "primary"}) { Create Invoice }
}
```

**Settings fieldset:**
```go
@popui.Fieldset(props.Fieldset{
    Legend: "Billing Settings",
    Variant: props.FieldsetVariantCard,
}) {
    @popui.Input(props.Input{Label: "Company Name", Placeholder: "Acme Corp"})
    @popui.Input(props.Input{Label: "Tax ID", Placeholder: "ESB12345678"})
}
```

**Titled settings section (tinted body):**
```go
@popui.FieldsetCard(props.FieldsetCard{
    Title: "URLs",
    Description: "Endpoints used by the integration",
}) {
    @popui.Input(props.Input{Label: "Config URL", Placeholder: "https://"})
    @popui.Input(props.Input{Label: "Launch URL", Placeholder: "https://"})
}
```

**Sidebar section:**
```go
@popui.SidebarSection(props.SidebarSection{Title: "Workflows"}) {
    @popui.SidebarItem(props.SidebarItem{Href: "/workflows", Selected: true}) { All Workflows }
    @popui.SidebarItem(props.SidebarItem{Href: "/workflows/active"}) { Active }
}
```
