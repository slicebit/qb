package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDelete(t *testing.T) {
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

	statement = Delete(users).
		Where(Eq(users.C("id"), 5)).
		Build(sqlite)

	assert.Equal(t, statement.SQL(), "DELETE FROM users\nWHERE (users.id = ?);")
	assert.Equal(t, statement.Bindings(), []interface{}{5})

	statement = Delete(users).
		Where(Eq(users.C("id"), 5)).
		Build(mysql)

	assert.Equal(t, statement.SQL(), "DELETE FROM `users`\nWHERE (`users`.`id` = ?);")
	assert.Equal(t, statement.Bindings(), []interface{}{5})

	statement = Delete(users).
		Where(Eq(users.C("id"), 5)).
		Returning(users.C("id")).
		Build(postgres)

	assert.Equal(t, statement.SQL(), "DELETE FROM \"users\"\nWHERE (\"users\".\"id\" = $1)\nRETURNING \"id\";")
	assert.Equal(t, statement.Bindings(), []interface{}{5})

	statement = Delete(users).Build(sqlite)
	assert.Equal(t, statement.SQL(), "DELETE FROM users;")
}
