package props

import "github.com/a-h/templ"

// CardDeck Templ component props. CardDeck is a tinted container that
// stacks regular Card children (optionally led by a CardDeckHead row).
type CardDeck struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

// CardDeckHead Templ component props. A small muted row at the top of a
// CardDeck for a label and an optional trailing action.
type CardDeckHead struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}
