// Package period reads the windows of time picked with popui.Period: a
// rolling number of days up to today, every time, or calendar days, weeks,
// months, quarters or years, alone or in a run of the same unit.
//
// A period is kept as a short string, as in a URL's query:
//
//	365d                     the last 365 days, today included
//	all                      every time
//	2026-11-02               a day
//	2026-W44                 an ISO week, from Monday
//	2026-11                  a month
//	2026-Q4                  a quarter
//	2026                     a year
//	2026-01..2026-03         a run of windows of the same unit
//
// popui.js's periodSelect controller reads and writes the same strings.
package period

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Units of calendar windows.
const (
	Day     = "day"
	Week    = "week"
	Month   = "month"
	Quarter = "quarter"
	Year    = "year"
)

// All is the period of every time.
const All = "all"

// Default is the period shown when none is picked: the last 365 days.
const Default = "365d"

// Period is a window of time. Exactly one of All, Days or Unit is set.
type Period struct {
	// All is every time, with neither start nor end.
	All bool
	// Days is a rolling number of days up to and including today.
	Days int
	// Unit is the unit of calendar windows, and From and To the first and
	// last of them, counted from the Unix epoch: days, ISO weeks from the
	// one starting 29 December 1969, months and quarters from year 0, and
	// years.
	Unit     string
	From, To int
}

var (
	reDays    = regexp.MustCompile(`^([1-9][0-9]{0,4})d$`)
	reDay     = regexp.MustCompile(`^([0-9]{4})-([0-9]{2})-([0-9]{2})$`)
	reWeek    = regexp.MustCompile(`^([0-9]{4})-W([0-9]{2})$`)
	reMonth   = regexp.MustCompile(`^([0-9]{4})-([0-9]{2})$`)
	reQuarter = regexp.MustCompile(`^([0-9]{4})-Q([1-4])$`)
	reYear    = regexp.MustCompile(`^([0-9]{4})$`)
)

// Parse reads a period, or Default when s is empty.
func Parse(s string) (Period, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		s = Default
	}
	if s == All {
		return Period{All: true}, nil
	}
	if m := reDays.FindStringSubmatch(s); m != nil {
		n, _ := strconv.Atoi(m[1])
		return Period{Days: n}, nil
	}
	a, b, ok := strings.Cut(s, "..")
	if !ok {
		b = a
	}
	ua, from, err := parseWindow(a)
	if err != nil {
		return Period{}, err
	}
	ub, to, err := parseWindow(b)
	if err != nil {
		return Period{}, err
	}
	if ua != ub {
		return Period{}, fmt.Errorf("%q mixes %ss and %ss", s, ua, ub)
	}
	if to < from {
		from, to = to, from
	}
	return Period{Unit: ua, From: from, To: to}, nil
}

func parseWindow(s string) (string, int, error) {
	atoi := func(v string) int { n, _ := strconv.Atoi(v); return n }
	if m := reDay.FindStringSubmatch(s); m != nil {
		t, err := time.Parse(time.DateOnly, s)
		if err != nil {
			return "", 0, fmt.Errorf("%q isn't a date", s)
		}
		return Day, dayIndex(t), nil
	}
	if m := reWeek.FindStringSubmatch(s); m != nil {
		y, w := atoi(m[1]), atoi(m[2])
		if w < 1 || w > weeksIn(y) {
			return "", 0, fmt.Errorf("%q isn't a week", s)
		}
		return Week, weekOf(dayIndex(date(y, 1, 4))) + w - 1, nil
	}
	if m := reMonth.FindStringSubmatch(s); m != nil {
		y, mo := atoi(m[1]), atoi(m[2])
		if mo < 1 || mo > 12 {
			return "", 0, fmt.Errorf("%q isn't a month", s)
		}
		return Month, y*12 + mo - 1, nil
	}
	if m := reQuarter.FindStringSubmatch(s); m != nil {
		return Quarter, atoi(m[1])*4 + atoi(m[2]) - 1, nil
	}
	if m := reYear.FindStringSubmatch(s); m != nil {
		return Year, atoi(m[1]), nil
	}
	return "", 0, fmt.Errorf("%q isn't a period", s)
}

