package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/a1ex3/l10n/internal/parser"
)

func createTempFile(t *testing.T, dir, filename string, content []byte) {
	t.Helper()
	filePath := filepath.Join(dir, filename)
	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file %s: %v", filename, err)
	}
}

func TestGeneratorOnEscapeSymbols(t *testing.T) {
	dir := t.TempDir()

	enContent := []byte(`{"helloWorld": "\tHello\nWorld", "welcomeUserOne": "Welcome user \"one\""}`)

	createTempFile(t, dir, "l10n_en.json", enContent)

	parser, err := parser.NewParser(dir, "en")
	if err != nil {
		t.Fatalf("Failed to create new parser: %v", err)
	}

	enLocale := parser.ArrayOfLocalizations[0]
	if enLocale.LanguageCode != "en" {
		t.Errorf("Expected language code 'en', got '%s'", enLocale.LanguageCode)
	}

	gen, err := NewGenerator("python", "AppLocalization", "en", "", parser)

	if err != nil {
		t.Errorf("Generator error: %v", err)
	}

	gen_data, err := gen.Get()

	if err != nil {
		t.Errorf("Generator error (Get()): %v", err)
	}

	if len(gen_data) == 0 {
		t.Errorf("Generated code length is 0!")
	}

	substrs := []string{`"\tHello\nWorld"`, `"Welcome user \"one\""`}

	for _, substr := range substrs {
		if !strings.Contains(gen_data, substr) {
			t.Errorf("Generated code doesn't contain: %s", substr)
		}
	}
}
