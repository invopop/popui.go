package props

import "github.com/a-h/templ"

// SegmentedTabView is a self-contained tabbed view: it owns the Alpine scope
// that tracks the active tab, renders a popui.SegmentedTabs control, and hosts
// one SegmentedTabPanel per tab in its children. It removes the hand-rolled
// x-data / x-show / padding wrappers consumers used to write around
// SegmentedTabs.
//
//	@popui.SegmentedTabView(props.SegmentedTabView{Tabs: []props.SegmentedTab{
//	    {Value: "a", Label: "First"}, {Value: "b", Label: "Second"},
//	}}) {
//	    @popui.SegmentedTabPanel(props.SegmentedTabPanel{For: "a"}) { … }
//	    @popui.SegmentedTabPanel(props.SegmentedTabPanel{For: "b"}) { … }
//	}
type SegmentedTabView struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Model is the Alpine variable name holding the active tab Value.
	// Defaults to "tab" — match it on the child SegmentedTabPanels (their
	// default is also "tab", so usually neither needs setting).
	Model string
	// Active is the tab Value selected on first render. Defaults to the
	// first tab's Value.
	Active string
	// FullWidth stretches the control to equal-width triggers.
	FullWidth bool
	Tabs      []SegmentedTab
}

// SegmentedTabPanel is one tab's content inside a SegmentedTabView. It shows
// only while the view's Model equals For. Place its content in the children.
type SegmentedTabPanel struct {
	Class      string
	Attributes templ.Attributes
	// For is the tab Value this panel is shown for (matches a SegmentedTab.Value).
	For string
	// Model is the Alpine variable name to test. Defaults to "tab"; set it
	// only if the parent SegmentedTabView uses a custom Model.
	Model string
}
