package popui

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/invopop/popui.go/props"
)

func renderWith(t *testing.T, c, children templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	ctx := context.Background()
	if children != nil {
		ctx = templ.WithChildren(ctx, children)
	}
	if err := c.Render(ctx, &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// TestDetailsStacked checks that a stacked Detail shows its children as its
// value, and a dash without any.
func TestDetailsStacked(t *testing.T) {
	html := renderWith(t, Details(props.Details{Layout: props.DetailsLayoutStacked}), templ.Join(
		renderedDetail(props.Detail{Label: "Bio"}, templ.Raw("<em>Engineer</em>")),
		Detail(props.Detail{Label: "Folder"}),
	))
	for _, want := range []string{"Bio</dt>", "<em>Engineer</em></dd>", "—</dd>", "space-y-4"} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in:\n%s", want, html)
		}
	}
	if strings.Contains(html, "popui-detail-row") {
		t.Errorf("stacked details have rows in:\n%s", html)
	}
}

// TestDescriptionListDeprecated checks the deprecated components look as
// they did: a 400px stacked list, whose items place their children between
// the label and the value.
func TestDescriptionListDeprecated(t *testing.T) {
	html := renderWith(t, DescriptionList(), templ.Join(
		renderedItem(props.Detail{Value: "Administrator"}, DT()),
		DescriptionListItem(props.Detail{Label: "Name"}),
	))
	for _, want := range []string{"w-[400px]", "space-y-4"} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in:\n%s", want, html)
		}
	}
	if strings.Contains(html, "—") || strings.Contains(html, "popui-detail-row") {
		t.Errorf("deprecated list changed its items in:\n%s", html)
	}
	if strings.Count(html, "<dt") != 2 || !strings.Contains(html, "</dt><dd") {
		t.Errorf("children not between label and value in:\n%s", html)
	}
}

func renderedDetail(p props.Detail, children templ.Component) templ.Component {
	return withKids(Detail(p), children)
}

func renderedItem(p props.Detail, children templ.Component) templ.Component {
	return withKids(DescriptionListItem(p), children)
}

func withKids(c, children templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		return c.Render(templ.WithChildren(ctx, children), w)
	})
}
