# Claude Notes — popui.go

## Removed components

When a component or prop type no longer exists (in this repo or in a consumer
upgrading popui.go), check [MIGRATIONS.md](./MIGRATIONS.md) — it maps every
removed/renamed identifier to its replacement with before/after snippets.
When removing or renaming a component, add an entry there in the same pass.

## Design

**Read [DESIGN.md](./DESIGN.md) before producing or editing any UI.** It defines the component inventory, semantic token system, layout templates, and the do/don't list. Treat it as authoritative — when DESIGN.md and this file disagree, DESIGN.md wins.

**When given a Figma URL, use the Figma MCP** (`mcp__Figma__get_design_context`, `mcp__Figma__get_screenshot`, `mcp__Figma__get_variable_defs`) to fetch the design *before* writing or editing code. The Figma file is the source of truth for pixel sizes, colors, spacing, and structure; DESIGN.md is the source of truth for the *vocabulary* (which PopUI component to use, which semantic token, etc.). Translate Figma values into PopUI tokens — never paste raw hex/px from Figma. The Figma desktop app must be open with the file as the active tab for the MCP to resolve nodes.

Non-negotiables (from DESIGN.md, repeated here so they don't get missed):

- Use existing PopUI components instead of writing raw `<button>`, `<input>`, `<select>`, `<form>`, `<table>`, `<nav>`, etc.
- Use semantic tokens (`text-foreground`, `bg-background`, `border-border`, `text-foreground-default-secondary`, `bg-background-default-secondary`, …). No `bg-red-50`, `text-gray-500`, `text-white`, raw hex, etc.
- Three font weights only: 400, 500, 600. No `font-bold`.
- "Three-dot menu" is `popui.ContextMenu`, not a hand-rolled button.
- Company/user logos with initials are `popui.Avatar`, not a styled `<div>`.
- Token names are exactly as listed in DESIGN.md §2. There is no `text-foreground-default`, no `border-border-default` — it's `text-foreground` and `border-border`.

Exception for component internals: the files under the repo root (e.g. `button.templ`, `card.templ`) are where the primitives are *defined*, so they use raw HTML and Tailwind utilities directly. DESIGN.md's "use PopUI components" rule applies to *consumers* of those primitives — including everything under `internal/docs/examples/` and `internal/docs/components/`.

## Repo workflow

This is a templ-based component library. Editing flow:

1. Edit `.templ` files.
2. Run `templ generate` to regenerate `*_templ.go`.
3. If you touched anything that affects compiled CSS (new utility classes in templates), run `go generate ./...` to rebuild `assets/popui.css` (uses local `tailwindcss` via the directive in `assets.go`).
4. `go build ./...` to verify.

Generated artifacts (`*_templ.go`, `assets/popui.css`) are checked in — commit them alongside source changes.

## Adding a new component to the docs site

To make a component show up in the sidebar:

1. Component source: `<name>.templ` + `props/<name>.go` at repo root.
2. Example for the docs site: `internal/docs/examples/<name>.templ` defining `templ <Name>Example()`.
3. Docs page: `internal/docs/components/<name>.templ` defining `templ <Name>()` — wraps each example in `modules.Example` and ends with an `modules.Section("API Reference", "api")` of `modules.APITable` entries.
4. Sidebar entry: append a `Page` to the appropriate group in `internal/docs/ui.go`.
5. Run `templ generate` and `go build ./...`.

## Dev server

`air` (config in `.air.toml`) — runs `go generate ./... && templ generate && go build ./cmd/popui` then `./popui serve -p 3000`.
