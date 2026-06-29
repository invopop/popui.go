package props

import "github.com/a-h/templ"

// TagStatusIcon Templ component props
type TagStatusIcon struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Label      string
	Status     string // success, failed, warning, running
}
