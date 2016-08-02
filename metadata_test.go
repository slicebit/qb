package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetadataCreateAllDropAllError(t *testing.T) {
	type Account struct {
		ID string `qb:"type:uuid; constraints:primary_key"`
	}
	qb, err := New("postgres", postgresDsn)

	qb.Dialect().SetEscaping(true)
	assert.Nil(t, err)
	qb.AddTable(&Account{})
	err = qb.Metadata().CreateAll(qb.Engine())
	assert.Nil(t, err)

	qbNew, err := New("postgres", postgresDsn)
	qbNew.Dialect().SetEscaping(true)
	assert.Nil(t, err)
	qbNew.AddTable(&Account{})
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
	dialect := NewDialect("postgres")
	metadata := MetaData(dialect)

	assert.Panics(t, func() { metadata.Add(UserMetadataError{}) })
	assert.Equal(t, len(metadata.Tables()), 0)
}

func TestMetadataAddTable(t *testing.T) {
	dialect := NewDialect("postgres")
	metadata := MetaData(dialect)

	table := Table("user", Column("id", BigInt()))

	metadata.AddTable(table)

	assert.Equal(t, metadata.Tables(), []TableElem{table})

	assert.Equal(t, metadata.Table("user").Name, "user")
}

func TestMetadataTable(t *testing.T) {
	dialect := NewDialect("postgres")
	metadata := MetaData(dialect)

	assert.Panics(t, func() { metadata.Table("invalid-table") })
}

func TestMetadataFailCreateDropAll(t *testing.T) {
	engine, _ := NewEngine("postgres", postgresDsn)
	dialect := NewDialect(engine.Driver())
	metadata := MetaData(dialect)

	var err error

	err = metadata.CreateAll(engine)
	assert.NotNil(t, err)

	err = metadata.DropAll(engine)
	assert.NotNil(t, err)
}

func TestMetadataWithNoConnection(t *testing.T) {
	engine, _ := NewEngine("postgres", postgresDsn)
	engine.DB().Close()

	metadata := MetaData(NewDialect(engine.Driver()))
	assert.NotNil(t, metadata.CreateAll(engine))
	assert.NotNil(t, metadata.DropAll(engine))
}
