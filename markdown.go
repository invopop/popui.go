package popui

import (
	"context"
	"io"

	"github.com/a-h/templ"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

// markdownRenderer converts the Markdown accepted in component copy into HTML.
// Hard wraps make a single newline a line break, so multi-line text reads as
// written without the author having to think about Markdown's paragraph rules,
// and Linkify turns bare URLs into links. Raw HTML in the source is deliberately
// escaped rather than passed through (goldmark's default, i.e. no WithUnsafe):
// component copy should not be able to inject markup.
var markdownRenderer = goldmark.New(
	goldmark.WithExtensions(
		extension.Linkify,
	),
	goldmark.WithRendererOptions(
		html.WithHardWraps(),
		html.WithXHTML(),
	),
)

// markdown renders s as Markdown, writing the HTML straight to the output. It is
// used for the component text that benefits from light formatting — emphasis,
// code, links, lists and line breaks — such as a Tooltip's description.
func markdown(s string) templ.Component {
	return templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		return markdownRenderer.Convert([]byte(s), w)
	})
}
