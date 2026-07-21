package props

import "github.com/a-h/templ"

// TabItem is one entry in a popui.Tabs trigger row: a Value that identifies the
// tab (assigned to the Alpine Model on click and matched by the content Tab), a
// visible Label, an optional leading Icon, and an optional Count.
type TabItem struct {
	Value string
	Label string
	Icon  templ.Component
	// Count renders a counter after the label when greater than zero: "(n)" on
	// the default variant, a small badge on the pill variant.
	Count int
}

// Tabs is a self-contained tabbed view: it owns the Alpine scope that tracks
// the active tab, renders the trigger row from Tabs, and hosts one popui.Tab
// content panel per tab in its children. All the switching logic lives here, so
// consumers just declare the tab list and drop the content inside each Tab.
//
//	@popui.Tabs(props.Tabs{Variant: "pill", Tabs: []props.TabItem{
//	    {Value: "a", Label: "First"}, {Value: "b", Label: "Second"},
//	}}) {
//	    @popui.Tab(props.Tab{Value: "a"}) { …first panel… }
//	    @popui.Tab(props.Tab{Value: "b"}) { …second panel… }
//	}
//
// Omit the children for a stand-alone trigger row (no panels).
type Tabs struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Variant    string // "default" (underline) or "pill"
	// Model is the Alpine variable holding the active tab Value. Defaults to
	// "tab" — match it on the child Tabs (their default is also "tab", so
	// usually neither needs setting).
	Model string
	// Active is the tab Value selected on first render. Defaults to the first
	// tab's Value.
	Active string
	// Tabs is the ordered trigger row.
	Tabs []TabItem
}

// Tab is one tab's content panel inside a popui.Tabs. It is shown only while
// the view's Model equals Value. Place the panel content in the children.
type Tab struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Value is the tab this panel belongs to (matches a TabItem.Value).
	Value string
	// Model is the Alpine variable to test. Defaults to "tab"; set it only when
	// the parent Tabs uses a custom Model.
	Model string
}
