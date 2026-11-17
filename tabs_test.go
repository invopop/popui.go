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

func renderTabs(t *testing.T, p props.Tabs) string {
	t.Helper()
	var buf bytes.Buffer
	if err := Tabs(p).Render(context.Background(), &buf); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

func TestTabsLinks(t *testing.T) {
	out := renderTabs(t, props.Tabs{Variant: tabsPrimary, Active: "doc", Tabs: []props.TabItem{
		{Value: "summary", Label: "Summary", Href: "/e/1"},
		{Value: "doc", Label: "Document", Href: "/e/1/document", Attributes: templ.Attributes{"hx-target": "#main"}},
	}})
	for _, want := range []string{
		`role="tablist" class="[^"]*\bh-10\b`,
		`<a href="/e/1/document" role="tab"[^>]*data-state="active" aria-selected="true" aria-current="page" hx-target="#main"`,
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

// Primary tabs are a row of cells: the active one has the content's
// background and covers the row's line, and the row is taller when large.
func TestTabsPrimary(t *testing.T) {
	tabs := []props.TabItem{{Value: "a", Label: "First"}, {Value: "b", Label: "Second", Count: 2}}
	out := renderTabs(t, props.Tabs{Variant: tabsPrimary, Size: "lg", Tabs: tabs})
	for _, want := range []string{
		`role="tablist" class="[^"]*\bh-12\b`,
		`role="tab" class="[^"]*data-\[state=active\]:-mb-px[^"]*"[^>]*data-state="active"`,
		`<span class="tabular-nums shrink-0 font-normal text-foreground-default-secondary">\(2\)</span>`,
	} {
		if !regexp.MustCompile(want).MatchString(out) {
			t.Errorf("missing %s in:\n%s", want, out)
		}
	}
	if regexp.MustCompile(`role="tablist" class="[^"]*\bh-10\b`).MatchString(out) {
		t.Errorf("a large row isn't also the default height:\n%s", out)
	}
}
