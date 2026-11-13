package gobltemplate_test

import (
	"encoding/json"
	"strings"
	"testing"
	"text/template"

	"github.com/invopop/popui.go/goblfields"
	"github.com/invopop/popui.go/gobltemplate"
	"github.com/invopop/popui.go/props"
)

const (
	jobID        = "job-1"
	supplierName = "Provide One"
	fullCode     = "full_code"
	lineCodes    = "line_codes"
	taxID        = "tax_id"
)

const invoiceJSON = `{
	"type": "standard",
	"series": "F2026",
	"code": "0001",
	"issue_date": "2026-11-03",
	"currency": "EUR",
	"preceding": [{"series": "F2025", "code": "0099"}, {"code": "0100"}],
	"supplier": {"name": "Provide One", "emails": [{"addr": ""}, {"addr": "billing@example.com"}]},
	"customer": {"name": "Sample Consumer"},
	"payment": {"terms": {"due_dates": [{"date": "2026-12-01", "percent": "50%"}, {"date": "2027-01-01", "percent": "50%"}]}}
}`

const statusJSON = `{
	"type": "response",
	"code": "ST-1",
	"issue_date": "2026-11-03",
	"supplier": {"name": "Provide One"},
	"lines": [
		{"index": 1, "key": "accepted", "doc": {"series": "F2026", "code": "0001"}},
		{"index": 2, "key": "rejected", "doc": {"code": "0002"}, "reasons": [{"key": "references", "description": "Unknown buyer"}, {"key": "other"}]},
		{"index": 3, "key": "error"}
	]
}`

const partyJSON = `{"name": "Provide One", "alias": "P1", "tax_id": {"country": "ES", "code": "B98602642"}, "emails": [{"addr": "one@example.com"}]}`

