package godef

type GoType interface {
	IsPrimitive() bool
	IsMap() bool
	IsArray() bool
	IsStruct() bool
	IsPointer() bool

	GetImports() []string
	GetSyntax() string
}
