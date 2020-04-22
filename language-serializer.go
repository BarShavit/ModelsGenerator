package main

type languageType int

const (
	LanguageTypeCSharp     = languageType(0)
	LanguageTypeGo         = languageType(1)
	LanguageTypeKotlin     = languageType(2)
	LanguageTypeTypescript = languageType(3)
)

type serializerInfo struct {
	packageName string
}

type generatedCode struct {
	fileName string
	code     string
}

func newGeneratedCode(fileName string, code string) *generatedCode {
	return &generatedCode{fileName: fileName, code: code}
}

type languageSerializer interface {
	getType() languageType
	getTypeName() string
	generateCode(objects []middleware, serializerInfo *serializerInfo) ([]*generatedCode, error)
}