// String writes the period as Parse reads it.
func (p Period) String() string {
	switch {
	case p.All:
		return All
	case p.Days > 0:
		return strconv.Itoa(p.Days) + "d"
	case p.From == p.To:
		return p.key(p.From)
	}
	return p.key(p.From) + ".." + p.key(p.To)
}

// key writes window n of the period's unit.
func (p Period) key(n int) string {
	switch p.Unit {
	case Day:
		return fromDay(n).Format(time.DateOnly)
	case Week:
		y, w := isoWeek(n)
		return fmt.Sprintf("%04d-W%02d", y, w)
	case Month:
		return fmt.Sprintf("%04d-%02d", n/12, n%12+1)
	case Quarter:
		return fmt.Sprintf("%04d-Q%d", n/4, n%4+1)
	}
	return fmt.Sprintf("%04d", n)
}

// Range returns the start of the period's first day and the start of the
// day after its last, in today's location. Both are zero for All. Rolling
// days end with today.
func (p Period) Range(today time.Time) (from, to time.Time) {
	if p.All {
		return time.Time{}, time.Time{}
	}
	first, last := p.days(dayIndex(today))
	loc := today.Location()
	at := func(n int) time.Time {
		d := fromDay(n)
		return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, loc)
	}
	return at(first), at(last + 1)
}

// RangeWithin returns the part of the period's Range inside outer's, as
// for a window from outer.Split, whose first and last windows may reach
// past outer's ends. Both are zero when they don't overlap. All has no
// bounds, so an outer All leaves the Range as it is.
func (p Period) RangeWithin(outer Period, today time.Time) (from, to time.Time) {
	from, to = p.Range(today)
	if outer.All {
		return from, to
	}
	of, ot := outer.Range(today)
	if p.All || from.Before(of) {
		from = of
	}
	if p.All || to.After(ot) {
		to = ot
	}
	if !from.Before(to) {
		return time.Time{}, time.Time{}
	}
	return from, to
}

// days returns the period's first and last days, from the Unix epoch.
func (p Period) days(today int) (int, int) {
	if p.Days > 0 {
		return today - p.Days + 1, today
	}
	return p.start(p.From), p.start(p.To+1) - 1
}

// start returns the first day of window n of the period's unit.
func (p Period) start(n int) int {
	switch p.Unit {
	case Day:
		return n
	case Week:
		return n*7 - 3
	case Month:
		return dayIndex(date(floorDiv(n, 12), n-floorDiv(n, 12)*12+1, 1))
	case Quarter:
		return dayIndex(date(floorDiv(n, 4), (n-floorDiv(n, 4)*4)*3+1, 1))
	}
	return dayIndex(date(n, 1, 1))
}

var monthNames = [...]string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}

// Label names the period for people, as "Last 365 days", "Q4 2026" or
// "Jan – Mar 2026".
func (p Period) Label() string {
	switch {
	case p.All:
		return "All time"
	case p.Days == 1:
		return "Today"
	case p.Days > 0:
		return "Last " + strconv.Itoa(p.Days) + " days"
	}
	if p.From == p.To {
		return p.name(p.From, true)
	}
	if p.year(p.From) != p.year(p.To) {
		return p.name(p.From, true) + " – " + p.name(p.To, true)
	}
	switch p.Unit {
	case Week:
		_, a := isoWeek(p.From)
		y, b := isoWeek(p.To)
		return fmt.Sprintf("Weeks %d – %d, %d", a, b, y)
	case Year:
		return p.name(p.From, true) + " – " + p.name(p.To, true)
	}
	return p.name(p.From, false) + " – " + p.name(p.To, true)
}

// year returns the year window n of the period's unit is in, ISO years for
// weeks.
func (p Period) year(n int) int {
	switch p.Unit {
	case Day:
		return fromDay(n).Year()
	case Week:
		y, _ := isoWeek(n)
		return y
	case Month:
		return floorDiv(n, 12)
	case Quarter:
		return floorDiv(n, 4)
	}
	return n
}

