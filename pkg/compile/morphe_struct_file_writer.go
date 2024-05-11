package compile

import (
	"fmt"
	"strings"

	"github.com/kaloseia/plugin-morphe-go-struct/pkg/core"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/godef"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/gofile"
)

type MorpheStructWriter struct {
	Type          MorpheStructType
	TargetDirPath string
}

func (w *MorpheStructWriter) WriteStruct(structDefinition *godef.Struct) error {
	allStructLines, allLinesErr := w.getAllStructLines(structDefinition)
	if allLinesErr != nil {
		return allLinesErr
	}

	structFileContents, structContentsErr := core.LinesToString(allStructLines)
	if structContentsErr != nil {
		return structContentsErr
	}

	return gofile.WriteGoStructFile(w.TargetDirPath, structDefinition.Name, structFileContents)
}

func (w *MorpheStructWriter) getAllStructLines(structDefinition *godef.Struct) ([]string, error) {
	allStructLines := []string{}

	packageLine := fmt.Sprintf("package %s", structDefinition.Package.Name)
	allStructLines = append(allStructLines, packageLine)
	allStructLines = append(allStructLines, "")

	importLines, importsErr := w.getAllStructImportLines(structDefinition)
	if importsErr != nil {
		return nil, importsErr
	}
	allStructLines = append(allStructLines, importLines...)
	allStructLines = append(allStructLines, "")

	typeLines, typeErr := w.getAllStructTypeLines(structDefinition)
	if typeErr != nil {
		return nil, typeErr
	}
	allStructLines = append(allStructLines, typeLines...)
	allStructLines = append(allStructLines, "")

	methodLines, methodErr := w.getAllStructMethodLines(structDefinition.Package, structDefinition.Methods)
	if methodErr != nil {
		return nil, methodErr
	}
	allStructLines = append(allStructLines, methodLines...)
	allStructLines = append(allStructLines, "")

	return allStructLines, nil
}

func (w *MorpheStructWriter) getAllStructImportLines(structDefinition *godef.Struct) ([]string, error) {
	if len(structDefinition.Imports) == 0 {
		return nil, nil
	}

	filteredImportsMap := map[string]any{}
	for _, structImport := range structDefinition.Imports {
		filteredImportsMap[structImport] = nil
	}

	allImportLines := []string{
		"import (",
	}

	filteredImports := core.MapKeysSorted(filteredImportsMap)
	for _, structImport := range filteredImports {
		allImportLines = append(allImportLines, `"`+structImport+`"`)
	}

	allImportLines = append(allImportLines, ")")
	return allImportLines, nil
}

func (w *MorpheStructWriter) getAllStructTypeLines(structDefinition *godef.Struct) ([]string, error) {

	allTypeLines := []string{
		fmt.Sprintf("type %s struct {", structDefinition.Name),
	}

	for _, structField := range structDefinition.Fields {
		structFieldTypeSyntax := structField.Type.GetSyntax()
		structFieldTags := strings.Join(structField.Tags, " ")
		if structFieldTags == "" {
			structFieldLine := fmt.Sprintf("\t%s %s", structField.Name, structFieldTypeSyntax)
			allTypeLines = append(allTypeLines, structFieldLine)
			continue
		}
		structFieldLine := fmt.Sprintf("\t%s %s `%s`", structField.Name, structFieldTypeSyntax, structFieldTags)
		allTypeLines = append(allTypeLines, structFieldLine)
	}

	allTypeLines = append(allTypeLines, "}")
	return allTypeLines, nil
}

func (w *MorpheStructWriter) getAllStructMethodLines(currentPackage godef.Package, structMethods []godef.StructMethod) ([]string, error) {
	allMethodLines := []string{}

	for _, structMethod := range structMethods {
		structMethodLines, methodErr := w.getStructMethodLines(currentPackage, structMethod)
		if methodErr != nil {
			return nil, methodErr
		}
		allMethodLines = append(allMethodLines, structMethodLines...)
		allMethodLines = append(allMethodLines, "")
	}

	return allMethodLines, nil
}

func (w *MorpheStructWriter) getStructMethodLines(currentPackage godef.Package, structMethod godef.StructMethod) ([]string, error) {
	receiverName := structMethod.ReceiverName
	receiverType := structMethod.ReceiverType.GetSyntaxLocal()
	methodParamBlock := w.getStructMethodParameterString(structMethod.Parameters)
	methodReturnBlock := w.getStructMethodReturnString(currentPackage, structMethod.ReturnTypes)

	methodHeader := fmt.Sprintf("func (%s %s) %s(%s) %s {", receiverName, receiverType, structMethod.Name, methodParamBlock, methodReturnBlock)
	methodLines := []string{
		methodHeader,
	}
	for _, bodyLine := range structMethod.BodyLines {
		methodLines = append(methodLines, "\t"+bodyLine)
	}
	methodLines = append(methodLines, "}")
	return methodLines, nil
}

func (w *MorpheStructWriter) getStructMethodParameterString(parameters map[string]godef.GoType) string {
	if parameters == nil {
		return ""
	}

	parameterStrings := []string{}
	sortedParamNames := core.MapKeysSorted(parameters)
	for _, paramName := range sortedParamNames {
		paramType := parameters[paramName]
		parameterStrings = append(parameterStrings, fmt.Sprintf("%s %s", paramName, paramType.GetSyntax()))
	}

	return strings.Join(parameterStrings, ", ")
}

func (w *MorpheStructWriter) getStructMethodReturnString(currentPackage godef.Package, returnTypes []godef.GoType) string {
	if returnTypes == nil {
		return ""
	}

	returnStrings := []string{}
	for _, returnType := range returnTypes {
		returnTypeStruct, isStructType := returnType.(godef.GoTypeStruct)
		if isStructType && returnTypeStruct.PackagePath == currentPackage.Path {
			returnStrings = append(returnStrings, returnType.GetSyntaxLocal())
			continue
		}
		returnStrings = append(returnStrings, returnType.GetSyntax())
	}

	if len(returnStrings) == 1 {
		return returnStrings[0]
	}

	return fmt.Sprintf("(%s)", strings.Join(returnStrings, ", "))
}
