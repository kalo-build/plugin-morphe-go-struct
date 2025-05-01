package hook

import (
	"github.com/kalo-build/go/pkg/godef"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/compile/write"
)

type WriteGoStruct struct {
	OnWriteGoStructStart   OnWriteGoStructStartHook
	OnWriteGoStructSuccess OnWriteGoStructSuccessHook
	OnWriteGoStructFailure OnWriteGoStructFailureHook
}

type OnWriteGoStructStartHook = func(writer write.GoStructWriter, goStruct *godef.Struct) (write.GoStructWriter, *godef.Struct, error)
type OnWriteGoStructSuccessHook = func(goStruct *godef.Struct, goStructContents []byte) (*godef.Struct, []byte, error)
type OnWriteGoStructFailureHook = func(writer write.GoStructWriter, goStruct *godef.Struct, failureErr error) error
