package gofile

import (
	"go/format"
	"os"
	"path/filepath"

	"github.com/kalo-build/go-util/strcase"
)

func WriteGoDefinitionFile(dirPath string, definitionName string, goFileContents string) ([]byte, error) {
	formattedStructContents, formatErr := format.Source([]byte(goFileContents))
	if formatErr != nil {
		return nil, formatErr
	}

	definitionFileName := strcase.ToSnakeCaseLower(definitionName)
	definitionFilePath := filepath.Join(dirPath, definitionFileName+".go")
	if _, readErr := os.ReadDir(dirPath); readErr != nil && os.IsNotExist(readErr) {
		mkDirErr := os.MkdirAll(dirPath, 0644)
		if mkDirErr != nil {
			return nil, mkDirErr
		}
	}
	return formattedStructContents, os.WriteFile(definitionFilePath, formattedStructContents, 0644)
}
