package cfg

import (
	"fmt"

	"github.com/kalo-build/go/pkg/godef"
)

type MorpheModelsConfig struct {
	Package godef.Package

	// ReceiverName is the standard receiver name for the compiled model receiver methods, ie "m" in "func (m *MyModel) MyMethod(){}"
	ReceiverName string

	// FieldCasing specifies the casing for serialization (JSON struct tags). Empty means no tags.
	// Valid values: "camel", "snake", "pascal", or "" (none)
	FieldCasing Casing
}

func (config MorpheModelsConfig) Validate() error {
	if config.Package.Path == "" {
		return fmt.Errorf("models %w", ErrNoPackagePath)
	}
	if config.Package.Name == "" {
		return fmt.Errorf("models %w", ErrNoPackageName)
	}
	if config.ReceiverName == "" {
		return fmt.Errorf("models %w", ErrNoReceiverName)
	}
	if !config.FieldCasing.IsValid() {
		return fmt.Errorf("models: invalid fieldCasing value %q, must be one of: camel, snake, pascal, or empty", config.FieldCasing)
	}
	return nil
}

func (config MorpheModelsConfig) GetPackage() godef.Package {
	return config.Package
}

func (config MorpheModelsConfig) GetReceiverName() string {
	return config.ReceiverName
}
