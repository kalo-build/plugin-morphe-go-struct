package compile

import (
	"errors"
	"fmt"

	"github.com/kalo-build/morphe-go/pkg/yaml"
)

var ErrNoEnumType = errors.New("no enum type provided")
var ErrNoEnum = errors.New("no enum provided")

func ErrUnsupportedEnumType(enumType yaml.EnumType) error {
	return fmt.Errorf("unsupported enum type: %s", enumType)
}

func ErrEnumEntryNotFound(entryName string) error {
	return fmt.Errorf("enum entry not found: %s", entryName)
}
