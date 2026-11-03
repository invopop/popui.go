## PopUI Go

This Go package uses Templ for creating template wrappers to define Popui components and use as much as possible the shared styles with the base popui svelte components.

See the [live demo site](https://popui-go.netlify.app/) for code examples and usage details.

Upgrading and a component no longer exists? [MIGRATIONS.md](./MIGRATIONS.md) maps every removed or renamed component to its replacement.

### Using the Go Package

Ensure the public assets are loaded using something like the following example for the go echo library:

```go
e.StaticFS(popui.AssetPath, popui.Assets)
```

Build your own layout using the `popui.App` as a base:

```go
import (
	"github.com/invopop/popui.go"
	"github.com/invopop/popui.go/props"
)

templ Page() {
    @popui.App(props.App{Title: "Test App"}) {
        @popui.Main() {
            @popui.Article() {
                @popui.Button(
                    props.Button{
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

Examples of each component are in `internal/docs/examples` and on the documentation site.

### JSON UI

The experimental `jsonui` package renders panels of PopUI components described in JSON, in Go or in the browser through WebAssembly. Its usage is described on the documentation site's JSON UI page, and the `/jsonui/` playground on the development server is the quickest way to try it.

## Development

JSON UI action and validation runtime tests use Node.js without additional dependencies:

```bash
node --test jsonui/tests/*.test.cjs
```

### Building

In order to generate the CSS, the `tailwindcss` application must be installed in your system. In general we recommend using brew for this:

```bash
brew install tailwindcss
```

We're not expecting many theme changes, but periodically it may be useful to copy the `tailwind.theme.css` file from the [popui](https://github.com/invopop/popui) project to the root of this one so that resource is available.

To generate all the resources, run:

```bash
go generate ./...
```

### Documentation site

The documentation site is served locally with `go run ./cmd/popui serve`, and built into `public/` with `go run ./cmd/popui build`, as Netlify does (see `netlify.toml`). Both need the JSON UI playground's editor, a CodeMirror bundle built with npm that isn't committed. Build it after cloning, and again when its `package-lock.json` changes:

```bash
cd internal/docs/playground/editor && npm ci && npm run build
```

### Editing components

1. Edit the `.templ` files.
2. Run `templ generate` to regenerate the `*_templ.go` files. Use the templ version in `go.mod`, as other versions generate code that may not build against it:

   ```bash
   go run github.com/a-h/templ/cmd/templ@$(go list -m -f '{{.Version}}' github.com/a-h/templ) generate
   ```

3. If templates use new utility classes, run `go generate ./...` to rebuild `assets/popui.css`.
4. Run `go build ./...` and `go test ./...`.

Generated files (`*_templ.go`, `assets/popui.css` and `jsonui/schema/panel.json`) are checked in, so commit them with the source changes.

When removing or renaming a component or props type, add an entry to [MIGRATIONS.md](./MIGRATIONS.md) in the same change, with before and after snippets.

### Design

Read [DESIGN.md](./DESIGN.md) before producing or editing any UI. It defines the component inventory, semantic token system, layout templates, and the do/don't list, and is authoritative. In short:

- Use existing PopUI components instead of raw `<button>`, `<input>`, `<select>`, `<form>`, `<table>`, `<nav>`, etc.
- Use semantic tokens (`text-foreground`, `bg-background`, `border-border`, `text-foreground-default-secondary`, `bg-background-default-secondary`, …), never raw colours such as `bg-red-50`, `text-gray-500`, `text-white` or hex values. Token names are exactly as listed in DESIGN.md §2: there is no `text-foreground-default` or `border-border-default`.
- Use only the 400, 500 and 600 font weights. No `font-bold`.
- A "three-dot menu" is `popui.Menu`, and company or user logos with initials are `popui.Avatar`.

The templates at the repository root (`button.templ`, `card.templ`, …) define the primitives, so they use raw HTML and Tailwind utilities directly. The rules apply to everything that uses them, including `internal/docs/examples/` and `internal/docs/components/`.

When working from a Figma design, the Figma file is the source of truth for sizes, colours, spacing and structure, and DESIGN.md for the vocabulary: translate Figma values into PopUI components and tokens rather than copying raw hex or pixel values.

### Documentation site

Use the popui utility to compile the project into the `public` directory with a single `index.html` for the development site.

```bash
go build ./cmd/popui && ./popui build
```

Run a development server to be able to quickly test changes to the templates using the air utility, this will ensure the Tailwind classes, generated code, and binary are correctly prepared. It serves the site on port 3000, and must run from the repository root so the JSON UI playground can build its WebAssembly.

```bash
air
```

To show a component on the site:

1. Component source: `<name>.templ` and `props/<name>.go` at the repository root.
2. Example: `internal/docs/examples/<name>.templ` defining `templ <Name>Example()`.
3. Page: `internal/docs/components/<name>.templ` defining `templ <Name>()`, which wraps each example in `modules.Example` and ends with a `modules.Section("API Reference", "api")` of `modules.APITable` entries.
4. Sidebar entry: append a `Page` to the appropriate group in `internal/docs/ui.go`.
5. Run `templ generate` and `go build ./...`.

### Adding a component to JSON UI

The schema, the documentation site's component table and the playground's completion are generated from the jsonui registry and the props structs, so a component only needs registering:

1. Register it in `init()` in `jsonui/jsonui.go`, e.g. `register("tabs", popui.Tabs)`. This works for components with the standard `templ Name(p ...props.X)` signature; others need their own entry, like `upload_file` in `jsonui/upload.go`.
2. For form controls, add the component to `modelProps` with the prop that holds its value. For components that only accept certain events in `on`, list them in `nodeEvents`.
3. Make sure its props describe themselves, as the schema and strict mode read them:
   - Doc comments on the fields become the schema's descriptions.
   - String props with a fixed set of values list them in a tag, including `""` when the prop is optional: `` Variant string `jsonschema:"enum=,enum=default,enum=pill"` ``.
   - String props that hold URLs are tagged `` `jsonschema:"format=uri-reference"` ``, so strict mode rejects `javascript:` and similar URLs in them. `templ.SafeURL` props are treated as URLs without a tag.
4. Run `go generate ./jsonui/schema`. `go test ./...` fails while `jsonui/schema/panel.json` is out of date.
