package godef

import (
	"fmt"
)

type GoTypePointer struct {
	ValueType GoType
}

func (t GoTypePointer) IsPrimitive() bool {
	return false
}

func (t GoTypePointer) IsMap() bool {
	return false
}

func (t GoTypePointer) IsArray() bool {
	return false
}

func (t GoTypePointer) IsStruct() bool {
	return false
}

func (t GoTypePointer) IsInterface() bool {
	return false
}

func (t GoTypePointer) IsPointer() bool {
	return true
}

func (t GoTypePointer) GetImports() []string {
	return t.ValueType.GetImports()
}

func (t GoTypePointer) GetSyntax() string {
	return fmt.Sprintf("*%s", t.ValueType.GetSyntax())
}

func (t GoTypePointer) GetSyntaxLocal() string {
	return fmt.Sprintf("*%s", t.ValueType.GetSyntaxLocal())
}

func (t GoTypePointer) DeepClone() GoTypePointer {
	return GoTypePointer{
		ValueType: DeepCloneGoType(t.ValueType),
	}
}
