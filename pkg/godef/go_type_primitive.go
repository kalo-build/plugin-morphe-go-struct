package godef

type GoTypePrimitive struct {
	Syntax string
}

func (t GoTypePrimitive) IsPrimitive() bool {
	return true
}

func (t GoTypePrimitive) IsMap() bool {
	return false
}

func (t GoTypePrimitive) IsArray() bool {
	return false
}

func (t GoTypePrimitive) IsStruct() bool {
	return false
}

func (t GoTypePrimitive) IsPointer() bool {
	return false
}

func (t GoTypePrimitive) GetImports() []string {
	return nil
}

func (t GoTypePrimitive) GetSyntax() string {
	return t.Syntax
}
