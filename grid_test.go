package popui

import (
	"slices"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/invopop/popui.go/props"
)

// TestGrid checks the grid's columns start at its breakpoint, two of them
// by default, and that its cells span them and draw their borders.
func TestGrid(t *testing.T) {
	html := renderWith(t, Grid(), templ.Join(
		renderedCell(props.GridCell{Rows: 2, Borders: []props.Border{{Edge: props.EdgeRight, Style: props.BorderDashed}}}),
		renderedCell(props.GridCell{Borders: []props.Border{{Edge: props.EdgeBottom, Style: props.BorderDashed}}}),
		renderedCell(props.GridCell{Span: 2, Borders: []props.Border{{Edge: props.EdgeTop}}}),
	))
	for _, want := range []string{
		"grid grid-cols-1 md:grid-cols-2",
		"md:row-span-2",
		"md:border-r md:[border-right-style:dashed]",
		"max-md:border-b max-md:[border-bottom-style:dashed]",
		"border-b [border-bottom-style:dashed]",
		"md:col-span-2",
		"border-t",
	} {
		if !hasClasses(html, want) {
			t.Errorf("missing %s in:\n%s", want, html)
		}
	}
}

// TestCellStacked checks that a side border moves to the matching end of
// the cell while the grid is stacked, unless the cell has its own there.
func TestCellStacked(t *testing.T) {
	g := gridLayout{bp: gridBreakpoints[props.GridStackLarge], columns: 2}
	tests := []struct {
		name    string
		borders []props.Border
		want    []string
		not     []string
	}{
		{
			name:    "right",
			borders: []props.Border{{Edge: props.EdgeRight}},
			want:    []string{"lg:border-r max-lg:border-b"},
		},
		{
			name:    "left dashed",
			borders: []props.Border{{Edge: props.EdgeLeft, Style: props.BorderDashed}},
			want:    []string{"lg:[border-left-style:dashed]", "max-lg:border-t max-lg:[border-top-style:dashed]"},
		},
		{
			name:    "right with bottom",
			borders: []props.Border{{Edge: props.EdgeRight, Style: props.BorderDashed}, {Edge: props.EdgeBottom}},
			want:    []string{"lg:border-r lg:[border-right-style:dashed] border-b"},
			not:     []string{"max-lg:border-b", "border-bottom-style"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := cellClass(g, props.GridCell{Borders: tt.borders})
			for _, w := range tt.want {
				if !hasClasses(c, w) {
					t.Errorf("missing %s in %q", w, c)
				}
			}
			for _, n := range tt.not {
				if strings.Contains(c, n) {
					t.Errorf("unexpected %s in %q", n, c)
				}
			}
		})
	}
}

// TestGridStackNone checks a grid that never stacks has its columns and
// side borders at every width.
func TestGridStackNone(t *testing.T) {
	html := renderWith(t, Grid(props.Grid{Columns: 3, Stack: props.GridStackNone}), templ.Join(
		renderedCell(props.GridCell{Borders: []props.Border{{Edge: props.EdgeRight}}}),
	))
	for _, want := range []string{"grid grid-cols-3", "border-r border-border"} {
		if !hasClasses(html, want) {
			t.Errorf("missing %s in:\n%s", want, html)
		}
	}
	if strings.Contains(html, "max-") || strings.Contains(html, "md:") {
		t.Errorf("unexpected breakpoint classes in:\n%s", html)
	}
}

// TestCellSpan checks a cell spans no more columns than its grid has.
func TestCellSpan(t *testing.T) {
	html := renderWith(t, Grid(props.Grid{Columns: 3}), templ.Join(
		renderedCell(props.GridCell{Span: 4, Rows: 9}),
	))
	if !hasClasses(html, "md:col-span-3 md:row-span-4") || strings.Contains(html, "col-span-4") {
		t.Errorf("want the span clamped to 3 columns in:\n%s", html)
	}
	html = renderWith(t, Grid(props.Grid{Columns: 1}), templ.Join(renderedCell(props.GridCell{Span: 2})))
	if strings.Contains(html, "col-span") {
		t.Errorf("want no span in a single column in:\n%s", html)
	}
}

// renderedCell renders a cell without children, so a grid's children
// can be joined.
func renderedCell(p props.GridCell) templ.Component {
	return Cell(p)
}

// hasClasses reports whether one of the class attributes in the html, or
// the html itself as a list of classes, has all the classes in want.
func hasClasses(html, want string) bool {
	lists := []string{html}
	for _, part := range strings.Split(html, `class="`)[1:] {
		lists = append(lists, part[:strings.IndexByte(part, '"')])
	}
	for _, list := range lists {
		have := strings.Fields(list)
		all := true
		for _, w := range strings.Fields(want) {
			if !slices.Contains(have, w) {
				all = false
				break
			}
		}
		if all {
			return true
		}
	}
	return false
}
