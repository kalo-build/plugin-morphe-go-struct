package compile_test

import (
	"fmt"
	"testing"

	"github.com/kalo-build/go/pkg/godef"
	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/compile"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/compile/hook"
	"github.com/stretchr/testify/suite"
)

type CompileStructuresTestSuite struct {
	suite.Suite
}

func TestCompileStructuresTestSuite(t *testing.T) {
	suite.Run(t, new(CompileStructuresTestSuite))
}

func (suite *CompileStructuresTestSuite) TestMorpheStructureToGoStruct() {
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheStructuresConfig: structuresConfig,
		},
		StructureHooks: hook.CompileMorpheStructure{},
	}

	structure0 := yaml.Structure{
		Name: "Basic",
		Fields: map[string]yaml.StructureField{
			"AutoIncrement": {
				Type: yaml.StructureFieldTypeAutoIncrement,
			},
			"Boolean": {
				Type: yaml.StructureFieldTypeBoolean,
			},
			"Date": {
				Type: yaml.StructureFieldTypeDate,
			},
			"Float": {
				Type: yaml.StructureFieldTypeFloat,
			},
			"Integer": {
				Type: yaml.StructureFieldTypeInteger,
			},
			"Protected": {
				Type: yaml.StructureFieldTypeProtected,
			},
			"Sealed": {
				Type: yaml.StructureFieldTypeSealed,
			},
			"String": {
				Type: yaml.StructureFieldTypeString,
			},
			"Time": {
				Type: yaml.StructureFieldTypeTime,
			},
			"UUID": {
				Type: yaml.StructureFieldTypeUUID,
				Attributes: []string{
					"immutable",
				},
			},
		},
	}

	r := registry.NewRegistry()

	structureStruct, structErr := compile.MorpheStructureToGoStruct(config, r, structure0)

	suite.Nil(structErr)
	suite.NotNil(structureStruct)

	suite.Equal(structureStruct.Package.Path, structuresConfig.Package.Path)
	suite.Equal(structureStruct.Package.Name, structuresConfig.Package.Name)
	suite.Equal(structureStruct.Name, "Basic")

	structFields := structureStruct.Fields
	suite.Len(structFields, 10)

	structField0 := structFields[0]
	suite.Equal(structField0.Name, "AutoIncrement")
	suite.Equal(structField0.Type, godef.GoTypeUint)
	suite.Len(structField0.Tags, 0)

	structField1 := structFields[1]
	suite.Equal(structField1.Name, "Boolean")
	suite.Equal(structField1.Type, godef.GoTypeBool)
	suite.Len(structField1.Tags, 0)

	structField2 := structFields[2]
	suite.Equal(structField2.Name, "Date")
	suite.Equal(structField2.Type, godef.GoTypeTime)
	suite.Len(structField2.Tags, 0)

	structField3 := structFields[3]
	suite.Equal(structField3.Name, "Float")
	suite.Equal(structField3.Type, godef.GoTypeFloat)
	suite.Len(structField3.Tags, 0)

	structField4 := structFields[4]
	suite.Equal(structField4.Name, "Integer")
	suite.Equal(structField4.Type, godef.GoTypeInt)
	suite.Len(structField4.Tags, 0)

	structField5 := structFields[5]
	suite.Equal(structField5.Name, "Protected")
	suite.Equal(structField5.Type, godef.GoTypeString)
	suite.Len(structField5.Tags, 0)

	structField6 := structFields[6]
	suite.Equal(structField6.Name, "Sealed")
	suite.Equal(structField6.Type, godef.GoTypeString)
	suite.Len(structField6.Tags, 0)

	structField7 := structFields[7]
	suite.Equal(structField7.Name, "String")
	suite.Equal(structField7.Type, godef.GoTypeString)
	suite.Len(structField7.Tags, 0)

	structField8 := structFields[8]
	suite.Equal(structField8.Name, "Time")
	suite.Equal(structField8.Type, godef.GoTypeTime)
	suite.Len(structField8.Tags, 0)

	structField9 := structFields[9]
	suite.Equal(structField9.Name, "UUID")
	suite.Equal(structField9.Type, godef.GoTypeString)
	suite.Equal(structField9.Tags, []string{`morphe:"immutable"`})

	suite.Len(structureStruct.Imports, 1)
	suite.Equal(structureStruct.Imports[0], "time")
}

func (suite *CompileStructuresTestSuite) TestMorpheStructureToGoStruct_NoPackagePath() {
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheStructuresConfig: structuresConfig,
		},
		StructureHooks: hook.CompileMorpheStructure{},
	}

	structure0 := yaml.Structure{
		Name: "Basic",
		Fields: map[string]yaml.StructureField{
			"String": {
				Type: yaml.StructureFieldTypeString,
			},
		},
	}

	r := registry.NewRegistry()

	structureStruct, structErr := compile.MorpheStructureToGoStruct(config, r, structure0)

	suite.NotNil(structErr)
	suite.ErrorContains(structErr, "structures package path cannot be empty")
	suite.Nil(structureStruct)
}

