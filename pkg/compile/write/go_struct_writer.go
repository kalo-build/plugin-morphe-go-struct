package write

import "github.com/kaloseia/plugin-morphe-go-struct/pkg/godef"

type GoStructWriter interface {
	WriteStruct(*godef.Struct) ([]byte, error)
}
