package props

import (
	"strings"
	"testing"
)

func TestShowPopover(t *testing.T) {
	tests := map[string]string{
		"file-preview":    `document.getElementById("file-preview").showPopover()`,
		`x');alert(1);('`: `document.getElementById("x');alert(1);('").showPopover()`,
		`x");alert(1);("`: `document.getElementById("x\");alert(1);(\"").showPopover()`,
	}
	for id, want := range tests {
		if got := ShowPopover(id); got != want {
			t.Errorf("ShowPopover(%q) = %s, want %s", id, got, want)
		}
	}
	// Markup can't close a script element the expression is placed in.
	if got := ShowPopover("</script>"); strings.ContainsAny(got, "<>") {
		t.Errorf("ShowPopover kept markup: %s", got)
	}
}
