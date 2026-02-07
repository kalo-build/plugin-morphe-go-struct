package cfg

import (
	"fmt"

	"github.com/kalo-build/go/pkg/godef"
)

type MorpheEntitiesConfig struct {
	Package godef.Package

	// ReceiverName is the standard receiver name for the compiled entity receiver methods, ie "e" in "func (e *MyEntity) MyMethod(){}"
	ReceiverName string

	// FieldCasing specifies the casing for serialization (JSON struct tags). Empty means no tags.
	// Valid values: "camel", "snake", "pascal", or "" (none)
	FieldCasing Casing
}

func (config MorpheEntitiesConfig) Validate() error {
	if config.Package.Path == "" {
		return fmt.Errorf("entities %w", ErrNoPackagePath)
	}
	if config.Package.Name == "" {
		return fmt.Errorf("entities %w", ErrNoPackageName)
	}
	if config.ReceiverName == "" {
		return fmt.Errorf("entities %w", ErrNoReceiverName)
	}
	if !config.FieldCasing.IsValid() {
		return fmt.Errorf("entities: invalid fieldCasing value %q, must be one of: camel, snake, pascal, or empty", config.FieldCasing)
	}
	return nil
}

func (config MorpheEntitiesConfig) GetPackage() godef.Package {
	return config.Package
}

func (config MorpheEntitiesConfig) GetReceiverName() string {
	return config.ReceiverName
}
