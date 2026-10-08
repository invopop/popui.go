package apps

// Listing is the marketplace-facing content of an application: what a
// workspace admin reads in the Console's app directory before enrolling.
// Modelled on the Stripe App Marketplace listing form; see HANDOFF.md.
type Listing struct {
	// About is the long description (up to 1000 characters).
	About string
	// Publisher is the "Built by" line; defaults to the owning org's name.
	Publisher string
	// Features are up to three key features, in display order.
	Features []Feature
	// Public links.
	WebsiteURL string
	PrivacyURL string
	TermsURL   string
	FAQURL     string
	// Support channels.
	SupportEmail string
	SupportURL   string
	// Countries the app supports (ISO 3166-1 alpha-2), empty for worldwide.
	Countries []string
}

// Feature is one key feature of a listing: a value-focused title, a short
// description and a screenshot.
type Feature struct {
	Title       string
	Description string
	ImageURL    string
}

// maxFeatures is how many key features a listing can carry.
const maxFeatures = 3

// countries is a short stand-in for gobl's l10n country definitions, enough
// for the multiselect to be meaningful.
var countries = []*Definition{
	{Key: "ES", Name: "Spain"}, {Key: "PT", Name: "Portugal"}, {Key: "FR", Name: "France"},
	{Key: "IT", Name: "Italy"}, {Key: "DE", Name: "Germany"}, {Key: "NL", Name: "Netherlands"},
	{Key: "BE", Name: "Belgium"}, {Key: "PL", Name: "Poland"}, {Key: "GR", Name: "Greece"},
	{Key: "DK", Name: "Denmark"}, {Key: "FI", Name: "Finland"}, {Key: "SE", Name: "Sweden"},
	{Key: "GB", Name: "United Kingdom"}, {Key: "IE", Name: "Ireland"}, {Key: "US", Name: "United States"},
	{Key: "CA", Name: "Canada"}, {Key: "MX", Name: "Mexico"}, {Key: "CO", Name: "Colombia"},
	{Key: "BR", Name: "Brazil"}, {Key: "AR", Name: "Argentina"}, {Key: "SA", Name: "Saudi Arabia"},
}

// sampleListings holds the listing content per app; apps without one start
// empty with the publisher defaulted to the org.
var sampleListings = map[string]*Listing{
	"gov-es": {
		About:     "Invopop connects your invoicing with the Spanish tax agency. The Spain app registers invoices with VERI*FACTU and SII, signs TicketBAI invoices for the Basque provinces, and generates Facturae XML for public administrations, all from the same GOBL document.",
		Publisher: "Invopop",
		Features: []Feature{
			{Title: "Register invoices with the AEAT in real time", Description: "Every issued invoice is signed, chained and sent to VERI*FACTU or SII as part of the workflow, with the response stored on the document.", ImageURL: ""},
			{Title: "One supplier registration for all three systems", Description: "Upload the signed agreement once and the app handles the registration flow, approval wait and seat allocation per channel.", ImageURL: ""},
			{Title: "Corrections and cancellations that stay compliant", Description: "Credit notes, corrective invoices and cancellation records reference the original automatically so the tax agency accepts them first time.", ImageURL: ""},
		},
		WebsiteURL:   "https://www.invopop.com",
		PrivacyURL:   "https://www.invopop.com/privacy",
		TermsURL:     "https://www.invopop.com/terms",
		FAQURL:       "https://docs.invopop.com/apps/spain#faq",
		SupportEmail: "support@invopop.com",
		SupportURL:   "https://docs.invopop.com/apps/spain",
		Countries:    []string{"ES"},
	},
	"pdf": {
		About:        "Render any GOBL document as a branded PDF with your logo, colours and the languages your customers read.",
		Publisher:    "Invopop",
		WebsiteURL:   "https://www.invopop.com",
		PrivacyURL:   "https://www.invopop.com/privacy",
		SupportEmail: "support@invopop.com",
	},
	"stripe": {
		Publisher: "Invopop",
	},
}

func listingFor(app *Application) *Listing {
	if l := sampleListings[app.ID]; l != nil {
		return l
	}
	l := &Listing{Publisher: sampleOrg.Name}
	sampleListings[app.ID] = l
	return l
}
