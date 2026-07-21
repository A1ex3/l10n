package generatorjava

import (
	"bytes"
	"strconv"
	"strings"
	"text/template"

	"github.com/a1ex3/l10n/internal/generator/common"
	"github.com/a1ex3/l10n/internal/parser"
	"golang.org/x/text/cases"
	lang "golang.org/x/text/language"
)

type generatorJava struct {
	localeData *common.LocalizationData
}

const localizationTemplate = `package {{ .PackageName }};

import java.text.MessageFormat;
import java.util.HashMap;
import java.util.Map;

interface {{.BaseClassName}} {
    {{- range .Methods }}
    /**
     * Description: <b> {{ .Description }} </b>
     * Example: <b> {{ .Example }} </b>
     */
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
        return "{{ escapeSymbolsString $translation }}";
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

func formatTranslation(translation string, parameters []common.Parameter) string {
	// Replace variable names with {0}, {1}, ...
	for i, param := range parameters {
		translation = strings.ReplaceAll(translation, "{"+param.Name+"}", "{"+strconv.Itoa(i)+"}")
	}
	return translation
}

func (g *generatorJava) transform(className, packageName, defaultLanguageCode string, data *parser.Parser) {
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
						paramType = "Integer"
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

func (g *generatorJava) Get() (string, error) {
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

func NewGeneratorJava(className, packageName, defaultLanguageCode string, data *parser.Parser) *generatorJava {
	gp := &generatorJava{
		localeData: nil,
	}
	gp.transform(className, packageName, defaultLanguageCode, data)
	return gp
}
