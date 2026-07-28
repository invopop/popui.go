package props

import "github.com/a-h/templ"

// Label defines the tag properties for labels.
type Label struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Hint shows a question mark icon after the label that reveals a small
	// Tooltip card with this text on hover. Parsed as Markdown, like any
	// Tooltip description.
	Hint string
	// Tooltip shows a question mark icon after the label that reveals a full
	// Tooltip card (title, description, optional illustration) on hover.
	// Takes precedence over Hint.
	Tooltip Tooltip
}

// LabelTooltip resolves the Tooltip card shown next to the label: the
// Tooltip prop when set, otherwise a description-only card from Hint.
//
// A label sits at the leading edge of its field, so a card centered on the icon
// hangs off that edge — the reason it defaults to TooltipPositionTopStart, which
// opens the card inwards. An explicit Position still wins.
func (l Label) LabelTooltip() Tooltip {
	t := l.Tooltip
	if t.Empty() && l.Hint != "" {
		t = Tooltip{Description: l.Hint}
	}
	if !t.Empty() && t.Position == "" {
		t.Position = TooltipPositionTopStart
	}
	return t
}
