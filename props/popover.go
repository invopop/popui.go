package props

import (
	"fmt"
	"math/rand"

	"github.com/a-h/templ"
)

// Popover Templ component props for a centered modal popover
type Popover struct {
	ID         string
	Class      string
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
