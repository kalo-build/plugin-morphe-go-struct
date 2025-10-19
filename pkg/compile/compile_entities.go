package compile

import (
	"fmt"
	"strings"

	"github.com/kalo-build/go-util/core"

	"github.com/kalo-build/go/pkg/godef"
	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/morphe-go/pkg/yamlops"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/compile/cfg"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/compile/hook"
	"github.com/kalo-build/plugin-morphe-go-struct/pkg/typemap"
)

func AllMorpheEntitiesToGoStructs(config MorpheCompileConfig, r *registry.Registry) (map[string][]*godef.Struct, error) {
	allEntityStructDefs := map[string][]*godef.Struct{}
	for entityName, entity := range r.GetAllEntities() {
		entityStructs, entityStructsErr := MorpheEntityToGoStructs(config.EntityHooks, config.MorpheConfig, r, entity)
		if entityStructsErr != nil {
			return nil, entityStructsErr
		}
		allEntityStructDefs[entityName] = entityStructs
	}
	return allEntityStructDefs, nil
}

func MorpheEntityToGoStructs(entityHooks hook.CompileMorpheEntity, config cfg.MorpheConfig, r *registry.Registry, entity yaml.Entity) ([]*godef.Struct, error) {
	if r == nil {
		return nil, triggerCompileMorpheEntityFailure(entityHooks, config, entity, ErrNoRegistry)
	}

	config, entity, compileStartErr := triggerCompileMorpheEntityStart(entityHooks, config, entity)
	if compileStartErr != nil {
		return nil, triggerCompileMorpheEntityFailure(entityHooks, config, entity, compileStartErr)
	}

	allEntityStructs, structsErr := morpheEntityToGoStructs(config, r, entity)
	if structsErr != nil {
		return nil, triggerCompileMorpheEntityFailure(entityHooks, config, entity, structsErr)
	}

	allEntityStructs, compileSuccessErr := triggerCompileMorpheEntitySuccess(entityHooks, allEntityStructs)
	if compileSuccessErr != nil {
		return nil, triggerCompileMorpheEntityFailure(entityHooks, config, entity, compileSuccessErr)
	}

	return allEntityStructs, nil
}

func morpheEntityToGoStructs(config cfg.MorpheConfig, r *registry.Registry, entity yaml.Entity) ([]*godef.Struct, error) {
	validateConfigErr := config.Validate()
	if validateConfigErr != nil {
		return nil, validateConfigErr

	}
	validateMorpheErr := entity.Validate(r.GetAllEntities(), r.GetAllModels(), r.GetAllEnums())
	if validateMorpheErr != nil {
		return nil, validateMorpheErr
	}

	entityStruct, entityStructErr := getEntityStruct(config, r, entity)
	if entityStructErr != nil {
		return nil, entityStructErr
	}

	allIdentifierStructs, identifierStructsErr := getAllEntityIdentifierStructs(config.MorpheEntitiesConfig, entity, entityStruct)
	if identifierStructsErr != nil {
		return nil, identifierStructsErr
	}

	allEntityStructs := []*godef.Struct{
		entityStruct,
	}
	allEntityStructs = append(allEntityStructs, allIdentifierStructs...)
	return allEntityStructs, nil
}

func getEntityStruct(config cfg.MorpheConfig, r *registry.Registry, entity yaml.Entity) (*godef.Struct, error) {
	entityStruct := godef.Struct{
		Package: config.MorpheEntitiesConfig.Package,
		Name:    entity.Name,
	}

	structFields, fieldsErr := getGoFieldsForMorpheEntity(config, r, entity)
	if fieldsErr != nil {
		return nil, fieldsErr
	}
	entityStruct.Fields = structFields

	structImports, importsErr := getImportsForStructFields(config.MorpheEntitiesConfig.Package, structFields)
	if importsErr != nil {
		return nil, importsErr
	}
	entityStruct.Imports = structImports

	return &entityStruct, nil
}

