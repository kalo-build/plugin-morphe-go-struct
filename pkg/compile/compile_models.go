package compile

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kalo-build/clone"
	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/go/pkg/godef"
	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/morphe-go/pkg/yamlops"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/compile/hook"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/typemap"
)

func AllMorpheModelsToGoStructs(config MorpheCompileConfig, r *registry.Registry) (map[string][]*godef.Struct, error) {
	allModelStructDefs := map[string][]*godef.Struct{}
	for modelName, model := range r.GetAllModels() {
		modelStructs, modelErr := MorpheModelToGoStructs(config, r, model)
		if modelErr != nil {
			return nil, modelErr
		}
		allModelStructDefs[modelName] = modelStructs
	}
	return allModelStructDefs, nil
}

func MorpheModelToGoStructs(config MorpheCompileConfig, r *registry.Registry, model yaml.Model) ([]*godef.Struct, error) {
	morpheConfig, model, compileStartErr := triggerCompileMorpheModelStart(config.ModelHooks, config.MorpheConfig, model)
	if compileStartErr != nil {
		return nil, triggerCompileMorpheModelFailure(config.ModelHooks, morpheConfig, model, compileStartErr)
	}
	config.MorpheConfig = morpheConfig

	allModelStructs, structsErr := morpheModelToGoStructs(config.MorpheConfig, r, model)
	if structsErr != nil {
		return nil, triggerCompileMorpheModelFailure(config.ModelHooks, morpheConfig, model, structsErr)
	}

	allModelStructs, compileSuccessErr := triggerCompileMorpheModelSuccess(config.ModelHooks, allModelStructs)
	if compileSuccessErr != nil {
		return nil, triggerCompileMorpheModelFailure(config.ModelHooks, morpheConfig, model, compileSuccessErr)
	}
	return allModelStructs, nil
}

func morpheModelToGoStructs(config cfg.MorpheConfig, r *registry.Registry, model yaml.Model) ([]*godef.Struct, error) {
	validateConfigErr := config.Validate()
	if validateConfigErr != nil {
		return nil, validateConfigErr
	}
	validateMorpheErr := model.Validate(r.GetAllEnums())
	if validateMorpheErr != nil {
		return nil, validateMorpheErr
	}

	// Validate aliased relationships
	validateAliasErr := validateAliasedRelations(r, model)
	if validateAliasErr != nil {
		return nil, validateAliasErr
	}

	modelStruct, modelStructErr := getModelStruct(config, r, model)
	if modelStructErr != nil {
		return nil, modelStructErr
	}
	allModelStructs := []*godef.Struct{
		modelStruct,
	}

	identifierStructs, identifierErr := getAllModelIdentifierStructs(config.MorpheModelsConfig, model, modelStruct)
	if identifierErr != nil {
		return nil, identifierErr
	}
	allModelStructs = append(allModelStructs, identifierStructs...)
	return allModelStructs, nil
}

func getModelStruct(config cfg.MorpheConfig, r *registry.Registry, model yaml.Model) (*godef.Struct, error) {
	if r == nil {
		return nil, ErrNoRegistry
	}

	modelStruct := godef.Struct{
		Package: config.MorpheModelsConfig.Package,
		Name:    model.Name,
	}
	structFields, fieldsErr := getGoFieldsForMorpheModel(config.MorpheEnumsConfig.Package, r, model, config.MorpheModelsConfig.FieldCasing)
	if fieldsErr != nil {
		return nil, fieldsErr
	}
	modelStruct.Fields = structFields

	structImports, importsErr := getImportsForStructFields(config.MorpheModelsConfig.Package, structFields)
	if importsErr != nil {
		return nil, importsErr
	}
	modelStruct.Imports = structImports

	return &modelStruct, nil
}

