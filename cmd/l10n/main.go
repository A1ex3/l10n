package main

import (
	"flag"
	"log"
	"os"

	"github.com/a1ex3/l10n/internal/config"
	generatorpython "github.com/a1ex3/l10n/internal/generator/python"
	"github.com/a1ex3/l10n/internal/parser"
)

var (
	availableProgrammingLanguages = []string{
		"Python",
	}
)

func writeToFile(path string, data string) error {
	bytes := []byte(data)
	err := os.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var (
		dirFlag                 = flag.String("DIR", "", "")
		templateFlag            = flag.String("TEMPLATE", "", "")
		outputLocalizationFlag  = flag.String("OUTPUT_LOCALIZATION_FILE", "", "")
		programmingLanguageFlag = flag.String("PROGRAMMING_LANGUAGE", "", "")
		classNameFlag           = flag.String("CLASS_NAME", "", "")
		fileFlag                = flag.String("FILE", "", "Use this parameter if you do not want to specify startup parameters, but still have a \"yaml\" file")
	)
	flag.Parse()

	cfg := config.NewConfig(availableProgrammingLanguages)
	if len(*fileFlag) > 0 {
		if err := cfg.FromFileYaml(*fileFlag); err != nil {
			log.Fatalln(err)
		}
	} else {

		if *dirFlag == "" || *templateFlag == "" || *outputLocalizationFlag == "" || *programmingLanguageFlag == "" || *classNameFlag == "" {
			log.Fatalln("enter the -h parameter, to get details about the startup parameters")
		}

		if err := cfg.FromArgs(*dirFlag, *templateFlag, *outputLocalizationFlag, *programmingLanguageFlag, *classNameFlag); err != nil {
			log.Fatalln(err)
		}
	}

	prs, errPrs := parser.NewParser(cfg.GetConfig().Dir, cfg.GetConfig().Template)
	if errPrs != nil {
		log.Fatalln(errPrs)
	}

	switch cfg := cfg.GetConfig(); cfg.ProgrammingLanguage {
	case availableProgrammingLanguages[0]: // Python
		gen, errGen := generatorpython.NewGeneratorPython(cfg.ClassName, cfg.Template, prs).Get()
		if errGen != nil {
			log.Fatalln(errGen)
		}
		if errWriteToFile := writeToFile(cfg.OutputLocalizationFile, gen); errWriteToFile != nil {
			log.Fatalln(errWriteToFile)
		}
	default:
		log.Fatalf("This programming language is not supported!")
	}
}
