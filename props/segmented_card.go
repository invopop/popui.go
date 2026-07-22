package props

import "github.com/a-h/templ"

// SegmentedCard Templ component props.
//
// Deprecated: use CardDeck instead.
type SegmentedCard struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

// SegmentedCardHead Templ component props.
//
// Deprecated: use CardDeckHead instead.
type SegmentedCardHead struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

// SegmentedCardBody Templ component props.
//
// Deprecated: place Card children directly inside CardDeck instead.
type SegmentedCardBody struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

// SegmentedCardContent Templ component props.
//
// Deprecated: use a plain Card inside CardDeck instead.
type SegmentedCardContent struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}
