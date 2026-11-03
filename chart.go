package popui

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/invopop/popui.go/props"
)

// chartTicks returns the top of a value axis reaching max, and the round
// values to mark on it from zero: steps of 1, 2 or 5 times a power of ten,
// three or four of them.
func chartTicks(peak float64) (float64, []float64) {
	if peak <= 0 || math.IsNaN(peak) || math.IsInf(peak, 0) {
		return 1, []float64{0}
	}
	raw := peak / 3
	mag := math.Pow(10, math.Floor(math.Log10(raw)))
	step := mag
	for _, m := range []float64{1, 2, 5, 10} {
		if m*mag >= raw {
			step = m * mag
			break
		}
	}
	// Counts don't have fractions to mark.
	if step < 1 && peak >= 1 {
		step = 1
	}
	top := math.Ceil(peak/step) * step
	var ticks []float64
	for v := 0.0; v <= top+step/2; v += step {
		ticks = append(ticks, v)
	}
	return top, ticks
}

// FormatNumber writes a figure with thousands separators, and with one
// decimal when it isn't whole: 1,284 or 12.5.
func FormatNumber(v float64) string {
	neg := v < 0
	v = math.Abs(v)
	whole := math.Round(v*10)/10 == math.Round(v)
	s := strconv.FormatFloat(v, 'f', 1, 64)
	if whole {
		s = strconv.FormatFloat(math.Round(v), 'f', 0, 64)
	}
	intPart, frac, _ := strings.Cut(s, ".")
	var b strings.Builder
	for i, r := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			b.WriteByte(',')
		}
		b.WriteRune(r)
	}
	if frac != "" {
		b.WriteString("." + frac)
	}
	if neg {
		return "-" + b.String()
	}
	return b.String()
}

// CompactNumber writes a figure briefly for an axis or a tight space: in
// full under 10,000, and in thousands or millions above, as 12K or 4.2M.
func CompactNumber(v float64) string {
	a := math.Abs(v)
	switch {
	case a < 10_000:
		return FormatNumber(v)
	case a < 1_000_000:
		return compactUnit(v/1_000) + "K"
	case a < 1_000_000_000:
		return compactUnit(v/1_000_000) + "M"
	}
	return compactUnit(v/1_000_000_000) + "B"
}

func compactUnit(v float64) string {
	if math.Abs(v) >= 100 {
		return strconv.FormatFloat(math.Round(v), 'f', 0, 64)
	}
	return strings.TrimSuffix(strconv.FormatFloat(math.Round(v*10)/10, 'f', 1, 64), ".0")
}

// chartDisplay is how a value is shown.
func chartDisplay(v props.ChartValue) string {
	if v.Display != "" {
		return v.Display
	}
	return FormatNumber(v.Value)
}

// chartTitle is how a value is named in full.
func chartTitle(v props.ChartValue) string {
	if v.Title != "" {
		return v.Title
	}
	return v.Label
}

// colorGrey is TagStatus's neutral color.
const colorGrey = "grey"

// chartColor returns the class that fills a mark with the value's color,
// matching TagStatus's dots.
func chartColor(color string) string {
	switch color {
	case "green":
		return "bg-icon-status-paid"
	case "yellow":
		return "bg-icon-status-processing"
	case "red":
		return "bg-icon-status-error"
	case "orange":
		return "bg-icon-status-draft"
	case "blue":
		return "bg-icon-status-sent"
	case "purple":
		return "bg-icon-status-registered"
	case "olive":
		return "bg-sherwood-50"
	case "teal":
		return "bg-icon-status-completed"
	case "crimson":
		return "bg-icon-status-rejected"
	case "steelBlue":
		return "bg-icon-status-received"
	case "empty":
		return "bg-icon-status-empty"
	case colorGrey:
		return "bg-icon-status-void"
	}
	return "bg-background-accent"
}

// chartMax is the largest of the values.
func chartMax(vs []props.ChartValue) float64 {
	m := 0.0
	for _, v := range vs {
		m = max(m, v.Value)
	}
	return m
}

// chartPercent is v as a percentage of whole, for CSS.
func chartPercent(v, whole float64) string {
	if whole <= 0 || v <= 0 {
		return "0%"
	}
	return strconv.FormatFloat(math.Min(v/whole*100, 100), 'f', 2, 64) + "%"
}

// chartLabelEvery returns how often to label the columns along the axis,
// so that at most about six labels share it.
func chartLabelEvery(n int) int {
	return max(1, int(math.Ceil(float64(n)/6)))
}

// chartLabelled says whether column i of n is labelled on the axis: every
// so often counting back from the latest, which is always labelled.
func chartLabelled(i, n int) bool {
	return (n-1-i)%chartLabelEvery(n) == 0
}

// sparklinePoints plots the values in a 100 by 32 box, for an SVG
// polyline stretched over its element.
func sparklinePoints(vs []float64) string {
	if len(vs) == 0 {
		return ""
	}
	lo, hi := vs[0], vs[0]
	for _, v := range vs {
		lo, hi = min(lo, v), max(hi, v)
	}
	// Counts start from zero, so their trend isn't exaggerated.
	lo = min(lo, 0)
	span := hi - lo
	pts := make([]string, len(vs))
	for i, v := range vs {
		x := 50.0
		if len(vs) > 1 {
			x = float64(i) / float64(len(vs)-1) * 100
		}
		// A flat line of nothing lies on the bottom, and of something
		// across the middle.
		y := 16.0
		if hi == 0 {
			y = 30
		}
		if span > 0 {
			// Two units of room keep the line's width inside the box.
			y = 30 - (v-lo)/span*28
		}
		pts[i] = fmt.Sprintf("%.2f,%.2f", x, y)
	}
	return strings.Join(pts, " ")
}

// sparklineArea closes the line along the bottom of the box, for the wash
// beneath it.
func sparklineArea(vs []float64) string {
	if len(vs) < 2 {
		return ""
	}
	return "0,32 " + sparklinePoints(vs) + " 100,32"
}

// metricChange writes a change as a signed percentage, "+12%" or "−3.5%".
func metricChange(c float64) string {
	pct := math.Round(c*1000) / 10
	sign := "+"
	if pct < 0 {
		sign = "−"
	}
	if pct == 0 {
		sign = ""
	}
	s := strconv.FormatFloat(math.Abs(pct), 'f', 1, 64)
	if math.Abs(pct) >= 10 {
		s = strconv.FormatFloat(math.Round(math.Abs(pct)), 'f', 0, 64)
	}
	return sign + strings.TrimSuffix(s, ".0") + "%"
}

// Tones of a metric's change.
const (
	toneGood = "good"
	toneBad  = "bad"
	toneFlat = "flat"
)

// metricTone says whether a change is good news, bad news or neither,
// for its color.
func metricTone(p props.Metric) string {
	pct := math.Round(p.Change * 1000)
	switch {
	case pct == 0:
		return toneFlat
	case (pct > 0) == (p.Direction != props.MetricUpBad):
		return toneGood
	}
	return toneBad
}

// chartPositive returns the values above zero, which are the ones drawn.
func chartPositive(values []props.ChartValue) []props.ChartValue {
	out := make([]props.ChartValue, 0, len(values))
	for _, v := range values {
		if v.Value > 0 {
			out = append(out, v)
		}
	}
	return out
}
