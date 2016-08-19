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

	assert.Equal(t, "UPDATE users\nSET email = ?;", statement.SQL())
	assert.Equal(t, []interface{}{"robert@de.niro"}, statement.Bindings())

	statement = Update(users).
		Values(map[string]interface{}{"email": "robert@de.niro"}).
		Build(mysql)

	assert.Equal(t, "UPDATE `users`\nSET `email` = ?;", statement.SQL())
	assert.Equal(t, []interface{}{"robert@de.niro"}, statement.Bindings())

	statement = Update(users).
		Values(map[string]interface{}{"email": "robert@de.niro"}).
		Where(Eq(users.C("email"), "al@pacino")).
		Returning(users.C("id"), users.C("email")).
		Build(postgres)

	assert.Equal(t, "UPDATE \"users\"\nSET \"email\" = $1\nWHERE \"users\".\"email\" = $2\nRETURNING \"id\", \"email\";", statement.SQL())
	assert.Equal(t, []interface{}{"robert@de.niro", "al@pacino"}, statement.Bindings())
}
