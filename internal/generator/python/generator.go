package generatorpython

import (
	"bytes"
	"text/template"

	"github.com/a1ex3/l10n/internal/generator/common"
	"github.com/a1ex3/l10n/internal/parser"
	"golang.org/x/text/cases"
	lang "golang.org/x/text/language"
)

type generatorPython struct {
	localeData *common.LocalizationData
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
        return f"{{ escapeSymbolsString $translation }}"
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
        return {{ .ClassName }}.languages[{{ .ClassName }}.current_language_code]
`

func (g *generatorPython) transform(className, defaultLanguageCode string, data *parser.Parser) {
	const baseClassNamePrefix string = "Base"

	ld := &common.LocalizationData{
		ClassName:       className,
		BaseClassName:   baseClassNamePrefix + cases.Title(lang.English).String(className),
		DefaultLangCode: defaultLanguageCode,
		CurrentLangCode: defaultLanguageCode,
		MethodsByName:   make(map[string]common.Method),
	}

	languages := make([]common.Language, 0)
	methods := make(map[string]common.Method)

	for _, loc := range data.ArrayOfLocalizations {
		for methodName, translateData := range loc.Data {
			parameters := make([]common.Parameter, 0)
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
					parameters = append(parameters, common.Parameter{
						Name: variable.VariableName,
						Type: paramType,
					})
				}
			}

			// Check if method exists and create if it doesn't
			if _, exists := methods[methodName]; !exists {
				m := common.Method{
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
		lang := common.Language{
			Name:         className + common.CapitalizeAfterHyphen(cases.Title(lang.English).String(loc.LanguageCode)),
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
	err := template.Must(template.New("localization").Funcs(template.FuncMap{
		"escapeSymbolsString": common.EscapeSymbolsString,
	}).Parse(localizationTemplate)).Execute(&tpl, g.localeData)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}

func NewGeneratorPython(className, defaultLanguageCode string, data *parser.Parser) *generatorPython {
	gp := &generatorPython{}
	gp.transform(className, defaultLanguageCode, data)
	return gp
}
