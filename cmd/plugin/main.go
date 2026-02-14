package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kalo-build/plugin-morphe-go-struct/pkg/compile"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/compile/cfg"
)

type CompileConfigEntryStruct struct {
	PackagePath  string `json:"PackagePath"`
	ReceiverName string `json:"ReceiverName"`
}

type CompileConfigEntryEnum struct {
	PackagePath string `json:"PackagePath"`
}

type CompileConfigEntries struct {
	// FieldCasing applies to all sections (models, structures, entities).
	// Valid values: "camel", "snake", "pascal", or "" (none / no JSON tags).
	FieldCasing string `json:"fieldCasing,omitempty"`

	Models     CompileConfigEntryStruct `json:"models"`
	Enums      CompileConfigEntryEnum   `json:"enums"`
	Structures CompileConfigEntryStruct `json:"structures"`
	Entities   CompileConfigEntryStruct `json:"entities"`
}

type CompileConfig struct {
	InputPath  string               `json:"inputPath"`
	OutputPath string               `json:"outputPath"`
	Config     CompileConfigEntries `json:"config"`
	Verbose    bool                 `json:"verbose,omitempty"`
}

const (
	ErrMissingConfig       = 3
	ErrInvalidConfig       = 4
	ErrInputPathRequired   = 12
	ErrOutputPathRequired  = 13
	ErrPackagePathRequired = 14
	ErrCompileFailed       = 1
)

// logInfo prints info messages only when verbose mode is enabled
func logInfo(verbose bool, format string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stdout, format+"\n", args...)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: plugin-morphe-go-types <config>")
		fmt.Fprintln(os.Stderr, "  config: JSON string with inputPath, outputPath, and optional config parameters")
		os.Exit(ErrMissingConfig)
	}

	rawConfig := os.Args[1]
	var compileConfig CompileConfig
	if err := json.Unmarshal([]byte(rawConfig), &compileConfig); err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing config JSON:", err)
		fmt.Fprintln(os.Stderr, "Expected format: {\"inputPath\":\"...\",\"outputPath\":\"...\",\"config\":{...},\"verbose\":false}")
		os.Exit(ErrInvalidConfig)
	}

	if compileConfig.InputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: Input path is required")
		os.Exit(ErrInputPathRequired)
	}

	if compileConfig.OutputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: Output path is required")
		os.Exit(ErrOutputPathRequired)
	}

	if compileConfig.Config.Models.PackagePath == "" {
		fmt.Fprintln(os.Stderr, "Error: Models package path is required")
		os.Exit(ErrPackagePathRequired)
	}

	if compileConfig.Config.Enums.PackagePath == "" {
		fmt.Fprintln(os.Stderr, "Error: Enums package path is required")
		os.Exit(ErrPackagePathRequired)
	}

	if compileConfig.Config.Structures.PackagePath == "" {
		fmt.Fprintln(os.Stderr, "Error: Structures package path is required")
		os.Exit(ErrPackagePathRequired)
	}

	if compileConfig.Config.Entities.PackagePath == "" {
		fmt.Fprintln(os.Stderr, "Error: Entities package path is required")
		os.Exit(ErrPackagePathRequired)
	}

	inputAbs, err := filepath.Abs(compileConfig.InputPath)
	if err == nil {
		compileConfig.InputPath = inputAbs
	}

	outputAbs, err := filepath.Abs(compileConfig.OutputPath)
	if err == nil {
		compileConfig.OutputPath = outputAbs
	}

	logInfo(compileConfig.Verbose, "Processing Morphe registry from: '%s'", compileConfig.InputPath)
	logInfo(compileConfig.Verbose, "Output Go types to: '%s'", compileConfig.OutputPath)

	logInfo(compileConfig.Verbose, "Initializing compile config...")
	// Initialize the compile config with default values
	morpheConfig := compile.DefaultMorpheCompileConfig(
		compileConfig.InputPath,
		compileConfig.OutputPath,
	)

	logInfo(compileConfig.Verbose, "Setting package paths...")
	// Set the package paths (mandatory)
	morpheConfig.MorpheModelsConfig.Package.Path = compileConfig.Config.Models.PackagePath
	morpheConfig.MorpheEnumsConfig.Package.Path = compileConfig.Config.Enums.PackagePath
	morpheConfig.MorpheStructuresConfig.Package.Path = compileConfig.Config.Structures.PackagePath
	morpheConfig.MorpheEntitiesConfig.Package.Path = compileConfig.Config.Entities.PackagePath

	logInfo(compileConfig.Verbose, "Setting receiver names...")
	// Set the receiver names (optional)
	if compileConfig.Config.Models.ReceiverName != "" {
		morpheConfig.MorpheModelsConfig.ReceiverName = compileConfig.Config.Models.ReceiverName
	}
	if compileConfig.Config.Structures.ReceiverName != "" {
		morpheConfig.MorpheStructuresConfig.ReceiverName = compileConfig.Config.Structures.ReceiverName
	}
	if compileConfig.Config.Entities.ReceiverName != "" {
		morpheConfig.MorpheEntitiesConfig.ReceiverName = compileConfig.Config.Entities.ReceiverName
	}

	// Set field casing for JSON struct tags (applies to all sections)
	if compileConfig.Config.FieldCasing != "" {
		casing := cfg.Casing(compileConfig.Config.FieldCasing)
		logInfo(compileConfig.Verbose, "Setting field casing to: %s", casing)
		morpheConfig.MorpheModelsConfig.FieldCasing = casing
		morpheConfig.MorpheStructuresConfig.FieldCasing = casing
		morpheConfig.MorpheEntitiesConfig.FieldCasing = casing
	}

	logInfo(compileConfig.Verbose, "Starting compilation process...")
	compileErr := compile.MorpheToGo(morpheConfig)
	if compileErr != nil {
		fmt.Fprintln(os.Stderr, "Compilation failed:", compileErr)
		os.Exit(ErrCompileFailed)
	}

	logInfo(compileConfig.Verbose, "Compilation completed successfully")
	os.Exit(0)
}
