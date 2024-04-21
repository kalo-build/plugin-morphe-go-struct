package compile

import (
	"fmt"
	"sort"

	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/core"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/godef"
"github.com/kaloseia/plugin-morphe-go-struct/pkg/strcase"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/typemap"
)

func MorpheModelToGoStructs(config ModelsConfig, model yaml.Model) ([]*godef.Struct, error) {
	validateConfigErr := validateModelsConfig(config)
	if validateConfigErr != nil {
		return nil, validateConfigErr
	}
	validateMorpheErr := validateMorpheModelDefinition(model)
	if validateMorpheErr != nil {
		return nil, validateMorpheErr
	}

	modelStruct, modelStructErr := getModelStruct(config.PackageName, model)
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

		identStruct, identStructErr := getIdentifierStruct(config.PackageName, modelStruct.Name, identifierName, allIdentFieldDefs)
		if identStructErr != nil {
			return nil, identStructErr
		}
		allModelStructs = append(allModelStructs, identStruct)

		modelIdentGetter, modelIdentErr := getModelIdentifierGetter(config, identifierName, identStruct)
		if modelIdentErr != nil {
			return nil, modelIdentErr
		}
		modelStruct.Methods = append(modelStruct.Methods, modelIdentGetter)
	}
	return allModelStructs, nil
}

func validateModelsConfig(config ModelsConfig) error {
	if config.PackagePath == "" {
		return fmt.Errorf("models %w", ErrNoPackagePath)
	}
	if config.PackageName == "" {
		return fmt.Errorf("models %w", ErrNoPackageName)
	}
	if config.ReceiverName == "" {
		return fmt.Errorf("models %w", ErrNoReceiverName)
	}
	return nil
}

func validateMorpheModelDefinition(model yaml.Model) error {
	if model.Name == "" {
		return ErrNoMorpheModelName
	}
	if len(model.Fields) == 0 {
		return ErrNoMorpheModelFields
	}
	if len(model.Identifiers) == 0 {
		return ErrNoMorpheModelIdentifiers
	}
	return nil
}

func getModelStruct(packageName string, model yaml.Model) (*godef.Struct, error) {
	modelStruct := godef.Struct{
		Package: packageName,
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

func getIdentifierStruct(packageName string, modelName string, identifierName string, allIdentFieldDefs []godef.StructField) (*godef.Struct, error) {
	identifierStructImports, identifierImportsErr := getImportsForStructFields(allIdentFieldDefs)
	if identifierImportsErr != nil {
		return nil, identifierImportsErr
	}
	identifierStruct := godef.Struct{
		Package: packageName,
		Imports: identifierStructImports,
		Name:    fmt.Sprintf("%sID%s", modelName, strcase.ToCamelCase(identifierName)),
		Fields:  allIdentFieldDefs,
	}
	return &identifierStruct, nil
}

func getModelIdentifierGetter(config ModelsConfig, identifierName string, identStruct *godef.Struct) (godef.StructMethod, error) {
	identStructType := godef.GoTypeStruct{
		PackagePath: config.PackagePath,
		Name:        identStruct.Name,
	}

	bodyLines := getModelIdentifierGetterBodyLines(identStruct, config.ReceiverName)

	modelIdentGetter := godef.StructMethod{
		ReceiverName: config.ReceiverName,
		ReceiverType: identStructType,
		Name:         fmt.Sprintf("GetID%s", strcase.ToCamelCase(identifierName)),
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
