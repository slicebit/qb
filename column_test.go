package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestColumn(t *testing.T) {
	sqlite := NewDialect("sqlite3")
	sqlite.SetEscaping(true)

	mysql := NewDialect("mysql")
	mysql.SetEscaping(true)

	postgres := NewDialect("postgres")
	postgres.SetEscaping(true)

	col := Column("id", Varchar().Size(40))
	assert.Equal(t, "id", col.Name)
	assert.Equal(t, Varchar().Size(40), col.Type)

	assert.Equal(t, "id VARCHAR(40)", col.String(sqlite))
	assert.Equal(t, "`id` VARCHAR(40)", col.String(mysql))
	assert.Equal(t, "\"id\" VARCHAR(40)", col.String(postgres))

	col = Column("s", Varchar().Size(255)).Unique().NotNull().Default("hello")
	assert.Equal(t, "s VARCHAR(255) UNIQUE NOT NULL DEFAULT 'hello'", col.String(sqlite))

	precisionCol := Column("f", Type("FLOAT").Precision(2, 5)).Null()
	assert.Equal(t, "f FLOAT(2, 5) NULL", precisionCol.String(sqlite))

	col = Column("id", Int()).PrimaryKey().AutoIncrement()

	assert.Equal(t, "id INTEGER PRIMARY KEY", col.String(sqlite))
	assert.Equal(t, "`id` INT PRIMARY KEY AUTO_INCREMENT", col.String(mysql))
	assert.Equal(t, "\"id\" SERIAL PRIMARY KEY", col.String(postgres))

	assert.Equal(t, "c INT TEST", Column("c", Int()).Constraint("TEST").String(sqlite))

	// like
	like := col.Like("s%")
	likeSqlite, likeBindingsSqlite := like.Build(sqlite)
	likeMysql, likeBindingsMysql := like.Build(mysql)
	likePostgres, likeBindingsPostgres := like.Build(postgres)

	assert.Equal(t, "id LIKE 's%'", likeSqlite)
	assert.Equal(t, []interface{}{}, likeBindingsSqlite)
	assert.Equal(t, "`id` LIKE 's%'", likeMysql)
	assert.Equal(t, []interface{}{}, likeBindingsMysql)
	assert.Equal(t, "\"id\" LIKE 's%'", likePostgres)
	assert.Equal(t, []interface{}{}, likeBindingsPostgres)

	// not in
	notIn := col.NotIn("id1", "id2")
	notInSqlite, likeBindingsSqlite := notIn.Build(sqlite)
	notInMysql, likeBindingsMysql := notIn.Build(mysql)
	notInPostgres, likeBindingsPostgres := notIn.Build(postgres)

	assert.Equal(t, "id NOT IN (?, ?)", notInSqlite)
	assert.Equal(t, []interface{}{"id1", "id2"}, likeBindingsSqlite)
	assert.Equal(t, "`id` NOT IN (?, ?)", notInMysql)
	assert.Equal(t, []interface{}{"id1", "id2"}, likeBindingsMysql)
	assert.Equal(t, "\"id\" NOT IN ($1, $2)", notInPostgres)
	assert.Equal(t, []interface{}{"id1", "id2"}, likeBindingsPostgres)

	postgres.Reset()

	// in
	in := col.In("id1", "id2")
	inSqlite, likeBindingsSqlite := in.Build(sqlite)
	inMysql, likeBindingsMysql := in.Build(mysql)
	inPostgres, likeBindingsPostgres := in.Build(postgres)

	assert.Equal(t, "id IN (?, ?)", inSqlite)
	assert.Equal(t, []interface{}{"id1", "id2"}, likeBindingsSqlite)
	assert.Equal(t, "`id` IN (?, ?)", inMysql)
	assert.Equal(t, []interface{}{"id1", "id2"}, likeBindingsMysql)
	assert.Equal(t, "\"id\" IN ($1, $2)", inPostgres)
	assert.Equal(t, []interface{}{"id1", "id2"}, likeBindingsPostgres)

	postgres.Reset()

	// not eq
	notEq := col.NotEq("id1")
	notEqSqlite, likeBindingsSqlite := notEq.Build(sqlite)
	notEqMysql, likeBindingsMysql := notEq.Build(mysql)
	notEqPostgres, likeBindingsPostgres := notEq.Build(postgres)

	assert.Equal(t, "id != ?", notEqSqlite)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsSqlite)
	assert.Equal(t, "`id` != ?", notEqMysql)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsMysql)
	assert.Equal(t, "\"id\" != $1", notEqPostgres)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsPostgres)

	postgres.Reset()

	// eq
	eq := col.Eq("id1")
	eqSqlite, likeBindingsSqlite := eq.Build(sqlite)
	eqMysql, likeBindingsMysql := eq.Build(mysql)
	eqPostgres, likeBindingsPostgres := eq.Build(postgres)

	assert.Equal(t, "id = ?", eqSqlite)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsSqlite)
	assert.Equal(t, "`id` = ?", eqMysql)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsMysql)
	assert.Equal(t, "\"id\" = $1", eqPostgres)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsPostgres)

	postgres.Reset()

	// gt
	gt := col.Gt("id1")
	gtSqlite, likeBindingsSqlite := gt.Build(sqlite)
	gtMysql, likeBindingsMysql := gt.Build(mysql)
	gtPostgres, likeBindingsPostgres := gt.Build(postgres)

	assert.Equal(t, "id > ?", gtSqlite)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsSqlite)
	assert.Equal(t, "`id` > ?", gtMysql)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsMysql)
	assert.Equal(t, "\"id\" > $1", gtPostgres)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsPostgres)

	postgres.Reset()

	// lt
	lt := col.Lt("id1")
	ltSqlite, likeBindingsSqlite := lt.Build(sqlite)
	ltMysql, likeBindingsMysql := lt.Build(mysql)
	ltPostgres, likeBindingsPostgres := lt.Build(postgres)

	assert.Equal(t, "id < ?", ltSqlite)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsSqlite)
	assert.Equal(t, "`id` < ?", ltMysql)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsMysql)
	assert.Equal(t, "\"id\" < $1", ltPostgres)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsPostgres)

	postgres.Reset()

	// gte
	gte := col.Gte("id1")
	gteSqlite, likeBindingsSqlite := gte.Build(sqlite)
	gteMysql, likeBindingsMysql := gte.Build(mysql)
	gtePostgres, likeBindingsPostgres := gte.Build(postgres)

	assert.Equal(t, "id >= ?", gteSqlite)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsSqlite)
	assert.Equal(t, "`id` >= ?", gteMysql)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsMysql)
	assert.Equal(t, "\"id\" >= $1", gtePostgres)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsPostgres)

	postgres.Reset()

	// lte
	lte := col.Lte("id1")
	lteSqlite, likeBindingsSqlite := lte.Build(sqlite)
	lteMysql, likeBindingsMysql := lte.Build(mysql)
	ltePostgres, likeBindingsPostgres := lte.Build(postgres)

	assert.Equal(t, "id <= ?", lteSqlite)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsSqlite)
	assert.Equal(t, "`id` <= ?", lteMysql)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsMysql)
	assert.Equal(t, "\"id\" <= $1", ltePostgres)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsPostgres)

	postgres.Reset()

	var sql string

	sql, _ = col.Build(sqlite)
	assert.Equal(t, "id", sql)

	sql, _ = col.Build(mysql)
	assert.Equal(t, "`id`", sql)

	sql, _ = col.Build(postgres)
	assert.Equal(t, "\"id\"", sql)
}
