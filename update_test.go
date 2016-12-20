package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdate(t *testing.T) {
	users := Table(
		"users",
		Column("id", BigInt()).NotNull(),
		Column("email", Varchar()).NotNull().Unique(),
		PrimaryKey("email"),
	)

	sql, binds := asDefSQLBinds(Update(users).
		Values(map[string]interface{}{"email": "robert@de.niro"}))

	assert.Equal(t, "UPDATE users\nSET email = ?", sql)
	assert.Equal(t, []interface{}{"robert@de.niro"}, binds)

	sql, binds = asDefSQLBinds(Update(users).
		Values(map[string]interface{}{"email": "robert@de.niro"}))

	assert.Equal(t, "UPDATE users\nSET email = ?", sql)
	assert.Equal(t, []interface{}{"robert@de.niro"}, binds)

	sql, binds = asDefSQLBinds(Update(users).
		Values(map[string]interface{}{"email": "robert@de.niro"}).
		Where(Eq(users.C("email"), "al@pacino")).
		Returning(users.C("id"), users.C("email")))

	assert.Equal(t, "UPDATE users\nSET email = ?\nWHERE email = ?\nRETURNING id, email", sql)
	assert.Equal(t, []interface{}{"robert@de.niro", "al@pacino"}, binds)
}
