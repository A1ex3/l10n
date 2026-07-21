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

type generatorKotlin struct {
	localeData *common.LocalizationData
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

func formatTranslation(translation string, parameters []common.Parameter) string {
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

	ld := &common.LocalizationData{
		PackageName:     packageName,
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
					parameters = append(parameters, common.Parameter{
						Name: variable.VariableName,
						Type: paramType,
					})
				}
			}

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

func (g *generatorKotlin) Get() (string, error) {
	var tpl bytes.Buffer
	err := template.Must(template.New("localization").Funcs(template.FuncMap{
		"formatTranslation": func(translation string, parameters []common.Parameter) string {
			return formatTranslation(translation, parameters)
		},
		"escapeSymbolsString": common.EscapeSymbolsString,
	}).Parse(localizationTemplate)).Execute(&tpl, g.localeData)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}

func NewGeneratorKotlin(className, packageName, defaultLanguageCode string, data *parser.Parser) *generatorKotlin {
	gp := &generatorKotlin{
		localeData: nil,
	}
	gp.transform(className, packageName, defaultLanguageCode, data)
	return gp
}