func getGoFieldsForMorpheModel(enumPackage godef.Package, r *registry.Registry, model yaml.Model, fieldCasing cfg.Casing) ([]godef.StructField, error) {
	allFields, fieldErr := getDirectGoFieldsForMorpheModel(enumPackage, r.GetAllEnums(), model.Fields, fieldCasing)
	if fieldErr != nil {
		return nil, fieldErr
	}

	allRelatedFields, relatedErr := getRelatedGoFieldsForMorpheModel(r, model.Related, fieldCasing)
	if relatedErr != nil {
		return nil, relatedErr
	}

	allFields = append(allFields, allRelatedFields...)
	return allFields, nil
}

func getDirectGoFieldsForMorpheModel(enumPackage godef.Package, allEnums map[string]yaml.Enum, modelFields map[string]yaml.ModelField, fieldCasing cfg.Casing) ([]godef.StructField, error) {
	allFields := []godef.StructField{}

	allFieldNames := core.MapKeysSorted(modelFields)
	for _, fieldName := range allFieldNames {
		fieldDef := modelFields[fieldName]

		goEnumField := getEnumFieldAsStructFieldType(enumPackage, allEnums, fieldName, string(fieldDef.Type), fieldCasing)
		if goEnumField.Name != "" && goEnumField.Type != nil {
			if hasAttribute(fieldDef.Attributes, "optional") {
				goEnumField.Type = godef.GoTypePointer{ValueType: goEnumField.Type}
			}
			allFields = append(allFields, goEnumField)
			continue
		}

		goFieldType, typeSupported := typemap.MorpheModelFieldToGoField[fieldDef.Type]
		if !typeSupported {
			return nil, ErrUnsupportedMorpheFieldType(fieldDef.Type)
		}

		// Model fields are required by default; wrap in pointer for "optional" attribute
		if hasAttribute(fieldDef.Attributes, "optional") {
			goFieldType = godef.GoTypePointer{ValueType: goFieldType}
		}

		tags := buildFieldTags(fieldName, fieldDef.Attributes, fieldCasing)

		goField := godef.StructField{
			Name: fieldName,
			Type: goFieldType,
			Tags: tags,
		}
		allFields = append(allFields, goField)
	}
	return allFields, nil
}

// hasAttribute checks if the given attribute is present in the attributes list.
func hasAttribute(attributes []string, attr string) bool {
	for _, a := range attributes {
		if a == attr {
			return true
		}
	}
	return false
}

// buildFieldTags constructs the struct tags for a field
func buildFieldTags(fieldName string, attributes []string, fieldCasing cfg.Casing) []string {
	var tags []string

	// Add morphe tag if there are attributes
	if len(attributes) > 0 {
		tags = append(tags, fmt.Sprintf("morphe:\"%s\"", strings.Join(attributes, ";")))
	}

	// Add JSON tag based on casing configuration
	if fieldCasing != cfg.CasingNone {
		jsonName := fieldCasing.Apply(fieldName)
		tags = append(tags, fmt.Sprintf("json:\"%s\"", jsonName))
	}

	return tags
}

