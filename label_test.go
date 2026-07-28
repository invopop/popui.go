package popui

import (
	"bytes"
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/invopop/popui.go/props"
)

func renderLabel(t *testing.T, p props.Label, text string) string {
	t.Helper()
	var buf bytes.Buffer
	ctx := templ.WithChildren(context.Background(), templ.Raw(text))
	if err := Label(p).Render(ctx, &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

// TestLabelHostsTheTooltip locks in that the label itself is the tooltip's
// trigger and anchor. Hanging the card off a wrapper around the 16px icon meant
// hovering only worked on the icon, and start alignment lined the card up with
// the icon at the end of the text rather than with the label.
func TestLabelHostsTheTooltip(t *testing.T) {
	openTag := regexp.MustCompile(`<label[^>]*>`)

	t.Run("with a hint", func(t *testing.T) {
		got := renderLabel(t, props.Label{ID: "field", Hint: hintText}, "Field")

		tag := openTag.FindString(got)
		for _, want := range []string{"group/tooltip", "relative"} {
			if !strings.Contains(tag, want) {
				t.Errorf("the label must carry %q to host the card, got %q", want, tag)
			}
		}

		// The card lives inside the label, so the label is what it is
		// positioned against and what hovering it responds to.
		if !regexp.MustCompile(`(?s)<label.*role="tooltip".*</label>`).MatchString(got) {
			t.Error("the card should be rendered inside the label")
		}
		// ...and there is no second wrapper around the icon doing the same job.
		if strings.Contains(got, "inline-block w-fit") {
			t.Error("the icon should not be wrapped in another tooltip container")
		}
		// The icon stays focusable so the card is reachable by keyboard.
		if !strings.Contains(got, `tabindex="0"`) {
			t.Error("the icon should remain focusable")
		}
	})

	t.Run("without a hint", func(t *testing.T) {
		got := renderLabel(t, props.Label{ID: "field"}, "Field")

		tag := openTag.FindString(got)
		if strings.Contains(tag, "group/tooltip") || strings.Contains(tag, "relative") {
			t.Errorf("a label with no hint should not become a positioning context, got %q", tag)
		}
		if strings.Contains(got, `role="tooltip"`) {
			t.Error("no hint means no card")
		}
	})
}
