## PopUI Go

This Go package uses Templ for creating template wrappers to define Popui components and use as much as possible the shared styles with the base popui svelte components.

See the [live demo site](https://popui-go.netlify.app/) for code examples and usage details.

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

Use the popui utility to compile the project into the `public` directory with a single `index.html` for the development site.

```bash
go build ./cmd/popui && ./popui build
```

Run a development server to be able to quickly test changes to the templates using the air utility, this will ensure the Tailwind classes, generated code, and binary are correctly prepared.

```bash
air
```

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
                        Variant: "primary",
                        Small: true,
                    }) {
                    Click Me
                }
            }
        }
    }
}
```

You can check out examples of each component in `examples` path and documentation project.
