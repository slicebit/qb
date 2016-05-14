package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type UserMetadata struct {
	ID int
}

func TestMetadata(t *testing.T) {
	engine, _ := NewEngine("postgres", "user=postgres dbname=qb_test sslmode=disable")
	builder := NewBuilder(engine.Driver())
	metadata := NewMetaData(engine, builder)

	metadata.Add(UserMetadata{})
}

func TestMetadataCreateAllDropAllError(t *testing.T) {
	type User struct {
		ID string `qb:"type:uuid; constraints:primary_key"`
	}
	qb, err := New("postgres", "user=postgres dbname=qb_test sslmode=disable")
	qb.Builder().SetEscaping(true)
	assert.Nil(t, err)
	qb.Metadata().Add(&User{})
	err = qb.Metadata().CreateAll()
	assert.Nil(t, err)

	qbNew, err := New("postgres", "user=postgres dbname=qb_test sslmode=disable")
	qbNew.Builder().SetEscaping(true)
	assert.Nil(t, err)
	qbNew.Metadata().Add(&User{})
	err = qbNew.Metadata().CreateAll()
	assert.NotNil(t, err)

	err = qb.Metadata().DropAll()
	assert.Nil(t, err)

	err = qb.Metadata().DropAll()
	assert.NotNil(t, err)
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
