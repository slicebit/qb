package qb

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"
)

type CombinerTestSuite struct {
	suite.Suite
	dialect Dialect
	ctx     Context
}

func (suite *CombinerTestSuite) SetupTest() {
	suite.dialect = NewDefaultDialect()
	suite.ctx = NewCompilerContext(suite.dialect)
}

func (suite *CombinerTestSuite) TestCombinerAnd() {
	email := Column("email", Varchar()).NotNull().Unique()
	id := Column("id", Int()).NotNull()

	and := And(Eq(email, "al@pacino.com"), NotEq(id, 1))
	sql := and.Accept(suite.ctx)
	binds := suite.ctx.Binds()

	assert.Equal(suite.T(), "(email = ? AND id != ?)", sql)
	assert.Equal(suite.T(), []interface{}{"al@pacino.com", 1}, binds)
}

func (suite *CombinerTestSuite) TestCombinerOr() {
	email := Column("email", Varchar()).NotNull().Unique()
	id := Column("id", Int()).NotNull()

	or := Or(Eq(email, "al@pacino.com"), NotEq(id, 1))
	sql := or.Accept(suite.ctx)
	binds := suite.ctx.Binds()

	assert.Equal(suite.T(), "(email = ? OR id != ?)", sql)
	assert.Equal(suite.T(), []interface{}{"al@pacino.com", 1}, binds)
}

func TestCombinerTestSuite(t *testing.T) {
	suite.Run(t, new(CombinerTestSuite))
}
