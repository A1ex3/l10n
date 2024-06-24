package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type entityConfig struct {
	Dir                    string `yaml:"dir"`
	Template               string `yaml:"template"`
	OutputLocalizationFile string `yaml:"output_localization_file"`
	ProgrammingLanguage    string `yaml:"programming_language"`
	ClassName              string `yaml:"class_name"`
}

type Config struct {
	availableProgrammingLanguages []string
	entityConfig                  *entityConfig
}

func readFile(filepath string) ([]byte, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (c *Config) argsValidator(
	dir string,
	template string,
	outputLocalizationFile string,
	programmingLanguage string,
	className string,
) error {
	if dir == "" {
		return errors.New("dir cannot be empty")
	}

	if template == "" {
		return errors.New("template cannot be empty")
	}

	if outputLocalizationFile == "" {
		return errors.New("output_localization_file cannot be empty")
	}

	if programmingLanguage == "" {
		return errors.New("programming_language cannot be empty")
	}
	validProgrammingLanguage := false
	for _, lang := range c.availableProgrammingLanguages {
		if lang == programmingLanguage {
			validProgrammingLanguage = true
			break
		}
	}
	if !validProgrammingLanguage {
		return errors.New("invalid programming_language")
	}

	if className == "" {
		return errors.New("class_name cannot be empty")
	}

	return nil
}

func (c *Config) FromArgs(
	dir string,
	template string,
	outputLocalizationFile string,
	programmingLanguage string,
	className string,
) error {
	if err := c.argsValidator(dir, template, outputLocalizationFile, programmingLanguage, className); err != nil {
		return err
	}

	c.entityConfig = &entityConfig{
		Dir:                    dir,
		Template:               template,
		OutputLocalizationFile: outputLocalizationFile,
		ProgrammingLanguage:    programmingLanguage,
		ClassName:              className,
	}

	return nil
}

func (c *Config) FromFileYaml(filepath string) error {
	dataFromFile, errDataFromFile := readFile(filepath)
	if errDataFromFile != nil {
		return errDataFromFile
	}

	cfg := &entityConfig{}
	err := yaml.Unmarshal([]byte(dataFromFile), cfg)
	if err != nil {
		return err
	}

	c.entityConfig = cfg

	if err := c.argsValidator(c.entityConfig.Dir, c.entityConfig.Template, c.entityConfig.OutputLocalizationFile, c.entityConfig.ProgrammingLanguage, c.entityConfig.ClassName); err != nil {
		return err
	}

	return nil
}

func (c *Config) GetConfig() *entityConfig {
	return c.entityConfig
}

func NewConfig(availableProgrammingLanguages []string) *Config {
	return &Config{
		availableProgrammingLanguages: availableProgrammingLanguages,
		entityConfig:                  nil,
	}
}
