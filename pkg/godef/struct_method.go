package godef

import "github.com/kaloseia/clone"

type StructMethod struct {
	ReceiverName string
	ReceiverType GoType
	Name         string
	Parameters   map[string]GoType
	ReturnTypes  []GoType
	BodyLines    []string
}

func (m StructMethod) DeepClone() StructMethod {
	return StructMethod{
		Name:         m.Name,
		ReceiverName: m.Name,
		ReceiverType: DeepCloneGoType(m.ReceiverType),
		Parameters:   DeepCloneGoTypeMap(m.Parameters),
		ReturnTypes:  DeepCloneGoTypeSlice(m.ReturnTypes),
		BodyLines:    clone.Slice(m.BodyLines),
	}
}
