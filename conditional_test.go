package qb

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"
)

type ConditionalTestSuite struct {
	suite.Suite
	dialect Dialect
	ctx     Context
	country ColumnElem
	score   ColumnElem
}

func (suite *ConditionalTestSuite) SetupTest() {
	suite.dialect = NewDefaultDialect()
	suite.ctx = NewCompilerContext(suite.dialect)
	suite.country = Column("country", Varchar()).NotNull()
	suite.score = Column("score", BigInt()).NotNull()
}

func (suite *ConditionalTestSuite) TestConditionalLike() {
	like := Like(suite.country, "%land%")
	sql := like.Accept(suite.ctx)
	bindings := suite.ctx.Binds()

	assert.Equal(suite.T(), "country LIKE ?", sql)
	assert.Equal(suite.T(), []interface{}{"%land%"}, bindings)
}

func (suite *ConditionalTestSuite) TestConditionalNotIn() {
	notIn := NotIn(suite.country, "USA", "England", "Sweden")
	sql := notIn.Accept(suite.ctx)
	bindings := suite.ctx.Binds()

	assert.Equal(suite.T(), "country NOT IN (?, ?, ?)", sql)
	assert.Equal(suite.T(), []interface{}{"USA", "England", "Sweden"}, bindings)
}

func (suite *ConditionalTestSuite) TestConditionalIn() {
	in := In(suite.country, "USA", "England", "Sweden")
	sql := in.Accept(suite.ctx)
	bindings := suite.ctx.Binds()
	assert.Equal(suite.T(), "country IN (?, ?, ?)", sql)
	assert.Equal(suite.T(), []interface{}{"USA", "England", "Sweden"}, bindings)
}

func (suite *ConditionalTestSuite) TestConditionalNotEq() {
	notEq := NotEq(suite.country, "USA")

	sql := notEq.Accept(suite.ctx)
	bindings := suite.ctx.Binds()

	assert.Equal(suite.T(), "country != ?", sql)
	assert.Equal(suite.T(), []interface{}{"USA"}, bindings)
}

func (suite *ConditionalTestSuite) TestConditionalEq() {
	eq := Eq(suite.country, "Turkey")

	sql := eq.Accept(suite.ctx)
	bindings := suite.ctx.Binds()

	assert.Equal(suite.T(), "country = ?", sql)
	assert.Equal(suite.T(), []interface{}{"Turkey"}, bindings)
}

func (suite *ConditionalTestSuite) TestConditionalGt() {
	gt := Gt(suite.score, 1500)

	sql := gt.Accept(suite.ctx)
	bindings := suite.ctx.Binds()

	assert.Equal(suite.T(), "score > ?", sql)
	assert.Equal(suite.T(), []interface{}{1500}, bindings)
}

func (suite *ConditionalTestSuite) TestConditionalLt() {
	lt := Lt(suite.score, 1500)

	sql := lt.Accept(suite.ctx)
	bindings := suite.ctx.Binds()

	assert.Equal(suite.T(), "score < ?", sql)
	assert.Equal(suite.T(), []interface{}{1500}, bindings)
}

func (suite *ConditionalTestSuite) TestConditionalGte() {
	gte := Gte(suite.score, 1500)

	sql := gte.Accept(suite.ctx)
	bindings := suite.ctx.Binds()

	assert.Equal(suite.T(), "score >= ?", sql)
	assert.Equal(suite.T(), []interface{}{1500}, bindings)
}

func (suite *ConditionalTestSuite) TestConditionalLte() {
	lte := Lte(suite.score, 1500)

	sql := lte.Accept(suite.ctx)
	bindings := suite.ctx.Binds()

	assert.Equal(suite.T(), "score <= ?", sql)
	assert.Equal(suite.T(), []interface{}{1500}, bindings)
}

func TestConditionalTestSuite(t *testing.T) {
	suite.Run(t, new(ConditionalTestSuite))
}