func getGoFieldsForMorpheEntity(config cfg.MorpheConfig, r *registry.Registry, entity yaml.Entity) ([]godef.StructField, error) {
	allFields := []godef.StructField{}

	allFieldNames := core.MapKeysSorted(entity.Fields)
	// Handle direct fields
	for _, fieldName := range allFieldNames {
		entityField := entity.Fields[fieldName]
		fieldType, fieldErr := getModelFieldType(config, r, entityField.Type)
		if fieldErr != nil {
			return nil, fieldErr
		}

		field := godef.StructField{
			Name: fieldName,
			Type: fieldType,
		}
		allFields = append(allFields, field)
	}

	// Handle related entities
	relatedFields, relatedErr := getRelatedGoFieldsForMorpheEntity(config, r, entity.Related)
	if relatedErr != nil {
		return nil, relatedErr
	}
	allFields = append(allFields, relatedFields...)

	return allFields, nil
}

func getRelatedGoFieldsForMorpheEntity(config cfg.MorpheConfig, r *registry.Registry, entityRelations map[string]yaml.EntityRelation) ([]godef.StructField, error) {
	allFields := []godef.StructField{}

	allRelatedEntityNames := core.MapKeysSorted(entityRelations)
	for _, relationshipName := range allRelatedEntityNames {
		relation := entityRelations[relationshipName]

		// Resolve actual target entity name (handles aliasing)
		targetAlias := yamlops.GetRelationTargetName(relationshipName, relation.Aliased)

		// For HasMany relationships with path aliases (e.g., "Person.WorkProject"),
		// extract just the entity name (first part before the dot)
		targetEntityName := targetAlias
		if yamlops.IsRelationMany(relation.Type) && yamlops.IsRelationHas(relation.Type) {
			if parts := strings.Split(targetAlias, "."); len(parts) > 1 {
				targetEntityName = parts[0]
			}
		}

		// Get target entity
		targetEntity, entityErr := r.GetEntity(targetEntityName)
		if entityErr != nil {
			return nil, fmt.Errorf("failed to get target entity for relation %s: %w", relationshipName, entityErr)
		}

		idField, idErr := getRelatedGoFieldForEntityPrimaryID(config, r, relationshipName, targetEntity, relation.Type)
		if idErr != nil {
			return nil, idErr
		}
		allFields = append(allFields, idField)

		// Add entity reference field
		entityField, entityErr := getRelatedGoFieldForEntity(relationshipName, targetEntity, relation.Type)
		if entityErr != nil {
			return nil, entityErr
		}
		allFields = append(allFields, entityField)
	}

	return allFields, nil
}

func getRelatedGoFieldForEntityPrimaryID(config cfg.MorpheConfig, r *registry.Registry, relationName string, targetEntity yaml.Entity, relationType string) (godef.StructField, error) {
	primaryID, hasPrimary := targetEntity.Identifiers["primary"]
	if !hasPrimary {
		return godef.StructField{}, fmt.Errorf("related entity %s has no primary identifier", targetEntity.Name)
	}

	if len(primaryID.Fields) != 1 {
		return godef.StructField{}, fmt.Errorf("related entity %s primary identifier must have exactly one field", targetEntity.Name)
	}

	targetPrimaryIdName := primaryID.Fields[0]
	targetPrimaryIdField, primaryFieldExists := targetEntity.Fields[targetPrimaryIdName]
	if !primaryFieldExists {
		return godef.StructField{}, fmt.Errorf("related entity %s primary identifier field %s not found", targetEntity.Name, targetPrimaryIdName)
	}

	fieldType, fieldErr := getModelFieldType(config, r, targetPrimaryIdField.Type)
	if fieldErr != nil {
		return godef.StructField{}, fieldErr
	}

	fieldName := relationName + "ID"
	if yamlops.IsRelationMany(relationType) {
		fieldName += "s"
		return godef.StructField{
			Name: fieldName,
			Type: godef.GoTypeArray{
				IsSlice:   true,
				ValueType: fieldType,
			},
		}, nil
	}

	return godef.StructField{
		Name: fieldName,
		Type: godef.GoTypePointer{
			ValueType: fieldType,
		},
	}, nil
}

func getRelatedGoFieldForEntity(relationName string, targetEntity yaml.Entity, relationType string) (godef.StructField, error) {
	var fieldType godef.GoType
	fieldName := relationName

	if yamlops.IsRelationOne(relationType) {
		fieldType = godef.GoTypePointer{
			ValueType: godef.GoTypeStruct{
				Name: targetEntity.Name,
			},
		}
	} else if yamlops.IsRelationMany(relationType) {
		fieldName += "s"
		fieldType = godef.GoTypeArray{
			IsSlice: true,
			ValueType: godef.GoTypeStruct{
				Name: targetEntity.Name,
			},
		}
	} else {
		return godef.StructField{}, fmt.Errorf("unknown entity relation type: %s", relationType)
	}

	return godef.StructField{
		Name: fieldName,
		Type: fieldType,
	}, nil
}

