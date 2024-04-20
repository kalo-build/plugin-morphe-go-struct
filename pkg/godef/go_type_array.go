package godef

import (
	"fmt"
)

type GoTypeArray struct {
	IsSlice   bool
	FixedSize uint
	ValueType GoType
}

func (t GoTypeArray) IsPrimitive() bool {
	return false
}

func (t GoTypeArray) IsMap() bool {
	return false
}

func (t GoTypeArray) IsArray() bool {
	return true
}

func (t GoTypeArray) IsStruct() bool {
	return false
}

func (t GoTypeArray) IsInterface() bool {
	return false
}

func (t GoTypeArray) IsPointer() bool {
	return false
}

func (t GoTypeArray) GetImports() []string {
	return t.ValueType.GetImports()
}

func (t GoTypeArray) GetSyntax() string {
	if t.IsSlice {
		return fmt.Sprintf("[]%s", t.ValueType.GetSyntax())
	}
	return fmt.Sprintf("[%v]%s", t.FixedSize, t.ValueType.GetSyntax())
}
