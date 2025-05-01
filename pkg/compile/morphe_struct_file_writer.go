package compile

import (
	"fmt"
	"strings"

	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/go/pkg/godef"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/gofile"
)

type MorpheStructFileWriter struct {
	Type          MorpheStructType
	TargetDirPath string
}

func (w *MorpheStructFileWriter) WriteStruct(structDefinition *godef.Struct) ([]byte, error) {
	allStructLines, allLinesErr := w.getAllStructLines(structDefinition)
	if allLinesErr != nil {
		return nil, allLinesErr
	}

	structFileContents, structContentsErr := core.LinesToString(allStructLines)
	if structContentsErr != nil {
		return nil, structContentsErr
	}

	return gofile.WriteGoDefinitionFile(w.TargetDirPath, structDefinition.Name, structFileContents)
}

func (w *MorpheStructFileWriter) getAllStructLines(structDefinition *godef.Struct) ([]string, error) {
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

func (w *MorpheStructFileWriter) getAllStructImportLines(structDefinition *godef.Struct) ([]string, error) {
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

func (w *MorpheStructFileWriter) getAllStructTypeLines(structDefinition *godef.Struct) ([]string, error) {
	allTypeLines := []string{
		fmt.Sprintf("type %s struct {", structDefinition.Name),
	}

	for _, structField := range structDefinition.Fields {
		structFieldTypeSyntax := structField.Type.GetSyntax()
		if len(structField.Tags) == 0 {
			structFieldLine := fmt.Sprintf("\t%s %s", structField.Name, structFieldTypeSyntax)
			allTypeLines = append(allTypeLines, structFieldLine)
			continue
		}
		structFieldLine := fmt.Sprintf("\t%s %s `%s`", structField.Name, structFieldTypeSyntax, strings.Join(structField.Tags, " "))
		allTypeLines = append(allTypeLines, structFieldLine)
	}

	allTypeLines = append(allTypeLines, "}")
	return allTypeLines, nil
}

func (w *MorpheStructFileWriter) getAllStructMethodLines(currentPackage godef.Package, structMethods []godef.StructMethod) ([]string, error) {
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

func (w *MorpheStructFileWriter) getStructMethodLines(currentPackage godef.Package, structMethod godef.StructMethod) ([]string, error) {
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

func (w *MorpheStructFileWriter) getStructMethodParameterString(parameters map[string]godef.GoType) string {
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

func (w *MorpheStructFileWriter) getStructMethodReturnString(currentPackage godef.Package, returnTypes []godef.GoType) string {
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
