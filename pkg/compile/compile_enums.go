package compile

import (
	"github.com/kaloseia/go-util/core"
	"github.com/kaloseia/go/pkg/godef"
	"github.com/kaloseia/morphe-go/pkg/registry"
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/hook"
)

func AllMorpheEnumsToGoEnums(config MorpheCompileConfig, r *registry.Registry) (map[string]*godef.Enum, error) {
	allEnumDefs := map[string]*godef.Enum{}
	for enumName, enum := range r.GetAllEnums() {
		enumType, enumErr := MorpheEnumToGoEnum(config.EnumHooks, config.MorpheEnumsConfig, enum)
		if enumErr != nil {
			return nil, enumErr
		}
		allEnumDefs[enumName] = enumType
	}
	return allEnumDefs, nil
}

func MorpheEnumToGoEnum(enumHooks hook.CompileMorpheEnum, config cfg.MorpheEnumsConfig, enum yaml.Enum) (*godef.Enum, error) {
	config, enum, compileStartErr := triggerCompileMorpheEnumStart(enumHooks, config, enum)
	if compileStartErr != nil {
		return nil, triggerCompileMorpheEnumFailure(enumHooks, config, enum, compileStartErr)
	}

	enumType, enumErr := morpheEnumToGoEnumType(config, enum)
	if enumErr != nil {
		return nil, triggerCompileMorpheEnumFailure(enumHooks, config, enum, enumErr)
	}

	enumType, compileSuccessErr := triggerCompileMorpheEnumSuccess(enumHooks, enumType)
	if compileSuccessErr != nil {
		return nil, triggerCompileMorpheEnumFailure(enumHooks, config, enum, compileSuccessErr)
	}
	return enumType, nil
}

func MorpheEnumTypeToGoType(enumPackage godef.Package, enumName string, morpheType yaml.EnumType) (godef.GoTypeDerived, error) {
	switch morpheType {
	case yaml.EnumTypeInteger:
		return godef.GoTypeDerived{
			PackagePath: enumPackage.Path,
			Name:        enumName,
			BaseType:    godef.GoTypeInt,
		}, nil
	case yaml.EnumTypeFloat:
		return godef.GoTypeDerived{
			PackagePath: enumPackage.Path,
			Name:        enumName,
			BaseType:    godef.GoTypeFloat,
		}, nil
	case yaml.EnumTypeString:
		return godef.GoTypeDerived{
			PackagePath: enumPackage.Path,
			Name:        enumName,
			BaseType:    godef.GoTypeString,
		}, nil
	default:
		return godef.GoTypeDerived{}, ErrUnsupportedEnumType(morpheType)
	}
}

func morpheEnumToGoEnumType(config cfg.MorpheEnumsConfig, enum yaml.Enum) (*godef.Enum, error) {
	validateConfigErr := config.Validate()
	if validateConfigErr != nil {
		return nil, validateConfigErr
	}
	validateMorpheErr := enum.Validate()
	if validateMorpheErr != nil {
		return nil, validateMorpheErr
	}

	enumType, enumTypeErr := getGoEnum(config, enum)
	if enumTypeErr != nil {
		return nil, enumTypeErr
	}

	return enumType, nil
}

func getGoEnum(config cfg.MorpheEnumsConfig, enum yaml.Enum) (*godef.Enum, error) {
	enumType := godef.Enum{
		Package: godef.Package{
			Path: config.Package.Path,
			Name: config.Package.Name,
		},
		Name: enum.Name,
	}

	goType, typeErr := MorpheEnumTypeToGoType(config.Package, enum.Name, enum.Type)
	if typeErr != nil {
		return nil, typeErr
	}
	enumType.Type = goType

	entries, entriesErr := getGoEntriesForMorpheEnum(enum.Name, enum.Entries)
	if entriesErr != nil {
		return nil, entriesErr
	}
	enumType.Entries = entries

	return &enumType, nil
}

func getGoEntriesForMorpheEnum(enumName string, entries map[string]any) ([]godef.EnumEntry, error) {
	goEntries := []godef.EnumEntry{}
	entryNames := core.MapKeysSorted(entries)

	for _, entryName := range entryNames {
		entryValue, entryExists := entries[entryName]
		if !entryExists {
			return nil, ErrEnumEntryNotFound(entryName)
		}
		goEntries = append(goEntries, godef.EnumEntry{
			Name:  enumName + entryName,
			Value: entryValue,
		})
	}
	return goEntries, nil
}

func triggerCompileMorpheEnumStart(hooks hook.CompileMorpheEnum, config cfg.MorpheEnumsConfig, enum yaml.Enum) (cfg.MorpheEnumsConfig, yaml.Enum, error) {
	if hooks.OnCompileMorpheEnumStart == nil {
		return config, enum, nil
	}

	updatedConfig, updatedEnum, startErr := hooks.OnCompileMorpheEnumStart(config, enum)
	if startErr != nil {
		return cfg.MorpheEnumsConfig{}, yaml.Enum{}, startErr
	}

	return updatedConfig, updatedEnum, nil
}

func triggerCompileMorpheEnumSuccess(hooks hook.CompileMorpheEnum, enumType *godef.Enum) (*godef.Enum, error) {
	if hooks.OnCompileMorpheEnumSuccess == nil {
		return enumType, nil
	}
	if enumType == nil {
		return nil, ErrNoEnumType
	}
	enumTypeClone := enumType.DeepClone()

	enumType, successErr := hooks.OnCompileMorpheEnumSuccess(&enumTypeClone)
	if successErr != nil {
		return nil, successErr
	}
	return enumType, nil
}

func triggerCompileMorpheEnumFailure(hooks hook.CompileMorpheEnum, config cfg.MorpheEnumsConfig, enum yaml.Enum, failureErr error) error {
	if hooks.OnCompileMorpheEnumFailure == nil {
		return failureErr
	}

	return hooks.OnCompileMorpheEnumFailure(config, enum.DeepClone(), failureErr)
}
