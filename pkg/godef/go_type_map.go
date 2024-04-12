package godef

import (
	"fmt"
	"sort"
)

type GoTypeMap struct {
	KeyType   GoType
	ValueType GoType
}

func (t GoTypeMap) IsPrimitive() bool {
	return false
}

func (t GoTypeMap) IsMap() bool {
	return true
}

func (t GoTypeMap) IsArray() bool {
	return false
}

func (t GoTypeMap) IsStruct() bool {
	return false
}

func (t GoTypeMap) IsPointer() bool {
	return false
}

func (t GoTypeMap) GetImports() []string {
	mapImportMap := map[string]any{}

	keyImports := t.KeyType.GetImports()
	for _, keyImport := range keyImports {
		mapImportMap[keyImport] = nil
	}

	valueImports := t.ValueType.GetImports()
	for _, valueImport := range valueImports {
		mapImportMap[valueImport] = nil
	}

	allMapImports := []string{}
	for importPath := range mapImportMap {
		allMapImports = append(allMapImports, importPath)
	}
	sort.Strings(allMapImports)

	return allMapImports
}

func (t GoTypeMap) GetSyntax() string {
	return fmt.Sprintf("map[%s]%s", t.KeyType.GetSyntax(), t.ValueType.GetSyntax())
}
