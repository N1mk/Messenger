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

	return str
}

func SpaceSplittedToUpperCamelCase(str string) string {
	splittedStr := strings.Fields(str)

	str = ""

	for _, s := range splittedStr {
		if len(s) > 1 && s[1] > 'A' && s[1] < 'Z' {
			str += strings.ToUpper(s)
			continue
		}

		str += strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
	}

	return str
}

func SpaceSplittedToLowerCamelCase(str string) string {
	firstWord := strings.Fields(str)[0]

	str = SpaceSplittedToUpperCamelCase(str)

	return strings.ToLower(firstWord) + str[len(firstWord):]
}
