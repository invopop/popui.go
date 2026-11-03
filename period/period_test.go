package period

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	madrid, _ := time.LoadLocation("Europe/Madrid")
	today := time.Date(2026, 11, 2, 15, 4, 0, 0, madrid)
	day := func(y, m, d int) time.Time { return time.Date(y, time.Month(m), d, 0, 0, 0, 0, madrid) }
	tests := []struct {
		in, out, label string
		from, to       time.Time
	}{
		{"", "365d", "Last 365 days", day(2025, 11, 3), day(2026, 11, 3)},
		{"7d", "", "Last 7 days", day(2026, 10, 27), day(2026, 11, 3)},
		{"1d", "", "Today", day(2026, 11, 2), day(2026, 11, 3)},
		{All, "", "All time", time.Time{}, time.Time{}},
		{"2026-11-02", "", "2 Nov 2026", day(2026, 11, 2), day(2026, 11, 3)},
		{"2026-10-01..2026-11-15", "", "1 Oct – 15 Nov 2026", day(2026, 10, 1), day(2026, 11, 16)},
		{"2026-W44", "", "Week 44, 2026", day(2026, 10, 26), day(2026, 11, 2)},
		{"2026-W01", "", "Week 1, 2026", day(2025, 12, 29), day(2026, 1, 5)},
		{"2020-W53", "", "Week 53, 2020", day(2020, 12, 28), day(2021, 1, 4)},
		{"2026-W40..2026-W44", "", "Weeks 40 – 44, 2026", day(2026, 9, 28), day(2026, 11, 2)},
		{"2026-11", "", "November 2026", day(2026, 11, 1), day(2026, 12, 1)},
		{"2026-03..2026-01", "2026-01..2026-03", "Jan – Mar 2026", day(2026, 1, 1), day(2026, 4, 1)},
		{"2025-11..2026-02", "", "Nov 2025 – Feb 2026", day(2025, 11, 1), day(2026, 3, 1)},
		{"2026-Q4", "", "Q4 2026", day(2026, 10, 1), day(2027, 1, 1)},
		{"2026-Q1..2026-Q3", "", "Q1 – Q3 2026", day(2026, 1, 1), day(2026, 10, 1)},
		{"2025-Q4..2026-Q1", "", "Q4 2025 – Q1 2026", day(2025, 10, 1), day(2026, 4, 1)},
		{"2026", "", "2026", day(2026, 1, 1), day(2027, 1, 1)},
		{"2024..2026", "", "2024 – 2026", day(2024, 1, 1), day(2027, 1, 1)},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			p, err := Parse(tt.in)
			if err != nil {
				t.Fatal(err)
			}
			out := tt.out
			if out == "" {
				out = tt.in
			}
			if s := p.String(); s != out {
				t.Errorf("String() = %q, want %q", s, out)
			}
			if l := p.Label(); l != tt.label {
				t.Errorf("Label() = %q, want %q", l, tt.label)
			}
			from, to := p.Range(today)
			if !from.Equal(tt.from) || !to.Equal(tt.to) {
				t.Errorf("Range() = %v, %v, want %v, %v", from, to, tt.from, tt.to)
			}
		})
	}
}

func TestParseErrors(t *testing.T) {
	for _, in := range []string{"x", "0d", "2026-13", "2026-W54", "2026-02-30", "2026-Q5", "2026..2026-01", "2026-01..x"} {
		if _, err := Parse(in); err == nil {
			t.Errorf("Parse(%q) didn't fail", in)
		}
	}
}

func TestShift(t *testing.T) {
	today := time.Date(2026, 11, 2, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		in   string
		n    int
		want string
	}{
		{"2026-11-02", -1, "2026-11-01"},
		{"2026-01", -1, "2025-12"},
		{"2026-Q1..2026-Q2", -1, "2025-Q3..2025-Q4"},
		{"2026-W02", -2, "2025-W52"},
		{"7d", -1, "2026-10-20..2026-10-26"},
		{"2025", 1, "2026"},
		{"all", -1, "all"},
	}
	for _, tt := range tests {
		p, _ := Parse(tt.in)
		if got := p.Shift(tt.n, today).String(); got != tt.want {
			t.Errorf("Parse(%q).Shift(%d) = %q, want %q", tt.in, tt.n, got, tt.want)
		}
	}
}
