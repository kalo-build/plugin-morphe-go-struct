package compile

import "github.com/kaloseia/morphe-go/pkg/registry"

func MorpheToGo(config MorpheCompileConfig) error {
	r, rErr := registry.LoadMorpheRegistry(config.RegistryHooks, config.MorpheLoadRegistryConfig)
	if rErr != nil {
		return rErr
	}

	allEnumDefs, compileAllErr := AllMorpheEnumsToGoEnums(config, r)
	if compileAllErr != nil {
		return compileAllErr
	}

	_, writeEnumsErr := WriteAllEnumDefinitions(config, allEnumDefs)
	if writeEnumsErr != nil {
		return writeEnumsErr
	}

	allModelStructDefs, compileAllErr := AllMorpheModelsToGoStructs(config, r)
	if compileAllErr != nil {
		return compileAllErr
	}

	_, writeModelStructsErr := WriteAllModelStructDefinitions(config, allModelStructDefs)
	if writeModelStructsErr != nil {
		return writeModelStructsErr
	}

	allStructureStructDefs, compileAllErr := AllMorpheStructuresToGoStructs(config, r)
	if compileAllErr != nil {
		return compileAllErr
	}

	_, writeStructureStructsErr := WriteAllStructureStructDefinitions(config, allStructureStructDefs)
	if writeStructureStructsErr != nil {
		return writeStructureStructsErr
	}

	return nil
}
