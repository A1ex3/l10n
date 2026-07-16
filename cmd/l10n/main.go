package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/a1ex3/l10n/internal/config"
	"github.com/a1ex3/l10n/internal/generator"
	"github.com/a1ex3/l10n/internal/parser"
)

var (
	availableProgrammingLanguages = []string{
		"Python",
		"Java",
		"Cpp",
		"Kotlin",
		"TypeScript",
		"JavaScript",
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

	conf := cfg.GetConfig()
	packageName := "l10n"
	programmingLang := strings.ToLower(conf.ProgrammingLanguage)
	if index := strings.Index(conf.ProgrammingLanguage, "."); index > 0 {
		packageName = conf.ProgrammingLanguage[index+1:]
		programmingLang = strings.ToLower(conf.ProgrammingLanguage[:index])
	}

	gen, errGen := generator.NewGenerator(programmingLang, conf.ClassName, conf.Template, packageName, prs)
	if errGen != nil {
		log.Fatalln(errGen)
	}

	code, errCode := gen.Get()
	if errCode != nil {
		log.Fatalln(errCode)
	}
	if errWriteToFile := writeToFile(conf.OutputLocalizationFile, code); errWriteToFile != nil {
		log.Fatalln(errWriteToFile)
	}
}
