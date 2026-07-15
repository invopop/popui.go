package props

import (
	"github.com/a-h/templ"
)

// FileBadge Templ component props
type FileBadge struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Type is the file extension shown inside the badge, e.g. "xml",
	// "pdf" or "png". Recognized types are colored with their document
	// token; anything else uses the secondary foreground color.
	Type string
}
