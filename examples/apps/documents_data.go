package apps

import (
	"fmt"
	"strings"
)

// Document is an example GOBL document published in an app's documentation
// (the "Documents" tab of the docs page), e.g. a B2C simplified invoice for
// SII. There is no access model for this yet; see HANDOFF.md.
type Document struct {
	ID          string
	Title       string
	Description string
	// Schema is the GOBL schema of the example (bill/invoice, org/party…).
	Schema string
	// Channel is the app channel the example belongs to, empty for the
	// default channel.
	Channel string
	// Content is the minimal GOBL JSON shown to readers.
	Content string
}

// DocumentGroup is one deck on the Documents tab: the documents of one
// channel that share a schema, in display order.
type DocumentGroup struct {
	Channel   string
	Schema    string
	Documents []*Document
}

// ID identifies the group in routes ("<channel>/<schema>" with the default
// channel spelled "default" and the schema slash turned into a dash).
func (g *DocumentGroup) ID() string {
	ch := g.Channel
	if ch == "" {
		ch = "default"
	}
	return ch + "/" + strings.ReplaceAll(g.Schema, "/", "-")
}

// Reorder applies the order of document IDs the CardDeck submitted; IDs it
// does not mention keep their relative order at the end.
func (g *DocumentGroup) Reorder(ids []string) {
	seen := map[string]bool{}
	var out []*Document
	for _, id := range ids {
		for _, d := range g.Documents {
			if d.ID == id && !seen[id] {
				out = append(out, d)
				seen[id] = true
			}
		}
	}
	for _, d := range g.Documents {
		if !seen[d.ID] {
			out = append(out, d)
		}
	}
	g.Documents = out
}

// DocumentCollection is the ordered set of groups for one app.
type DocumentCollection struct {
	Groups []*DocumentGroup
}

// Count returns the number of documents across all groups.
func (dc *DocumentCollection) Count() int {
	n := 0
	for _, g := range dc.Groups {
		n += len(g.Documents)
	}
	return n
}

// Group finds a group by its route ID.
func (dc *DocumentCollection) Group(id string) *DocumentGroup {
	for _, g := range dc.Groups {
		if g.ID() == id {
			return g
		}
	}
	return nil
}

// ForChannel lists the groups of one channel in display order.
func (dc *DocumentCollection) ForChannel(channel string) []*DocumentGroup {
	var out []*DocumentGroup
	for _, g := range dc.Groups {
		if g.Channel == channel {
			out = append(out, g)
		}
	}
	return out
}

// Find locates a document and its group.
func (dc *DocumentCollection) Find(id string) (*DocumentGroup, *Document) {
	for _, g := range dc.Groups {
		for _, d := range g.Documents {
			if d.ID == id {
				return g, d
			}
		}
	}
	return nil, nil
}

// Add places a document in the group for its channel and schema, creating
// the group when it is the first of its kind.
func (dc *DocumentCollection) Add(d *Document) {
	for _, g := range dc.Groups {
		if g.Channel == d.Channel && g.Schema == d.Schema {
			g.Documents = append(g.Documents, d)
			return
		}
	}
	dc.Groups = append(dc.Groups, &DocumentGroup{Channel: d.Channel, Schema: d.Schema, Documents: []*Document{d}})
}

// Remove drops a document, and its group once empty.
func (dc *DocumentCollection) Remove(id string) {
	for gi, g := range dc.Groups {
		for i, d := range g.Documents {
			if d.ID != id {
				continue
			}
			g.Documents = append(g.Documents[:i], g.Documents[i+1:]...)
			if len(g.Documents) == 0 {
				dc.Groups = append(dc.Groups[:gi], dc.Groups[gi+1:]...)
			}
			return
		}
	}
}

// slugify turns a title into a key usable as an ID segment.
func slugify(s string) string {
	var b strings.Builder
	dash := false
	for _, r := range strings.ToLower(s) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			dash = false
		default:
			if !dash && b.Len() > 0 {
				b.WriteByte('-')
				dash = true
			}
		}
	}
	return strings.TrimSuffix(b.String(), "-")
}

