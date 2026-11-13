package gobltemplate_test

import (
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
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
	euro         = "EUR"
	docKey       = "doc"
	customerName = "Sample Consumer"
)

func supplier() *org.Party {
	return &org.Party{
		Name:   supplierName,
		Emails: []*org.Email{{Address: ""}, {Address: "billing@example.com"}},
	}
}

func customer() *org.Party {
	return &org.Party{Name: customerName, Emails: []*org.Email{{Address: "customer@example.com"}}}
}

func invoice() *bill.Invoice {
	return &bill.Invoice{
		Type:      bill.InvoiceTypeStandard,
		Series:    "F2026",
		Code:      "0001",
		IssueDate: cal.MakeDate(2026, time.November, 3),
		Preceding: []*org.DocumentRef{{Series: "F2025", Code: "0099"}, {Code: "0100"}},
		Supplier:  supplier(),
		Customer:  customer(),
		Payment: &bill.PaymentDetails{
			Terms: &pay.Terms{DueDates: []*pay.DueDate{
				{Date: cal.NewDate(2026, time.December, 1)},
				{Date: cal.NewDate(2027, time.January, 1)},
			}},
		},
	}
}

func status() *bill.Status {
	return &bill.Status{
		Type:      bill.StatusTypeResponse,
		Series:    "ST",
		Code:      "1",
		IssueDate: cal.MakeDate(2026, time.November, 3),
		Supplier:  supplier(),
		Customer:  customer(),
		Lines: []*bill.StatusLine{
			{Key: bill.StatusLineAccepted, Doc: &org.DocumentRef{Series: "F2026", Code: "0001"}},
			{
				Key: bill.StatusLineRejected,
				Doc: &org.DocumentRef{Code: "0002"},
				Reasons: []*bill.Reason{
					{Key: bill.ReasonKeyPrices, Description: "Wrong prices"},
					{Key: bill.ReasonKeyQuality},
				},
			},
			nil,
			{Key: bill.StatusLineRejected},
		},
	}
}

func party() *org.Party {
	return &org.Party{
		Name:   supplierName,
		Alias:  "P1",
		TaxID:  &tax.Identity{Country: l10n.ES.Tax(), Code: "B98602642"},
		Emails: []*org.Email{{Address: "one@example.com"}},
	}
}

func item() *org.Item {
	price := num.MakeAmount(1000, 2)
	return &org.Item{Ref: "W-1", Name: "Widget", Description: "A widget", Currency: euro, Price: &price}
}

