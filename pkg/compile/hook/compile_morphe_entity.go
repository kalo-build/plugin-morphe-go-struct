package hook

import (
	"github.com/kaloseia/go/pkg/godef"
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/cfg"
)

type CompileMorpheEntity struct {
	OnCompileMorpheEntityStart   OnCompileMorpheEntityStartHook
	OnCompileMorpheEntitySuccess OnCompileMorpheEntitySuccessHook
	OnCompileMorpheEntityFailure OnCompileMorpheEntityFailureHook
}

type OnCompileMorpheEntityStartHook = func(config cfg.MorpheConfig, entity yaml.Entity) (cfg.MorpheConfig, yaml.Entity, error)
type OnCompileMorpheEntitySuccessHook = func(entityStructs []*godef.Struct) ([]*godef.Struct, error)
type OnCompileMorpheEntityFailureHook = func(config cfg.MorpheConfig, entity yaml.Entity, compileFailure error) error
