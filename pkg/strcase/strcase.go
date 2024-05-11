package strcase

import (
	"regexp"

	"github.com/gobeam/stringy"
)

func ToCamelCase(input string) string {
	structuredInput := stringy.New(input)
	return structuredInput.CamelCase()
}

func ToSnakeCase(input string) string {
	structuredInput := stringy.New(input)
	snakeString := structuredInput.SnakeCase("?", "").Get()
	return separateAllCaseTransitions(snakeString, "_")
}

func ToSnakeCaseLower(input string) string {
	snakeString := ToSnakeCase(input)
	return ToLower(snakeString)
}

func ToLower(input string) string {
	structuredInput := stringy.New(input)
	return structuredInput.ToLower()
}

func separateAllCaseTransitions(input string, separator string) string {
	reTransitions := regexp.MustCompile(`([A-Z]|[0-9])([A-Z][^A-Z])`)
	matches := reTransitions.FindAllStringSubmatchIndex(input, -1)
	for _, match := range matches {
		if len(match) < 4 {
			continue
		}
		input = input[:match[2]+1] + separator + input[match[3]:]
	}
	return input
}
