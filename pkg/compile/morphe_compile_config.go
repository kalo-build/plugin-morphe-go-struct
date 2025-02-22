package compile

import (
	r "github.com/kaloseia/morphe-go/pkg/registry"
	rcfg "github.com/kaloseia/morphe-go/pkg/registry/cfg"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/hook"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/write"
)

type MorpheCompileConfig struct {
	rcfg.MorpheLoadRegistryConfig
	cfg.MorpheConfig

	RegistryHooks r.LoadMorpheRegistryHooks

	ModelWriter write.GoStructWriter
	ModelHooks  hook.CompileMorpheModel

	StructureWriter write.GoStructWriter
	StructureHooks  hook.CompileMorpheStructure

	EntityWriter write.GoStructWriter
	EntityHooks  hook.CompileMorpheEntity

	EnumWriter write.GoEnumWriter
	EnumHooks  hook.CompileMorpheEnum

	WriteStructHooks hook.WriteGoStruct
	WriteGoEnumHooks hook.WriteGoEnum
}
