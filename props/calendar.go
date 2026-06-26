package props

import "github.com/a-h/templ"

// Calendar renders a dual-month, range-selection calendar grid with a
// preset rail (this/last week, month, quarter, custom) and month
// navigation. See popui.Calendar.
//
// Interactivity is provided by the `rangeCalendar` Alpine controller in
// popui.js. There are two ways to supply that scope:
//
//   - Self-contained: set Name (and optionally From / To). The component
//     declares its own `x-data="rangeCalendar({...})"` and renders inside a
//     bordered popover-style container — drop it in anywhere:
//
//     @popui.Calendar(props.Calendar{Name: "period"})
//
//   - Embedded: leave Name empty. The component renders markup only and
//     inherits a `rangeCalendar` scope from an ancestor — this is how
//     popui.Filter's calendar chip hosts it (the chip owns the scope so its
//     summary box and the grid stay in sync).
// Calendar preset keys — each maps to a built-in date-range computation in the
// rangeCalendar controller (popui.js). Use them as CalendarPreset.Key to choose
// which shortcuts the preset rail shows; the date math stays in popui.js.
const (
	CalendarPresetThisWeek    string = "thisWeek"
	CalendarPresetLastWeek    string = "lastWeek"
	CalendarPresetThisMonth   string = "thisMonth"
	CalendarPresetLastMonth   string = "lastMonth"
	CalendarPresetThisQuarter string = "thisQuarter"
	CalendarPresetLastQuarter string = "lastQuarter"
	// CalendarPresetCustom is the "no preset / pick your own dates" entry; keep
	// it in the rail so users can always fall back to a manual range.
	CalendarPresetCustom string = "custom"
)

// CalendarPreset is one shortcut in the calendar's preset rail. Key selects the
// built-in range (a CalendarPreset* constant); Label is the rail text — falls
// back to the Key when empty.
type CalendarPreset struct {
	Key   string
	Label string
}

// DefaultCalendarPresets is the preset rail used when none is supplied:
// this/last week, this/last month, this/last quarter, plus custom. Callers can
// pass a subset (or reordering) to Calendar.Presets / FilterInput.Presets.
func DefaultCalendarPresets() []CalendarPreset {
	return []CalendarPreset{
		{Key: CalendarPresetThisWeek, Label: "This Week"},
		{Key: CalendarPresetLastWeek, Label: "Last Week"},
		{Key: CalendarPresetThisMonth, Label: "This month"},
		{Key: CalendarPresetLastMonth, Label: "Last month"},
		{Key: CalendarPresetThisQuarter, Label: "This quarter"},
		{Key: CalendarPresetLastQuarter, Label: "Last quarter"},
		{Key: CalendarPresetCustom, Label: "Custom"},
	}
}

type Calendar struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	// Name, when set, makes the calendar self-contained: it declares its own
	// `rangeCalendar` Alpine scope under this field name and wraps the grid in
	// a bordered container. Leave empty to render markup only inside an
	// ancestor-provided scope.
	Name string
	// From / To seed the initially selected range (ISO yyyy-mm-dd) when Name
	// is set. Empty means no initial selection.
	From string
	To   string

	// Presets selects which range shortcuts the preset rail shows, and in what
	// order. Empty uses DefaultCalendarPresets(). Each Key must be a
	// CalendarPreset* constant (the date math lives in popui.js).
	Presets []CalendarPreset
}
