package compile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/kaloseia/go/pkg/godef"
	rcfg "github.com/kaloseia/morphe-go/pkg/registry/cfg"
	"github.com/kaloseia/plugin-morphe-go-struct/internal/testutils"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/cfg"
)

type CompileTestSuite struct {
	testutils.FileSuite

	TestDirPath            string
	TestGroundTruthDirPath string

	ModelsDirPath   string
	EntitiesDirPath string
}

func TestCompileTestSuite(t *testing.T) {
	suite.Run(t, new(CompileTestSuite))
}

func (suite *CompileTestSuite) SetupTest() {
	suite.TestDirPath = testutils.GetTestDirPath()
	suite.TestGroundTruthDirPath = filepath.Join(suite.TestDirPath, "ground-truth", "compile-minimal")

	suite.ModelsDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "models")
	suite.EntitiesDirPath = filepath.Join(suite.TestDirPath, "registry", "minimal", "entities")
}

func (suite *CompileTestSuite) TearDownTest() {
	suite.TestDirPath = ""
}

func (suite *CompileTestSuite) TestMorpheToGoStructs() {
	workingDirPath := suite.TestDirPath + "/working"
	suite.Nil(os.Mkdir(workingDirPath, 0644))
	defer os.RemoveAll(workingDirPath)

	config := compile.MorpheCompileConfig{
		MorpheLoadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryModelsDirPath:   suite.ModelsDirPath,
			RegistryEntitiesDirPath: suite.EntitiesDirPath,
		},

		MorpheModelsConfig: cfg.MorpheModelsConfig{
			Package: godef.Package{
				Path: "github.com/kaloseia/dummy/models",
				Name: "models",
			},
			ReceiverName: "m",
		},

		ModelWriter: &compile.MorpheStructFileWriter{
			Type:          compile.MorpheStructTypeModels,
			TargetDirPath: workingDirPath + "/models",
		},
	}

	allWrittenModels, compileErr := compile.MorpheToGoStructs(config)

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

	modelPath1 := modelsDirPath + "/person.go"
	gtModelPath1 := gtModelsDirPath + "/person.go"
	suite.FileExists(modelPath1)
	suite.FileEquals(modelPath1, gtModelPath1)

	modelIDPath10 := modelsDirPath + "/person_id_name.go"
	gtModelIDPath10 := gtModelsDirPath + "/person_id_name.go"
	suite.FileExists(modelIDPath10)
	suite.FileEquals(modelIDPath10, gtModelIDPath10)

	modelIDPath11 := modelsDirPath + "/person_id_primary.go"
	gtModelIDPath11 := gtModelsDirPath + "/person_id_primary.go"
	suite.FileExists(modelIDPath11)
	suite.FileEquals(modelIDPath11, gtModelIDPath11)

	suite.Len(allWrittenModels, 2)

	// Contact Info
	model0 := allWrittenModels.GetAllCompiledModelStructs("ContactInfo")
	suite.Len(model0, 3)

	model00 := allWrittenModels.GetCompiledModelStruct("ContactInfo", "ContactInfo")
	suite.FileContentsEquals(gtModelPath0, model00.StructContents)
	suite.NotNil(model00.Struct)

	model01 := allWrittenModels.GetCompiledModelStruct("ContactInfo", "ContactInfoIDEmail")
	suite.FileContentsEquals(gtModelIDPath00, model01.StructContents)
	suite.NotNil(model01.Struct)

	model02 := allWrittenModels.GetCompiledModelStruct("ContactInfo", "ContactInfoIDPrimary")
	suite.FileContentsEquals(gtModelIDPath01, model02.StructContents)
	suite.NotNil(model02.Struct)

	// Person
	model1 := allWrittenModels.GetAllCompiledModelStructs("Person")
	suite.Len(model1, 3)

	model10 := allWrittenModels.GetCompiledModelStruct("Person", "Person")
	suite.FileContentsEquals(gtModelPath1, model10.StructContents)
	suite.NotNil(model10.Struct)

	model11 := allWrittenModels.GetCompiledModelStruct("Person", "PersonIDName")
	suite.FileContentsEquals(gtModelIDPath10, model11.StructContents)
	suite.NotNil(model11.Struct)

	model12 := allWrittenModels.GetCompiledModelStruct("Person", "PersonIDPrimary")
	suite.FileContentsEquals(gtModelIDPath11, model12.StructContents)
	suite.NotNil(model11.Struct)

}
