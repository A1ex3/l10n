package generatorcpp

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

type localizationData struct {
	Namespace       string
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
	Type string // can be "std::string", "float", or "int"
}

type generatorCpp struct {
	localeData *localizationData
}

const localizationTemplate = `#pragma once

#include <string>
#include <unordered_map>
#include <vector>
#include <sstream>

namespace {{ .Namespace }} {

class {{ .BaseClassName }} {
public:
    virtual ~{{ .BaseClassName }}() = default;

    {{- range .Methods }}
    /**
     * Description: {{ .Description }}
     * Example: {{ .Example }}
     */
    virtual std::string {{ .Name }}({{ range $i, $param := .Parameters }}{{ if $i }}, {{ end }}{{ $param.Type }} {{ $param.Name }}{{ end }}) const = 0;
    {{- end }}
};

{{- range .Languages }}

class {{ .Name }} : public {{ $.BaseClassName }} {
public:
    const std::string LANGUAGE_CODE = "{{ .Code }}";

    {{- range $methodName, $translation := .Translations }}
    std::string {{ $methodName }}({{ $methodParams := index $.MethodsByName $methodName }}{{ range $i, $param := $methodParams.Parameters }}{{ if $i }}, {{ end }}{{ $param.Type }} {{ $param.Name }}{{ end }}) const override {
        {{- if $methodParams.Parameters }}
        std::ostringstream oss;
        oss << "{{ formatTranslation $translation $methodParams.Parameters }}";
        return oss.str();
        {{- else }}
        return "{{ escapeSymbolsString $translation }}";
        {{- end }}
    }
    {{- end }}
};

{{- end }}

class {{ .ClassName }} {
public:
    static const std::string DEFAULT_LANGUAGE_CODE;
    static std::string currentLanguageCode;
    static std::unordered_map<std::string, {{ .BaseClassName }}*> languages;

    static {{ .BaseClassName }}* get() {
        return languages[currentLanguageCode];
    }

    static bool setCurrentLanguageCode(const std::string& languageCode) {
        if (languages.find(languageCode) != languages.end()) {
            currentLanguageCode = languageCode;
            return true;
        } else {
            return false;
        }
    }

    static std::string getCurrentLanguageCode() {
        return currentLanguageCode;
    }
};

const std::string {{ .ClassName }}::DEFAULT_LANGUAGE_CODE = "{{ .DefaultLangCode }}";
std::string {{ .ClassName }}::currentLanguageCode = "{{ .CurrentLangCode }}";
std::unordered_map<std::string, {{ .BaseClassName }}*> {{ .ClassName }}::languages = {
    {{- range $key, $lang := .Languages }}
    {"{{ $lang.Code }}", new {{ $lang.Name }}()},
    {{- end }}
};

}  // namespace {{ .Namespace }}
`

func formatTranslation(translation string, parameters []parameter) string {
	// Replace variable names with {0}, {1}, ...
	for i, param := range parameters {
		translation = strings.ReplaceAll(translation, "{"+param.Name+"}", "{"+strconv.Itoa(i)+"}")
	}
	return translation
}

func (g *generatorCpp) transform(className, namespaceName, defaultLanguageCode string, data *parser.Parser) {
	const baseClassNamePrefix = "Base"

	ld := &localizationData{
		Namespace:       namespaceName,
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
						paramType = "std::string"
					case parser.TypeFloatString:
						paramType = "float"
					case parser.TypeIntString:
						paramType = "int"
					default:
						paramType = "std::string"
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

func (g *generatorCpp) Get() (string, error) {
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

func NewGeneratorCpp(className, namespaceName, defaultLanguageCode string, data *parser.Parser) *generatorCpp {
	gp := &generatorCpp{
		localeData: nil,
	}
	gp.transform(className, namespaceName, defaultLanguageCode, data)
	return gp
}
