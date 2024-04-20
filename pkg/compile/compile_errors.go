package compile

import (
	"errors"
	"fmt"

	"github.com/kaloseia/morphe-go/pkg/yaml"
)

var ErrNoReceiverName = errors.New("method receiver name cannot be empty")
var ErrNoPackageName = errors.New("package name cannot be empty")
var ErrNoMorpheModelName = errors.New("morphe model has no name")
var ErrNoMorpheModelFields = errors.New("morphe model has no fields")
var ErrNoMorpheModelIdentifiers = errors.New("morphe model has no identifiers")

func ErrUnsupportedMorpheFieldType(unsupportedType yaml.ModelFieldType) error {
	return fmt.Errorf("unsupported morphe field type for go conversion: '%s'", unsupportedType)
}

func ErrMissingMorpheIdentifierField(modelName string, identifierName string, fieldName string) error {
	return fmt.Errorf("morphe model '%s' has no field '%s' referenced in identifiers ('%s')", modelName, identifierName, fieldName)
}
