package compile

import "github.com/kaloseia/go/pkg/godef"

// CompiledMorpheStructs maps Morphe.Name -> MorpheStruct.Name -> CompiledStruct
type CompiledMorpheStructs map[string]map[string]CompiledStruct

func (structs CompiledMorpheStructs) AddCompiledMorpheStruct(morpheName string, structDef *godef.Struct, structContents []byte) {
	if structs[morpheName] == nil {
		structs[morpheName] = make(map[string]CompiledStruct)
	}
	structs[morpheName][structDef.Name] = CompiledStruct{
		Struct:         structDef,
		StructContents: structContents,
	}
}

func (structs CompiledMorpheStructs) GetAllCompiledMorpheStructs(morpheName string) map[string]CompiledStruct {
	morpheStructs, morpheStructsExist := structs[morpheName]
	if !morpheStructsExist {
		return nil
	}
	return morpheStructs
}

func (structs CompiledMorpheStructs) GetCompiledMorpheStruct(morpheName string, structName string) CompiledStruct {
	morpheStructs, morpheStructsExist := structs[morpheName]
	if !morpheStructsExist {
		return CompiledStruct{}
	}
	compiledStruct, compiledStructExists := morpheStructs[structName]
	if !compiledStructExists {
		return CompiledStruct{}
	}
	return compiledStruct
}
