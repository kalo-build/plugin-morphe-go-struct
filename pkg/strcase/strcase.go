package strcase

import "github.com/gobeam/stringy"

func ToCamelCase(input string) string {
	inputStr := stringy.New(input)
	return inputStr.CamelCase()
}
