package hook

import (
	"github.com/kalo-build/go/pkg/godef"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/compile/cfg"
)

type CompileMorpheModel struct {
	OnCompileMorpheModelStart   OnCompileMorpheModelStartHook
	OnCompileMorpheModelSuccess OnCompileMorpheModelSuccessHook
	OnCompileMorpheModelFailure OnCompileMorpheModelFailureHook
}

type OnCompileMorpheModelStartHook = func(config cfg.MorpheConfig, model yaml.Model) (cfg.MorpheConfig, yaml.Model, error)
type OnCompileMorpheModelSuccessHook = func(allModelStructs []*godef.Struct) ([]*godef.Struct, error)
type OnCompileMorpheModelFailureHook = func(config cfg.MorpheConfig, model yaml.Model, compileFailure error) error
