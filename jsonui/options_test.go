package jsonui

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

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
			`{"template":[{"type":"button","props":{"attributes":{"x-on:click":"alert(1)"}}}]}`,
			`attribute "x-on:click" is not allowed`,
		},
		"alpine shorthand": {
			`{"template":[{"type":"button","props":{"attributes":{"@click":"alert(1)"}}}]}`,
			`attribute "@click" is not allowed`,
		},
		"alpine bind": {
			`{"template":[{"type":"button","props":{"attributes":{":class":"alert(1)"}}}]}`,
			`attribute ":class" is not allowed`,
		},
		"event handler, any case": {
			`{"template":[{"type":"button","props":{"Attributes":{"OnClick":"alert(1)"}}}]}`,
			`attribute "OnClick" is not allowed`,
		},
		"javascript href": {
			`{"template":[{"type":"button","props":{"href":" JavaScript:alert(1)"}}]}`,
			`unsafe URL`,
		},
		"javascript URL nested in actions": {
			`{"template":[{"type":"description_list_item","props":{"actions":[{"url":"javascript:alert(1)"}]}}]}`,
			`description_list_item props.actions[0].url: unsafe URL`,
		},
		"javascript URL in attributes": {
			`{"template":[{"type":"form","props":{"attributes":{"formaction":"javascript:alert(1)"}}}]}`,
			`unsafe URL`,
		},
		"text bind expression": {
			`{"data":{"a":1},"template":[{"type":"text","bind":"alert(1)"}]}`,
			`must be a data path in strict mode`,
		},
		"action not allowed": {
			`{"template":[{"type":"button","on":{"click":{"action":"steal"}}}]}`,
			`action "steal" is not allowed`,
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
		`{"data":{"a":1},"template":[{"type":"input","model":"a; alert(1)"}]}`:         `not a data path`,
		`{"template":[{"type":"button","on":{"click=\"x\" y":{"action":"alert"}}}]}`:   `invalid event name`,
		`{"data":{"a":1},"template":[{"type":"button","model":"a"}]}`:                  `only supported on form controls`,
		`{"template":[{"type":"button","on":{"click":{"message":"no action name"}}}]}`: `needs an "action" name`,
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
		"template": [
			{"type": "input", "props": {"attributes": {"aria-label": "Name", "data-x": "1"}}, "model": "customer.name"},
			{"type": "button", "props": {"href": "https://invopop.com"}, "on": {"click.prevent": [{"action": "alert", "message": "Hi {{customer.name}}"}, {"action": "ui.notify"}]}},
			{"type": "text", "bind": "customer.name"}
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
		doc := `{"template":[{"type":"button","on":{"click":{"action":"` + a + `"}}}]}`
		if _, err := render(t, Options{}, doc); err != nil {
			t.Errorf("%s: %v", a, err)
		}
	}
	doc := `{"template":[{"type":"button","on":{"click":{"action":"alert"}}}]}`
	if _, err := render(t, Options{Actions: []string{}}, doc); err == nil {
		t.Error("empty Actions should allow no actions")
	}
}

func TestTrustedAllowsScript(t *testing.T) {
	doc := `{"data":{"a":1},"template":[
		{"type":"button","props":{"attributes":{"@click":"a++"}},"on":{"click":{"action":"anything"}}},
		{"type":"text","bind":"a * 2"}
	]}`
	html, err := render(t, Options{Trusted: true}, doc)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, `@click="a++"`) || !strings.Contains(html, `x-text="a * 2"`) {
		t.Errorf("trusted output missing script:\n%s", html)
	}
}
