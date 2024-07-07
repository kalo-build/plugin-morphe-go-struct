package gofile

import (
	"go/format"
	"os"
	"path/filepath"

	"github.com/kaloseia/plugin-morphe-go-struct/pkg/strcase"
)

func WriteGoStructFile(dirPath string, structName string, structFileContents string) ([]byte, error) {
	formattedStructContents, formatErr := format.Source([]byte(structFileContents))
	if formatErr != nil {
		return nil, formatErr
	}

	structFileName := strcase.ToSnakeCaseLower(structName)
	structFilePath := filepath.Join(dirPath, structFileName+".go")
	if _, readErr := os.ReadDir(dirPath); readErr != nil && os.IsNotExist(readErr) {
		mkDirErr := os.MkdirAll(dirPath, 0644)
		if mkDirErr != nil {
			return nil, mkDirErr
		}
	}
	return formattedStructContents, os.WriteFile(structFilePath, formattedStructContents, 0644)
}
