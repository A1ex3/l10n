package generatorjava

import (
	"bytes"
	"strconv"
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
	Type string // can be "String", "Float", or "Integer"
}

type generatorJava struct {
	localeData *localizationData
}

const localizationTemplate = `package {{ .PackageName }};

import java.text.MessageFormat;
import java.util.HashMap;
import java.util.Map;

interface {{.BaseClassName}} {
    {{- range .Methods }}
    String {{ .Name }}({{ range $i, $param := .Parameters }}{{ if $i }}, {{ end }}{{ $param.Type }} {{ $param.Name }}{{ end }});
    {{- end }}
}

{{- range .Languages }}

final class {{ .Name }} implements {{ $.BaseClassName }} {
    public final String LANGUAGE_CODE = "{{ .Code }}";

    {{- range $methodName, $translation := .Translations }}
    @Override
    public String {{ $methodName }}({{ $methodParams := index $.MethodsByName $methodName }}{{ range $i, $param := $methodParams.Parameters }}{{ if $i }}, {{ end }}{{ $param.Type }} {{ $param.Name }}{{ end }}) {
        {{- if $methodParams.Parameters }}
        return MessageFormat.format("{{ formatTranslation $translation $methodParams.Parameters }}", {{ range $i, $param := $methodParams.Parameters }}{{ if $i }}, {{ end }}{{ $param.Name }}{{ end }});
        {{- else }}
        return "{{ $translation }}";
        {{- end }}
    }
    {{- end }}
}
{{- end }}

public final class {{ .ClassName }} {
    public static final String DEFAULT_LANGUAGE_CODE = "{{ .DefaultLangCode }}";
    private static String currentLanguageCode = "{{ .CurrentLangCode }}";
    public static final Map<String, {{ .BaseClassName }}> languages = new HashMap<>();

    static {
        {{- range $key, $lang := .Languages }}
        languages.put("{{ $lang.Code }}", new {{ $lang.Name }}());
        {{- end }}
    }

    public static {{ .BaseClassName }} get() {
        return languages.get(currentLanguageCode);
    }

    public static boolean setCurrentLanguageCode(String languageCode) {
        if (languages.containsKey(languageCode)) {
            currentLanguageCode = languageCode;
            return true;
        } else {
            return false;
        }
    }

    public static String getCurrentLanguageCode() {
        return currentLanguageCode;
    }
}
`

func formatTranslation(translation string, parameters []parameter) string {
	// Replace variable names with {0}, {1}, ...
	for i, param := range parameters {
		translation = strings.ReplaceAll(translation, "{"+param.Name+"}", "{"+strconv.Itoa(i)+"}")
	}
	return translation
}

func (g *generatorJava) transform(className, packageName, defaultLanguageCode string, data *parser.Parser) {
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
						paramType = "Integer"
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

func (g *generatorJava) Get() (string, error) {
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

func NewGeneratorJava(className, packageName, defaultLanguageCode string, data *parser.Parser) *generatorJava {
	gp := &generatorJava{
		localeData: nil,
	}
	gp.transform(className, packageName, defaultLanguageCode, data)
	return gp
}
