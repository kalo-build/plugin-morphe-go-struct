package compile

import (
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/core"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/godef"
)

func MorpheToGoStructs(config MorpheCompileConfig) error {
	r, rErr := loadMorpheRegistry(config.RegistryHooks, config.MorpheLoadRegistryConfig)
	if rErr != nil {
		return rErr
	}

	allModelStructs := map[string][]*godef.Struct{}
	for modelName, model := range r.GetAllModels() {
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
			writeStructErr := config.ModelWriter.WriteStruct(modelStruct)
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
	for entityName, entity := range r.GetAllEntities() {
		entityStruct, entityErr := MorpheEntityToGoStruct(entitiesPackage, entity)
		if entityErr != nil {
			return entityErr
		}
		allEntityStructs[entityName] = entityStruct
	}

	return nil
}
