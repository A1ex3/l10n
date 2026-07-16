package generatorkotlin

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/a1ex3/l10n/internal/generator/common"
	"github.com/a1ex3/l10n/internal/parser"
	"golang.org/x/text/cases"
	lang "golang.org/x/text/language"
)

type localizationData struct {
	PackageName     string
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
	Type string // can be "String", "Float", or "Int"
}

type generatorKotlin struct {
	localeData *localizationData
}

const localizationTemplate = `package {{ .PackageName }}

interface {{ .BaseClassName }} {
    {{- range .Methods }}
    /**
     * Description: <b>{{ .Description }}</b>
     * Example: <b>{{ .Example }}</b>
     */
    fun {{ .Name }}({{ range $i, $param := .Parameters }}{{ if $i }}, {{ end }}{{ $param.Name }}: {{ $param.Type }}{{ end }}): String
    {{- end }}
}

{{- range .Languages }}

class {{ .Name }} : {{ $.BaseClassName }} {
    val LANGUAGE_CODE: String = "{{ .Code }}"

    {{- range $methodName, $translation := .Translations }}
    override fun {{ $methodName }}({{ $methodParams := index $.MethodsByName $methodName }}{{ range $i, $param := $methodParams.Parameters }}{{ if $i }}, {{ end }}{{ $param.Name }}: {{ $param.Type }}{{ end }}): String {
        {{- if $methodParams.Parameters }}
        return "{{ formatTranslation $translation $methodParams.Parameters }}".format({{ range $i, $param := $methodParams.Parameters }}{{ if $i }}, {{ end }}{{ $param.Name }}{{ end }})
        {{- else }}
        return "{{ escapeSymbolsString $translation }}"
        {{- end }}
    }
    {{- end }}
}
{{- end }}

object {{ .ClassName }} {
    const val DEFAULT_LANGUAGE_CODE = "{{ .DefaultLangCode }}"
    private var currentLanguageCode = "{{ .CurrentLangCode }}"
    private val languages = mapOf(
        {{- range $key, $lang := .Languages }}
        "{{ $lang.Code }}" to {{ $lang.Name }}(),
        {{- end }}
    )

    fun get(): {{ .BaseClassName }}? {
        return languages[currentLanguageCode]
    }

    fun setCurrentLanguageCode(languageCode: String): Boolean {
        return if (languages.containsKey(languageCode)) {
            currentLanguageCode = languageCode
            true
        } else {
            false
        }
    }

    fun getCurrentLanguageCode(): String {
        return currentLanguageCode
    }
}
`

func formatTranslation(translation string, parameters []parameter) string {
	// Replace variable names with %s, %f, %d, etc., depending on the parameter type
	for _, param := range parameters {
		placeholder := "%s" // default
		switch param.Type {
		case "Float":
			placeholder = "%f"
		case "Int":
			placeholder = "%d"
		}
		translation = strings.ReplaceAll(translation, "{"+param.Name+"}", placeholder)
	}
	return translation
}

func (g *generatorKotlin) transform(className, packageName, defaultLanguageCode string, data *parser.Parser) {
	const baseClassNamePrefix = "Base"

	ld := &localizationData{
		PackageName:     packageName,
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
					var paramType string
					switch variable.VariableType {
					case parser.TypeStringString:
						paramType = "String"
					case parser.TypeFloatString:
						paramType = "Float"
					case parser.TypeIntString:
						paramType = "Int"
					default:
						paramType = "String"
					}
					parameters = append(parameters, parameter{
						Name: variable.VariableName,
						Type: paramType,
					})
				}
			}

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

func (g *generatorKotlin) Get() (string, error) {
	var tpl bytes.Buffer
	err := template.Must(template.New("localization").Funcs(template.FuncMap{
		"formatTranslation": func(translation string, parameters []parameter) string {
			return formatTranslation(translation, parameters)
		},
		"escapeSymbolsString": common.EscapeSymbolsString,
	}).Parse(localizationTemplate)).Execute(&tpl, g.localeData)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}

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

func NewGeneratorKotlin(className, packageName, defaultLanguageCode string, data *parser.Parser) *generatorKotlin {
	gp := &generatorKotlin{
		localeData: nil,
	}
	gp.transform(className, packageName, defaultLanguageCode, data)
	return gp
}
