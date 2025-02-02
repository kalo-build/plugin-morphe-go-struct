package compile

import (
	"fmt"
	"sort"

	"github.com/kaloseia/clone"
	"github.com/kaloseia/go-util/core"
	"github.com/kaloseia/go-util/strcase"
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

	allModelStructs, structsErr := morpheModelToGoStructs(config, r, model)
	if structsErr != nil {
		return nil, triggerCompileMorpheModelFailure(config.ModelHooks, morpheConfig, model, structsErr)
	}

	allModelStructs, compileSuccessErr := triggerCompileMorpheModelSuccess(config.ModelHooks, allModelStructs)
	if compileSuccessErr != nil {
		return nil, triggerCompileMorpheModelFailure(config.ModelHooks, morpheConfig, model, compileSuccessErr)
	}
	return allModelStructs, nil
}

func morpheModelToGoStructs(config MorpheCompileConfig, r *registry.Registry, model yaml.Model) ([]*godef.Struct, error) {
	validateConfigErr := config.MorpheModelsConfig.Validate()
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

	modelIdentifiers := model.Identifiers
	allIdentifierNames := core.MapKeysSorted(modelIdentifiers)
	for _, identifierName := range allIdentifierNames {
		identifierDef := modelIdentifiers[identifierName]

		allIdentFieldDefs, identFieldDefsErr := getIdentifierStructFieldSubset(*modelStruct, identifierName, identifierDef)
		if identFieldDefsErr != nil {
			return nil, identFieldDefsErr
		}

		identStruct, identStructErr := getIdentifierStruct(config.MorpheModelsConfig.Package, modelStruct.Name, identifierName, allIdentFieldDefs)
		if identStructErr != nil {
			return nil, identStructErr
		}
		allModelStructs = append(allModelStructs, identStruct)

		modelIdentGetter, modelIdentErr := getModelIdentifierGetter(config.MorpheModelsConfig, modelStruct.Name, identifierName, identStruct)
		if modelIdentErr != nil {
			return nil, modelIdentErr
		}
		modelStruct.Methods = append(modelStruct.Methods, modelIdentGetter)
	}
	return allModelStructs, nil
}

func getModelStruct(config MorpheCompileConfig, r *registry.Registry, model yaml.Model) (*godef.Struct, error) {
	if r == nil {
		return nil, ErrNoRegistry
	}

	modelStruct := godef.Struct{
		Package: config.MorpheModelsConfig.Package,
		Name:    model.Name,
	}
	structFields, fieldsErr := getGoFieldsForMorpheModel(config, r, model)
	if fieldsErr != nil {
		return nil, fieldsErr
	}
	modelStruct.Fields = structFields

	structImports, importsErr := getImportsForStructFields(structFields)
	if importsErr != nil {
		return nil, importsErr
	}
	modelStruct.Imports = structImports

	return &modelStruct, nil
}

func getGoFieldsForMorpheModel(config MorpheCompileConfig, r *registry.Registry, model yaml.Model) ([]godef.StructField, error) {
	allFields, fieldErr := getDirectGoFieldsForMorpheModel(config.MorpheEnumsConfig.Package, r.GetAllEnums(), model.Fields)
	if fieldErr != nil {
		return nil, fieldErr
	}

	allRelatedFields, relatedErr := getRelatedGoFieldsForMorpheModel(config.MorpheModelsConfig.Package, r, model.Related)
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

		goFieldType, typeSupported := typemap.MorpheFieldToGoField[fieldDef.Type]
		if !typeSupported {
			return nil, ErrUnsupportedMorpheFieldType(fieldDef.Type)
		}

		goField := godef.StructField{
			Name: fieldName,
			Type: goFieldType,
			Tags: fieldDef.Attributes,
		}
		allFields = append(allFields, goField)
	}
	return allFields, nil
}

func getRelatedGoFieldsForMorpheModel(modelPackage godef.Package, r *registry.Registry, modelRelations map[string]yaml.ModelRelation) ([]godef.StructField, error) {
	allFields := []godef.StructField{}

	allRelatedModelNames := core.MapKeysSorted(modelRelations)
	for _, relatedModelName := range allRelatedModelNames {
		relatedModelDef, relatedModelDefErr := r.GetModel(relatedModelName)
		if relatedModelDefErr != nil {
			return nil, relatedModelDefErr
		}

		goIDField, goIDErr := getRelatedGoFieldForMorpheModelPrimaryID(relatedModelName, relatedModelDef)
		if goIDErr != nil {
			return nil, goIDErr
		}
		allFields = append(allFields, goIDField)

		goRelatedField := getRelatedGoFieldForMorpheModel(modelPackage.Path, relatedModelName)
		allFields = append(allFields, goRelatedField)
	}
	return allFields, nil
}

