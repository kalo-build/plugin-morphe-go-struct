package compile

import (
	"path"

	"github.com/kalo-build/go/pkg/godef"
	r "github.com/kalo-build/morphe-go/pkg/registry"
	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/compile/hook"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/compile/write"
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

func DefaultMorpheCompileConfig(
	yamlRegistryPath string,
	baseOutputDirPath string,
) MorpheCompileConfig {
	return MorpheCompileConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryEnumsDirPath:      path.Join(yamlRegistryPath, "enums"),
			RegistryModelsDirPath:     path.Join(yamlRegistryPath, "models"),
			RegistryStructuresDirPath: path.Join(yamlRegistryPath, "structures"),
			RegistryEntitiesDirPath:   path.Join(yamlRegistryPath, "entities"),
		},
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig: cfg.MorpheModelsConfig{
				Package: godef.Package{
					Name: "models",
				},
				ReceiverName: "m",
			},
			MorpheEnumsConfig: cfg.MorpheEnumsConfig{
				Package: godef.Package{
					Name: "enums",
				},
			},
			MorpheStructuresConfig: cfg.MorpheStructuresConfig{
				Package: godef.Package{
					Name: "structures",
				},
				ReceiverName: "s",
			},
			MorpheEntitiesConfig: cfg.MorpheEntitiesConfig{
				Package: godef.Package{
					Name: "entities",
				},
				ReceiverName: "e",
			},
		},

		RegistryHooks: r.LoadMorpheRegistryHooks{},

		EnumWriter: &MorpheEnumFileWriter{
			TargetDirPath: path.Join(baseOutputDirPath, "enums"),
		},
		EnumHooks: hook.CompileMorpheEnum{},

		ModelWriter: &MorpheStructFileWriter{
			TargetDirPath: path.Join(baseOutputDirPath, "models"),
		},
		ModelHooks: hook.CompileMorpheModel{},

		EntityWriter: &MorpheStructFileWriter{
			TargetDirPath: path.Join(baseOutputDirPath, "entities"),
		},
		EntityHooks: hook.CompileMorpheEntity{},

		WriteStructHooks: hook.WriteGoStruct{},
		WriteGoEnumHooks: hook.WriteGoEnum{},

		StructureWriter: &MorpheStructFileWriter{
			TargetDirPath: path.Join(baseOutputDirPath, "structures"),
		},
		StructureHooks: hook.CompileMorpheStructure{},
	}
}
