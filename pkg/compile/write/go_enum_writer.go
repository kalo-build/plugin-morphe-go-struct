package write

import "github.com/kaloseia/go/pkg/godef"

type GoEnumWriter interface {
	WriteEnum(*godef.Enum) ([]byte, error)
}
