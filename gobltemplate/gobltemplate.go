// Package gobltemplate is what an app needs to let its users write Go
// templates about GOBL documents, such as an email body or a Slack message:
//
//   - Variables lists the values a template prints by key: {{.full_code}}.
//   - Data builds what a template executes against: the variables, and the
//     document itself under doc, as in {{.doc.supplier.name}}.
//   - Render executes a template.
//   - Fields builds the FieldPicker tree that inserts both.
//
// A template reads one shape of data in every app, so a template that works
// in one app works in the others:
//
//	data, err := gobltemplate.Data(env.Extract(), gobltemplate.Job{ID: task.JobId})
//	if err != nil {
//		return err
//	}
//	text, err := gobltemplate.Render("body", body, data)
package gobltemplate

import (
	"strings"

	"github.com/a-h/templ"
	"github.com/invopop/icons"
)

// Schema is a kind of GOBL document that a template can be about.
type Schema struct {
	// ID is the short path of the schema, such as bill/invoice.
	ID string
	// Label names the schema's group in the picker.
	Label string
	// Icon is drawn before the schema's group in the picker.
	Icon func() templ.Component
}

// The schemas that have variables of their own.
var (
	Invoice  = Schema{ID: "bill/invoice", Label: "Invoice", Icon: icons.Invoice}
	Status   = Schema{ID: "bill/status", Label: "Status", Icon: icons.Status}
	Payment  = Schema{ID: "bill/payment", Label: "Payment", Icon: icons.Payment}
	Order    = Schema{ID: "bill/order", Label: "Order", Icon: icons.Order}
	Delivery = Schema{ID: "bill/delivery", Label: "Delivery", Icon: icons.Van}
	Party    = Schema{ID: "org/party", Label: "Party", Icon: icons.Team}
	Item     = Schema{ID: "org/item", Label: "Item", Icon: icons.Box}
)

// billSchemas are the documents with a type, a code, a supplier and a
// customer.
var billSchemas = []Schema{Invoice, Status, Payment, Order, Delivery}

// Schemas returns every schema that has variables of its own.
func Schemas() []Schema {
	return []Schema{Invoice, Status, Payment, Order, Delivery, Party, Item}
}

// Lookup finds a schema by its short path or its full id, such as
// https://gobl.org/draft-0/bill/invoice.
func Lookup(id string) (Schema, bool) {
	for _, s := range Schemas() {
		if s.matches(id) {
			return s, true
		}
	}
	return Schema{}, false
}

func (s Schema) matches(id string) bool {
	return id == s.ID || strings.HasSuffix(id, "/"+s.ID)
}
