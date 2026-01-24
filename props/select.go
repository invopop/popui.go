package props

import "github.com/a-h/templ"

// SelectOption defines an option for a Select component
type SelectOption struct {
	Value    string
	Label    string
	Selected bool
	Disabled bool
}

// Select Templ component props
type Select struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Name       string
	Label      string
	Disabled   bool
	Autofocus  bool
	Multiple   bool
	Required   bool
	Error      Error
	Options    []SelectOption
}
