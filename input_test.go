package popui

import (
	"bytes"
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/invopop/popui.go/props"
)

func renderInput(t *testing.T, p props.Input) string {
	t.Helper()
	var buf bytes.Buffer
	if err := Input(p).Render(context.Background(), &buf); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

var inputTag = regexp.MustCompile(`<input[^>]*>`)

// TestInputISODefault checks date and time inputs are ISO text fields with
// a hidden native picker of their type, and Local ones are native controls.
func TestInputISODefault(t *testing.T) {
	for typ := range inputISOPlaceholders {
		inputs := inputTag.FindAllString(renderInput(t, props.Input{Type: typ}), -1)
		if len(inputs) != 2 {
			t.Fatalf("%s: want field and picker, got %q", typ, inputs)
		}
		field, picker := inputs[0], inputs[1]
		for _, want := range []string{`type="text"`, `data-iso="` + typ + `"`, `pattern="`, `class="` /* font-mono checked below */} {
			if !strings.Contains(field, want) {
				t.Errorf("%s field: want %s in %s", typ, want, field)
			}
		}
		if !strings.Contains(field, "font-mono") || strings.Contains(field, "data-pickable") || strings.Contains(field, "data-clearable") {
			t.Errorf("%s field: %s", typ, field)
		}
		if !strings.Contains(picker, `type="`+typ+`"`) || !strings.Contains(picker, "data-iso-picker") || strings.Contains(picker, "name=") {
			t.Errorf("%s picker: %s", typ, picker)
		}

		local := inputTag.FindAllString(renderInput(t, props.Input{Type: typ, Local: true}), -1)
		if len(local) != 1 || !strings.Contains(local[0], `type="`+typ+`"`) || strings.Contains(local[0], "data-iso") {
			t.Errorf("%s local: %q", typ, local)
		}
	}
	// Disabled fields have no picker.
	if n := len(inputTag.FindAllString(renderInput(t, props.Input{Type: inputDate, Disabled: true}), -1)); n != 1 {
		t.Errorf("disabled date: %d inputs", n)
	}
}

// TestInputISOPlaceholder checks ISO fields get one placeholder, their own
// or the format's, and no stray attributes.
func TestInputISOPlaceholder(t *testing.T) {
	for placeholder, want := range map[string]string{"": inputISOPlaceholders[inputDateTimeLocal], "When": "When"} {
		field := inputTag.FindString(renderInput(t, props.Input{Type: inputDateTimeLocal, Placeholder: placeholder}))
		if strings.Contains(field, " else") {
			t.Errorf("stray else attribute in %s", field)
		}
		if n := strings.Count(field, "placeholder="); n != 1 {
			t.Errorf("%d placeholders in %s", n, field)
		}
		if !strings.Contains(field, `placeholder="`+want+`"`) {
			t.Errorf("want placeholder %q in %s", want, field)
		}
	}
}

// TestInputISOPickerAttributes checks min, max and step reach the picker.
func TestInputISOPickerAttributes(t *testing.T) {
	html := renderInput(t, props.Input{Type: inputDate, Attributes: templ.Attributes{pickerAttributes[0]: "2026-01-01", pickerAttributes[1]: "2026-12-31"}})
	picker := inputTag.FindAllString(html, -1)[1]
	for _, want := range []string{`min="2026-01-01"`, `max="2026-12-31"`} {
		if !strings.Contains(picker, want) {
			t.Errorf("want %s in %s", want, picker)
		}
	}
	picker = inputTag.FindAllString(renderInput(t, props.Input{Type: inputTime}), -1)[1]
	if !strings.Contains(picker, `step="1"`) {
		t.Errorf("want step in %s", picker)
	}
}

// TestInputISOPatterns checks the patterns as browsers apply them, to the
// whole value.
func TestInputISOPatterns(t *testing.T) {
	for typ, cases := range map[string]struct{ ok, bad []string }{
		inputDate: {
			ok:  []string{"2026-09-20", "0000-01-01"},
			bad: []string{"2026-13-01", "2026-00-10", "2026-01-32", "20/09/2026", "2026-9-20"},
		},
		inputTime: {
			ok:  []string{"12:30", "23:59:59", "00:00:00.123"},
			bad: []string{"24:00", "12:60", "12:30:60", "9:30", "12:30Z"},
		},
		inputDateTimeLocal: {
			ok: []string{
				"2026-09-20T12:30:00Z", "2026-09-20T12:30Z", "2026-12-01T09:05:00+01:00",
				"2026-12-01T23:59:59.123-05:30", "2026-06-30t23:59:60z",
			},
			bad: []string{
				"2026-99-99T99:99:99+99:99", "2026-01-01T24:00:00Z", "2026-01-01T12:00:61Z",
				"2026-01-01T12:00:00+24:00", "2026-09-20 12:30:00Z", "2026-09-20T12:30:00", "garbage",
			},
		},
		inputWeek:  {ok: []string{"2026-W01", "2026-W53"}, bad: []string{"2026-W00", "2026-W54", "2026W01", "2026-w01"}},
		inputMonth: {ok: []string{"2026-09", "2026-12"}, bad: []string{"2026-13", "2026-00", "2026-9"}},
	} {
		re := regexp.MustCompile(`^(?:` + inputISOPatterns[typ] + `)$`)
		for _, v := range cases.ok {
			if !re.MatchString(v) {
				t.Errorf("%s: %s rejected", typ, v)
			}
		}
		for _, v := range cases.bad {
			if re.MatchString(v) {
				t.Errorf("%s: %s accepted", typ, v)
			}
		}
	}
}
