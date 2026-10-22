package apps

import (
	"encoding/json"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ApplicationForm mirrors admin.ApplicationForm in access: the fields the
// application forms present and bind. The `form` tags are the posted names.
type ApplicationForm struct {
	ID               string   `form:"id"`
	Rev              string   `form:"rev"`
	Name             string   `form:"name"`
	Description      string   `form:"description"`
	Categories       []string `form:"categories"`
	Visibility       string   `form:"visibility"`
	GateFeature      string   `form:"gate_feature"`
	ChannelKeys      []string `form:"channel_keys[]"`
	ChannelLabels    []string `form:"channel_labels[]"`
	ClientID         string   `form:"client_id"`
	ClientSecret     string   `form:"client_secret"`
	ClientSecretHash string   `form:"client_secret_hash"`
	RedirectURL      string   `form:"redirect_url"`
	Scopes           []string `form:"scopes"`
	IconURL          string   `form:"icon_url"`
	LogoURL          string   `form:"logo_url"`
	ConfigURL        string   `form:"config_url"`
	LaunchURL        string   `form:"launch_url"`
	InfoURL          string   `form:"info_url"`
	APIURL           string   `form:"api_url"`
	Transport        string   `form:"transport"`
	SigningSecret    string   `form:"signing_secret"`
	SettingsSchema   string   `form:"settings_schema"`
	AllowedOrigins   []string `form:"allowed_origins[]"`
	CreatedAt        string   `form:"created_at"`
	UpdatedAt        string   `form:"updated_at"`
	Disabled         bool     `form:"disabled"`
	Tags             []string `form:"tags[]"`
	AllowOrgIDs      []string `form:"allow_org_ids[]"`
}

func (f *ApplicationForm) persisted() bool {
	return f.Rev != ""
}

func (f *ApplicationForm) transportOrDefault() string {
	if f.Transport == "" {
		return TransportNATS
	}
	return f.Transport
}

func (f *ApplicationForm) idOrPlaceholder() string {
	if f.ID == "" {
		return "<app-id>"
	}
	return f.ID
}

func (f *ApplicationForm) clientIDOrPlaceholder() string {
	if f.ClientID == "" {
		return "<client id>"
	}
	return f.ClientID
}

func (f *ApplicationForm) hasCategory(key string) bool {
	return contains(f.Categories, key)
}

func (f *ApplicationForm) hasScope(key string) bool {
	return contains(f.Scopes, key)
}

// ApplicationModelToForm mirrors admin.ApplicationModelToForm.
func ApplicationModelToForm(a *Application) *ApplicationForm {
	chKeys := make([]string, len(a.Channels))
	chLabels := make([]string, len(a.Channels))
	for i, ch := range a.Channels {
		chKeys[i] = ch.Key
		chLabels[i] = ch.Label
	}
	return &ApplicationForm{
		ID:               a.ID,
		Rev:              a.Rev,
		Name:             a.Name,
		Description:      a.Description,
		Categories:       a.Categories,
		Visibility:       a.Visibility,
		GateFeature:      a.GateFeature,
		ChannelKeys:      chKeys,
		ChannelLabels:    chLabels,
		ClientID:         a.ClientID,
		ClientSecret:     a.ClientSecret,
		ClientSecretHash: a.ClientSecretHash,
		RedirectURL:      a.RedirectURL,
		Scopes:           a.Scopes,
		IconURL:          a.IconURL,
		LogoURL:          a.LogoURL,
		ConfigURL:        a.ConfigURL,
		LaunchURL:        a.LaunchURL,
		InfoURL:          a.InfoURL,
		APIURL:           a.APIURL,
		Transport:        a.Transport,
		SigningSecret:    a.SigningSecret,
		SettingsSchema:   a.SettingsSchema,
		AllowedOrigins:   a.AllowedOrigins,
		CreatedAt:        timeToString(a.CreatedAt),
		UpdatedAt:        timeToString(a.UpdatedAt),
		Disabled:         a.DisabledAt != nil,
		Tags:             a.Tags,
		AllowOrgIDs:      a.AllowOrgIDs,
	}
}

// newApplicationForm is what models.NewApplication("") produces for the
// "Create app" page: private visibility and freshly generated OAuth
// credentials, which the form shows once.
func newApplicationForm() *ApplicationForm {
	return &ApplicationForm{
		Visibility:       AppVisibilityPrivate,
		Transport:        TransportNATS,
		ClientID:         "ci_" + randomHex(8),
		ClientSecret:     "cs_" + randomHex(24),
		ClientSecretHash: "sha256:" + randomHex(16),
		SigningSecret:    "whsec_" + randomHex(24),
	}
}

// channelPair is the JSON shape consumed by the appChannelList Alpine component.
type channelPair struct {
	Key    string `json:"key"`
	Label  string `json:"label"`
	Locked bool   `json:"locked"`
}

// channelPairs seeds the channel editor. Rows from a persisted application
// load with a locked key so existing channel IDs aren't accidentally renamed.
func channelPairs(f *ApplicationForm) []channelPair {
	out := make([]channelPair, 0, len(f.ChannelKeys))
	locked := f.persisted()
	for i, key := range f.ChannelKeys {
		var label string
		if i < len(f.ChannelLabels) {
			label = f.ChannelLabels[i]
		}
		out = append(out, channelPair{Key: key, Label: label, Locked: locked})
	}
	return out
}

func timeToString(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("2006-01-02T15:04:05Z")
}

// ActionForm mirrors admin.ActionForm in access.
type ActionForm struct {
	ID            string `form:"id"`
	Rev           string `form:"rev"`
	Name          string `form:"name"`
	Category      string `form:"category"`
	FeatureID     string `form:"feature_id"`
	Pops          string `form:"pops"`
	Seat          bool   `form:"seat"`
	RenewalPeriod string `form:"renewal_period"`
	Channel       string `form:"channel"`
	// ActionVerb and Qualifier compose the ID of a new action:
	// "<app_id>[.<channel>].<action_verb>[.<qualifier>]". Read on create only.
	ActionVerb   string   `form:"action_verb"`
	Qualifier    string   `form:"qualifier"`
	Unroll       bool     `form:"unroll"`
	Description  string   `form:"description"`
	IconURL      string   `form:"icon_url"`
	ConfigURL    string   `form:"config_url"`
	InfoURL      string   `form:"info_url"`
	ConfigSchema string   `form:"config_schema"`
	Schemas      []string `form:"schemas"`
	Countries    string   `form:"countries"`
	TaxRegimes   string   `form:"tax_regimes"`
	Addons       []string `form:"addons"`

	ExcludeSiloEntry bool `form:"exclude_silo_entry"`

	// Silo entry options
	ConvertIntoCurrency string `form:"convert_into_currency"`
	SharedMeta          bool   `form:"shared_meta"`
	RequireSignature    bool   `form:"require_signature"`
	RemoveIncludedTax   bool   `form:"remove_included_tax"`
	OptionalSiloEntry   bool   `form:"optional_silo_entry"`

	// Results are bound by hand from "results[i].<field>" (see bindResults);
	// access uses go-playground/form for this.
	Results []ActionResultForm `form:"-"`

	WorkerCount string `form:"worker_count"`
	Disabled    bool   `form:"disabled"`

	CreatedAt string `form:"created_at"`
	UpdatedAt string `form:"updated_at"`

	disabledAt string
}

// ActionResultForm mirrors admin.ActionResultForm.
type ActionResultForm struct {
	Status      string `json:"status,omitempty"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	Default     bool   `json:"default,omitempty"`
}

func (f *ActionForm) isNew() bool {
	return f.Rev == ""
}

func (f *ActionForm) hasSchema(s string) bool { return contains(f.Schemas, s) }
func (f *ActionForm) hasAddon(a string) bool  { return contains(f.Addons, a) }

// composeActionID mirrors admin.composeActionID.
func composeActionID(appID, channel, action, qualifier string) string {
	parts := []string{appID}
	for _, p := range []string{channel, action, qualifier} {
		if p != "" {
			parts = append(parts, p)
		}
	}
	return strings.Join(parts, ".")
}

// ToAction mirrors admin.ActionForm.ToGatewayAction.
func (f *ActionForm) ToAction(appID string) *Action {
	id := f.ID
	if f.isNew() {
		id = composeActionID(appID, f.Channel, f.ActionVerb, f.Qualifier)
	}
	act := &Action{
		Id:                  id,
		Rev:                 f.Rev,
		ApplicationId:       appID,
		Name:                f.Name,
		Pops:                parseNumber(f.Pops),
		Category:            f.Category,
		FeatureId:           f.FeatureID,
		Seat:                f.Seat,
		RenewalPeriod:       f.RenewalPeriod,
		Unroll:              f.Unroll,
		Description:         f.Description,
		IconUrl:             f.IconURL,
		ConfigUrl:           f.ConfigURL,
		InfoUrl:             f.InfoURL,
		ConfigSchema:        f.ConfigSchema,
		Schemas:             f.Schemas,
		Countries:           parseStringArray(f.Countries),
		TaxRegimes:          parseStringArray(f.TaxRegimes),
		Addons:              f.Addons,
		ExcludeSiloEntry:    f.ExcludeSiloEntry,
		ConvertIntoCurrency: f.ConvertIntoCurrency,
		SharedMeta:          f.SharedMeta,
		RequireSignature:    f.RequireSignature,
		RemoveIncludedTax:   f.RemoveIncludedTax,
		OptionalSiloEntry:   f.OptionalSiloEntry,
		WorkerCount:         parseNumber(f.WorkerCount),
		CreatedAt:           f.CreatedAt,
		UpdatedAt:           f.UpdatedAt,
	}
	for _, r := range f.Results {
		act.Results = append(act.Results, &ActionResult{
			Status:      parseNumber(r.Status),
			Code:        r.Code,
			Description: r.Description,
			Default:     r.Default,
		})
	}
	if f.Disabled {
		act.DisabledAt = timeToString(time.Now())
	}
	return act
}

// actionFormFromAction mirrors admin.actionFormFromGateway.
func actionFormFromAction(act *Action) *ActionForm {
	f := &ActionForm{
		ID:                  act.Id,
		Rev:                 act.Rev,
		Name:                act.Name,
		Category:            act.Category,
		FeatureID:           act.FeatureId,
		Seat:                act.Seat,
		RenewalPeriod:       act.RenewalPeriod,
		Unroll:              act.Unroll,
		Description:         act.Description,
		IconURL:             act.IconUrl,
		ConfigURL:           act.ConfigUrl,
		InfoURL:             act.InfoUrl,
		ConfigSchema:        act.ConfigSchema,
		ExcludeSiloEntry:    act.ExcludeSiloEntry,
		Schemas:             act.Schemas,
		Countries:           strings.Join(act.Countries, ", "),
		TaxRegimes:          strings.Join(act.TaxRegimes, ", "),
		ConvertIntoCurrency: act.ConvertIntoCurrency,
		SharedMeta:          act.SharedMeta,
		RequireSignature:    act.RequireSignature,
		RemoveIncludedTax:   act.RemoveIncludedTax,
		OptionalSiloEntry:   act.OptionalSiloEntry,
		Results:             []ActionResultForm{},
		Addons:              act.Addons,
		WorkerCount:         strconv.Itoa(int(act.WorkerCount)),
		Pops:                strconv.Itoa(int(act.Pops)),
		Disabled:            act.DisabledAt != "",
		CreatedAt:           act.CreatedAt,
		UpdatedAt:           act.UpdatedAt,
		disabledAt:          act.DisabledAt,
	}
	for _, r := range act.Results {
		f.Results = append(f.Results, ActionResultForm{
			Status:      strconv.Itoa(int(r.Status)),
			Code:        r.Code,
			Description: r.Description,
			Default:     r.Default,
		})
	}
	return f
}

// bindResults reads the "results[i].<field>" inputs the action results
// manager submits. Rows are numbered from zero without gaps.
func bindResults(values url.Values) []ActionResultForm {
	var out []ActionResultForm
	for i := 0; ; i++ {
		prefix := "results[" + strconv.Itoa(i) + "]."
		if !values.Has(prefix+"status") && !values.Has(prefix+"code") {
			return out
		}
		out = append(out, ActionResultForm{
			Status:      values.Get(prefix + "status"),
			Code:        values.Get(prefix + "code"),
			Description: values.Get(prefix + "description"),
			Default:     values.Get(prefix+"default") == "true",
		})
	}
}

func actionResultsToJSON(results []ActionResultForm) string {
	if results == nil {
		results = []ActionResultForm{}
	}
	out, _ := json.Marshal(results)
	return string(out)
}

var arrayBadCharsRegexp = regexp.MustCompile(`[^A-Z,]`)

func parseNumber(src string) int32 {
	n, err := strconv.Atoi(src)
	if err != nil {
		return 0
	}
	return int32(n)
}

func parseStringArray(src string) []string {
	if src == "" {
		return nil
	}
	src = arrayBadCharsRegexp.ReplaceAllString(src, "")
	return strings.Split(src, ",")
}

func contains(list []string, v string) bool {
	for _, s := range list {
		if s == v {
			return true
		}
	}
	return false
}

func randomHex(n int) string {
	const hex = "0123456789abcdef"
	b := make([]byte, n)
	seed := time.Now().UnixNano()
	for i := range b {
		seed = seed*6364136223846793005 + 1442695040888963407
		b[i] = hex[(seed>>33)&0xf]
	}
	return string(b)
}

// transportOptions are the delivery modes offered on the App tab.
var transportOptions = []*Definition{
	{Key: TransportNATS, Name: "NATS request/reply", Desc: "The app runs a worker that subscribes to gw.<app>.task with the credentials issued on the Credentials tab."},
	{Key: TransportHTTP, Name: "HTTP webhook", Desc: "The gateway POSTs each task to the API URL, signed with the app's signing secret; the result comes back in the response."},
}

// actionResultStatuses are the statuses an action can report. The gateway's
// enum today exposes only NA/OK/KO in the admin while the SDK also returns
// SKIP, ERR and QUEUED; the numeric values for the last three are a proposal
// to confirm against gateway.ActionResultStatus.
var actionResultStatuses = []*Definition{
	{Key: "0", Name: "NA", Desc: "not applicable"},
	{Key: "1", Name: "OK", Desc: "success"},
	{Key: "2", Name: "KO", Desc: "permanent failure"},
	{Key: "3", Name: "SKIP", Desc: "nothing to do"},
	{Key: "4", Name: "ERR", Desc: "temporary failure, retried"},
	{Key: "5", Name: "QUEUED", Desc: "result reported later"},
}
