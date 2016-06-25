package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStatement(t *testing.T) {
	statement := Statement()

	statement.AddClause("SELECT name")
	statement.AddClause("FROM user")
	statement.AddClause("WHERE id = ?")
	statement.AddBinding(5)

	assert.Equal(t, statement.Clauses(), []string{"SELECT name", "FROM user", "WHERE id = ?"})
	assert.Equal(t, statement.Bindings(), []interface{}{5})
	assert.Equal(t, statement.SQL(), "SELECT name\nFROM user\nWHERE id = ?;")
}

func TestStatementRaw(t *testing.T) {

	statement := Statement()
	sql := `
		SELECT name
		FROM user
		WHERE id = ?;
		`
	statement.Text(sql)
	assert.Equal(t, statement.Clauses(), []string{"SELECT name", "FROM user", "WHERE id = ?"})
	assert.Equal(t, statement.SQL(), "SELECT name\nFROM user\nWHERE id = ?;")
}

func TestStatementWithCustomDelimiter(t *testing.T) {
	statement := Statement()

	assert.Equal(t, statement.SQL(), "")

	statement.SetDelimiter(" ")

	statement.AddClause("SELECT name")
	statement.AddClause("FROM user")

	statement.AddClause("WHERE id = ?")
	statement.AddBinding(5)

	assert.Equal(t, statement.Clauses(), []string{"SELECT name", "FROM user", "WHERE id = ?"})
	assert.Equal(t, statement.Bindings(), []interface{}{5})
	assert.Equal(t, statement.SQL(), "SELECT name FROM user WHERE id = ?;")
}