func (suite *CompileStructuresTestSuite) TestMorpheStructureToGoStruct_NoPackageName() {
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/structures",
			Name: "",
		},
		ReceiverName: "s",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheStructuresConfig: structuresConfig,
		},
		StructureHooks: hook.CompileMorpheStructure{},
	}

	structure0 := yaml.Structure{
		Name: "Basic",
		Fields: map[string]yaml.StructureField{
			"String": {
				Type: yaml.StructureFieldTypeString,
			},
		},
	}

	r := registry.NewRegistry()

	structureStruct, structErr := compile.MorpheStructureToGoStruct(config, r, structure0)

	suite.NotNil(structErr)
	suite.ErrorContains(structErr, "structures package name cannot be empty")
	suite.Nil(structureStruct)
}

func (suite *CompileStructuresTestSuite) TestMorpheStructureToGoStruct_NoReceiverName() {
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheStructuresConfig: structuresConfig,
		},
		StructureHooks: hook.CompileMorpheStructure{},
	}

	structure0 := yaml.Structure{
		Name: "Basic",
		Fields: map[string]yaml.StructureField{
			"String": {
				Type: yaml.StructureFieldTypeString,
			},
		},
	}

	r := registry.NewRegistry()

	structureStruct, structErr := compile.MorpheStructureToGoStruct(config, r, structure0)

	suite.NotNil(structErr)
	suite.ErrorContains(structErr, "structures method receiver name cannot be empty")
	suite.Nil(structureStruct)
}

func (suite *CompileStructuresTestSuite) TestMorpheStructureToGoStruct_NoStructureName() {
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheStructuresConfig: structuresConfig,
		},
		StructureHooks: hook.CompileMorpheStructure{},
	}

	structure0 := yaml.Structure{
		Name: "",
		Fields: map[string]yaml.StructureField{
			"String": {
				Type: yaml.StructureFieldTypeString,
			},
		},
	}

	r := registry.NewRegistry()

	structureStruct, structErr := compile.MorpheStructureToGoStruct(config, r, structure0)

	suite.NotNil(structErr)
	suite.ErrorContains(structErr, "morphe structure has no name")
	suite.Nil(structureStruct)
}

func (suite *CompileStructuresTestSuite) TestMorpheStructureToGoStruct_NoFields() {
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheStructuresConfig: structuresConfig,
		},
		StructureHooks: hook.CompileMorpheStructure{},
	}

	structure0 := yaml.Structure{
		Name:   "Basic",
		Fields: map[string]yaml.StructureField{},
	}

	r := registry.NewRegistry()

	structureStruct, structErr := compile.MorpheStructureToGoStruct(config, r, structure0)

	suite.NotNil(structErr)
	suite.ErrorContains(structErr, "morphe structure has no fields")
	suite.Nil(structureStruct)
}

func (suite *CompileStructuresTestSuite) TestMorpheStructureToGoStruct_StartHook_Successful() {
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheStructuresConfig: structuresConfig,
		},
		StructureHooks: hook.CompileMorpheStructure{
			OnCompileMorpheStructureStart: func(config cfg.MorpheConfig, structure yaml.Structure) (cfg.MorpheConfig, yaml.Structure, error) {
				config.MorpheStructuresConfig.ReceiverName = "CHANGED"
				structure.Name = structure.Name + "CHANGED"
				delete(structure.Fields, "Float")
				return config, structure, nil
			},
		},
	}

	structure0 := yaml.Structure{
		Name: "Basic",
		Fields: map[string]yaml.StructureField{
			"String": {
				Type: yaml.StructureFieldTypeString,
			},
			"Float": {
				Type: yaml.StructureFieldTypeFloat,
			},
		},
	}

	r := registry.NewRegistry()

	structureStruct, structErr := compile.MorpheStructureToGoStruct(config, r, structure0)

	suite.Nil(structErr)
	suite.NotNil(structureStruct)
	suite.Equal(structureStruct.Name, "BasicCHANGED")
	suite.Len(structureStruct.Fields, 1)

	structField0 := structureStruct.Fields[0]
	suite.Equal(structField0.Name, "String")
	suite.Equal(structField0.Type, godef.GoTypeString)
}

func (suite *CompileStructuresTestSuite) TestMorpheStructureToGoStruct_StartHook_Failure() {
	var featureFlag = "otherName"
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheStructuresConfig: structuresConfig,
		},
		StructureHooks: hook.CompileMorpheStructure{
			OnCompileMorpheStructureStart: func(config cfg.MorpheConfig, structure yaml.Structure) (cfg.MorpheConfig, yaml.Structure, error) {
				if featureFlag != "otherName" {
					return config, structure, nil
				}
				return config, structure, fmt.Errorf("compile structure start hook error")
			},
		},
	}

	structure0 := yaml.Structure{
		Name: "Basic",
		Fields: map[string]yaml.StructureField{
			"String": {
				Type: yaml.StructureFieldTypeString,
			},
		},
	}

	r := registry.NewRegistry()

	structureStruct, structErr := compile.MorpheStructureToGoStruct(config, r, structure0)

	suite.NotNil(structErr)
	suite.ErrorContains(structErr, "compile structure start hook error")
	suite.Nil(structureStruct)
}

