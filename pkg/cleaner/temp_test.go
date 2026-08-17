package cleaner

import (
	"testing"
)

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		input    int64
		expected string
	}{
		{500, "500 B"},
		{1024, "1.00 KB"},
		{1048576, "1.00 MB"},
		{1073741824, "1.00 GB"},
	}

	for _, tt := range tests {
		got := FormatBytes(tt.input)
		if got != tt.expected {
			t.Errorf("FormatBytes(%d) = %s; want %s", tt.input, got, tt.expected)
		}
	}
}

func TestCleanTempFiles(t *testing.T) {
	res, err := CleanTempFiles()
	if err != nil {
		t.Fatalf("CleanTempFiles falhou com erro: %v", err)
	}

	if res.FilesRemoved < 0 {
		t.Errorf("FilesRemoved não pode ser negativo: %d", res.FilesRemoved)
	}
}
