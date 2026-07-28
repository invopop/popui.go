package popui

import (
	"strings"
	"testing"

	"github.com/invopop/popui.go/props"
)

// TestTooltipPositionClasses covers the placement variants, in particular that
// the -start and -end ones align the card with an edge of the trigger instead of
// centering it: centering is what makes a wide card overflow when the trigger
// sits near the edge of its container.
func TestTooltipPositionClasses(t *testing.T) {
	const (
		centered = "-translate-x-1/2"
		above    = "bottom-full" // the card sits above, so its bottom meets the trigger's top
		below    = "top-full"
	)

	for _, tc := range []struct {
		position string
		want     []string
		centered bool
	}{
		{props.TooltipPositionTop, []string{above}, true},
		{props.TooltipPositionBottom, []string{below}, true},
		{props.TooltipPositionTopStart, []string{above, "start-0"}, false},
		{props.TooltipPositionTopEnd, []string{above, "end-0"}, false},
		{props.TooltipPositionBottomStart, []string{below, "start-0"}, false},
		{props.TooltipPositionBottomEnd, []string{below, "end-0"}, false},
		{props.TooltipPositionLeft, []string{"right-full"}, false},
		{props.TooltipPositionRight, []string{"left-full"}, false},
		{"", []string{above}, true}, // defaults to top
	} {
		t.Run(tc.position, func(t *testing.T) {
			got := tooltipPositionClasses(tc.position)
			for _, want := range tc.want {
				if !strings.Contains(got, want) {
					t.Errorf("position %q: expected %q in %q", tc.position, want, got)
				}
			}
			if isCentered := strings.Contains(got, centered); isCentered != tc.centered {
				t.Errorf("position %q: centered=%v, want %v (%q)", tc.position, isCentered, tc.centered, got)
			}
		})
	}
}

// TestLabelTooltipDefaultsToStart covers the placement a label's help icon gets
// by default. The icon sits at the leading edge of a field, so a centered card
// would hang off that edge — every docs example needed a padding workaround
// before this default existed.
func TestLabelTooltipDefaultsToStart(t *testing.T) {
	t.Run("hint", func(t *testing.T) {
		got := props.Label{Hint: "explain"}.LabelTooltip()
		if got.Description != "explain" {
			t.Errorf("hint should become the description, got %q", got.Description)
		}
		if got.Position != props.TooltipPositionTopStart {
			t.Errorf("position = %q, want %q", got.Position, props.TooltipPositionTopStart)
		}
	})

	t.Run("tooltip without a position", func(t *testing.T) {
		got := props.Label{Tooltip: props.Tooltip{Title: "t"}}.LabelTooltip()
		if got.Position != props.TooltipPositionTopStart {
			t.Errorf("position = %q, want %q", got.Position, props.TooltipPositionTopStart)
		}
	})

	t.Run("an explicit position wins", func(t *testing.T) {
		got := props.Label{Tooltip: props.Tooltip{
			Title:    "t",
			Position: props.TooltipPositionBottomEnd,
		}}.LabelTooltip()
		if got.Position != props.TooltipPositionBottomEnd {
			t.Errorf("position = %q, want it left alone", got.Position)
		}
	})

	t.Run("no hint means no card", func(t *testing.T) {
		if got := (props.Label{}).LabelTooltip(); !got.Empty() {
			t.Errorf("expected an empty tooltip, got %+v", got)
		}
	})
}
