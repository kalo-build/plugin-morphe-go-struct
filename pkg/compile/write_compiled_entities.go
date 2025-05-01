package compile

import (
	"github.com/kalo-build/clone"
	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/go/pkg/godef"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/compile/hook"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/compile/write"
)

func WriteAllEntityStructDefinitions(config MorpheCompileConfig, allEntityStructDefs map[string][]*godef.Struct) (CompiledMorpheStructs, error) {
	allWrittenEntities := CompiledMorpheStructs{}

	sortedEntityNames := core.MapKeysSorted(allEntityStructDefs)
	for _, entityName := range sortedEntityNames {
		entityStructs := allEntityStructDefs[entityName]
		for _, entityStruct := range entityStructs {
			entityStruct, entityStructContents, writeErr := WriteEntityStructDefinition(config.WriteStructHooks, config.EntityWriter, entityStruct)
			if writeErr != nil {
				return nil, writeErr
			}
			allWrittenEntities.AddCompiledMorpheStruct(entityName, entityStruct, entityStructContents)
		}
	}
	return allWrittenEntities, nil
}

func WriteEntityStructDefinition(hooks hook.WriteGoStruct, writer write.GoStructWriter, entityStruct *godef.Struct) (*godef.Struct, []byte, error) {
	writer, entityStruct, writeStartErr := triggerWriteEntityStructStart(hooks, writer, entityStruct)
	if writeStartErr != nil {
		return nil, nil, triggerWriteEntityStructFailure(hooks, writer, entityStruct, writeStartErr)
	}

	entityStructContents, writeStructErr := writer.WriteStruct(entityStruct)
	if writeStructErr != nil {
		return nil, nil, triggerWriteEntityStructFailure(hooks, writer, entityStruct, writeStructErr)
	}

	entityStruct, entityStructContents, writeSuccessErr := triggerWriteEntityStructSuccess(hooks, entityStruct, entityStructContents)
	if writeSuccessErr != nil {
		return nil, nil, triggerWriteEntityStructFailure(hooks, writer, entityStruct, writeSuccessErr)
	}
	return entityStruct, entityStructContents, nil
}

func triggerWriteEntityStructStart(hooks hook.WriteGoStruct, writer write.GoStructWriter, entityStruct *godef.Struct) (write.GoStructWriter, *godef.Struct, error) {
	if hooks.OnWriteGoStructStart == nil {
		return writer, entityStruct, nil
	}
	if entityStruct == nil {
		return nil, nil, ErrNoEntityStruct
	}
	entityStructClone := entityStruct.DeepClone()

	updatedWriter, updatedEntityStruct, startErr := hooks.OnWriteGoStructStart(writer, &entityStructClone)
	if startErr != nil {
		return nil, nil, startErr
	}

	return updatedWriter, updatedEntityStruct, nil
}

func triggerWriteEntityStructSuccess(hooks hook.WriteGoStruct, entityStruct *godef.Struct, entityStructContents []byte) (*godef.Struct, []byte, error) {
	if hooks.OnWriteGoStructSuccess == nil {
		return entityStruct, entityStructContents, nil
	}
	if entityStruct == nil {
		return nil, nil, ErrNoEntityStruct
	}
	entityStructClone := entityStruct.DeepClone()
	entityStructContentsClone := clone.Slice(entityStructContents)

	updatedEntityStruct, updatedEntityStructContents, successErr := hooks.OnWriteGoStructSuccess(&entityStructClone, entityStructContentsClone)
	if successErr != nil {
		return nil, nil, successErr
	}
	return updatedEntityStruct, updatedEntityStructContents, nil
}

func triggerWriteEntityStructFailure(hooks hook.WriteGoStruct, writer write.GoStructWriter, entityStruct *godef.Struct, failureErr error) error {
	if hooks.OnWriteGoStructFailure == nil {
		return failureErr
	}

	entityStructClone := entityStruct.DeepClone()
	return hooks.OnWriteGoStructFailure(writer, &entityStructClone, failureErr)
}