func (suite *CompileStructuresTestSuite) TestMorpheStructureToGoStruct_SuccessHook_Successful() {
	var featureFlag = "otherName"
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheStructuresConfig: structuresConfig,
		},
		StructureHooks: hook.CompileMorpheStructure{
			OnCompileMorpheStructureSuccess: func(structureStruct *godef.Struct) (*godef.Struct, error) {
				if featureFlag != "otherName" {
					return structureStruct, nil
				}
				structureStruct.Name = structureStruct.Name + "CHANGED"
				newFields := []godef.StructField{}
				for _, structField := range structureStruct.Fields {
					if structField.Name == "Float" {
						continue
					}
					newFields = append(newFields, structField)
				}
				structureStruct.Fields = newFields
				return structureStruct, nil
			},
		},
	}

	structure0 := yaml.Structure{
		Name: "Basic",
		Fields: map[string]yaml.StructureField{
			"String": {
				Type: yaml.StructureFieldTypeString,
			},
			"Float": {
				Type: yaml.StructureFieldTypeFloat,
			},
		},
	}

	r := registry.NewRegistry()

	structureStruct, structErr := compile.MorpheStructureToGoStruct(config, r, structure0)

	suite.Nil(structErr)
	suite.NotNil(structureStruct)
	suite.Equal(structureStruct.Name, "BasicCHANGED")
	suite.Len(structureStruct.Fields, 1)

	structField0 := structureStruct.Fields[0]
	suite.Equal(structField0.Name, "String")
	suite.Equal(structField0.Type, godef.GoTypeString)
}

func (suite *CompileStructuresTestSuite) TestMorpheStructureToGoStruct_SuccessHook_Failure() {
	var featureFlag = "otherName"
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheStructuresConfig: structuresConfig,
		},
		StructureHooks: hook.CompileMorpheStructure{
			OnCompileMorpheStructureSuccess: func(structureStruct *godef.Struct) (*godef.Struct, error) {
				if featureFlag != "otherName" {
					return structureStruct, nil
				}
				return nil, fmt.Errorf("compile structure success hook error")
			},
		},
	}

	structure0 := yaml.Structure{
		Name: "Basic",
		Fields: map[string]yaml.StructureField{
			"String": {
				Type: yaml.StructureFieldTypeString,
			},
		},
	}

	r := registry.NewRegistry()

	structureStruct, structErr := compile.MorpheStructureToGoStruct(config, r, structure0)

	suite.NotNil(structErr)
	suite.ErrorContains(structErr, "compile structure success hook error")
	suite.Nil(structureStruct)
}

func (suite *CompileStructuresTestSuite) TestMorpheStructureToGoStruct_FailureHook_NoPackagePath() {
	structureHooks := hook.CompileMorpheStructure{
		OnCompileMorpheStructureFailure: func(config cfg.MorpheConfig, structure yaml.Structure, compileFailure error) error {
			return fmt.Errorf("Structure %s: %w", structure.Name, compileFailure)
		},
	}
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{

			MorpheStructuresConfig: structuresConfig,
		},
		StructureHooks: structureHooks,
	}

	structure0 := yaml.Structure{
		Name: "Basic",
		Fields: map[string]yaml.StructureField{
			"String": {
				Type: yaml.StructureFieldTypeString,
			},
		},
	}

	r := registry.NewRegistry()

	structureStruct, structErr := compile.MorpheStructureToGoStruct(config, r, structure0)

	suite.NotNil(structErr)
	suite.ErrorContains(structErr, "Structure Basic: structures package path cannot be empty")
	suite.Nil(structureStruct)
}

func (suite *CompileStructuresTestSuite) TestMorpheStructureToGoStruct_EnumField() {
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/enums",
			Name: "enums",
		},
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheStructuresConfig: structuresConfig,
			MorpheEnumsConfig:      enumsConfig,
		},
		StructureHooks: hook.CompileMorpheStructure{},
	}

	structure0 := yaml.Structure{
		Name: "Basic",
		Fields: map[string]yaml.StructureField{
			"String": {
				Type: yaml.StructureFieldTypeString,
			},
			"Nationality": {
				Type: "Nationality",
			},
		},
	}

	enum0 := yaml.Enum{
		Name: "Nationality",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"US": "American",
			"DE": "German",
			"FR": "French",
		},
	}

	r := registry.NewRegistry()
	r.SetEnum("Nationality", enum0)

	structureStruct, structErr := compile.MorpheStructureToGoStruct(config, r, structure0)

	suite.Nil(structErr)
	suite.NotNil(structureStruct)
	suite.Equal(structureStruct.Name, "Basic")

	fields := structureStruct.Fields
	suite.Len(fields, 2)

	field0 := fields[0]
	suite.Equal(field0.Name, "Nationality")
	suite.Equal(field0.Type, godef.GoTypeDerived{
		PackagePath: enumsConfig.Package.Path,
		Name:        "Nationality",
		BaseType:    godef.GoTypeString,
	})

	field1 := fields[1]
	suite.Equal(field1.Name, "String")
	suite.Equal(field1.Type, godef.GoTypeString)
}