// invoiceJSON is a compact minimal GOBL invoice used as example content.
func invoiceJSON(addon, docType, series string, simplified bool) string {
	tags := ""
	if simplified {
		tags = `
  "$tags": ["simplified"],`
	}
	return fmt.Sprintf(`{
  "$schema": "https://gobl.org/draft-0/bill/invoice",
  "$addons": ["%s"],%s
  "series": "%s",
  "code": "0001",
  "issue_date": "2026-09-01",
  "currency": "EUR",
  "supplier": {
    "name": "Provide One S.L.",
    "tax_id": { "country": "ES", "code": "B98602642" }
  },
  "customer": {
    "name": "Sample Consumer",
    "tax_id": { "country": "ES", "code": "54387763P" }
  },
  "lines": [
    {
      "quantity": "20",
      "item": { "name": "Development services", "price": "90.00" },
      "taxes": [{ "cat": "VAT", "rate": "standard" }]
    }
  ]%s
}`, addon, tags, series, docTypeExt(docType))
}

func docTypeExt(docType string) string {
	if docType == "" {
		return ""
	}
	return fmt.Sprintf(`,
  "tax": { "ext": { "es-sii-doc-type": "%s" } }`, docType)
}

const partyJSON = `{
  "$schema": "https://gobl.org/draft-0/org/party",
  "name": "Provide One S.L.",
  "tax_id": { "country": "ES", "code": "B98602642" },
  "addresses": [
    {
      "num": "42",
      "street": "Calle Pradillo",
      "locality": "Madrid",
      "region": "Madrid",
      "code": "28002",
      "country": "ES"
    }
  ],
  "emails": [{ "addr": "billing@provideone.com" }]
}`

