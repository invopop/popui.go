package jsonui

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/santhosh-tekuri/jsonschema/v6/kind"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Documents describe what their data must be with JSON Schemas, one for each
// part of the data they validate, e.g. the "config" a step saves:
//
//	{
//	  "data": {"config": {"to": ""}},
//	  "schemas": {
//	    "config": {
//	      "type": "object",
//	      "required": ["to"],
//	      "properties": {"to": {"type": "string", "format": "email"}}
//	    }
//	  },
//	  "components": [{"is": "input", "model": "config.to", "props": {"label": "To"}}]
//	}
//
// The data is validated as it changes, and each error is shown with the form
// control whose model is the path it's about. The same schemas can be used
// wherever the data is saved, so the browser and the server agree, and tell
// anyone writing the data without the form, such as an LLM, what it must be.
// Labels and options belong to the components, independently of the schemas.

// FieldError is a problem with the value at a data path, such as
// "config.to", or "config" for the object itself.
type FieldError struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

// Messages used for several kinds of error.
const (
	msgRequired = "Required"
	msgEmail    = "Must be an email address"
)

// schemaURL is the base URL of the schemas in a document, which never
// reach outside it.
const schemaURL = "jsonui:schemas"

// Validator validates document data against the document's schemas.
type Validator struct {
	paths   []string
	schemas map[string]*jsonschema.Schema
	// raw are the schemas as JSON values, used to resolve errorMessage overrides.
	raw map[string]any
}

// NewValidator compiles a document's schemas, keyed by data path. Schemas
// can only refer to definitions inside themselves with "$ref", as documents
// may come from anywhere and must not make the renderer fetch anything.
func NewValidator(schemas map[string]json.RawMessage) (*Validator, error) {
	v := &Validator{schemas: map[string]*jsonschema.Schema{}, raw: map[string]any{}}
	for path, raw := range schemas {
		if !safePath(path) {
			return nil, fmt.Errorf("schemas: %q is not a data path", path)
		}
		doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(raw))
		if err != nil {
			return nil, fmt.Errorf("schemas %s: %w", path, err)
		}
		c := jsonschema.NewCompiler()
		c.UseLoader(noLoader{})
		c.AssertFormat()
		url := schemaURL + "/" + path
		if err := c.AddResource(url, doc); err != nil {
			return nil, fmt.Errorf("schemas %s: %w", path, err)
		}
		s, err := c.Compile(url)
		if err != nil {
			return nil, fmt.Errorf("schemas %s: %w", path, schemaError(err))
		}
		v.paths = append(v.paths, path)
		v.schemas[path] = s
		v.raw[path] = doc
	}
	sort.Strings(v.paths) // stable error order
	return v, nil
}

// noLoader refuses to load schemas from anywhere: documents' schemas must be
// complete.
type noLoader struct{}

func (noLoader) Load(url string) (any, error) {
	return nil, fmt.Errorf("can't load %s: schemas may only refer to themselves", url)
}

// schemaError keeps the cause of a compile error, as the rest repeats the
// internal URL.
func schemaError(err error) error {
	var le *jsonschema.LoadURLError
	if errors.As(err, &le) {
		return le.Err
	}
	var se *jsonschema.SchemaValidationError
	if errors.As(err, &se) {
		return se.Err
	}
	return err
}

