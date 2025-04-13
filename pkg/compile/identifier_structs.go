package compile

import (
	"fmt"
	"sort"

	"github.com/kalo-build/go-util/strcase"
	"github.com/kalo-build/go/pkg/godef"
)

type IdentifierConfig interface {
	GetPackage() godef.Package
	GetReceiverName() string
}

type Identifier interface {
	GetFields() []string
}

// Common function to get identifier structs for both models and entities
func GetIdentifierStructs(
	config IdentifierConfig,
	parentName string,
	parentStruct *godef.Struct,
	identifiers map[string]Identifier,
) ([]*godef.Struct, error) {
	allIdentifierStructs := []*godef.Struct{}
	allIdentifierNames := getSortedIdentifierNames(identifiers)

	for _, identifierName := range allIdentifierNames {
		identifierDef := identifiers[identifierName]

		allIdentFieldDefs, identFieldDefsErr := getIdentifierStructFieldSubset(
			*parentStruct,
			identifierName,
			identifierDef.GetFields(),
		)
		if identFieldDefsErr != nil {
			return nil, identFieldDefsErr
		}

		identStruct, identStructErr := getIdentifierStruct(
			config.GetPackage(),
			parentName,
			identifierName,
			allIdentFieldDefs,
		)
		if identStructErr != nil {
			return nil, identStructErr
		}
		allIdentifierStructs = append(allIdentifierStructs, identStruct)

		getter, getterErr := getIdentifierGetter(
			config,
			parentName,
			identifierName,
			identStruct,
		)
		if getterErr != nil {
			return nil, getterErr
		}
		parentStruct.Methods = append(parentStruct.Methods, getter)
	}

	return allIdentifierStructs, nil
}

func getIdentifierStructFieldSubset(
	parentStruct godef.Struct,
	identifierName string,
	fields []string,
) ([]godef.StructField, error) {
	identifierFieldDefs := []godef.StructField{}

	for _, fieldName := range fields {
		field, found := findStructFieldByName(parentStruct.Fields, fieldName)
		if !found {
			return nil, fmt.Errorf("identifier %s references unknown field: %s", identifierName, fieldName)
		}
		identifierFieldDefs = append(identifierFieldDefs, field)
	}

	return identifierFieldDefs, nil
}

func getIdentifierStruct(
	pkg godef.Package,
	parentName string,
	identifierName string,
	fields []godef.StructField,
) (*godef.Struct, error) {
	structImports, importsErr := getImportsForStructFields(pkg, fields)
	if importsErr != nil {
		return nil, importsErr
	}

	return &godef.Struct{
		Package: pkg,
		Imports: structImports,
		Name:    fmt.Sprintf("%sID%s", parentName, strcase.ToPascalCase(identifierName)),
		Fields:  fields,
	}, nil
}

func getIdentifierGetter(
	config IdentifierConfig,
	parentName string,
	identifierName string,
	identStruct *godef.Struct,
) (godef.StructMethod, error) {
	returnType := godef.GoTypeStruct{
		PackagePath: config.GetPackage().Path,
		Name:        identStruct.Name,
	}

	bodyLines := getIdentifierGetterBodyLines(identStruct, config.GetReceiverName())

	return godef.StructMethod{
		ReceiverName: config.GetReceiverName(),
		ReceiverType: godef.GoTypeStruct{
			PackagePath: config.GetPackage().Path,
			Name:        parentName,
		},
		Name:        fmt.Sprintf("GetID%s", strcase.ToPascalCase(identifierName)),
		ReturnTypes: []godef.GoType{returnType},
		BodyLines:   bodyLines,
	}, nil
}

// Helper functions
func getSortedIdentifierNames(identifiers map[string]Identifier) []string {
	names := make([]string, 0, len(identifiers))
	for name := range identifiers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func findStructFieldByName(fields []godef.StructField, name string) (godef.StructField, bool) {
	for _, field := range fields {
		if field.Name == name {
			return field, true
		}
	}
	return godef.StructField{}, false
}

func getIdentifierGetterBodyLines(identStruct *godef.Struct, receiverName string) []string {
	bodyLines := []string{
		fmt.Sprintf("\treturn %s{", identStruct.Name),
	}
	for _, fieldDef := range identStruct.Fields {
		fieldLine := fmt.Sprintf("\t\t%s: %s.%s,", fieldDef.Name, receiverName, fieldDef.Name)
		bodyLines = append(bodyLines, fieldLine)
	}
	bodyLines = append(bodyLines, "\t}")
	return bodyLines
}
