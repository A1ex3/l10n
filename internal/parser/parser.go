package parser

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// typeVariable represents the type of a variable.
type typeVariable string

const (
	TypeStringString typeVariable = "string"
	TypeIntString    typeVariable = "int"
	TypeFloatString  typeVariable = "float"
)

// translateParamsVariables represents variables within translation parameters.
type translateParamsVariables struct {
	VariableName string       // Name of the variable
	VariableType typeVariable // Type of the variable (string, int, float)
	DefaultValue interface{}  // Default value of the variable
}

// translateParams represents parameters for translation, including description, example, and variables.
type translateParams struct {
	Description string                     // Description of the translation
	Example     string                     // Example of the translation
	Variables   []translateParamsVariables // Variables used in the translation
}

// translateData represents the actual translation text and its associated parameters.
type translateData struct {
	Text   string           // Actual translated text
	Params *translateParams // Optional parameters for the translation
}

// localizationData represents all translations for a specific language.
type localizationData struct {
	LanguageCode string                   // Language code (e.g., "en", "fr")
	Data         map[string]translateData // Map of translation IDs to translateData
}

// Parser manages the parsing and handling of localization data.
type Parser struct {
	ArrayOfLocalizations []localizationData // Array of all localization data for different languages
}

// readLocales reads localization JSON files from a directory.
// It returns a map of language codes to translation data and any encountered error.
func (p *Parser) readLocales(pathToTranslates string) (map[string]map[string]interface{}, error) {
	locales := make(map[string]map[string]interface{})
	files, err := os.ReadDir(pathToTranslates)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filename := file.Name()
		if !strings.HasPrefix(filename, "l10n_") || !strings.HasSuffix(filename, ".json") {
			continue
		}
		filePath := filepath.Join(pathToTranslates, filename)
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		var jsonContent map[string]interface{}
		if err := json.Unmarshal(data, &jsonContent); err != nil {
			return nil, err
		}
		languageCode := strings.TrimSuffix(strings.TrimPrefix(filename, "l10n_"), ".json")
		locales[languageCode] = jsonContent
	}

	return locales, nil
}

// templateLanguageCodeInLocales checks if a given template language code exists in the provided data.
// It returns true if found, false otherwise.
func (p *Parser) templateLanguageCodeInLocales(templateLanguageCode string, data map[string]map[string]interface{}) bool {
	_, ok := data[templateLanguageCode]
	return ok
}

// defineVariableType determines the type of a given value based on optional parameters.
// It returns the determined type and a boolean indicating success.
//
// params(...interface{}) - [0](typeVariable): Expected type
func (p *Parser) defineVariableType(value interface{}, params ...interface{}) (typeVariable, bool) {
	lenParams := len(params)

	if lenParams == 0 {
		switch value.(type) {
		case string:
			return TypeStringString, true
		case int:
			return TypeIntString, true
		case float32, float64:
			return TypeFloatString, true
		default:
			return "", false
		}
	} else if lenParams == 1 {
		expectedType := params[0].(typeVariable)

		if reflect.TypeOf(value) == reflect.TypeOf(0.0) && expectedType == TypeIntString {
			value = int(value.(float64))
		}

		if dataType, ok := p.defineVariableType(value); ok && dataType == expectedType {
			return dataType, ok
		} else {
			return "", false
		}
	} else {
		log.Fatalf("Exceeded expected number of parameters! Expected: %d, got: %d", 1, lenParams)
		return "", false
	}
}

