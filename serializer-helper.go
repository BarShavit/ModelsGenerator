package main

import (
	"fmt"
	"strings"
)

/**
Check if the data member is list.
We expect list to write in the gen file like "memberName list<Type>".
If it's a list, we extract the list type.
Return true if map, the first string the the key type and the
second string is the value type
*/
func isList(memberType string) (bool, string) {
	if strings.Contains(memberType, "list") {
		memberType = strings.Replace(memberType, "<", " ", -1)
		memberType = strings.Replace(memberType, ">", " ", -1)

		splittedList := strings.Split(memberType, " ")

		if len(splittedList) < 2 {
			return false, ""
		}

		return true, splittedList[1]
	}

	return false, ""
}

/**
Check if the data member is map.
We expect map to write in the gen file like "memberName map<Type,Type>".
If it's a map, we extract the map types
*/
func isMap(memberType string) (bool, string, string) {
	if strings.Contains(memberType, "map") {
		memberType = strings.Replace(memberType, "<", " ", -1)
		memberType = strings.Replace(memberType, ">", " ", -1)

		splittedList := strings.Split(memberType, " ")

		if len(splittedList) < 2 {
			return false, "", ""
		}

		typesSplit := strings.Split(splittedList[1], ",")

		if len(typesSplit) < 2 {
			return false, "", ""
		}

		return true, typesSplit[0], typesSplit[1]
	}

	return false, "", ""
}

func toCamelCase(value string) string {
	if len(value) < 2 {
		return strings.ToLower(value)
	}

	return fmt.Sprintf("%s%s", strings.ToLower(string(value[0])), value[1:])
}

func toFirstCharUpper(value string) string {
	if len(value) < 2 {
		return strings.ToUpper(value)
	}

	return fmt.Sprintf("%s%s", strings.ToUpper(string(value[0])), value[1:])
}

func appendUnique(strings []string, str string) []string {
	for _, s := range strings {
		if s == str {
			return strings
		}
	}

	return append(strings, str)
}

func memberUnique(members []*dataMember, member *dataMember) bool {
	for _, m := range members {
		if m.name == member.name {
			return false
		}
	}

	return true
}

func enumValueUnique(values []*enumValue, value *enumValue) bool {
	for _, v := range values {
		if v.name == value.name {
			return false
		}
	}

	return true
}
