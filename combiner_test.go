package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCombiners(t *testing.T) {
	sqlite := NewDialect("sqlite3")
	sqlite.SetEscaping(true)

	mysql := NewDialect("mysql")
	mysql.SetEscaping(true)

	postgres := NewDialect("postgres")
	postgres.SetEscaping(true)

	email := Column("email", Varchar().NotNull().Unique())
	id := Column("id", Int().NotNull())

	and := And(Eq(email, "al@pacino.com"), NotEq(id, 1))
	or := Or(Eq(email, "al@pacino.com"), NotEq(id, 1))

	var sql string
	var bindings []interface{}
	sql, bindings = and.Build(sqlite)

	assert.Equal(t, "(email = ? AND id != ?)", sql)
	assert.Equal(t, []interface{}{"al@pacino.com", 1}, bindings)

	sql, bindings = and.Build(mysql)

	assert.Equal(t, "(`email` = ? AND `id` != ?)", sql)
	assert.Equal(t, []interface{}{"al@pacino.com", 1}, bindings)

	sql, bindings = and.Build(postgres)

	assert.Equal(t, "(\"email\" = $1 AND \"id\" != $2)", sql)
	assert.Equal(t, []interface{}{"al@pacino.com", 1}, bindings)

	sql, bindings = or.Build(sqlite)

	assert.Equal(t, "(email = ? OR id != ?)", sql)
	assert.Equal(t, []interface{}{"al@pacino.com", 1}, bindings)

	sql, bindings = or.Build(mysql)

	assert.Equal(t, "(`email` = ? OR `id` != ?)", sql)
	assert.Equal(t, []interface{}{"al@pacino.com", 1}, bindings)

	sql, bindings = or.Build(postgres)

	assert.Equal(t, "(\"email\" = $3 OR \"id\" != $4)", sql)
	assert.Equal(t, []interface{}{"al@pacino.com", 1}, bindings)
}
