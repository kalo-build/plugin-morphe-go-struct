package compile

import "github.com/kaloseia/plugin-morphe-go-struct/pkg/godef"

type CompiledStruct struct {
	Struct         *godef.Struct
	StructContents []byte
}
