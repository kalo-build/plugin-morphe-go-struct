package hook

import (
	"github.com/kaloseia/go/pkg/godef"
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/cfg"
)

type CompileMorpheEnum struct {
	OnCompileMorpheEnumStart   func(cfg.MorpheEnumsConfig, yaml.Enum) (cfg.MorpheEnumsConfig, yaml.Enum, error)
	OnCompileMorpheEnumSuccess func(*godef.Enum) (*godef.Enum, error)
	OnCompileMorpheEnumFailure func(cfg.MorpheEnumsConfig, yaml.Enum, error) error
}
