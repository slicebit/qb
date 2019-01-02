package qb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	users := Table(
		"users",
		Column("id", Varchar().Size(36)),
		Column("email", Varchar()).Unique(),
	)

	ins := Insert(users).Values(map[string]interface{}{
		"id":    "9883cf81-3b56-4151-ae4e-3903c5bc436d",
		"email": "al@pacino.com",
	})

	dialect := NewDefaultDialect()
	ctx := NewCompilerContext(dialect)

	sql := ins.Accept(ctx)
	binds := ctx.Binds()

	assert.Contains(t, sql, "INSERT INTO users")
	assert.Contains(t, sql, "id", "email")
	assert.Contains(t, sql, "VALUES(?, ?)")
	assert.Contains(t, binds, "9883cf81-3b56-4151-ae4e-3903c5bc436d")
	assert.Contains(t, binds, "al@pacino.com")

	sql = Insert(users).
		Values(map[string]interface{}{
			"id":    "9883cf81-3b56-4151-ae4e-3903c5bc436d",
			"email": "al@pacino.com",
		}).
		Returning(users.C("id"), users.C("email")).Accept(ctx)
	binds = ctx.Binds()

	assert.Contains(t, sql, "INSERT INTO users")
	assert.Contains(t, sql, "id", "email")
	assert.Contains(t, sql, "VALUES(?, ?)")
	assert.Contains(t, sql, "RETURNING id, email")
	assert.Contains(t, binds, "9883cf81-3b56-4151-ae4e-3903c5bc436d", "al@pacino.com")
}
