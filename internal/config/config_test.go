package config

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestArgsValidator(t *testing.T) {
	c := NewConfig([]string{"go", "python", "java"})

	tests := []struct {
		dir                    string
		template               string
		outputLocalizationFile string
		programmingLanguage    string
		className              string
		expectedError          string
	}{
		{"", "template", "output", "go", "ClassName", "dir cannot be empty"},
		{"dir", "", "output", "go", "ClassName", "template cannot be empty"},
		{"dir", "template", "", "go", "ClassName", "output_localization_file cannot be empty"},
		{"dir", "template", "output", "", "ClassName", "programming_language cannot be empty"},
		{"dir", "template", "output", "ruby", "ClassName", "invalid programming_language"},
		{"dir", "template", "output", "go", "", "class_name cannot be empty"},
		{"dir", "template", "output", "go", "ClassName", ""},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			err := c.argsValidator(tt.dir, tt.template, tt.outputLocalizationFile, tt.programmingLanguage, tt.className)
			if (err == nil && tt.expectedError != "") || (err != nil && err.Error() != tt.expectedError) {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestFromArgs(t *testing.T) {
	c := NewConfig([]string{"go", "python", "java"})
	err := c.FromArgs("dir", "template", "output", "go", "ClassName")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	cfg := c.GetConfig()
	if cfg.Dir != "dir" || cfg.Template != "template" || cfg.OutputLocalizationFile != "output" || cfg.ProgrammingLanguage != "go" || cfg.ClassName != "ClassName" {
		t.Errorf("config values not set correctly: %+v", cfg)
	}
}

func TestFromArgsInvalid(t *testing.T) {
	c := NewConfig([]string{"go", "python", "java"})
	err := c.FromArgs("", "template", "output", "go", "ClassName")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestFromFileYaml(t *testing.T) {
	c := NewConfig([]string{"go", "python", "java"})

	configData := entityConfig{
		Dir:                    "dir",
		Template:               "template",
		OutputLocalizationFile: "output",
		ProgrammingLanguage:    "go",
		ClassName:              "ClassName",
	}

	data, err := yaml.Marshal(&configData)
	if err != nil {
		t.Fatalf("failed to marshal yaml: %v", err)
	}

	filepath := "config_test.yaml"
	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	defer os.Remove(filepath)

	err = c.FromFileYaml(filepath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	cfg := c.GetConfig()
	if cfg.Dir != "dir" || cfg.Template != "template" || cfg.OutputLocalizationFile != "output" || cfg.ProgrammingLanguage != "go" || cfg.ClassName != "ClassName" {
		t.Errorf("config values not set correctly: %+v", cfg)
	}
}

func TestFromFileYamlInvalid(t *testing.T) {
	c := NewConfig([]string{"go", "python", "java"})

	invalidYaml := `
dir: dir
template: template
output_localization_file: output
programming_language: invalid_language
class_name: ClassName
`
	filepath := "invalid_config_test.yaml"
	err := os.WriteFile(filepath, []byte(invalidYaml), 0644)
	if err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	defer os.Remove(filepath)

	err = c.FromFileYaml(filepath)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
