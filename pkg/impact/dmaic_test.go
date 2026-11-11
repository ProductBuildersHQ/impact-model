package impact

import (
	"fmt"
	"testing"
)

func TestTimesPerYear(t *testing.T) {
	tests := []struct {
		iso  string
		want float64
		ok   bool
	}{
		{"P1Y", 1, true},
		{"P1M", 12, true},
		{"P3M", 4, true},
		{"P1W", 52.18, true},
		{"P2W", 26.09, true},
		{"P1D", 365.24, true},
		{"PT12H", 730.49, true},
		{"", 0, false},
		{"P0D", 0, false},
		{"monthly", 0, false},
	}
	for _, tt := range tests {
		got, ok := TimesPerYear(tt.iso)
		if ok != tt.ok || (ok && fmt.Sprintf("%.2f", got) != fmt.Sprintf("%.2f", tt.want)) {
			t.Errorf("TimesPerYear(%q) = %.2f, %v; want %.2f, %v", tt.iso, got, ok, tt.want, tt.ok)
		}
	}
}
