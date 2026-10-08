// Package goblfields builds the field trees consumed by popui's FieldPicker
// from the GOBL JSON Schemas embedded in github.com/invopop/gobl.
//
// Everything is resolved locally: there are no requests to the GOBL schema
// servers, and the field list a consumer renders is therefore pinned to the
// GOBL release its go.mod already selects.
//
//	fields, err := goblfields.Build("bill/invoice")
//	if err != nil {
//		return err
//	}
//	@popui.FieldPicker(props.FieldPicker{Root: "bill/invoice", Fields: fields, Format: "{{.doc.%s}}", Target: "#body"})
//
// Building a tree walks a few dozen schema documents, so consumers that render
// a picker on every request should build once and reuse the result: the tree
// only changes when the GOBL dependency does.
package goblfields

import (
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"

	"github.com/invopop/gobl/data"
	"github.com/invopop/popui.go/props"
)

// SchemaPrefix is the base every GOBL schema id is built on.
const SchemaPrefix = "https://gobl.org/draft-0/"

// Type labels the builder chooses itself, rather than copies from a schema:
// a field with properties, or with nothing more specific to say, is an
// object, and an array is just that. Scalar labels (string, integer, …) pass
// through from the schema's own type.
const (
	TypeArray  = "array"
	TypeObject = "object"
)

// DefaultMaxDepth is deep enough to reach every field GOBL can express: a
// branch stops growing on its own once it revisits a schema it already
// contains, so the cap is a backstop rather than the thing bounding the tree.
const DefaultMaxDepth = 6

// Schemas holds the GOBL schema documents, rooted so that bill/invoice.json
// is the document for bill/invoice. It defaults to the copy embedded in the
// gobl module; Options.Source overrides it.
var Schemas = mustSub(data.Content, "schemas")

// Options tune how a field tree is built.
type Options struct {
	// MaxDepth limits how far a branch is expanded, counting the root
	// document's own properties as depth 1. Defaults to DefaultMaxDepth.
	MaxDepth int
	// Descriptions keeps each field's schema description. They are long —
	// several lines for many GOBL fields — and multiply the size of the
	// payload the picker ships to the browser, so they are dropped unless
	// asked for.
	Descriptions bool
	// Source overrides where schema documents are read from. Defaults to
	// Schemas, the documents embedded in the gobl module.
	Source fs.FS
}

// Build returns the selectable field tree for a GOBL schema, given either as a
// short path (bill/invoice) or a full schema id.
//
// Branches stop at MaxDepth, and at any schema that already appears higher up
// the same branch, which is what keeps recursive schemas finite: a line with
// sub-lines, a party with a related party. A referenced document that is
// missing leaves that branch unexpanded rather than failing the build; only a
// missing root document is an error.
func Build(schema string, opts ...Options) ([]props.Field, error) {
	o := firstOption(opts)
	b := &builder{
		source:       source(o),
		maxDepth:     maxDepth(o),
		descriptions: o.Descriptions,
		registry:     make(map[string]*definition),
	}

	docID := SchemaID(schema)
	doc, err := b.document(docID)
	if err != nil {
		return nil, err
	}

	root := b.rootTarget(doc, docID)
	if root == nil {
		return nil, fmt.Errorf("goblfields: %s: could not resolve the root definition", SchemaPath(docID))
	}

	return b.walk(root, "", 1, map[string]bool{root.key: true}, true), nil
}

// MustBuild is Build for a schema known to be there, panicking when it is
// not. It is meant for package-level variables, where a tree is built once and
// reused for the life of the process:
//
//	var invoiceFields = goblfields.MustBuild("bill/invoice")
func MustBuild(schema string, opts ...Options) []props.Field {
	fields, err := Build(schema, opts...)
	if err != nil {
		panic(err)
	}
	return fields
}

// List returns the short paths of every schema document available to build
// from, in alphabetical order (bill/invoice, envelope, org/party, …).
func List(opts ...Options) ([]string, error) {
	src := source(firstOption(opts))

	var paths []string
	err := fs.WalkDir(src, ".", func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || path.Ext(name) != ".json" {
			return nil
		}
		paths = append(paths, strings.TrimSuffix(name, ".json"))
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("goblfields: listing schemas: %w", err)
	}

	sort.Strings(paths)
	return paths, nil
}

