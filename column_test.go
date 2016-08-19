package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestColumn(t *testing.T) {
	sqlite := NewDialect("sqlite3")
	sqlite.SetEscaping(true)

	mysql := NewDialect("mysql")
	mysql.SetEscaping(true)

	postgres := NewDialect("postgres")
	postgres.SetEscaping(true)

	col := Column("id", Varchar().Size(40))
	assert.Equal(t, "id", col.Name)
	assert.Equal(t, Varchar().Size(40), col.Type)

	assert.Equal(t, "id VARCHAR(40)", col.String(sqlite))
	assert.Equal(t, "`id` VARCHAR(40)", col.String(mysql))
	assert.Equal(t, "\"id\" VARCHAR(40)", col.String(postgres))

	col = Column("id", Int()).PrimaryKey().AutoIncrement()

	assert.Equal(t, "id INTEGER PRIMARY KEY", col.String(sqlite))
	assert.Equal(t, "`id` INT PRIMARY KEY AUTO_INCREMENT", col.String(mysql))
	assert.Equal(t, "\"id\" SERIAL PRIMARY KEY", col.String(postgres))

	var sql string

	sql, _ = col.Build(sqlite)
	assert.Equal(t, "id", sql)

	sql, _ = col.Build(mysql)
	assert.Equal(t, "`id`", sql)

	sql, _ = col.Build(postgres)
	assert.Equal(t, "\"id\"", sql)
}
