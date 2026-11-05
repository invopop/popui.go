package popui

import (
	"bytes"
	"context"
	"regexp"
	"strings"
	"testing"

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

// TestInputISOPlaceholder checks ISO fields get one placeholder, their own
// or the format's, and no stray attributes.
func TestInputISOPlaceholder(t *testing.T) {
	for placeholder, want := range map[string]string{"": inputISOPlaceholder, "When": "When"} {
		html := renderInput(t, props.Input{Type: inputDateTimeLocal, ISO: true, Placeholder: placeholder})
		input := regexp.MustCompile(`<input[^>]*>`).FindString(html)
		if strings.Contains(input, " else") {
			t.Errorf("stray else attribute in %s", input)
		}
		if n := strings.Count(input, "placeholder="); n != 1 {
			t.Errorf("%d placeholders in %s", n, input)
		}
		if !strings.Contains(input, `placeholder="`+want+`"`) {
			t.Errorf("want placeholder %q in %s", want, input)
		}
	}
}

// TestInputISOPattern checks the pattern as browsers apply it, to the
// whole value.
func TestInputISOPattern(t *testing.T) {
	re := regexp.MustCompile(`^(?:` + inputISOPattern + `)$`)
	for _, v := range []string{
		"2026-09-20T12:30:00Z",
		"2026-09-20T12:30Z",
		"2026-12-01T09:05:00+01:00",
		"2026-12-01T23:59:59.123-05:30",
		"2026-06-30t23:59:60z",
	} {
		if !re.MatchString(v) {
			t.Errorf("%s rejected", v)
		}
	}
	for _, v := range []string{
		"2026-99-99T99:99:99+99:99",
		"2026-13-01T00:00:00Z",
		"2026-00-10T00:00:00Z",
		"2026-01-32T00:00:00Z",
		"2026-01-01T24:00:00Z",
		"2026-01-01T12:60:00Z",
		"2026-01-01T12:00:61Z",
		"2026-01-01T12:00:00+24:00",
		"2026-09-20 12:30:00Z",
		"2026-09-20T12:30:00",
		"garbage",
	} {
		if re.MatchString(v) {
			t.Errorf("%s accepted", v)
		}
	}
}
