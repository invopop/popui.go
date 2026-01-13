package props

import "github.com/a-h/templ"

// PageState defines the props for the PageState component.
type PageState struct {
	ID           string
	Class        string
	Attributes   templ.Attributes
	Illustration templ.Component
	Title        string
	Description  string
}
