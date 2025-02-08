package hook

import (
	"github.com/kaloseia/go/pkg/godef"
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/cfg"
)

type CompileMorpheStructure struct {
	OnCompileMorpheStructureStart   OnCompileMorpheStructureStartHook
	OnCompileMorpheStructureSuccess OnCompileMorpheStructureSuccessHook
	OnCompileMorpheStructureFailure OnCompileMorpheStructureFailureHook
}

type OnCompileMorpheStructureStartHook = func(config cfg.MorpheConfig, structure yaml.Structure) (cfg.MorpheConfig, yaml.Structure, error)
type OnCompileMorpheStructureSuccessHook = func(structureStruct *godef.Struct) (*godef.Struct, error)
type OnCompileMorpheStructureFailureHook = func(config cfg.MorpheConfig, structure yaml.Structure, compileFailure error) error
