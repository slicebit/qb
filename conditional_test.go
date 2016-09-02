package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConditionals(t *testing.T) {

	sqlite := NewDialect("sqlite3")
	sqlite.SetEscaping(true)

	mysql := NewDialect("mysql")
	mysql.SetEscaping(true)

	postgres := NewDialect("postgres")
	postgres.SetEscaping(true)

	country := Column("country", Varchar()).NotNull()

	var sql string
	var bindings interface{}

	like := Like(country, "%land%")

	sql, _ = like.Build(sqlite)
	assert.Equal(t, "country LIKE '%land%'", sql)

	sql, _ = like.Build(mysql)
	assert.Equal(t, "`country` LIKE '%land%'", sql)

	sql, _ = like.Build(postgres)
	assert.Equal(t, "\"country\" LIKE '%land%'", sql)

	notIn := NotIn(country, "USA", "England", "Sweden")

	sql, bindings = notIn.Build(sqlite)
	assert.Equal(t, "country NOT IN (?, ?, ?)", sql)
	assert.Equal(t, []interface{}{"USA", "England", "Sweden"}, bindings)

	sql, bindings = notIn.Build(mysql)
	assert.Equal(t, "`country` NOT IN (?, ?, ?)", sql)
	assert.Equal(t, []interface{}{"USA", "England", "Sweden"}, bindings)

	sql, bindings = notIn.Build(postgres)
	assert.Equal(t, "\"country\" NOT IN ($1, $2, $3)", sql)
	assert.Equal(t, []interface{}{"USA", "England", "Sweden"}, bindings)

	in := In(country, "USA", "England", "Sweden")

	sql, bindings = in.Build(sqlite)
	assert.Equal(t, "country IN (?, ?, ?)", sql)
	assert.Equal(t, []interface{}{"USA", "England", "Sweden"}, bindings)

	sql, bindings = in.Build(mysql)
	assert.Equal(t, "`country` IN (?, ?, ?)", sql)
	assert.Equal(t, []interface{}{"USA", "England", "Sweden"}, bindings)

	sql, bindings = in.Build(postgres)
	assert.Equal(t, "\"country\" IN ($4, $5, $6)", sql)
	assert.Equal(t, []interface{}{"USA", "England", "Sweden"}, bindings)

	notEq := NotEq(country, "USA")

	sql, bindings = notEq.Build(sqlite)
	assert.Equal(t, "country != ?", sql)
	assert.Equal(t, []interface{}{"USA"}, bindings)

	sql, bindings = notEq.Build(mysql)
	assert.Equal(t, "`country` != ?", sql)
	assert.Equal(t, []interface{}{"USA"}, bindings)

	sql, bindings = notEq.Build(postgres)
	assert.Equal(t, "\"country\" != $7", sql)
	assert.Equal(t, []interface{}{"USA"}, bindings)

	eq := Eq(country, "Turkey")

	sql, bindings = eq.Build(sqlite)
	assert.Equal(t, "country = ?", sql)
	assert.Equal(t, []interface{}{"Turkey"}, bindings)

	sql, bindings = eq.Build(mysql)
	assert.Equal(t, "`country` = ?", sql)
	assert.Equal(t, []interface{}{"Turkey"}, bindings)

	sql, bindings = eq.Build(postgres)
	assert.Equal(t, "\"country\" = $8", sql)
	assert.Equal(t, []interface{}{"Turkey"}, bindings)

	score := Column("score", BigInt()).NotNull()

	gt := Gt(score, 1500)

	sql, bindings = gt.Build(sqlite)
	assert.Equal(t, "score > ?", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

	sql, bindings = gt.Build(mysql)
	assert.Equal(t, "`score` > ?", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

	sql, bindings = gt.Build(postgres)
	assert.Equal(t, "\"score\" > $9", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

	lt := Lt(score, 1500)

	sql, bindings = lt.Build(sqlite)
	assert.Equal(t, "score < ?", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

	sql, bindings = lt.Build(mysql)
	assert.Equal(t, "`score` < ?", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

	sql, bindings = lt.Build(postgres)
	assert.Equal(t, "\"score\" < $10", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

	gte := Gte(score, 1500)

	sql, bindings = gte.Build(sqlite)
	assert.Equal(t, "score >= ?", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

	sql, bindings = gte.Build(mysql)
	assert.Equal(t, "`score` >= ?", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

	sql, bindings = gte.Build(postgres)
	assert.Equal(t, "\"score\" >= $11", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

	lte := Lte(score, 1500)

	sql, bindings = lte.Build(sqlite)
	assert.Equal(t, "score <= ?", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

	sql, bindings = lte.Build(mysql)
	assert.Equal(t, "`score` <= ?", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

	sql, bindings = lte.Build(postgres)
	assert.Equal(t, "\"score\" <= $12", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

}
