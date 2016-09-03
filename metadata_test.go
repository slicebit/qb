package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetadataCreateAllDropAllError(t *testing.T) {
	accounts := Table(
		"account",
		Column("id", Type("UUID")),
		PrimaryKey("id"),
	)
	engine, err := New("postgres", postgresDsn)
	metadata := MetaData()

	engine.Dialect().SetEscaping(true)
	assert.Nil(t, err)
	metadata.AddTable(accounts)
	err = metadata.CreateAll(engine)
	assert.Nil(t, err)

	engineNew, err := New("postgres", postgresDsn)
	engineNew.Dialect().SetEscaping(true)
	metadataNew := MetaData()
	assert.Nil(t, err)
	metadataNew.AddTable(accounts)
	err = metadataNew.CreateAll(engineNew)
	assert.NotNil(t, err)

	err = metadataNew.DropAll(engine)
	assert.Nil(t, err)

	err = metadataNew.DropAll(engineNew)
	assert.NotNil(t, err)
}

func TestMetadataAddError(t *testing.T) {
	metadata := MetaData()
	assert.Equal(t, len(metadata.Tables()), 0)
}

func TestMetadataAddTable(t *testing.T) {
	metadata := MetaData()

	table := Table("user", Column("id", BigInt()))

	metadata.AddTable(table)

	assert.Equal(t, []TableElem{table}, metadata.Tables())

	assert.Equal(t, "user", metadata.Table("user").Name)
}

func TestMetadataTable(t *testing.T) {
	metadata := MetaData()

	assert.Panics(t, func() { metadata.Table("invalid-table") })
}

func TestMetadataFailCreateDropAll(t *testing.T) {
	engine, _ := New("postgres", postgresDsn)
	metadata := MetaData()

	var err error

	err = metadata.CreateAll(engine)
	assert.NotNil(t, err)

	err = metadata.DropAll(engine)
	assert.NotNil(t, err)
}

func TestMetadataWithNoConnection(t *testing.T) {
	engine, _ := New("postgres", postgresDsn)
	engine.DB().Close()

	metadata := MetaData()
	assert.NotNil(t, metadata.CreateAll(engine))
	assert.NotNil(t, metadata.DropAll(engine))
}
