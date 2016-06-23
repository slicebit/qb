package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdate(t *testing.T) {
	sqlite := NewDialect("sqlite3")
	sqlite.SetEscaping(true)

	mysql := NewDialect("mysql")
	mysql.SetEscaping(true)

	postgres := NewDialect("postgres")
	postgres.SetEscaping(true)

	users := Table(
		"users",
		Column("id", BigInt().NotNull()),
		Column("email", Varchar().NotNull().Unique()),
		PrimaryKey("email"),
	)

	var statement *Stmt

	statement = Update(users).
		Values(map[string]interface{}{"email": "robert@de.niro"}).
		Build(sqlite)

	assert.Equal(t, statement.SQL(), "UPDATE users\nSET email = ?;")
	assert.Equal(t, statement.Bindings(), []interface{}{"robert@de.niro"})

	statement = Update(users).
		Values(map[string]interface{}{"email": "robert@de.niro"}).
		Build(mysql)

	assert.Equal(t, statement.SQL(), "UPDATE `users`\nSET `email` = ?;")
	assert.Equal(t, statement.Bindings(), []interface{}{"robert@de.niro"})

	statement = Update(users).
		Values(map[string]interface{}{"email": "robert@de.niro"}).
		Where(Eq(users.C("email"), "al@pacino")).
		Returning(users.C("id"), users.C("email")).
		Build(postgres)

	assert.Equal(t, statement.SQL(), "UPDATE \"users\"\nSET \"email\" = $1\nWHERE (\"users\".\"email\" = $2)\nRETURNING \"id\", \"email\";")
	assert.Equal(t, statement.Bindings(), []interface{}{"robert@de.niro", "al@pacino"})
}
