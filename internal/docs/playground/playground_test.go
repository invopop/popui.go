package playground

import (
	"strings"
	"testing"

	"github.com/invopop/popui.go/jsonui/schema"
)

// TestSamplesSchema checks the samples refer to the schema by its "$id", which
// the editor replaces with the playground's own panel.json.
func TestSamplesSchema(t *testing.T) {
	data, err := static.ReadFile("samples.js")
	if err != nil {
		t.Fatal(err)
	}
	if want := `const SCHEMA = "` + schema.ID + `";`; !strings.Contains(string(data), want) {
		t.Errorf("samples.js doesn't contain %s", want)
	}
}