func data(t *testing.T, schema, doc string, job gobltemplate.Job) map[string]any {
	t.Helper()
	d, err := gobltemplate.Data(schema, json.RawMessage(doc), job)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func expect(t *testing.T, d map[string]any, want map[string]string) {
	t.Helper()
	for key, value := range want {
		if d[key] != value {
			t.Errorf("%s = %q, want %q", key, d[key], value)
		}
	}
}

func TestDataInvoice(t *testing.T) {
	d := data(t, "https://gobl.org/draft-0/bill/invoice", invoiceJSON, gobltemplate.Job{
		ID:     jobID,
		Faults: []gobltemplate.Fault{{Code: "E1", Message: "first"}, {Message: "second"}},
	})
	expect(t, d, map[string]string{
		"type":            "standard",
		fullCode:          "F2026-0001",
		"issue_date":      "2026-11-03",
		"supplier_name":   supplierName,
		"supplier_email":  "billing@example.com",
		"customer_name":   "Sample Consumer",
		"customer_email":  "",
		"preceding_codes": "F2025-0099, 0100",
		"due_dates":       "2026-12-01, 2027-01-01",
		"job_id":          jobID,
		"silo_entry_id":   "",
		"faults":          "E1: first\nsecond\n",
	})
	for _, key := range []string{lineCodes, taxID, "total"} {
		if _, ok := d[key]; ok {
			t.Errorf("an invoice has %s", key)
		}
	}
	doc, _ := d["doc"].(map[string]any)
	if doc["code"] != "0001" {
		t.Errorf("doc = %v", d["doc"])
	}
}

func TestDataStatus(t *testing.T) {
	d := data(t, "bill/status", statusJSON, gobltemplate.Job{})
	expect(t, d, map[string]string{
		fullCode:        "ST-1",
		lineCodes:       "F2026-0001, 0002",
		"line_statuses": "accepted, rejected, error",
		"line_summary":  "F2026-0001: accepted, 0002: rejected, error",
		"line_reasons":  "Unknown buyer, other",
	})

	d = data(t, "bill/status", `{"code": "ST-2"}`, gobltemplate.Job{})
	expect(t, d, map[string]string{lineCodes: "", "line_summary": ""})
}

func TestDataParty(t *testing.T) {
	d := data(t, "org/party", partyJSON, gobltemplate.Job{})
	expect(t, d, map[string]string{
		"name":  supplierName,
		"alias": "P1",
		taxID:   "ESB98602642",
		"email": "one@example.com",
	})

	d = data(t, "org/party", `{"name": "No Tax ID"}`, gobltemplate.Job{})
	expect(t, d, map[string]string{taxID: "", "email": ""})
}

func TestDataTypedDocument(t *testing.T) {
	type invoice struct {
		Series string `json:"series"`
		Code   string `json:"code"`
	}
	d, err := gobltemplate.Data("bill/invoice", &invoice{Series: "F", Code: "1"}, gobltemplate.Job{})
	if err != nil {
		t.Fatal(err)
	}
	expect(t, d, map[string]string{fullCode: "F-1"})
}

func TestDataUnknownSchema(t *testing.T) {
	d := data(t, "bill/unknown", `{"code": "X"}`, gobltemplate.Job{ID: jobID})
	var keys []string
	for key := range d {
		keys = append(keys, key)
	}
	if len(d) != 5 || d["job_id"] != jobID || d["doc"] == nil {
		t.Errorf("keys = %v", keys)
	}
}

func TestDataNoDocument(t *testing.T) {
	d, err := gobltemplate.Data("bill/invoice", nil, gobltemplate.Job{})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := d["doc"]; ok {
		t.Error("doc without a document")
	}
	expect(t, d, map[string]string{fullCode: ""})
}

func TestRender(t *testing.T) {
	d := data(t, "bill/invoice", invoiceJSON, gobltemplate.Job{})
	tests := []struct{ name, text, want string }{
		{"variable", "{{.full_code}}", "F2026-0001"},
		{"document", "{{.doc.supplier.name}}", supplierName},
		{"missing value", "[{{.doc.customer.tax_id.code}}]", "[]"},
		{"missing variable", "[{{.line_codes}}]", "[]"},
		{"inside a block", "{{if .doc.notes}}notes{{else}}[{{.doc.ordering.code}}]{{end}}", "[]"},
		{"range", "{{range $i, $r := .doc.preceding}}{{if $i}}, {{end}}{{$r.code}}{{end}}", "0099, 0100"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gobltemplate.Render("body", tt.text, d)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
	if d["doc"].(map[string]any)["supplier"].(map[string]any)["name"] != supplierName {
		t.Error("Render changed a value")
	}
}

func TestRenderErrors(t *testing.T) {
	for _, text := range []string{
		"{{.doc.lines[].doc.code}}",
		"{{.doc.$regime}}",
		"{{.full_code",
	} {
		if _, err := gobltemplate.Render("body", text, map[string]any{}); err == nil {
			t.Errorf("%s: no error", text)
		}
	}
	if _, err := gobltemplate.Render("body", "{{index .doc.lines 0}}", map[string]any{"doc": map[string]any{"lines": []any{}}}); err == nil {
		t.Error("no error for an index out of range")
	}
}

func TestLookup(t *testing.T) {
	if s, ok := gobltemplate.Lookup("https://gobl.org/draft-0/bill/status"); !ok || s.Label != "Status" {
		t.Errorf("got %v %v", s, ok)
	}
	if _, ok := gobltemplate.Lookup("bill/unknown"); ok {
		t.Error("found bill/unknown")
	}
}

func TestFieldsAnyGroup(t *testing.T) {
	fields := gobltemplate.Fields(gobltemplate.Invoice, gobltemplate.Status, gobltemplate.Party)
	if len(fields) != 4 || fields[0].Name != gobltemplate.AnyGroup || !fields[0].Group {
		t.Fatalf("groups = %v", names(fields))
	}
	shared := map[string]string{}
	for _, f := range fields[0].Children {
		shared[f.Name] = f.Description
	}
	if shared["job_id"] != "All schemas" {
		t.Errorf("job_id = %q", shared["job_id"])
	}
	if shared[fullCode] != "Invoice, Status" {
		t.Errorf("full_code = %q", shared[fullCode])
	}
	if _, ok := shared[taxID]; ok {
		t.Error("tax_id is under Any")
	}

	for i, s := range []gobltemplate.Schema{gobltemplate.Invoice, gobltemplate.Status, gobltemplate.Party} {
		group := fields[i+1]
		if group.Name != s.Label || group.Description != s.ID || !group.Group {
			t.Errorf("group %d = %s %s", i+1, group.Name, group.Description)
		}
		keys := map[string]bool{}
		for _, f := range group.Children {
			keys[f.Name] = true
		}
		for _, v := range gobltemplate.Variables(s.ID) {
			if !keys[v.Key] {
				t.Errorf("%s has no %s", s.Label, v.Key)
			}
		}
	}
}

// TestFieldsEveryPickParses checks every path that the picker can insert.
// The gobl release of this module has no bill/status schema, so that group
// has no document tree here.
func TestFieldsEveryPickParses(t *testing.T) {
	fields := gobltemplate.Fields(gobltemplate.Schemas()...)
	paths := map[string]bool{}
	for _, group := range fields {
		for _, f := range goblfields.Flatten(group.Children) {
			paths[f.Path] = true
			if f.Type == goblfields.TypeObject {
				continue
			}
			if _, err := template.New("").Parse(strings.Replace(gobltemplate.Format, "%s", f.Path, 1)); err != nil {
				t.Errorf("%s: %v", f.Path, err)
			}
		}
	}
	for _, path := range []string{"doc.supplier.name", "doc.tax_id.code", lineCodes} {
		if !paths[path] {
			t.Errorf("no %s", path)
		}
	}
}

// TestFieldsMatchData checks that every variable in the picker has a value
// in the data.
func TestFieldsMatchData(t *testing.T) {
	for _, s := range gobltemplate.Schemas() {
		d := data(t, s.ID, `{}`, gobltemplate.Job{})
		for _, f := range gobltemplate.Fields(s)[1].Children {
			if f.Name == "doc" {
				continue
			}
			if _, ok := d[f.Path]; !ok {
				t.Errorf("%s: no data for %s", s.ID, f.Path)
			}
		}
	}
}

func names(fields []props.Field) []string {
	out := make([]string, len(fields))
	for i, f := range fields {
		out[i] = f.Name
	}
	return out
}
