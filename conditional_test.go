package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConditionals(t *testing.T) {
	country := Column("country", Varchar()).NotNull()

	var sql string
	var bindings interface{}

	like := Like(country, "%land%")

	sql, bindings = asDefSQLBinds(like)
	assert.Equal(t, "country LIKE ?", sql)
	assert.Equal(t, []interface{}{"%land%"}, bindings)

	notIn := NotIn(country, "USA", "England", "Sweden")

	sql, bindings = asDefSQLBinds(notIn)
	assert.Equal(t, "country NOT IN (?, ?, ?)", sql)
	assert.Equal(t, []interface{}{"USA", "England", "Sweden"}, bindings)

	in := In(country, "USA", "England", "Sweden")

	sql, bindings = asDefSQLBinds(in)
	assert.Equal(t, "country IN (?, ?, ?)", sql)
	assert.Equal(t, []interface{}{"USA", "England", "Sweden"}, bindings)

	notEq := NotEq(country, "USA")

	sql, bindings = asDefSQLBinds(notEq)
	assert.Equal(t, "country != ?", sql)
	assert.Equal(t, []interface{}{"USA"}, bindings)

	eq := Eq(country, "Turkey")

	sql, bindings = asDefSQLBinds(eq)
	assert.Equal(t, "country = ?", sql)
	assert.Equal(t, []interface{}{"Turkey"}, bindings)

	score := Column("score", BigInt()).NotNull()

	gt := Gt(score, 1500)

	sql, bindings = asDefSQLBinds(gt)
	assert.Equal(t, "score > ?", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

	lt := Lt(score, 1500)

	sql, bindings = asDefSQLBinds(lt)
	assert.Equal(t, "score < ?", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

	gte := Gte(score, 1500)

	sql, bindings = asDefSQLBinds(gte)
	assert.Equal(t, "score >= ?", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

	lte := Lte(score, 1500)

	sql, bindings = asDefSQLBinds(lte)
	assert.Equal(t, "score <= ?", sql)
	assert.Equal(t, []interface{}{1500}, bindings)
}
