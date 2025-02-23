package compile

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kaloseia/clone"
	"github.com/kaloseia/go-util/core"
	"github.com/kaloseia/go/pkg/godef"
	"github.com/kaloseia/morphe-go/pkg/registry"
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/morphe-go/pkg/yamlops"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/hook"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/typemap"
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
	structFields, fieldsErr := getGoFieldsForMorpheModel(config.MorpheEnumsConfig.Package, r, model)
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

func getGoFieldsForMorpheModel(enumPackage godef.Package, r *registry.Registry, model yaml.Model) ([]godef.StructField, error) {
	allFields, fieldErr := getDirectGoFieldsForMorpheModel(enumPackage, r.GetAllEnums(), model.Fields)
	if fieldErr != nil {
		return nil, fieldErr
	}

	allRelatedFields, relatedErr := getRelatedGoFieldsForMorpheModel(r, model.Related)
	if relatedErr != nil {
		return nil, relatedErr
	}

	allFields = append(allFields, allRelatedFields...)
	return allFields, nil
}

func getDirectGoFieldsForMorpheModel(enumPackage godef.Package, allEnums map[string]yaml.Enum, modelFields map[string]yaml.ModelField) ([]godef.StructField, error) {
	allFields := []godef.StructField{}

	allFieldNames := core.MapKeysSorted(modelFields)
	for _, fieldName := range allFieldNames {
		fieldDef := modelFields[fieldName]

		goEnumField := getEnumFieldAsStructFieldType(enumPackage, allEnums, fieldName, string(fieldDef.Type))
		if goEnumField.Name != "" && goEnumField.Type != nil {
			allFields = append(allFields, goEnumField)
			continue
		}

		goFieldType, typeSupported := typemap.MorpheModelFieldToGoField[fieldDef.Type]
		if !typeSupported {
			return nil, ErrUnsupportedMorpheFieldType(fieldDef.Type)
		}

		var tags []string
		if len(fieldDef.Attributes) > 0 {
			tags = []string{fmt.Sprintf("morphe:\"%s\"", strings.Join(fieldDef.Attributes, ";"))}
		}

		goField := godef.StructField{
			Name: fieldName,
			Type: goFieldType,
			Tags: tags,
		}
		allFields = append(allFields, goField)
	}
	return allFields, nil
}

func getRelatedGoFieldsForMorpheModel(r *registry.Registry, modelRelations map[string]yaml.ModelRelation) ([]godef.StructField, error) {
	allFields := []godef.StructField{}

	allRelatedModelNames := core.MapKeysSorted(modelRelations)
	for _, relatedModelName := range allRelatedModelNames {
		relationDef := modelRelations[relatedModelName]
		relatedModelDef, relatedModelDefErr := r.GetModel(relatedModelName)
		if relatedModelDefErr != nil {
			return nil, relatedModelDefErr
		}

		goIDField, goIDErr := getRelatedGoFieldForMorpheModelPrimaryID(relatedModelName, relatedModelDef, relationDef.Type)
		if goIDErr != nil {
			return nil, goIDErr
		}
		allFields = append(allFields, goIDField)

		goRelatedField := getRelatedGoFieldForMorpheModel(relatedModelName, relationDef.Type)
		allFields = append(allFields, goRelatedField)
	}
	return allFields, nil
}

func getRelatedGoFieldForMorpheModelPrimaryID(relatedModelName string, relatedModelDef yaml.Model, relationType string) (godef.StructField, error) {
	relatedPrimaryIDFieldName, relatedIDFieldNameErr := yamlops.GetModelPrimaryIdentifierFieldName(relatedModelDef)
	if relatedIDFieldNameErr != nil {
		return godef.StructField{}, fmt.Errorf("related %w", relatedIDFieldNameErr)
	}

	idFieldName := fmt.Sprintf("%s%s", relatedModelName, relatedPrimaryIDFieldName)
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
		}, nil
	}

	return godef.StructField{
		Name: idFieldName,
		Type: godef.GoTypePointer{
			ValueType: idFieldType,
		},
	}, nil
}

func getRelatedGoFieldForMorpheModel(relatedModelName string, relationType string) godef.StructField {
	fieldName := relatedModelName
	if yamlops.IsRelationMany(relationType) {
		fieldName += "s"
	}

	valueType := godef.GoTypeStruct{
		Name: relatedModelName,
	}

	if yamlops.IsRelationMany(relationType) {
		return godef.StructField{
			Name: fieldName,
			Type: godef.GoTypeArray{
				IsSlice:   true,
				ValueType: valueType,
			},
		}
	}

	return godef.StructField{
		Name: fieldName,
		Type: godef.GoTypePointer{
			ValueType: valueType,
		},
	}
}

func getEnumFieldAsStructFieldType(enumPackage godef.Package, allEnums map[string]yaml.Enum, fieldName string, enumName string) godef.StructField {
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
