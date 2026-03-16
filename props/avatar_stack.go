package props

import "github.com/a-h/templ"

// AvatarStack defines the properties for a stack of overlapping avatars.
type AvatarStack struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

// AvatarStackItem defines the properties for a single item in an avatar stack.
type AvatarStackItem struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}
