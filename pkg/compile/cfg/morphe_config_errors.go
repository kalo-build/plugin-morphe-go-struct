package cfg

import "errors"

var ErrNoPackagePath = errors.New("package path cannot be empty")
var ErrNoPackageName = errors.New("package name cannot be empty")
var ErrNoReceiverName = errors.New("method receiver name cannot be empty")
