package godef

type StructMethod struct {
	ReceiverName string
	ReceiverType GoType
	Name         string
	ReturnTypes  []GoType
	BodyLines    []string
}
