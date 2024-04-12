package godef

var (
	GoTypeByte = GoTypePrimitive{
		Syntax: "byte",
	}
	GoTypeRune = GoTypePrimitive{
		Syntax: "rune",
	}
	GoTypeString = GoTypePrimitive{
		Syntax: "string",
	}
	GoTypeInt = GoTypePrimitive{
		Syntax: "int",
	}
	GoTypeUint = GoTypePrimitive{
		Syntax: "uint",
	}
	GoTypeFloat = GoTypePrimitive{
		Syntax: "float",
	}
	GoTypeBool = GoTypePrimitive{
		Syntax: "bool",
	}
	GoTypeTime = GoTypeStruct{
		PackagePath: "time",
		Name:        "Time",
	}
)
