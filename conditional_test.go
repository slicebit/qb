package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConditionals(t *testing.T) {

	compile := func(c Clause, d Dialect) (string, []interface{}) {
		ctx := NewCompilerContext(d)
		return c.Accept(ctx), ctx.Binds
	}

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

	sql, bindings = compile(like, sqlite)
	assert.Equal(t, "country LIKE ?", sql)
	assert.Equal(t, []interface{}{"%land%"}, bindings)

	sql, _ = compile(like, mysql)
	assert.Equal(t, "`country` LIKE ?", sql)
	assert.Equal(t, []interface{}{"%land%"}, bindings)

	sql, _ = compile(like, postgres)
	assert.Equal(t, "\"country\" LIKE $1", sql)
	assert.Equal(t, []interface{}{"%land%"}, bindings)

	notIn := NotIn(country, "USA", "England", "Sweden")

	sql, bindings = compile(notIn, sqlite)
	assert.Equal(t, "country NOT IN (?, ?, ?)", sql)
	assert.Equal(t, []interface{}{"USA", "England", "Sweden"}, bindings)

	sql, bindings = compile(notIn, mysql)
	assert.Equal(t, "`country` NOT IN (?, ?, ?)", sql)
	assert.Equal(t, []interface{}{"USA", "England", "Sweden"}, bindings)

	sql, bindings = compile(notIn, postgres)
	assert.Equal(t, "\"country\" NOT IN ($1, $2, $3)", sql)
	assert.Equal(t, []interface{}{"USA", "England", "Sweden"}, bindings)

	in := In(country, "USA", "England", "Sweden")

	sql, bindings = compile(in, sqlite)
	assert.Equal(t, "country IN (?, ?, ?)", sql)
	assert.Equal(t, []interface{}{"USA", "England", "Sweden"}, bindings)

	sql, bindings = compile(in, mysql)
	assert.Equal(t, "`country` IN (?, ?, ?)", sql)
	assert.Equal(t, []interface{}{"USA", "England", "Sweden"}, bindings)

	sql, bindings = compile(in, postgres)
	assert.Equal(t, "\"country\" IN ($1, $2, $3)", sql)
	assert.Equal(t, []interface{}{"USA", "England", "Sweden"}, bindings)

	notEq := NotEq(country, "USA")

	sql, bindings = compile(notEq, sqlite)
	assert.Equal(t, "country != ?", sql)
	assert.Equal(t, []interface{}{"USA"}, bindings)

	sql, bindings = compile(notEq, mysql)
	assert.Equal(t, "`country` != ?", sql)
	assert.Equal(t, []interface{}{"USA"}, bindings)

	sql, bindings = compile(notEq, postgres)
	assert.Equal(t, "\"country\" != $1", sql)
	assert.Equal(t, []interface{}{"USA"}, bindings)

	eq := Eq(country, "Turkey")

	sql, bindings = compile(eq, sqlite)
	assert.Equal(t, "country = ?", sql)
	assert.Equal(t, []interface{}{"Turkey"}, bindings)

	sql, bindings = compile(eq, mysql)
	assert.Equal(t, "`country` = ?", sql)
	assert.Equal(t, []interface{}{"Turkey"}, bindings)

	sql, bindings = compile(eq, postgres)
	assert.Equal(t, "\"country\" = $1", sql)
	assert.Equal(t, []interface{}{"Turkey"}, bindings)

	score := Column("score", BigInt()).NotNull()

	gt := Gt(score, 1500)

	sql, bindings = compile(gt, sqlite)
	assert.Equal(t, "score > ?", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

	sql, bindings = compile(gt, mysql)
	assert.Equal(t, "`score` > ?", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

	sql, bindings = compile(gt, postgres)
	assert.Equal(t, "\"score\" > $1", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

	lt := Lt(score, 1500)

	sql, bindings = compile(lt, sqlite)
	assert.Equal(t, "score < ?", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

	sql, bindings = compile(lt, mysql)
	assert.Equal(t, "`score` < ?", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

	sql, bindings = compile(lt, postgres)

	assert.Equal(t, "\"score\" < $1", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

	gte := Gte(score, 1500)

	sql, bindings = compile(gte, sqlite)
	assert.Equal(t, "score >= ?", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

	sql, bindings = compile(gte, mysql)
	assert.Equal(t, "`score` >= ?", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

	sql, bindings = compile(gte, postgres)
	assert.Equal(t, "\"score\" >= $1", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

	lte := Lte(score, 1500)

	sql, bindings = compile(lte, sqlite)
	assert.Equal(t, "score <= ?", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

	sql, bindings = compile(lte, mysql)
	assert.Equal(t, "`score` <= ?", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

	sql, bindings = compile(lte, postgres)

	assert.Equal(t, "\"score\" <= $1", sql)
	assert.Equal(t, []interface{}{1500}, bindings)

}
