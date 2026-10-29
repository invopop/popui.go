package jsonui

import (
	"strings"
	"testing"
)

func mustRender(t *testing.T, doc string) string {
	t.Helper()
	html, err := render(t, Options{}, doc)
	if err != nil {
		t.Fatal(err)
	}
	return html
}

func assertContains(t *testing.T, html string, wants ...string) {
	t.Helper()
	for _, want := range wants {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in:\n%s", want, html)
		}
	}
}

// TestEachOptions checks that options rendered from a list keep the
// select's model selected as the list or the model changes.
func TestEachOptions(t *testing.T) {
	html := mustRender(t, `{
		"data": {"config": {"connection_id": "c2"}, "connections": [], "loading": false},
		"components": [{
			"is": "select",
			"model": "config.connection_id",
			"bind": {"disabled": "loading"},
			"children": [
				{"is": "option", "props": {"value": "", "label": "Default"}},
				{"is": "option", "each": "connections", "bind": {"value": "item.id", "label": "item.name"}}
			]
		}]
	}`)
	assertContains(t, html,
		`<template x-for="(item, index) in connections">`,
		`:value="item.id"`,
		`x-text="item.name"`,
		`:selected="String(item.id) === String(config.connection_id)"`,
		`:disabled="loading"`,
		`x-model="config.connection_id"`,
	)
	// The saved choice is shown while the list is empty or without it.
	assertContains(t, html,
		`<template x-if="config.connection_id != null &amp;&amp; ![&#34;&#34;,&#34;&#34;].includes(String(config.connection_id)) &amp;&amp; !(connections || []).some((item, index) =&gt; String(item.id) === String(config.connection_id))">`,
		`<option :value="config.connection_id" x-text="String(config.connection_id) + ((connections || []).length ? &#39; (unavailable)&#39; : &#39;&#39;)" :selected="true" selected disabled>`,
		`:selected="String(&#34;&#34;) === String(config.connection_id)"`,
	)
	// Static options are selected at render time.
	if strings.Contains(html, `<option value="" selected`) {
		t.Errorf("default option selected in:\n%s", html)
	}
	html = mustRender(t, `{"data": {"v": "a"}, "components": [{"is": "select", "model": "v", "children": [
		{"is": "option", "props": {"value": "a", "label": "A"}},
		{"is": "option", "props": {"value": "b", "label": "B"}}
	]}]}`)
	assertContains(t, html, `<option value="a" selected`)
	if strings.Contains(html, "x-if") {
		t.Errorf("unavailable option without lists in:\n%s", html)
	}
}

func TestEachText(t *testing.T) {
	html := mustRender(t, `{"data": {"rows": [{"name": "A"}]}, "components": [
		{"is": "p", "each": "rows", "children": [{"is": "text", "bind": "item.name"}, {"is": "text", "text": " #"}, {"is": "text", "bind": "index"}]},
		{"is": "text", "each": "rows", "text": "row"}
	]}`)
	assertContains(t, html,
		`<span x-text="item.name"></span> #<span x-text="index"></span>`,
		`<template x-for="(item, index) in rows"><span>row</span></template>`,
	)
}

func TestShow(t *testing.T) {
	html := mustRender(t, `{"data": {"sandbox": true, "workflow": {"schema": "bill/invoice"}, "n": 2}, "components": [
		{"is": "p", "show": "sandbox", "children": [{"is": "text", "text": "a"}]},
		{"is": "p", "show": "!sandbox", "children": [{"is": "text", "text": "b"}]},
		{"is": "text", "show": {"workflow.schema": ["bill/invoice", "org/party"], "n": [1, 2]}, "text": "c"}
	]}`)
	assertContains(t, html,
		`<template x-if="sandbox">`,
		`<template x-if="!sandbox">`,
		`<template x-if="[1,2].includes(n) &amp;&amp; [&#34;bill/invoice&#34;,&#34;org/party&#34;].includes(workflow.schema)"><span>c</span></template>`,
	)

	// Inside each, conditions can use the item.
	html = mustRender(t, `{"data": {"rows": []}, "components": [
		{"is": "p", "each": "rows", "show": {"item.active": [true]}}
	]}`)
	assertContains(t, html, `<template x-for="(item, index) in rows"><template x-if="[true].includes(item.active)">`)
}

// TestListsStrict checks that lists and conditions only read the document's
// own data.
func TestListsStrict(t *testing.T) {
	bad := map[string]string{
		`{"data": {}, "components": [{"is": "p", "each": "rows"}]}`:                                                                                     `data has no "rows"`,
		`{"data": {"rows": 1}, "components": [{"is": "p", "each": "rows"}]}`:                                                                            `"rows" is not a list`,
		`{"components": [{"is": "p", "each": "rows"}]}`:                                                                                                 `document data must be set`,
		`{"data": {}, "components": [{"is": "p", "show": "document.cookie"}]}`:                                                                          `data has no "document.cookie"`,
		`{"data": {"a": 1}, "components": [{"is": "p", "show": "a || alert(1)"}]}`:                                                                      `is not a data path`,
		`{"data": {"a": 1}, "components": [{"is": "p", "show": {"a": "x"}}]}`:                                                                           `a: expected a list of values`,
		`{"data": {"a": 1}, "components": [{"is": "p", "show": {"a": [{"b": 1}]}}]}`:                                                                    `values must be strings`,
		`{"data": {"a": 1}, "components": [{"is": "p", "show": {}}]}`:                                                                                   `expected at least one data path`,
		`{"data": {"a": 1}, "components": [{"is": "p", "show": 1}]}`:                                                                                    `expected a data path or an object`,
		`{"data": {"a": 1}, "components": [{"is": "p", "show": "item.x"}]}`:                                                                             `data has no "item.x"`,
		`{"data": {"rows": []}, "components": [{"is": "card_header", "each": "rows", "bind": {"title": "item.name"}}]}`:                                 `card_header bind title: can't be bound to the item; no props can`,
		`{"data": {"rows": []}, "components": [{"is": "option", "each": "rows", "bind": {"label": "item.name"}}]}`:                                      `option: must be a child of a select`,
		`{"data": {"v": []}, "components": [{"is": "select", "model": "v", "props": {"multiple": true}, "children": [{"is": "option"}]}]}`:              `children are only supported by single selects without options`,
		`{"data": {"rows": []}, "components": [{"is": "select", "children": [{"is": "option", "each": "rows", "bind": {"value": "item.__proto__"}}]}]}`: errNotDataPath,
	}
	for doc, want := range bad {
		if _, err := render(t, Options{}, doc); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: want error containing %q, got %v", doc, want, err)
		}
	}
}

// TestLiveBind checks that live props follow the data as well as being
// rendered from it.
func TestLiveBind(t *testing.T) {
	html := mustRender(t, `{"data": {"busy": true}, "components": [
		{"is": "button", "bind": {"disabled": "busy"}, "children": [{"is": "text", "text": "Save"}]}
	]}`)
	assertContains(t, html, `:disabled="busy"`, ` disabled`)
}

func TestEachModel(t *testing.T) {
	html := mustRender(t, `{"data": {"to": [{"email": "a@b.c"}]}, "components": [
		{"is": "input", "each": "to", "model": "item.email"}
	]}`)
	assertContains(t, html, `x-model="item.email"`)
}
