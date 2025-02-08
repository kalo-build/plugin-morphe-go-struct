package compile

import (
	"github.com/kaloseia/clone"
	"github.com/kaloseia/go-util/core"
	"github.com/kaloseia/go/pkg/godef"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/hook"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/write"
)

func WriteAllStructureStructDefinitions(config MorpheCompileConfig, allStructureStructDefs map[string]*godef.Struct) (CompiledMorpheStructs, error) {
	allWrittenStructures := CompiledMorpheStructs{}

	sortedStructureNames := core.MapKeysSorted(allStructureStructDefs)
	for _, structureName := range sortedStructureNames {
		structureStruct := allStructureStructDefs[structureName]
		structureStruct, structureStructContents, writeErr := WriteStructureStructDefinition(config.WriteStructHooks, config.StructureWriter, structureStruct)
		if writeErr != nil {
			return nil, writeErr
		}
		allWrittenStructures.AddCompiledMorpheStruct(structureName, structureStruct, structureStructContents)
	}
	return allWrittenStructures, nil
}

func WriteStructureStructDefinition(hooks hook.WriteGoStruct, writer write.GoStructWriter, structureStruct *godef.Struct) (*godef.Struct, []byte, error) {
	writer, structureStruct, writeStartErr := triggerWriteStructureStructStart(hooks, writer, structureStruct)
	if writeStartErr != nil {
		return nil, nil, triggerWriteStructureStructFailure(hooks, writer, structureStruct, writeStartErr)
	}

	structureStructContents, writeStructErr := writer.WriteStruct(structureStruct)
	if writeStructErr != nil {
		return nil, nil, triggerWriteStructureStructFailure(hooks, writer, structureStruct, writeStructErr)
	}

	structureStruct, structureStructContents, writeSuccessErr := triggerWriteStructureStructSuccess(hooks, structureStruct, structureStructContents)
	if writeSuccessErr != nil {
		return nil, nil, triggerWriteStructureStructFailure(hooks, writer, structureStruct, writeSuccessErr)
	}
	return structureStruct, structureStructContents, nil
}

func triggerWriteStructureStructStart(hooks hook.WriteGoStruct, writer write.GoStructWriter, structureStruct *godef.Struct) (write.GoStructWriter, *godef.Struct, error) {
	if hooks.OnWriteGoStructStart == nil {
		return writer, structureStruct, nil
	}
	if structureStruct == nil {
		return nil, nil, ErrNoStructureStruct
	}
	structureStructClone := structureStruct.DeepClone()

	updatedWriter, updatedStructureStruct, startErr := hooks.OnWriteGoStructStart(writer, &structureStructClone)
	if startErr != nil {
		return nil, nil, startErr
	}

	return updatedWriter, updatedStructureStruct, nil
}

func triggerWriteStructureStructSuccess(hooks hook.WriteGoStruct, structureStruct *godef.Struct, structureStructContents []byte) (*godef.Struct, []byte, error) {
	if hooks.OnWriteGoStructSuccess == nil {
		return structureStruct, structureStructContents, nil
	}
	if structureStruct == nil {
		return nil, nil, ErrNoStructureStruct
	}
	structureStructClone := structureStruct.DeepClone()
	structureStructContentsClone := clone.Slice(structureStructContents)

	updatedStructureStruct, updatedStructureStructContents, successErr := hooks.OnWriteGoStructSuccess(&structureStructClone, structureStructContentsClone)
	if successErr != nil {
		return nil, nil, successErr
	}
	return updatedStructureStruct, updatedStructureStructContents, nil
}

func triggerWriteStructureStructFailure(hooks hook.WriteGoStruct, writer write.GoStructWriter, structureStruct *godef.Struct, failureErr error) error {
	if hooks.OnWriteGoStructFailure == nil {
		return failureErr
	}

	structureStructClone := structureStruct.DeepClone()
	return hooks.OnWriteGoStructFailure(writer, &structureStructClone, failureErr)
}
