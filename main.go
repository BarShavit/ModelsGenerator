package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

var serializers = make(map[languageType]languageSerializer)
var languageMap = make(map[string]languageType)

type languageParameter struct {
	languageType languageType
	packageName  string
}

func main() {
	serializers[LanguageTypeGo] = newGoLanguageSerializer()
	serializers[LanguageTypeKotlin] = newKotlinLanguageSerializer()

	languageMap["go"] = LanguageTypeGo
	languageMap["kotlin"] = LanguageTypeKotlin

	if err := handleCommand(); err != nil {
		panic(err)
	}
}

func handleCommand() error {
	if len(os.Args) > 1 && os.Args[1] == "help" {
		println("Hello! This tool helps you to convert a \"gen\" file into different languages models.\n" +
			"The format should be \"gen-file-path.gen language:package\". You can avoid from adding package name.\n" +
			"The supported languages are Go, Kotlin, C# and Typescript.\n" +
			"For more info and documentation for gen files, visit https://github.com/BarShavit/ModelsGenerator.\n" +
			"Thank you!")

		return nil
	}

	if len(os.Args) < 3 {
		return errors.New("didn't receive enough arguments. expected \"*filepath* language:package")
	}

	return handleGenerate()
}

func handleGenerate() error {
	filePath := os.Args[1]
	languages := make([]*languageParameter, 0)

	// Get all the languages we should generate from arguments
	for _, arg := range os.Args[2:] {
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

	generatedTime := time.Now().Format(time.RFC3339)
	generatedTime = strings.Replace(generatedTime, ":", "-", -1)

	for _, lang := range languages {
		generatedCode, err := serializers[lang.languageType].generateCode(meddlers,
			&serializerInfo{packageName: lang.packageName})

		if err != nil {
			return err
		}

		for _, code := range generatedCode {
			if err := saveGeneratedCode(code, serializers[lang.languageType].getTypeName(), generatedTime); err != nil {
				return err
			}
		}
	}

	return nil
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
