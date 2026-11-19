package props

import "github.com/a-h/templ"

// Grid breakpoints, below which a grid's cells stack in one column.
const (
	GridStackSmall  = "sm"
	GridStackMedium = "md"
	GridStackLarge  = "lg"
	GridStackXLarge = "xl"
	// GridStackNone keeps the columns at every width.
	GridStackNone = "none"
)

// GridCell border edges.
const (
	EdgeTop    = "top"
	EdgeRight  = "right"
	EdgeBottom = "bottom"
	EdgeLeft   = "left"
)

// GridCell border styles.
const (
	BorderSolid  = "solid"
	BorderDashed = "dashed"
)

// Grid defines the props for the Grid component.
type Grid struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Columns is the number of columns, from 1 to 4. Default is 2.
	Columns int `jsonschema:"minimum=1,maximum=4"`
	// Stack is the breakpoint below which the cells stack in one column:
	// "sm", "md", "lg" or "xl", or "none" to keep the columns at every
	// width. Default is "md".
	Stack string `jsonschema:"enum=,enum=sm,enum=md,enum=lg,enum=xl,enum=none"`
}

// GridCell defines the props for the Cell component, a Grid's section.
type GridCell struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Span is the number of columns the cell takes, at most the grid's
	// Columns. Default is 1.
	Span int `jsonschema:"minimum=1,maximum=4"`
	// Rows is the number of rows the cell takes, beside cells stacked in
	// the other columns. Default is 1.
	Rows int `jsonschema:"minimum=1,maximum=4"`
	// Borders are the lines on the cell's edges, between it and the cells
	// beside it. While the grid is stacked, a right border is drawn at the
	// bottom and a left one at the top, unless the cell has its own there.
	Borders []Border
}

// Border is a line on one edge of a grid's Cell.
type Border struct {
	// Edge is the side of the cell the line is on.
	Edge string `jsonschema:"enum=top,enum=right,enum=bottom,enum=left"`
	// Style is "solid" or "dashed". Default is "solid".
	Style string `jsonschema:"enum=,enum=solid,enum=dashed"`
}
