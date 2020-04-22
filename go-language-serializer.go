package main

import (
	"errors"
	"fmt"
	"time"
)

type goLanguageSerializer struct {
	typesMap map[string]string
}

func newGoLanguageSerializer() *goLanguageSerializer {
	result := &goLanguageSerializer{typesMap: make(map[string]string, 0)}

	result.typesMap["bool"] = "bool"
	result.typesMap["int"] = "int"
	result.typesMap["string"] = "string"
	result.typesMap["double"] = "float64"
	result.typesMap["float"] = "float32"
	result.typesMap["char"] = "byte"
	result.typesMap["byte"] = "byte"

	return result
}

func (g *goLanguageSerializer) getType() languageType {
	return LanguageTypeGo
}

func (g *goLanguageSerializer) generateCode(objects []middleware, serializerInfo *serializerInfo) ([]*generatedCode, error) {
	result := make([]*generatedCode, 0)

	for _, object := range objects {
		serialized, err := g.serializeMiddleware(object, serializerInfo)
		if err != nil {
			return nil, err
		}

		result = append(result, serialized)
	}

	return result, nil
}

func (g *goLanguageSerializer) serializeMiddleware(middleware middleware, serializerInfo *serializerInfo) (*generatedCode, error) {
	if class, ok := middleware.(*class); ok {
		return g.serializeClass(class, serializerInfo)
	}

	if enum, ok := middleware.(*enum); ok {
		return g.serializeEnum(enum, serializerInfo)
	}

	return nil, errors.New("tried to serialize unknown middleware type")
}

func (g *goLanguageSerializer) serializeDeclaration(serializerInfo *serializerInfo) string {
	generatedMark := "// **********************************\n" +
		"Generated by ModelsGenerator\n" +
		time.Now().String() +
		"// **********************************\n\n"

	if serializerInfo.packageName == "" {
		return generatedMark
	}

	return fmt.Sprintf("package %s\n\n", serializerInfo.packageName) + generatedMark
}

func (g *goLanguageSerializer) serializeClass(class *class, serializerInfo *serializerInfo) (*generatedCode, error) {
	serializedCode := g.serializeDeclaration(serializerInfo)
	fileName := fmt.Sprintf("%s.go", class.name)

	serializedCode += fmt.Sprintf("type %s struct {\n", class.name)

	for _, member := range class.dataMembers {
		if isList, listType := isList(member.memberType); isList {
			serializedCode += fmt.Sprintf("\t%s []*%s", member.name, listType)
			continue
		}

		if isMap, mapKeyType, mapValueType := isMap(member.memberType); isMap {
			serializedCode += fmt.Sprintf("\t%s map[%s]*%s", member.name, mapKeyType, mapValueType)
			continue
		}

		if primitiveType, ok := g.typesMap[member.memberType]; ok {
			serializedCode += fmt.Sprintf("\t%s %s", member.name, primitiveType)
			continue
		}

		// It's not a language type, so we create it as a pointer
		serializedCode += fmt.Sprintf("\t%s *%s", member.name, member.memberType)
	}

	return newGeneratedCode(fileName, serializedCode), nil
}

func (g *goLanguageSerializer) serializeEnum(enum *enum, serializerInfo *serializerInfo) (*generatedCode, error) {
	serializedCode := g.serializeDeclaration(serializerInfo)
	fileName := fmt.Sprintf("%s.go", enum.name)

	serializedCode += fmt.Sprintf("type %s int\n\n"+
		"const (\n", enum.name)

	for _, value := range enum.enumValues {
		serializedCode += fmt.Sprintf("\t%s%s = %s(%v)\n", enum.name, value.name, enum.name, value.value)
	}

	serializedCode += ")"

	return newGeneratedCode(fileName, serializedCode), nil
}
