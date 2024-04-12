package typemap

import (
	"github.com/kaloseia/morphe-go/pkg/yaml"

	"github.com/kaloseia/plugin-morphe-go-struct/pkg/godef"
)

var MorpheFieldToGoField = map[yaml.ModelFieldType]godef.GoType{
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
