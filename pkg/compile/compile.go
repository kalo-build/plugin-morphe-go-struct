package compile

import (
	"github.com/kaloseia/morphe-go/pkg/registry"

	"github.com/kaloseia/plugin-morphe-go-struct/pkg/godef"
)

func MorpheToGoStructs(modelsDir string, entitiesDir string) error {
	r, rErr := loadRegistry(modelsDir, entitiesDir)
	if rErr != nil {
		return rErr
	}

	modelsPackageName := "models"
	allModelStructs := map[string]*godef.Struct{}
	for modelName, model := range r.Models {
		modelStruct, modelErr := MorpheModelToGoStruct(modelsPackageName, model)
		if modelErr != nil {
			return modelErr
		}
		allModelStructs[modelName] = modelStruct
	}

	entitiesPackageName := "entities"
	allEntityStructs := map[string]*godef.Struct{}
	for entityName, entity := range r.Entities {
		entityStruct, entityErr := MorpheEntityToGoStruct(entitiesPackageName, entity)
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
