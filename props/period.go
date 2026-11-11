package props

import (
	"time"

	"github.com/a-h/templ"
)

// Period renders a picker for a window of time to show data over, as in a
// list's top bar: a rolling number of days up to today, every time, or
// calendar days, weeks, months, quarters or years, alone or a run of them
// picked by dragging across them or with a shift-click. Arrows beside it
// step to the window after or before. Windows run newest first, as lists
// do: the latest at the top left, and the left arrow steps later, up to
// now.
//
// The value is a string the period package reads, as "365d", "2026-Q4"
// or "2026-01..2026-03", submitted under Name. Picking a period submits
// the form the picker is in or, outside a form, fires a change event from
// its input for hx-trigger="change".
type Period struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	// Name is the field the value is submitted under. Defaults to
	// "period".
	Name string

	// Value is the period picked, as the period package writes it. Empty
	// uses Default.
	Value string

	// Default is the period shown when none is picked, which isn't
	// submitted so that the address stays short. Defaults to
	// period.Default, the last 365 days. A server reading the value
	// should use the same default when none is submitted, as
	// period.Parse("") gives period.Default.
	Default string

	// Today is the current time in the time zone the data is read in,
	// which decides the days rolling periods cover and the windows still
	// to come, which can't be picked. Zero uses the server's clock.
	Today time.Time

	// Rolling lists the numbers of days offered as "Last N days". Empty
	// offers 7, 30, 90 and 365.
	Rolling []int

	// NoAll leaves out "All time".
	NoAll bool

	// Disabled greys the picker out, as while another mode sets the
	// period, such as a list following new rows. The value is still
	// submitted.
	Disabled bool
}
