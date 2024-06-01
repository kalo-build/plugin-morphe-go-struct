package godef

import "github.com/kaloseia/clone"

type Struct struct {
	Package Package
	Imports []string
	Name    string
	Fields  []StructField
	Methods []StructMethod
}

func (s Struct) DeepClone() Struct {
	return Struct{
		Package: s.Package,
		Imports: clone.Slice(s.Imports),
		Name:    s.Name,
		Fields:  clone.DeepCloneSlice(s.Fields),
		Methods: clone.DeepCloneSlice(s.Methods),
	}
}
