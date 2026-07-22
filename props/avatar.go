package props

import "github.com/a-h/templ"

// Avatar sizes
const (
	AvatarSizeLarge string = "lg"
)

// Avatar provides props for the avatar component
type Avatar struct {
	ID         string
	Class      string
	Attributes templ.Attributes
	// Initial is the text shown when no image child is provided. Up to
	// three characters are rendered; anything longer is truncated.
	Initial string
	// Color is a foreground color token name applied to the initials:
	// accent, success, warning, critical, info, document-xml,
	// document-pdf or document-png. Empty uses the default secondary
	// foreground.
	Color string
	Size  string
}

// Initials returns Initial truncated to at most three characters.
func (a Avatar) Initials() string {
	r := []rune(a.Initial)
	if len(r) > 3 {
		r = r[:3]
	}
	return string(r)
}
