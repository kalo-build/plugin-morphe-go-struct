package cfg

import (
	"fmt"

	"github.com/kaloseia/go/pkg/godef"
)

type MorpheEnumsConfig struct {
	Package godef.Package
}

func (config MorpheEnumsConfig) Validate() error {
	if config.Package.Path == "" {
		return fmt.Errorf("enums %w", ErrNoPackagePath)
	}
	if config.Package.Name == "" {
		return fmt.Errorf("enums %w", ErrNoPackageName)
	}
	return nil
}
