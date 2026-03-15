package compile

import (
	"testing"

	"github.com/kalo-build/go/pkg/godef"
	"github.com/stretchr/testify/assert"
)

type mockIdentifier struct {
	fields []string
}

func (m mockIdentifier) GetFields() []string {
	return m.fields
}

type mockConfig struct {
	pkg          godef.Package
	receiverName string
}

func (m mockConfig) GetPackage() godef.Package {
	return m.pkg
}

func (m mockConfig) GetReceiverName() string {
	return m.receiverName
}

func TestGetIdentifierStructs(t *testing.T) {
	config := mockConfig{
		pkg: godef.Package{
			Path: "test/path",
			Name: "test",
		},
		receiverName: "m",
	}

	parentStruct := &godef.Struct{
		Name: "TestStruct",
		Fields: []godef.StructField{
			{
				Name: "ID",
				Type: godef.GoTypeInt,
			},
			{
				Name: "Name",
				Type: godef.GoTypeString,
			},
		},
	}

	identifiers := map[string]Identifier{
		"primary": mockIdentifier{
			fields: []string{"ID"},
		},
		"byName": mockIdentifier{
			fields: []string{"Name"},
		},
	}

	// Run test
	identifierStructs, err := GetIdentifierStructs(
		config,
		parentStruct.Name,
		parentStruct,
		identifiers,
	)

	assert.Nil(t, err)
	assert.Len(t, identifierStructs, 2)

	identifierStruct0 := identifierStructs[0]
	assert.Equal(t, "TestStructIDByName", identifierStruct0.Name)
	assert.Len(t, identifierStruct0.Fields, 1)
	assert.Equal(t, "Name", identifierStruct0.Fields[0].Name)

	identifierStruct1 := identifierStructs[1]
	assert.Equal(t, "TestStructIDPrimary", identifierStruct1.Name)
	assert.Len(t, identifierStruct1.Fields, 1)
	assert.Equal(t, "ID", identifierStruct1.Fields[0].Name)

	assert.Len(t, parentStruct.Methods, 2)
	assert.Equal(t, "GetIDByName", parentStruct.Methods[0].Name)
	assert.Equal(t, "GetIDPrimary", parentStruct.Methods[1].Name)
}

func TestGetIdentifierStructs_RelPrefixedFields(t *testing.T) {
	config := mockConfig{
		pkg: godef.Package{
			Path: "test/path",
			Name: "test",
		},
		receiverName: "m",
	}

	parentStruct := &godef.Struct{
		Name: "TaskTag",
		Fields: []godef.StructField{
			{Name: "ID", Type: godef.GoTypeInt},
			{Name: "TaskID", Type: godef.GoTypeInt},
			{Name: "TagID", Type: godef.GoTypeInt},
		},
	}

	identifiers := map[string]Identifier{
		"primary": mockIdentifier{
			fields: []string{"ID"},
		},
		"taskTag": mockIdentifier{
			fields: []string{"rel:Task", "rel:Tag"},
		},
	}

	identifierStructs, err := GetIdentifierStructs(config, parentStruct.Name, parentStruct, identifiers)

	assert.Nil(t, err)
	assert.Len(t, identifierStructs, 2)

	identifierStruct0 := identifierStructs[0]
	assert.Equal(t, "TaskTagIDPrimary", identifierStruct0.Name)
	assert.Len(t, identifierStruct0.Fields, 1)
	assert.Equal(t, "ID", identifierStruct0.Fields[0].Name)

	identifierStruct1 := identifierStructs[1]
	assert.Equal(t, "TaskTagIDTaskTag", identifierStruct1.Name)
	assert.Len(t, identifierStruct1.Fields, 2)
	assert.Equal(t, "TaskID", identifierStruct1.Fields[0].Name)
	assert.Equal(t, "TagID", identifierStruct1.Fields[1].Name)
}

func TestGetIdentifierStructs_RelPrefixedPolyFields(t *testing.T) {
	config := mockConfig{
		pkg: godef.Package{
			Path: "test/path",
			Name: "test",
		},
		receiverName: "m",
	}

	parentStruct := &godef.Struct{
		Name: "Comment",
		Fields: []godef.StructField{
			{Name: "ID", Type: godef.GoTypeInt},
			{Name: "CommentableID", Type: godef.GoTypeInt},
			{Name: "CommentableType", Type: godef.GoTypeString},
		},
	}

	identifiers := map[string]Identifier{
		"primary": mockIdentifier{
			fields: []string{"ID"},
		},
		"commentable": mockIdentifier{
			fields: []string{"rel:Commentable"},
		},
	}

	identifierStructs, err := GetIdentifierStructs(config, parentStruct.Name, parentStruct, identifiers)

	assert.Nil(t, err)
	assert.Len(t, identifierStructs, 2)

	commentableStruct := identifierStructs[0]
	assert.Equal(t, "CommentIDCommentable", commentableStruct.Name)
	assert.Len(t, commentableStruct.Fields, 2)
	assert.Equal(t, "CommentableID", commentableStruct.Fields[0].Name)
	assert.Equal(t, "CommentableType", commentableStruct.Fields[1].Name)
}
