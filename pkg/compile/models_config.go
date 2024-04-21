package compile

type ModelsConfig struct {
	// PackagePath is the package import path for the compiled models package, ie. "github.com/myorg/myproject/models"
	PackagePath string

	// PackageName is the package import name for the compiled models package, ie. "models"
	PackageName string

	// ReceiverName is the standard receiver name for the compiled model receiver methods, ie "m" in "func (m *MyModel) MyMethod(){}"
	ReceiverName string
}
