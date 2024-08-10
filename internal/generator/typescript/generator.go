package generatorts

import (
	"bytes"
	"strings"
	"text/template"

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
	Type string // can be "string", "number"
}

type generatorTS struct {
	localeData *localizationData
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
        return "{{ $translation }}";
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

func formatTranslation(translation string, parameters []parameter) string {
	// Replace variable names with ${parameter} for TypeScript template literals
	for _, param := range parameters {
		translation = strings.ReplaceAll(translation, "{"+param.Name+"}", "${"+param.Name+"}")
	}
	return translation
}

func (g *generatorTS) transform(className, packageName, defaultLanguageCode string, data *parser.Parser) {
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
						paramType = "string"
					case parser.TypeFloatString:
						paramType = "number"
					case parser.TypeIntString:
						paramType = "number"
					default:
						paramType = "string"
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

func (g *generatorTS) Get() (string, error) {
	var tpl bytes.Buffer
	err := template.Must(template.New("localization").Funcs(template.FuncMap{
		"formatTranslation": func(translation string, parameters []parameter) string {
			return formatTranslation(translation, parameters)
		},
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

func NewGeneratorTS(className, packageName, defaultLanguageCode string, data *parser.Parser) *generatorTS {
	gp := &generatorTS{
		localeData: nil,
	}
	gp.transform(className, packageName, defaultLanguageCode, data)
	return gp
}
