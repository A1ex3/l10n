package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func createTempFile(t *testing.T, dir, filename string, content []byte) {
	t.Helper()
	filePath := filepath.Join(dir, filename)
	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file %s: %v", filename, err)
	}
}

func TestReadLocales(t *testing.T) {
	dir := t.TempDir()

	enContent := []byte(`{"hello_world": "Hello, World!", "#hello_world": {"description": "Greeting message", "example": "Hello, World!"}}`)
	frContent := []byte(`{"hello_world": "Bonjour le monde!", "#hello_world": {"description": "Message de salutation", "example": "Bonjour le monde!"}}`)

	createTempFile(t, dir, "l10n_en.json", enContent)
	createTempFile(t, dir, "l10n_fr.json", frContent)

	p := &Parser{}
	locales, err := p.readLocales(dir)
	if err != nil {
		t.Fatalf("Failed to read locales: %v", err)
	}

	if len(locales) != 2 {
		t.Errorf("Expected 2 locales, got %d", len(locales))
	}

	if _, ok := locales["en"]; !ok {
		t.Error("Locale 'en' not found")
	}

	if _, ok := locales["fr"]; !ok {
		t.Error("Locale 'fr' not found")
	}
}

func TestTemplateLanguageCodeInLocales(t *testing.T) {
	data := map[string]map[string]interface{}{
		"en": {"hello_world": "Hello, World!"},
		"fr": {"hello_world": "Bonjour le monde!"},
	}

	p := &Parser{}

	if !p.templateLanguageCodeInLocales("en", data) {
		t.Error("Expected 'en' template language code to be found")
	}

	if p.templateLanguageCodeInLocales("es", data) {
		t.Error("Did not expect 'es' template language code to be found")
	}
}

func TestDefineVariableType(t *testing.T) {
	p := &Parser{}

	tests := []struct {
		value       interface{}
		params      []interface{}
		expected    typeVariable
		expectError bool
	}{
		{"string", nil, TypeStringString, false},
		{123, nil, TypeIntString, false},
		{123.45, nil, TypeFloatString, false},
		{"string", []interface{}{TypeStringString}, TypeStringString, false},
		{123.45, []interface{}{TypeIntString}, TypeIntString, false},
		{true, nil, "", true},
	}

	for _, test := range tests {
		dataType, ok := p.defineVariableType(test.value, test.params...)
		if ok != !test.expectError {
			t.Errorf("Expected error: %v, got: %v", test.expectError, !ok)
		}

		if dataType != test.expected {
			t.Errorf("Expected type: %v, got: %v", test.expected, dataType)
		}
	}
}

func TestParse(t *testing.T) {
	data := map[string]map[string]interface{}{
		"en": {
			"hello_world": "Hello, World!",
			"#hello_world": map[string]interface{}{
				"description": "Greeting message",
				"example":     "Hello, World!",
				"variables": map[string]interface{}{
					"name": map[string]interface{}{
						"type":         "string",
						"defaultValue": "World",
					},
				},
			},
		},
	}

	p := &Parser{}
	if err := p.parse(data); err != nil {
		t.Fatalf("Failed to parse data: %v", err)
	}

	if len(p.ArrayOfLocalizations) != 1 {
		t.Errorf("Expected 1 localization, got %d", len(p.ArrayOfLocalizations))
	}

	locale := p.ArrayOfLocalizations[0]
	if locale.LanguageCode != "en" {
		t.Errorf("Expected language code 'en', got '%s'", locale.LanguageCode)
	}

	translateData := locale.Data["hello_world"]
	if translateData.Text != "Hello, World!" {
		t.Errorf("Expected translation text 'Hello, World!', got '%s'", translateData.Text)
	}

	if translateData.Params == nil {
		t.Fatalf("Expected translation params, got nil")
	}

	if translateData.Params.Description != "Greeting message" {
		t.Errorf("Expected description 'Greeting message', got '%s'", translateData.Params.Description)
	}

	if len(translateData.Params.Variables) != 1 {
		t.Fatalf("Expected 1 variable, got %d", len(translateData.Params.Variables))
	}

	variable := translateData.Params.Variables[0]
	if variable.VariableName != "name" {
		t.Errorf("Expected variable name 'name', got '%s'", variable.VariableName)
	}

	if variable.VariableType != TypeStringString {
		t.Errorf("Expected variable type 'string', got '%s'", variable.VariableType)
	}

	if variable.DefaultValue != "World" {
		t.Errorf("Expected default value 'World', got '%v'", variable.DefaultValue)
	}
}

func TestEqualization(t *testing.T) {
	data := map[string]map[string]interface{}{
		"en": {"hello_world": "Hello, World!"},
		"fr": {},
	}

	p := &Parser{}
	if err := p.parse(data); err != nil {
		t.Fatalf("Failed to parse data: %v", err)
	}

	p.equalization("en")

	frLocale := p.ArrayOfLocalizations[1]
	if _, ok := frLocale.Data["hello_world"]; !ok {
		t.Error("Expected 'hello_world' key in French locale")
	}
}

func TestNewParser(t *testing.T) {
	dir := t.TempDir()

	enContent := []byte(`{"hello_world": "Hello, World!", "#hello_world": {"description": "Greeting message", "example": "Hello, World!"}}`)
	frContent := []byte(`{"hello_world": "Bonjour le monde!", "#hello_world": {"description": "Message de salutation", "example": "Bonjour le monde!"}}`)

	createTempFile(t, dir, "l10n_en.json", enContent)
	createTempFile(t, dir, "l10n_fr.json", frContent)

	parser, err := NewParser(dir, "en")
	if err != nil {
		t.Fatalf("Failed to create new parser: %v", err)
	}

	if len(parser.ArrayOfLocalizations) != 2 {
		t.Errorf("Expected 2 localizations, got %d", len(parser.ArrayOfLocalizations))
	}

	enLocale := parser.ArrayOfLocalizations[0]
	if enLocale.LanguageCode != "en" {
		t.Errorf("Expected language code 'en', got '%s'", enLocale.LanguageCode)
	}

	frLocale := parser.ArrayOfLocalizations[1]
	if frLocale.LanguageCode != "fr" {
		t.Errorf("Expected language code 'fr', got '%s'", frLocale.LanguageCode)
	}

	if _, ok := frLocale.Data["hello_world"]; !ok {
		t.Error("Expected 'hello_world' key in French locale")
	}
}
