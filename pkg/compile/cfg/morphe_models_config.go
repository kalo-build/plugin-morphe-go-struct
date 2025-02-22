package cfg

import (
	"fmt"

	"github.com/kaloseia/go/pkg/godef"
)

type MorpheModelsConfig struct {
	Package godef.Package

	// ReceiverName is the standard receiver name for the compiled model receiver methods, ie "m" in "func (m *MyModel) MyMethod(){}"
	ReceiverName string
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
	return nil
}

func (config MorpheModelsConfig) GetPackage() godef.Package {
	return config.Package
}

func (config MorpheModelsConfig) GetReceiverName() string {
	return config.ReceiverName
}
