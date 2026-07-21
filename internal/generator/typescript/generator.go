package generatorts

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/a1ex3/l10n/internal/generator/common"
	"github.com/a1ex3/l10n/internal/parser"
	"golang.org/x/text/cases"
	lang "golang.org/x/text/language"
)

type generatorTS struct {
	localeData *common.LocalizationData
}

const localizationTemplate = `// {{ .PackageName }}

interface {{ .BaseClassName }} {
    {{- range .Methods }}
    /**
     * Description: <b>{{ .Description }}</b>
     * Example: <b>{{ .Example }}</b>
     */
    {{ .Name }}({{ range $i, $param := .Parameters }}{{ if $i }}, {{ end }}{{ $param.Name }}: {{ $param.Type }}{{ end }}): string;
    {{- end }}
}

{{- range .Languages }}

class {{ .Name }} implements {{ $.BaseClassName }} {
    LANGUAGE_CODE: string = "{{ .Code }}";

    {{- range $methodName, $translation := .Translations }}
    {{ $methodName }}({{ $methodParams := index $.MethodsByName $methodName }}{{ range $i, $param := $methodParams.Parameters }}{{ if $i }}, {{ end }}{{ $param.Name }}: {{ $param.Type }}{{ end }}): string {
        {{- if $methodParams.Parameters }}
        return ` + "`" + `{{ formatTranslation $translation $methodParams.Parameters }}` + "`" + `;
        {{- else }}
        return "{{ escapeSymbolsString $translation }}";
        {{- end }}
    }
    {{- end }}
}
{{- end }}

class {{ .ClassName }} {
    static readonly DEFAULT_LANGUAGE_CODE: string = "{{ .DefaultLangCode }}";
    static currentLanguageCode: string = "{{ .CurrentLangCode }}";
    static readonly languages: Record<string, {{ .BaseClassName }}> = {
        {{- range $key, $lang := .Languages }}
        "{{ $lang.Code }}": new {{ $lang.Name }}(),
        {{- end }}
    };

    static get(): {{ .BaseClassName }} {
        return this.languages[this.currentLanguageCode];
    }

    static setCurrentLanguageCode(languageCode: string): boolean {
        if (this.languages.hasOwnProperty(languageCode)) {
            this.currentLanguageCode = languageCode;
            return true;
        } else {
            return false;
        }
    }

    static getCurrentLanguageCode(): string {
        return this.currentLanguageCode;
    }
}
`

func formatTranslation(translation string, parameters []common.Parameter) string {
	// Replace variable names with ${parameter} for TypeScript template literals
	for _, param := range parameters {
		translation = strings.ReplaceAll(translation, "{"+param.Name+"}", "${"+param.Name+"}")
	}
	return translation
}

func (g *generatorTS) transform(className, packageName, defaultLanguageCode string, data *parser.Parser) {
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
						paramType = "string"
					case parser.TypeFloatString:
						paramType = "number"
					case parser.TypeIntString:
						paramType = "number"
					default:
						paramType = "string"
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

func (g *generatorTS) Get() (string, error) {
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

func NewGeneratorTS(className, packageName, defaultLanguageCode string, data *parser.Parser) *generatorTS {
	gp := &generatorTS{
		localeData: nil,
	}
	gp.transform(className, packageName, defaultLanguageCode, data)
	return gp
}
