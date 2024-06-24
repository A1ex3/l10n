package generatorpython

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/a1ex3/l10n/internal/parser"
)

func TestGeneratorPython(t *testing.T) {
	// Prepare test data
	pathToTranslates := t.TempDir()
	templateLanguageCode := "en"

	// Create temporary JSON files for testing
	createTestJSONFile(t, pathToTranslates, "l10n_en.json", `{
		"greet": "Hello, {name}!",
		"#greet": {
			"variables": {
				"name": {
					"type": "string"
				}
			}
		}
	}`)
	createTestJSONFile(t, pathToTranslates, "l10n_es.json", `{
		"greet": "Hola, {name}!",
		"#greet": {
			"variables": {
				"name": {
					"type": "string"
				}
			}
		}
	}`)

	// Create a new parser and parse the data
	p, err := parser.NewParser(pathToTranslates, templateLanguageCode)
	if err != nil {
		t.Fatalf("Unexpected error while parsing: %v", err)
	}

	// Create Python generator
	className := "Greeting"
	defaultLangCode := "en"
	generator := NewGeneratorPython(className, defaultLangCode, p)

	// Get generated code
	result, err := generator.Get()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Expected results
	expectedBaseClass := `class BaseGreeting:
    @staticmethod
    def greet(name: str) -> str:
        raise NotImplementedError("Method must be implemented in subclass!")
`

	expectedEnClass := `class GreetingEn(BaseGreeting):
    LANGUAGE_CODE: str = "en"

    @staticmethod
    def greet(name: str) -> str:
        return f"Hello, {name}!"
`

	expectedEsClass := `class GreetingEs(BaseGreeting):
    LANGUAGE_CODE: str = "es"

    @staticmethod
    def greet(name: str) -> str:
        return f"Hola, {name}!"
`

	expectedMainClass := `class Greeting:
    DEFAULT_LANGUAGE_CODE: str = "en"
    current_language_code: str = "en"
    language_codes: list[str] = ["en", "es"]
    languages: dict[str, BaseGreeting] = {
        "en": GreetingEn(),
        "es": GreetingEs(),
    }

    @staticmethod
    def get() -> BaseGreeting:
        if Greeting.current_language_code in Greeting.languages:
            return Greeting.languages[Greeting.current_language_code]
        else:
            raise NotImplementedError(f"Such localization does not exist: { Greeting.current_language_code }")
`

	expectedResult := expectedBaseClass + "\n" + expectedEnClass + "\n" + expectedEsClass + "\n" + expectedMainClass

	// Remove all whitespace for comparison
	cleanedExpectedResult := strings.ReplaceAll(expectedResult, " ", "")
	cleanedExpectedResult = strings.ReplaceAll(cleanedExpectedResult, "\n", "")
	cleanedResult := strings.ReplaceAll(result, " ", "")
	cleanedResult = strings.ReplaceAll(cleanedResult, "\n", "")

	if cleanedResult != cleanedExpectedResult {
		t.Errorf("Generated code does not match expected result.\nExpected:\n%s\nGot:\n%s", expectedResult, result)
	}
}

func createTestJSONFile(t *testing.T, dir, filename, content string) {
	t.Helper()
	fullPath := filepath.Join(dir, filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("Unable to create test directory: %v", err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		t.Fatalf("Unable to write test file: %v", err)
	}
}
