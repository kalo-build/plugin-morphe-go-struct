package compile

import (
	"fmt"
	"time"

	"github.com/kalo-build/go-util/core"
	"github.com/kalo-build/go-util/strcase"
	"github.com/kalo-build/go/pkg/godef"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/gofile"
)

type MorpheEnumFileWriter struct {
	TargetDirPath string
}

func (w *MorpheEnumFileWriter) WriteEnum(enumDefinition *godef.Enum) ([]byte, error) {
	allEnumLines, allLinesErr := w.getAllEnumLines(enumDefinition)
	if allLinesErr != nil {
		return nil, allLinesErr
	}

	enumFileContents, enumContentsErr := core.LinesToString(allEnumLines)
	if enumContentsErr != nil {
		return nil, enumContentsErr
	}

	return gofile.WriteGoDefinitionFile(w.TargetDirPath, enumDefinition.Name, enumFileContents)
}

func (w *MorpheEnumFileWriter) getAllEnumLines(enumDefinition *godef.Enum) ([]string, error) {
	allEnumLines := []string{
		fmt.Sprintf("package %s", enumDefinition.Package.Name),
		fmt.Sprintf("type %s %s", enumDefinition.Name, enumDefinition.Type.BaseType.GetSyntaxLocal()),
		"const (",
	}

	for _, enumEntry := range enumDefinition.Entries {
		entryName := strcase.ToPascalCase(enumEntry.Name)
		entryValue := w.formatEnumValue(enumEntry.Value)
		enumEntryLine := fmt.Sprintf("\t%s %s = %v",
			entryName, enumDefinition.Name, entryValue)
		allEnumLines = append(allEnumLines, enumEntryLine)
	}

	allEnumLines = append(allEnumLines, ")")
	return allEnumLines, nil
}

func (w *MorpheEnumFileWriter) formatEnumValue(value any) string {
	switch typedValue := value.(type) {
	case string:
		return fmt.Sprintf("%q", typedValue)
	case time.Time:
		formattedValue := ""
		if typedValue.Hour() == 0 && typedValue.Minute() == 0 && typedValue.Second() == 0 && typedValue.Nanosecond() == 0 {
			formattedValue = typedValue.Format("2006-01-02")
		} else {
			formattedValue = typedValue.Format(time.RFC3339)
		}
		return fmt.Sprintf("%q", formattedValue)
	default:
		return fmt.Sprintf("%v", typedValue)
	}
}
