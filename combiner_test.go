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

	email := Column("email", Varchar()).NotNull().Unique()
	id := Column("id", Int()).NotNull()

	and := And(Eq(email, "al@pacino.com"), NotEq(id, 1))
	or := Or(Eq(email, "al@pacino.com"), NotEq(id, 1))

	var sql string
	ctx := NewCompilerContext(sqlite)
	sql = and.Accept(ctx)

	assert.Equal(t, "(email = ? AND id != ?)", sql)
	assert.Equal(t, []interface{}{"al@pacino.com", 1}, ctx.Binds)

	ctx = NewCompilerContext(mysql)
	sql = and.Accept(ctx)

	assert.Equal(t, "(`email` = ? AND `id` != ?)", sql)
	assert.Equal(t, []interface{}{"al@pacino.com", 1}, ctx.Binds)

	ctx = NewCompilerContext(postgres)
	sql = and.Accept(ctx)

	assert.Equal(t, "(\"email\" = $1 AND \"id\" != $2)", sql)
	assert.Equal(t, []interface{}{"al@pacino.com", 1}, ctx.Binds)

	ctx = NewCompilerContext(sqlite)
	sql = or.Accept(ctx)

	assert.Equal(t, "(email = ? OR id != ?)", sql)
	assert.Equal(t, []interface{}{"al@pacino.com", 1}, ctx.Binds)

	ctx = NewCompilerContext(mysql)
	sql = or.Accept(ctx)

	assert.Equal(t, "(`email` = ? OR `id` != ?)", sql)
	assert.Equal(t, []interface{}{"al@pacino.com", 1}, ctx.Binds)

	ctx = NewCompilerContext(postgres)
	sql = or.Accept(ctx)

	assert.Equal(t, "(\"email\" = $1 OR \"id\" != $2)", sql)
	assert.Equal(t, []interface{}{"al@pacino.com", 1}, ctx.Binds)
}
