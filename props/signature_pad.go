package props

// SignaturePad defines the properties for the signature dialog component.
type SignaturePad struct {
	// Action is the form POST action URL.
	Action string
	// Signer is the full name shown in the typed-signature preview.
	Signer string
	// Locale sets the UI language. Supported: "en" (default), "fr", "es".
	Locale string
}

// WithDefaults returns a SignaturePad with the default locale set if empty.
func (s SignaturePad) WithDefaults() SignaturePad {
	if s.Locale == "" {
		s.Locale = "en"
	}
	return s
}