func getRelatedGoFieldsForMorpheModel(r *registry.Registry, modelRelations map[string]yaml.ModelRelation, fieldCasing cfg.Casing) ([]godef.StructField, error) {
	allFields := []godef.StructField{}

	allRelatedModelNames := core.MapKeysSorted(modelRelations)
	for _, relationshipName := range allRelatedModelNames {
		relationDef := modelRelations[relationshipName]

		// For polymorphic For* relationships (ForOnePoly/ForManyPoly),
		// we need to generate type and ID fields instead of looking up a model
		if yamlops.IsRelationPoly(relationDef.Type) && yamlops.IsRelationFor(relationDef.Type) {
			// Validate that For property is provided and has at least one model
			if len(relationDef.For) == 0 {
				return nil, fmt.Errorf("polymorphic relation '%s' must have at least one model in 'for' property", relationshipName)
			}

			// Generate polymorphic type field
			typeFieldName := relationshipName + "Type"
			typeField := godef.StructField{
				Name: typeFieldName,
				Type: godef.GoTypeString,
				Tags: buildFieldTags(typeFieldName, nil, fieldCasing),
			}
			allFields = append(allFields, typeField)

			// Generate polymorphic ID field
			idFieldName := relationshipName + "ID"
			idField := godef.StructField{
				Name: idFieldName,
				Type: godef.GoTypeString,
				Tags: buildFieldTags(idFieldName, nil, fieldCasing),
			}
			allFields = append(allFields, idField)

			// No need to generate the relationship field for ForOnePoly/ForManyPoly
			// as it's a polymorphic relationship and can't be strongly typed
			continue
		}

		// For polymorphic Has* relationships (HasOnePoly/HasManyPoly),
		// the relationship name is used as field name, and aliased specifies the target model
		if yamlops.IsRelationPoly(relationDef.Type) && yamlops.IsRelationHas(relationDef.Type) {
			// Use the aliased value if provided, otherwise use the relationship name (same as regular relationships)
			targetModelName := yamlops.GetRelationTargetName(relationshipName, relationDef.Aliased)

			// Validate that Through property references a valid polymorphic relationship
			if relationDef.Through != "" {
				// Get the target model to validate the Through relationship exists
				relatedModelDef, relatedModelDefErr := r.GetModel(targetModelName)
				if relatedModelDefErr != nil {
					return nil, relatedModelDefErr
				}

				// Check if the Through relationship exists on the target model
				throughRelation, throughExists := relatedModelDef.Related[relationDef.Through]
				if !throughExists {
					return nil, fmt.Errorf("polymorphic relation '%s' has invalid 'through' property: relation '%s' not found on model '%s'", relationshipName, relationDef.Through, targetModelName)
				}

				// Verify the Through relationship is a polymorphic For* relationship
				if !yamlops.IsRelationPoly(throughRelation.Type) || !yamlops.IsRelationFor(throughRelation.Type) {
					return nil, fmt.Errorf("polymorphic relation '%s' has invalid 'through' property: relation '%s' must be a polymorphic For* relationship", relationshipName, relationDef.Through)
				}
			}

			relatedModelDef, relatedModelDefErr := r.GetModel(targetModelName)
			if relatedModelDefErr != nil {
				return nil, relatedModelDefErr
			}

			goIDField, goIDErr := getRelatedGoFieldForMorpheModelPrimaryID(relationshipName, targetModelName, relatedModelDef, relationDef.Type, fieldCasing)
			if goIDErr != nil {
				return nil, goIDErr
			}
			allFields = append(allFields, goIDField)

			goRelatedField := getRelatedGoFieldForMorpheModel(relationshipName, targetModelName, relationDef.Type, fieldCasing)
			allFields = append(allFields, goRelatedField)
			continue
		}

		// Resolve actual target model name (handles aliasing)
		targetAlias := yamlops.GetRelationTargetName(relationshipName, relationDef.Aliased)

		// For HasMany relationships with path aliases (e.g., "Person.WorkProject"),
		// extract just the model name (first part before the dot)
		targetModelName := targetAlias
		if yamlops.IsRelationMany(relationDef.Type) && yamlops.IsRelationHas(relationDef.Type) {
			if parts := strings.Split(targetAlias, "."); len(parts) > 1 {
				targetModelName = parts[0]
			}
		}

		relatedModelDef, relatedModelDefErr := r.GetModel(targetModelName)
		if relatedModelDefErr != nil {
			return nil, fmt.Errorf("failed to get model '%s' for relation '%s': %w", targetModelName, relationshipName, relatedModelDefErr)
		}

		goIDField, goIDErr := getRelatedGoFieldForMorpheModelPrimaryID(relationshipName, targetModelName, relatedModelDef, relationDef.Type, fieldCasing)
		if goIDErr != nil {
			return nil, goIDErr
		}
		allFields = append(allFields, goIDField)

		goRelatedField := getRelatedGoFieldForMorpheModel(relationshipName, targetModelName, relationDef.Type, fieldCasing)
		allFields = append(allFields, goRelatedField)
	}
	return allFields, nil
}

