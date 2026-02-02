package props

import "github.com/a-h/templ"

// Head Templ component properties.
type Head struct {
	Title       string
	Description string

	// AlpineJS when true includes Alpine.js in the head.
	AlpineJS bool

	// Auth when true enables authentication token handling in included scripts.
	Auth bool

	// HTMX when true loads the htmx library including special authentication
	// token handling.
	HTMX bool
	// Axios when true loads the axios javascript library alongside an interceptor
	// to add authentication tokens to requests automatically.
	Axios bool

	// Scripts is a list of additional script paths to include.
	Scripts []Script
	// Stylesheets is a list of additional stylesheet links to include.
	Stylesheets []Link
}

// Body Templ component props.
type Body struct {
	ID    string
	Class string

	// Data adds the x-data attribute to the body element.
	Data string

	Attributes templ.Attributes
}

// PopupLayout Templ component props
type PopupLayout struct {
	Title string
}

// Script defines a script to include in the head section.
type Script struct {
	Src   string
	Defer bool
	Async bool
}

// Link defines a link to include in the head section.
type Link struct {
	Href string
	Rel  string // default rel will be "stylesheet"
}

// RelOrStylesheet returns the rel or "stylesheet" value if not set.
func (l Link) RelOrStylesheet() string {
	if l.Rel == "" {
		return "stylesheet"
	}
	return l.Rel
}
