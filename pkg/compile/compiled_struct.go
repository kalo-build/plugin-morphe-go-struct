package compile

import "github.com/kaloseia/go/pkg/godef"

type CompiledStruct struct {
	Struct         *godef.Struct
	StructContents []byte
}
