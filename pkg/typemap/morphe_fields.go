package typemap

import (
	"github.com/kalo-build/morphe-go/pkg/yaml"

	"github.com/kalo-build/go/pkg/godef"
)

var MorpheModelFieldToGoField = map[yaml.ModelFieldType]godef.GoType{
	yaml.ModelFieldTypeUUID:          godef.GoTypeString,
	yaml.ModelFieldTypeAutoIncrement: godef.GoTypeUint,
	yaml.ModelFieldTypeString:        godef.GoTypeString,
	yaml.ModelFieldTypeInteger:       godef.GoTypeInt,
	yaml.ModelFieldTypeFloat:         godef.GoTypeFloat,
	yaml.ModelFieldTypeBoolean:       godef.GoTypeBool,
	yaml.ModelFieldTypeTime:          godef.GoTypeTime,
	yaml.ModelFieldTypeDate:          godef.GoTypeTime,
	yaml.ModelFieldTypeProtected:     godef.GoTypeString,
	yaml.ModelFieldTypeSealed:        godef.GoTypeString,
}

var MorpheStructureFieldToGoField = map[yaml.StructureFieldType]godef.GoType{
	yaml.StructureFieldTypeUUID:          godef.GoTypeString,
	yaml.StructureFieldTypeAutoIncrement: godef.GoTypeUint,
	yaml.StructureFieldTypeString:        godef.GoTypeString,
	yaml.StructureFieldTypeInteger:       godef.GoTypeInt,
	yaml.StructureFieldTypeFloat:         godef.GoTypeFloat,
	yaml.StructureFieldTypeBoolean:       godef.GoTypeBool,
	yaml.StructureFieldTypeTime:          godef.GoTypeTime,
	yaml.StructureFieldTypeDate:          godef.GoTypeTime,
	yaml.StructureFieldTypeProtected:     godef.GoTypeString,
	yaml.StructureFieldTypeSealed:        godef.GoTypeString,
}
