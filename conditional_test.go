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

	country := Column("country", Varchar().NotNull())

	var sql string
	var bindings interface{}

	like := Like(country, "%land%")

	sql, _ = like.Build(sqlite)
	assert.Equal(t, sql, "country LIKE '%land%'")

	sql, _ = like.Build(mysql)
	assert.Equal(t, sql, "`country` LIKE '%land%'")

	sql, _ = like.Build(postgres)
	assert.Equal(t, sql, "\"country\" LIKE '%land%'")

	notIn := NotIn(country, "USA", "England", "Sweden")

	sql, bindings = notIn.Build(sqlite)
	assert.Equal(t, sql, "country NOT IN (?, ?, ?)")
	assert.Equal(t, bindings, []interface{}{"USA", "England", "Sweden"})

	sql, bindings = notIn.Build(mysql)
	assert.Equal(t, sql, "`country` NOT IN (?, ?, ?)")
	assert.Equal(t, bindings, []interface{}{"USA", "England", "Sweden"})

	sql, bindings = notIn.Build(postgres)
	assert.Equal(t, sql, "\"country\" NOT IN ($1, $2, $3)")
	assert.Equal(t, bindings, []interface{}{"USA", "England", "Sweden"})

	in := In(country, "USA", "England", "Sweden")

	sql, bindings = in.Build(sqlite)
	assert.Equal(t, sql, "country IN (?, ?, ?)")
	assert.Equal(t, bindings, []interface{}{"USA", "England", "Sweden"})

	sql, bindings = in.Build(mysql)
	assert.Equal(t, sql, "`country` IN (?, ?, ?)")
	assert.Equal(t, bindings, []interface{}{"USA", "England", "Sweden"})

	sql, bindings = in.Build(postgres)
	assert.Equal(t, sql, "\"country\" IN ($4, $5, $6)")
	assert.Equal(t, bindings, []interface{}{"USA", "England", "Sweden"})

	notEq := NotEq(country, "USA")

	sql, bindings = notEq.Build(sqlite)
	assert.Equal(t, sql, "country != ?")
	assert.Equal(t, bindings, []interface{}{"USA"})

	sql, bindings = notEq.Build(mysql)
	assert.Equal(t, sql, "`country` != ?")
	assert.Equal(t, bindings, []interface{}{"USA"})

	sql, bindings = notEq.Build(postgres)
	assert.Equal(t, sql, "\"country\" != $7")
	assert.Equal(t, bindings, []interface{}{"USA"})

	eq := Eq(country, "Turkey")

	sql, bindings = eq.Build(sqlite)
	assert.Equal(t, sql, "country = ?")
	assert.Equal(t, bindings, []interface{}{"Turkey"})

	sql, bindings = eq.Build(mysql)
	assert.Equal(t, sql, "`country` = ?")
	assert.Equal(t, bindings, []interface{}{"Turkey"})

	sql, bindings = eq.Build(postgres)
	assert.Equal(t, sql, "\"country\" = $8")
	assert.Equal(t, bindings, []interface{}{"Turkey"})

	score := Column("score", BigInt().NotNull())

	gt := Gt(score, 1500)

	sql, bindings = gt.Build(sqlite)
	assert.Equal(t, sql, "score > ?")
	assert.Equal(t, bindings, []interface{}{1500})

	sql, bindings = gt.Build(mysql)
	assert.Equal(t, sql, "`score` > ?")
	assert.Equal(t, bindings, []interface{}{1500})

	sql, bindings = gt.Build(postgres)
	assert.Equal(t, sql, "\"score\" > $9")
	assert.Equal(t, bindings, []interface{}{1500})

	st := St(score, 1500)

	sql, bindings = st.Build(sqlite)
	assert.Equal(t, sql, "score < ?")
	assert.Equal(t, bindings, []interface{}{1500})

	sql, bindings = st.Build(mysql)
	assert.Equal(t, sql, "`score` < ?")
	assert.Equal(t, bindings, []interface{}{1500})

	sql, bindings = st.Build(postgres)
	assert.Equal(t, sql, "\"score\" < $10")
	assert.Equal(t, bindings, []interface{}{1500})

	gte := Gte(score, 1500)

	sql, bindings = gte.Build(sqlite)
	assert.Equal(t, sql, "score >= ?")
	assert.Equal(t, bindings, []interface{}{1500})

	sql, bindings = gte.Build(mysql)
	assert.Equal(t, sql, "`score` >= ?")
	assert.Equal(t, bindings, []interface{}{1500})

	sql, bindings = gte.Build(postgres)
	assert.Equal(t, sql, "\"score\" >= $11")
	assert.Equal(t, bindings, []interface{}{1500})

	ste := Ste(score, 1500)

	sql, bindings = ste.Build(sqlite)
	assert.Equal(t, sql, "score <= ?")
	assert.Equal(t, bindings, []interface{}{1500})

	sql, bindings = ste.Build(mysql)
	assert.Equal(t, sql, "`score` <= ?")
	assert.Equal(t, bindings, []interface{}{1500})

	sql, bindings = ste.Build(postgres)
	assert.Equal(t, sql, "\"score\" <= $12")
	assert.Equal(t, bindings, []interface{}{1500})

}
