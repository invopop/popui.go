package props

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/a-h/templ"
)

// DropdownSelect mirrors @invopop/popui Svelte DropdownSelect.svelte: an
// input-style trigger that opens a popover panel listing the available
// options. Single or multi selection, with optional per-option Color
// rendered as a TagStatus dot.
//
// In single mode (default) the trigger shows the selected option's dot
// next to its label. In multi mode the trigger shows a stack of dots
// (up to MaxStackedDots) for the selected options, then "{N} {MultipleLabel}".
//
// Form integration: the component emits one or more hidden `<input
// name=Name>` elements driven by the live Alpine state, so the form
// submits the selection naturally. Multi-mode submission is debounced to
// popover-close (only if the selection actually changed) so users can
// tick several boxes without each click hitting the server.
type DropdownSelect struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	Name string

	// Value is the currently-selected option in single mode.
	Value string

	// Values are the currently-selected options in multi mode.
	Values []string

	Multiple      bool
	MultipleLabel string

	Placeholder string

	Options []DropdownSelectOption

	// AutoOpen pops the dropdown open on mount.
	AutoOpen bool
}

// DropdownSelectOption is one row in the dropdown. Color, if set, matches
// the TagStatus.Status vocabulary (green, orange, blue, etc.) so the
// option renders with the same coloured tag as the corresponding
// TagStatus elsewhere on the page.
type DropdownSelectOption struct {
	Value string
	Label string
	Color string
}

// GenerateID generates a unique ID for the DropdownSelect if none is set.
func (d DropdownSelect) GenerateID() DropdownSelect {
	if d.ID != "" {
		return d
	}
	d.ID = fmt.Sprintf("dropdown-select-%06d", rand.Intn(1000000))
	return d
}

// InitialValues returns the starting selection as a slice — Values when
// Multiple, single-element slice from Value otherwise (or empty when
// Value is unset).
func (d DropdownSelect) InitialValues() []string {
	if d.Multiple {
		return d.Values
	}
	if d.Value == "" {
		return nil
	}
	return []string{d.Value}
}

// InitialValuesJSON returns the InitialValues slice as a JSON literal,
// suitable for embedding directly inside an x-data Alpine expression.
func (d DropdownSelect) InitialValuesJSON() string {
	vs := d.InitialValues()
	if vs == nil {
		vs = []string{}
	}
	b, _ := json.Marshal(vs)
	return string(b)
}

// MultipleLabelOr returns MultipleLabel or a sensible default.
func (d DropdownSelect) MultipleLabelOr() string {
	if d.MultipleLabel != "" {
		return d.MultipleLabel
	}
	return "items"
}

// SelectedOption returns the option matching Value (single mode), or nil.
func (d DropdownSelect) SelectedOption() *DropdownSelectOption {
	for i := range d.Options {
		if d.Options[i].Value == d.Value {
			return &d.Options[i]
		}
	}
	return nil
}

// MaxStackedDots caps how many overlapping dots are rendered on the
// trigger when many options are selected. Matches the Svelte
// DropdownSelect (`selectedColors.slice(0, 3)`).
const MaxStackedDots = 3
