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
	builder := NewBuilder(engine.Driver())
	metadata := NewMetaData(engine, builder)

	metadata.Add(UserMetadata{})
}

type UserMetadataError struct {
	ID int `qb:"constraints:i:nvalid"`
}

func TestMetadataAddError(t *testing.T) {

	engine, _ := NewEngine("postgres", "user=root dbname=pqtest")
	builder := NewBuilder("postgres")
	metadata := NewMetaData(engine, builder)

	assert.Panics(t, func() { metadata.Add(UserMetadataError{}) })
	assert.Equal(t, len(metadata.Tables()), 0)
}

func TestMetadataAddTable(t *testing.T) {

	engine, _ := NewEngine("postgres", "user=root dbname=pqtest")
	builder := NewBuilder("postgres")
	metadata := NewMetaData(engine, builder)

	table := NewTable(
		builder,
		"user",
		[]Column{
			NewColumn("id", NewType("BIGINT"), []Constraint{}),
		},
	)

	metadata.AddTable(table)

	assert.Equal(t, metadata.Tables(), []*Table{table})

	assert.Equal(t, metadata.Table("user").Name(), "user")
}

func TestMetadataTable(t *testing.T) {
	engine, _ := NewEngine("postgres", "user=root dbname=pqtest")
	builder := NewBuilder("postgres")
	metadata := NewMetaData(engine, builder)

	assert.Nil(t, metadata.Table("invalid-table"))
}

func TestMetadataFailCreateDropAll(t *testing.T) {
	engine, _ := NewEngine("postgres", "user=postgres dbname=qb_test")
	builder := NewBuilder("postgres")
	metadata := NewMetaData(engine, builder)
	metadata.Engine().DB().Close()
	assert.NotNil(t, metadata.CreateAll())
	assert.NotNil(t, metadata.DropAll())
}
