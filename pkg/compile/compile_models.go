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
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/hook"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/typemap"
)

func AllMorpheModelsToGoStructs(config MorpheCompileConfig, r *registry.Registry) (map[string][]*godef.Struct, error) {
	allModelStructDefs := map[string][]*godef.Struct{}
	for modelName, model := range r.GetAllModels() {
		modelStructs, modelErr := MorpheModelToGoStructs(config.ModelHooks, config.MorpheModelsConfig, model)
		if modelErr != nil {
			return nil, modelErr
		}
		allModelStructDefs[modelName] = modelStructs
	}
	return allModelStructDefs, nil
}

func MorpheModelToGoStructs(modelHooks hook.CompileMorpheModel, config cfg.MorpheModelsConfig, model yaml.Model) ([]*godef.Struct, error) {
	config, model, compileStartErr := triggerCompileMorpheModelStart(modelHooks, config, model)
	if compileStartErr != nil {
		return nil, triggerCompileMorpheModelFailure(modelHooks, config, model, compileStartErr)
	}
	allModelStructs, structsErr := morpheModelToGoStructs(config, model)
	if structsErr != nil {
		return nil, triggerCompileMorpheModelFailure(modelHooks, config, model, structsErr)
	}

	allModelStructs, compileSuccessErr := triggerCompileMorpheModelSuccess(modelHooks, allModelStructs)
	if compileSuccessErr != nil {
		return nil, triggerCompileMorpheModelFailure(modelHooks, config, model, compileSuccessErr)
	}
	return allModelStructs, nil
}

func morpheModelToGoStructs(config cfg.MorpheModelsConfig, model yaml.Model) ([]*godef.Struct, error) {
	validateConfigErr := config.Validate()
	if validateConfigErr != nil {
		return nil, validateConfigErr
	}
	validateMorpheErr := model.Validate()
	if validateMorpheErr != nil {
		return nil, validateMorpheErr
	}

	modelStruct, modelStructErr := getModelStruct(config.Package, model)
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

		identStruct, identStructErr := getIdentifierStruct(config.Package, modelStruct.Name, identifierName, allIdentFieldDefs)
		if identStructErr != nil {
			return nil, identStructErr
		}
		allModelStructs = append(allModelStructs, identStruct)

		modelIdentGetter, modelIdentErr := getModelIdentifierGetter(config, modelStruct.Name, identifierName, identStruct)
		if modelIdentErr != nil {
			return nil, modelIdentErr
		}
		modelStruct.Methods = append(modelStruct.Methods, modelIdentGetter)
	}
	return allModelStructs, nil
}

func getModelStruct(structPackage godef.Package, model yaml.Model) (*godef.Struct, error) {
	modelStruct := godef.Struct{
		Package: structPackage,
		Name:    model.Name,
	}
	structFields, fieldsErr := getGoFieldsForMorpheModel(model.Fields)
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

func getGoFieldsForMorpheModel(modelFields map[string]yaml.ModelField) ([]godef.StructField, error) {
	allFields := []godef.StructField{}

	allFieldNames := core.MapKeysSorted(modelFields)
	for _, fieldName := range allFieldNames {
		fieldDef := modelFields[fieldName]
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

func triggerCompileMorpheModelStart(hooks hook.CompileMorpheModel, config cfg.MorpheModelsConfig, model yaml.Model) (cfg.MorpheModelsConfig, yaml.Model, error) {
	if hooks.OnCompileMorpheModelStart == nil {
		return config, model, nil
	}

	updatedConfig, updatedModel, startErr := hooks.OnCompileMorpheModelStart(config, model)
	if startErr != nil {
		return cfg.MorpheModelsConfig{}, yaml.Model{}, startErr
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

func triggerCompileMorpheModelFailure(hooks hook.CompileMorpheModel, config cfg.MorpheModelsConfig, model yaml.Model, failureErr error) error {
	if hooks.OnCompileMorpheModelFailure == nil {
		return failureErr
	}

	return hooks.OnCompileMorpheModelFailure(config, model.DeepClone(), failureErr)
}
