package props

import (
	"fmt"

	"github.com/a-h/templ"
)

// Card Templ component props
type Card struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Href       templ.SafeURL
	Disabled   bool

	// Order positions the card inside a reorderable CardDeck, 1-based, and
	// renders as `data-order`. The deck sorts its cards by it on load and
	// renumbers every card to 1..n after each move, so the attribute, the
	// DOM sequence and the visible sequence are always the same list — read
	// it back to persist an order. Cards without one sink to the end of the
	// deck, keeping the sequence they were rendered in. Ignored outside a
	// reorderable deck, where the rendered sequence is already the order.
	Order int
}

// CardContent Templ component props
type CardContent struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

// CardHeader Templ component props
type CardHeader struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	Title      string
	Subtitle   string
}

// CardProgressBar Templ component props
type CardProgressBar struct {
	ID          string
	Class       string
	Attributes  templ.Attributes
	Title       string
	Subtitle    string
	Total       int64
	Current     int64
	HideCounter bool
}

// PercentValue returns the progress percentage, capped at 100.
func (p CardProgressBar) PercentValue() int64 {
	if p.Total == 0 {
		return 0
	}
	percent := int64((float64(p.Current) / float64(p.Total)) * 100)
	if percent > 100 {
		return 100
	}
	return percent
}

// PercentColor returns the expected progress bar color based on PercentValue
func (p CardProgressBar) PercentColor() string {
	if p.PercentValue() == 100 {
		return "#C92D45"
	}
	return "#008852"
}

// FormatAmount returns a formatted string representation of an amount
func (p CardProgressBar) FormatAmount(amount int64) string {
	return fmt.Sprintf("%d", amount)
}
