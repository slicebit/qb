package qb

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestColumn(t *testing.T) {
	sqlite := NewDialect("sqlite3")
	sqlite.SetEscaping(true)

	mysql := NewDialect("mysql")
	mysql.SetEscaping(true)

	postgres := NewDialect("postgres")
	postgres.SetEscaping(true)

	col := Column("id", Varchar().Size(40))
	assert.Equal(t, col.Name, "id")
	assert.Equal(t, col.Type, Varchar().Size(40))

	assert.Equal(t, col.String(sqlite), "id VARCHAR(40)")
	assert.Equal(t, col.String(mysql), "`id` VARCHAR(40)")
	assert.Equal(t, col.String(postgres), "\"id\" VARCHAR(40)")

	var sql string

	sql, _ = col.Build(sqlite)
	assert.Equal(t, sql, "id")

	sql, _ = col.Build(mysql)
	assert.Equal(t, sql, "`id`")

	sql, _ = col.Build(postgres)
	assert.Equal(t, sql, "\"id\"")
}
