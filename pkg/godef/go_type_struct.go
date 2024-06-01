package godef

import (
	"strings"
)

type GoTypeStruct struct {
	// PackagePath is the fully qualified import path, ie. "time", "github.com/<org>/<package>", ...
	PackagePath string

	// Name is the name of the struct, ie. "Time", ...
	Name string
}

func (t GoTypeStruct) IsPrimitive() bool {
	return false
}

func (t GoTypeStruct) IsMap() bool {
	return false
}

func (t GoTypeStruct) IsArray() bool {
	return false
}

func (t GoTypeStruct) IsStruct() bool {
	return true
}

func (t GoTypeStruct) IsInterface() bool {
	return false
}

func (t GoTypeStruct) IsPointer() bool {
	return false
}

func (t GoTypeStruct) GetImports() []string {
	return []string{t.PackagePath}
}

func (t GoTypeStruct) GetSyntax() string {
	packageName := t.getPackageName()
	if packageName == "" {
		return t.GetSyntaxLocal()
	}
	return packageName + "." + t.Name
}

func (t GoTypeStruct) GetSyntaxLocal() string {
	return t.Name
}

func (t GoTypeStruct) getPackageName() string {
	if t.PackagePath == "" {
		return ""
	}
	packagePieces := strings.Split(t.PackagePath, "/")
	return packagePieces[len(packagePieces)-1]
}

func (t GoTypeStruct) DeepClone() GoTypeStruct {
	return GoTypeStruct{
		PackagePath: t.PackagePath,
		Name:        t.Name,
	}
}
