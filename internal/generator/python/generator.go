package generatorpython

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/a1ex3/l10n/internal/parser"
	"golang.org/x/text/cases"
	lang "golang.org/x/text/language"
)

type localizationData struct {
	Languages       []language
	ClassName       string
	BaseClassName   string
	DefaultLangCode string
	CurrentLangCode string
	Methods         []method
	MethodsByName   map[string]method
}

type language struct {
	Name         string
	Code         string
	Translations map[string]string
}

type method struct {
	Name        string
	Parameters  []parameter
	Description string
	Example     string
}

type parameter struct {
	Name string
	Type string // can be "str", "float", or "int"
}

type generatorPython struct {
	localeData *localizationData
}

const localizationTemplate = `class {{.BaseClassName}}:
    {{- range .Methods }}
    def {{ .Name }}(self, {{ range $i, $param := .Parameters }}{{ if $i }}, {{ end }}{{ $param.Name }}: {{ $param.Type }}{{ end }}) -> str:
        raise NotImplementedError("Method must be implemented in subclass!")
    {{- end }}

{{- range .Languages }}

class {{ .Name }}({{ $.BaseClassName }}):
    def __init__(self) -> None:
        self.LANGUAGE_CODE: str = "{{ .Code }}"
    {{- range $methodName, $translation := .Translations }}
    def {{ $methodName }}(self, {{ $methodParams := index $.MethodsByName $methodName }}{{ range $i, $param := $methodParams.Parameters }}{{ if $i }}, {{ end }}{{ $param.Name }}: {{ $param.Type }}{{ end }}) -> str:
        """Description: {{ $methodParams.Description }}
        Example: {{ $methodParams.Example }}
        """
        return f"{{ $translation }}"
    {{- end }}
{{- end }}

class {{ .ClassName }}:
    DEFAULT_LANGUAGE_CODE: str = "{{ .DefaultLangCode }}"
    current_language_code: str = "{{ .CurrentLangCode }}"
    languages: dict[str, {{ .BaseClassName }}] = {
        {{- range $key, $lang := .Languages }}
        "{{ $lang.Code }}": {{ $lang.Name }}(),
        {{- end }}
    }

    @staticmethod
    def get() -> {{ .BaseClassName }}:
        if {{ .ClassName }}.current_language_code in {{ .ClassName }}.languages:
            return {{ .ClassName }}.languages[{{ .ClassName }}.current_language_code]
        else:
            raise NotImplementedError(f"Such localization does not exist: {{ .ClassName }}.current_language_code")
`

func capitalizeAfterHyphen(input string) string {
	parts := strings.Split(input, "-")
	var capitalizedParts []string

	for _, part := range parts {
		if part != "" {
			capitalized := cases.Title(lang.English).String(part)
			capitalizedParts = append(capitalizedParts, capitalized)
		}
	}
	return strings.Join(capitalizedParts, "")
}

func (g *generatorPython) transform(className, defaultLanguageCode string, data *parser.Parser) {
	const baseClassNamePrefix string = "Base"

	ld := &localizationData{
		ClassName:       className,
		BaseClassName:   baseClassNamePrefix + cases.Title(lang.English).String(className),
		DefaultLangCode: defaultLanguageCode,
		CurrentLangCode: defaultLanguageCode,
		MethodsByName:   make(map[string]method),
	}

	languages := make([]language, 0)
	methods := make(map[string]method)

	for _, loc := range data.ArrayOfLocalizations {
		for methodName, translateData := range loc.Data {
			parameters := make([]parameter, 0)
			description := ""
			example := ""
			if params := translateData.Params; params != nil {
				description = params.Description
				example = params.Example
				for _, variable := range params.Variables {
					// Determine parameter type based on variable type
					var paramType string
					switch variable.VariableType {
					case parser.TypeStringString:
						paramType = "str"
					case parser.TypeFloatString:
						paramType = "float"
					case parser.TypeIntString:
						paramType = "int"
					default:
						paramType = "str" // Default to string if type is unknown
					}
					parameters = append(parameters, parameter{
						Name: variable.VariableName,
						Type: paramType,
					})
				}
			}

			// Check if method exists and create if it doesn't
			if _, exists := methods[methodName]; !exists {
				m := method{
					Name:        methodName,
					Parameters:  parameters,
					Description: description,
					Example:     example,
				}
				methods[methodName] = m
			}
		}

		translations := make(map[string]string)
		for key, translateData := range loc.Data {
			translations[key] = translateData.Text
		}
		lang := language{
			Name:         className + capitalizeAfterHyphen(cases.Title(lang.English).String(loc.LanguageCode)),
			Code:         loc.LanguageCode,
			Translations: translations,
		}
		languages = append(languages, lang)
	}

	ld.Languages = languages

	for _, m := range methods {
		ld.Methods = append(ld.Methods, m)
		ld.MethodsByName[m.Name] = m
	}

	g.localeData = ld
}

func (g *generatorPython) Get() (string, error) {
	var tpl bytes.Buffer
	err := template.Must(template.New("localization").Parse(localizationTemplate)).Execute(&tpl, g.localeData)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}

func NewGeneratorPython(className, defaultLanguageCode string, data *parser.Parser) *generatorPython {
	gp := &generatorPython{
		localeData: nil,
	}
	gp.transform(className, defaultLanguageCode, data)
	return gp
}
