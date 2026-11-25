package props

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/a-h/templ"
)

// Popover Templ component props for a centered modal popover
type Popover struct {
	ID         string
	Class      string
	Backdrop   bool // Enable greyed background backdrop
	Attributes templ.Attributes
}

// GenerateID generates a unique ID for the Popover if none is provided.
func (p Popover) GenerateID() Popover {
	if p.ID != "" {
		return p
	}
	// generate a short random identifier
	p.ID = fmt.Sprintf("popover-%06d", rand.Intn(100000))
	return p
}

// ShowPopover returns a JavaScript expression that opens the popover with the
// given ID. The ID is written as a JavaScript string, so no ID can add script
// of its own.
func ShowPopover(id string) string {
	return "document.getElementById(" + jsString(id) + ").showPopover()"
}

// jsString writes s as a JavaScript string literal. JSON's encoding is one,
// and escapes <, > and & too.
func jsString(s string) string {
	b, _ := json.Marshal(s) // a string always encodes
	return string(b)
}
