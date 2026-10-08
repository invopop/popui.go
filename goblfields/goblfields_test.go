package goblfields_test

import (
	"strings"
	"testing"
	"testing/fstest"

	"github.com/invopop/popui.go/goblfields"
	"github.com/invopop/popui.go/props"
)

// testSchemas is a cut-down stand-in for the real GOBL documents, keeping the
// shapes that matter: refs between documents, arrays of refs, local $defs, a
// schema that is itself an array, and a definition that contains itself.
var testSchemas = fstest.MapFS{
	"bill/invoice.json": &fstest.MapFile{Data: []byte(`{
		"$id": "https://gobl.org/draft-0/bill/invoice",
		"$ref": "#/$defs/bill.Invoice",
		"$defs": {
			"bill.Invoice": {
				"type": "object",
				"required": ["currency", "supplier", "lines"],
				"properties": {
					"$addons": {"$ref": "https://gobl.org/draft-0/tax/addon-list", "title": "Addons"},
					"uuid": {"type": "string", "title": "UUID"},
					"currency": {"$ref": "https://gobl.org/draft-0/currency/code", "title": "Currency"},
					"supplier": {"$ref": "https://gobl.org/draft-0/org/party", "title": "Supplier"},
					"lines": {"type": "array", "items": {"$ref": "https://gobl.org/draft-0/bill/line"}, "title": "Lines"},
					"complements": {"type": "array", "items": {"$ref": "https://gobl.org/draft-0/schema/object"}, "title": "Complements"}
				}
			}
		}
	}`)},
	"currency/code.json": &fstest.MapFile{Data: []byte(`{
		"$ref": "#/$defs/currency.Code",
		"$defs": {"currency.Code": {"type": "string"}}
	}`)},
	"tax/addon-list.json": &fstest.MapFile{Data: []byte(`{
		"$ref": "#/$defs/tax.Addons",
		"$defs": {"tax.Addons": {"type": "array", "items": {"$ref": "https://gobl.org/draft-0/cbc/key"}}}
	}`)},
	"cbc/key.json": &fstest.MapFile{Data: []byte(`{
		"$ref": "#/$defs/cbc.Key",
		"$defs": {"cbc.Key": {"type": "string"}}
	}`)},
	"org/party.json": &fstest.MapFile{Data: []byte(`{
		"$ref": "#/$defs/org.Party",
		"$defs": {
			"org.Party": {
				"type": "object",
				"required": ["name"],
				"properties": {
					"name": {"type": "string", "title": "Name"},
					"people": {"type": "array", "items": {"$ref": "https://gobl.org/draft-0/org/person"}, "title": "People"}
				}
			}
		}
	}`)},
	"org/person.json": &fstest.MapFile{Data: []byte(`{
		"$ref": "#/$defs/org.Person",
		"$defs": {"org.Person": {"type": "object", "properties": {"name": {"$ref": "https://gobl.org/draft-0/org/name"}}}}
	}`)},
	"org/name.json": &fstest.MapFile{Data: []byte(`{
		"$ref": "#/$defs/org.Name",
		"$defs": {"org.Name": {"type": "object", "properties": {"given": {"type": "string"}}}}
	}`)},
	"bill/line.json": &fstest.MapFile{Data: []byte(`{
		"$ref": "#/$defs/bill.Line",
		"$defs": {
			"bill.Line": {
				"type": "object",
				"required": ["item"],
				"properties": {
					"item": {"$ref": "https://gobl.org/draft-0/org/item", "title": "Item"},
					"breakdown": {"type": "array", "items": {"$ref": "#/$defs/bill.SubLine"}}
				}
			},
			"bill.SubLine": {
				"type": "object",
				"required": ["item"],
				"properties": {
					"item": {"$ref": "https://gobl.org/draft-0/org/item"},
					"breakdown": {"type": "array", "items": {"$ref": "#/$defs/bill.SubLine"}}
				}
			}
		}
	}`)},
	"envelope.json": &fstest.MapFile{Data: []byte(`{
		"$ref": "#/$defs/gobl.Envelope",
		"$defs": {
			"gobl.Envelope": {
				"type": "object",
				"required": ["$schema", "head", "doc"],
				"properties": {
					"$schema": {"type": "string", "title": "JSON Schema ID"},
					"head": {"$ref": "https://gobl.org/draft-0/head/header", "title": "Header"},
					"doc": {"$ref": "https://gobl.org/draft-0/schema/object", "title": "Document"}
				}
			}
		}
	}`)},
	"head/header.json": &fstest.MapFile{Data: []byte(`{
		"$ref": "#/$defs/head.Header",
		"$defs": {"head.Header": {"type": "object", "required": ["uuid"], "properties": {"uuid": {"type": "string", "title": "UUID"}}}}
	}`)},
	"schema/object.json": &fstest.MapFile{Data: []byte(`{
		"$ref": "#/$defs/schema.Object",
		"$defs": {"schema.Object": {"type": "object", "title": "Object"}}
	}`)},
	"org/item.json": &fstest.MapFile{Data: []byte(`{
		"$ref": "#/$defs/org.Item",
		"$defs": {
			"org.Item": {
				"type": "object",
				"required": ["name"],
				"properties": {
					"name": {"type": "string", "title": "Name"},
					"price": {"type": "string", "title": "Price"}
				}
			}
		}
	}`)},
}

