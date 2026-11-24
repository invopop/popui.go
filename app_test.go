package popui

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/invopop/popui.go/props"
)

func TestAppOverscroll(t *testing.T) {
	const rule = "html, body { overscroll-behavior: none; }"
	tests := []struct {
		name string
		p    props.App
	}{
		{"default", props.App{}},
		{"without PopUI's stylesheet", props.App{SkipCSS: true}},
		{"filling the viewport", props.App{FillViewport: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := App(tt.p).Render(context.Background(), &buf); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(buf.String(), rule) {
				t.Errorf("the page may bounce: %q isn't in the head", rule)
			}
		})
	}
}
