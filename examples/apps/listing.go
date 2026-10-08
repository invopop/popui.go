package apps

import (
	"strings"

	"github.com/labstack/echo/v4"
)

// ListingForm is the Listing tab's fields. The three key feature slots are
// posted as parallel arrays in deck order; empty slots are dropped.
type ListingForm struct {
	About               string   `form:"about"`
	Publisher           string   `form:"publisher"`
	FeatureTitles       []string `form:"feature_titles[]"`
	FeatureDescriptions []string `form:"feature_descriptions[]"`
	FeatureImageURLs    []string `form:"feature_image_urls[]"`
	FeatureOrder        string   `form:"feature_order"`
	WebsiteURL          string   `form:"website_url"`
	PrivacyURL          string   `form:"privacy_url"`
	TermsURL            string   `form:"terms_url"`
	FAQURL              string   `form:"faq_url"`
	SupportEmail        string   `form:"support_email"`
	SupportURL          string   `form:"support_url"`
	Countries           []string `form:"countries"`
}

func listingFormFrom(l *Listing) *ListingForm {
	f := &ListingForm{
		About:        l.About,
		Publisher:    l.Publisher,
		WebsiteURL:   l.WebsiteURL,
		PrivacyURL:   l.PrivacyURL,
		TermsURL:     l.TermsURL,
		FAQURL:       l.FAQURL,
		SupportEmail: l.SupportEmail,
		SupportURL:   l.SupportURL,
		Countries:    l.Countries,
	}
	for _, ft := range l.Features {
		f.FeatureTitles = append(f.FeatureTitles, ft.Title)
		f.FeatureDescriptions = append(f.FeatureDescriptions, ft.Description)
		f.FeatureImageURLs = append(f.FeatureImageURLs, ft.ImageURL)
	}
	return f
}

// features zips the parallel feature fields, applying the deck's order
// ("feature-1,feature-3,feature-2") and skipping slots without a title.
func (f *ListingForm) features() []Feature {
	n := len(f.FeatureTitles)
	idx := make([]int, 0, n)
	if f.FeatureOrder != "" {
		for _, id := range strings.Split(f.FeatureOrder, ",") {
			if i := parseNumber(strings.TrimPrefix(id, "feature-")) - 1; i >= 0 && int(i) < n {
				idx = append(idx, int(i))
			}
		}
	}
	if len(idx) != n {
		idx = idx[:0]
		for i := 0; i < n; i++ {
			idx = append(idx, i)
		}
	}
	var out []Feature
	for _, i := range idx {
		if strings.TrimSpace(f.FeatureTitles[i]) == "" {
			continue
		}
		ft := Feature{Title: f.FeatureTitles[i]}
		if i < len(f.FeatureDescriptions) {
			ft.Description = f.FeatureDescriptions[i]
		}
		if i < len(f.FeatureImageURLs) {
			ft.ImageURL = f.FeatureImageURLs[i]
		}
		out = append(out, ft)
	}
	return out
}

// slots returns exactly maxFeatures features for rendering: the saved ones
// followed by empty slots.
func (f *ListingForm) slots() []Feature {
	out := make([]Feature, maxFeatures)
	for i := 0; i < maxFeatures && i < len(f.FeatureTitles); i++ {
		out[i].Title = f.FeatureTitles[i]
		if i < len(f.FeatureDescriptions) {
			out[i].Description = f.FeatureDescriptions[i]
		}
		if i < len(f.FeatureImageURLs) {
			out[i].ImageURL = f.FeatureImageURLs[i]
		}
	}
	return out
}

func (f *ListingForm) hasCountry(code string) bool { return contains(f.Countries, code) }

func listing(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	return render(c, AppListing(sampleOrg, app, listingFormFrom(listingFor(app)), nil))
}

// updateListing re-renders what was posted, as the app form does. The feature
// slots are re-packed so the saved order and gaps read back correctly.
func updateListing(c echo.Context) error {
	app, err := findApp(c)
	if err != nil {
		return err
	}
	form := new(ListingForm)
	if err := c.Bind(form); err != nil {
		return err
	}
	l := &Listing{Features: form.features()}
	packed := listingFormFrom(l)
	form.FeatureTitles, form.FeatureDescriptions, form.FeatureImageURLs = packed.FeatureTitles, packed.FeatureDescriptions, packed.FeatureImageURLs
	form.FeatureOrder = ""
	success(c, "Listing saved.")
	return render(c, AppListing(sampleOrg, app, form, nil))
}
