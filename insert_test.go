package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsert(t *testing.T) {
	sqlite := NewDialect("sqlite3")
	sqlite.SetEscaping(true)

	mysql := NewDialect("mysql")
	mysql.SetEscaping(true)

	postgres := NewDialect("postgres")
	postgres.SetEscaping(true)

	users := Table(
		"users",
		Column("id", Varchar().Size(36)),
		Column("email", Varchar().Unique()),
	)

	var statement *Stmt

	ins := Insert(users).Values(map[string]interface{}{
		"id":    "9883cf81-3b56-4151-ae4e-3903c5bc436d",
		"email": "al@pacino.com",
	})

	statement = ins.Build(sqlite)
	assert.Contains(t, statement.SQL(), "INSERT INTO users")
	assert.Contains(t, statement.SQL(), "id", "email")
	assert.Contains(t, statement.SQL(), "VALUES(?, ?)")
	assert.Contains(t, statement.Bindings(), "9883cf81-3b56-4151-ae4e-3903c5bc436d")
	assert.Contains(t, statement.Bindings(), "al@pacino.com")

	statement = ins.Build(mysql)
	assert.Contains(t, statement.SQL(), "INSERT INTO `users`")
	assert.Contains(t, statement.SQL(), "`id`", "`email`")
	assert.Contains(t, statement.SQL(), "VALUES(?, ?)")
	assert.Contains(t, statement.Bindings(), "9883cf81-3b56-4151-ae4e-3903c5bc436d")
	assert.Contains(t, statement.Bindings(), "al@pacino.com")

	statement = ins.Build(postgres)
	assert.Contains(t, statement.SQL(), "INSERT INTO \"users\"")
	assert.Contains(t, statement.SQL(), "\"id\"", "\"email\"")
	assert.Contains(t, statement.SQL(), "VALUES($1, $2)")
	assert.Contains(t, statement.Bindings(), "9883cf81-3b56-4151-ae4e-3903c5bc436d")
	assert.Contains(t, statement.Bindings(), "al@pacino.com")

	statement = Insert(users).
		Values(map[string]interface{}{
			"id":    "9883cf81-3b56-4151-ae4e-3903c5bc436d",
			"email": "al@pacino.com",
		}).
		Returning(users.C("id"), users.C("email")).
		Build(postgres)

	assert.Contains(t, statement.SQL(), "INSERT INTO \"users\"")
	assert.Contains(t, statement.SQL(), "\"id\"", "\"email\"")
	assert.Contains(t, statement.SQL(), "VALUES($1, $2)")
	assert.Contains(t, statement.SQL(), "RETURNING \"id\", \"email\";")
	assert.Contains(t, statement.Bindings(), "9883cf81-3b56-4151-ae4e-3903c5bc436d", "al@pacino.com")
}
