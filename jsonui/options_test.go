package jsonui

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

const errUnsafeURL = "unsafe URL"

func render(t *testing.T, opts Options, doc string) (string, error) {
	t.Helper()
	var buf bytes.Buffer
	err := opts.Render(context.Background(), []byte(doc), &buf)
	return buf.String(), err
}

// TestStrictRejectsScript locks in that, by default, a document has no way
// to run script other than the actions the host allows.
func TestStrictRejectsScript(t *testing.T) {
	cases := map[string]struct {
		doc  string
		want string
	}{
		"alpine directive": {
			`{"components":[{"is":"button","props":{"attributes":{"x-on:click":"alert(1)"}}}]}`,
			`attribute "x-on:click" is not allowed`,
		},
		"tooltip image URL": {
			`{"components":[{"is":"input","props":{"label":"Name","tooltip":{"image":"javascript:alert(1)"}}}]}`,
			`input props.tooltip.image: unsafe URL`,
		},
		"description list action URL": {
			`{"components":[{"is":"description_list_item","props":{"label":"Site","actions":[{"url":"javascript:alert(1)"}]}}]}`,
			`description_list_item props.actions[0].url: unsafe URL`,
		},
		"text bind to a global": {
			`{"data":{"a":1},"components":[{"is":"text","bind":"document.cookie"}]}`,
			`text bind: data has no "document.cookie"`,
		},
		"text bind without data": {
			`{"components":[{"is":"text","bind":"document.cookie"}]}`,
			`text bind: document data must be set`,
		},
		"text bind through a prototype": {
			`{"data":{"a":{}},"components":[{"is":"text","bind":"a.constructor.name"}]}`,
			`must be a data path in strict mode`,
		},
		"alpine shorthand": {
			`{"components":[{"is":"button","props":{"attributes":{"@click":"alert(1)"}}}]}`,
			`attribute "@click" is not allowed`,
		},
		"alpine bind": {
			`{"components":[{"is":"button","props":{"attributes":{":class":"alert(1)"}}}]}`,
			`attribute ":class" is not allowed`,
		},
		"event handler, any case": {
			`{"components":[{"is":"button","props":{"Attributes":{"OnClick":"alert(1)"}}}]}`,
			`attribute "OnClick" is not allowed`,
		},
		"javascript href": {
			`{"components":[{"is":"button","props":{"href":" JavaScript:alert(1)"}}]}`,
			errUnsafeURL,
		},
		"javascript URL nested in actions": {
			`{"components":[{"is":"description_list_item","props":{"actions":[{"url":"javascript:alert(1)"}]}}]}`,
			`description_list_item props.actions[0].url: unsafe URL`,
		},
		"javascript URL in attributes": {
			`{"components":[{"is":"form","props":{"attributes":{"formaction":"javascript:alert(1)"}}}]}`,
			errUnsafeURL,
		},
		"text bind expression": {
			`{"data":{"a":1},"components":[{"is":"text","bind":"alert(1)"}]}`,
			`must be a data path in strict mode`,
		},
		"action not allowed": {
			`{"components":[{"is":"button","on":{"click":{"action":"steal"}}}]}`,
			`action "steal" is not allowed`,
		},
		"htmx event handler": {
			`{"components":[{"is":"button","props":{"attributes":{"hx-on:click":"alert(1)"}}}]}`,
			`attribute "hx-on:click" is not allowed`,
		},
		"htmx script values": {
			`{"components":[{"is":"button","props":{"attributes":{"data-hx-vals":"js:{a: alert(1)}"}}}]}`,
			`attribute "data-hx-vals" is not allowed`,
		},
		"bound javascript href": {
			`{"data":{"url":"javascript:alert(1)"},"components":[{"is":"button","bind":{"href":"url"}}]}`,
			errUnsafeURL,
		},
		"bound attributes": {
			`{"data":{"attrs":{"x-on:click":"alert(1)"}},"components":[{"is":"button","bind":{"attributes":"attrs"}}]}`,
			`attribute "x-on:click" is not allowed`,
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := render(t, Options{}, c.doc)
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("want error containing %q, got %v", c.want, err)
			}
		})
	}
}

