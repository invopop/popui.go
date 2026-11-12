package popui

import (
	"bytes"
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/invopop/popui.go/props"
)

func renderTabs(t *testing.T, p props.Tabs) string {
	t.Helper()
	var buf bytes.Buffer
	if err := Tabs(p).Render(context.Background(), &buf); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

func TestTabsLinks(t *testing.T) {
	out := renderTabs(t, props.Tabs{Variant: "primary", Active: "doc", Tabs: []props.TabItem{
		{Value: "summary", Label: "Summary", Href: "/e/1"},
		{Value: "doc", Label: "Document", Href: "/e/1/document"},
	}})
	for _, want := range []string{
		`role="tablist" class="[^"]*popui-tabs-primary`,
		`<a href="/e/1/document" role="tab"[^>]*data-state="active" aria-selected="true" aria-current="page"`,
		`<a href="/e/1" role="tab"[^>]*data-state="inactive" aria-selected="false">`,
	} {
		if !regexp.MustCompile(want).MatchString(out) {
			t.Errorf("missing %s in:\n%s", want, out)
		}
	}
	if strings.Contains(out, "<button") {
		t.Errorf("links don't switch panels, so there are no buttons:\n%s", out)
	}
}

// The active tab is drawn before Alpine starts, so the row doesn't change
// as the page loads.
func TestTabsStateBeforeAlpine(t *testing.T) {
	out := renderTabs(t, props.Tabs{Tabs: []props.TabItem{
		{Value: "a", Label: "First"},
		{Value: "b", Label: "Second"},
	}})
	for _, want := range []string{
		`<button type="button" role="tab"[^>]*data-state="active"[^>]*>\s*<span class="truncate">First`,
		`<button type="button" role="tab"[^>]*data-state="inactive"[^>]*>\s*<span class="truncate">Second`,
	} {
		if !regexp.MustCompile(want).MatchString(out) {
			t.Errorf("missing %s in:\n%s", want, out)
		}
	}
}
