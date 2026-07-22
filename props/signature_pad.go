package props

// SignaturePad defines the properties for the signature dialog component.
type SignaturePad struct {
	// Action is the form POST action URL.
	Action string
	// Signer is the full name shown in the typed-signature preview.
	Signer string
	// Title is the dialog heading.
	Title string
	// TabTyped is the label for the typed-name tab.
	TabTyped string
	// TabDraw is the label for the freehand drawing tab.
	TabDraw string
	// UseLabel is the label for the confirm button.
	UseLabel string
	// ClearLabel is the label for the clear button.
	ClearLabel string
}
