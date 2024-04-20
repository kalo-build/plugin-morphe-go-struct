package compile

import (
	"fmt"
	"sort"

	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/core"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/godef"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/typemap"
)

func MorpheModelToGoStruct(packageName string, model yaml.Model) (*godef.Struct, error) {
	if packageName == "" {
		return nil, fmt.Errorf("model %w", ErrNoPackageName)
	}
	validateErr := validateMorpheModelDefinition(model)
	if validateErr != nil {
		return nil, validateErr
	}

	morpheStruct := godef.Struct{
		Package: packageName,
		Name:    model.Name,
	}
	structFields, fieldsErr := compileMorpheModelFieldsToGo(model.Fields)
	if fieldsErr != nil {
		return nil, fieldsErr
	}
	morpheStruct.Fields = structFields

	structImports, importsErr := compileMorpheModelImportsFromFields(structFields)
	if importsErr != nil {
		return nil, importsErr
	}
	morpheStruct.Imports = structImports

	return &morpheStruct, nil
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

func compileMorpheModelFieldsToGo(modelFields map[string]yaml.ModelField) ([]godef.StructField, error) {
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

func compileMorpheModelImportsFromFields(allFields []godef.StructField) ([]string, error) {
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
