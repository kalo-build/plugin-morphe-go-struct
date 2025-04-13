package compile

import (
	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/go/pkg/godef"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/compile/hook"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/compile/write"
)

type CompiledEnums struct {
	Enums    map[string]*godef.Enum
	Contents map[string][]byte
}

func (c *CompiledEnums) AddCompiledEnum(enum *godef.Enum, contents []byte) {
	if c.Enums == nil {
		c.Enums = map[string]*godef.Enum{}
	}
	if c.Contents == nil {
		c.Contents = map[string][]byte{}
	}
	c.Enums[enum.Name] = enum
	c.Contents[enum.Name] = contents
}

func WriteAllEnumDefinitions(config MorpheCompileConfig, allEnumDefs map[string]*godef.Enum) (CompiledEnums, error) {
	allWrittenEnums := CompiledEnums{}

	sortedEnumNames := core.MapKeysSorted(allEnumDefs)
	for _, enumName := range sortedEnumNames {
		enumDef := allEnumDefs[enumName]
		enumDef, enumContents, writeErr := WriteEnumDefinition(config.WriteGoEnumHooks, config.EnumWriter, enumDef)
		if writeErr != nil {
			return CompiledEnums{}, writeErr
		}
		allWrittenEnums.AddCompiledEnum(enumDef, enumContents)
	}
	return allWrittenEnums, nil
}

func WriteEnumDefinition(hooks hook.WriteGoEnum, writer write.GoEnumWriter, enum *godef.Enum) (*godef.Enum, []byte, error) {
	writer, enum, writeStartErr := triggerWriteEnumStart(hooks, writer, enum)
	if writeStartErr != nil {
		return nil, nil, triggerWriteEnumFailure(hooks, writer, enum, writeStartErr)
	}

	enumContents, writeEnumErr := writer.WriteEnum(enum)
	if writeEnumErr != nil {
		return nil, nil, triggerWriteEnumFailure(hooks, writer, enum, writeEnumErr)
	}

	enum, enumContents, writeSuccessErr := triggerWriteEnumSuccess(hooks, enum, enumContents)
	if writeSuccessErr != nil {
		return nil, nil, triggerWriteEnumFailure(hooks, writer, enum, writeSuccessErr)
	}
	return enum, enumContents, nil
}

func triggerWriteEnumStart(hooks hook.WriteGoEnum, writer write.GoEnumWriter, enum *godef.Enum) (write.GoEnumWriter, *godef.Enum, error) {
	if hooks.OnWriteGoEnumStart == nil {
		return writer, enum, nil
	}
	if enum == nil {
		return nil, nil, ErrNoEnum
	}
	enumClone := enum.DeepClone()

	updatedWriter, updatedEnum, startErr := hooks.OnWriteGoEnumStart(writer, &enumClone)
	if startErr != nil {
		return nil, nil, startErr
	}

	return updatedWriter, updatedEnum, nil
}

func triggerWriteEnumSuccess(hooks hook.WriteGoEnum, enum *godef.Enum, enumContents []byte) (*godef.Enum, []byte, error) {
	if hooks.OnWriteGoEnumSuccess == nil {
		return enum, enumContents, nil
	}
	if enum == nil {
		return nil, nil, ErrNoEnum
	}
	enumClone := enum.DeepClone()
	enumContentsClone := make([]byte, len(enumContents))
	copy(enumContentsClone, enumContents)

	updatedEnum, updatedEnumContents, successErr := hooks.OnWriteGoEnumSuccess(&enumClone, enumContentsClone)
	if successErr != nil {
		return nil, nil, successErr
	}
	return updatedEnum, updatedEnumContents, nil
}

func triggerWriteEnumFailure(hooks hook.WriteGoEnum, writer write.GoEnumWriter, enum *godef.Enum, failureErr error) error {
	if hooks.OnWriteGoEnumFailure == nil {
		return failureErr
	}

	enumClone := enum.DeepClone()
	return hooks.OnWriteGoEnumFailure(writer, &enumClone, failureErr)
}
