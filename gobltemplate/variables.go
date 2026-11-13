package gobltemplate

import (
	"strings"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
)

// Variable is a value that a template prints by its key, such as
// {{.full_code}}. Saved templates use the keys, so a key is never renamed or
// removed.
type Variable struct {
	Key         string `json:"key"`
	Description string `json:"description"`

	// schemas carry the variable. Nil means every document, of any schema.
	schemas []Schema
	// always marks a value that every job has.
	always bool
	value  func(doc any, job Job) string
}

// variables is the table of every variable, in the order the picker lists
// them.
var variables = []Variable{
	{Key: "type", Description: "Type", schemas: billSchemas, value: fromHeader(func(h header) string {
		return h.Type.String()
	})},
	{Key: "series", Description: "Series", schemas: billSchemas, value: fromHeader(func(h header) string {
		return h.Series.String()
	})},
	{Key: "code", Description: "Code", schemas: billSchemas, value: fromHeader(func(h header) string {
		return h.Code.String()
	})},
	{Key: "full_code", Description: "Full code (series-code)", schemas: billSchemas, value: fromHeader(func(h header) string {
		return h.Series.Join(h.Code).String()
	})},
	{Key: "issue_date", Description: "Issue date", schemas: billSchemas, value: fromHeader(func(h header) string {
		return h.IssueDate.String()
	})},
	{Key: "supplier_name", Description: "Supplier name", schemas: billSchemas, value: fromHeader(func(h header) string {
		return partyName(h.Supplier)
	})},
	{Key: "supplier_email", Description: "Supplier email", schemas: billSchemas, value: fromHeader(func(h header) string {
		return firstEmail(h.Supplier)
	})},
	{Key: "customer_name", Description: "Customer name", schemas: billSchemas, value: fromHeader(func(h header) string {
		return partyName(h.Customer)
	})},
	{Key: "customer_email", Description: "Customer email", schemas: billSchemas, value: fromHeader(func(h header) string {
		return firstEmail(h.Customer)
	})},

	{Key: "total", Description: "Total", schemas: []Schema{Payment}, value: from(func(p *bill.Payment) string {
		return p.Total.String()
	})},
	{Key: "currency", Description: "Currency", schemas: []Schema{Payment, Item}, value: firstOf(
		from(func(p *bill.Payment) string { return p.Currency.String() }),
		from(func(i *org.Item) string { return i.Currency.String() }),
	)},
	{Key: "preceding_codes", Description: "Preceding documents", schemas: []Schema{Invoice, Payment, Order, Delivery}, value: fromHeader(func(h header) string {
		codes := make([]string, 0, len(h.Preceding))
		for _, ref := range h.Preceding {
			if ref != nil {
				codes = append(codes, ref.Series.Join(ref.Code).String())
			}
		}
		return joinList(codes)
	})},
	{Key: "due_dates", Description: "Due dates", schemas: []Schema{Invoice}, value: from(func(inv *bill.Invoice) string {
		if inv.Payment == nil || inv.Payment.Terms == nil {
			return ""
		}
		var dates []string
		for _, due := range inv.Payment.Terms.DueDates {
			if due != nil && due.Date != nil {
				dates = append(dates, due.Date.String())
			}
		}
		return joinList(dates)
	})},

	{Key: "line_codes", Description: "Line documents", schemas: []Schema{Status}, value: from(func(s *bill.Status) string {
		var codes []string
		for _, line := range statusLines(s) {
			codes = append(codes, lineCode(line))
		}
		return joinList(codes)
	})},
	{Key: "line_statuses", Description: "Line statuses", schemas: []Schema{Status}, value: from(func(s *bill.Status) string {
		var keys []string
		for _, line := range statusLines(s) {
			keys = append(keys, line.Key.String())
		}
		return joinList(keys)
	})},
	{Key: "line_summary", Description: "Line summary", schemas: []Schema{Status}, value: from(func(s *bill.Status) string {
		var summary []string
		for _, line := range statusLines(s) {
			summary = append(summary, joinNonEmpty(": ", lineCode(line), line.Key.String()))
		}
		return joinList(summary)
	})},
	{Key: "line_reasons", Description: "Line reasons", schemas: []Schema{Status}, value: from(func(s *bill.Status) string {
		var reasons []string
		for _, line := range statusLines(s) {
			for _, reason := range line.Reasons {
				switch {
				case reason == nil:
				case reason.Description != "":
					reasons = append(reasons, reason.Description)
				default:
					reasons = append(reasons, reason.Key.String())
				}
			}
		}
		return joinList(reasons)
	})},

	{Key: "name", Description: "Name", schemas: []Schema{Party, Item}, value: firstOf(
		from(func(p *org.Party) string { return p.Name }),
		from(func(i *org.Item) string { return i.Name }),
	)},
	{Key: "tax_id", Description: "Tax ID", schemas: []Schema{Party}, value: from(func(p *org.Party) string {
		if p.TaxID == nil {
			return ""
		}
		return p.TaxID.String()
	})},
	{Key: "email", Description: "Email", schemas: []Schema{Party}, value: from(firstEmail)},
	{Key: "alias", Description: "Alias", schemas: []Schema{Party}, value: from(func(p *org.Party) string {
		return p.Alias
	})},
	{Key: "ref", Description: "Reference", schemas: []Schema{Item}, value: from(func(i *org.Item) string {
		return i.Ref.String()
	})},
	{Key: "description", Description: "Description", schemas: []Schema{Item}, value: from(func(i *org.Item) string {
		return i.Description
	})},
	{Key: "price", Description: "Price", schemas: []Schema{Item}, value: from(func(i *org.Item) string {
		if i.Price == nil {
			return ""
		}
		return i.Price.String()
	})},

	{Key: "job_id", Description: "Job ID", always: true, value: func(_ any, j Job) string { return j.ID }},
	{Key: "silo_entry_id", Description: "Silo entry ID", always: true, value: func(_ any, j Job) string { return j.SiloEntryID }},
	{Key: "faults", Description: "Faults", value: func(_ any, j Job) string { return faults(j.Faults) }},
	{Key: "meta_link", Description: "Link", value: func(_ any, j Job) string { return j.Link }},
}