func getRelatedGoFieldForMorpheModelPrimaryID(relationshipName, targetModelName string, relatedModelDef yaml.Model, relationType string, fieldCasing cfg.Casing) (godef.StructField, error) {
	relatedPrimaryIDFieldName, relatedIDFieldNameErr := yamlops.GetModelPrimaryIdentifierFieldName(relatedModelDef)
	if relatedIDFieldNameErr != nil {
		return godef.StructField{}, fmt.Errorf("related %w", relatedIDFieldNameErr)
	}

	// Use relationship name for field naming (semantic), not target model name
	idFieldName := fmt.Sprintf("%s%s", relationshipName, relatedPrimaryIDFieldName)
	if yamlops.IsRelationMany(relationType) {
		idFieldName += "s"
	}

	relatedPrimaryIDFieldDef, relatedIDFieldDefErr := yamlops.GetModelFieldDefinitionByName(relatedModelDef, relatedPrimaryIDFieldName)
	if relatedIDFieldDefErr != nil {
		return godef.StructField{}, fmt.Errorf("related %w (primary identifier)", relatedIDFieldDefErr)
	}

	idFieldType, typeSupported := typemap.MorpheModelFieldToGoField[relatedPrimaryIDFieldDef.Type]
	if !typeSupported {
		return godef.StructField{}, ErrUnsupportedMorpheFieldType(relatedPrimaryIDFieldDef.Type)
	}

	if yamlops.IsRelationMany(relationType) {
		return godef.StructField{
			Name: idFieldName,
			Type: godef.GoTypeArray{
				IsSlice:   true,
				ValueType: idFieldType,
			},
			Tags: buildFieldTags(idFieldName, nil, fieldCasing),
		}, nil
	}

	return godef.StructField{
		Name: idFieldName,
		Type: godef.GoTypePointer{
			ValueType: idFieldType,
		},
		Tags: buildFieldTags(idFieldName, nil, fieldCasing),
	}, nil
}

func getRelatedGoFieldForMorpheModel(relationshipName, targetModelName string, relationType string, fieldCasing cfg.Casing) godef.StructField {
	// Use relationship name for field naming (semantic)
	fieldName := relationshipName
	if yamlops.IsRelationMany(relationType) {
		fieldName += "s"
	}

	// Use target model name for struct type reference
	valueType := godef.GoTypeStruct{
		Name: targetModelName,
	}

	if yamlops.IsRelationMany(relationType) {
		return godef.StructField{
			Name: fieldName,
			Type: godef.GoTypeArray{
				IsSlice:   true,
				ValueType: valueType,
			},
			Tags: buildFieldTags(fieldName, nil, fieldCasing),
		}
	}

	return godef.StructField{
		Name: fieldName,
		Type: godef.GoTypePointer{
			ValueType: valueType,
		},
		Tags: buildFieldTags(fieldName, nil, fieldCasing),
	}
}

func getEnumFieldAsStructFieldType(enumPackage godef.Package, allEnums map[string]yaml.Enum, fieldName string, enumName string, fieldCasing cfg.Casing) godef.StructField {
	if len(allEnums) == 0 {
		return godef.StructField{}
	}

	enumType, enumTypeExists := allEnums[enumName]
	if !enumTypeExists {
		return godef.StructField{}
	}

	goFieldType, conversionErr := MorpheEnumTypeToGoType(enumPackage, enumType.Name, enumType.Type)
	if conversionErr != nil {
		return godef.StructField{}
	}

	goField := godef.StructField{
		Name: fieldName,
		Type: goFieldType,
		Tags: buildFieldTags(fieldName, nil, fieldCasing),
	}

	return goField
}

func getAllModelIdentifierStructs(config cfg.MorpheModelsConfig, model yaml.Model, modelStruct *godef.Struct) ([]*godef.Struct, error) {
	return GetIdentifierStructs(
		config,
		modelStruct.Name,
		modelStruct,
		wrapModelIdentifiers(model.Identifiers),
	)
}

