package jsonui

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

const patchPanel = `{
  "components": [
    {"is": "card", "key": "summary", "children": [
      {"is": "text", "key": "code", "text": "Code"},
      {"is": "text", "text": "Plain"},
      {"is": "text", "key": "terms", "text": "Terms"},
      {"is": "text", "key": "tags", "text": "Tags"}
    ]},
    {"is": "card", "key": "totals", "children": [{"is": "text", "text": "Totals"}]},
    {"is": "card", "key": "customer", "children": [{"is": "text", "text": "Customer"}]}
  ]
}`

func patchTestPanel(t *testing.T, src string) *Panel {
	t.Helper()
	p := new(Panel)
	if err := decodeStrict([]byte(src), p); err != nil {
		t.Fatal(err)
	}
	return p
}

// texts lists the text of the panel's text nodes, in order.
func texts(p *Panel) string {
	var out []string
	walk(p.Components, func(n *Node) {
		if n.Component == textType {
			out = append(out, n.Text)
		}
	})
	return strings.Join(out, " ")
}

func TestPatchApply(t *testing.T) {
	p := patchTestPanel(t, patchPanel)
	before, _ := json.Marshal(p)

	patch := new(Patch)
	src := `{"hide": ["terms", "gone"], "order": ["customer", "tags", "totals", "code", "gone", "lost"]}`
	if err := decodeStrict([]byte(src), patch); err != nil {
		t.Fatal(err)
	}
	out, missing, err := p.Apply(patch)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := texts(out), "Tags Plain Code Customer Totals"; got != want {
		t.Errorf("patched panel: got %q, want %q", got, want)
	}
	if got := strings.Join(missing, " "); got != "gone lost" {
		t.Errorf("missing keys: got %q", got)
	}
	if after, _ := json.Marshal(p); !bytes.Equal(before, after) {
		t.Errorf("the panel itself changed:\n%s", after)
	}

	// A patched panel renders as any other, and keys aren't rendered.
	c, err := out.Component(Options{})
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := c.Render(context.Background(), &buf); err != nil {
		t.Fatal(err)
	}
	html := buf.String()
	if strings.Contains(html, "Terms") || strings.Contains(html, "summary") {
		t.Errorf("hidden node or key rendered:\n%s", html)
	}
	if strings.Index(html, "Customer") > strings.Index(html, "Totals") {
		t.Errorf("order not applied:\n%s", html)
	}
}

func TestPatchEmpty(t *testing.T) {
	p := patchTestPanel(t, patchPanel)
	for _, patch := range []*Patch{nil, {}} {
		out, missing, err := p.Apply(patch)
		if err != nil || len(missing) > 0 {
			t.Fatal(err, missing)
		}
		if texts(out) != texts(p) {
			t.Errorf("an empty patch changed the panel: %q", texts(out))
		}
	}
}

func TestPanelKeys(t *testing.T) {
	keys, err := patchTestPanel(t, patchPanel).Keys()
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(keys, " "); got != "summary code terms tags totals customer" {
		t.Errorf("keys: got %q", got)
	}

	for name, src := range map[string]string{
		"already used": `{"components": [{"is": "card", "key": "a", "children": [{"is": "text", "key": "a", "text": "x"}]}]}`,
		"must be":      `{"components": [{"is": "text", "key": "1 a", "text": "x"}]}`,
	} {
		p := patchTestPanel(t, src)
		if _, err := p.Keys(); err == nil || !strings.Contains(err.Error(), name) {
			t.Errorf("keys: want %q error, got %v", name, err)
		}
		if _, _, err := p.Apply(nil); err == nil {
			t.Errorf("apply: want %q error", name)
		}
		if _, err := Build([]byte(src)); err == nil || !strings.Contains(err.Error(), name) {
			t.Errorf("build: want %q error, got %v", name, err)
		}
	}
}
