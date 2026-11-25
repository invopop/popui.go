package popui

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/invopop/popui.go/props"
)

// TestHandlersSurviveSanitization guards the on* handlers popui builds
// itself. templ replaces an on* handler given as a plain string with a
// comment, which leaves the control doing nothing.
func TestHandlersSurviveSanitization(t *testing.T) {
	tests := []struct {
		name string
		c    templ.Component
		want string
	}{
		{"file upload", FileUpload(props.InputFile{ID: "avatar"}), "document.getElementById(this.dataset.popuiFile).click()"},
		{"file download preview", FileDownload(props.FileDownload{Preview: props.ShowPopover("preview")}), "if (!event.target.closest("},
		{"signature pad", SignaturePad(props.SignaturePad{}), "popuiSignPad.clear()"},
		{"signature pad use", SignaturePad(props.SignaturePad{}), "popuiSignPad.use()"},
		{"detail copy", detailActionButton(props.DetailAction{Copy: "secret-value"}, "sm"), "navigator.clipboard.writeText(this.getAttribute('data-popui-copy'))"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := tt.c.Render(context.Background(), &buf); err != nil {
				t.Fatal(err)
			}
			out := buf.String()
			if strings.Contains(out, templ.FailedSanitizationJS) {
				t.Errorf("a handler was sanitized away:\n%s", out)
			}
			if !strings.Contains(out, templ.EscapeString(tt.want)) {
				t.Errorf("handler %q not rendered:\n%s", tt.want, out)
			}
		})
	}
}
