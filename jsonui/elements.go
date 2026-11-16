package jsonui

import (
	"context"
	"io"
	"maps"
	"strconv"

	"github.com/a-h/templ"
	popui "github.com/invopop/popui.go"
	"github.com/invopop/popui.go/props"
)

// Element is the props of the plain HTML elements that structure a table
// or a list, which popui styles from the component around them, such as
// "table".
type Element struct {
	Class string
	// Attributes are added to the element.
	Attributes templ.Attributes
}

// Cell is the props of a table's "th" and "td" cells.
type Cell struct {
	Class string
	// Colspan is the number of columns the cell spans, e.g. for a row that
	// says the table is empty.
	Colspan int
	// Align places the cell's content, see the popui Align constants.
	// Default is "start", or "end" for numeric cells.
	Align string `jsonschema:"enum=,enum=start,enum=center,enum=end"`
	// Numeric is for cells of amounts and quantities: aligned to the end,
	// kept on one line, with digits in equal widths so they line up. Set it
	// on the column's header cell as well.
	Numeric bool
	// Attributes are added to the cell.
	Attributes templ.Attributes
}

// attributes adds the cell's alignment to its attributes, for the table's
// styles.
func (c Cell) attributes() templ.Attributes {
	align := c.Align
	if align == "" && c.Numeric {
		align = props.AlignEnd
	}
	if align == "" {
		return c.Attributes
	}
	attrs := templ.Attributes{"data-align": align}
	if c.Numeric {
		attrs["data-numeric"] = true
	}
	maps.Copy(attrs, c.Attributes)
	return attrs
}

// elementComponent renders the named element with its children.
func elementComponent(tag string) func(...Element) templ.Component {
	return func(p ...Element) templ.Component {
		var e Element
		if len(p) > 0 {
			e = p[0]
		}
		return element(tag, e.Class, 0, e.Attributes)
	}
}

// cellComponent renders the named table cell with its children.
func cellComponent(tag string) func(...Cell) templ.Component {
	return func(p ...Cell) templ.Component {
		var c Cell
		if len(p) > 0 {
			c = p[0]
		}
		return element(tag, c.Class, c.Colspan, c.attributes())
	}
}

// element renders an element from the fixed list registered below, with
// its children.
func element(tag, class string, colspan int, attrs templ.Attributes) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, "<"+tag); err != nil {
			return err
		}
		if class != "" {
			if _, err := io.WriteString(w, ` class="`+templ.EscapeString(class)+`"`); err != nil {
				return err
			}
		}
		if colspan > 1 {
			if _, err := io.WriteString(w, ` colspan="`+strconv.Itoa(colspan)+`"`); err != nil {
				return err
			}
		}
		if err := templ.RenderAttributes(ctx, w, attrs); err != nil {
			return err
		}
		if _, err := io.WriteString(w, ">"); err != nil {
			return err
		}
		if err := templ.GetChildren(ctx).Render(ctx, w); err != nil {
			return err
		}
		_, err := io.WriteString(w, "</"+tag+">")
		return err
	})
}

func init() {
	// Tables
	register("table", popui.Table)
	for _, tag := range []string{"thead", "tbody", "tr"} {
		register(tag, elementComponent(tag))
	}
	for _, tag := range []string{"th", "td"} {
		register(tag, cellComponent(tag))
	}

	// Lists
	for _, tag := range []string{"ul", "ol", "li"} {
		register(tag, elementComponent(tag))
	}
}
