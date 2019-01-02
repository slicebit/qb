package qb

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"
)

type ColumnTestSuite struct {
	suite.Suite
	dialect Dialect
	ctx     Context
}

func (suite *ColumnTestSuite) SetupTest() {
	suite.dialect = NewDefaultDialect()
	suite.ctx = NewCompilerContext(suite.dialect)
}

func (suite *ColumnTestSuite) TestColumnVarcharSpecificSize() {
	col := Column("id", Varchar().Size(40))
	assert.Equal(suite.T(), "id", col.Name)
	assert.Equal(suite.T(), Varchar().Size(40), col.Type)
	assert.Equal(suite.T(), "id VARCHAR(40)", col.String(suite.dialect))
}

func (suite *ColumnTestSuite) TestColumnVarcharUniqueNotNullDefault() {
	col := Column("s", Varchar().Size(255)).Unique().NotNull().Default("hello")
	assert.Equal(suite.T(), "s VARCHAR(255) UNIQUE NOT NULL DEFAULT 'hello'", col.String(suite.dialect))
}

func (suite *ColumnTestSuite) TestColumnFloatPrecision() {
	col := Column("f", Type("FLOAT").Precision(2, 5)).Null()
	assert.Equal(suite.T(), "f FLOAT(2, 5) NULL", col.String(suite.dialect))
}

func (suite *ColumnTestSuite) TestColumnIntInlinePrimaryKeyAutoIncrement() {
	col := Column("id", Int()).PrimaryKey().AutoIncrement().inlinePrimaryKey()
	assert.Equal(suite.T(), "id INT PRIMARY KEY AUTO INCREMENT", col.String(suite.dialect))
	assert.Equal(suite.T(), "c INT TEST", Column("c", Int()).Constraint("TEST").String(suite.dialect))
}

func (suite *ColumnTestSuite) TestColumnLike() {
	col := Column("id", Int()).PrimaryKey().AutoIncrement().inlinePrimaryKey()
	like := col.Like("s%")

	sql := like.Accept(suite.ctx)
	binds := suite.ctx.Binds()

	assert.Equal(suite.T(), "id LIKE ?", sql)
	assert.Equal(suite.T(), []interface{}{"s%"}, binds)
}

func (suite *ColumnTestSuite) TestColumnNotIn() {
	col := Column("id", Int()).PrimaryKey().AutoIncrement().inlinePrimaryKey()
	notIn := col.NotIn("id1", "id2")
	sql := notIn.Accept(suite.ctx)
	binds := suite.ctx.Binds()
	assert.Equal(suite.T(), "id NOT IN (?, ?)", sql)
	assert.Equal(suite.T(), []interface{}{"id1", "id2"}, binds)
}

func (suite *ColumnTestSuite) TestColumnIn() {
	col := Column("id", Int()).PrimaryKey().AutoIncrement().inlinePrimaryKey()
	in := col.In("id1", "id2")
	sql := in.Accept(suite.ctx)
	binds := suite.ctx.Binds()

	assert.Equal(suite.T(), "id IN (?, ?)", sql)
	assert.Equal(suite.T(), []interface{}{"id1", "id2"}, binds)
}

func (suite *ColumnTestSuite) TestColumnNotEq() {
	col := Column("id", Int()).PrimaryKey().AutoIncrement().inlinePrimaryKey()
	notEq := col.NotEq("id1")
	sql := notEq.Accept(suite.ctx)
	binds := suite.ctx.Binds()

	assert.Equal(suite.T(), "id != ?", sql)
	assert.Equal(suite.T(), []interface{}{"id1"}, binds)
}

func (suite *ColumnTestSuite) TestColumnEq() {
	col := Column("id", Int()).PrimaryKey().AutoIncrement().inlinePrimaryKey()
	eq := col.Eq("id1")
	sql := eq.Accept(suite.ctx)
	binds := suite.ctx.Binds()

	assert.Equal(suite.T(), "id = ?", sql)
	assert.Equal(suite.T(), []interface{}{"id1"}, binds)
}

func (suite *ColumnTestSuite) TestColumnGt() {
	col := Column("id", Int()).PrimaryKey().AutoIncrement().inlinePrimaryKey()
	gt := col.Gt("id1")
	sql := gt.Accept(suite.ctx)
	binds := suite.ctx.Binds()

	assert.Equal(suite.T(), "id > ?", sql)
	assert.Equal(suite.T(), []interface{}{"id1"}, binds)
}

func (suite *ColumnTestSuite) TestColumnLt() {
	col := Column("id", Int()).PrimaryKey().AutoIncrement().inlinePrimaryKey()
	lt := col.Lt("id1")
	sql := lt.Accept(suite.ctx)
	binds := suite.ctx.Binds()

	assert.Equal(suite.T(), "id < ?", sql)
	assert.Equal(suite.T(), []interface{}{"id1"}, binds)
}

func (suite *ColumnTestSuite) TestcolumnGte() {
	col := Column("id", Int()).PrimaryKey().AutoIncrement().inlinePrimaryKey()
	gte := col.Gte("id1")
	sql := gte.Accept(suite.ctx)
	binds := suite.ctx.Binds()

	assert.Equal(suite.T(), "id >= ?", sql)
	assert.Equal(suite.T(), []interface{}{"id1"}, binds)
}

func (suite *ColumnTestSuite) TestColumnLte() {
	col := Column("id", Int()).PrimaryKey().AutoIncrement().inlinePrimaryKey()
	lte := col.Lte("id1")
	sql := lte.Accept(suite.ctx)
	binds := suite.ctx.Binds()

	assert.Equal(suite.T(), "id <= ?", sql)
	assert.Equal(suite.T(), []interface{}{"id1"}, binds)
}

func TestColumnTestSuite(t *testing.T) {
	suite.Run(t, new(ColumnTestSuite))
}
