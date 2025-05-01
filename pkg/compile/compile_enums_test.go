package compile_test

import (
	"fmt"
	"testing"

	"github.com/kalo-build/go/pkg/godef"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/compile"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/compile/hook"
	"github.com/stretchr/testify/suite"
)

type CompileEnumsTestSuite struct {
	suite.Suite
}

func TestCompileEnumsTestSuite(t *testing.T) {
	suite.Run(t, new(CompileEnumsTestSuite))
}

func (suite *CompileEnumsTestSuite) SetupTest() {
}

func (suite *CompileEnumsTestSuite) TearDownTest() {
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToGoEnum_String() {
	enumHooks := hook.CompileMorpheEnum{}

	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/enums",
			Name: "enums",
		},
	}

	enum0 := yaml.Enum{
		Name: "Color",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Red":   "rgb(255,0,0)",
			"Green": "rgb(0,255,0)",
			"Blue":  "rgb(0,0,255)",
		},
	}

	goEnum, goEnumErr := compile.MorpheEnumToGoEnum(enumHooks, enumsConfig, enum0)

	suite.Nil(goEnumErr)

	suite.Equal(goEnum.Package.Path, enumsConfig.Package.Path)
	suite.Equal(goEnum.Package.Name, enumsConfig.Package.Name)

	suite.Equal(goEnum.Name, "Color")
	suite.Equal(goEnum.Type, godef.GoTypeDerived{
		PackagePath: enumsConfig.Package.Path,
		Name:        "Color",
		BaseType:    godef.GoTypeString,
	})

	enumEntries := goEnum.Entries
	suite.Len(enumEntries, 3)

	enumEntry0 := enumEntries[0]
	suite.Equal(enumEntry0.Name, "ColorBlue")
	suite.Equal(enumEntry0.Value, "rgb(0,0,255)")

	enumEntry1 := enumEntries[1]
	suite.Equal(enumEntry1.Name, "ColorGreen")
	suite.Equal(enumEntry1.Value, "rgb(0,255,0)")

	enumEntry2 := enumEntries[2]
	suite.Equal(enumEntry2.Name, "ColorRed")
	suite.Equal(enumEntry2.Value, "rgb(255,0,0)")
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToGoEnum_Float() {
	enumHooks := hook.CompileMorpheEnum{}

	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/enums",
			Name: "enums",
		},
	}

	enum0 := yaml.Enum{
		Name: "Analytics",
		Type: yaml.EnumTypeFloat,
		Entries: map[string]any{
			"Pi":    3.141,
			"Euler": 2.718,
		},
	}

	goEnum, goEnumErr := compile.MorpheEnumToGoEnum(enumHooks, enumsConfig, enum0)

	suite.Nil(goEnumErr)

	suite.Equal(goEnum.Package.Path, enumsConfig.Package.Path)
	suite.Equal(goEnum.Package.Name, enumsConfig.Package.Name)

	suite.Equal(goEnum.Name, "Analytics")
	suite.Equal(goEnum.Type, godef.GoTypeDerived{
		PackagePath: enumsConfig.Package.Path,
		Name:        "Analytics",
		BaseType:    godef.GoTypeFloat,
	})

	enumEntries := goEnum.Entries
	suite.Len(enumEntries, 2)

	enumEntry0 := enumEntries[0]
	suite.Equal(enumEntry0.Name, "AnalyticsEuler")
	suite.Equal(enumEntry0.Value, 2.718)

	enumEntry1 := enumEntries[1]
	suite.Equal(enumEntry1.Name, "AnalyticsPi")
	suite.Equal(enumEntry1.Value, 3.141)
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToGoEnum_Integer() {
	enumHooks := hook.CompileMorpheEnum{}

	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/enums",
			Name: "enums",
		},
	}

	enum0 := yaml.Enum{
		Name: "Analytics",
		Type: yaml.EnumTypeInteger,
		Entries: map[string]any{
			"AnswerToLife":  42,
			"FineStructure": 317,
		},
	}

	goEnum, goEnumErr := compile.MorpheEnumToGoEnum(enumHooks, enumsConfig, enum0)

	suite.Nil(goEnumErr)

	suite.Equal(goEnum.Package.Path, enumsConfig.Package.Path)
	suite.Equal(goEnum.Package.Name, enumsConfig.Package.Name)

	suite.Equal(goEnum.Name, "Analytics")
	suite.Equal(goEnum.Type, godef.GoTypeDerived{
		PackagePath: enumsConfig.Package.Path,
		Name:        "Analytics",
		BaseType:    godef.GoTypeInt,
	})

	enumEntries := goEnum.Entries
	suite.Len(enumEntries, 2)

	enumEntry0 := enumEntries[0]
	suite.Equal(enumEntry0.Name, "AnalyticsAnswerToLife")
	suite.Equal(enumEntry0.Value, 42)

	enumEntry1 := enumEntries[1]
	suite.Equal(enumEntry1.Name, "AnalyticsFineStructure")
	suite.Equal(enumEntry1.Value, 317)
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToGoEnum_NoName() {
	enumHooks := hook.CompileMorpheEnum{}

	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/enums",
			Name: "enums",
		},
	}

	enum0 := yaml.Enum{
		Name: "",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Red":   "rgb(255,0,0)",
			"Green": "rgb(0,255,0)",
			"Blue":  "rgb(0,0,255)",
		},
	}

	goEnum, goEnumErr := compile.MorpheEnumToGoEnum(enumHooks, enumsConfig, enum0)

	suite.ErrorIs(goEnumErr, yaml.ErrNoMorpheEnumName)
	suite.Nil(goEnum)
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToGoEnum_NoType() {
	enumHooks := hook.CompileMorpheEnum{}

	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/enums",
			Name: "enums",
		},
	}

	enum0 := yaml.Enum{
		Name: "Color",
		Type: "",
		Entries: map[string]any{
			"Red":   "rgb(255,0,0)",
			"Green": "rgb(0,255,0)",
			"Blue":  "rgb(0,0,255)",
		},
	}

	goEnum, goEnumErr := compile.MorpheEnumToGoEnum(enumHooks, enumsConfig, enum0)

	suite.ErrorIs(goEnumErr, yaml.ErrNoMorpheEnumType)
	suite.Nil(goEnum)
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToGoEnum_NoEntries() {
	enumHooks := hook.CompileMorpheEnum{}

	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/enums",
			Name: "enums",
		},
	}

	enum0 := yaml.Enum{
		Name:    "Color",
		Type:    yaml.EnumTypeString,
		Entries: map[string]any{},
	}

	goEnum, goEnumErr := compile.MorpheEnumToGoEnum(enumHooks, enumsConfig, enum0)

	suite.ErrorIs(goEnumErr, yaml.ErrNoMorpheEnumEntries)
	suite.Nil(goEnum)
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToGoEnum_EntryTypeMismatch() {
	enumHooks := hook.CompileMorpheEnum{}

	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/enums",
			Name: "enums",
		},
	}

	enum0 := yaml.Enum{
		Name: "Color",
		Type: yaml.EnumTypeInteger,
		Entries: map[string]any{
			"Red":   "rgb(255,0,0)",
			"Green": "rgb(0,255,0)",
			"Blue":  "rgb(0,0,255)",
		},
	}

	goEnum, goEnumErr := compile.MorpheEnumToGoEnum(enumHooks, enumsConfig, enum0)

	suite.ErrorContains(goEnumErr, "enum entry 'Blue' value 'rgb(0,0,255)' with type 'string' does not match the enum type of 'Integer'")
	suite.Nil(goEnum)
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToGoEnum_StartHook_Successful() {
	var featureFlag = "otherName"
	enumHooks := hook.CompileMorpheEnum{
		OnCompileMorpheEnumStart: func(config cfg.MorpheEnumsConfig, enum yaml.Enum) (cfg.MorpheEnumsConfig, yaml.Enum, error) {
			if featureFlag != "otherName" {
				return config, enum, nil
			}
			enum.Name = enum.Name + "CHANGED"
			delete(enum.Entries, "Green")
			return config, enum, nil
		},
	}

	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/enums",
			Name: "enums",
		},
	}

	enum0 := yaml.Enum{
		Name: "Color",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Red":   "rgb(255,0,0)",
			"Green": "rgb(0,255,0)",
			"Blue":  "rgb(0,0,255)",
		},
	}

	goEnum, goEnumErr := compile.MorpheEnumToGoEnum(enumHooks, enumsConfig, enum0)

	suite.Nil(goEnumErr)

	suite.Equal(goEnum.Name, "ColorCHANGED")
	suite.Equal(goEnum.Type, godef.GoTypeDerived{
		PackagePath: enumsConfig.Package.Path,
		Name:        "ColorCHANGED",
		BaseType:    godef.GoTypeString,
	})

	enumEntries := goEnum.Entries
	suite.Len(enumEntries, 2)

	enumEntry0 := enumEntries[0]
	suite.Equal(enumEntry0.Name, "ColorCHANGEDBlue")
	suite.Equal(enumEntry0.Value, "rgb(0,0,255)")

	enumEntry1 := enumEntries[1]
	suite.Equal(enumEntry1.Name, "ColorCHANGEDRed")
	suite.Equal(enumEntry1.Value, "rgb(255,0,0)")
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToGoEnum_StartHook_Failure() {
	var featureFlag = "otherName"
	enumHooks := hook.CompileMorpheEnum{
		OnCompileMorpheEnumStart: func(config cfg.MorpheEnumsConfig, enum yaml.Enum) (cfg.MorpheEnumsConfig, yaml.Enum, error) {
			if featureFlag != "otherName" {
				return config, enum, nil
			}
			return config, enum, fmt.Errorf("compile enum start hook error")
		},
	}

	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/enums",
			Name: "enums",
		},
	}

	enum0 := yaml.Enum{
		Name: "Color",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Red":   "rgb(255,0,0)",
			"Green": "rgb(0,255,0)",
			"Blue":  "rgb(0,0,255)",
		},
	}

	goEnum, goEnumErr := compile.MorpheEnumToGoEnum(enumHooks, enumsConfig, enum0)

	suite.ErrorContains(goEnumErr, "compile enum start hook error")
	suite.Nil(goEnum)
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToGoEnum_SuccessHook_Successful() {
	var featureFlag = "otherName"
	enumHooks := hook.CompileMorpheEnum{
		OnCompileMorpheEnumSuccess: func(enum *godef.Enum) (*godef.Enum, error) {
			if featureFlag != "otherName" {
				return enum, nil
			}
			enum.Name = enum.Name + "CHANGED"
			enum.Type.Name = enum.Type.Name + "CHANGED"
			newEntries := []godef.EnumEntry{}
			for _, enumEntry := range enum.Entries {
				if enumEntry.Name == "ColorGreen" {
					continue
				}
				newEntries = append(newEntries, enumEntry)
			}
			enum.Entries = newEntries
			return enum, nil
		},
	}

	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/enums",
			Name: "enums",
		},
	}

	enum0 := yaml.Enum{
		Name: "Color",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Red":   "rgb(255,0,0)",
			"Green": "rgb(0,255,0)",
			"Blue":  "rgb(0,0,255)",
		},
	}

	goEnum, goEnumErr := compile.MorpheEnumToGoEnum(enumHooks, enumsConfig, enum0)

	suite.Nil(goEnumErr)

	suite.Equal(goEnum.Name, "ColorCHANGED")
	suite.Equal(goEnum.Type, godef.GoTypeDerived{
		PackagePath: enumsConfig.Package.Path,
		Name:        "ColorCHANGED",
		BaseType:    godef.GoTypeString,
	})

	enumEntries := goEnum.Entries
	suite.Len(enumEntries, 2)

	enumEntry0 := enumEntries[0]
	suite.Equal(enumEntry0.Name, "ColorBlue")
	suite.Equal(enumEntry0.Value, "rgb(0,0,255)")

	enumEntry1 := enumEntries[1]
	suite.Equal(enumEntry1.Name, "ColorRed")
	suite.Equal(enumEntry1.Value, "rgb(255,0,0)")
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToGoEnum_SuccessHook_Failure() {
	var featureFlag = "otherName"
	enumHooks := hook.CompileMorpheEnum{
		OnCompileMorpheEnumSuccess: func(enum *godef.Enum) (*godef.Enum, error) {
			if featureFlag != "otherName" {
				return enum, nil
			}
			return nil, fmt.Errorf("compile enum success hook error")
		},
	}

	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/enums",
			Name: "enums",
		},
	}

	enum0 := yaml.Enum{
		Name: "Color",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Red":   "rgb(255,0,0)",
			"Green": "rgb(0,255,0)",
			"Blue":  "rgb(0,0,255)",
		},
	}

	goEnum, goEnumErr := compile.MorpheEnumToGoEnum(enumHooks, enumsConfig, enum0)

	suite.ErrorContains(goEnumErr, "compile enum success hook error")
	suite.Nil(goEnum)
}

func (suite *CompileEnumsTestSuite) TestMorpheEnumToGoEnum_FailureHook() {
	var featureFlag = "otherName"
	var failureErr error
	enumHooks := hook.CompileMorpheEnum{
		OnCompileMorpheEnumFailure: func(config cfg.MorpheEnumsConfig, enum yaml.Enum, err error) error {
			if featureFlag != "otherName" {
				return err
			}
			failureErr = err
			return fmt.Errorf("compile enum failure hook error: %w", err)
		},
	}

	enumsConfig := cfg.MorpheEnumsConfig{
		Package: godef.Package{
			Path: "github.com/kalo-build/project/domain/enums",
			Name: "enums",
		},
	}

	enum0 := yaml.Enum{
		Name: "",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Red":   "rgb(255,0,0)",
			"Green": "rgb(0,255,0)",
			"Blue":  "rgb(0,0,255)",
		},
	}

	goEnum, goEnumErr := compile.MorpheEnumToGoEnum(enumHooks, enumsConfig, enum0)

	suite.ErrorIs(failureErr, yaml.ErrNoMorpheEnumName)
	suite.ErrorContains(goEnumErr, "compile enum failure hook error")
	suite.Nil(goEnum)
}
