package hook

import (
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/godef"
)

type CompileMorpheModel struct {
	OnCompileMorpheModelStart   OnCompileMorpheModelStartHook
	OnCompileMorpheModelSuccess OnCompileMorpheModelSuccessHook
	OnCompileMorpheModelFailure OnCompileMorpheModelFailureHook
}

type OnCompileMorpheModelStartHook = func(config cfg.MorpheModelsConfig, model yaml.Model) (cfg.MorpheModelsConfig, yaml.Model, error)
type OnCompileMorpheModelSuccessHook = func(allModelStructs []*godef.Struct) ([]*godef.Struct, error)
type OnCompileMorpheModelFailureHook = func(config cfg.MorpheModelsConfig, model yaml.Model, compileFailure error) error
