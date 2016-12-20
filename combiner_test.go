package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCombiners(t *testing.T) {
	email := Column("email", Varchar()).NotNull().Unique()
	id := Column("id", Int()).NotNull()

	and := And(Eq(email, "al@pacino.com"), NotEq(id, 1))
	or := Or(Eq(email, "al@pacino.com"), NotEq(id, 1))

	sql, binds := asDefSQLBinds(and)

	assert.Equal(t, "(email = ? AND id != ?)", sql)
	assert.Equal(t, []interface{}{"al@pacino.com", 1}, binds)

	sql, binds = asDefSQLBinds(or)

	assert.Equal(t, "(email = ? OR id != ?)", sql)
	assert.Equal(t, []interface{}{"al@pacino.com", 1}, binds)
}
