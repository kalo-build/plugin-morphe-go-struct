package compile

import "github.com/kalo-build/go/pkg/godef"

type CompiledStruct struct {
	Struct         *godef.Struct
	StructContents []byte
}
