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
	result.typesMap["date"] = "time.Time"

	return result
}

func (g *goLanguageSerializer) getType() languageType {
	return LanguageTypeGo
}

func (g *goLanguageSerializer) getTypeName() string {
	return "go"
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
		"//\tGenerated by ModelsGenerator\n//\t" +
		time.Now().Format(time.RFC3339) +
		"\n// **********************************\n\n"

	if serializerInfo.packageName == "" {
		return generatedMark
	}

	return fmt.Sprintf("package %s\n\n", serializerInfo.packageName) + generatedMark
}

func (g *goLanguageSerializer) serializeClass(class *class, serializerInfo *serializerInfo) (*generatedCode, error) {
	serializedCode := g.serializeDeclaration(serializerInfo)
	fileName := fmt.Sprintf("%s.go", class.name)

	serializedCode += fmt.Sprintf("type %s struct {\n", toFirstCharUpper(class.name))

	for _, member := range class.dataMembers {
		if isList, listType := isList(member.memberType); isList {
			// If the list type isn't primitive, we put it as a pointer
			pointerMark := ""
			if _, shouldNotBePointer := g.typesMap[listType]; !shouldNotBePointer {
				pointerMark = "*"
			}

			primitiveType, isPrimitive := g.typesMap[listType]
			if isPrimitive {
				listType = primitiveType
			}

			serializedCode += fmt.Sprintf("\t%s []%s%s `json:\"%s\"`\n",
				toFirstCharUpper(member.name), pointerMark, listType, toCamelCase(member.name))
			continue
		}

		if isMap, mapKeyType, mapValueType := isMap(member.memberType); isMap {
			// If the map value type isn't primitive, we put it as a pointer
			// We assume the key is a primitive
			pointerMark := ""
			if _, shouldNotBePointer := g.typesMap[mapValueType]; !shouldNotBePointer {
				pointerMark = "*"
			}

			primitiveType, isPrimitive := g.typesMap[mapKeyType]
			if isPrimitive {
				mapKeyType = primitiveType
			}

			primitiveType, isPrimitive = g.typesMap[mapValueType]
			if isPrimitive {
				mapValueType = primitiveType
			}

			serializedCode += fmt.Sprintf("\t%s map[%s]%s%s `json:\"%s\"`\n", toFirstCharUpper(member.name),
				mapKeyType, pointerMark, mapValueType, toCamelCase(member.name))
			continue
		}

		if primitiveType, ok := g.typesMap[member.memberType]; ok {
			serializedCode += fmt.Sprintf("\t%s %s `json:\"%s\"`\n",
				toFirstCharUpper(member.name), primitiveType, toCamelCase(member.name))
			continue
		}

		// It's not a language type, so we create it as a pointer
		serializedCode += fmt.Sprintf("\t%s *%s `json:\"%s\"`\n",
			toFirstCharUpper(member.name), member.memberType, toCamelCase(member.name))
	}

	serializedCode += "}"

	return newGeneratedCode(fileName, serializedCode), nil
}

func (g *goLanguageSerializer) serializeEnum(enum *enum, serializerInfo *serializerInfo) (*generatedCode, error) {
	serializedCode := g.serializeDeclaration(serializerInfo)
	fileName := fmt.Sprintf("%s.go", enum.name)

	serializedCode += fmt.Sprintf("type %s int\n\n"+
		"const (\n", enum.name)

	for _, value := range enum.enumValues {
		serializedCode += fmt.Sprintf("\t%s%s = %s(%v)\n", enum.name,
			toFirstCharUpper(value.name), enum.name, value.value)
	}

	serializedCode += ")"

	return newGeneratedCode(fileName, serializedCode), nil
}
