package godef

type StructMethod struct {
	ReceiverName string
	ReceiverType GoType
	Name         string
	Parameters   map[string]GoType
	ReturnTypes  []GoType
	BodyLines    []string
}
