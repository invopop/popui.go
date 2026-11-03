package popui

import (
	"slices"

	"github.com/invopop/popui.go/classes"
	"github.com/invopop/popui.go/props"
)

// gridBreakpoint holds a grid's classes for the width its columns start
// at. They're written out in full so Tailwind finds them. Arrays of
// counts are indexed by the count, and of borders by style: solid, then
// dashed.
type gridBreakpoint struct {
	cols, span, rows [5]string
	// right and left are side borders while there are columns, and
	// bottom and top are what they become while the cells are stacked.
	right, left, bottom, top [2]string
}

var gridBreakpoints = map[string]*gridBreakpoint{
	props.GridStackSmall: {
		cols:   [5]string{"", "sm:grid-cols-1", "sm:grid-cols-2", "sm:grid-cols-3", "sm:grid-cols-4"},
		span:   [5]string{"", "", "sm:col-span-2", "sm:col-span-3", "sm:col-span-4"},
		rows:   [5]string{"", "", "sm:row-span-2", "sm:row-span-3", "sm:row-span-4"},
		right:  [2]string{"sm:border-r", "sm:border-r sm:[border-right-style:dashed]"},
		left:   [2]string{"sm:border-l", "sm:border-l sm:[border-left-style:dashed]"},
		bottom: [2]string{"max-sm:border-b", "max-sm:border-b max-sm:[border-bottom-style:dashed]"},
		top:    [2]string{"max-sm:border-t", "max-sm:border-t max-sm:[border-top-style:dashed]"},
	},
	props.GridStackMedium: {
		cols:   [5]string{"", "md:grid-cols-1", "md:grid-cols-2", "md:grid-cols-3", "md:grid-cols-4"},
		span:   [5]string{"", "", "md:col-span-2", "md:col-span-3", "md:col-span-4"},
		rows:   [5]string{"", "", "md:row-span-2", "md:row-span-3", "md:row-span-4"},
		right:  [2]string{"md:border-r", "md:border-r md:[border-right-style:dashed]"},
		left:   [2]string{"md:border-l", "md:border-l md:[border-left-style:dashed]"},
		bottom: [2]string{"max-md:border-b", "max-md:border-b max-md:[border-bottom-style:dashed]"},
		top:    [2]string{"max-md:border-t", "max-md:border-t max-md:[border-top-style:dashed]"},
	},
	props.GridStackLarge: {
		cols:   [5]string{"", "lg:grid-cols-1", "lg:grid-cols-2", "lg:grid-cols-3", "lg:grid-cols-4"},
		span:   [5]string{"", "", "lg:col-span-2", "lg:col-span-3", "lg:col-span-4"},
		rows:   [5]string{"", "", "lg:row-span-2", "lg:row-span-3", "lg:row-span-4"},
		right:  [2]string{"lg:border-r", "lg:border-r lg:[border-right-style:dashed]"},
		left:   [2]string{"lg:border-l", "lg:border-l lg:[border-left-style:dashed]"},
		bottom: [2]string{"max-lg:border-b", "max-lg:border-b max-lg:[border-bottom-style:dashed]"},
		top:    [2]string{"max-lg:border-t", "max-lg:border-t max-lg:[border-top-style:dashed]"},
	},
	props.GridStackXLarge: {
		cols:   [5]string{"", "xl:grid-cols-1", "xl:grid-cols-2", "xl:grid-cols-3", "xl:grid-cols-4"},
		span:   [5]string{"", "", "xl:col-span-2", "xl:col-span-3", "xl:col-span-4"},
		rows:   [5]string{"", "", "xl:row-span-2", "xl:row-span-3", "xl:row-span-4"},
		right:  [2]string{"xl:border-r", "xl:border-r xl:[border-right-style:dashed]"},
		left:   [2]string{"xl:border-l", "xl:border-l xl:[border-left-style:dashed]"},
		bottom: [2]string{"max-xl:border-b", "max-xl:border-b max-xl:[border-bottom-style:dashed]"},
		top:    [2]string{"max-xl:border-t", "max-xl:border-t max-xl:[border-top-style:dashed]"},
	},
	// A grid that never stacks has its columns and side borders at
	// every width.
	props.GridStackNone: {
		cols:  [5]string{"", "grid-cols-1", "grid-cols-2", "grid-cols-3", "grid-cols-4"},
		span:  [5]string{"", "", "col-span-2", "col-span-3", "col-span-4"},
		rows:  [5]string{"", "", "row-span-2", "row-span-3", "row-span-4"},
		right: [2]string{"border-r", "border-r [border-right-style:dashed]"}, //nolint:goconst // written out for Tailwind
		left:  [2]string{"border-l", "border-l [border-left-style:dashed]"},
	},
}

// Borders on the top and bottom edges, which stay put when the grid
// stacks.
var (
	gridBorderTop    = [2]string{"border-t", "border-t [border-top-style:dashed]"}
	gridBorderBottom = [2]string{"border-b", "border-b [border-bottom-style:dashed]"}
)

// gridStackKey is the context key for the breakpoint of the grid a cell
// is in.
type gridStackKey struct{}

// gridStack returns the grid's breakpoint, defaulting to "md".
func gridStack(p props.Grid) *gridBreakpoint {
	if bp := gridBreakpoints[p.Stack]; bp != nil {
		return bp
	}
	return gridBreakpoints[props.GridStackMedium]
}

// gridClass returns the classes of the grid's columns.
func gridClass(p props.Grid) string {
	n := p.Columns
	if n < 1 || n > 4 {
		n = 2
	}
	if p.Stack == props.GridStackNone {
		return classes.Join("grid", gridStack(p).cols[n])
	}
	return classes.Join("grid grid-cols-1", gridStack(p).cols[n])
}

// cellClass returns the classes of a cell in a grid with the breakpoint:
// its spans and borders.
func cellClass(bp *gridBreakpoint, p props.GridCell) string {
	list := []string{bp.span[clampCount(p.Span)], bp.rows[clampCount(p.Rows)]}
	has := func(edge string) bool {
		return slices.ContainsFunc(p.Borders, func(b props.Border) bool { return b.Edge == edge })
	}
	for _, b := range p.Borders {
		s := 0
		if b.Style == props.BorderDashed {
			s = 1
		}
		switch b.Edge {
		case props.EdgeTop:
			list = append(list, gridBorderTop[s])
		case props.EdgeBottom:
			list = append(list, gridBorderBottom[s])
		case props.EdgeRight:
			list = append(list, bp.right[s])
			if !has(props.EdgeBottom) {
				list = append(list, bp.bottom[s])
			}
		case props.EdgeLeft:
			list = append(list, bp.left[s])
			if !has(props.EdgeTop) {
				list = append(list, bp.top[s])
			}
		}
	}
	return classes.Join(list...)
}

// clampCount keeps a span within the 1 to 4 a grid has columns for.
func clampCount(n int) int {
	return min(max(n, 1), 4)
}
