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
	likeSqlite, likeBindingsSqlite := asSQL(like, sqlite)
	likeMysql, likeBindingsMysql := asSQL(like, mysql)
	likePostgres, likeBindingsPostgres := asSQL(like, postgres)

	assert.Equal(t, "id LIKE 's%'", likeSqlite)
	assert.Equal(t, []interface{}{}, likeBindingsSqlite)
	assert.Equal(t, "`id` LIKE 's%'", likeMysql)
	assert.Equal(t, []interface{}{}, likeBindingsMysql)
	assert.Equal(t, "\"id\" LIKE 's%'", likePostgres)
	assert.Equal(t, []interface{}{}, likeBindingsPostgres)

	// not in
	notIn := col.NotIn("id1", "id2")
	notInSqlite, likeBindingsSqlite := asSQL(notIn, sqlite)
	notInMysql, likeBindingsMysql := asSQL(notIn, mysql)
	notInPostgres, likeBindingsPostgres := asSQL(notIn, postgres)

	assert.Equal(t, "id NOT IN (?, ?)", notInSqlite)
	assert.Equal(t, []interface{}{"id1", "id2"}, likeBindingsSqlite)
	assert.Equal(t, "`id` NOT IN (?, ?)", notInMysql)
	assert.Equal(t, []interface{}{"id1", "id2"}, likeBindingsMysql)
	assert.Equal(t, "\"id\" NOT IN ($1, $2)", notInPostgres)
	assert.Equal(t, []interface{}{"id1", "id2"}, likeBindingsPostgres)

	postgres.Reset()

	// in
	in := col.In("id1", "id2")
	inSqlite, likeBindingsSqlite := asSQL(in, sqlite)
	inMysql, likeBindingsMysql := asSQL(in, mysql)
	inPostgres, likeBindingsPostgres := asSQL(in, postgres)

	assert.Equal(t, "id IN (?, ?)", inSqlite)
	assert.Equal(t, []interface{}{"id1", "id2"}, likeBindingsSqlite)
	assert.Equal(t, "`id` IN (?, ?)", inMysql)
	assert.Equal(t, []interface{}{"id1", "id2"}, likeBindingsMysql)
	assert.Equal(t, "\"id\" IN ($1, $2)", inPostgres)
	assert.Equal(t, []interface{}{"id1", "id2"}, likeBindingsPostgres)

	postgres.Reset()

	// not eq
	notEq := col.NotEq("id1")
	notEqSqlite, likeBindingsSqlite := asSQL(notEq, sqlite)
	notEqMysql, likeBindingsMysql := asSQL(notEq, mysql)
	notEqPostgres, likeBindingsPostgres := asSQL(notEq, postgres)

	assert.Equal(t, "id != ?", notEqSqlite)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsSqlite)
	assert.Equal(t, "`id` != ?", notEqMysql)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsMysql)
	assert.Equal(t, "\"id\" != $1", notEqPostgres)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsPostgres)

	postgres.Reset()

	// eq
	eq := col.Eq("id1")
	eqSqlite, likeBindingsSqlite := asSQL(eq, sqlite)
	eqMysql, likeBindingsMysql := asSQL(eq, mysql)
	eqPostgres, likeBindingsPostgres := asSQL(eq, postgres)

	assert.Equal(t, "id = ?", eqSqlite)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsSqlite)
	assert.Equal(t, "`id` = ?", eqMysql)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsMysql)
	assert.Equal(t, "\"id\" = $1", eqPostgres)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsPostgres)

	postgres.Reset()

	// gt
	gt := col.Gt("id1")
	gtSqlite, likeBindingsSqlite := asSQL(gt, sqlite)
	gtMysql, likeBindingsMysql := asSQL(gt, mysql)
	gtPostgres, likeBindingsPostgres := asSQL(gt, postgres)

	assert.Equal(t, "id > ?", gtSqlite)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsSqlite)
	assert.Equal(t, "`id` > ?", gtMysql)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsMysql)
	assert.Equal(t, "\"id\" > $1", gtPostgres)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsPostgres)

	postgres.Reset()

	// lt
	lt := col.Lt("id1")
	ltSqlite, likeBindingsSqlite := asSQL(lt, sqlite)
	ltMysql, likeBindingsMysql := asSQL(lt, mysql)
	ltPostgres, likeBindingsPostgres := asSQL(lt, postgres)

	assert.Equal(t, "id < ?", ltSqlite)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsSqlite)
	assert.Equal(t, "`id` < ?", ltMysql)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsMysql)
	assert.Equal(t, "\"id\" < $1", ltPostgres)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsPostgres)

	postgres.Reset()

	// gte
	gte := col.Gte("id1")
	gteSqlite, likeBindingsSqlite := asSQL(gte, sqlite)
	gteMysql, likeBindingsMysql := asSQL(gte, mysql)
	gtePostgres, likeBindingsPostgres := asSQL(gte, postgres)

	assert.Equal(t, "id >= ?", gteSqlite)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsSqlite)
	assert.Equal(t, "`id` >= ?", gteMysql)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsMysql)
	assert.Equal(t, "\"id\" >= $1", gtePostgres)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsPostgres)

	postgres.Reset()

	// lte
	lte := col.Lte("id1")
	lteSqlite, likeBindingsSqlite := asSQL(lte, sqlite)
	lteMysql, likeBindingsMysql := asSQL(lte, mysql)
	ltePostgres, likeBindingsPostgres := asSQL(lte, postgres)

	assert.Equal(t, "id <= ?", lteSqlite)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsSqlite)
	assert.Equal(t, "`id` <= ?", lteMysql)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsMysql)
	assert.Equal(t, "\"id\" <= $1", ltePostgres)
	assert.Equal(t, []interface{}{"id1"}, likeBindingsPostgres)

	postgres.Reset()

	var sql string

	sql = col.Accept(NewCompilerContext(sqlite))
	assert.Equal(t, "id", sql)

	sql = col.Accept(NewCompilerContext(mysql))
	assert.Equal(t, "`id`", sql)

	sql = col.Accept(NewCompilerContext(postgres))
	assert.Equal(t, "\"id\"", sql)
}
