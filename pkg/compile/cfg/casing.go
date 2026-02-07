package cfg

import "github.com/kalo-build/go-util/strcase"

// Casing represents a field naming convention
type Casing string

const (
	// CasingNone means no JSON tags will be added
	CasingNone Casing = ""
	// CasingCamel converts to camelCase (e.g., "firstName")
	CasingCamel Casing = "camel"
	// CasingSnake converts to snake_case (e.g., "first_name")
	CasingSnake Casing = "snake"
	// CasingPascal keeps PascalCase (e.g., "FirstName")
	CasingPascal Casing = "pascal"
)

// IsValid returns true if the casing is a valid option
func (c Casing) IsValid() bool {
	switch c {
	case CasingNone, CasingCamel, CasingSnake, CasingPascal:
		return true
	default:
		return false
	}
}

// Apply converts a field name to the specified casing
func (c Casing) Apply(fieldName string) string {
	switch c {
	case CasingCamel:
		return strcase.ToCamelCase(fieldName)
	case CasingSnake:
		return strcase.ToSnakeCaseLower(fieldName)
	case CasingPascal:
		return strcase.ToPascalCase(fieldName)
	default:
		return fieldName
	}
}
