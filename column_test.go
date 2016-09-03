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

	col = Column("s", Varchar().Size(255)).Unique().NotNull().Default("hello")
	assert.Equal(t, "s VARCHAR(255) UNIQUE NOT NULL DEFAULT 'hello'", col.String(sqlite))

	precisionCol := Column("f", Type("FLOAT").Precision(2, 5)).Null()
	assert.Equal(t, "f FLOAT(2, 5) NULL", precisionCol.String(sqlite))

	col = Column("id", Int()).PrimaryKey().AutoIncrement()

	assert.Equal(t, "id INTEGER PRIMARY KEY", col.String(sqlite))
	assert.Equal(t, "`id` INT PRIMARY KEY AUTO_INCREMENT", col.String(mysql))
	assert.Equal(t, "\"id\" SERIAL PRIMARY KEY", col.String(postgres))

	assert.Equal(t, "c INT TEST", Column("c", Int()).Constraint("TEST").String(sqlite))

	var sql string

	sql = col.Accept(NewCompilerContext(sqlite))
	assert.Equal(t, "id", sql)

	sql = col.Accept(NewCompilerContext(mysql))
	assert.Equal(t, "`id`", sql)

	sql = col.Accept(NewCompilerContext(postgres))
	assert.Equal(t, "\"id\"", sql)
}