// TestAlwaysRejected covers checks that apply to trusted documents too, as
// "model" and "on" keys end up inside generated Alpine attributes.
func TestAlwaysRejected(t *testing.T) {
	cases := map[string]string{
		`{"data":{"a":1},"components":[{"is":"input","model":"a; alert(1)"}]}`:                              `not a data path`,
		`{"components":[{"is":"button","on":{"click=\"x\" y":{"action":"alert"}}}]}`:                        `invalid event name`,
		`{"data":{"a":1},"components":[{"is":"button","model":"a"}]}`:                                       `only supported on form controls`,
		`{"components":[{"is":"button","on":{"click":{"message":"no action name"}}}]}`:                      `needs an "action" name`,
		`{"data":{"a":{}},"components":[{"is":"input","model":"a.__proto__.x"}]}`:                           `not a data path`,
		`{"components":[{"is":"button","on":{"click":{"action":"set","path":"constructor.prototype.x"}}}]}`: `path "constructor.prototype.x" is not a data path`,
		`{"components":[{"is":"button","on":{"click":{"action":"alert","into":"__proto__.x"}}}]}`:           `into "__proto__.x" is not a data path`,
		`{"components":[{"is":"button","on":{"click":{"action":"alert","busy":1}}}]}`:                       `busy 1 is not a data path`,
		`{"components":[{"is":"button","on":{"click":null}}]}`:                                              `expected an action`,
		`{"handlers":{"h":null},"components":[]}`:                                                           `expected an action`,
		`{"components":[]} {"components":[]}`:                                                               `unexpected data after the JSON value`,
	}
	for doc, want := range cases {
		_, err := render(t, Options{Trusted: true}, doc)
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: want error containing %q, got %v", doc, want, err)
		}
	}
}

func TestStrictAllows(t *testing.T) {
	doc := `{
		"data": {"customer": {"name": "Acme", "site": "https://acme.example"}},
		"components": [
			{"is": "input", "props": {"attributes": {"aria-label": "Name", "data-x": "1"}}, "model": "customer.name"},
			{"is": "button", "props": {"href": "https://invopop.com"}, "on": {"click.prevent": [{"action": "alert", "message": "Hi {{customer.name}}"}, {"action": "ui.notify"}]}},
			{"is": "text", "bind": "customer.name"}
		]
	}`
	html, err := render(t, Options{Actions: append(BuiltinActions, "ui.notify")}, doc)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`x-model="customer.name"`,
		`value="Acme"`,
		`aria-label="Name"`,
		`x-on:click.prevent="$actions([`,
		`<span x-text="customer.name">Acme</span>`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in:\n%s", want, html)
		}
	}
}

func TestDefaultActionsAreBuiltins(t *testing.T) {
	for _, a := range BuiltinActions {
		doc := `{"components":[{"is":"button","on":{"click":{"action":"` + a + `"}}}]}`
		if _, err := render(t, Options{}, doc); err != nil {
			t.Errorf("%s: %v", a, err)
		}
	}
	doc := `{"components":[{"is":"button","on":{"click":{"action":"alert"}}}]}`
	if _, err := render(t, Options{Actions: []string{}}, doc); err == nil {
		t.Error("empty Actions should allow no actions")
	}
}

func TestTrustedAllowsScript(t *testing.T) {
	doc := `{"data":{"a":1},"components":[
		{"is":"button","props":{"attributes":{"@click":"a++"}},"on":{"click":{"action":"anything"}}},
		{"is":"text","bind":"a * 2"}
	]}`
	html, err := render(t, Options{Trusted: true}, doc)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, `@click="a++"`) || !strings.Contains(html, `x-text="a * 2"`) {
		t.Errorf("trusted output missing script:\n%s", html)
	}
}

func TestTabsModelMustBeIdentifier(t *testing.T) {
	doc := `{"components":[{"is":"tabs","props":{"model":"x = alert(1), y","tabs":[{"value":"a","label":"A"}]}}]}`
	if _, err := render(t, Options{}, doc); err == nil || !strings.Contains(err.Error(), "must be a plain identifier") {
		t.Fatalf("want identifier error, got %v", err)
	}
	doc = `{"components":[{"is":"tabs","props":{"model":"view","tabs":[{"value":"a","label":"A"}]},"children":[{"is":"tab","props":{"model":"view","value":"a"}}]}]}`
	if _, err := render(t, Options{}, doc); err != nil {
		t.Fatal(err)
	}
}

