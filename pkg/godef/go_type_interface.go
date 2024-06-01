package godef

import (
	"strings"
)

type GoTypeInterface struct {
	// PackagePath is the fully qualified import path, ie. "time", "github.com/<org>/<package>", ...
	PackagePath string

	// Name is the name of the interface, ie. "error", ...
	Name string
}

func (t GoTypeInterface) IsPrimitive() bool {
	return false
}

func (t GoTypeInterface) IsMap() bool {
	return false
}

func (t GoTypeInterface) IsArray() bool {
	return false
}

func (t GoTypeInterface) IsStruct() bool {
	return false
}

func (t GoTypeInterface) IsInterface() bool {
	return true
}

func (t GoTypeInterface) IsPointer() bool {
	return false
}

func (t GoTypeInterface) GetImports() []string {
	return []string{t.PackagePath}
}

func (t GoTypeInterface) GetSyntax() string {
	packageName := t.getPackageName()
	if packageName == "" {
		return t.GetSyntaxLocal()
	}
	return packageName + "." + t.Name
}

func (t GoTypeInterface) GetSyntaxLocal() string {
	return t.Name
}

func (t GoTypeInterface) getPackageName() string {
	if t.PackagePath == "" {
		return ""
	}
	packagePieces := strings.Split(t.PackagePath, "/")
	return packagePieces[len(packagePieces)-1]
}

func (t GoTypeInterface) DeepClone() GoTypeInterface {
	return GoTypeInterface{
		PackagePath: t.PackagePath,
		Name:        t.Name,
	}
}