func data(t *testing.T, doc any, job gobltemplate.Job) map[string]any {
	t.Helper()
	d, err := gobltemplate.Data(doc, job)
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
	d := data(t, invoice(), gobltemplate.Job{
		ID:     jobID,
		Faults: []gobltemplate.Fault{{Code: "E1", Message: "first"}, {Message: "second"}},
	})
	expect(t, d, map[string]string{
		"type":            "standard",
		fullCode:          "F2026-0001",
		"issue_date":      "2026-11-03",
		"supplier_name":   supplierName,
		"supplier_email":  "billing@example.com",
		"customer_name":   customerName,
		"customer_email":  "customer@example.com",
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
	doc, _ := d[docKey].(map[string]any)
	if doc["code"] != "0001" {
		t.Errorf("doc = %v", d[docKey])
	}
}

func TestDataStatus(t *testing.T) {
	d := data(t, status(), gobltemplate.Job{})
	expect(t, d, map[string]string{
		fullCode:        "ST-1",
		lineCodes:       "F2026-0001, 0002",
		"line_statuses": "accepted, rejected, rejected",
		"line_summary":  "F2026-0001: accepted, 0002: rejected, rejected",
		"line_reasons":  "Wrong prices, quality",
	})

	d = data(t, &bill.Status{Code: "ST-2"}, gobltemplate.Job{})
	expect(t, d, map[string]string{lineCodes: "", "line_summary": ""})
}

func TestDataParty(t *testing.T) {
	d := data(t, party(), gobltemplate.Job{})
	expect(t, d, map[string]string{
		"name":  supplierName,
		"alias": "P1",
		taxID:   "ESB98602642",
		"email": "one@example.com",
	})

	d = data(t, &org.Party{Name: "No Tax ID"}, gobltemplate.Job{})
	expect(t, d, map[string]string{taxID: "", "email": ""})
}

func TestDataItem(t *testing.T) {
	expect(t, data(t, item(), gobltemplate.Job{}), map[string]string{
		"name":     "Widget",
		"ref":      "W-1",
		"price":    "10.00",
		"currency": euro,
	})
}

func TestDataOtherDocument(t *testing.T) {
	d := data(t, &struct {
		Code string `json:"code"`
	}{Code: "X"}, gobltemplate.Job{ID: jobID})
	if len(d) != 5 || d["job_id"] != jobID || d[docKey] == nil {
		t.Errorf("data = %v", d)
	}
}

func TestDataNoDocument(t *testing.T) {
	for _, doc := range []any{nil, (*bill.Invoice)(nil)} {
		d := data(t, doc, gobltemplate.Job{})
		if _, ok := d[docKey]; ok {
			t.Errorf("%T: doc without a document", doc)
		}
	}
}

// TestEveryVariableHasAValue checks that each variable reads every schema it
// lists.
func TestEveryVariableHasAValue(t *testing.T) {
	total := num.MakeAmount(5000, 2)
	docs := map[string]any{
		gobltemplate.Invoice.ID: invoice(),
		gobltemplate.Status.ID:  status(),
		gobltemplate.Payment.ID: &bill.Payment{
			Type: bill.PaymentTypeReceipt, Series: "P", Code: "1", IssueDate: cal.MakeDate(2026, time.November, 3), Currency: euro, Total: total,
			Preceding: []*org.DocumentRef{{Code: "0001"}}, Supplier: supplier(), Customer: party(),
		},
		gobltemplate.Order.ID: &bill.Order{
			Type: bill.OrderTypePurchase, Series: "O", Code: "1", IssueDate: cal.MakeDate(2026, time.November, 3),
			Preceding: []*org.DocumentRef{{Code: "Q-1"}}, Supplier: supplier(), Customer: party(),
		},
		gobltemplate.Delivery.ID: &bill.Delivery{
			Type: bill.DeliveryTypeNote, Series: "D", Code: "1", IssueDate: cal.MakeDate(2026, time.November, 3),
			Preceding: []*org.DocumentRef{{Code: "O-1"}}, Supplier: supplier(), Customer: party(),
		},
		gobltemplate.Party.ID: party(),
		gobltemplate.Item.ID:  item(),
	}
	job := gobltemplate.Job{ID: jobID, SiloEntryID: "entry-1", Link: "https://example.com", Faults: []gobltemplate.Fault{{Message: "x"}}}
	for _, s := range gobltemplate.Schemas() {
		d := data(t, docs[s.ID], job)
		for _, v := range gobltemplate.Variables(s.ID) {
			if d[v.Key] == "" {
				t.Errorf("%s: %s is empty", s.ID, v.Key)
			}
		}
	}
}

func TestRender(t *testing.T) {
	d := data(t, invoice(), gobltemplate.Job{})
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
	if d[docKey].(map[string]any)["supplier"].(map[string]any)["name"] != supplierName {
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
	if _, err := gobltemplate.Render("body", "{{index .doc.lines 0}}", map[string]any{docKey: map[string]any{"lines": []any{}}}); err == nil {
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
		if last := group.Children[len(group.Children)-1]; last.Name != docKey {
			t.Errorf("%s ends in %s", s.Label, last.Name)
		}
	}
}

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
	for _, path := range []string{"doc.supplier.name", "doc.tax_id.code", "doc.price", lineCodes} {
		if !paths[path] {
			t.Errorf("no %s", path)
		}
	}
}

// TestFieldsMatchData checks that every variable in the picker has a value
// in the data.
func TestFieldsMatchData(t *testing.T) {
	docs := map[string]any{
		gobltemplate.Invoice.ID:  &bill.Invoice{},
		gobltemplate.Status.ID:   &bill.Status{},
		gobltemplate.Payment.ID:  &bill.Payment{},
		gobltemplate.Order.ID:    &bill.Order{},
		gobltemplate.Delivery.ID: &bill.Delivery{},
		gobltemplate.Party.ID:    &org.Party{},
		gobltemplate.Item.ID:     &org.Item{},
	}
	for _, s := range gobltemplate.Schemas() {
		d := data(t, docs[s.ID], gobltemplate.Job{})
		for _, f := range gobltemplate.Fields(s)[1].Children {
			if f.Name == docKey {
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
