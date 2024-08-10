package main

import (
	"flag"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/a1ex3/l10n/internal/config"
	generatorcpp "github.com/a1ex3/l10n/internal/generator/cpp"
	generatorjava "github.com/a1ex3/l10n/internal/generator/java"
	generatorjs "github.com/a1ex3/l10n/internal/generator/javascript"
	generatorkotlin "github.com/a1ex3/l10n/internal/generator/kotlin"
	generatorpython "github.com/a1ex3/l10n/internal/generator/python"
	generatorts "github.com/a1ex3/l10n/internal/generator/typescript"
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
	if conf.ProgrammingLanguage == availableProgrammingLanguages[0] { // Python
		gen, errGen := generatorpython.NewGeneratorPython(conf.ClassName, conf.Template, prs).Get()
		if errGen != nil {
			log.Fatalln(errGen)
		}
		if errWriteToFile := writeToFile(conf.OutputLocalizationFile, gen); errWriteToFile != nil {
			log.Fatalln(errWriteToFile)
		}
	} else if regexp.MustCompile(`^([A-Za-z0-9]+\.)*[A-Za-z0-9]+$`).MatchString(conf.ProgrammingLanguage) { // Used only if you need to specify a package when generating a file.
		packageName := ""
		programmingLang := ""

		if index := strings.Index(conf.ProgrammingLanguage, "."); index > 0 {
			programmingLang = conf.ProgrammingLanguage[:index]
			packageName = conf.ProgrammingLanguage[index+1:]
		} else {
			programmingLang = conf.ProgrammingLanguage
			packageName = "l10n"
		}

		if programmingLang == availableProgrammingLanguages[1] { // Java
			gen, errGen := generatorjava.NewGeneratorJava(conf.ClassName, packageName, conf.Template, prs).Get()
			if errGen != nil {
				log.Fatalln(errGen)
			}
			if errWriteToFile := writeToFile(conf.OutputLocalizationFile, gen); errWriteToFile != nil {
				log.Fatalln(errWriteToFile)
			}
		} else if programmingLang == availableProgrammingLanguages[2] {
			gen, errGen := generatorcpp.NewGeneratorCpp(conf.ClassName, packageName, conf.Template, prs).Get()
			if errGen != nil {
				log.Fatalln(errGen)
			}
			if errWriteToFile := writeToFile(conf.OutputLocalizationFile, gen); errWriteToFile != nil {
				log.Fatalln(errWriteToFile)
			}
		} else if programmingLang == availableProgrammingLanguages[3] {
			gen, errGen := generatorkotlin.NewGeneratorKotlin(conf.ClassName, packageName, conf.Template, prs).Get()
			if errGen != nil {
				log.Fatalln(errGen)
			}
			if errWriteToFile := writeToFile(conf.OutputLocalizationFile, gen); errWriteToFile != nil {
				log.Fatalln(errWriteToFile)
			}
		} else if programmingLang == availableProgrammingLanguages[4] {
			gen, errGen := generatorts.NewGeneratorTS(conf.ClassName, packageName, conf.Template, prs).Get()
			if errGen != nil {
				log.Fatalln(errGen)
			}
			if errWriteToFile := writeToFile(conf.OutputLocalizationFile, gen); errWriteToFile != nil {
				log.Fatalln(errWriteToFile)
			}
		} else if programmingLang == availableProgrammingLanguages[5] {
			gen, errGen := generatorjs.NewGeneratorJS(conf.ClassName, packageName, conf.Template, prs).Get()
			if errGen != nil {
				log.Fatalln(errGen)
			}
			if errWriteToFile := writeToFile(conf.OutputLocalizationFile, gen); errWriteToFile != nil {
				log.Fatalln(errWriteToFile)
			}
		} else {
			log.Fatalf("This programming language is not supported!")
		}
	}
}
