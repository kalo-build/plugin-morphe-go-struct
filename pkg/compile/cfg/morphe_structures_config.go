package cfg

import (
	"fmt"

	"github.com/kalo-build/go/pkg/godef"
)

type MorpheStructuresConfig struct {
	Package godef.Package

	// ReceiverName is the standard receiver name for the compiled model receiver methods, ie "m" in "func (m *MyModel) MyMethod(){}"
	ReceiverName string

	// FieldCasing specifies the casing for serialization (JSON struct tags). Empty means no tags.
	// Valid values: "camel", "snake", "pascal", or "" (none)
	FieldCasing Casing
}

func (config MorpheStructuresConfig) Validate() error {
	if config.Package.Path == "" {
		return fmt.Errorf("structures %w", ErrNoPackagePath)
	}
	if config.Package.Name == "" {
		return fmt.Errorf("structures %w", ErrNoPackageName)
	}
	if config.ReceiverName == "" {
		return fmt.Errorf("structures %w", ErrNoReceiverName)
	}
	if !config.FieldCasing.IsValid() {
		return fmt.Errorf("structures: invalid fieldCasing value %q, must be one of: camel, snake, pascal, or empty", config.FieldCasing)
	}
	return nil
}
