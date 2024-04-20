package godef

type GoType interface {
	IsPrimitive() bool
	IsMap() bool
	IsArray() bool
	IsStruct() bool
	IsInterface() bool
	IsPointer() bool

	GetImports() []string
	GetSyntax() string
}