func getRelatedGoFieldForMorpheModelPrimaryID(relatedModelName string, relatedModelDef yaml.Model) (godef.StructField, error) {
	relatedPrimaryIDFieldName, relatedIDFieldNameErr := yamlops.GetModelPrimaryIdentifierFieldName(relatedModelDef)
	if relatedIDFieldNameErr != nil {
		return godef.StructField{}, fmt.Errorf("related %w", relatedIDFieldNameErr)
	}

	idFieldName := fmt.Sprintf("%s%s", relatedModelName, relatedPrimaryIDFieldName)
	relatedPrimaryIDFieldDef, relatedIDFieldDefErr := yamlops.GetModelFieldDefinitionByName(relatedModelDef, relatedPrimaryIDFieldName)
	if relatedIDFieldDefErr != nil {
		return godef.StructField{}, fmt.Errorf("related %w (primary identifier)", relatedIDFieldDefErr)
	}

	idFieldType, typeSupported := typemap.MorpheFieldToGoField[relatedPrimaryIDFieldDef.Type]
	if !typeSupported {
		return godef.StructField{}, ErrUnsupportedMorpheFieldType(relatedPrimaryIDFieldDef.Type)
	}

	return godef.StructField{
		Name: idFieldName,
		Type: idFieldType,
	}, nil
}

func getRelatedGoFieldForMorpheModel(modelPackagePath string, relatedModelName string) godef.StructField {
	return godef.StructField{
		Name: relatedModelName,
		Type: godef.GoTypePointer{
			ValueType: godef.GoTypeStruct{
				PackagePath: modelPackagePath,
				Name:        relatedModelName,
			},
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

func getIdentifierStructFieldSubset(modelStruct godef.Struct, identifierName string, identifier yaml.ModelIdentifier) ([]godef.StructField, error) {
	identifierFieldDefs := []godef.StructField{}
	for _, fieldName := range identifier.Fields {
		identifierFieldDef := godef.StructField{}
		for _, modelFieldDef := range modelStruct.Fields {
			if modelFieldDef.Name != fieldName {
				continue
			}
			identifierFieldDef = godef.StructField{
				Name: modelFieldDef.Name,
				Type: modelFieldDef.Type,
			}
		}
		if identifierFieldDef.Name == "" {
			return nil, ErrMissingMorpheIdentifierField(modelStruct.Name, identifierName, fieldName)
		}
		identifierFieldDefs = append(identifierFieldDefs, identifierFieldDef)
	}
	return identifierFieldDefs, nil
}

func getIdentifierStruct(structPackage godef.Package, modelName string, identifierName string, allIdentFieldDefs []godef.StructField) (*godef.Struct, error) {
	identifierStructImports, identifierImportsErr := getImportsForStructFields(allIdentFieldDefs)
	if identifierImportsErr != nil {
		return nil, identifierImportsErr
	}
	identifierStruct := godef.Struct{
		Package: structPackage,
		Imports: identifierStructImports,
		Name:    fmt.Sprintf("%sID%s", modelName, strcase.ToPascalCase(identifierName)),
		Fields:  allIdentFieldDefs,
	}
	return &identifierStruct, nil
}

func getModelIdentifierGetter(config cfg.MorpheModelsConfig, modelName string, identifierName string, identStruct *godef.Struct) (godef.StructMethod, error) {
	identStructType := godef.GoTypeStruct{
		PackagePath: config.Package.Path,
		Name:        identStruct.Name,
	}

	bodyLines := getModelIdentifierGetterBodyLines(identStruct, config.ReceiverName)

	modelIdentGetter := godef.StructMethod{
		ReceiverName: config.ReceiverName,
		ReceiverType: godef.GoTypeStruct{
			PackagePath: config.Package.Path,
			Name:        modelName,
		},
		Name: fmt.Sprintf("GetID%s", strcase.ToPascalCase(identifierName)),
		ReturnTypes: []godef.GoType{
			identStructType,
		},
		BodyLines: bodyLines,
	}
	return modelIdentGetter, nil
}

func getModelIdentifierGetterBodyLines(identStruct *godef.Struct, receiverName string) []string {
	bodyLines := []string{
		fmt.Sprintf(`	return %s{`, identStruct.Name),
	}
	for _, fieldDef := range identStruct.Fields {
		fieldLine := fmt.Sprintf(`		%s: %s.%s,`, fieldDef.Name, receiverName, fieldDef.Name)
		bodyLines = append(bodyLines, fieldLine)
	}
	bodyLines = append(bodyLines, `	}`)
	return bodyLines
}

func getImportsForStructFields(allFields []godef.StructField) ([]string, error) {
	structImportMap := map[string]any{}
	for _, fieldDef := range allFields {
		allFieldImports := fieldDef.Type.GetImports()
		for _, fieldImport := range allFieldImports {
			structImportMap[fieldImport] = nil
		}
	}

	allStructImports := []string{}
	for importPath := range structImportMap {
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
