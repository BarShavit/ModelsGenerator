package ModelsGenerator

import "strings"

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

		return true, typesSplit[0], typesSplit[1]
	}

	return false, "", ""
}