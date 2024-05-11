package strcase_test

import (
	"testing"

	"github.com/kaloseia/plugin-morphe-go-struct/pkg/strcase"
	"github.com/stretchr/testify/suite"
)

type StringCaseTestSuite struct {
	suite.Suite
}

func TestStringCaseTestSuite(t *testing.T) {
	suite.Run(t, new(StringCaseTestSuite))
}

func (suite *StringCaseTestSuite) TestToSnakeCase() {
	suite.Equal(strcase.ToSnakeCase("UpperCamelCase"), "Upper_Camel_Case")
	suite.Equal(strcase.ToSnakeCase("lowerCamelCase"), "lower_Camel_Case")

	suite.Equal(strcase.ToSnakeCase("Upper-Kebab-Case"), "Upper_Kebab_Case")
	suite.Equal(strcase.ToSnakeCase("lower-kebab-case"), "lower_kebab_case")

	suite.Equal(strcase.ToSnakeCase("Upper_Snake_Case"), "Upper_Snake_Case")
	suite.Equal(strcase.ToSnakeCase("lower_snake_case"), "lower_snake_case")

	suite.Equal(strcase.ToSnakeCase("CaseIDThing"), "Case_ID_Thing")
	suite.Equal(strcase.ToSnakeCase("CaseID"), "Case_ID")
	suite.Equal(strcase.ToSnakeCase("IDCase"), "ID_Case")

	suite.Equal(strcase.ToSnakeCase("CaseSQLThing"), "Case_SQL_Thing")
	suite.Equal(strcase.ToSnakeCase("CaseSQL"), "Case_SQL")
	suite.Equal(strcase.ToSnakeCase("SQLCase"), "SQL_Case")

	suite.Equal(strcase.ToSnakeCase("Single"), "Single")
	suite.Equal(strcase.ToSnakeCase("SINGLE"), "SINGLE")
	suite.Equal(strcase.ToSnakeCase("single"), "single")

	suite.Equal(strcase.ToSnakeCase("num01Word"), "num01_Word")
	suite.Equal(strcase.ToSnakeCase("01Word"), "01_Word")
	suite.Equal(strcase.ToSnakeCase("Word01"), "Word01")
}

func (suite *StringCaseTestSuite) TestToSnakeCaseLower() {
	suite.Equal(strcase.ToSnakeCaseLower("UpperCamelCase"), "upper_camel_case")
	suite.Equal(strcase.ToSnakeCaseLower("lowerCamelCase"), "lower_camel_case")

	suite.Equal(strcase.ToSnakeCaseLower("Upper-Kebab-Case"), "upper_kebab_case")
	suite.Equal(strcase.ToSnakeCaseLower("lower-kebab-case"), "lower_kebab_case")

	suite.Equal(strcase.ToSnakeCaseLower("Upper_Snake_Case"), "upper_snake_case")
	suite.Equal(strcase.ToSnakeCaseLower("lower_snake_case"), "lower_snake_case")

	suite.Equal(strcase.ToSnakeCaseLower("CaseIDThing"), "case_id_thing")
	suite.Equal(strcase.ToSnakeCaseLower("CaseID"), "case_id")
	suite.Equal(strcase.ToSnakeCaseLower("IDCase"), "id_case")

	suite.Equal(strcase.ToSnakeCaseLower("CaseSQLThing"), "case_sql_thing")
	suite.Equal(strcase.ToSnakeCaseLower("CaseSQL"), "case_sql")
	suite.Equal(strcase.ToSnakeCaseLower("SQLCase"), "sql_case")

	suite.Equal(strcase.ToSnakeCaseLower("Single"), "single")
	suite.Equal(strcase.ToSnakeCaseLower("SINGLE"), "single")
	suite.Equal(strcase.ToSnakeCaseLower("single"), "single")

	suite.Equal(strcase.ToSnakeCaseLower("num01Word"), "num01_word")
	suite.Equal(strcase.ToSnakeCaseLower("01Word"), "01_word")
	suite.Equal(strcase.ToSnakeCaseLower("Word01"), "word01")
}
