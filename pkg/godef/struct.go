package godef

type Struct struct {
	Package Package
	Imports []string
	Name    string
	Fields  []StructField
	Methods []StructMethod
}
