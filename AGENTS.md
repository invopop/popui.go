# PopUI Build & Architecture

## Project Overview

Go component library using [templ](https://templ.guide) for HTML templating and Tailwind CSS v4 for styling. The dev server runs via `air` on port 3000.

## Build Pipeline

There are two independent build steps:

### 1. Tailwind CSS (manual)

```bash
tailwindcss -i styles.css -o assets/popui.css --minify
```

Run this whenever you add new Tailwind classes that aren't already in `assets/popui.css`. The CSS is **embedded** into the Go binary via `//go:embed assets/*` in `assets.go`, so a Go rebuild is also needed after regenerating CSS.

After rebuilding CSS, trigger an air rebuild by touching any `.templ` file:

```bash
tailwindcss -i styles.css -o assets/popui.css --minify && touch layout.templ
```

If dynamically constructed class names are needed (e.g. `"shadow-" + name`), add them to the `@source inline(...)` safelist in `styles.css`.

### 2. Go + Templ (automatic via air)

Air watches `.templ`, `.go`, and other source files and runs:

```
templ generate && go build ./cmd/popui
```

Air does **not** rebuild CSS. It excludes `assets/popui.css` from its watch list.

### Prerequisites

The `tailwindcss` CLI must be installed on your system. Install via npm globally:

```bash
npm install -g @tailwindcss/cli
```

There are **no npm dependencies** in this project. The Tailwind plugins (`@tailwindcss/forms`, `@tailwindcss/typography`) and their runtime dependencies are vendored in the `tailwind/` directory. The `styles.css` entry point references them via local paths:

```css
@import "./tailwind/node_modules/tailwindcss" source(none);
@plugin "./tailwind/forms.js";
@plugin "./tailwind/typography/index.js";
```

To update vendored plugins, replace the files in `tailwind/` with newer versions from npm and test with a full CSS rebuild.

## Key Directories

- `*.templ` (root) - Core PopUI components
- `props/` - Go structs for component props
- `icons/` - Auto-generated icon components (via `go generate` from SVG source files)
- `internal/docs/` - Documentation site pages and component examples
- `internal/docs/components/` - Individual component doc pages
- `internal/docs/examples/` - Example usage for each component
- `examples/` - Full application examples (admin, app, console, wizard, etc.)
- `cmd/popui/serve.go` - Dev server routes
- `assets/popui.css` - Compiled CSS (embedded in binary, do not edit directly)
- `styles.css` - Tailwind CSS entry point
- `tailwind.theme.css` - Design tokens (colors, shadows, spacing)
- `tailwind/` - Vendored Tailwind plugins (@tailwindcss/forms, @tailwindcss/typography) and their dependencies
- `components.css` - Custom CSS beyond Tailwind utilities

## Icons

Icons are generated from SVG files. The generator at `icons/generate.go` produces:
- `icons/icons_list.templ` - Icon templ components
- `internal/docs/components/icons_list.templ` - Docs icon gallery page

Run with `go generate ./icons/` (requires the popui SVG source repo at `../../popui/icons/themes`).

When adding icons manually to `icons/icons_list.templ`, they are automatically included in the docs gallery.

## Templ Conventions

- `.templ` files contain the source templates
- `_templ.go` files are auto-generated (never edit manually)
- Components accept props via variadic `...props.Type` pattern
- Use `@popui.ComponentName(props)` syntax to compose components
