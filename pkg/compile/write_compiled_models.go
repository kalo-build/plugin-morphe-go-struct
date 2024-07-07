package compile

import (
	"github.com/kaloseia/clone"
	"github.com/kaloseia/go-util/core"
	"github.com/kaloseia/go/pkg/godef"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/hook"
	"github.com/kaloseia/plugin-morphe-go-struct/pkg/compile/write"
)

func WriteAllModelStructDefinitions(config MorpheCompileConfig, allModelStructDefs map[string][]*godef.Struct) (CompiledModelStructs, error) {
	allWrittenModels := CompiledModelStructs{}

	sortedModelNames := core.MapKeysSorted(allModelStructDefs)
	for _, modelName := range sortedModelNames {
		modelStructs := allModelStructDefs[modelName]
		for _, modelStruct := range modelStructs {
			modelStruct, modelStructContents, writeErr := WriteModelStructDefinition(config.WriteStructHooks, config.ModelWriter, modelStruct)
			if writeErr != nil {
				return nil, writeErr
			}
			allWrittenModels.AddCompiledModelStruct(modelName, modelStruct, modelStructContents)
		}
	}
	return allWrittenModels, nil
}

func WriteModelStructDefinition(hooks hook.WriteGoStruct, writer write.GoStructWriter, modelStruct *godef.Struct) (*godef.Struct, []byte, error) {
	writer, modelStruct, writeStartErr := triggerWriteModelStructStart(hooks, writer, modelStruct)
	if writeStartErr != nil {
		return nil, nil, triggerWriteModelStructFailure(hooks, writer, modelStruct, writeStartErr)
	}

	modelStructContents, writeStructErr := writer.WriteStruct(modelStruct)
	if writeStructErr != nil {
		return nil, nil, triggerWriteModelStructFailure(hooks, writer, modelStruct, writeStructErr)
	}

	modelStruct, modelStructContents, writeSuccessErr := triggerWriteModelStructSuccess(hooks, modelStruct, modelStructContents)
	if writeSuccessErr != nil {
		return nil, nil, triggerWriteModelStructFailure(hooks, writer, modelStruct, writeSuccessErr)
	}
	return modelStruct, modelStructContents, nil
}

func triggerWriteModelStructStart(hooks hook.WriteGoStruct, writer write.GoStructWriter, modelStruct *godef.Struct) (write.GoStructWriter, *godef.Struct, error) {
	if hooks.OnWriteGoStructStart == nil {
		return writer, modelStruct, nil
	}
	if modelStruct == nil {
		return nil, nil, ErrNoModelStruct
	}
	modelStructClone := modelStruct.DeepClone()

	updatedWriter, updatedModelStruct, startErr := hooks.OnWriteGoStructStart(writer, &modelStructClone)
	if startErr != nil {
		return nil, nil, startErr
	}

	return updatedWriter, updatedModelStruct, nil
}

func triggerWriteModelStructSuccess(hooks hook.WriteGoStruct, modelStruct *godef.Struct, modelStructContents []byte) (*godef.Struct, []byte, error) {
	if hooks.OnWriteGoStructSuccess == nil {
		return modelStruct, modelStructContents, nil
	}
	if modelStruct == nil {
		return nil, nil, ErrNoModelStruct
	}
	modelStructClone := modelStruct.DeepClone()
	modelStructContentsClone := clone.Slice(modelStructContents)

	updatedModelStruct, updatedModelStructContents, successErr := hooks.OnWriteGoStructSuccess(&modelStructClone, modelStructContentsClone)
	if successErr != nil {
		return nil, nil, successErr
	}
	return updatedModelStruct, updatedModelStructContents, nil
}

func triggerWriteModelStructFailure(hooks hook.WriteGoStruct, writer write.GoStructWriter, modelStruct *godef.Struct, failureErr error) error {
	if hooks.OnWriteGoStructFailure == nil {
		return failureErr
	}

	modelStructClone := modelStruct.DeepClone()
	return hooks.OnWriteGoStructFailure(writer, &modelStructClone, failureErr)
}
