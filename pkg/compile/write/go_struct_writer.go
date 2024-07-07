package write

import "github.com/kaloseia/go/pkg/godef"

type GoStructWriter interface {
	WriteStruct(*godef.Struct) ([]byte, error)
}
