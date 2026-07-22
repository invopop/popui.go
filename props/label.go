package props

import "github.com/a-h/templ"

// Label defines the tag properties for labels.
type Label struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Hint shows a question mark icon after the label that reveals a small
	// Tooltip card with this text on hover.
	Hint string
	// Tooltip shows a question mark icon after the label that reveals a full
	// Tooltip card (title, description, optional illustration) on hover.
	// Takes precedence over Hint.
	Tooltip Tooltip
}

// LabelTooltip resolves the Tooltip card shown next to the label: the
// Tooltip prop when set, otherwise a description-only card from Hint.
func (l Label) LabelTooltip() Tooltip {
	if !l.Tooltip.Empty() {
		return l.Tooltip
	}
	if l.Hint != "" {
		return Tooltip{Description: l.Hint}
	}
	return Tooltip{}
}
