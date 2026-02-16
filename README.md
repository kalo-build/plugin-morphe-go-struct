# plugin-morphe-go-struct

Generates Go struct definitions from Morphe schema definitions (`KA:MO1:YAML1`). Produces typed Go models, entities, enums, and structures with relationship handling, identifier accessors, and configurable JSON tags.

## What it generates

| Morphe artifact | Go output                                                                 |
|-----------------|---------------------------------------------------------------------------|
| **Model**       | Struct with fields, relationship fields (IDs + pointers/slices), identifier getter methods |
| **Entity**      | Struct with resolved fields, `morphe:` attribute tags, identifier getter methods |
| **Enum**        | Named `string` type with typed constants                                   |
| **Structure**   | Plain struct with typed fields                                             |

### Example output

**Model** (`person.go`):

```go
package models

import "github.com/myapp/enums"

type Person struct {
    FirstName   string
    ID          uint
    LastName    string
    Nationality enums.Nationality
    CompanyID   *uint
    Company     *Company
}

func (m Person) GetIDPrimary() PersonIDPrimary {
    return PersonIDPrimary{ID: m.ID}
}

func (m Person) GetIDName() PersonIDName {
    return PersonIDName{FirstName: m.FirstName, LastName: m.LastName}
}
```

**Enum** (`nationality.go`):

```go
package enums

type Nationality string

const (
    NationalityDE Nationality = "German"
    NationalityFR Nationality = "French"
    NationalityUS Nationality = "American"
)
```

**Structure** (`address.go`):

```go
package structures

type Address struct {
    City    string
    HouseNr string
    Street  string
    ZipCode string
}
```

### Relationship handling

| Relationship type | Generated fields                              |
|-------------------|-----------------------------------------------|
| `ForOne`          | `{Rel}ID *uint` + `{Rel} *{Target}`            |
| `ForMany`         | `{Rel}IDs []uint` + `{Rel} []{Target}`         |
| `HasOne`          | `{Rel}ID *uint` + `{Rel} *{Target}`            |
| `HasMany`         | `{Rel}IDs []uint` + `{Rel} []{Target}`         |
| Polymorphic       | `{Rel}Type *string` + `{Rel}ID *uint` + pointer |

### Type mappings

| Morphe type     | Go type     |
|-----------------|-------------|
| `UUID`          | `string`    |
| `AutoIncrement` | `uint`      |
| `String`        | `string`    |
| `Integer`       | `int`       |
| `Float`         | `float64`   |
| `Boolean`       | `bool`      |
| `Time`          | `time.Time` |
| `Date`          | `time.Time` |
| `Protected`     | `string`    |
| `Sealed`        | `string`    |

## Input / output

| Direction | Format         | Store suggestion | Description                        |
|-----------|----------------|------------------|------------------------------------|
| Input     | `KA:MO1:YAML1` | `KA_MO_YAML`   | Morphe registry (models, enums, structures, entities) |
| Output    | `KA:MO1:GO1`   | `KA_GO_TYPES`  | Go source files (structs, enums)   |

## Configuration

| Key                       | Type   | Required | Default | Description                                    |
|---------------------------|--------|----------|---------|------------------------------------------------|
| `config.fieldCasing`      | string | no       | `""`    | JSON struct tag casing: `"camel"`, `"snake"`, `"pascal"`, or `""` (no tags) |
| `config.models.PackagePath`     | string | yes | —       | Go import path for the generated models package |
| `config.enums.PackagePath`      | string | yes | —       | Go import path for the generated enums package  |
| `config.structures.PackagePath` | string | yes | —       | Go import path for the generated structures package |
| `config.entities.PackagePath`   | string | yes | —       | Go import path for the generated entities package |

## Pipeline context

```yaml
stores:
  KA_MO_YAML:
    format: "KA:MO1:YAML1"
    type: "localFileSystem"
    options:
      path: "./morphe"

  KA_GO_TYPES:
    format: "KA:MO1:GO1"
    type: "localFileSystem"
    options:
      path: "./types"

plugins:
  "@kalo-build/plugin-morphe-go-struct":
    version: "v1.0.0"
    inputs:
      morphe:
        format: "KA:MO1:YAML1"
        store: "KA_MO_YAML"
    output:
      format: "KA:MO1:GO1"
      store: "KA_GO_TYPES"
    config:
      fieldCasing: "camel"
      models:
        PackagePath: "github.com/myapp/internal/types/models"
      enums:
        PackagePath: "github.com/myapp/internal/types/enums"
      structures:
        PackagePath: "github.com/myapp/internal/types/structures"
      entities:
        PackagePath: "github.com/myapp/internal/types/entities"
```

## Project structure

```
plugin-morphe-go-struct/
├── cmd/plugin/             # WASM entry point
├── pkg/
│   ├── compile/            # Compilation pipeline
│   │   ├── compile.go      # MorpheToGo entry point
│   │   ├── compile_models.go
│   │   ├── compile_entities.go
│   │   ├── compile_structures.go
│   │   ├── compile_enums.go
│   │   ├── identifier_structs.go  # Identifier struct + getter generation
│   │   ├── cfg/            # Configuration structs and casing
│   │   ├── hook/           # Extensibility hooks
│   │   └── write/          # File writers
│   └── typemap/            # Morphe → Go type mappings
├── testdata/
│   ├── registry/           # Sample Morphe registry input
│   └── ground-truth/       # Expected Go output for integration tests
├── dist/                   # WASM output
└── plugin.yaml             # Kalo plugin manifest
```

## Building

```bash
# Native binary
go build ./cmd/plugin

# WASM (for Kalo CLI)
GOOS=wasip1 GOARCH=wasm go build -o dist/plugin.wasm cmd/plugin/main.go
```

## Testing

```bash
go test ./...
```
