package props

import "github.com/a-h/templ"

// DrawerContextItem renders one row inside a DrawerContext / DropdownSelect
// option list. Mirrors @invopop/popui Svelte DrawerContextItem.svelte:
// optional colour-tagged prefix (TagStatus dot-only), label, and a
// right-side affordance — either a reactive checkbox (multi-mode) or a
// tick icon (single-mode, shown only when this row is selected).
//
// Selection state is REACTIVE: the parent owns an Alpine `values` array
// (or single `value`) and passes IsSelectedExpr — an Alpine expression
// evaluating to bool — so the checkbox / tick toggles without a server
// round-trip. OnClickExpr is the Alpine handler invoked on row click.
type DrawerContextItem struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	Value string
	Label string

	// Color is a TagStatus.Status value (green, orange, blue, etc.). When
	// non-empty, a TagStatus dot-only pill is rendered to the left of the
	// label, giving the row the same colour vocabulary as the matching
	// TagStatus rendered elsewhere on the page.
	Color string

	// Multiple switches the right-side affordance from a tick (single
	// mode, shown only when selected) to a checkbox (multi mode, always
	// rendered, reactive to IsSelectedExpr).
	Multiple bool

	// IsSelectedExpr is an Alpine expression evaluating to bool — used to
	// drive the right-side checkbox/tick visibility. Example:
	//   "values.includes('production')"
	IsSelectedExpr string

	// OnClickExpr is the Alpine handler run when the row is clicked.
	// Example: "toggle('production')"
	OnClickExpr string
}
