package godef

type Package struct {
	// Path is the package import path for the compiled models package, ie. "github.com/myorg/myproject/models"
	Path string

	// Name is the package import name for the compiled models package, ie. "models"
	Name string
}
