package cfg

import (
	"fmt"

	"github.com/kaloseia/go/pkg/godef"
)

type MorpheEntitiesConfig struct {
	Package godef.Package

	// ReceiverName is the standard receiver name for the compiled entity receiver methods, ie "e" in "func (e *MyEntity) MyMethod(){}"
	ReceiverName string
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
	return nil
}

func (config MorpheEntitiesConfig) GetPackage() godef.Package {
	return config.Package
}

func (config MorpheEntitiesConfig) GetReceiverName() string {
	return config.ReceiverName
}