// parse processes the provided translation data and populates p.ArrayOfLocalizations.
// It returns an error if any issue occurs during parsing.
func (p *Parser) parse(data map[string]map[string]interface{}) error {
	const translateParamsSymbol string = "#"

	p.ArrayOfLocalizations = make([]localizationData, 0)

	for langCode, translates := range data {
		locale := localizationData{
			LanguageCode: langCode,
			Data:         make(map[string]translateData),
		}

		for id, data_ := range translates {
			if !strings.HasPrefix(id, translateParamsSymbol) {
				if _, okDefineVariableType := p.defineVariableType(data_, TypeStringString); okDefineVariableType {
					translateData := translateData{
						Text: data_.(string),
					}

					if translateParams_, okTranslateParams := translates[translateParamsSymbol+id]; okTranslateParams {
						translateParams := &translateParams{}
						params := translateParams_.(map[string]interface{})

						if description, okDescription := params["description"]; okDescription {
							translateParams.Description = description.(string)
						}

						if example, okExample := params["example"]; okExample {
							translateParams.Example = example.(string)
						}

						if variables, okVariables := params["variables"]; okVariables {
							vars := variables.(map[string]interface{})

							for varName, varParams_ := range vars {
								varParams := varParams_.(map[string]interface{})

								variable := translateParamsVariables{VariableName: varName}
								var currentVarTypeName typeVariable = ""

								if varTypeName, ok := varParams["type"]; ok {
									varTypeNameStr := varTypeName.(string)
									if varTypeNameStr == string(TypeStringString) {
										currentVarTypeName = TypeStringString
									} else if varTypeNameStr == string(TypeIntString) {
										currentVarTypeName = TypeIntString
									} else if varTypeNameStr == string(TypeFloatString) {
										currentVarTypeName = TypeFloatString
									} else {
										return fmt.Errorf(
											"this type of variable does not exist. Available types: %s, %s, %s. Got: %s. Locale: %s, Id: %s, Text: %s, Param: %s, Variable: %s",
											TypeStringString,
											TypeIntString,
											TypeFloatString,
											varTypeNameStr,
											locale.LanguageCode,
											id,
											data_,
											translateParamsSymbol+id,
											variable.VariableName,
										)
									}
									variable.VariableType = currentVarTypeName

									if varDefaultValue, ok := varParams["defaultValue"]; ok {
										if _, okDefineVariableType := p.defineVariableType(varDefaultValue, variable.VariableType); okDefineVariableType {
											variable.DefaultValue = varDefaultValue
										} else {
											return fmt.Errorf(
												"the default value type does not match the specified variable type! Expected Type: %s. Got: %T. Locale: %s, Id: %s, Text: %s, Param: %s, Variable: %s",
												variable.VariableType,
												varDefaultValue,
												locale.LanguageCode,
												id,
												data_,
												translateParamsSymbol+id,
												variable.VariableName,
											)
										}
									} else {
										variable.DefaultValue = nil
									}
								} else {
									return fmt.Errorf(
										"no data type has been selected for the specified variable! Locale: %s, Id: %s, Text: %s, Param: %s, Variable: %s",
										locale.LanguageCode,
										id,
										data_,
										translateParamsSymbol+id,
										variable.VariableName,
									)
								}
								translateParams.Variables = append(translateParams.Variables, variable)
							}
						} else {
							translateParams.Variables = nil
						}
						translateData.Params = translateParams
					} else {
						translateData.Params = nil
					}
					if _, okData := locale.Data[id]; !okData {
						locale.Data[id] = translateData
					} else {
						return fmt.Errorf("such an identifier already exists. Locale: %s, Id: %s", locale.LanguageCode, id)
					}
				} else {
					return fmt.Errorf("the translation text must be of type string! Locale: %s, Id: %s, Text: %v, Type: %T", locale.LanguageCode, id, data_, data_)
				}
			}
		}
		p.ArrayOfLocalizations = append(p.ArrayOfLocalizations, locale)
	}

	return nil
}

// equalization ensures all language translations have the same set of keys as the template language.
func (p *Parser) equalization(templateLanguageCode string) {
	var defaultTranslate map[string]translateData

	for _, val := range p.ArrayOfLocalizations {
		if val.LanguageCode == templateLanguageCode {
			defaultTranslate = val.Data
			break
		}
	}

	for key, value := range defaultTranslate {
		for i := range p.ArrayOfLocalizations {
			if _, ok := p.ArrayOfLocalizations[i].Data[key]; !ok {
				p.ArrayOfLocalizations[i].Data[key] = value
			}
		}
	}
}

// NewParser creates a new instance of the Parser, initializes it with localization data,
// and performs equalization of translations based on a template language code.
// It reads localization JSON files from the specified directory `pathToTranslates`, uses
// the `templateLanguageCode` for validation and initialization purposes,
// then parses these data and equalizes translations for all supported languages.
//
// Example usage:
//
//	p, err := NewParser("/path/to/translations", "en")
//	if err != nil {
//		log.Fatalf("Error initializing Parser: %v", err)
//	}
//	// Now you have access to the array of localizations through `p.ArrayOfLocalizations`
//
// Example JSON data structure it processes:
//
//	{
//		"l10n_en.json": {
//			"hello_world": "Hello, World!",
//			"#hello_world": {
//				"description": "Greeting message",
//				"example": "Hello, World!"
//			}
//		},
//		"l10n_fr.json": {
//			"hello_world": "Bonjour le monde!",
//			"#hello_world": {
//				"description": "Message de salutation",
//				"example": "Bonjour le monde!"
//			}
//		}
//	}
//
// In this example, `l10n_en.json` and `l10n_fr.json` contain translations for different languages,
// and `#hello_world` represents translation parameters with a description and usage example.
func NewParser(pathToTranslates, templateLanguageCode string) (*Parser, error) {
	pars := &Parser{}

	data, errData := pars.readLocales(pathToTranslates)
	if errData != nil {
		return nil, errData
	}

	if !pars.templateLanguageCodeInLocales(templateLanguageCode, data) {
		return nil, fmt.Errorf("template: %s, not contained in localizations", templateLanguageCode)
	}

	if err := pars.parse(data); err != nil {
		return nil, err
	} else {
		pars.equalization(templateLanguageCode)
		return pars, nil
	}
}
