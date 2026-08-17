package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("User home dir não disponível para teste")
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"/var/log", "/var/log"},
		{"~", home},
		{"~/test", filepath.Join(home, "test")},
		{`~\test`, filepath.Join(home, "test")},
	}

	for _, tt := range tests {
		got := ExpandHome(tt.input)
		if got != tt.expected {
			t.Errorf("ExpandHome(%q) = %q; esperado %q", tt.input, got, tt.expected)
		}
	}
}

func TestGetAuthor(t *testing.T) {
	os.Setenv("TOMMY_AUTHOR", "TestAuthor")
	defer os.Unsetenv("TOMMY_AUTHOR")

	if got := GetAuthor(); got != "TestAuthor" {
		t.Errorf("GetAuthor() = %q; esperado 'TestAuthor'", got)
	}
}

func TestGetOutputDir(t *testing.T) {
	os.Setenv("TOMMY_OUTPUT_DIR", "./CustomOutput")
	defer os.Unsetenv("TOMMY_OUTPUT_DIR")

	if got := GetOutputDir(); got != "./CustomOutput" {
		t.Errorf("GetOutputDir() = %q; esperado './CustomOutput'", got)
	}
}

func TestGetTempDirs(t *testing.T) {
	os.Setenv("TOMMY_TEMP_DIR", "/custom/temp1, /custom/temp2")
	defer os.Unsetenv("TOMMY_TEMP_DIR")

	dirs := GetTempDirs()
	if len(dirs) != 2 || dirs[0] != "/custom/temp1" || dirs[1] != "/custom/temp2" {
		t.Errorf("GetTempDirs() = %v; esperado [/custom/temp1 /custom/temp2]", dirs)
	}
}
