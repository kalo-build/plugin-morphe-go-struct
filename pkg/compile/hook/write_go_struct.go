package hook

import (
	"github.com/kaloseia/go/pkg/godef"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/write"
)

type WriteGoStruct struct {
	OnWriteGoStructStart   OnWriteGoStructStartHook
	OnWriteGoStructSuccess OnWriteGoStructSuccessHook
	OnWriteGoStructFailure OnWriteGoStructFailureHook
}

type OnWriteGoStructStartHook = func(writer write.GoStructWriter, modelStruct *godef.Struct) (write.GoStructWriter, *godef.Struct, error)
type OnWriteGoStructSuccessHook = func(modelStruct *godef.Struct, modelStructContents []byte) (*godef.Struct, []byte, error)
type OnWriteGoStructFailureHook = func(writer write.GoStructWriter, modelStruct *godef.Struct, failureErr error) error
