package cfg

type MorpheConfig struct {
	MorpheModelsConfig
	MorpheStructuresConfig
	MorpheEnumsConfig
	MorpheEntitiesConfig
}

func (config MorpheConfig) Validate() error {
	modelsErr := config.MorpheModelsConfig.Validate()
	if modelsErr != nil {
		return modelsErr
	}

	structuresErr := config.MorpheStructuresConfig.Validate()
	if structuresErr != nil {
		return structuresErr
	}

	enumsErr := config.MorpheEnumsConfig.Validate()
	if enumsErr != nil {
		return enumsErr
	}

	entitiesErr := config.MorpheEntitiesConfig.Validate()
	if entitiesErr != nil {
		return entitiesErr
	}

	return nil
}