// Flatten returns every field in the tree, depth first, in the order the
// picker renders them.
func Flatten(fields []props.Field, into ...props.Field) []props.Field {
	out := into
	for _, field := range fields {
		out = append(out, field)
		out = Flatten(field.Children, out...)
	}
	return out
}

// Nest puts a tree under a single field, rewriting every path beneath it:
// Nest("doc", invoice) gives one object row, doc, whose children are the
// invoice's fields at doc.supplier.name and so on. It is for the hybrid
// list an app offers — its own variables at the top level, the document
// one row down — so that a single Format such as "{{.%s}}" serves both.
func Nest(name string, fields []props.Field) props.Field {
	return props.Field{
		Name:          name,
		Path:          name,
		Type:          TypeObject,
		AlwaysPresent: true,
		Children:      prefixPaths(name+".", fields),
	}
}

// Group puts fields under a heading of the picker's tree: it is opened like
// an object, but unlike Nest it leaves the fields' paths — and so what they
// emit — as they are. Use it to offer several trees in one picker, one per
// document schema a message might be about, say:
//
//	goblfields.Group("Invoice", invoiceFields)
//
// The heading has no Description; set one on the result when the schema's
// name should show beside it.
func Group(name string, fields []props.Field) props.Field {
	return props.Field{
		Name:     name,
		Path:     name,
		Type:     TypeObject,
		Group:    true,
		Children: fields,
	}
}

func prefixPaths(prefix string, fields []props.Field) []props.Field {
	out := make([]props.Field, len(fields))
	for i, field := range fields {
		field.Path = prefix + field.Path
		field.Children = prefixPaths(prefix, field.Children)
		out[i] = field
	}
	return out
}

// Trail returns the chain of fields leading to a path, root first. An unknown
// path resolves to the deepest ancestors that do exist, so a picker can still
// open near a value it cannot place exactly. A Group is looked through — its
// name is not part of the paths beneath it — and appears in the trail when
// the path is found under it; when several groups hold the same path, the
// first one wins.
func Trail(fields []props.Field, path string) []props.Field {
	if path == "" {
		return nil
	}
	return trail(fields, path)
}

func trail(level []props.Field, path string) []props.Field {
	for i := range level {
		field := &level[i]
		if field.Group {
			if below := trail(field.Children, path); len(below) > 0 {
				return append([]props.Field{*field}, below...)
			}
			continue
		}
		if field.Path == path {
			return []props.Field{*field}
		}
		if strings.HasPrefix(path, field.Path+".") {
			return append([]props.Field{*field}, trail(field.Children, path)...)
		}
	}
	return nil
}

type builder struct {
	source       fs.FS
	maxDepth     int
	descriptions bool
	// registry caches documents by schema id. A document that could not be
	// read is cached as nil so it is only looked for once.
	registry map[string]*definition
}

// walk turns the properties of a resolved definition into fields, recursing
// into the ones that can be expanded.
func (b *builder) walk(t *target, parentPath string, depth int, branch map[string]bool, ancestorsPresent bool) []props.Field {
	set := t.def.Properties
	if set.len() == 0 {
		return nil
	}

	fields := make([]props.Field, 0, set.len())
	for _, name := range set.names {
		fields = append(fields, b.field(t, set.defs[name], name, parentPath, depth, branch, ancestorsPresent))
	}
	return fields
}

func (b *builder) field(owner *target, property *definition, name, parentPath string, depth int, branch map[string]bool, ancestorsPresent bool) props.Field {
	isArray, node, resolved := b.resolveProperty(property, owner.docID)

	path := parentPath + name
	if isArray {
		path += "[]"
	}

	// "Always present" means required at every step down from the root.
	// Within an array it reads per item: lines[].item is marked when every
	// line that exists carries an item.
	required := owner.def.requires(name)
	alwaysPresent := ancestorsPresent && required

	field := props.Field{
		Name:          name,
		Path:          path,
		Type:          typeLabel(node, isArray),
		Array:         isArray,
		Required:      required,
		AlwaysPresent: alwaysPresent,
	}
	if b.descriptions {
		field.Description = strings.TrimSpace(firstString(property.Description, node.Description))
	}

	expandable := node.Properties.len() > 0 && depth < b.maxDepth
	if expandable && resolved == nil {
		// An object described inline, with no $ref. It cannot recur — only a
		// reference can point back up the tree — so it needs no cycle key.
		resolved = &target{def: node, docID: owner.docID}
	}
	if expandable && resolved.key != "" && branch[resolved.key] {
		expandable = false
	}
	if expandable {
		next := make(map[string]bool, len(branch)+1)
		for key := range branch {
			next[key] = true
		}
		if resolved.key != "" {
			next[resolved.key] = true
		}
		field.Children = b.walk(resolved, path+".", depth+1, next, alwaysPresent)
	}

	return field
}

