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

type CompileEntitiesTestSuite struct {
	suite.Suite
}

func TestCompileEntitiesTestSuite(t *testing.T) {
	suite.Run(t, new(CompileEntitiesTestSuite))
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToGoStructs() {
	modelsConfig := cfg.MorpheModelsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
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
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/entities",
			Name: "entities",
		},
		ReceiverName: "e",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig:     modelsConfig,
			MorpheStructuresConfig: structuresConfig,
			MorpheEnumsConfig:      enumsConfig,
			MorpheEntitiesConfig:   entitiesConfig,
		},
		EntityHooks: hook.CompileMorpheEntity{},
	}

	entity0 := yaml.Entity{
		Name: "User",
		Fields: map[string]yaml.EntityField{
			"UUID": {
				Type:       "User.UUID",
				Attributes: []string{"immutable"},
			},
			"AutoIncrement": {
				Type: "User.Child.AutoIncrement",
			},
			"Boolean": {
				Type: "User.Child.Boolean",
			},
			"Date": {
				Type: "User.Child.Date",
			},
			"Float": {
				Type: "User.Child.Float",
			},
			"Integer": {
				Type: "User.Child.Integer",
			},
			"String": {
				Type: "User.Child.String",
			},
			"Time": {
				Type: "User.Child.Time",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()

	userModel := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"UUID": {
				Type:       yaml.ModelFieldTypeUUID,
				Attributes: []string{"immutable"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Child": {
				Type: "HasOne",
			},
		},
	}
	r.SetModel("User", userModel)

	childModel := yaml.Model{
		Name: "Child",
		Fields: map[string]yaml.ModelField{
			"UUID": {
				Type:       yaml.ModelFieldTypeUUID,
				Attributes: []string{"immutable"},
			},
			"AutoIncrement": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"Boolean": {
				Type: yaml.ModelFieldTypeBoolean,
			},
			"Date": {
				Type: yaml.ModelFieldTypeDate,
			},
			"Float": {
				Type: yaml.ModelFieldTypeFloat,
			},
			"Integer": {
				Type: yaml.ModelFieldTypeInteger,
			},
			"String": {
				Type: yaml.ModelFieldTypeString,
			},
			"Time": {
				Type: yaml.ModelFieldTypeTime,
			},
		},
		Related: map[string]yaml.ModelRelation{
			"User": {
				Type: "ForOne",
			},
		},
	}
	r.SetModel("Child", childModel)

	allGoStructs, allStructsErr := compile.MorpheEntityToGoStructs(config.EntityHooks, config.MorpheConfig, r, entity0)

	suite.Nil(allStructsErr)
	suite.Len(allGoStructs, 2)

	goStruct0 := allGoStructs[0]
	suite.Equal(goStruct0.Package.Path, entitiesConfig.Package.Path)
	suite.Equal(goStruct0.Package.Name, entitiesConfig.Package.Name)
	suite.Equal(goStruct0.Name, "User")

	structFields0 := goStruct0.Fields
	suite.Len(structFields0, 8)

	tsField00 := structFields0[0]
	suite.Equal("AutoIncrement", tsField00.Name)
	suite.Equal(godef.GoTypeUint, tsField00.Type)
	suite.Empty(tsField00.Tags)

	tsField01 := structFields0[1]
	suite.Equal("Boolean", tsField01.Name)
	suite.Equal(godef.GoTypeBool, tsField01.Type)
	suite.Empty(tsField01.Tags)

	tsField02 := structFields0[2]
	suite.Equal("Date", tsField02.Name)
	suite.Equal(godef.GoTypeTime, tsField02.Type)
	suite.Empty(tsField02.Tags)

	tsField03 := structFields0[3]
	suite.Equal("Float", tsField03.Name)
	suite.Equal(godef.GoTypeFloat, tsField03.Type)
	suite.Empty(tsField03.Tags)

	tsField04 := structFields0[4]
	suite.Equal("Integer", tsField04.Name)
	suite.Equal(godef.GoTypeInt, tsField04.Type)
	suite.Empty(tsField04.Tags)

	tsField05 := structFields0[5]
	suite.Equal("String", tsField05.Name)
	suite.Equal(godef.GoTypeString, tsField05.Type)
	suite.Empty(tsField05.Tags)

	tsField06 := structFields0[6]
	suite.Equal("Time", tsField06.Name)
	suite.Equal(godef.GoTypeTime, tsField06.Type)
	suite.Empty(tsField06.Tags)

	tsField07 := structFields0[7]
	suite.Equal("UUID", tsField07.Name)
	suite.Equal(godef.GoTypeString, tsField07.Type)
	suite.Len(tsField07.Tags, 1)
	suite.Equal(`morphe:"immutable"`, tsField07.Tags[0])

	goStruct1 := allGoStructs[1]
	suite.Equal(goStruct1.Package.Path, entitiesConfig.Package.Path)
	suite.Equal(goStruct1.Package.Name, entitiesConfig.Package.Name)
	suite.Equal(goStruct1.Name, "UserIDPrimary")

	structFields1 := goStruct1.Fields
	suite.Len(structFields1, 1)

	structFields10 := structFields1[0]
	suite.Equal(structFields10.Name, "UUID")
	suite.Equal(structFields10.Type, godef.GoTypeString)
	suite.Len(structFields10.Tags, 1)
	suite.Equal(`morphe:"immutable"`, structFields10.Tags[0])
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToGoStructs_NoEntityName() {
	modelsConfig := cfg.MorpheModelsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
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
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/entities",
			Name: "entities",
		},
		ReceiverName: "e",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig:     modelsConfig,
			MorpheStructuresConfig: structuresConfig,
			MorpheEnumsConfig:      enumsConfig,
			MorpheEntitiesConfig:   entitiesConfig,
		},
		EntityHooks: hook.CompileMorpheEntity{},
	}

	entity0 := yaml.Entity{
		Name: "",
		Fields: map[string]yaml.EntityField{
			"AutoIncrement": {
				Type: "User.AutoIncrement",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()

	allGoStructs, goStructErr := compile.MorpheEntityToGoStructs(config.EntityHooks, config.MorpheConfig, r, entity0)

	suite.NotNil(goStructErr)
	suite.ErrorContains(goStructErr, "entity has no name")
	suite.Nil(allGoStructs)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToGoStructs_NoFields() {
	modelsConfig := cfg.MorpheModelsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
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
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/entities",
			Name: "entities",
		},
		ReceiverName: "e",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig:     modelsConfig,
			MorpheStructuresConfig: structuresConfig,
			MorpheEnumsConfig:      enumsConfig,
			MorpheEntitiesConfig:   entitiesConfig,
		},
		EntityHooks: hook.CompileMorpheEntity{},
	}

	entity0 := yaml.Entity{
		Name:   "Basic",
		Fields: map[string]yaml.EntityField{},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("Basic", model0)

	allGoStructs, goStructErr := compile.MorpheEntityToGoStructs(config.EntityHooks, config.MorpheConfig, r, entity0)

	suite.NotNil(goStructErr)
	suite.ErrorContains(goStructErr, "morphe entity Basic has no fields")
	suite.Nil(allGoStructs)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToGoStructs_NoIdentifiers() {
	modelsConfig := cfg.MorpheModelsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
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
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/entities",
			Name: "entities",
		},
		ReceiverName: "e",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig:     modelsConfig,
			MorpheStructuresConfig: structuresConfig,
			MorpheEnumsConfig:      enumsConfig,
			MorpheEntitiesConfig:   entitiesConfig,
		},
		EntityHooks: hook.CompileMorpheEntity{},
	}

	entity0 := yaml.Entity{
		Name: "Basic",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Basic.ID",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{},
		Related:     map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("Basic", model0)

	allGoStructs, goStructErr := compile.MorpheEntityToGoStructs(config.EntityHooks, config.MorpheConfig, r, entity0)

	suite.NotNil(goStructErr)
	suite.ErrorContains(goStructErr, "entity 'Basic' has no identifiers")
	suite.Nil(allGoStructs)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToGoStructs_EnumField() {
	modelsConfig := cfg.MorpheModelsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
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
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/entities",
			Name: "entities",
		},
		ReceiverName: "e",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig:     modelsConfig,
			MorpheStructuresConfig: structuresConfig,
			MorpheEnumsConfig:      enumsConfig,
			MorpheEntitiesConfig:   entitiesConfig,
		},
		EntityHooks: hook.CompileMorpheEntity{},
	}

	entity0 := yaml.Entity{
		Name: "User",
		Fields: map[string]yaml.EntityField{
			"UUID": {
				Type: "User.UUID",
				Attributes: []string{
					"immutable",
				},
			},
			"Nationality": {
				Type: "User.Nationality",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()

	userModel := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
				},
			},
			"Nationality": {
				Type: "Nationality",
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("User", userModel)

	enum0 := yaml.Enum{
		Name: "Nationality",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"US": "American",
			"DE": "German",
			"FR": "French",
		},
	}
	r.SetEnum("Nationality", enum0)

	allGoStructs, goStructErr := compile.MorpheEntityToGoStructs(config.EntityHooks, config.MorpheConfig, r, entity0)

	suite.Nil(goStructErr)
	suite.Len(allGoStructs, 2)

	goStruct0 := allGoStructs[0]
	suite.Equal(goStruct0.Package.Path, entitiesConfig.Package.Path)
	suite.Equal(goStruct0.Package.Name, entitiesConfig.Package.Name)
	suite.Equal(goStruct0.Name, "User")

	structFields0 := goStruct0.Fields
	suite.Len(structFields0, 2)

	field00 := structFields0[0]
	suite.Equal(field00.Name, "Nationality")
	suite.Equal(field00.Type, godef.GoTypeDerived{
		PackagePath: enumsConfig.Package.Path,
		Name:        "Nationality",
		BaseType:    godef.GoTypeString,
	})
	suite.Empty(field00.Tags)

	field01 := structFields0[1]
	suite.Equal(field01.Name, "UUID")
	suite.Equal(field01.Type, godef.GoTypeString)
	suite.Len(field01.Tags, 1)
	suite.Equal(`morphe:"immutable"`, field01.Tags[0])

	goStruct1 := allGoStructs[1]
	suite.Equal(goStruct1.Package.Path, entitiesConfig.Package.Path)
	suite.Equal(goStruct1.Package.Name, entitiesConfig.Package.Name)
	suite.Equal(goStruct1.Name, "UserIDPrimary")

	structFields1 := goStruct1.Fields
	suite.Len(structFields1, 1)

	field10 := structFields1[0]
	suite.Equal(field10.Name, "UUID")
	suite.Equal(field10.Type, godef.GoTypeString)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToGoStructs_EnumField_EnumNotFound() {
	modelsConfig := cfg.MorpheModelsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
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
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/entities",
			Name: "entities",
		},
		ReceiverName: "e",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig:     modelsConfig,
			MorpheStructuresConfig: structuresConfig,
			MorpheEnumsConfig:      enumsConfig,
			MorpheEntitiesConfig:   entitiesConfig,
		},
		EntityHooks: hook.CompileMorpheEntity{},
	}

	entity0 := yaml.Entity{
		Name: "User",
		Fields: map[string]yaml.EntityField{
			"UUID": {
				Type: "User.UUID",
				Attributes: []string{
					"immutable",
				},
			},
			"Nationality": {
				Type: "User.Nationality",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()

	userModel := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
				},
			},
			"Nationality": {
				Type: "Nationality",
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("User", userModel)

	allGoStructs, goStructErr := compile.MorpheEntityToGoStructs(config.EntityHooks, config.MorpheConfig, r, entity0)

	suite.NotNil(goStructErr)
	suite.ErrorContains(goStructErr, "morphe entity 'User' field 'Nationality' has unknown non-primitive type 'Nationality'")
	suite.Nil(allGoStructs)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToGoStructs_Related_ForOne() {
	modelsConfig := cfg.MorpheModelsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
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
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/entities",
			Name: "entities",
		},
		ReceiverName: "e",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig:     modelsConfig,
			MorpheStructuresConfig: structuresConfig,
			MorpheEnumsConfig:      enumsConfig,
			MorpheEntitiesConfig:   entitiesConfig,
		},
		EntityHooks: hook.CompileMorpheEntity{},
	}

	entity0 := yaml.Entity{
		Name: "Basic",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Basic.ID",
			},
			"String": {
				Type: "Basic.String",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"BasicParent": {
				Type: "ForOne",
			},
		},
	}

	r := registry.NewRegistry()

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"String": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"BasicParent": {
				Type: "ForOne",
			},
		},
	}
	r.SetModel("Basic", model0)

	model1 := yaml.Model{
		Name: "BasicParent",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"String": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Basic": {
				Type: "HasOne",
			},
		},
	}
	r.SetModel("BasicParent", model1)

	entity1 := yaml.Entity{
		Name: "BasicParent",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "BasicParent.ID",
			},
			"String": {
				Type: "BasicParent.String",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"Basic": {
				Type: "HasMany",
			},
		},
	}
	r.SetEntity("BasicParent", entity1)

	allGoStructs, goStructErr := compile.MorpheEntityToGoStructs(config.EntityHooks, config.MorpheConfig, r, entity0)

	suite.Nil(goStructErr)
	suite.Len(allGoStructs, 2)

	goStruct0 := allGoStructs[0]
	suite.Equal(goStruct0.Package.Path, entitiesConfig.Package.Path)
	suite.Equal(goStruct0.Package.Name, entitiesConfig.Package.Name)
	suite.Equal(goStruct0.Name, "Basic")

	structFields0 := goStruct0.Fields
	suite.Len(structFields0, 4)

	field00 := structFields0[0]
	suite.Equal(field00.Name, "ID")
	suite.Equal(field00.Type, godef.GoTypeUint)

	field01 := structFields0[1]
	suite.Equal(field01.Name, "String")
	suite.Equal(field01.Type, godef.GoTypeString)

	field02 := structFields0[2]
	suite.Equal(field02.Name, "BasicParentID")
	suite.Equal(field02.Type, godef.GoTypePointer{
		ValueType: godef.GoTypeUint,
	})

	field03 := structFields0[3]
	suite.Equal(field03.Name, "BasicParent")
	suite.Equal(field03.Type, godef.GoTypePointer{
		ValueType: godef.GoTypeStruct{
			PackagePath: "",
			Name:        "BasicParent",
		},
	})

	goStruct1 := allGoStructs[1]
	suite.Equal(goStruct1.Package.Path, entitiesConfig.Package.Path)
	suite.Equal(goStruct1.Package.Name, entitiesConfig.Package.Name)
	suite.Equal(goStruct1.Name, "BasicIDPrimary")

	structFields1 := goStruct1.Fields
	suite.Len(structFields1, 1)

	field10 := structFields1[0]
	suite.Equal(field10.Name, "ID")
	suite.Equal(field10.Type, godef.GoTypeUint)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToGoStructs_Related_ForMany() {
	modelsConfig := cfg.MorpheModelsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
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
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/entities",
			Name: "entities",
		},
		ReceiverName: "e",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig:     modelsConfig,
			MorpheStructuresConfig: structuresConfig,
			MorpheEnumsConfig:      enumsConfig,
			MorpheEntitiesConfig:   entitiesConfig,
		},
		EntityHooks: hook.CompileMorpheEntity{},
	}

	entity0 := yaml.Entity{
		Name: "Basic",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Basic.ID",
			},
			"String": {
				Type: "Basic.String",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"BasicParent": {
				Type: "ForMany",
			},
		},
	}

	r := registry.NewRegistry()

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"String": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"BasicParent": {
				Type: "ForMany",
			},
		},
	}
	r.SetModel("Basic", model0)

	model1 := yaml.Model{
		Name: "BasicParent",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"String": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Basic": {
				Type: "HasMany",
			},
		},
	}
	r.SetModel("BasicParent", model1)

	entity1 := yaml.Entity{
		Name: "BasicParent",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "BasicParent.ID",
			},
			"String": {
				Type: "BasicParent.String",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"Basic": {
				Type: "HasMany",
			},
		},
	}
	r.SetEntity("BasicParent", entity1)

	allGoStructs, goStructErr := compile.MorpheEntityToGoStructs(config.EntityHooks, config.MorpheConfig, r, entity0)

	suite.Nil(goStructErr)
	suite.Len(allGoStructs, 2)

	goStruct0 := allGoStructs[0]
	suite.Equal(goStruct0.Package.Path, entitiesConfig.Package.Path)
	suite.Equal(goStruct0.Package.Name, entitiesConfig.Package.Name)
	suite.Equal(goStruct0.Name, "Basic")

	structFields0 := goStruct0.Fields
	suite.Len(structFields0, 4)

	field00 := structFields0[0]
	suite.Equal(field00.Name, "ID")
	suite.Equal(field00.Type, godef.GoTypeUint)

	field01 := structFields0[1]
	suite.Equal(field01.Name, "String")
	suite.Equal(field01.Type, godef.GoTypeString)

	field02 := structFields0[2]
	suite.Equal(field02.Name, "BasicParentIDs")
	suite.Equal(field02.Type, godef.GoTypeArray{
		IsSlice:   true,
		ValueType: godef.GoTypeUint,
	})

	field03 := structFields0[3]
	suite.Equal(field03.Name, "BasicParents")
	suite.Equal(field03.Type, godef.GoTypeArray{
		IsSlice: true,
		ValueType: godef.GoTypeStruct{
			PackagePath: "",
			Name:        "BasicParent",
		},
	})

	goStruct1 := allGoStructs[1]
	suite.Equal(goStruct1.Package.Path, entitiesConfig.Package.Path)
	suite.Equal(goStruct1.Package.Name, entitiesConfig.Package.Name)
	suite.Equal(goStruct1.Name, "BasicIDPrimary")

	structFields1 := goStruct1.Fields
	suite.Len(structFields1, 1)

	field10 := structFields1[0]
	suite.Equal(field10.Name, "ID")
	suite.Equal(field10.Type, godef.GoTypeUint)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToGoStructs_Related_HasOne() {
	modelsConfig := cfg.MorpheModelsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
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
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/entities",
			Name: "entities",
		},
		ReceiverName: "e",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig:     modelsConfig,
			MorpheStructuresConfig: structuresConfig,
			MorpheEnumsConfig:      enumsConfig,
			MorpheEntitiesConfig:   entitiesConfig,
		},
		EntityHooks: hook.CompileMorpheEntity{},
	}

	entity0 := yaml.Entity{
		Name: "BasicParent",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "BasicParent.ID",
			},
			"String": {
				Type: "BasicParent.String",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"Basic": {
				Type: "HasOne",
			},
		},
	}

	r := registry.NewRegistry()

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"String": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"BasicParent": {
				Type: "ForOne",
			},
		},
	}
	r.SetModel("Basic", model0)

	model1 := yaml.Model{
		Name: "BasicParent",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"String": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Basic": {
				Type: "HasOne",
			},
		},
	}
	r.SetModel("BasicParent", model1)

	entity1 := yaml.Entity{
		Name: "Basic",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Basic.ID",
			},
			"String": {
				Type: "Basic.String",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"BasicParent": {
				Type: "ForOne",
			},
		},
	}
	r.SetEntity("Basic", entity1)

	allGoStructs, goStructErr := compile.MorpheEntityToGoStructs(config.EntityHooks, config.MorpheConfig, r, entity0)

	suite.Nil(goStructErr)
	suite.Len(allGoStructs, 2)

	goStruct0 := allGoStructs[0]
	suite.Equal(goStruct0.Package.Path, entitiesConfig.Package.Path)
	suite.Equal(goStruct0.Package.Name, entitiesConfig.Package.Name)
	suite.Equal(goStruct0.Name, "BasicParent")

	structFields0 := goStruct0.Fields
	suite.Len(structFields0, 4)

	field00 := structFields0[0]
	suite.Equal(field00.Name, "ID")
	suite.Equal(field00.Type, godef.GoTypeUint)

	field01 := structFields0[1]
	suite.Equal(field01.Name, "String")
	suite.Equal(field01.Type, godef.GoTypeString)

	field02 := structFields0[2]
	suite.Equal(field02.Name, "BasicID")
	suite.Equal(field02.Type, godef.GoTypePointer{
		ValueType: godef.GoTypeUint,
	})

	field03 := structFields0[3]
	suite.Equal(field03.Name, "Basic")
	suite.Equal(field03.Type, godef.GoTypePointer{
		ValueType: godef.GoTypeStruct{
			PackagePath: "",
			Name:        "Basic",
		},
	})

	goStruct1 := allGoStructs[1]
	suite.Equal(goStruct1.Package.Path, entitiesConfig.Package.Path)
	suite.Equal(goStruct1.Package.Name, entitiesConfig.Package.Name)
	suite.Equal(goStruct1.Name, "BasicParentIDPrimary")

	structFields1 := goStruct1.Fields
	suite.Len(structFields1, 1)

	field10 := structFields1[0]
	suite.Equal(field10.Name, "ID")
	suite.Equal(field10.Type, godef.GoTypeUint)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToGoStructs_Related_HasMany() {
	modelsConfig := cfg.MorpheModelsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
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
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/entities",
			Name: "entities",
		},
		ReceiverName: "e",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig:     modelsConfig,
			MorpheStructuresConfig: structuresConfig,
			MorpheEnumsConfig:      enumsConfig,
			MorpheEntitiesConfig:   entitiesConfig,
		},
		EntityHooks: hook.CompileMorpheEntity{},
	}

	entity0 := yaml.Entity{
		Name: "BasicParent",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "BasicParent.ID",
			},
			"String": {
				Type: "BasicParent.String",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"Basic": {
				Type: "HasMany",
			},
		},
	}

	r := registry.NewRegistry()

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"String": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"BasicParent": {
				Type: "ForOne",
			},
		},
	}
	r.SetModel("Basic", model0)

	model1 := yaml.Model{
		Name: "BasicParent",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"String": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Basic": {
				Type: "HasMany",
			},
		},
	}
	r.SetModel("BasicParent", model1)

	entity1 := yaml.Entity{
		Name: "Basic",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Basic.ID",
			},
			"String": {
				Type: "Basic.String",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"BasicParent": {
				Type: "ForOne",
			},
		},
	}
	r.SetEntity("Basic", entity1)

	allGoStructs, goStructErr := compile.MorpheEntityToGoStructs(config.EntityHooks, config.MorpheConfig, r, entity0)

	suite.Nil(goStructErr)
	suite.Len(allGoStructs, 2)

	goStruct0 := allGoStructs[0]
	suite.Equal(goStruct0.Package.Path, entitiesConfig.Package.Path)
	suite.Equal(goStruct0.Package.Name, entitiesConfig.Package.Name)
	suite.Equal(goStruct0.Name, "BasicParent")

	structFields0 := goStruct0.Fields
	suite.Len(structFields0, 4)

	field00 := structFields0[0]
	suite.Equal(field00.Name, "ID")
	suite.Equal(field00.Type, godef.GoTypeUint)

	field01 := structFields0[1]
	suite.Equal(field01.Name, "String")
	suite.Equal(field01.Type, godef.GoTypeString)

	field02 := structFields0[2]
	suite.Equal(field02.Name, "BasicIDs")
	suite.Equal(field02.Type, godef.GoTypeArray{
		IsSlice:   true,
		ValueType: godef.GoTypeUint,
	})

	field03 := structFields0[3]
	suite.Equal(field03.Name, "Basics")
	suite.Equal(field03.Type, godef.GoTypeArray{
		IsSlice: true,
		ValueType: godef.GoTypeStruct{
			PackagePath: "",
			Name:        "Basic",
		},
	})

	goStruct1 := allGoStructs[1]
	suite.Equal(goStruct1.Package.Path, entitiesConfig.Package.Path)
	suite.Equal(goStruct1.Package.Name, entitiesConfig.Package.Name)
	suite.Equal(goStruct1.Name, "BasicParentIDPrimary")

	structFields1 := goStruct1.Fields
	suite.Len(structFields1, 1)

	field10 := structFields1[0]
	suite.Equal(field10.Name, "ID")
	suite.Equal(field10.Type, godef.GoTypeUint)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToGoStructs_StartHook_Successful() {
	var hookCalled = false
	entityHooks := hook.CompileMorpheEntity{
		OnCompileMorpheEntityStart: func(config cfg.MorpheConfig, entity yaml.Entity) (cfg.MorpheConfig, yaml.Entity, error) {
			hookCalled = true
			return config, entity, nil
		},
	}

	modelsConfig := cfg.MorpheModelsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
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
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/entities",
			Name: "entities",
		},
		ReceiverName: "e",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig:     modelsConfig,
			MorpheStructuresConfig: structuresConfig,
			MorpheEnumsConfig:      enumsConfig,
			MorpheEntitiesConfig:   entitiesConfig,
		},
		EntityHooks: entityHooks,
	}

	entity0 := yaml.Entity{
		Name: "Basic",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Basic.ID",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("Basic", model0)

	allGoStructs, goStructErr := compile.MorpheEntityToGoStructs(config.EntityHooks, config.MorpheConfig, r, entity0)

	suite.True(hookCalled)
	suite.Nil(goStructErr)
	suite.Len(allGoStructs, 2)

	goStruct0 := allGoStructs[0]
	suite.Equal(goStruct0.Package.Path, entitiesConfig.Package.Path)
	suite.Equal(goStruct0.Package.Name, entitiesConfig.Package.Name)
	suite.Equal(goStruct0.Name, "Basic")

	structFields0 := goStruct0.Fields
	suite.Len(structFields0, 1)

	field00 := structFields0[0]
	suite.Equal(field00.Name, "ID")
	suite.Equal(field00.Type, godef.GoTypeUint)

	goStruct1 := allGoStructs[1]
	suite.Equal(goStruct1.Package.Path, entitiesConfig.Package.Path)
	suite.Equal(goStruct1.Package.Name, entitiesConfig.Package.Name)
	suite.Equal(goStruct1.Name, "BasicIDPrimary")

	structFields1 := goStruct1.Fields
	suite.Len(structFields1, 1)

	field10 := structFields1[0]
	suite.Equal(field10.Name, "ID")
	suite.Equal(field10.Type, godef.GoTypeUint)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToGoStructs_StartHook_Failure() {
	var featureFlag = "otherName"
	entityHooks := hook.CompileMorpheEntity{
		OnCompileMorpheEntityStart: func(config cfg.MorpheConfig, entity yaml.Entity) (cfg.MorpheConfig, yaml.Entity, error) {
			if featureFlag != "otherName" {
				return config, entity, nil
			}
			return config, entity, fmt.Errorf("compile entity start hook error")
		},
	}

	modelsConfig := cfg.MorpheModelsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
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
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/entities",
			Name: "entities",
		},
		ReceiverName: "e",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig:     modelsConfig,
			MorpheStructuresConfig: structuresConfig,
			MorpheEnumsConfig:      enumsConfig,
			MorpheEntitiesConfig:   entitiesConfig,
		},
		EntityHooks: entityHooks,
	}

	entity0 := yaml.Entity{
		Name: "Basic",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Basic.ID",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("Basic", model0)

	allGoStructs, goStructErr := compile.MorpheEntityToGoStructs(config.EntityHooks, config.MorpheConfig, r, entity0)

	suite.NotNil(goStructErr)
	suite.ErrorContains(goStructErr, "compile entity start hook error")
	suite.Nil(allGoStructs)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToGoStructs_SuccessHook_Successful() {
	var featureFlag = "otherName"
	entityHooks := hook.CompileMorpheEntity{
		OnCompileMorpheEntitySuccess: func(allEntityStructs []*godef.Struct) ([]*godef.Struct, error) {
			if featureFlag != "otherName" {
				return allEntityStructs, nil
			}
			for structIdx := range allEntityStructs {
				allEntityStructs[structIdx].Name = allEntityStructs[structIdx].Name + "CHANGED"
			}
			return allEntityStructs, nil
		},
	}

	modelsConfig := cfg.MorpheModelsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
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
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/entities",
			Name: "entities",
		},
		ReceiverName: "e",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig:     modelsConfig,
			MorpheStructuresConfig: structuresConfig,
			MorpheEnumsConfig:      enumsConfig,
			MorpheEntitiesConfig:   entitiesConfig,
		},
		EntityHooks: entityHooks,
	}

	entity0 := yaml.Entity{
		Name: "Basic",
		Fields: map[string]yaml.EntityField{
			"UUID": {
				Type: "Basic.UUID",
				Attributes: []string{
					"immutable",
				},
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()

	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
				},
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("Basic", model0)

	allGoStructs, goStructErr := compile.MorpheEntityToGoStructs(config.EntityHooks, config.MorpheConfig, r, entity0)

	suite.Nil(goStructErr)
	suite.Len(allGoStructs, 2)

	goStruct0 := allGoStructs[0]
	suite.Equal(goStruct0.Package.Path, entitiesConfig.Package.Path)
	suite.Equal(goStruct0.Package.Name, entitiesConfig.Package.Name)
	suite.Equal(goStruct0.Name, "BasicCHANGED")

	structFields0 := goStruct0.Fields
	suite.Len(structFields0, 1)

	field00 := structFields0[0]
	suite.Equal(field00.Name, "UUID")
	suite.Equal(field00.Type, godef.GoTypeString)

	goStruct1 := allGoStructs[1]
	suite.Equal(goStruct1.Package.Path, entitiesConfig.Package.Path)
	suite.Equal(goStruct1.Package.Name, entitiesConfig.Package.Name)
	suite.Equal(goStruct1.Name, "BasicIDPrimaryCHANGED")

	structFields1 := goStruct1.Fields
	suite.Len(structFields1, 1)

	field10 := structFields1[0]
	suite.Equal(field10.Name, "UUID")
	suite.Equal(field10.Type, godef.GoTypeString)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToGoStructs_SuccessHook_Failure() {
	var featureFlag = "otherName"
	entityHooks := hook.CompileMorpheEntity{
		OnCompileMorpheEntitySuccess: func(allEntityStructs []*godef.Struct) ([]*godef.Struct, error) {
			if featureFlag != "otherName" {
				return allEntityStructs, nil
			}
			return nil, fmt.Errorf("compile entity success hook error")
		},
	}

	modelsConfig := cfg.MorpheModelsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
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
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/entities",
			Name: "entities",
		},
		ReceiverName: "e",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig:     modelsConfig,
			MorpheStructuresConfig: structuresConfig,
			MorpheEnumsConfig:      enumsConfig,
			MorpheEntitiesConfig:   entitiesConfig,
		},
		EntityHooks: entityHooks,
	}

	entity0 := yaml.Entity{
		Name: "User",
		Fields: map[string]yaml.EntityField{
			"UUID": {
				Type: "User.UUID",
				Attributes: []string{
					"immutable",
				},
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()

	userModel := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
				},
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	r.SetModel("User", userModel)

	allGoStructs, goStructErr := compile.MorpheEntityToGoStructs(config.EntityHooks, config.MorpheConfig, r, entity0)

	suite.NotNil(goStructErr)
	suite.ErrorContains(goStructErr, "compile entity success hook error")
	suite.Nil(allGoStructs)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToGoStructs_FailureHook_UnknownRootModel() {
	entityHooks := hook.CompileMorpheEntity{
		OnCompileMorpheEntityFailure: func(config cfg.MorpheConfig, entity yaml.Entity, compileFailure error) error {
			return fmt.Errorf("Entity %s: %w", entity.Name, compileFailure)
		},
	}

	modelsConfig := cfg.MorpheModelsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
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
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/entities",
			Name: "entities",
		},
		ReceiverName: "e",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig:     modelsConfig,
			MorpheStructuresConfig: structuresConfig,
			MorpheEnumsConfig:      enumsConfig,
			MorpheEntitiesConfig:   entitiesConfig,
		},
		EntityHooks: entityHooks,
	}

	entity0 := yaml.Entity{
		Name: "User",
		Fields: map[string]yaml.EntityField{
			"UUID": {
				Type: "NonExistentModel.UUID",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"UUID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	r := registry.NewRegistry()

	allGoStructs, goStructErr := compile.MorpheEntityToGoStructs(config.EntityHooks, config.MorpheConfig, r, entity0)

	suite.NotNil(goStructErr)
	suite.ErrorContains(goStructErr, "Entity User: morphe entity User field UUID references unknown root model: NonExistentModel")
	suite.Nil(allGoStructs)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToGoStructs_Related_ForOnePoly() {
	modelsConfig := cfg.MorpheModelsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	enumConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/enums",
			Name: "enums",
		},
	}
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/entities",
			Name: "entities",
		},
		ReceiverName: "e",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig:     modelsConfig,
			MorpheStructuresConfig: structuresConfig,
			MorpheEnumsConfig:      enumConfig,
			MorpheEntitiesConfig:   entitiesConfig,
		},
		EntityHooks: hook.CompileMorpheEntity{},
	}

	// Comment entity for a model that has ForOnePoly relationships
	commentEntity := yaml.Entity{
		Name: "Comment",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Comment.ID",
			},
			"Content": {
				Type: "Comment.Content",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	// Comment model with ForOnePoly relationship
	commentModel := yaml.Model{
		Name: "Comment",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"Content": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Post", "Article"},
			},
		},
	}

	postModel := yaml.Model{
		Name: "Post",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeUUID,
			},
			"Title": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
	}

	articleModel := yaml.Model{
		Name: "Article",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeUUID,
			},
			"Content": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
	}

	r := registry.NewRegistry()
	r.SetModel("Comment", commentModel)
	r.SetModel("Post", postModel)
	r.SetModel("Article", articleModel)
	r.SetEntity("Comment", commentEntity)

	allGoStructs, goStructErr := compile.MorpheEntityToGoStructs(config.EntityHooks, config.MorpheConfig, r, commentEntity)

	suite.Nil(goStructErr)
	suite.Len(allGoStructs, 2)

	goStruct0 := allGoStructs[0]
	suite.Equal(goStruct0.Name, "Comment")

	structFields0 := goStruct0.Fields
	suite.Len(structFields0, 2)

	// Entity fields only include the declared fields, not the polymorphic columns
	field00 := structFields0[0]
	suite.Equal(field00.Name, "Content")
	suite.Equal(field00.Type, godef.GoTypeString)

	field01 := structFields0[1]
	suite.Equal(field01.Name, "ID")
	suite.Equal(field01.Type, godef.GoTypeUint)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToGoStructs_Related_ForManyPoly() {
	modelsConfig := cfg.MorpheModelsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	enumConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/enums",
			Name: "enums",
		},
	}
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/entities",
			Name: "entities",
		},
		ReceiverName: "e",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig:     modelsConfig,
			MorpheStructuresConfig: structuresConfig,
			MorpheEnumsConfig:      enumConfig,
			MorpheEntitiesConfig:   entitiesConfig,
		},
		EntityHooks: hook.CompileMorpheEntity{},
	}

	// Tag entity for a model that has ForManyPoly relationships
	tagEntity := yaml.Entity{
		Name: "Tag",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Tag.ID",
			},
			"Name": {
				Type: "Tag.Name",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	// Tag model with ForManyPoly relationship
	tagModel := yaml.Model{
		Name: "Tag",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
			"Name": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Taggable": {
				Type: "ForManyPoly",
				For:  []string{"Post", "Product"},
			},
		},
	}

	postModel := yaml.Model{
		Name: "Post",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeUUID,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
	}

	productModel := yaml.Model{
		Name: "Product",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
	}

	r := registry.NewRegistry()
	r.SetModel("Tag", tagModel)
	r.SetModel("Post", postModel)
	r.SetModel("Product", productModel)
	r.SetEntity("Tag", tagEntity)

	allGoStructs, goStructErr := compile.MorpheEntityToGoStructs(config.EntityHooks, config.MorpheConfig, r, tagEntity)

	suite.Nil(goStructErr)
	suite.Len(allGoStructs, 2)

	goStruct0 := allGoStructs[0]
	suite.Equal(goStruct0.Name, "Tag")

	structFields0 := goStruct0.Fields
	suite.Len(structFields0, 2)

	// Entity fields only include the declared fields, not the polymorphic columns
	field00 := structFields0[0]
	suite.Equal(field00.Name, "ID")
	suite.Equal(field00.Type, godef.GoTypeUint)

	field01 := structFields0[1]
	suite.Equal(field01.Name, "Name")
	suite.Equal(field01.Type, godef.GoTypeString)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToGoStructs_Related_HasOnePoly() {
	modelsConfig := cfg.MorpheModelsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	enumConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/enums",
			Name: "enums",
		},
	}
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/entities",
			Name: "entities",
		},
		ReceiverName: "e",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig:     modelsConfig,
			MorpheStructuresConfig: structuresConfig,
			MorpheEnumsConfig:      enumConfig,
			MorpheEntitiesConfig:   entitiesConfig,
		},
		EntityHooks: hook.CompileMorpheEntity{},
	}

	// Post entity that has a polymorphic Comment relationship
	postEntity := yaml.Entity{
		Name: "Post",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Post.ID",
			},
			"Title": {
				Type: "Post.Title",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"Note": {
				Type:    "HasOnePoly",
				Through: "Commentable",
				Aliased: "Comment",
			},
		},
	}

	// Comment model with the forward ForOnePoly relationship
	commentModel := yaml.Model{
		Name: "Comment",
		Fields: map[string]yaml.ModelField{
			"ID":      {Type: yaml.ModelFieldTypeAutoIncrement},
			"Content": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Post", "Article"},
			},
		},
	}

	// Comment entity
	commentEntity := yaml.Entity{
		Name: "Comment",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Comment.ID",
			},
			"Text": {
				Type: "Comment.Text",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Post", "Article"},
			},
		},
	}

	postModel := yaml.Model{
		Name: "Post",
		Fields: map[string]yaml.ModelField{
			"ID":    {Type: yaml.ModelFieldTypeUUID},
			"Title": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Note": {
				Type:    "HasOnePoly",
				Through: "Commentable",
				Aliased: "Comment",
			},
		},
	}

	r := registry.NewRegistry()
	r.SetModel("Comment", commentModel)
	r.SetModel("Post", postModel)
	r.SetEntity("Comment", commentEntity)
	r.SetEntity("Post", postEntity)

	allGoStructs, goStructErr := compile.MorpheEntityToGoStructs(config.EntityHooks, config.MorpheConfig, r, postEntity)

	suite.Nil(goStructErr)
	suite.Len(allGoStructs, 2)

	goStruct0 := allGoStructs[0]
	suite.Equal(goStruct0.Name, "Post")

	structFields0 := goStruct0.Fields
	suite.Len(structFields0, 4)

	// Regular fields
	field00 := structFields0[0]
	suite.Equal(field00.Name, "ID")
	suite.Equal(field00.Type, godef.GoTypeString)

	field01 := structFields0[1]
	suite.Equal(field01.Name, "Title")
	suite.Equal(field01.Type, godef.GoTypeString)

	// HasOnePoly generates regular relationship fields
	field02 := structFields0[2]
	suite.Equal(field02.Name, "NoteID")
	suite.Equal(field02.Type, godef.GoTypePointer{
		ValueType: godef.GoTypeUint,
	})

	field03 := structFields0[3]
	suite.Equal(field03.Name, "Note")
	suite.Equal(field03.Type, godef.GoTypePointer{
		ValueType: godef.GoTypeStruct{
			Name: "Comment",
		},
	})
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToGoStructs_Related_HasManyPoly() {
	modelsConfig := cfg.MorpheModelsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	enumConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/enums",
			Name: "enums",
		},
	}
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/entities",
			Name: "entities",
		},
		ReceiverName: "e",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig:     modelsConfig,
			MorpheStructuresConfig: structuresConfig,
			MorpheEnumsConfig:      enumConfig,
			MorpheEntitiesConfig:   entitiesConfig,
		},
		EntityHooks: hook.CompileMorpheEntity{},
	}

	// Post entity that has polymorphic Comments relationship
	postEntity := yaml.Entity{
		Name: "Post",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Post.ID",
			},
			"Title": {
				Type: "Post.Title",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"Note": {
				Type:    "HasManyPoly",
				Through: "Commentable",
				Aliased: "Comment",
			},
		},
	}

	// Comment model with the forward ForOnePoly relationship
	commentModel := yaml.Model{
		Name: "Comment",
		Fields: map[string]yaml.ModelField{
			"ID":      {Type: yaml.ModelFieldTypeAutoIncrement},
			"Content": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Post", "Article"},
			},
		},
	}

	// Comment entity
	commentEntity := yaml.Entity{
		Name: "Comment",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Comment.ID",
			},
			"Text": {
				Type: "Comment.Text",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Post", "Article"},
			},
		},
	}

	postModel := yaml.Model{
		Name: "Post",
		Fields: map[string]yaml.ModelField{
			"ID":    {Type: yaml.ModelFieldTypeUUID},
			"Title": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Note": {
				Type:    "HasManyPoly",
				Through: "Commentable",
				Aliased: "Comment",
			},
		},
	}

	r := registry.NewRegistry()
	r.SetModel("Comment", commentModel)
	r.SetModel("Post", postModel)
	r.SetEntity("Comment", commentEntity)
	r.SetEntity("Post", postEntity)

	allGoStructs, goStructErr := compile.MorpheEntityToGoStructs(config.EntityHooks, config.MorpheConfig, r, postEntity)

	suite.Nil(goStructErr)
	suite.Len(allGoStructs, 2)

	goStruct0 := allGoStructs[0]
	suite.Equal(goStruct0.Name, "Post")

	structFields0 := goStruct0.Fields
	suite.Len(structFields0, 4)

	// Regular fields
	field00 := structFields0[0]
	suite.Equal(field00.Name, "ID")
	suite.Equal(field00.Type, godef.GoTypeString)

	field01 := structFields0[1]
	suite.Equal(field01.Name, "Title")
	suite.Equal(field01.Type, godef.GoTypeString)

	// HasManyPoly generates regular relationship array fields
	field02 := structFields0[2]
	suite.Equal(field02.Name, "NoteIDs")
	suite.Equal(field02.Type, godef.GoTypeArray{
		IsSlice:   true,
		ValueType: godef.GoTypeUint,
	})

	field03 := structFields0[3]
	suite.Equal(field03.Name, "Notes")
	suite.Equal(field03.Type, godef.GoTypeArray{
		IsSlice: true,
		ValueType: godef.GoTypeStruct{
			Name: "Comment",
		},
	})
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToGoStructs_Mixed_Polymorphic_And_Regular() {
	modelsConfig := cfg.MorpheModelsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	enumConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/enums",
			Name: "enums",
		},
	}
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/entities",
			Name: "entities",
		},
		ReceiverName: "e",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig:     modelsConfig,
			MorpheStructuresConfig: structuresConfig,
			MorpheEnumsConfig:      enumConfig,
			MorpheEntitiesConfig:   entitiesConfig,
		},
		EntityHooks: hook.CompileMorpheEntity{},
	}

	// Comment entity with both polymorphic and regular relationships
	commentEntity := yaml.Entity{
		Name: "Comment",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Comment.ID",
			},
			"Content": {
				Type: "Comment.Content",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"User": {
				Type: "ForOne", // Regular relationship
			},
		},
	}

	// User model and entity
	userModel := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: yaml.ModelFieldTypeAutoIncrement},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	}

	userEntity := yaml.Entity{
		Name: "User",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "User.ID",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	// Comment model with mixed relationships
	commentModel := yaml.Model{
		Name: "Comment",
		Fields: map[string]yaml.ModelField{
			"ID":      {Type: yaml.ModelFieldTypeAutoIncrement},
			"Content": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Post", "Article"},
			},
			"User": {
				Type: "ForOne", // Regular relationship
			},
		},
	}

	postModel := yaml.Model{
		Name: "Post",
		Fields: map[string]yaml.ModelField{
			"id": {Type: yaml.ModelFieldTypeUUID},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
	}

	articleModel := yaml.Model{
		Name: "Article",
		Fields: map[string]yaml.ModelField{
			"id": {Type: yaml.ModelFieldTypeUUID},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"id"}},
		},
	}

	r := registry.NewRegistry()
	r.SetModel("Comment", commentModel)
	r.SetModel("User", userModel)
	r.SetModel("Post", postModel)
	r.SetModel("Article", articleModel)
	r.SetEntity("Comment", commentEntity)
	r.SetEntity("User", userEntity)

	allGoStructs, goStructErr := compile.MorpheEntityToGoStructs(config.EntityHooks, config.MorpheConfig, r, commentEntity)

	suite.Nil(goStructErr)
	suite.Len(allGoStructs, 2)

	goStruct0 := allGoStructs[0]
	suite.Equal(goStruct0.Name, "Comment")

	structFields0 := goStruct0.Fields
	suite.Len(structFields0, 4)

	// Regular entity fields
	field00 := structFields0[0]
	suite.Equal(field00.Name, "Content")
	suite.Equal(field00.Type, godef.GoTypeString)

	field01 := structFields0[1]
	suite.Equal(field01.Name, "ID")
	suite.Equal(field01.Type, godef.GoTypeUint)

	// Regular relationship fields
	field02 := structFields0[2]
	suite.Equal(field02.Name, "UserID")
	suite.Equal(field02.Type, godef.GoTypePointer{
		ValueType: godef.GoTypeUint,
	})

	field03 := structFields0[3]
	suite.Equal(field03.Name, "User")
	suite.Equal(field03.Type, godef.GoTypePointer{
		ValueType: godef.GoTypeStruct{
			Name: "User",
		},
	})
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToGoStructs_Related_ForOnePoly_Aliased() {
	modelsConfig := cfg.MorpheModelsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	enumConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/enums",
			Name: "enums",
		},
	}
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/entities",
			Name: "entities",
		},
		ReceiverName: "e",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig:     modelsConfig,
			MorpheStructuresConfig: structuresConfig,
			MorpheEnumsConfig:      enumConfig,
			MorpheEntitiesConfig:   entitiesConfig,
		},
		EntityHooks: hook.CompileMorpheEntity{},
	}

	// Comment entity
	commentEntity := yaml.Entity{
		Name: "Comment",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Comment.ID",
			},
			"Text": {
				Type: "Comment.Text",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}

	// Document model
	documentModel := yaml.Model{
		Name: "Document",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: yaml.ModelFieldTypeAutoIncrement},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	}

	// Video model
	videoModel := yaml.Model{
		Name: "Video",
		Fields: map[string]yaml.ModelField{
			"ID": {Type: yaml.ModelFieldTypeAutoIncrement},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
	}

	// Comment model with ForOnePoly using alias
	commentModel := yaml.Model{
		Name: "Comment",
		Fields: map[string]yaml.ModelField{
			"ID":   {Type: yaml.ModelFieldTypeAutoIncrement},
			"Text": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"CommentableResource": {
				Type:    "ForOnePoly",
				For:     []string{"Document", "Video"},
				Aliased: "Resource", // Alias that doesn't map to any model
			},
		},
	}

	r := registry.NewRegistry()
	r.SetModel("Document", documentModel)
	r.SetModel("Video", videoModel)
	r.SetModel("Comment", commentModel)
	r.SetEntity("Comment", commentEntity)

	allGoStructs, goStructErr := compile.MorpheEntityToGoStructs(config.EntityHooks, config.MorpheConfig, r, commentEntity)

	suite.Nil(goStructErr)
	suite.Len(allGoStructs, 2)

	goStruct0 := allGoStructs[0]
	suite.Equal(goStruct0.Name, "Comment")

	structFields0 := goStruct0.Fields
	suite.Len(structFields0, 2)

	// Entity only includes declared fields
	suite.Equal("ID", structFields0[0].Name)
	suite.Equal(godef.GoTypeUint, structFields0[0].Type)

	suite.Equal("Text", structFields0[1].Name)
	suite.Equal(godef.GoTypeString, structFields0[1].Type)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToGoStructs_Related_HasOnePoly_Aliased() {
	modelsConfig := cfg.MorpheModelsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	enumConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/enums",
			Name: "enums",
		},
	}
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/entities",
			Name: "entities",
		},
		ReceiverName: "e",
	}
	config := compile.MorpheCompileConfig{
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig:     modelsConfig,
			MorpheStructuresConfig: structuresConfig,
			MorpheEnumsConfig:      enumConfig,
			MorpheEntitiesConfig:   entitiesConfig,
		},
		EntityHooks: hook.CompileMorpheEntity{},
	}

	// Comment model with ForOnePoly
	commentModel := yaml.Model{
		Name: "Comment",
		Fields: map[string]yaml.ModelField{
			"ID":   {Type: yaml.ModelFieldTypeAutoIncrement},
			"Text": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Post", "Task"},
			},
		},
	}

	// Post model with HasOnePoly using different field name
	postModel := yaml.Model{
		Name: "Post",
		Fields: map[string]yaml.ModelField{
			"ID":    {Type: yaml.ModelFieldTypeAutoIncrement},
			"Title": {Type: yaml.ModelFieldTypeString},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"FeaturedComment": {
				Type:    "HasOnePoly",
				Through: "Commentable",
				Aliased: "Comment",
			},
		},
	}

	// Post entity with the aliased relationship
	postEntity := yaml.Entity{
		Name: "Post",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Post.ID",
			},
			"Title": {
				Type: "Post.Title",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"FeaturedComment": {
				Type:    "HasOnePoly",
				Through: "Commentable",
				Aliased: "Comment",
			},
		},
	}

	// Comment entity
	commentEntity := yaml.Entity{
		Name: "Comment",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "Comment.ID",
			},
			"Text": {
				Type: "Comment.Text",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Post", "Task"},
			},
		},
	}

	r := registry.NewRegistry()
	r.SetModel("Comment", commentModel)
	r.SetModel("Post", postModel)
	r.SetEntity("Post", postEntity)
	r.SetEntity("Comment", commentEntity)

	allGoStructs, goStructErr := compile.MorpheEntityToGoStructs(config.EntityHooks, config.MorpheConfig, r, postEntity)

	suite.Nil(goStructErr)
	suite.Len(allGoStructs, 2)

	goStruct0 := allGoStructs[0]
	suite.Equal(goStruct0.Name, "Post")

	structFields0 := goStruct0.Fields
	suite.Len(structFields0, 4)

	// Regular fields
	suite.Equal("ID", structFields0[0].Name)
	suite.Equal(godef.GoTypeUint, structFields0[0].Type)

	suite.Equal("Title", structFields0[1].Name)
	suite.Equal(godef.GoTypeString, structFields0[1].Type)

	// Aliased polymorphic relationship uses the field name
	suite.Equal("FeaturedCommentID", structFields0[2].Name)
	suite.Equal(godef.GoTypePointer{
		ValueType: godef.GoTypeUint,
	}, structFields0[2].Type)

	suite.Equal("FeaturedComment", structFields0[3].Name)
	suite.Equal(godef.GoTypePointer{
		ValueType: godef.GoTypeStruct{
			Name: "Comment",
		},
	}, structFields0[3].Type)
}
