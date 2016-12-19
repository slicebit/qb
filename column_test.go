package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestColumn(t *testing.T) {
	dialect := NewDialect("default")

	col := Column("id", Varchar().Size(40))
	assert.Equal(t, "id", col.Name)
	assert.Equal(t, Varchar().Size(40), col.Type)

	assert.Equal(t, "id VARCHAR(40)", col.String(dialect))

	col = Column("s", Varchar().Size(255)).Unique().NotNull().Default("hello")
	assert.Equal(t, "s VARCHAR(255) UNIQUE NOT NULL DEFAULT 'hello'", col.String(dialect))

	precisionCol := Column("f", Type("FLOAT").Precision(2, 5)).Null()
	assert.Equal(t, "f FLOAT(2, 5) NULL", precisionCol.String(dialect))

	col = Column("id", Int()).PrimaryKey().AutoIncrement().inlinePrimaryKey()

	assert.Equal(t, "id INT PRIMARY KEY AUTO INCREMENT", col.String(dialect))

	assert.Equal(t, "c INT TEST", Column("c", Int()).Constraint("TEST").String(dialect))

	// like
	like := col.Like("s%")
	sql, binds := asDefSQLBinds(like)

	assert.Equal(t, "id LIKE ?", sql)
	assert.Equal(t, []interface{}{"s%"}, binds)

	// not in
	notIn := col.NotIn("id1", "id2")
	sql, binds = asDefSQLBinds(notIn)

	assert.Equal(t, "id NOT IN (?, ?)", sql)
	assert.Equal(t, []interface{}{"id1", "id2"}, binds)

	// in
	in := col.In("id1", "id2")
	sql, binds = asDefSQLBinds(in)

	assert.Equal(t, "id IN (?, ?)", sql)
	assert.Equal(t, []interface{}{"id1", "id2"}, binds)

	// not eq
	notEq := col.NotEq("id1")
	sql, binds = asDefSQLBinds(notEq)

	assert.Equal(t, "id != ?", sql)
	assert.Equal(t, []interface{}{"id1"}, binds)

	// eq
	eq := col.Eq("id1")
	sql, binds = asDefSQLBinds(eq)

	assert.Equal(t, "id = ?", sql)
	assert.Equal(t, []interface{}{"id1"}, binds)

	// gt
	gt := col.Gt("id1")
	sql, binds = asDefSQLBinds(gt)

	assert.Equal(t, "id > ?", sql)
	assert.Equal(t, []interface{}{"id1"}, binds)

	// lt
	lt := col.Lt("id1")
	sql, binds = asDefSQLBinds(lt)

	assert.Equal(t, "id < ?", sql)
	assert.Equal(t, []interface{}{"id1"}, binds)

	// gte
	gte := col.Gte("id1")
	sql, binds = asDefSQLBinds(gte)

	assert.Equal(t, "id >= ?", sql)
	assert.Equal(t, []interface{}{"id1"}, binds)

	// lte
	lte := col.Lte("id1")
	sql, binds = asDefSQLBinds(lte)

	assert.Equal(t, "id <= ?", sql)
	assert.Equal(t, []interface{}{"id1"}, binds)

	sql = col.Accept(NewCompilerContext(dialect))
	assert.Equal(t, "id", sql)
}