var testOptions = goblfields.Options{Source: testSchemas}

// Paths several tests meet: an array of schema/object, and a required field
// three levels down inside an array.
const (
	invoiceSchema   = "bill/invoice"
	addonsPath      = "$addons[]"
	complementsPath = "complements[]"
	linesPath       = "lines[]"
	lineItemName    = "lines[].item.name"
	uuidPath        = "uuid"
	// The one scalar label the tests check for; the builder passes it
	// through from the schema rather than naming it itself.
	typeString = "string"
)

// build is the invoice tree every test works from.
func build(t *testing.T, schema string, opts ...goblfields.Options) []props.Field {
	t.Helper()
	if len(opts) == 0 {
		opts = []goblfields.Options{testOptions}
	}
	fields, err := goblfields.Build(schema, opts...)
	if err != nil {
		t.Fatalf("building %s: %v", schema, err)
	}
	return fields
}

func paths(fields []props.Field) []string {
	out := make([]string, len(fields))
	for i, f := range fields {
		out[i] = f.Path
	}
	return out
}

// field returns the flattened tree entry at a path, failing when it is absent.
func field(t *testing.T, fields []props.Field, path string) props.Field {
	t.Helper()
	for _, f := range goblfields.Flatten(fields) {
		if f.Path == path {
			return f
		}
	}
	t.Fatalf("no field at %s", path)
	return props.Field{}
}

func has(fields []props.Field, path string) bool {
	for _, f := range goblfields.Flatten(fields) {
		if f.Path == path {
			return true
		}
	}
	return false
}

