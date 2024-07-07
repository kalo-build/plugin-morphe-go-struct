package compile

import "github.com/kaloseia/plugin-morphe-go-struct/pkg/godef"

// CompiledModelStructs maps Model.Name -> ModelStruct.Name -> CompiledStruct
type CompiledModelStructs map[string]map[string]CompiledStruct

func (structs CompiledModelStructs) AddCompiledModelStruct(modelName string, structDef *godef.Struct, structContents []byte) {
	if structs[modelName] == nil {
		structs[modelName] = make(map[string]CompiledStruct)
	}
	structs[modelName][structDef.Name] = CompiledStruct{
		Struct:         structDef,
		StructContents: structContents,
	}
}

func (structs CompiledModelStructs) GetAllCompiledModelStructs(modelName string) map[string]CompiledStruct {
	modelStructs, modelStructsExist := structs[modelName]
	if !modelStructsExist {
		return nil
	}
	return modelStructs
}

func (structs CompiledModelStructs) GetCompiledModelStruct(modelName string, structName string) CompiledStruct {
	modelStructs, modelStructsExist := structs[modelName]
	if !modelStructsExist {
		return CompiledStruct{}
	}
	compiledStruct, compiledStructExists := modelStructs[structName]
	if !compiledStructExists {
		return CompiledStruct{}
	}
	return compiledStruct
}
