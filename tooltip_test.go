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
