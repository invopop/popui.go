package jsonui

import (
	"strings"
	"testing"
)

// TestNumericCells checks that numeric cells are marked for the table's
// styles, aligned to the end unless they say otherwise.
func TestNumericCells(t *testing.T) {
	html := mustRender(t, `{"components": [{"is": "table", "children": [{"is": "tbody", "children": [{"is": "tr", "children": [
		{"is": "td", "props": {"numeric": true}, "children": [{"is": "text", "text": "1.00"}]},
		{"is": "td", "props": {"align": "center"}, "children": [{"is": "text", "text": "x"}]},
		{"is": "td", "props": {"numeric": true, "align": "start"}, "children": [{"is": "text", "text": "2"}]},
		{"is": "td", "children": [{"is": "text", "text": "plain"}]}
	]}]}]}]}`)
	assertContains(t, html,
		`<td data-align="end" data-numeric>1.00</td>`,
		`<td data-align="center">x</td>`,
		`<td data-align="start" data-numeric>2</td>`,
		`<td>plain</td>`,
	)
}

// TestDetails checks that details are inline rows by default, and that a
// detail without a value shows its children, or a dash without either.
func TestDetails(t *testing.T) {
	html := mustRender(t, `{"components": [{"is": "details", "children": [
		{"is": "detail", "props": {"label": "Name", "value": "Acme"}},
		{"is": "detail", "props": {"label": "State"}, "children": [
			{"is": "tag_status", "props": {"label": "Sent"}}
		]},
		{"is": "detail", "props": {"label": "Folder"}}
	]}]}`)
	assertContains(t, html, `w-full`, `space-y-1`, `>Acme</span>`, `Sent`, `>—</span>`)
	if n := strings.Count(html, `popui-detail-row`); n != 3 {
		t.Errorf("want 3 detail rows, got %d in:\n%s", n, html)
	}
	if strings.Contains(html, "w-[400px]") {
		t.Errorf("inline details have a fixed width in:\n%s", html)
	}

	// Stacked details keep their own layout, with their children as values.
	html = mustRender(t, `{"components": [{"is": "details", "props": {"layout": "stacked"}, "children": [
		{"is": "detail", "props": {"label": "Name", "value": "Acme"}}
	]}]}`)
	if strings.Contains(html, `popui-detail-row`) {
		t.Errorf("stacked details have rows in:\n%s", html)
	}
}

// TestStats checks that the headline figure of a group is marked to lead
// it, and that bound values are shown.
func TestStats(t *testing.T) {
	html := mustRender(t, `{"data": {"totals": {"payable": "€ 121,00", "tax": "€ 21,00"}}, "components": [{"is": "stats", "props": {"align": "end"}, "children": [
		{"is": "stat", "props": {"label": "Total payable", "size": "lg"}, "bind": {"value": "totals.payable"}},
		{"is": "stat", "props": {"label": "Tax"}, "bind": {"value": "totals.tax"}}
	]}]}`)
	assertContains(t, html, `data-size="lg"`, `justify-end`, `Total payable`, `€ 121,00`, `€ 21,00`, `text-xl font-medium`)
	if n := strings.Count(html, `data-size="lg"`); n != 1 {
		t.Errorf("want 1 large stat, got %d in:\n%s", n, html)
	}
}

// TestTextOptions checks the text sizes and tones of paragraphs and titles.
func TestTextOptions(t *testing.T) {
	html := mustRender(t, `{"components": [
		{"is": "title", "props": {"size": "base"}, "children": [{"is": "text", "text": "Supplier"}]},
		{"is": "description", "props": {"size": "sm"}, "children": [{"is": "text", "text": "Note"}]},
		{"is": "p", "props": {"tone": "tertiary", "numeric": true}, "children": [{"is": "text", "text": "12"}]}
	]}`)
	assertContains(t, html, `font-medium`, `text-sm`, `text-foreground-default-tertiary`, `tabular-nums`)
	for _, gone := range []string{"text-lg", "font-semibold"} {
		if strings.Contains(html, gone) {
			t.Errorf("%s not replaced in:\n%s", gone, html)
		}
	}
}
