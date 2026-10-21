package stringoperations

import "strings"

func SpaceSplittedToUnderscoreSplitted(str string) string {
	splittedStr := strings.Fields(str)

	str = ""

	for i, s := range splittedStr {
		str += strings.ToLower(s)
		if i != len(splittedStr)-1 {
			str += "_"
		}
	}

	str += ".go"

	return str
}

func SpaceSplittedToUpperCamelCase(str string) string {
	splittedStr := strings.Fields(str)

	for _, s := range splittedStr {
		str += strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
	}

	return str
}

func SpaceSplittedToLowerCamelCase(str string) string {
	str = SpaceSplittedToUpperCamelCase(str)

	return strings.ToLower(string(str[0])) + str[1:]
}
