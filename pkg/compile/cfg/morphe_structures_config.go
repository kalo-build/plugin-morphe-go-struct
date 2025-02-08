package cfg

import (
	"fmt"

	"github.com/kaloseia/go/pkg/godef"
)

type MorpheStructuresConfig struct {
	Package godef.Package

	// ReceiverName is the standard receiver name for the compiled model receiver methods, ie "m" in "func (m *MyModel) MyMethod(){}"
	ReceiverName string
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
	return nil
}
