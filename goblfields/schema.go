package goblfields

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"strings"
)

// definition is a single node of a GOBL JSON Schema document: the document
// root, one of its $defs entries, or a property inside one. GOBL only ever
// uses the subset described here, plus hints (`calculated`, `recommended`)
// the field tree ignores — `calculated` in particular does not mean "always
// present": issue_time is calculated but optional.
type definition struct {
	Defs        map[string]*definition `json:"$defs"`
	ID          string                 `json:"$id"`
	Ref         string                 `json:"$ref"`
	Description string                 `json:"description"`
	Items       *definition            `json:"items"`
	Properties  *propertySet           `json:"properties"`
	Required    []string               `json:"required"`
	Title       string                 `json:"title"`
	Type        string                 `json:"type"`
}

// requires reports whether the definition lists the named property in its
// `required` array.
func (d *definition) requires(name string) bool {
	if d == nil {
		return false
	}
	for _, r := range d.Required {
		if r == name {
			return true
		}
	}
	return false
}

// propertySet is a JSON object of properties that keeps the order they were
// written in. GOBL orders properties meaningfully ($schema, uuid, type,
// series, code, issue_date, …) and the picker lists them in that order, which
// a plain map[string]*definition would lose.
type propertySet struct {
	names []string
	defs  map[string]*definition
}

// UnmarshalJSON decodes the object one key at a time so the order survives.
func (p *propertySet) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))

	tok, err := dec.Token()
	if err != nil {
		return err
	}
	if tok != json.Delim('{') {
		return fmt.Errorf("properties: expected object, got %v", tok)
	}

	p.defs = make(map[string]*definition)
	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			return err
		}
		name, ok := tok.(string)
		if !ok {
			return fmt.Errorf("properties: expected name, got %v", tok)
		}
		def := new(definition)
		if err := dec.Decode(def); err != nil {
			return fmt.Errorf("properties: %s: %w", name, err)
		}
		p.names = append(p.names, name)
		p.defs[name] = def
	}

	_, err = dec.Token() // closing brace
	return err
}

// len reports the number of properties, tolerating a nil set so callers can
// test "has properties" without a nil check of their own.
func (p *propertySet) len() int {
	if p == nil {
		return 0
	}
	return len(p.names)
}

// target is a resolved definition together with the document it came from.
// Key identifies it across documents and is what stops a branch expanding the
// same schema twice: local refs (#/$defs/SubLine) resolve against the owning
// document, so it has to carry the document id as well as the name.
type target struct {
	def   *definition
	docID string
	key   string
}

// SchemaID accepts either a short path (bill/invoice) or a full schema id and
// always returns the full id, without the fragment or any query modifiers
// (e.g. ?tax_regime=) that GOBL allows on a schema reference.
func SchemaID(schema string) string {
	base, _, _ := strings.Cut(schema, "#")
	path, _, _ := strings.Cut(base, "?")
	if strings.HasPrefix(path, SchemaPrefix) {
		return path
	}
	return SchemaPrefix + strings.TrimLeft(path, "/")
}

// SchemaPath is the short form of a schema id, as used in workflows and silo
// entries: https://gobl.org/draft-0/bill/invoice becomes bill/invoice.
func SchemaPath(schema string) string {
	return strings.TrimPrefix(SchemaID(schema), SchemaPrefix)
}

// readDocument loads and decodes the schema document for a full schema id.
func readDocument(source fs.FS, id string) (*definition, error) {
	name := SchemaPath(id) + ".json"

	content, err := fs.ReadFile(source, name)
	if err != nil {
		return nil, fmt.Errorf("goblfields: %s: %w", SchemaPath(id), err)
	}

	doc := new(definition)
	if err := json.Unmarshal(content, doc); err != nil {
		return nil, fmt.Errorf("goblfields: %s: %w", SchemaPath(id), err)
	}

	return doc, nil
}
