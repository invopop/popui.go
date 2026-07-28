package popui

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/invopop/popui.go/props"
)

func TestMarkdown(t *testing.T) {
	render := func(t *testing.T, src string) string {
		t.Helper()
		var buf bytes.Buffer
		if err := markdown(src).Render(context.Background(), &buf); err != nil {
			t.Fatalf("render: %v", err)
		}
		return buf.String()
	}

	t.Run("a single newline is a line break", func(t *testing.T) {
		// The whole point of hard wraps: multi-line copy reads as written,
		// instead of Markdown joining the lines into one paragraph.
		got := render(t, "first line\nsecond line")
		if !strings.Contains(got, "<br />") {
			t.Errorf("expected a line break, got %q", got)
		}
	})

	t.Run("lists, emphasis and code render", func(t *testing.T) {
		got := render(t, "- **bold** item\n- `code` item")
		for _, want := range []string{"<ul>", "<li>", "<strong>bold</strong>", "<code>code</code>"} {
			if !strings.Contains(got, want) {
				t.Errorf("expected %q in %q", want, got)
			}
		}
	})

	t.Run("bare urls are linkified", func(t *testing.T) {
		got := render(t, "see https://invopop.com for more")
		if !strings.Contains(got, `<a href="https://invopop.com"`) {
			t.Errorf("expected a link, got %q", got)
		}
	})

	t.Run("raw html is escaped, not passed through", func(t *testing.T) {
		got := render(t, "<script>alert(1)</script>")
		if strings.Contains(got, "<script>") {
			t.Errorf("raw HTML must not survive: %q", got)
		}
	})
}

// TestTooltipDescriptionIsMarkdown covers the wiring: a Tooltip's description
// goes through the Markdown renderer, and the card holds it in a div so block
// level output is valid.
func TestTooltipDescriptionIsMarkdown(t *testing.T) {
	var buf bytes.Buffer
	err := Tooltip(props.Tooltip{
		Title:       "Pricing strategies",
		Description: "first line\nsecond line\n\n- one\n- two",
	}).Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	got := buf.String()
	for _, want := range []string{"<br />", "<ul>", "<li>one</li>"} {
		if !strings.Contains(got, want) {
			t.Errorf("expected %q in the card, got %q", want, got)
		}
	}
	if strings.Contains(got, "<p class=\"text-sm text-foreground-inverse-secondary\">") {
		t.Error("the description must not be wrapped in a p: markdown emits block elements")
	}
}
