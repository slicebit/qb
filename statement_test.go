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

	assert.Equal(t, []string{"SELECT name", "FROM user", "WHERE id = ?"}, statement.Clauses())
	assert.Equal(t, []interface{}{5}, statement.Bindings())
	assert.Equal(t, "SELECT name\nFROM user\nWHERE id = ?;", statement.SQL())
}

func TestStatementRaw(t *testing.T) {

	statement := Statement()
	sql := `
		SELECT name
		FROM user
		WHERE id = ?;
		`
	statement.Text(sql)
	assert.Equal(t, []string{"SELECT name", "FROM user", "WHERE id = ?"}, statement.Clauses())
	assert.Equal(t, "SELECT name\nFROM user\nWHERE id = ?;", statement.SQL())
}

func TestStatementWithCustomDelimiter(t *testing.T) {
	statement := Statement()

	assert.Equal(t, "", statement.SQL())

	statement.SetDelimiter(" ")

	statement.AddClause("SELECT name")
	statement.AddClause("FROM user")

	statement.AddClause("WHERE id = ?")
	statement.AddBinding(5)

	assert.Equal(t, []string{"SELECT name", "FROM user", "WHERE id = ?"}, statement.Clauses())
	assert.Equal(t, []interface{}{5}, statement.Bindings())
	assert.Equal(t, "SELECT name FROM user WHERE id = ?;", statement.SQL())
}
