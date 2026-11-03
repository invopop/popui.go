package popui

import (
	"bytes"
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/invopop/popui.go/props"
)

func renderTable(t *testing.T, p props.Table) string {
	t.Helper()
	var buf bytes.Buffer
	ctx := templ.WithChildren(context.Background(), templ.Raw("<tbody><tr><td>1</td></tr></tbody>"))
	if err := Table(p).Render(ctx, &buf); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}

var tableTag = regexp.MustCompile(`<table[^>]*>`)

func invoiceColumns() []props.TableColumn {
	return []props.TableColumn{
		{Key: "series", Title: "Series", Width: 90, MinWidth: 60},
		{Key: "code", Title: "Code", Width: 180, Frozen: true, Sort: "asc",
			SortAscending:  templ.Attributes{"hx-get": "/?sort=code"},
			SortDescending: templ.Attributes{"hx-get": "/?sort=code&desc"}},
		{Key: "supplier", Title: "Supplier"},
		{Key: "total", Title: "Total", Align: props.TableColumnAlignRight},
	}
}

// TestTableColumns locks in what the server renders for a table with
// Columns, which popui.js relies on: the widths as custom properties on the
// table, the frozen columns by position, and a header cell per column
// carrying its key.
func TestTableColumns(t *testing.T) {
	got := renderTable(t, props.Table{Columns: invoiceColumns(), Resizable: true, Freezable: true})
	tag := tableTag.FindString(got)

	for _, want := range []string{
		`--popui-col-1:90px;--popui-col-2:180px;--popui-col-3:150px;--popui-col-4:150px;`,
		`data-popui-columns`,
		// Freezing Code freezes Series before it.
		`data-frozen="1 2"`,
		`data-frozen-last="2"`,
		"popui-table-columns",
	} {
		if !strings.Contains(tag, want) {
			t.Errorf("table tag missing %q: %s", want, tag)
		}
	}
	if n := strings.Count(got, "<col>"); n != 4 {
		t.Errorf("want 4 <col>, got %d", n)
	}
	if n := strings.Count(got, `class="popui-table-col-resizer"`); n != 4 {
		t.Errorf("want a resize handle per column, got %d", n)
	}
	series := thTag(got, "series")
	for _, want := range []string{`data-min-width="60"`, "data-frozen"} {
		if !strings.Contains(series, want) {
			t.Errorf("series header missing %q: %s", want, series)
		}
	}
	supplier := thTag(got, "supplier")
	// The default minimum.
	if !strings.Contains(supplier, `data-min-width="60"`) {
		t.Errorf("supplier header has no default minimum: %s", supplier)
	}
	if strings.Contains(supplier, "data-frozen") {
		t.Error("a column after the last frozen one must not be frozen")
	}
	if !strings.Contains(thTag(got, "code"), `aria-sort="ascending"`) {
		t.Error("the sorted column's header must carry aria-sort")
	}
	for _, want := range []string{
		`hx-get="/?sort=code"`,
		"Sort ascending",
		"Freeze column",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q", want)
		}
	}
}

// thTag is the opening tag of the header cell of the column with the key.
func thTag(html, key string) string {
	return regexp.MustCompile(`<th[^>]*data-key="` + key + `"[^>]*>`).FindString(html)
}

// TestTableColumnsCSS checks the stylesheet: the last column takes the
// space left, so it has no width of its own, frozen columns are offset by
// the widths before them, and right aligned columns align their cells.
func TestTableColumnsCSS(t *testing.T) {
	cols := tableColumnsOf(invoiceColumns())
	css := tableColumnsCSS(cols)
	c := "." + tableColumnsClass(cols)
	for _, want := range []string{
		c + "{width:max(100%,calc(var(--popui-col-1) + var(--popui-col-2) + var(--popui-col-3) + var(--popui-col-4)))}",
		c + ">colgroup>:nth-child(3){width:var(--popui-col-3)}",
		c + `[data-frozen~="2"]>*>tr>:nth-child(2){position:sticky;left:calc(var(--popui-col-1));z-index:5}`,
		c + `[data-frozen~="1"]>*>tr>:nth-child(1){position:sticky;left:0px;z-index:5}`,
		c + ">tbody>tr>:nth-child(4){text-align:right}",
	} {
		if !strings.Contains(css, want) {
			t.Errorf("css missing %q", want)
		}
	}
	if strings.Contains(css, ">colgroup>:nth-child(4)") {
		t.Error("the last column must not have a width")
	}
	if tableColumnsClass(cols) != tableColumnsClass(tableColumnsOf(invoiceColumns())) {
		t.Error("the class must be the same for tables of the same shape")
	}
}

// TestTableWithoutColumns checks tables without Columns keep their markup.
func TestTableWithoutColumns(t *testing.T) {
	got := renderTable(t, props.Table{})
	if strings.Contains(got, "<colgroup") || strings.Contains(got, "<style") || strings.Contains(got, "data-popui-columns") {
		t.Errorf("a table without columns must not be laid out from them: %s", got)
	}
	if !strings.Contains(tableTag.FindString(got), "popui-table-full-width") {
		t.Error("a table without columns keeps its auto layout")
	}
}

// TestTableColumnsKeepsPageStyle checks a style the page sets on the table
// follows the widths rather than replacing them.
func TestTableColumnsKeepsPageStyle(t *testing.T) {
	got := renderTable(t, props.Table{
		Columns:    []props.TableColumn{{Title: "A", Width: 100}},
		Attributes: templ.Attributes{styleAttribute: "color:red"},
	})
	if !strings.Contains(tableTag.FindString(got), `style="--popui-col-1:100px;color:red"`) {
		t.Errorf("got %s", tableTag.FindString(got))
	}
}
