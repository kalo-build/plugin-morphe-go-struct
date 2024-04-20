package compile

import (
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/godef"
)

func MorpheEntityToGoStruct(packageName string, entity yaml.Entity) (*godef.Struct, error) {
	morpheStruct := godef.Struct{
		Package: packageName,
		Name:    entity.Name,
	}

	// TODO: Imports, Fields, ...
	return &morpheStruct, nil
}
