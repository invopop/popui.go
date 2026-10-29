package props

import "github.com/a-h/templ"

// StatusBadge Templ component props
type StatusBadge struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Label      string
	// Status is the state the badge shows.
	Status string `jsonschema:"enum=success,enum=failed,enum=warning,enum=running"`
}
