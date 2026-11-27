package props

import "github.com/a-h/templ"

// Step statuses.
const (
	StepStatusUpcoming  = "upcoming"
	StepStatusCurrent   = "current"
	StepStatusCompleted = "completed"
)

// Steps orientations.
const (
	StepsOrientationVertical   = "vertical"
	StepsOrientationHorizontal = "horizontal"
	StepsOrientationResponsive = "responsive"
)

// Steps defines the props for the Steps component.
type Steps struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Label names the list for screen readers, such as "Progress".
	Label string
	// Orientation is "vertical" for a timeline with a label beside each
	// step and its sub-steps below it, "horizontal" for a bar of the
	// steps' markers alone, with a line between them, or "responsive"
	// for the timeline on wide screens and the bar on narrow ones.
	// Default is "vertical".
	Orientation string `jsonschema:"enum=,enum=vertical,enum=horizontal,enum=responsive"`
}

// Step defines the props for the Step component.
type Step struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Label says what the step is, shown beside its marker.
	Label string
	// Status is "completed" for a step that is done, whose marker is a
	// tick, "current" for the one in progress, or "upcoming" for those
	// still to come. Default is "upcoming".
	Status string `jsonschema:"enum=,enum=upcoming,enum=current,enum=completed"`
	// StatusLabel says the status to screen readers, which don't see the
	// marker, in the page's language. Default is "Completed" for a
	// completed step and "Upcoming" for an upcoming one. The current
	// step is announced as such without one.
	StatusLabel string
}

// StatusText is what a screen reader is told about the step's status:
// its StatusLabel, or the default for its Status.
func (s Step) StatusText() string {
	switch {
	case s.StatusLabel != "":
		return s.StatusLabel
	case s.Status == StepStatusCompleted:
		return "Completed"
	case s.Status == StepStatusCurrent:
		return ""
	default:
		return "Upcoming"
	}
}

// SubStep defines the props for the SubStep component.
type SubStep struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Label says what the sub-step is.
	Label string
	// Current marks the sub-step in progress. Its label stands out in
	// place of its step's.
	Current bool
}
