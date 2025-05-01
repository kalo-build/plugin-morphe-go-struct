package compile

import (
	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/go/pkg/godef"
	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/compile/hook"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/typemap"
)

func AllMorpheStructuresToGoStructs(config MorpheCompileConfig, r *registry.Registry) (map[string]*godef.Struct, error) {
	allStructureStructDefs := map[string]*godef.Struct{}
	for structureName, structure := range r.GetAllStructures() {
		structureStruct, structureErr := MorpheStructureToGoStruct(config, r, structure)
		if structureErr != nil {
			return nil, structureErr
		}
		allStructureStructDefs[structureName] = structureStruct
	}
	return allStructureStructDefs, nil
}

func MorpheStructureToGoStruct(config MorpheCompileConfig, r *registry.Registry, structure yaml.Structure) (*godef.Struct, error) {
	morpheConfig, structure, compileStartErr := triggerCompileMorpheStructureStart(config.StructureHooks, config.MorpheConfig, structure)
	if compileStartErr != nil {
		return nil, triggerCompileMorpheStructureFailure(config.StructureHooks, morpheConfig, structure, compileStartErr)
	}
	config.MorpheConfig = morpheConfig

	structureStruct, structErr := morpheStructureToGoStruct(config, r, structure)
	if structErr != nil {
		return nil, triggerCompileMorpheStructureFailure(config.StructureHooks, morpheConfig, structure, structErr)
	}

	structureStruct, compileSuccessErr := triggerCompileMorpheStructureSuccess(config.StructureHooks, structureStruct)
	if compileSuccessErr != nil {
		return nil, triggerCompileMorpheStructureFailure(config.StructureHooks, morpheConfig, structure, compileSuccessErr)
	}
	return structureStruct, nil
}

func morpheStructureToGoStruct(config MorpheCompileConfig, r *registry.Registry, structure yaml.Structure) (*godef.Struct, error) {
	validateConfigErr := config.MorpheStructuresConfig.Validate()
	if validateConfigErr != nil {
		return nil, validateConfigErr
	}
	validateMorpheErr := structure.Validate(r.GetAllEnums())
	if validateMorpheErr != nil {
		return nil, validateMorpheErr
	}

	structureStruct := godef.Struct{
		Package: config.MorpheStructuresConfig.Package,
		Name:    structure.Name,
	}

	structFields, fieldsErr := getGoFieldsForMorpheStructure(config.MorpheEnumsConfig.Package, r, structure)
	if fieldsErr != nil {
		return nil, fieldsErr
	}
	structureStruct.Fields = structFields

	structImports, importsErr := getImportsForStructFields(config.MorpheStructuresConfig.Package, structFields)
	if importsErr != nil {
		return nil, importsErr
	}
	structureStruct.Imports = structImports

	return &structureStruct, nil
}

func getGoFieldsForMorpheStructure(enumPackage godef.Package, r *registry.Registry, structure yaml.Structure) ([]godef.StructField, error) {
	if r == nil {
		return nil, ErrNoRegistry
	}

	allFields, fieldsErr := getDirectGoFieldsForMorpheStructure(enumPackage, r.GetAllEnums(), structure.Fields)
	if fieldsErr != nil {
		return nil, fieldsErr
	}

	return allFields, nil
}

func getDirectGoFieldsForMorpheStructure(enumPackage godef.Package, allEnums map[string]yaml.Enum, structureFields map[string]yaml.StructureField) ([]godef.StructField, error) {
	allFields := []godef.StructField{}

	allFieldNames := core.MapKeysSorted(structureFields)
	for _, fieldName := range allFieldNames {
		fieldDef := structureFields[fieldName]

		goEnumField := getEnumFieldAsStructFieldType(enumPackage, allEnums, fieldName, string(fieldDef.Type))
		if goEnumField.Name != "" && goEnumField.Type != nil {
			allFields = append(allFields, goEnumField)
			continue
		}

		goFieldType, typeSupported := typemap.MorpheStructureFieldToGoField[fieldDef.Type]
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

func triggerCompileMorpheStructureStart(hooks hook.CompileMorpheStructure, config cfg.MorpheConfig, structure yaml.Structure) (cfg.MorpheConfig, yaml.Structure, error) {
	if hooks.OnCompileMorpheStructureStart == nil {
		return config, structure, nil
	}
	return hooks.OnCompileMorpheStructureStart(config, structure)
}

func triggerCompileMorpheStructureSuccess(hooks hook.CompileMorpheStructure, structureStruct *godef.Struct) (*godef.Struct, error) {
	if hooks.OnCompileMorpheStructureSuccess == nil {
		return structureStruct, nil
	}

	if structureStruct == nil {
		return nil, ErrNoStructureStruct
	}

	structureStructClone := structureStruct.DeepClone()
	return hooks.OnCompileMorpheStructureSuccess(&structureStructClone)
}

func triggerCompileMorpheStructureFailure(hooks hook.CompileMorpheStructure, config cfg.MorpheConfig, structure yaml.Structure, compileErr error) error {
	if hooks.OnCompileMorpheStructureFailure == nil {
		return compileErr
	}
	return hooks.OnCompileMorpheStructureFailure(config, structure, compileErr)
}
