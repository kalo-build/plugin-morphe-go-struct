package compile

import (
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/hook"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/write"
)

type MorpheCompileConfig struct {
	cfg.MorpheLoadRegistryConfig
	cfg.MorpheModelsConfig

	RegistryHooks hook.LoadMorpheRegistry

	ModelWriter write.GoStructWriter
	ModelHooks  hook.CompileMorpheModel

	WriteStructHooks hook.WriteGoStruct
}