func getModelFieldType(config cfg.MorpheConfig, r *registry.Registry, fieldType yaml.ModelFieldPath) (godef.GoType, error) {
	fieldPath := strings.Split(string(fieldType), ".")
	if len(fieldPath) < 2 {
		return nil, fmt.Errorf("invalid field type path: %s", fieldType)
	}

	// Get root model
	rootModelName := fieldPath[0]
	currentModel, modelErr := r.GetModel(rootModelName)
	if modelErr != nil {
		return nil, fmt.Errorf("morphe entity field %s references unknown root model: %s", fieldType, rootModelName)
	}

	// Traverse through related models
	for fieldIdx := 1; fieldIdx < len(fieldPath)-1; fieldIdx++ {
		relationshipName := fieldPath[fieldIdx]
		relationDef, exists := currentModel.Related[relationshipName]
		if !exists {
			return nil, fmt.Errorf("morphe entity field %s references unknown related model: %s", fieldType, relationshipName)
		}

		// Resolve actual target model name (handles aliasing)
		targetModelName := yamlops.GetRelationTargetName(relationshipName, relationDef.Aliased)

		relatedModel, relatedErr := r.GetModel(targetModelName)
		if relatedErr != nil {
			return nil, fmt.Errorf("morphe entity field %s references invalid related model: %s", fieldType, relationshipName)
		}
		currentModel = relatedModel
	}

	// Get terminal field
	terminalFieldName := fieldPath[len(fieldPath)-1]
	terminalField, exists := currentModel.Fields[terminalFieldName]
	if !exists {
		return nil, fmt.Errorf("morphe entity field %s references unknown model field: %s", fieldType, terminalFieldName)
	}

	goEnumField := getEnumFieldAsStructFieldType(
		config.MorpheEnumsConfig.Package,
		r.GetAllEnums(),
		terminalFieldName,
		string(terminalField.Type),
	)
	if goEnumField.Name != "" && goEnumField.Type != nil {
		return goEnumField.Type, nil
	}

	goFieldType, supported := typemap.MorpheModelFieldToGoField[terminalField.Type]
	if !supported {
		return nil, fmt.Errorf("morphe entity field %s has unsupported type: %s", fieldType, terminalField.Type)
	}

	return goFieldType, nil
}

func triggerCompileMorpheEntityStart(hooks hook.CompileMorpheEntity, config cfg.MorpheConfig, entity yaml.Entity) (cfg.MorpheConfig, yaml.Entity, error) {
	if hooks.OnCompileMorpheEntityStart == nil {
		return config, entity, nil
	}
	return hooks.OnCompileMorpheEntityStart(config, entity)
}

func triggerCompileMorpheEntitySuccess(hooks hook.CompileMorpheEntity, entityStructs []*godef.Struct) ([]*godef.Struct, error) {
	if hooks.OnCompileMorpheEntitySuccess == nil {
		return entityStructs, nil
	}
	return hooks.OnCompileMorpheEntitySuccess(entityStructs)
}

func triggerCompileMorpheEntityFailure(hooks hook.CompileMorpheEntity, config cfg.MorpheConfig, entity yaml.Entity, failureErr error) error {
	if hooks.OnCompileMorpheEntityFailure == nil {
		return failureErr
	}
	return hooks.OnCompileMorpheEntityFailure(config, entity, failureErr)
}

func getAllEntityIdentifierStructs(config cfg.MorpheEntitiesConfig, entity yaml.Entity, entityStruct *godef.Struct) ([]*godef.Struct, error) {
	return GetIdentifierStructs(
		config,
		entityStruct.Name,
		entityStruct,
		wrapEntityIdentifiers(entity.Identifiers),
	)
}

// Adapter to make EntityIdentifier implement Identifier interface
type entityIdentifierWrapper struct {
	yaml.EntityIdentifier
}

func (e entityIdentifierWrapper) GetFields() []string {
	return e.Fields
}

func wrapEntityIdentifiers(identifiers map[string]yaml.EntityIdentifier) map[string]Identifier {
	wrapped := make(map[string]Identifier)
	for k, v := range identifiers {
		wrapped[k] = entityIdentifierWrapper{v}
	}
	return wrapped
}
