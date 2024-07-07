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
	cfg.MorpheModelsConfig

	RegistryHooks r.LoadMorpheRegistryHooks

	ModelWriter write.GoStructWriter
	ModelHooks  hook.CompileMorpheModel

	WriteStructHooks hook.WriteGoStruct
}
