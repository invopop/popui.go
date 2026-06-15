package props

import (
	"fmt"
	"math/rand"

	"github.com/a-h/templ"
)

// MultiselectOption defines an option for a Multiselect component.
type MultiselectOption struct {
	Value    string
	Label    string
	Selected bool
	Disabled bool
}

// Multiselect Templ component props. The Multiselect renders a select-like
// control where any number of options can be chosen. Selected options are
// shown as removable tag chips inside the field.
type Multiselect struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Name is used for form submission. Each selected value is submitted under
	// this name (e.g. name=a&name=b), matching native <select multiple> behaviour.
	Name string
	// Label renders a simple text label above the field.
	Label string
	// Placeholder is shown inside the field when nothing is selected.
	Placeholder string
	Disabled    bool
	Required    bool
	// Searchable adds a search box inside the dropdown to filter options.
	Searchable bool
	Error      Error
	Options    []MultiselectOption
}

// GenerateID returns a new Multiselect instance with either the existing ID
// or a new randomly generated one. This ensures the label and dropdown can be
// correctly linked when no ID is provided, and is designed to be used inline:
//
//	@popui.Multiselect(props.Multiselect{Name: "tags", Label: "Tags"}.GenerateID())
func (m Multiselect) GenerateID() Multiselect {
	if m.ID != "" {
		return m
	}
	// generate a short random identifier
	m.ID = fmt.Sprintf("multiselect-%06d", rand.Intn(100000))
	return m
}
