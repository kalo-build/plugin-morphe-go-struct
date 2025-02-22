package hook

import (
	"github.com/kaloseia/go/pkg/godef"
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/cfg"
)

type CompileMorpheEnum struct {
	OnCompileMorpheEnumStart   OnCompileMorpheEnumStartHook
	OnCompileMorpheEnumSuccess OnCompileMorpheEnumSuccessHook
	OnCompileMorpheEnumFailure OnCompileMorpheEnumFailureHook
}

type OnCompileMorpheEnumStartHook = func(config cfg.MorpheEnumsConfig, enum yaml.Enum) (cfg.MorpheEnumsConfig, yaml.Enum, error)
type OnCompileMorpheEnumSuccessHook = func(enum *godef.Enum) (*godef.Enum, error)
type OnCompileMorpheEnumFailureHook = func(config cfg.MorpheEnumsConfig, enum yaml.Enum, compileFailure error) error
