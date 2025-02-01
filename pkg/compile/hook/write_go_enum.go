package hook

import (
	"github.com/kaloseia/go/pkg/godef"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/write"
)

type WriteGoEnum struct {
	OnWriteGoEnumStart   OnWriteGoEnumStartHook
	OnWriteGoEnumSuccess OnWriteGoEnumSuccessHook
	OnWriteGoEnumFailure OnWriteGoEnumFailureHook
}

type OnWriteGoEnumStartHook = func(writer write.GoEnumWriter, enum *godef.Enum) (write.GoEnumWriter, *godef.Enum, error)
type OnWriteGoEnumSuccessHook = func(enum *godef.Enum, enumContents []byte) (*godef.Enum, []byte, error)
type OnWriteGoEnumFailureHook = func(writer write.GoEnumWriter, enum *godef.Enum, failureErr error) error
