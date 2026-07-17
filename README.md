## PopUI Go

PopUI Go is Invopop's UI toolkit for Go web applications. It uses
[Templ](https://templ.guide) to define PopUI components as template wrappers,
sharing styles with the base [popui](https://github.com/invopop/popui) Svelte
components.

See the [live demo site](https://popui-go.netlify.app/) for every component
with code examples and API references.

Two companion documents are worth knowing about:

- [DESIGN.md](./DESIGN.md) — the design system reference: component
  inventory, semantic token vocabulary, and layout rules. Authoritative when
  producing or reviewing UI built with PopUI.
- [MIGRATIONS.md](./MIGRATIONS.md) — upgrading and a component no longer
  exists? It maps every removed or renamed component to its replacement,
  with before/after snippets.

### Using the Go Package

```bash
go get github.com/invopop/popui.go
```

Serve the embedded assets (CSS, JS, fonts), for example with echo:

```go
e.StaticFS(popui.AssetPath, popui.Assets)
```

Then build your layout on `popui.App`:

```go
import (
	"github.com/invopop/popui.go"
	"github.com/invopop/popui.go/props"
)

templ Page() {
	@popui.App(props.App{Title: "Test App"}) {
		@popui.Main() {
			@popui.Article() {
				@popui.Button(props.Button{
					Variant: props.ButtonVariantPrimary,
					Size:    props.ButtonSizeSmall,
				}) {
					Click Me
				}
			}
		}
	}
}
```

### Editing components

1. Edit the `.templ` files and run `templ generate`.
2. Run `go generate ./...` to rebuild the compiled CSS (requires the
   `tailwindcss` CLI, e.g. `brew install tailwindcss`).
3. `go build ./...` to verify.

Generated artifacts (`*_templ.go`, `assets/popui.css`) are checked in —
commit them alongside source changes. The design tokens come from
`tailwind.theme.css`, copied periodically from the
[popui](https://github.com/invopop/popui) project.

### Documentation site

Compile the docs site into the `public` directory with a single
`index.html`:

```bash
go build ./cmd/popui && ./popui build
```

Or run a development server with the [air](https://github.com/air-verse/air)
utility, which regenerates Tailwind classes, templ code, and the binary on
change:

```bash
air
```

### Examples

The `examples` directory contains standalone pages served by the dev server
under `/examples/...`:

- **Admin** — an admin panel with sidebar navigation, tables, and actions.
- **App** — a simple app layout with header, main content, and footer.
- **Console** — a recreation of the Invopop Console: document tables with
  filtering and an entry view embedding a CodeMirror GOBL editor.
- **Prose** — a prose-styled article page.
- **Wizard** — a multi-step onboarding flow.
