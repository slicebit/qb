package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuery(t *testing.T) {
	query := Statement()

	query.AddClause("SELECT name")
	query.AddClause("FROM user")
	query.AddClause("WHERE id = ?")
	query.AddBinding(5)

	assert.Equal(t, query.Clauses(), []string{"SELECT name", "FROM user", "WHERE id = ?"})
	assert.Equal(t, query.Bindings(), []interface{}{5})
	assert.Equal(t, query.SQL(), "SELECT name\nFROM user\nWHERE id = ?;")
}

func TestQueryWithDelimiter(t *testing.T) {
	query := Statement()

	assert.Equal(t, query.SQL(), "")

	query.SetDelimiter(" ")

	query.AddClause("SELECT name")
	query.AddClause("FROM user")

	query.AddClause("WHERE id = ?")
	query.AddBinding(5)

	assert.Equal(t, query.Clauses(), []string{"SELECT name", "FROM user", "WHERE id = ?"})
	assert.Equal(t, query.Bindings(), []interface{}{5})
	assert.Equal(t, query.SQL(), "SELECT name FROM user WHERE id = ?;")
}
