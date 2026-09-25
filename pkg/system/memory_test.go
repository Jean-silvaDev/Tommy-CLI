package system

import (
	"strings"
	"testing"
)

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		input    uint64
		expected string
	}{
		{500, "500 B"},
		{1024, "1.00 KB"},
		{1048576, "1.00 MB"},
		{1073741824, "1.00 GB"},
		{1610612736, "1.50 GB"},
		{1099511627776, "1.00 TB"},
	}

	for _, tt := range tests {
		got := FormatBytes(tt.input)
		if got != tt.expected {
			t.Errorf("FormatBytes(%d) = %q; esperado %q", tt.input, got, tt.expected)
		}
	}
}

func TestCalculatePercent(t *testing.T) {
	tests := []struct {
		used     uint64
		total    uint64
		expected float64
	}{
		{0, 100, 0.0},
		{50, 100, 50.0},
		{75, 100, 75.0},
		{100, 100, 100.0},
		{0, 0, 0.0},
		{150, 100, 100.0},
	}

	for _, tt := range tests {
		got := CalculatePercent(tt.used, tt.total)
		if got != tt.expected {
			t.Errorf("CalculatePercent(%d, %d) = %f; esperado %f", tt.used, tt.total, got, tt.expected)
		}
	}
}

func TestRenderProgressBar(t *testing.T) {
	bar := RenderProgressBar(50.0, 10)
	if !strings.Contains(bar, "50.0%") {
		t.Errorf("RenderProgressBar esperado conter '50.0%%', obtido: %s", bar)
	}
}
