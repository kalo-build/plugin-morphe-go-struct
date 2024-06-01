package godef

import "github.com/kaloseia/clone"

type StructField struct {
	Name string
	Type GoType
	Tags []string
}

func (f StructField) DeepClone() StructField {
	return StructField{
		Name: f.Name,
		Type: DeepCloneGoType(f.Type),
		Tags: clone.Slice(f.Tags),
	}
}
