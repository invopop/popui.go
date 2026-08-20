package props

import "github.com/a-h/templ"

// CardDeck Templ component props. CardDeck is a tinted container that
// stacks regular Card children (optionally led by a CardDeckHead row).
type CardDeck struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	// Reorderable turns the deck into a drag-to-reorder list. Place a
	// CardDeckReorderButton in the head to toggle reorder mode: the child
	// cards shake briefly, take a grab cursor, become draggable along
	// the vertical axis only, and have their own links and menus made inert.
	// Every button in the deck except the toggle is disabled until reorder
	// mode is switched off again.
	//
	// Give every child Card a stable `data-card-id` attribute so the new
	// order can be read back — decks without one fall back to the card's
	// element id, then to its 1-based position when the deck loaded, which
	// stays with the card as it moves so the permutation is still readable.
	//
	// The cards' Order fields are the deck's source of truth: it sorts by
	// them on load and renumbers every card to 1..n after each move, so
	// `data-order` always matches the visible sequence.
	Reorderable bool

	// Name, on a reorderable deck, renders a hidden input holding the card
	// ids in their current order, comma separated, so an enclosing form or
	// HTMX submit carries the order. Regardless of Name, every committed
	// move also dispatches a bubbling `popui-card-deck-reorder` event whose
	// `detail.order` is the ordered id list.
	Name string
}

// CardDeckHead Templ component props. A small muted row at the top of a
// CardDeck for a label and an optional trailing action.
type CardDeckHead struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

// CardDeckReorderButton Templ component props. The button toggles its
// deck's reorder mode, so it only works inside a CardDeck built with
// Reorderable set.
type CardDeckReorderButton struct {
	ID         string
	Class      string
	Attributes templ.Attributes

	// Label is the resting label, "Reorder" by default.
	Label string
	// DoneLabel is the label shown while reorder mode is active, "Done" by
	// default. It renders on a success-green button to read as the way out.
	DoneLabel string
}

// ReorderLabel returns the resting label, defaulting to "Reorder".
func (b CardDeckReorderButton) ReorderLabel() string {
	if b.Label == "" {
		return "Reorder"
	}
	return b.Label
}

// ReorderDoneLabel returns the active-state label, defaulting to "Done".
func (b CardDeckReorderButton) ReorderDoneLabel() string {
	if b.DoneLabel == "" {
		return "Done"
	}
	return b.DoneLabel
}
