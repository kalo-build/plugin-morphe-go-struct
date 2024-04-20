package compile_test

import (
	"testing"

	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/godef"
	"github.com/stretchr/testify/suite"
)

type CompileModelsTestSuite struct {
	suite.Suite
}

func TestCompileModelsTestSuite(t *testing.T) {
	suite.Run(t, new(CompileModelsTestSuite))
}

func (suite *CompileModelsTestSuite) SetupTest() {
}

func (suite *CompileModelsTestSuite) TearDownTest() {
}

func (suite *CompileModelsTestSuite) TestMorpheModelToGoStruct() {
	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
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
			"Protected": {
				Type: yaml.ModelFieldTypeProtected,
			},
			"Sealed": {
				Type: yaml.ModelFieldTypeSealed,
			},
			"String": {
				Type: yaml.ModelFieldTypeString,
			},
			"Time": {
				Type: yaml.ModelFieldTypeTime,
			},
			"UUID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"immutable",
				},
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{
					"UUID",
				},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	const packageName = "models"

	goStruct, structErr := compile.MorpheModelToGoStruct(packageName, model0)

	suite.Nil(structErr)

	suite.Equal(goStruct.Package, packageName)

	structImports0 := goStruct.Imports
	suite.Len(structImports0, 1)

	structImports00 := structImports0[0]
	suite.Equal(structImports00, "time")

	suite.Equal(goStruct.Name, "Basic")

	structFields0 := goStruct.Fields
	suite.Len(structFields0, 10)

	structFields00 := structFields0[0]
	suite.Equal(structFields00.Name, "AutoIncrement")
	suite.Equal(structFields00.Type, godef.GoTypeUint)
	suite.Len(structFields00.Tags, 0)

	structFields01 := structFields0[1]
	suite.Equal(structFields01.Name, "Boolean")
	suite.Equal(structFields01.Type, godef.GoTypeBool)
	suite.Len(structFields01.Tags, 0)

	structFields02 := structFields0[2]
	suite.Equal(structFields02.Name, "Date")
	suite.Equal(structFields02.Type, godef.GoTypeTime)
	suite.Len(structFields02.Tags, 0)

	structFields03 := structFields0[3]
	suite.Equal(structFields03.Name, "Float")
	suite.Equal(structFields03.Type, godef.GoTypeFloat)
	suite.Len(structFields03.Tags, 0)

	structFields04 := structFields0[4]
	suite.Equal(structFields04.Name, "Integer")
	suite.Equal(structFields04.Type, godef.GoTypeInt)
	suite.Len(structFields04.Tags, 0)

	structFields05 := structFields0[5]
	suite.Equal(structFields05.Name, "Protected")
	suite.Equal(structFields05.Type, godef.GoTypeString)
	suite.Len(structFields05.Tags, 0)

	structFields06 := structFields0[6]
	suite.Equal(structFields06.Name, "Sealed")
	suite.Equal(structFields06.Type, godef.GoTypeString)
	suite.Len(structFields06.Tags, 0)

	structFields07 := structFields0[7]
	suite.Equal(structFields07.Name, "String")
	suite.Equal(structFields07.Type, godef.GoTypeString)
	suite.Len(structFields07.Tags, 0)

	structFields08 := structFields0[8]
	suite.Equal(structFields08.Name, "Time")
	suite.Equal(structFields08.Type, godef.GoTypeTime)
	suite.Len(structFields08.Tags, 0)

	structFields09 := structFields0[9]
	suite.Equal(structFields09.Name, "UUID")
	suite.Equal(structFields09.Type, godef.GoTypeString)
	suite.Len(structFields09.Tags, 1)
	suite.Equal(structFields09.Tags[0], "immutable")
}

func (suite *CompileModelsTestSuite) TestMorpheModelToGoStruct_NoPackageName() {
	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"AutoIncrement": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{
					"AutoIncrement",
				},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	const packageName = ""

	goStruct, structErr := compile.MorpheModelToGoStruct(packageName, model0)

	suite.NotNil(structErr)
	suite.ErrorContains(structErr, "model package name cannot be empty")

	suite.Nil(goStruct)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToGoStruct_NoModelName() {
	model0 := yaml.Model{
		Name: "",
		Fields: map[string]yaml.ModelField{
			"AutoIncrement": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{
					"AutoIncrement",
				},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	const packageName = "models"

	goStruct, structErr := compile.MorpheModelToGoStruct(packageName, model0)

	suite.NotNil(structErr)
	suite.ErrorContains(structErr, "morphe model has no name")

	suite.Nil(goStruct)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToGoStruct_NoFields() {
	model0 := yaml.Model{
		Name:   "Basic",
		Fields: map[string]yaml.ModelField{},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{
					"AutoIncrement",
				},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	const packageName = "models"

	goStruct, structErr := compile.MorpheModelToGoStruct(packageName, model0)

	suite.NotNil(structErr)
	suite.ErrorContains(structErr, "morphe model has no fields")

	suite.Nil(goStruct)
}

func (suite *CompileModelsTestSuite) TestMorpheModelToGoStruct_NoIdentifiers() {
	model0 := yaml.Model{
		Name: "Basic",
		Fields: map[string]yaml.ModelField{
			"AutoIncrement": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{},
		Related:     map[string]yaml.ModelRelation{},
	}
	const packageName = "models"

	goStruct, structErr := compile.MorpheModelToGoStruct(packageName, model0)

	suite.NotNil(structErr)
	suite.ErrorContains(structErr, "morphe model has no identifiers")

	suite.Nil(goStruct)
}
