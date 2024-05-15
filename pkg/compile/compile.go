package compile

import (
	"github.com/kaloseia/morphe-go/pkg/registry"

	"github.com/kaloseia/plugin-morphe-go-struct/pkg/core"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/event"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/godef"
)

type MorpheCompileConfig struct {
	MorpheLoadRegistryConfig
	MorpheModelsConfig

	ModelsWriter   StructWriter
	EntitiesWriter StructWriter
	EventBus       event.EventBusable
}

func MorpheToGoStructs(config MorpheCompileConfig) error {
	r, rErr := loadRegistry(config.RegistryModelsDirPath, config.RegistryEntitiesDirPath)
	if rErr != nil {
		return rErr
	}

	allModelStructs := map[string][]*godef.Struct{}
	for modelName, model := range r.Models {
		modelStructs, modelErr := MorpheModelToGoStructs(config.MorpheModelsConfig, model)
		if modelErr != nil {
			return modelErr
		}
		allModelStructs[modelName] = modelStructs
	}
	sortedModelNames := core.MapKeysSorted(allModelStructs)
	for _, modelName := range sortedModelNames {
		modelStructs := allModelStructs[modelName]
		for _, modelStruct := range modelStructs {
			writeStructErr := config.ModelsWriter.WriteStruct(modelStruct)
			if writeStructErr != nil {
				return writeStructErr
			}
		}
	}

	entitiesPackage := godef.Package{
		Path: "placeholder",
		Name: "entities",
	}
	allEntityStructs := map[string]*godef.Struct{}
	for entityName, entity := range r.Entities {
		entityStruct, entityErr := MorpheEntityToGoStruct(entitiesPackage, entity)
		if entityErr != nil {
			return entityErr
		}
		allEntityStructs[entityName] = entityStruct
	}

	return nil
}

func loadRegistry(modelsDir string, entitiesDir string) (*registry.Registry, error) {
	r := registry.NewRegistry()

	modelsErr := r.LoadModelsFromDirectory(modelsDir)
	if modelsErr != nil {
		return nil, modelsErr
	}
	entitiesErr := r.LoadEntitiesFromDirectory(entitiesDir)
	if entitiesErr != nil {
		return nil, entitiesErr
	}

	return r, nil
}