// Variables returns the variables that a template about the schema can
// print, given as a short path or a full id. A schema with no variables of
// its own gets the job's variables only.
func Variables(schema string) []Variable {
	var out []Variable
	for _, v := range variables {
		if v.carriedBy(schema) {
			out = append(out, v)
		}
	}
	return out
}

func (v Variable) carriedBy(schema string) bool {
	if v.schemas == nil {
		return true
	}
	for _, s := range v.schemas {
		if s.matches(schema) {
			return true
		}
	}
	return false
}

// header holds the fields that every bill document has. A status has no
// preceding documents.
type header struct {
	Type      cbc.Key
	Series    cbc.Code
	Code      cbc.Code
	IssueDate cal.Date
	Supplier  *org.Party
	Customer  *org.Party
	Preceding []*org.DocumentRef
}

func headerOf(doc any) (header, bool) {
	switch d := doc.(type) {
	case *bill.Invoice:
		if d != nil {
			return header{Type: d.Type, Series: d.Series, Code: d.Code, IssueDate: d.IssueDate, Supplier: d.Supplier, Customer: d.Customer, Preceding: d.Preceding}, true
		}
	case *bill.Status:
		if d != nil {
			return header{Type: d.Type, Series: d.Series, Code: d.Code, IssueDate: d.IssueDate, Supplier: d.Supplier, Customer: d.Customer}, true
		}
	case *bill.Payment:
		if d != nil {
			return header{Type: d.Type, Series: d.Series, Code: d.Code, IssueDate: d.IssueDate, Supplier: d.Supplier, Customer: d.Customer, Preceding: d.Preceding}, true
		}
	case *bill.Order:
		if d != nil {
			return header{Type: d.Type, Series: d.Series, Code: d.Code, IssueDate: d.IssueDate, Supplier: d.Supplier, Customer: d.Customer, Preceding: d.Preceding}, true
		}
	case *bill.Delivery:
		if d != nil {
			return header{Type: d.Type, Series: d.Series, Code: d.Code, IssueDate: d.IssueDate, Supplier: d.Supplier, Customer: d.Customer, Preceding: d.Preceding}, true
		}
	}
	return header{}, false
}

// fromHeader reads a value from a bill document.
func fromHeader(fn func(header) string) func(any, Job) string {
	return func(doc any, _ Job) string {
		h, ok := headerOf(doc)
		if !ok {
			return ""
		}
		return fn(h)
	}
}

// from reads a value from a document of type T.
func from[T any](fn func(*T) string) func(any, Job) string {
	return func(doc any, _ Job) string {
		d, ok := doc.(*T)
		if !ok || d == nil {
			return ""
		}
		return fn(d)
	}
}

// firstOf returns the first value that is not empty.
func firstOf(fns ...func(any, Job) string) func(any, Job) string {
	return func(doc any, job Job) string {
		for _, fn := range fns {
			if v := fn(doc, job); v != "" {
				return v
			}
		}
		return ""
	}
}

func statusLines(s *bill.Status) []*bill.StatusLine {
	lines := make([]*bill.StatusLine, 0, len(s.Lines))
	for _, line := range s.Lines {
		if line != nil {
			lines = append(lines, line)
		}
	}
	return lines
}

func lineCode(line *bill.StatusLine) string {
	if line.Doc == nil {
		return ""
	}
	return line.Doc.Series.Join(line.Doc.Code).String()
}

func partyName(party *org.Party) string {
	if party == nil {
		return ""
	}
	return party.Name
}

// firstEmail returns a party's first email address.
func firstEmail(party *org.Party) string {
	if party == nil {
		return ""
	}
	for _, email := range party.Emails {
		if email != nil && email.Address != "" {
			return email.Address
		}
	}
	return ""
}

func faults(list []Fault) string {
	var b strings.Builder
	for _, f := range list {
		if f.Code != "" {
			b.WriteString(f.Code)
			b.WriteString(": ")
		}
		b.WriteString(f.Message)
		b.WriteString("\n")
	}
	return b.String()
}

// joinList joins the non-empty values of a list, in order.
func joinList(values []string) string {
	return joinNonEmpty(", ", values...)
}

func joinNonEmpty(sep string, values ...string) string {
	kept := make([]string, 0, len(values))
	for _, v := range values {
		if v != "" {
			kept = append(kept, v)
		}
	}
	return strings.Join(kept, sep)
}