// resolveProperty follows a property to the schema that actually describes its
// value: through `items` for arrays, and through `$ref` for the GOBL types
// that back most fields. A $ref can itself point at an array schema
// (tax/addon-list, for one), so that is unwrapped too.
func (b *builder) resolveProperty(property *definition, ownerDocID string) (bool, *definition, *target) {
	var items *definition
	if property.Type == TypeArray {
		items = property.Items
	}

	source := property
	if items != nil {
		source = items
	}

	var resolved *target
	if source.Ref != "" {
		resolved = b.resolveRef(source.Ref, ownerDocID)
	}

	if resolved != nil && resolved.def.Type == TypeArray && resolved.def.Items != nil {
		nested := resolved.def.Items
		var nestedTarget *target
		if nested.Ref != "" {
			nestedTarget = b.resolveRef(nested.Ref, resolved.docID)
		}
		if nestedTarget != nil {
			return true, nestedTarget.def, nestedTarget
		}
		return true, nested, nil
	}

	if resolved != nil {
		return items != nil, resolved.def, resolved
	}
	return items != nil, source, nil
}

func (b *builder) resolveRef(ref, ownerDocID string) *target {
	docRef, fragment, _ := strings.Cut(ref, "#")

	docID := ownerDocID
	if docRef != "" {
		docID = SchemaID(docRef)
	}

	doc, err := b.document(docID)
	if err != nil {
		return nil
	}

	if fragment == "" {
		return b.rootTarget(doc, docID)
	}

	name := strings.TrimPrefix(fragment, "/$defs/")
	def, ok := doc.Defs[name]
	if !ok {
		return nil
	}

	return &target{def: def, docID: docID, key: docID + "#/$defs/" + name}
}

// rootTarget resolves the definition a document's own $ref points at, which is
// where GOBL keeps the document's properties.
func (b *builder) rootTarget(doc *definition, docID string) *target {
	if doc.Ref == "" {
		return &target{def: doc, docID: docID, key: docID}
	}
	return b.resolveRef(doc.Ref, docID)
}

func (b *builder) document(id string) (*definition, error) {
	if doc, ok := b.registry[id]; ok {
		if doc == nil {
			return nil, fmt.Errorf("goblfields: %s: schema not available", SchemaPath(id))
		}
		return doc, nil
	}

	doc, err := readDocument(b.source, id)
	if err != nil {
		b.registry[id] = nil
		return nil, err
	}

	b.registry[id] = doc
	return doc, nil
}

// typeLabel names the JSON Schema type of a field's value. An array is just
// "array": the element type would repeat what the path's [] already says for
// all but a handful of fields (a bill/invoice has 192 arrays, and 182 of them
// hold objects), and what is actually worth knowing — whether there is
// anything underneath — the row shows with its drill-in chevron.
func typeLabel(node *definition, isArray bool) string {
	if isArray {
		return TypeArray
	}
	if node.Properties.len() > 0 || node.Type == "" {
		return TypeObject
	}
	return node.Type
}

func firstString(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func firstOption(opts []Options) Options {
	for _, o := range opts {
		return o
	}
	return Options{}
}

func source(o Options) fs.FS {
	if o.Source != nil {
		return o.Source
	}
	return Schemas
}

func maxDepth(o Options) int {
	if o.MaxDepth > 0 {
		return o.MaxDepth
	}
	return DefaultMaxDepth
}

func mustSub(source fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(source, dir)
	if err != nil {
		panic(err)
	}
	return sub
}
