package props

import "github.com/a-h/templ"

// Toast defines the props for the Toast component.
type Toast struct {
	ID          string
	Class       string
	Attributes  templ.Attributes
	Variant     string // default, success, error, warning, info
	Title       string
	Description string
	Dismissible bool
	Icon        bool
	Duration    int  // milliseconds, 0 for no auto-dismiss
	Animated    bool // if true, starts hidden and animates in (for HTMX injection)
}

// ToastContainer defines the props for the ToastContainer component.
type ToastContainer struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Position   string // top-right, top-left, top-center, bottom-right, bottom-left, bottom-center
}