// Validate checks the document data against the schemas, returning every
// error found, one per path at most, in order of path.
func (v *Validator) Validate(data map[string]any) []FieldError {
	var out []FieldError
	seen := map[string]bool{}
	add := func(e FieldError) {
		if !seen[e.Path] {
			seen[e.Path] = true
			out = append(out, e)
		}
	}
	for _, path := range v.paths {
		value, err := lookup(data, path)
		if err != nil {
			value = nil // missing, so the schema decides whether it may be
		}
		err = v.schemas[path].Validate(value)
		var ve *jsonschema.ValidationError
		if !errors.As(err, &ve) {
			continue
		}
		for _, e := range v.fieldErrors(path, ve) {
			add(e)
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}

// ValidateJSON validates document data given as JSON, as the browser sends
// it.
func (v *Validator) ValidateJSON(data []byte) ([]FieldError, error) {
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("data: %w", err)
	}
	obj, ok := doc.(map[string]any)
	if !ok {
		return nil, errors.New("data: must be an object")
	}
	return v.Validate(obj), nil
}

// fieldErrors flattens a validation error into an error for each path the
// problems are at. Groups of alternatives (anyOf, oneOf) that fail on the
// value itself are reported once, as listing why every alternative failed
// says more about the schema than the value.
func (v *Validator) fieldErrors(root string, ve *jsonschema.ValidationError) []FieldError {
	switch ve.ErrorKind.(type) {
	case *kind.AnyOf, *kind.OneOf:
		if allAt(ve.Causes, ve.InstanceLocation) {
			return []FieldError{{Path: joinPath(root, ve.InstanceLocation), Message: v.message(ve, oneOfMessage(ve))}}
		}
	}
	if len(ve.Causes) > 0 {
		var out []FieldError
		for _, c := range ve.Causes {
			out = append(out, v.fieldErrors(root, c)...)
		}
		return out
	}
	at := joinPath(root, ve.InstanceLocation)
	switch k := ve.ErrorKind.(type) {
	case *kind.Required:
		out := make([]FieldError, len(k.Missing))
		for i, prop := range k.Missing {
			out[i] = FieldError{Path: at + "." + prop, Message: v.requiredMessage(ve, prop)}
		}
		return out
	case *kind.AdditionalProperties:
		out := make([]FieldError, len(k.Properties))
		for i, prop := range k.Properties {
			out[i] = FieldError{Path: at + "." + prop, Message: v.message(ve, "Not allowed")}
		}
		return out
	}
	return []FieldError{{Path: at, Message: v.message(ve, defaultMessage(ve.ErrorKind))}}
}

// allAt reports whether every error in a tree is about the value at loc.
func allAt(errs []*jsonschema.ValidationError, loc []string) bool {
	for _, e := range errs {
		if !slices.Equal(e.InstanceLocation, loc) || !allAt(e.Causes, loc) {
			return false
		}
	}
	return true
}

func joinPath(root string, loc []string) string {
	if len(loc) == 0 {
		return root
	}
	return root + "." + strings.Join(loc, ".")
}

// message is the error's "errorMessage" from its schema when there is one,
// or def. "errorMessage" is either one message for every error in the
// schema, or an object with a message for each keyword, e.g.
// {"format": "Enter an email address"}.
func (v *Validator) message(ve *jsonschema.ValidationError, def string) string {
	kw := ve.ErrorKind.KeywordPath()
	if len(kw) == 0 {
		return def
	}
	s, _ := v.subschema(ve.SchemaURL).(map[string]any)
	switch m := s["errorMessage"].(type) {
	case string:
		return m
	case map[string]any:
		if msg, ok := m[kw[0]].(string); ok {
			return msg
		}
	}
	return def
}

// requiredMessage is the message for a missing property: the "required"
// entry of its parent's "errorMessage", which may also be an object with a
// message for each property, then its own schema's "errorMessage", then its
// parent's when that is one message for every error.
func (v *Validator) requiredMessage(ve *jsonschema.ValidationError, prop string) string {
	parent, _ := v.subschema(ve.SchemaURL).(map[string]any)
	if m, ok := parent["errorMessage"].(map[string]any); ok {
		switch r := m["required"].(type) {
		case string:
			return r
		case map[string]any:
			if msg, ok := r[prop].(string); ok {
				return msg
			}
		}
	}
	props, _ := parent["properties"].(map[string]any)
	child, _ := props[prop].(map[string]any)
	switch m := child["errorMessage"].(type) {
	case string:
		return m
	case map[string]any:
		if msg, ok := m["required"].(string); ok {
			return msg
		}
	}
	if m, ok := parent["errorMessage"].(string); ok {
		return m
	}
	return msgRequired
}

// subschema finds a schema in the raw schemas by its URL, such as
// "jsonui:schemas/config#/properties/to".
func (v *Validator) subschema(url string) any {
	base, ptr, _ := strings.Cut(url, "#")
	path, ok := strings.CutPrefix(base, schemaURL+"/")
	if !ok {
		return nil
	}
	return pointer(v.raw[path], ptr)
}

// pointer resolves a JSON pointer, such as "/properties/to", in a value.
func pointer(v any, ptr string) any {
	if ptr == "" {
		return v
	}
	for _, tok := range strings.Split(strings.TrimPrefix(ptr, "/"), "/") {
		tok = strings.NewReplacer("~1", "/", "~0", "~").Replace(tok)
		switch cur := v.(type) {
		case map[string]any:
			v = cur[tok]
		case []any:
			i := 0
			if _, err := fmt.Sscan(tok, &i); err != nil || i < 0 || i >= len(cur) {
				return nil
			}
			v = cur[i]
		default:
			return nil
		}
	}
	return v
}

// oneOfMessage describes a value matching none of the alternatives, listing
// them when they are all consts.
func oneOfMessage(ve *jsonschema.ValidationError) string {
	var want []string
	for _, c := range ve.Causes {
		k, ok := c.ErrorKind.(*kind.Const)
		for _, cc := range c.Causes {
			if ok {
				break
			}
			k, ok = cc.ErrorKind.(*kind.Const)
		}
		if !ok {
			return "Not one of the allowed values"
		}
		if empty(k.Got) {
			return msgRequired
		}
		want = append(want, display(k.Want))
	}
	return mustBeOneOf(want)
}

var printer = message.NewPrinter(language.English)

// defaultMessage describes an error to the person filling in the form,
// rather than to the author of the schema.
func defaultMessage(k jsonschema.ErrorKind) string {
	switch k := k.(type) {
	case *kind.Type:
		if k.Got == jsonNull {
			return msgRequired
		}
		return "Must be " + article(strings.Join(k.Want, " or "))
	case *kind.MinLength:
		if k.Want == 1 {
			return msgRequired
		}
		return fmt.Sprintf("Must be at least %d characters", k.Want)
	case *kind.MaxLength:
		return fmt.Sprintf("Must be at most %d characters", k.Want)
	case *kind.Pattern:
		return "Not in the expected format"
	case *kind.Format:
		switch k.Want {
		case "email", "idn-email":
			return msgEmail
		case "uri", "uri-reference", "iri", "iri-reference":
			return "Must be a URL"
		case "date":
			return "Must be a date"
		case "date-time":
			return "Must be a date and time"
		}
		return "Must be " + article(k.Want)
	case *kind.Minimum:
		return "Must be at least " + number(k.Want)
	case *kind.Maximum:
		return "Must be at most " + number(k.Want)
	case *kind.ExclusiveMinimum:
		return "Must be more than " + number(k.Want)
	case *kind.ExclusiveMaximum:
		return "Must be less than " + number(k.Want)
	case *kind.MultipleOf:
		return "Must be a multiple of " + number(k.Want)
	case *kind.MinItems:
		if k.Want == 1 {
			return msgRequired
		}
		return fmt.Sprintf("Must have at least %d items", k.Want)
	case *kind.MaxItems:
		return fmt.Sprintf("Must have at most %d items", k.Want)
	case *kind.UniqueItems:
		return "Must not have duplicate items"
	case *kind.Enum:
		if empty(k.Got) {
			return msgRequired
		}
		want := make([]string, len(k.Want))
		for i, w := range k.Want {
			want[i] = display(w)
		}
		return mustBeOneOf(want)
	case *kind.Const:
		return "Must be " + display(k.Want)
	case *kind.FalseSchema:
		return "Not allowed"
	}
	return capitalize(k.LocalizedString(printer))
}

// empty reports whether a value is what an untouched control holds.
func empty(v any) bool {
	return v == nil || v == ""
}

func mustBeOneOf(want []string) string {
	if len(want) == 1 {
		return "Must be " + want[0]
	}
	return "Must be one of " + strings.Join(want, ", ")
}

func article(s string) string {
	if s == "" {
		return s
	}
	if strings.ContainsRune("aeiou", rune(s[0])) {
		return "an " + s
	}
	return "a " + s
}

func number(r *big.Rat) string {
	if r.IsInt() {
		return r.Num().String()
	}
	f, _ := r.Float64()
	return fmt.Sprint(f)
}

func capitalize(s string) string {
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}
