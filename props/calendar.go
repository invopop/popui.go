package props

import "github.com/a-h/templ"

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

	// Forward-looking presets, for validity periods and scheduling rather than
	// filtering history. Unlike the calendar-aligned this/last presets, they all
	// roll from today to the same day of the month N months ahead (17 Sep →
	// 17 Oct for a month), clamped to the last day when that month is shorter.
	//
	// CalendarPresetNextMonth is today to the same day next month.
	CalendarPresetNextMonth string = "nextMonth"
	// CalendarPresetNext3Months is today to the same day three months out.
	CalendarPresetNext3Months string = "next3Months"
	// CalendarPresetNext6Months is today to the same day six months out.
	CalendarPresetNext6Months string = "next6Months"
	// CalendarPresetNext12Months is today to the same day twelve months out.
	CalendarPresetNext12Months string = "next12Months"
	// CalendarPresetIndefinite is a start date with no end date. It selects
	// from today onwards (or from a start day already clicked without an end);
	// clicking a day moves the start. The grid paints every
	// day after the start as selected, and the value is submitted as
	// "YYYY-MM-DD.." (empty end). Confirm is enabled as soon as a start is set.
	CalendarPresetIndefinite string = "indefinite"

	// CalendarPresetCustom is the "no preset / pick your own dates" entry; keep
	// it in the rail so users can always fall back to a manual range.
	CalendarPresetCustom string = "custom"
)

// CalendarPreset is one shortcut in the calendar's preset rail. Key selects the
// built-in range (a CalendarPreset* constant); Label is the rail text — falls
// back to CalendarPresetLabel(Key), then to the Key itself, when empty.
type CalendarPreset struct {
	Key   string
	Label string
}

// calendarPresetLabels holds the standard rail text for every built-in key.
var calendarPresetLabels = map[string]string{
	CalendarPresetThisWeek:     "This Week",
	CalendarPresetLastWeek:     "Last Week",
	CalendarPresetThisMonth:    "This month",
	CalendarPresetLastMonth:    "Last month",
	CalendarPresetThisQuarter:  "This quarter",
	CalendarPresetLastQuarter:  "Last quarter",
	CalendarPresetNextMonth:    "Next month",
	CalendarPresetNext3Months:  "Next 3 months",
	CalendarPresetNext6Months:  "Next 6 months",
	CalendarPresetNext12Months: "Next 12 months",
	CalendarPresetIndefinite:   "Indefinite",
	CalendarPresetCustom:       "Custom",
}

// CalendarPresetLabel returns the standard rail text for a built-in preset key,
// or "" for an unknown key.
func CalendarPresetLabel(key string) string {
	return calendarPresetLabels[key]
}

// DefaultCalendarPresets is the preset rail used when none is supplied:
// this/last week, this/last month, this/last quarter, plus custom. Callers can
// pass a subset (or reordering) to Calendar.Presets / FilterInput.Presets.
func DefaultCalendarPresets() []CalendarPreset {
	return calendarPresets{
		{Key: CalendarPresetThisWeek},
		{Key: CalendarPresetLastWeek},
		{Key: CalendarPresetThisMonth},
		{Key: CalendarPresetLastMonth},
		{Key: CalendarPresetThisQuarter},
		{Key: CalendarPresetLastQuarter},
		{Key: CalendarPresetCustom},
	}.withLabels()
}

// FutureCalendarPresets is a forward-looking rail for validity periods and
// scheduling: next month, next 3 / 6 / 12 months, indefinite, plus custom.
func FutureCalendarPresets() []CalendarPreset {
	return calendarPresets{
		{Key: CalendarPresetNextMonth},
		{Key: CalendarPresetNext3Months},
		{Key: CalendarPresetNext6Months},
		{Key: CalendarPresetNext12Months},
		{Key: CalendarPresetIndefinite},
		{Key: CalendarPresetCustom},
	}.withLabels()
}

// calendarPresets is a helper type for filling in standard labels.
type calendarPresets []CalendarPreset

func (ps calendarPresets) withLabels() []CalendarPreset {
	for i := range ps {
		if ps[i].Label == "" {
			ps[i].Label = CalendarPresetLabel(ps[i].Key)
		}
	}
	return ps
}

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
	// is set. Both empty means no initial selection; From without To seeds an
	// indefinite range (start date, no end — see CalendarPresetIndefinite).
	From string
	To   string

	// Presets selects which range shortcuts the preset rail shows, and in what
	// order. Empty uses DefaultCalendarPresets(). Each Key must be a
	// CalendarPreset* constant (the date math lives in popui.js). Ignored when
	// Single is set — single-date pickers have no preset rail.
	Presets []CalendarPreset

	// Single switches the calendar to single-date selection: one month grid,
	// no preset rail, and clicking a day selects just that day. Seed the
	// initial selection with From (To is ignored).
	Single bool
}
