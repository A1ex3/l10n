package generator

import (
	"fmt"
	"strings"

	generatorcpp "github.com/a1ex3/l10n/internal/generator/cpp"
	generatorjava "github.com/a1ex3/l10n/internal/generator/java"
	generatorjs "github.com/a1ex3/l10n/internal/generator/javascript"
	generatorkotlin "github.com/a1ex3/l10n/internal/generator/kotlin"
	generatorpython "github.com/a1ex3/l10n/internal/generator/python"
	generatorts "github.com/a1ex3/l10n/internal/generator/typescript"
	"github.com/a1ex3/l10n/internal/parser"
)

type Generator interface {
	Get() (string, error)
}

// Creates a generator for the requested programming language.
func NewGenerator(language, className, defaultLanguageCode, packageName string, data *parser.Parser) (Generator, error) {
	switch strings.ToLower(language) {
	case "python":
		return generatorpython.NewGeneratorPython(className, defaultLanguageCode, data), nil
	case "java":
		return generatorjava.NewGeneratorJava(className, packageName, defaultLanguageCode, data), nil
	case "cpp":
		return generatorcpp.NewGeneratorCpp(className, packageName, defaultLanguageCode, data), nil
	case "kotlin":
		return generatorkotlin.NewGeneratorKotlin(className, packageName, defaultLanguageCode, data), nil
	case "typescript":
		return generatorts.NewGeneratorTS(className, packageName, defaultLanguageCode, data), nil
	case "javascript":
		return generatorjs.NewGeneratorJS(className, packageName, defaultLanguageCode, data), nil
	default:
		return nil, fmt.Errorf("unsupported programming language %q", language)
	}
}
