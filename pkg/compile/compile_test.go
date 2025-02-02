package compile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kaloseia/go-util/assertfile"
	"github.com/kaloseia/go/pkg/godef"
	rcfg "github.com/kaloseia/morphe-go/pkg/registry/cfg"
	"github.com/kaloseia/plugin-morphe-go-struct/internal/testutils"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/cfg"
)

type CompileTestSuite struct {
	assertfile.FileSuite

	TestDirPath            string
	TestGroundTruthDirPath string

	ModelsDirPath     string
	EnumsDirPath      string
	StructuresDirPath string
	EntitiesDirPath   string
}

func TestCompileTestSuite(t *testing.T) {
	suite.Run(t, new(CompileTestSuite))
}

func (suite *CompileTestSuite) SetupTest() {
	suite.TestDirPath = testutils.GetTestDirPath()
	suite.TestGroundTruthDirPath = filepath.Join(suite.TestDirPath, "ground-truth", "compile-minimal")

	suite.ModelsDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "models")
	suite.EnumsDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "enums")
	suite.StructuresDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "structures")
	suite.EntitiesDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "entities")
}

func (suite *CompileTestSuite) TearDownTest() {
	suite.TestDirPath = ""
}

func (suite *CompileTestSuite) TestMorpheToGo() {
	workingDirPath := suite.TestDirPath + "/working"
	suite.Nil(os.Mkdir(workingDirPath, 0644))
	defer os.RemoveAll(workingDirPath)

	config := compile.MorpheCompileConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryEnumsDirPath:      suite.EnumsDirPath,
			RegistryStructuresDirPath: suite.StructuresDirPath,
			RegistryModelsDirPath:     suite.ModelsDirPath,
			RegistryEntitiesDirPath:   suite.EntitiesDirPath,
		},
		MorpheConfig: cfg.MorpheConfig{
			MorpheModelsConfig: cfg.MorpheModelsConfig{
				Package: godef.Package{
					Path: "github.com/kaloseia/dummy/models",
					Name: "models",
				},
				ReceiverName: "m",
			},
			MorpheEnumsConfig: cfg.MorpheEnumsConfig{
				Package: godef.Package{
					Path: "github.com/kaloseia/dummy/enums",
					Name: "enums",
				},
			},
		},

		ModelWriter: &compile.MorpheStructFileWriter{
			Type:          compile.MorpheStructTypeModels,
			TargetDirPath: workingDirPath + "/models",
		},

		EnumWriter: &compile.MorpheEnumFileWriter{
			TargetDirPath: workingDirPath + "/enums",
		},
	}

	compileErr := compile.MorpheToGo(config)

	suite.NoError(compileErr)

	modelsDirPath := workingDirPath + "/models"
	gtModelsDirPath := suite.TestGroundTruthDirPath + "/models"
	suite.DirExists(modelsDirPath)

	modelPath0 := modelsDirPath + "/contact_info.go"
	gtModelPath0 := gtModelsDirPath + "/contact_info.go"
	suite.FileExists(modelPath0)
	suite.FileEquals(modelPath0, gtModelPath0)

	modelIDPath00 := modelsDirPath + "/contact_info_id_email.go"
	gtModelIDPath00 := gtModelsDirPath + "/contact_info_id_email.go"
	suite.FileExists(modelIDPath00)
	suite.FileEquals(modelIDPath00, gtModelIDPath00)

	modelIDPath01 := modelsDirPath + "/contact_info_id_primary.go"
	gtModelIDPath01 := gtModelsDirPath + "/contact_info_id_primary.go"
	suite.FileExists(modelIDPath01)
	suite.FileEquals(modelIDPath01, gtModelIDPath01)

	modelPath1 := modelsDirPath + "/company.go"
	gtModelPath1 := gtModelsDirPath + "/company.go"
	suite.FileExists(modelPath1)
	suite.FileEquals(modelPath1, gtModelPath1)

	modelIDPath10 := modelsDirPath + "/company_id_name.go"
	gtModelIDPath10 := gtModelsDirPath + "/company_id_name.go"
	suite.FileExists(modelIDPath10)
	suite.FileEquals(modelIDPath10, gtModelIDPath10)

	modelIDPath11 := modelsDirPath + "/company_id_primary.go"
	gtModelIDPath11 := gtModelsDirPath + "/company_id_primary.go"
	suite.FileExists(modelIDPath11)
	suite.FileEquals(modelIDPath11, gtModelIDPath11)

	modelPath2 := modelsDirPath + "/person.go"
	gtModelPath2 := gtModelsDirPath + "/person.go"
	suite.FileExists(modelPath2)
	suite.FileEquals(modelPath2, gtModelPath2)

	modelIDPath20 := modelsDirPath + "/person_id_name.go"
	gtModelIDPath20 := gtModelsDirPath + "/person_id_name.go"
	suite.FileExists(modelIDPath20)
	suite.FileEquals(modelIDPath20, gtModelIDPath20)

	modelIDPath21 := modelsDirPath + "/person_id_primary.go"
	gtModelIDPath21 := gtModelsDirPath + "/person_id_primary.go"
	suite.FileExists(modelIDPath21)
	suite.FileEquals(modelIDPath21, gtModelIDPath21)

	enumsDirPath := workingDirPath + "/enums"
	gtEnumsDirPath := suite.TestGroundTruthDirPath + "/enums"
	suite.DirExists(enumsDirPath)

	enumPath0 := enumsDirPath + "/nationality.go"
	gtEnumPath0 := gtEnumsDirPath + "/nationality.go"
	suite.FileExists(enumPath0)
	suite.FileEquals(enumPath0, gtEnumPath0)

	enumPath1 := enumsDirPath + "/universal_number.go"
	gtEnumPath1 := gtEnumsDirPath + "/universal_number.go"
	suite.FileExists(enumPath1)
	suite.FileEquals(enumPath1, gtEnumPath1)
}
