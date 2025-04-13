package hook

import (
	"github.com/kalo-build/go/pkg/godef"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/compile/cfg"
)

type CompileMorpheEntity struct {
	OnCompileMorpheEntityStart   OnCompileMorpheEntityStartHook
	OnCompileMorpheEntitySuccess OnCompileMorpheEntitySuccessHook
	OnCompileMorpheEntityFailure OnCompileMorpheEntityFailureHook
}

type OnCompileMorpheEntityStartHook = func(config cfg.MorpheConfig, entity yaml.Entity) (cfg.MorpheConfig, yaml.Entity, error)
type OnCompileMorpheEntitySuccessHook = func(entityStructs []*godef.Struct) ([]*godef.Struct, error)
type OnCompileMorpheEntityFailureHook = func(config cfg.MorpheConfig, entity yaml.Entity, compileFailure error) error