func equal(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func TestBuildAcceptsShortPathOrFullID(t *testing.T) {
	byID := build(t, "https://gobl.org/draft-0/"+invoiceSchema)
	equal(t, paths(byID), paths(build(t, invoiceSchema)))
}

func TestBuildFailsOnAnUnknownSchema(t *testing.T) {
	if _, err := goblfields.Build("bill/unknown", testOptions); err == nil {
		t.Fatal("expected an error for a schema that is not there")
	}
}

func TestBuildListsTheRootProperties(t *testing.T) {
	equal(t, paths(build(t, invoiceSchema)), []string{
		addonsPath, uuidPath, "currency", "supplier", linesPath, complementsPath,
	})
}

func TestBuildMarksArraysInThePathAndType(t *testing.T) {
	fields := build(t, invoiceSchema)
	if lines := field(t, fields, linesPath); !lines.Array || lines.Type != goblfields.TypeArray {
		t.Fatalf("lines[]: array=%v type=%q", lines.Array, lines.Type)
	}
	if got := field(t, fields, uuidPath).Type; got != typeString {
		t.Fatalf("uuid type = %q", got)
	}
}

func TestBuildFollowsARefThatPointsAtAnArraySchema(t *testing.T) {
	// The property is a plain $ref, so it only comes out as an array at all
	// because the document it points at is itself one.
	addons := field(t, build(t, invoiceSchema), addonsPath)
	if !addons.Array || addons.Type != goblfields.TypeArray {
		t.Fatalf("$addons[]: array=%v type=%q", addons.Array, addons.Type)
	}
}

func TestBuildResolvesRefsIntoOtherDocuments(t *testing.T) {
	fields := build(t, invoiceSchema)
	if got := field(t, fields, "supplier").Type; got != "object" {
		t.Fatalf("supplier type = %q", got)
	}
	// supplier.name only exists at all if org/party was resolved.
	if got := field(t, fields, "supplier.name").Type; got != typeString {
		t.Fatalf("supplier.name type = %q", got)
	}
	if got := field(t, fields, "currency").Type; got != typeString {
		t.Fatalf("currency type = %q", got)
	}
}

func TestBuildMarksAlwaysPresentOnlyWhenEveryStepIsRequired(t *testing.T) {
	fields := build(t, invoiceSchema)
	for path, want := range map[string]bool{
		"currency":                     true,
		uuidPath:                       false,
		"supplier.name":                true,
		lineItemName:                   true,
		"lines[].item.price":           false, // price is optional on the item,
		"supplier.people[].name.given": false, // and people optional on the supplier
	} {
		if got := field(t, fields, path).AlwaysPresent; got != want {
			t.Errorf("%s: always present = %v, want %v", path, got, want)
		}
	}
}

func TestBuildKeepsRequiredSeparateFromTheInheritedAnswer(t *testing.T) {
	fields := build(t, invoiceSchema)
	for path, want := range map[string]bool{
		"supplier.people[].name.given": false,
		"lines[].item.price":           false,
		lineItemName:                   true,
	} {
		if got := field(t, fields, path).Required; got != want {
			t.Errorf("%s: required = %v, want %v", path, got, want)
		}
	}
}

func TestBuildStopsABranchThatRepeatsASchemaItAlreadyContains(t *testing.T) {
	fields := build(t, invoiceSchema)
	if got := field(t, fields, "lines[].breakdown[]"); len(got.Children) == 0 {
		t.Fatal("lines[].breakdown[] should expand once")
	}
	if got := field(t, fields, "lines[].breakdown[].breakdown[]"); len(got.Children) != 0 {
		t.Fatalf("lines[].breakdown[].breakdown[] should be a leaf, got %v", paths(got.Children))
	}
}

func TestBuildStopsAtTheRequestedDepth(t *testing.T) {
	shallow := build(t, invoiceSchema, goblfields.Options{Source: testSchemas, MaxDepth: 2})
	if !has(shallow, "supplier.people[]") {
		t.Fatal("supplier.people[] should be within two levels")
	}
	if has(shallow, "supplier.people[].name") {
		t.Fatal("supplier.people[].name is three levels down")
	}
	for _, f := range goblfields.Flatten(shallow) {
		if len(f.Children) > 0 && len(f.Children[0].Children) > 0 {
			t.Fatalf("%s expands past the depth cap", f.Path)
		}
	}
}

// GOBL types a value whose schema is only known at runtime as schema/object,
// which carries no properties. There is nothing to drill into, so it is a leaf
// wherever it turns up: the envelope's doc, an invoice complement.
func TestSchemaObjectIsALeaf(t *testing.T) {
	envelope := build(t, "envelope")
	if !has(envelope, "doc") {
		t.Fatal("the envelope should have a doc")
	}
	for _, f := range goblfields.Flatten(envelope) {
		if len(f.Path) > 4 && f.Path[:4] == "doc." {
			t.Fatalf("doc should be a leaf, found %s", f.Path)
		}
	}
}

func TestSchemaObjectDoesNotStopTheRestOfTheSchemaExpanding(t *testing.T) {
	if !has(build(t, "envelope"), "head.uuid") {
		t.Fatal("head should still expand")
	}
}

func TestSchemaObjectIsALeafInsideAnArrayToo(t *testing.T) {
	fields := build(t, invoiceSchema)
	if !has(fields, complementsPath) {
		t.Fatal("complements[] should be there")
	}
	for _, f := range goblfields.Flatten(fields) {
		if strings.HasPrefix(f.Path, complementsPath+".") {
			t.Fatalf("complements[] should be a leaf, found %s", f.Path)
		}
	}
}

func TestTrailReturnsTheChainDownToTheField(t *testing.T) {
	trail := goblfields.Trail(build(t, invoiceSchema), lineItemName)
	equal(t, paths(trail), []string{linesPath, "lines[].item", lineItemName})
}

func TestTrailFallsBackToTheDeepestPartOfThePathThatExists(t *testing.T) {
	trail := goblfields.Trail(build(t, invoiceSchema), "lines[].item.nope")
	equal(t, paths(trail), []string{linesPath, "lines[].item"})
}

func TestTrailHasNothingToReturnForAnEmptyPath(t *testing.T) {
	if trail := goblfields.Trail(build(t, invoiceSchema), ""); len(trail) != 0 {
		t.Fatalf("got %v", paths(trail))
	}
}

func TestListReturnsTheAvailableSchemas(t *testing.T) {
	list, err := goblfields.List(testOptions)
	if err != nil {
		t.Fatal(err)
	}
	equal(t, list[:3], []string{invoiceSchema, "bill/line", "cbc/key"})
}

// The embedded GOBL schemas are the ones consumers actually build from, so the
// tree is checked against them too: property order, depth and the recursion
// guard all depend on the real documents.
func TestBuildFromTheEmbeddedSchemas(t *testing.T) {
	fields, err := goblfields.Build(invoiceSchema)
	if err != nil {
		t.Fatal(err)
	}

	equal(t, paths(fields)[:6], []string{"$regime", "$addons[]", "$tags[]", uuidPath, "type", "series"})

	// org/party has no `required` of its own, so nothing under the supplier
	// is guaranteed even though the supplier itself is.
	if got := field(t, fields, "supplier"); !got.AlwaysPresent {
		t.Error("an invoice always has a supplier")
	}
	if got := field(t, fields, "supplier.name"); got.AlwaysPresent {
		t.Error("org/party requires nothing, so a supplier name is not guaranteed")
	}
	// issue_time carries GOBL's `calculated` hint but is not required.
	if got := field(t, fields, "issue_time"); got.AlwaysPresent {
		t.Error("issue_time is calculated, not always present")
	}
	if got := field(t, fields, linesPath); got.Type != goblfields.TypeArray || !got.Array {
		t.Errorf("lines[]: type=%q array=%v", got.Type, got.Array)
	}
}

func TestNestPutsATreeUnderOneField(t *testing.T) {
	doc := goblfields.Nest("doc", build(t, invoiceSchema))
	if doc.Path != "doc" || doc.Type != goblfields.TypeObject || !doc.AlwaysPresent {
		t.Fatalf("root = %+v", doc)
	}
	equal(t, paths(doc.Children)[:3], []string{"doc.$addons[]", "doc.uuid", "doc.currency"})
	// Every path beneath is rewritten, at any depth, and nothing else changes.
	if got := field(t, []props.Field{doc}, "doc."+lineItemName); !got.AlwaysPresent || got.Name != "name" {
		t.Fatalf("nested leaf = %+v", got)
	}
	// The input is left alone.
	if got := paths(build(t, invoiceSchema))[0]; got != addonsPath {
		t.Fatalf("original tree changed: %s", got)
	}
}
