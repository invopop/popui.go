package popui

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/invopop/popui.go/props"
)

const testCopyID = "01a0fd49-5b16-7ecc-9560-15f78de8f0ef"

func TestCopyText(t *testing.T) {
	tests := []struct {
		value          string
		prefix, suffix int
		want           string
	}{
		{testCopyID, 8, 4, "01a0fd49...f0ef"},
		{testCopyID, 8, 0, "01a0fd49..."},
		{testCopyID, 0, 4, "...f0ef"},
		{testCopyID, 0, 0, testCopyID},
		{"short", 3, 2, "short"},
		{"shorter", 3, 2, "sho...er"},
		{"ñandú-señor", 2, 2, "ña...or"},
		{"", 8, 4, ""},
	}
	for _, tt := range tests {
		if got := copyText(tt.value, tt.prefix, tt.suffix); got != tt.want {
			t.Errorf("copyText(%q, %d, %d) = %q, want %q", tt.value, tt.prefix, tt.suffix, got, tt.want)
		}
	}
}

// The copy text is rendered by the server so it's in place before popui.js
// runs, rather than appearing a moment later and shifting the layout.
func TestButtonCopyRendersText(t *testing.T) {
	var buf bytes.Buffer
	c := Button(props.Button{Copy: testCopyID, CopyPrefixLength: 8, CopySuffixLength: 4})
	if err := c.Render(context.Background(), &buf); err != nil {
		t.Fatal(err)
	}
	want := `<span data-copy-text class="font-mono">01a0fd49...f0ef</span>`
	if !strings.Contains(buf.String(), want) {
		t.Errorf("rendered %s, want it to contain %s", buf.String(), want)
	}
}
