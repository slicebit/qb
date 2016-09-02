package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSQLText(t *testing.T) {
	text := SQLText("1")
	assert.Equal(t, "1", text.Text)
}

func TestGetClauseFrom(t *testing.T) {
	var c Clause
	c = SQLText("1")
	assert.Equal(t, c, GetClauseFrom(c))

	c = GetClauseFrom(2)
	b, ok := c.(BindClause)
	assert.True(t, ok, "Should have returned a BindClause")
	assert.Equal(t, 2, b.Value)
}

func TestGetListFrom(t *testing.T) {
	var c Clause
	c = ListClause{}
	assert.Equal(t, c, GetListFrom(c))

	text := SQLText("SOME SQL")
	c = GetListFrom(text)
	l, ok := c.(ListClause)
	assert.True(t, ok, "Should have returned a ListClause")
	assert.Equal(t, 1, len(l.Clauses))
	assert.Equal(t, text, l.Clauses[0])

	c = GetListFrom([]int{2})
	l, ok = c.(ListClause)
	assert.True(t, ok, "Should have returned a ListClause")
	assert.Equal(t, 1, len(l.Clauses))
	assert.Equal(t, 2, l.Clauses[0].(BindClause).Value)

	c = GetListFrom([]interface{}{2, Bind(4)})
	l, ok = c.(ListClause)
	assert.True(t, ok, "Should have returned a ListClause")
	assert.Equal(t, 2, len(l.Clauses))
	assert.Equal(t, 2, l.Clauses[0].(BindClause).Value)
	assert.Equal(t, 4, l.Clauses[1].(BindClause).Value)
}
