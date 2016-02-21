package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type UserMetadata struct {
	ID int
}

func TestMetadata(t *testing.T) {

	engine, _ := NewEngine("postgres", "user=root dbname=pqtest")
	metadata := NewMetaData(engine)

	metadata.Add(UserMetadata{})
}

type UserMetadataError struct {
	ID int `qb:"constraints:i:nvalid"`
}

func TestMetadataAddError(t *testing.T) {

	engine, _ := NewEngine("postgres", "user=root dbname=pqtest")
	metadata := NewMetaData(engine)

	assert.Panics(t, func() { metadata.Add(UserMetadataError{}) })
	assert.Equal(t, len(metadata.Tables()), 0)
}

func TestMetadataAddTable(t *testing.T) {

	engine, _ := NewEngine("postgres", "user=root dbname=pqtest")
	metadata := NewMetaData(engine)

	table := NewTable(
		engine.Driver(),
		"user",
		[]Column{
			NewColumn("id", BigInt(), []Constraint{}),
		},
		[]Constraint{},
	)

	metadata.AddTable(table)

	assert.Equal(t, metadata.Tables(), []*Table{table})
}