// sampleDocuments mirrors the "Documents" tab of the Spain app page on
// docs.invopop.com (apps/spain.mdx), regrouped by channel: the shared party
// records and Facturae live in the default channel, the rest under their
// channel, with the SII corrections at the end of the SII invoices.
var sampleDocuments = map[string]*DocumentCollection{
	"gov-es": {Groups: []*DocumentGroup{
		{Channel: "", Schema: "org/party", Documents: []*Document{
			{ID: "supplier", Title: "Supplier", Schema: "org/party", Description: "A Spanish company with its NIF and the legal representative required for supplier registration.", Content: partyJSON},
			{ID: "customer", Title: "Customer", Schema: "org/party", Description: "A Spanish business customer identified by NIF with a full address.", Content: partyJSON},
		}},
		{Channel: "", Schema: "bill/invoice", Documents: []*Document{
			{ID: "fe-b2g", Title: "Facturae B2G Invoice", Schema: "bill/invoice", Description: "Invoice to a public body, converted to Facturae 3.2.x XML.", Content: invoiceJSON("es-facturae-v3", "", "FE", false)},
			{ID: "fe-credit", Title: "Facturae Credit Note", Schema: "bill/invoice", Description: "Credit note in Facturae format referencing the original invoice.", Content: invoiceJSON("es-facturae-v3", "", "FECR", false)},
			{ID: "fe-face", Title: "Facturae FACe Invoice (with Administrative Centers)", Schema: "bill/invoice", Description: "Includes the DIR3 administrative centre codes FACe requires.", Content: invoiceJSON("es-facturae-v3", "", "FACE", false)},
		}},
		{Channel: "ticketbai", Schema: "bill/invoice", Documents: []*Document{
			{ID: "tbai-b2b", Title: "B2B Invoice", Schema: "bill/invoice", Channel: "ticketbai", Description: "Standard invoice between two businesses, signed and registered with TicketBAI.", Content: invoiceJSON("es-tbai-v1", "", "SAMPLE", false)},
			{ID: "tbai-b2c", Title: "B2C Simplified Invoice", Schema: "bill/invoice", Channel: "ticketbai", Description: "Simplified invoice for a consumer, under the €400 threshold.", Content: invoiceJSON("es-tbai-v1", "", "SIMP", true)},
			{ID: "tbai-credit", Title: "Credit Note (Factura Rectificativa)", Schema: "bill/invoice", Channel: "ticketbai", Description: "Corrects a previously registered invoice by reference.", Content: invoiceJSON("es-tbai-v1", "", "CR", false)},
		}},
		{Channel: "sii", Schema: "bill/invoice", Documents: []*Document{
			{ID: "sii-b2c", Title: "B2C Simplified Invoice", Schema: "bill/invoice", Channel: "sii", Description: "Simplified invoice (F2) for a consumer; totals, regime and VAT rate are computed on build.", Content: invoiceJSON("es-sii-v1", "F2", "SIMP", true)},
			{ID: "sii-b2b", Title: "B2B Standard Invoice", Schema: "bill/invoice", Channel: "sii", Description: "Complete invoice (F1) between two Spanish businesses.", Content: invoiceJSON("es-sii-v1", "F1", "SAMPLE", false)},
			{ID: "sii-b2b-services-eu", Title: "B2B Services EU Client (Reverse Charge)", Schema: "bill/invoice", Channel: "sii", Description: "Services to an EU business; VAT is reverse charged.", Content: invoiceJSON("es-sii-v1", "F1", "EU", false)},
			{ID: "sii-b2b-goods-eu", Title: "B2B Goods EU Client (Intra-Community)", Schema: "bill/invoice", Channel: "sii", Description: "Intra-community supply of goods, exempt with the E5 code.", Content: invoiceJSON("es-sii-v1", "F1", "EUG", false)},
			{ID: "sii-b2b-services-no-eu", Title: "B2B Services Non-EU Client (Outside Scope)", Schema: "bill/invoice", Channel: "sii", Description: "Services to a business outside the EU, outside the scope of Spanish VAT.", Content: invoiceJSON("es-sii-v1", "F1", "INT", false)},
			{ID: "sii-b2b-goods-no-eu", Title: "B2B Goods Non-EU Client (Export)", Schema: "bill/invoice", Channel: "sii", Description: "Export of goods outside the EU, exempt with the E2 code.", Content: invoiceJSON("es-sii-v1", "F1", "EXP", false)},
			{ID: "sii-exempt-e1", Title: "B2B Exempt E1 Invoice", Schema: "bill/invoice", Channel: "sii", Description: "Domestic operation exempt under article 20 (E1).", Content: invoiceJSON("es-sii-v1", "F1", "EX", false)},
			{ID: "sii-b2c-oss", Title: "B2C One-Stop-Shop Invoice", Schema: "bill/invoice", Channel: "sii", Description: "Distance sale to an EU consumer taxed under the OSS scheme.", Content: invoiceJSON("es-sii-v1", "F1", "OSS", false)},
			{ID: "sii-ieet", Title: "B2C Hotel + IEET (Catalonia tourist tax)", Schema: "bill/invoice", Channel: "sii", Description: "Accommodation invoice including the Catalan tourist tax as a charge.", Content: invoiceJSON("es-sii-v1", "F1", "HOTEL", false)},
			{ID: "sii-replacement", Title: "Replacement Invoice (Factura de canje)", Schema: "bill/invoice", Channel: "sii", Description: "Full invoice (F3) replacing previously issued simplified invoices.", Content: invoiceJSON("es-sii-v1", "F3", "CANJE", false)},
			{ID: "credit-note", Title: "Credit Note", Schema: "bill/invoice", Channel: "sii", Description: "Negative invoice crediting part or all of the original.", Content: invoiceJSON("es-sii-v1", "R1", "CR", false)},
			{ID: "corrective", Title: "Corrective Invoice", Schema: "bill/invoice", Channel: "sii", Description: "Rectifying invoice that supersedes the original in full.", Content: invoiceJSON("es-sii-v1", "R1", "COR", false)},
		}},
	}},
	"pdf": {Groups: []*DocumentGroup{
		{Channel: "", Schema: "bill/invoice", Documents: []*Document{
			{ID: "pdf-invoice", Title: "Standard Invoice", Schema: "bill/invoice", Description: "A complete invoice with notes and payment terms.", Content: invoiceJSON("", "", "SAMPLE", false)},
		}},
	}},
}
