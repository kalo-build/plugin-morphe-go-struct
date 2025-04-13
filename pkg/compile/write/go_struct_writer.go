package write

import "github.com/kalo-build/go/pkg/godef"

type GoStructWriter interface {
	WriteStruct(*godef.Struct) ([]byte, error)
}