// Adapter to make ModelIdentifier implement Identifier interface
type modelIdentifierWrapper struct {
	yaml.ModelIdentifier
}

func (m modelIdentifierWrapper) GetFields() []string {
	return m.Fields
}

func wrapModelIdentifiers(identifiers map[string]yaml.ModelIdentifier) map[string]Identifier {
	wrapped := make(map[string]Identifier)
	for k, v := range identifiers {
		wrapped[k] = modelIdentifierWrapper{v}
	}
	return wrapped
}

func getImportsForStructFields(structPackage godef.Package, allFields []godef.StructField) ([]string, error) {
	structImportMap := map[string]any{}
	for _, fieldDef := range allFields {
		allFieldImports := fieldDef.Type.GetImports()
		for _, fieldImport := range allFieldImports {
			if fieldImport == "" || fieldImport == structPackage.Path {
				continue
			}
			structImportMap[fieldImport] = nil
		}
	}

	allStructImports := []string{}
	for importPath := range structImportMap {
		if importPath == "" || importPath == structPackage.Path {
			continue
		}
		allStructImports = append(allStructImports, importPath)
	}
	sort.Strings(allStructImports)

	return allStructImports, nil
}

func triggerCompileMorpheModelStart(modelHooks hook.CompileMorpheModel, config cfg.MorpheConfig, model yaml.Model) (cfg.MorpheConfig, yaml.Model, error) {
	if modelHooks.OnCompileMorpheModelStart == nil {
		return config, model, nil
	}

	updatedConfig, updatedModel, startErr := modelHooks.OnCompileMorpheModelStart(config, model)
	if startErr != nil {
		return cfg.MorpheConfig{}, yaml.Model{}, startErr
	}

	return updatedConfig, updatedModel, nil
}

func triggerCompileMorpheModelSuccess(hooks hook.CompileMorpheModel, allModelStructs []*godef.Struct) ([]*godef.Struct, error) {
	if hooks.OnCompileMorpheModelSuccess == nil {
		return allModelStructs, nil
	}
	if allModelStructs == nil {
		return nil, ErrNoModelStructs
	}
	allModelStructsClone := clone.DeepCloneSlicePointers(allModelStructs)

	allModelStructs, successErr := hooks.OnCompileMorpheModelSuccess(allModelStructsClone)
	if successErr != nil {
		return nil, successErr
	}
	return allModelStructs, nil
}

func triggerCompileMorpheModelFailure(hooks hook.CompileMorpheModel, morpheConfig cfg.MorpheConfig, model yaml.Model, failureErr error) error {
	if hooks.OnCompileMorpheModelFailure == nil {
		return failureErr
	}

	return hooks.OnCompileMorpheModelFailure(morpheConfig, model.DeepClone(), failureErr)
}

// validateAliasedRelations validates that all aliased target models exist in the registry
func validateAliasedRelations(r *registry.Registry, model yaml.Model) error {
	for relationshipName, relation := range model.Related {
		targetAlias := yamlops.GetRelationTargetName(relationshipName, relation.Aliased)

		// For HasMany relationships with path aliases (e.g., "Person.WorkProject"),
		// extract just the model name (first part before the dot)
		targetModelName := targetAlias
		if yamlops.IsRelationMany(relation.Type) && yamlops.IsRelationHas(relation.Type) {
			if parts := strings.Split(targetAlias, "."); len(parts) > 1 {
				targetModelName = parts[0]
			}
		}

		// Only validate if the target is different (i.e., aliased)
		if targetModelName != relationshipName {
			// Skip validation for polymorphic relationships - aliases are conceptual names
			if yamlops.IsRelationPoly(relation.Type) {
				continue
			}

			_, err := r.GetModel(targetModelName)
			if err != nil {
				return fmt.Errorf("aliased target model '%s' for relation '%s' in model '%s' not found",
					targetModelName, relationshipName, model.Name)
			}
		}
	}
	return nil
}
