package qb

import (
	"io/ioutil"
	"os"

	"testing"

	"github.com/aacanakin/qb"
	"github.com/aacanakin/qb/dialects/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestMetadataCreateAllDropAllError(t *testing.T) {

	qb.RegisterDialect("sqlite3", sqlite.NewDialect())

	tmpFile, err := ioutil.TempFile("", "qbtestdb")
	if err != nil {
		t.Fatalf("Cannot create a temporary file. Got '%s'", err)
	}
	defer os.Remove(tmpFile.Name())
	dsn := "file://" + tmpFile.Name()
	tmpFile.Close()

	accounts := Table(
		"account",
		Column("id", Type("UUID")),
		PrimaryKey("id"),
	)
	engine, err := New("sqlite3", dsn)
	metadata := MetaData()

	engine.Dialect().SetEscaping(true)
	assert.Nil(t, err)
	metadata.AddTable(accounts)
	err = metadata.CreateAll(engine)
	assert.Nil(t, err)

	engineNew, err := New("sqlite3", dsn)
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
	engine, _ := New("sqlite3", ":memory:")
	metadata := MetaData()

	var err error

	err = metadata.CreateAll(engine)
	assert.NotNil(t, err)

	err = metadata.DropAll(engine)
	assert.NotNil(t, err)
}

func TestMetadataWithNoConnection(t *testing.T) {
	engine, _ := New("sqlite3", ":memory:")
	engine.DB().Close()

	metadata := MetaData()
	assert.NotNil(t, metadata.CreateAll(engine))
	assert.NotNil(t, metadata.DropAll(engine))
}