func TestUploadFile(t *testing.T) {
	doc := `{
		"data": {"config": {"logo_url": ""}, "ui": {"uploading": false}},
		"components": [{
			"is": "upload_file",
			"model": "config.logo_url",
			"props": {"loading": "ui.uploading"},
			"on": {
				"select": {"action": "asset.choose", "into": "config.logo_url", "busy": "ui.uploading"},
				"remove": {"action": "set", "path": "config.logo_url", "value": ""}
			}
		}]
	}`
	html, err := render(t, Options{Actions: append(BuiltinActions, "asset.choose")}, doc)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`:src="config.logo_url"`,
		`x-text="ui.uploading ? &#39;Uploading...&#39; : &#39;Select logo&#39;"`,
		`@click="$actions([{&#34;action&#34;:&#34;asset.choose&#34;`,
		`:disabled="ui.uploading || !config.logo_url"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in:\n%s", want, html)
		}
	}

	bad := map[string]string{
		`{"data":{"a":""},"components":[{"is":"upload_file","model":"a","props":{"loading":"x; alert(1)"}}]}`:  `is not a data path`,
		`{"data":{"a":""},"components":[{"is":"upload_file","model":"a","on":{"click":{"action":"alert"}}}]}`:  `unknown event "click"`,
		`{"data":{"a":""},"components":[{"is":"upload_file","model":"a","on":{"select":{"action":"steal"}}}]}`: `action "steal" is not allowed`,
		`{"data":{"a":""},"components":[{"is":"upload_file"}]}`:                                                `model is required`,
	}
	for doc, want := range bad {
		if _, err := render(t, Options{}, doc); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: want error containing %q, got %v", doc, want, err)
		}
	}
}

func TestHandlers(t *testing.T) {
	doc := `{
		"data": {"a": ""},
		"handlers": {"update": {"action": "dispatch", "event": "changed"}},
		"components": [
			{"is": "input", "model": "a", "on": {"change": "update", "blur": [{"action": "set", "path": "a", "value": ""}, "update"]}}
		]
	}`
	html, err := render(t, Options{}, doc)
	if err != nil {
		t.Fatal(err)
	}
	update := `{&#34;action&#34;:&#34;dispatch&#34;,&#34;event&#34;:&#34;changed&#34;}`
	for _, want := range []string{
		`x-on:change="$actions([` + update + `], $event)"`,
		`x-on:blur="$actions([{&#34;action&#34;:&#34;set&#34;,&#34;path&#34;:&#34;a&#34;,&#34;value&#34;:&#34;&#34;},` + update + `], $event)"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s in:\n%s", want, html)
		}
	}

	bad := map[string]string{
		`{"components":[{"is":"button","on":{"click":"missing"}}]}`:                             `unknown handler "missing"`,
		`{"handlers":{"a":{"action":"alert"},"b":["a"]},"components":[]}`:                       `can't refer to other handlers`,
		`{"handlers":{"unused":{"action":"steal"}},"components":[]}`:                            `handler "unused": action "steal" is not allowed`,
		`{"handlers":{"n":42},"components":[]}`:                                                 `expected an action, a handler name or a list`,
		`{"handlers":{"x":{"message":"hi"}},"components":[{"is":"button","on":{"click":"x"}}]}`: `needs an "action" name`,
	}
	for doc, want := range bad {
		if _, err := render(t, Options{}, doc); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: want error containing %q, got %v", doc, want, err)
		}
	}
}

// TestNullProps checks that null props and attributes are treated as empty
// when generated bindings are added to them.
func TestNullProps(t *testing.T) {
	for _, doc := range []string{
		`{"data":{"a":"x"},"components":[{"is":"input","props":null,"model":"a"}]}`,
		`{"components":[{"is":"button","props":{"attributes":null},"on":{"click":{"action":"alert"}}}]}`,
	} {
		html, err := render(t, Options{}, doc)
		if err != nil {
			t.Errorf("%s: %v", doc, err)
		} else if !strings.Contains(html, "x-model") && !strings.Contains(html, "x-on:click") {
			t.Errorf("%s: missing binding in:\n%s", doc, html)
		}
	}
}

// TestBoundAttributesKeepHandlers checks that attributes taken from the data
// are combined with the ones generated for "on".
func TestBoundAttributesKeepHandlers(t *testing.T) {
	doc := `{"data":{"attrs":{"aria-label":"Save"}},"components":[
		{"is":"button","bind":{"attributes":"attrs"},"on":{"click":{"action":"alert"}}}
	]}`
	html, err := render(t, Options{}, doc)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, `aria-label="Save"`) || !strings.Contains(html, `x-on:click=`) {
		t.Errorf("want bound and generated attributes in:\n%s", html)
	}
}

func TestMultipleSelectModel(t *testing.T) {
	doc := `{"data":{"tags":["es","pt"]},"components":[{"is":"select","model":"tags","props":{"multiple":true,"options":[
		{"value":"es","label":"Spanish"},{"value":"fr","label":"French"},{"value":"pt","label":"Portuguese"}
	]}}]}`
	html, err := render(t, Options{}, doc)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range []string{"es", "pt"} {
		if !strings.Contains(html, `value="`+v+`" selected`) {
			t.Errorf("option %s not selected in:\n%s", v, html)
		}
	}
	if strings.Contains(html, `value="fr" selected`) {
		t.Errorf("option fr selected in:\n%s", html)
	}
}
