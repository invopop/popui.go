package props

import "github.com/a-h/templ"

// SegmentedTab is one option in a popui.SegmentedTabs control.
type SegmentedTab struct {
	Value string // identifies the tab; assigned to the Alpine Model on click
	Label string // visible text
}

// SegmentedTabs renders a pill-style segmented control (a row of equal-width
// toggle buttons in a tinted track). It's Alpine-driven and stateless on its
// own: the active tab lives in the Model variable, defined by a surrounding
// x-data scope, so the same variable can also drive x-show panels next to the
// control.
//
//	<div x-data="{ tab: 'a' }">
//	    @popui.SegmentedTabs(props.SegmentedTabs{Model: "tab", Tabs: []props.SegmentedTab{
//	        {Value: "a", Label: "First"}, {Value: "b", Label: "Second"},
//	    }})
//	    <div x-show="tab === 'a'">…</div>
//	    <div x-show="tab === 'b'" x-cloak>…</div>
//	</div>
type SegmentedTabs struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Model is the Alpine variable (from a surrounding x-data) holding the
	// active tab's Value. Clicking a trigger assigns it.
	Model string
	// FullWidth stretches the control to fill its container with equal-width
	// triggers. Off by default — the control sizes to its labels.
	FullWidth bool
	Tabs      []SegmentedTab
}
