package db

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Struct with lowercase id field
type lowerIDStruct struct {
	id uuid.UUID
}

// Struct with exported ID field
type upperIDStruct struct {
	ID uuid.UUID
}

// Struct without ID fields
type noIDStruct struct {
	Name string
}

func TestGetIdValue_LowercaseField(t *testing.T) {
	orm := MongoORM[lowerIDStruct]{}
	expected := uuid.New()
	doc := lowerIDStruct{id: expected}

	val, err := orm.getIdValue(&doc)
	assert.NoError(t, err)
	assert.Equal(t, expected, val)
}

func TestGetIdValue_UppercaseField(t *testing.T) {
	orm := MongoORM[upperIDStruct]{}
	expected := uuid.New()
	doc := upperIDStruct{ID: expected}

	val, err := orm.getIdValue(&doc)
	assert.NoError(t, err)
	assert.Equal(t, expected, val)
}

func TestGetIdValue_FieldMissing(t *testing.T) {
	orm := MongoORM[noIDStruct]{}
	doc := noIDStruct{Name: "test"}

	val, err := orm.getIdValue(&doc)
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, val)
}

func TestSetIdValue_LowercaseField(t *testing.T) {
	orm := MongoORM[lowerIDStruct]{}
	id := uuid.New()
	doc := lowerIDStruct{}

	err := orm.setIdValue(&doc, id)
	assert.NoError(t, err)
	assert.Equal(t, id, doc.id)
}

func TestSetIdValue_UppercaseField(t *testing.T) {
	orm := MongoORM[upperIDStruct]{}
	id := uuid.New()
	doc := upperIDStruct{}

	err := orm.setIdValue(&doc, id)
	assert.NoError(t, err)
	assert.Equal(t, id, doc.ID)
}

func TestSetIdValue_FieldMissing(t *testing.T) {
	orm := MongoORM[noIDStruct]{}
	id := uuid.New()
	doc := noIDStruct{}

	err := orm.setIdValue(&doc, id)
	assert.Error(t, err)
}
