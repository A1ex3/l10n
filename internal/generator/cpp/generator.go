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

type generatorCpp struct {
	localeData *common.LocalizationData
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

func formatTranslation(translation string, parameters []common.Parameter) string {
	// Replace variable names with {0}, {1}, ...
	for i, param := range parameters {
		translation = strings.ReplaceAll(translation, "{"+param.Name+"}", "{"+strconv.Itoa(i)+"}")
	}
	return translation
}

func (g *generatorCpp) transform(className, namespaceName, defaultLanguageCode string, data *parser.Parser) {
	const baseClassNamePrefix = "Base"

	ld := &common.LocalizationData{
		Namespace:       namespaceName,
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
						paramType = "std::string"
					case parser.TypeFloatString:
						paramType = "float"
					case parser.TypeIntString:
						paramType = "int"
					default:
						paramType = "std::string"
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

func (g *generatorCpp) Get() (string, error) {
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

func NewGeneratorCpp(className, namespaceName, defaultLanguageCode string, data *parser.Parser) *generatorCpp {
	gp := &generatorCpp{}
	gp.transform(className, namespaceName, defaultLanguageCode, data)
	return gp
}
