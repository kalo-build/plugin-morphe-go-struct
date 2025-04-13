package hook

import (
	"github.com/kalo-build/go/pkg/godef"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/compile/cfg"
)

type CompileMorpheStructure struct {
	OnCompileMorpheStructureStart   OnCompileMorpheStructureStartHook
	OnCompileMorpheStructureSuccess OnCompileMorpheStructureSuccessHook
	OnCompileMorpheStructureFailure OnCompileMorpheStructureFailureHook
}

type OnCompileMorpheStructureStartHook = func(config cfg.MorpheConfig, structure yaml.Structure) (cfg.MorpheConfig, yaml.Structure, error)
type OnCompileMorpheStructureSuccessHook = func(structureStruct *godef.Struct) (*godef.Struct, error)
type OnCompileMorpheStructureFailureHook = func(config cfg.MorpheConfig, structure yaml.Structure, compileFailure error) error