// name names window n of the period's unit, with its year or without.
func (p Period) name(n int, withYear bool) string {
	y := ""
	if withYear {
		y = " " + strconv.Itoa(p.year(n))
	}
	switch p.Unit {
	case Day:
		d := fromDay(n)
		return strconv.Itoa(d.Day()) + " " + monthNames[d.Month()-1][:3] + y
	case Week:
		_, w := isoWeek(n)
		if withYear {
			y = "," + y
		}
		return "Week " + strconv.Itoa(w) + y
	case Month:
		name := monthNames[n-floorDiv(n, 12)*12]
		if !withYear {
			name = name[:3]
		} else if p.From != p.To {
			name = name[:3]
		}
		return name + y
	case Quarter:
		return "Q" + strconv.Itoa(n-floorDiv(n, 4)*4+1) + y
	}
	return strconv.Itoa(n)
}

// Shift moves the period by n of its own lengths: a month to the month
// before for -1, or the last 30 days to the 30 before them. All doesn't
// move.
func (p Period) Shift(n int, today time.Time) Period {
	switch {
	case p.All:
		return p
	case p.Days > 0:
		first, last := p.days(dayIndex(today))
		return Period{Unit: Day, From: first + n*p.Days, To: last + n*p.Days}
	}
	span := p.To - p.From + 1
	return Period{Unit: p.Unit, From: p.From + n*span, To: p.To + n*span}
}

// Split divides the period into the calendar windows of the finest unit
// that gives at most limit of them, oldest first, as a chart's columns run:
// days, weeks, months, quarters, or else years however many there are.
// The first and last windows may reach past the period's ends: use each
// window's RangeWithin the period to count only what's in it. All has no
// start to split from, so gives none, and nor does a limit below one.
func (p Period) Split(today time.Time, limit int) []Period {
	if p.All || limit < 1 {
		return nil
	}
	first, last := p.days(dayIndex(today))
	for _, u := range []string{Day, Week, Month, Quarter, Year} {
		a, b := windowOf(u, first), windowOf(u, last)
		if b-a < limit || u == Year {
			out := make([]Period, 0, b-a+1)
			for n := a; n <= b; n++ {
				out = append(out, Period{Unit: u, From: n, To: n})
			}
			return out
		}
	}
	return nil
}

// Short names the period's first window briefly, for a chart's axis: "2
// Nov", "W44", "Nov", "Q4" or "2026". Label names it in full.
func (p Period) Short() string {
	switch {
	case p.All:
		return "All"
	case p.Days > 0:
		return strconv.Itoa(p.Days) + "d"
	case p.Unit == Week:
		_, w := isoWeek(p.From)
		return "W" + strconv.Itoa(w)
	case p.Unit == Month:
		return monthNames[p.From-floorDiv(p.From, 12)*12][:3]
	}
	return p.name(p.From, false)
}

// windowOf returns the window of the unit that day n is in.
func windowOf(unit string, n int) int {
	d := fromDay(n)
	switch unit {
	case Day:
		return n
	case Week:
		return weekOf(n)
	case Month:
		return d.Year()*12 + int(d.Month()) - 1
	case Quarter:
		return d.Year()*4 + (int(d.Month())-1)/3
	}
	return d.Year()
}

func date(y, m, d int) time.Time {
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
}

// dayIndex counts the days from the Unix epoch to t's date.
func dayIndex(t time.Time) int {
	return int(floorDiv64(date(t.Year(), int(t.Month()), t.Day()).Unix(), 86400))
}

func fromDay(n int) time.Time {
	return time.Unix(int64(n)*86400, 0).UTC()
}

// weekOf returns the ISO week day n is in. The epoch was a Thursday, so
// day -3 starts week 0.
func weekOf(n int) int {
	return floorDiv(n+3, 7)
}

// isoWeek returns the ISO year and week number of week n, the year of its
// Thursday.
func isoWeek(n int) (int, int) {
	th := fromDay(n*7 - 3 + 3)
	return th.Year(), (th.YearDay()-1)/7 + 1
}

// weeksIn returns how many ISO weeks year y has.
func weeksIn(y int) int {
	_, w := date(y, 12, 28).ISOWeek()
	return w
}

func floorDiv(a, b int) int {
	q := a / b
	if (a%b != 0) && ((a < 0) != (b < 0)) {
		q--
	}
	return q
}

func floorDiv64(a, b int64) int64 {
	q := a / b
	if (a%b != 0) && ((a < 0) != (b < 0)) {
		q--
	}
	return q
}
