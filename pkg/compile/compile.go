package compile

import "github.com/kaloseia/morphe-go/pkg/registry"

func MorpheToGoStructs(config MorpheCompileConfig) (CompiledModelStructs, error) {
	r, rErr := registry.LoadMorpheRegistry(config.RegistryHooks, config.MorpheLoadRegistryConfig)
	if rErr != nil {
		return nil, rErr
	}

	allModelStructDefs, compileAllErr := AllMorpheModelsToGoStructs(config, r)
	if compileAllErr != nil {
		return nil, compileAllErr
	}

	allWrittenModels, writeAllErr := WriteAllModelStructDefinitions(config, allModelStructDefs)
	if writeAllErr != nil {
		return nil, writeAllErr
	}
	return allWrittenModels, nil
}
