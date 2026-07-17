package props

import (
	"fmt"
	"math/rand"

	"github.com/a-h/templ"
)

// SelectOption defines an option for a Select component
type SelectOption struct {
	Value string
	Label string
	// Description renders next to the label in a dimmed color, giving extra
	// context for the option. Only shown by the Multiple variant's dropdown.
	Description string
	Selected    bool
	Disabled    bool
}

// Select Templ component props.
//
// By default the component renders a native single-choice dropdown. With
// Multiple it becomes a select-like combobox where any number of options can
// be chosen; selected options show as removable tag chips inside the field,
// and each value is submitted under Name (e.g. name=a&name=b) via a visually
// hidden native <select multiple>.
type Select struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Name       string
	Label      string
	Disabled   bool
	Autofocus  bool
	// Multiple switches to the tag-chip combobox variant.
	Multiple bool
	Required bool
	// Placeholder is shown inside the field when nothing is selected
	// (Multiple variant only).
	Placeholder string
	// Searchable adds a search box inside the dropdown to filter options
	// (Multiple variant only).
	Searchable bool
	Error      Error
	Options    []SelectOption
}

// GenerateID returns a new Select instance with either the existing ID or a
// new randomly generated one. The Multiple variant applies it to the combobox
// trigger so the label links to it when no ID is provided. It is designed to
// be used inline:
//
//	@popui.Select(props.Select{Multiple: true, Name: "tags"}.GenerateID())
func (s Select) GenerateID() Select {
	if s.ID != "" {
		return s
	}
	// generate a short random identifier
	s.ID = fmt.Sprintf("select-%06d", rand.Intn(100000))
	return s
}
