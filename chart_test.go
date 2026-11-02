package popui

import (
	"bytes"
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/invopop/popui.go/props"
)

func TestChartTicks(t *testing.T) {
	tests := []struct {
		max   float64
		top   float64
		ticks []float64
	}{
		{0, 1, []float64{0}},
		{1, 1, []float64{0, 1}},
		{7, 10, []float64{0, 5, 10}},
		{10, 10, []float64{0, 5, 10}},
		{95, 100, []float64{0, 50, 100}},
		{1284, 1500, []float64{0, 500, 1000, 1500}},
		{0.4, 0.4, []float64{0, 0.2, 0.4}},
	}
	for _, tt := range tests {
		top, ticks := chartTicks(tt.max)
		if top != tt.top || !slices.EqualFunc(ticks, tt.ticks, func(a, b float64) bool { return a-b < 1e-9 && b-a < 1e-9 }) {
			t.Errorf("chartTicks(%v) = %v, %v, want %v, %v", tt.max, top, ticks, tt.top, tt.ticks)
		}
	}
}

func TestFormatNumber(t *testing.T) {
	for v, want := range map[float64]string{
		0: "0", 7: "7", 1284: "1,284", 1234567: "1,234,567", 12.5: "12.5", -4200: "-4,200", 3.04: "3",
	} {
		if got := FormatNumber(v); got != want {
			t.Errorf("FormatNumber(%v) = %q, want %q", v, got, want)
		}
	}
	for v, want := range map[float64]string{
		9999: "9,999", 12_900: "12.9K", 150_000: "150K", 4_200_000: "4.2M", 20_000: "20K",
	} {
		if got := CompactNumber(v); got != want {
			t.Errorf("CompactNumber(%v) = %q, want %q", v, got, want)
		}
	}
}

func TestChartLabelled(t *testing.T) {
	var got []int
	for i := range 13 {
		if chartLabelled(i, 13) {
			got = append(got, i)
		}
	}
	if want := []int{0, 3, 6, 9, 12}; !slices.Equal(got, want) {
		t.Errorf("labelled %v, want %v", got, want)
	}
}

func TestMetricChange(t *testing.T) {
	tests := []struct {
		p    props.Metric
		text string
		tone string
	}{
		{props.Metric{Change: 0.123}, "+12%", toneGood},
		{props.Metric{Change: -0.035}, "−3.5%", toneBad},
		{props.Metric{Change: 0.05, Direction: props.MetricUpBad}, "+5%", toneBad},
		{props.Metric{Change: -0.5, Direction: props.MetricUpBad}, "−50%", toneGood},
		{props.Metric{Change: 0.0001}, "0%", toneFlat},
	}
	for _, tt := range tests {
		if got := metricChange(tt.p.Change); got != tt.text {
			t.Errorf("metricChange(%v) = %q, want %q", tt.p.Change, got, tt.text)
		}
		if got := metricTone(tt.p); got != tt.tone {
			t.Errorf("metricTone(%v) = %q, want %q", tt.p.Change, got, tt.tone)
		}
	}
}

func TestColumnChartRender(t *testing.T) {
	var buf bytes.Buffer
	err := ColumnChart(props.ColumnChart{
		Label: "Invoices",
		Columns: []props.ChartValue{
			{Label: "Oct", Title: "October 2026", Value: 40, Href: "/oct"},
			{Label: "Nov", Title: "November 2026", Value: 10, Partial: true},
		},
	}).Render(context.Background(), &buf)
	if err != nil {
		t.Fatal(err)
	}
	html := buf.String()
	for _, want := range []string{`href="/oct"`, `aria-label="October 2026: 40"`, "height: 25.00%", "opacity-50", "<caption>Invoices</caption>"} {
		if !strings.Contains(html, want) {
			t.Errorf("chart is missing %q", want)
		}
	}
}
