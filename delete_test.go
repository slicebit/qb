package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDelete(t *testing.T) {
	dialect := NewDialect("default")

	users := Table(
		"users",
		Column("id", Varchar().Size(36)),
		Column("email", Varchar()).Unique(),
	)

	var statement *Stmt

	statement = Delete(users).
		Where(Eq(users.C("id"), 5)).
		Build(dialect)

	assert.Equal(t, "DELETE FROM users\nWHERE users.id = ?;", statement.SQL())
	assert.Equal(t, []interface{}{5}, statement.Bindings())

	statement = Delete(users).
		Where(Eq(users.C("id"), 5)).
		Returning(users.C("id")).
		Build(dialect)

	assert.Equal(t, "DELETE FROM users\nWHERE users.id = ?\nRETURNING id;", statement.SQL())
	assert.Equal(t, []interface{}{5}, statement.Bindings())

	statement = Delete(users).Build(dialect)
	assert.Equal(t, "DELETE FROM users;", statement.SQL())
}
