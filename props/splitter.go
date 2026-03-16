package props

import (
	"github.com/a-h/templ"
)

// Splitter orientations
const (
	SplitterOrientationVertical   string = "vertical"
	SplitterOrientationHorizontal string = "horizontal"
)

// Splitter defines the properties for a resizable split pane component.
// Uses Alpine.js for drag interaction.
type Splitter struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	// Orientation defines the split direction: "vertical" (top/bottom, default)
	// or "horizontal" (left/right).
	Orientation string
	// InitialPercent sets the initial size of the first pane as a percentage (default 50).
	InitialPercent int
	// MinPercent sets the minimum size of the first pane (default 15).
	MinPercent int
	// MaxPercent sets the maximum size of the first pane (default 85).
	MaxPercent int

	// First is the content rendered in the first (top/left) panel.
	First templ.Component
	// Second is the content rendered in the second (bottom/right) panel.
	Second templ.Component

	// FirstClass adds CSS classes to the first panel.
	FirstClass string
	// SecondClass adds CSS classes to the second panel.
	SecondClass string
}

// GetInitialPercent returns the initial percent with a default of 50.
func (s Splitter) GetInitialPercent() int {
	if s.InitialPercent > 0 {
		return s.InitialPercent
	}
	return 50
}

// GetMinPercent returns the min percent with a default of 15.
func (s Splitter) GetMinPercent() int {
	if s.MinPercent > 0 {
		return s.MinPercent
	}
	return 15
}

// GetMaxPercent returns the max percent with a default of 85.
func (s Splitter) GetMaxPercent() int {
	if s.MaxPercent > 0 {
		return s.MaxPercent
	}
	return 85
}

// IsHorizontal returns true if the splitter is horizontal (left/right).
func (s Splitter) IsHorizontal() bool {
	return s.Orientation == SplitterOrientationHorizontal
}
