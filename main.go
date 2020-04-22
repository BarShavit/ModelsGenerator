package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var serializers = make(map[languageType]languageSerializer)
var languageMap = make(map[string]languageType)

type languageParameter struct {
	languageType languageType
	packageName  string
}

func main() {
	serializers[LanguageTypeGo] = newGoLanguageSerializer()

	languageMap["go"] = LanguageTypeGo

	if err := handleCommand(); err != nil {
		panic(err)
	}
}

func handleCommand() error {
	if len(os.Args) > 0 && os.Args[0] == "help" {
		println("Hello! This tool helps you to convert a \"gen\" file into different languages models.\n" +
			"The format should be \"gen-file-path.gen language:package\". You can avoid from adding package name.\n" +
			"The supported languages are Go, Kotlin, C# and Typescript.\n" +
			"For more info and documentation for gen files, visit https://github.com/BarShavit/ModelsGenerator.\n" +
			"Thank you!")
	}
	if len(os.Args) < 2 {
		return errors.New("didn't receive enough arguments. expected \"*filepath* language:package")
	}

	return handleGenerate()
}

func handleGenerate() error {
	filePath := os.Args[0]
	languages := make([]*languageParameter, 0)

	// Get all the languages we should generate from arguments
	for _, arg := range os.Args[1:] {
		if param, err := parseToLanguageParameter(arg); err != nil {
			return err
		} else {
			languages = append(languages, param)
		}
	}

	meddlers, err := parseFile(filePath)
	if err != nil {
		return err
	}

	for _, lang := range languages {
		serializers[lang.languageType].generateCode(meddlers,
			&serializerInfo{packageName: lang.packageName})
	}
}

func parseToLanguageParameter(parameter string) (*languageParameter, error) {
	splittedParameter := strings.Split(parameter, ":")
	packageName := ""

	if len(splittedParameter) > 1 {
		packageName = splittedParameter[1]
	}

	languageType, ok := languageMap[splittedParameter[0]]
	if !ok {
		return nil, errors.New(
			fmt.Sprintf("received unsupported language %s", splittedParameter[0]))
	}

	return &languageParameter{
		languageType: languageType,
		packageName:  packageName,
	}, nil
}
