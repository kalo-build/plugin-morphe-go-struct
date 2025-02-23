package compile_test

import (
	"fmt"
	"testing"

	"github.com/kaloseia/go/pkg/godef"
	"github.com/kaloseia/morphe-go/pkg/registry"
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/cfg"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/hook"
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
			Path: "github.com/kaloseia/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/enums",
			Name: "enums",
		},
	}
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/entities",
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
	suite.Len(tsField07.Tags, 0)

	goStruct1 := allGoStructs[1]
	suite.Equal(goStruct1.Package.Path, entitiesConfig.Package.Path)
	suite.Equal(goStruct1.Package.Name, entitiesConfig.Package.Name)
	suite.Equal(goStruct1.Name, "UserIDPrimary")

	structFields1 := goStruct1.Fields
	suite.Len(structFields1, 1)

	structFields10 := structFields1[0]
	suite.Equal(structFields10.Name, "UUID")
	suite.Equal(structFields10.Type, godef.GoTypeString)
}

func (suite *CompileEntitiesTestSuite) TestMorpheEntityToGoStructs_NoEntityName() {
	modelsConfig := cfg.MorpheModelsConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/enums",
			Name: "enums",
		},
	}
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/entities",
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
			Path: "github.com/kaloseia/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/enums",
			Name: "enums",
		},
	}
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/entities",
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
			Path: "github.com/kaloseia/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/enums",
			Name: "enums",
		},
	}
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/entities",
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
			Path: "github.com/kaloseia/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/enums",
			Name: "enums",
		},
	}
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/entities",
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
	suite.Len(field01.Tags, 0)

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
			Path: "github.com/kaloseia/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/enums",
			Name: "enums",
		},
	}
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/entities",
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
			Path: "github.com/kaloseia/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/enums",
			Name: "enums",
		},
	}
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/entities",
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
			Path: "github.com/kaloseia/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/enums",
			Name: "enums",
		},
	}
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/entities",
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
			Path: "github.com/kaloseia/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/enums",
			Name: "enums",
		},
	}
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/entities",
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
			Path: "github.com/kaloseia/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/enums",
			Name: "enums",
		},
	}
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/entities",
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
			Path: "github.com/kaloseia/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/enums",
			Name: "enums",
		},
	}
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/entities",
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
			Path: "github.com/kaloseia/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/enums",
			Name: "enums",
		},
	}
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/entities",
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
			Path: "github.com/kaloseia/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/enums",
			Name: "enums",
		},
	}
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/entities",
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
			Path: "github.com/kaloseia/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/enums",
			Name: "enums",
		},
	}
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/entities",
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
			Path: "github.com/kaloseia/project/domain/models",
			Name: "models",
		},
		ReceiverName: "m",
	}
	structuresConfig := cfg.MorpheStructuresConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/structures",
			Name: "structures",
		},
		ReceiverName: "s",
	}
	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/enums",
			Name: "enums",
		},
	}
	entitiesConfig := cfg.MorpheEntitiesConfig{
		Package: godef.Package{
			Path: "github.com/kaloseia/project/domain/entities",
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
