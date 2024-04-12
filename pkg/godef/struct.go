package godef

type Struct struct {
	Package string
	Imports []string
	Name    string
	Fields  []StructField
	Methods []StructMethod
}
