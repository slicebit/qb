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
	metadata := MetaData(builder)

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
	err = qb.Metadata().CreateAll(qb.Engine())
	assert.Nil(t, err)

	qbNew, err := New("postgres", "user=postgres dbname=qb_test sslmode=disable")
	qbNew.Builder().SetEscaping(true)
	assert.Nil(t, err)
	qbNew.AddTable(&User{})
	err = qbNew.Metadata().CreateAll(qbNew.Engine())
	assert.NotNil(t, err)

	err = qb.Metadata().DropAll(qb.Engine())
	assert.Nil(t, err)

	err = qbNew.Metadata().DropAll(qbNew.Engine())
	assert.NotNil(t, err)
}

type UserMetadataError struct {
	ID int `qb:"constraints:i:nvalid"`
}

func TestMetadataAddError(t *testing.T) {
	builder := NewBuilder("postgres")
	metadata := MetaData(builder)

	assert.Panics(t, func() { metadata.Add(UserMetadataError{}) })
	assert.Equal(t, len(metadata.Tables()), 0)
}

func TestMetadataAddTable(t *testing.T) {
	builder := NewBuilder("postgres")
	metadata := MetaData(builder)

	table := Table("user", Column("id", BigInt()))

	metadata.AddTable(table)

	assert.Equal(t, metadata.Tables(), []TableElem{table})

	assert.Equal(t, metadata.Table("user").Name, "user")
}

func TestMetadataTable(t *testing.T) {
	builder := NewBuilder("postgres")
	metadata := MetaData(builder)

	assert.Nil(t, metadata.Table("invalid-table"))
}

func TestMetadataFailCreateDropAll(t *testing.T) {
	engine, _ := NewEngine("postgres", "user=postgres dbname=qb_test sslmode=disable")
	builder := NewBuilder(engine.Driver())
	metadata := MetaData(builder)

	var err error

	err = metadata.CreateAll(engine)
	assert.NotNil(t, err)

	err = metadata.DropAll(engine)
	assert.NotNil(t, err)
}
