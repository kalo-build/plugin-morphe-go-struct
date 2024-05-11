package compile

import "github.com/kaloseia/plugin-morphe-go-struct/pkg/godef"

type StructWriter interface {
	WriteStruct(*godef.Struct) error
}
