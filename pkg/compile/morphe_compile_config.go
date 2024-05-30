package compile

import (
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/hook"
)

type MorpheCompileConfig struct {
	cfg.MorpheLoadRegistryConfig
	cfg.MorpheModelsConfig

	RegistryHooks hook.LoadMorpheRegistry

	ModelWriter StructWriter
	ModelHooks  hook.CompileMorpheModels

	EntityWriter StructWriter
	EntityHooks  hook.CompileMorpheEntities
}
