package props

import "github.com/a-h/templ"

// ChartValue is one value in a chart: a column of a ColumnChart, a bar of
// a BarList, or a part of a SegmentedBar.
type ChartValue struct {
	// Label names the value on the chart: a column's axis label ("Nov"),
	// a bar's or part's name.
	Label string
	// Title names it in full in its tooltip, as "November 2026". Defaults
	// to Label.
	Title string
	// Value is the number plotted, from zero.
	Value float64
	// Display is the value as shown, as "1,284" or "12.5%". Defaults to the
	// value with thousands separators.
	Display string
	// Href makes the value a link, such as to the rows it counts.
	Href templ.SafeURL
	// Color is one of the TagStatus colors, for a value that stands for a
	// status. Defaults to the accent color.
	Color string `jsonschema:"enum=,enum=grey,enum=green,enum=yellow,enum=red,enum=orange,enum=blue,enum=purple,enum=olive,enum=teal,enum=crimson,enum=steelBlue,enum=empty"`
	// Partial marks a value still adding up, such as the current month's,
	// drawn lighter.
	Partial bool
}

// ColumnChart shows values over time as columns from a common baseline,
// with a value axis of round numbers, and each value's title and figure
// in a tooltip on hover or focus. Columns with an Href are links.
type ColumnChart struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Label names what the chart shows, for screen readers, which read the
	// values from a table.
	Label string
	// Columns are the values, oldest first.
	Columns []ChartValue
	// Height of the plot in pixels. Defaults to 160.
	Height int
}

// BarList shows values to compare, such as counts by category, as
// horizontal bars under their labels, with figures at the end.
type BarList struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Label names what the list shows, for screen readers.
	Label string
	Bars  []ChartValue
}

// SegmentedBar shows a whole split into its parts, such as documents by
// state, as one bar with a legend of the parts' names and figures.
type SegmentedBar struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Label names the whole, for screen readers.
	Label string
	// Parts in the order they're drawn. Those of no value are left out.
	Parts []ChartValue
	// NoLegend leaves out the legend, as when the parts are listed
	// elsewhere.
	NoLegend bool
}

// Sparkline shows the shape of a series of values in a line, with no
// axes, such as a metric's recent trend.
type Sparkline struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Label describes the trend, for screen readers, which otherwise skip
	// it.
	Label  string
	Values []float64
}

// Metric directions, whether a rise is good news or bad.
const (
	MetricUpGood = ""
	MetricUpBad  = "bad"
)

// Metric shows a headline figure with its label, and optionally its change
// on a previous period and a sparkline of its trend.
type Metric struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Label      string
	// Value is the figure, as shown.
	Value string
	// Href makes the metric a link, such as to the rows it counts.
	Href templ.SafeURL
	// Note follows the value in smaller text, as "4.2% of jobs".
	Note string
	// Change is the relative change on the previous period, 0.12 for 12%
	// more. It's shown when Previous names that period.
	Change float64
	// Previous names the period Change compares with, as "previous 30
	// days" or "Q3 2026", after "vs".
	Previous string
	// NoChange shows that there's nothing to compare with, as when the
	// previous period had none, instead of a change.
	NoChange bool
	// Direction is "bad" when a rise is bad news, such as of errors, which
	// colors the change. Default is a rise being good.
	Direction string `jsonschema:"enum=,enum=bad"`
	// Trend is drawn as a sparkline below the figure.
	Trend []float64
}
