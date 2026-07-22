package props

import "github.com/a-h/templ"

// StatusBadge Templ component props
type StatusBadge struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Label      string
	Status     string // success, failed, warning, running
}
