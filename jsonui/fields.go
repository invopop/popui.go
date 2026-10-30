package jsonui

import "github.com/invopop/popui.go/classes"

// Form controls with a "model" are fields: they show the error for their
// model's path, which jsonui.js finds by validating the data against the
// document's schemas, or which actions report, such as a host saving the
// data. Errors from schemas are only shown once the field has been changed
// or left, or after the "validate" action, so a new form isn't covered in
// them.

// errorIDName is the name of the Alpine x-id scope of a field, and errorID
// the expression of its error's id, unique to the field on the page.
const (
	errorIDName = "jsonui-error"
	errorID     = "$id('" + errorIDName + "')"
)

// fieldAttrs adds the bindings that mark a control as invalid, and described
// by its error, while its path has an error shown.
func (b *docBuilder) fieldAttrs(path string, attrs map[string]any) {
	b.fields++
	shown := shownError(path)
	if _, ok := attrs[":class"]; !ok { // only trusted documents set their own
		attrs[":class"] = shown + " ? '" + classes.FormFieldError() + "' : ''"
	}
	attrs[":aria-invalid"] = "!!" + shown
	_, own := attrs["aria-describedby"]
	if _, bound := attrs[":aria-describedby"]; !own && !bound {
		attrs[":aria-describedby"] = shown + " ? " + errorID + " : null"
	}
}
