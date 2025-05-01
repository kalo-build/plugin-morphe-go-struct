package write

import "github.com/kalo-build/go/pkg/godef"

type GoEnumWriter interface {
	WriteEnum(*godef.Enum) ([]byte, error)
}
