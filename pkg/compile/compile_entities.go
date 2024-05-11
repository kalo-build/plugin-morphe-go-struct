package compile

import (
	"github.com/kaloseia/morphe-go/pkg/yaml"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/godef"
)

func MorpheEntityToGoStruct(structPackage godef.Package, entity yaml.Entity) (*godef.Struct, error) {
	morpheStruct := godef.Struct{
		Package: structPackage,
		Name:    entity.Name,
	}

	return &morpheStruct, nil
}
