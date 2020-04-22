package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func saveGeneratedCode(generatedCode *generatedCode, languageName string, generateDate string) error {
	folderPath := fmt.Sprintf("%s/%s", languageName, generateDate)
	filePath := fmt.Sprintf("%s/%s", folderPath, generatedCode.fileName)

	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, []byte(generatedCode.code), 0755)
}
