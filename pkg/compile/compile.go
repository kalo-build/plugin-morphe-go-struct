package compile

import (
	"fmt"

	"github.com/kalo-build/morphe-go/pkg/registry"
)

func MorpheToGo(config MorpheCompileConfig) error {
	r, rErr := registry.LoadMorpheRegistry(config.RegistryHooks, config.MorpheLoadRegistryConfig)
	if rErr != nil {
		return rErr
	}

	hasEnums := r.HasEnums()
	if hasEnums {
		allEnumDefs, compileAllErr := AllMorpheEnumsToGoEnums(config, r)
		if compileAllErr != nil {
			return compileAllErr
		}

		_, writeEnumsErr := WriteAllEnumDefinitions(config, allEnumDefs)
		if writeEnumsErr != nil {
			return writeEnumsErr
		}
	}

	hasModels := r.HasModels()
	if hasModels {
		allModelStructDefs, compileAllErr := AllMorpheModelsToGoStructs(config, r)
		if compileAllErr != nil {
			return compileAllErr
		}

		_, writeModelStructsErr := WriteAllModelStructDefinitions(config, allModelStructDefs)
		if writeModelStructsErr != nil {
			return writeModelStructsErr
		}
	}

	hasStructures := r.HasStructures()
	if hasStructures {
		allStructureStructDefs, compileAllErr := AllMorpheStructuresToGoStructs(config, r)
		if compileAllErr != nil {
			return compileAllErr
		}

		_, writeStructureStructsErr := WriteAllStructureStructDefinitions(config, allStructureStructDefs)
		if writeStructureStructsErr != nil {
			return writeStructureStructsErr
		}
	}

	hasEntities := r.HasEntities()
	if hasEntities {
		if !hasModels {
			return fmt.Errorf("entities compilation requires models to be compiled")
		}

		allEntityStructDefs, compileAllErr := AllMorpheEntitiesToGoStructs(config, r)
		if compileAllErr != nil {
			return compileAllErr
		}

		_, writeEntityStructsErr := WriteAllEntityStructDefinitions(config, allEntityStructDefs)
		if writeEntityStructsErr != nil {
			return writeEntityStructsErr
		}
	}

	return nil
}
