package qb

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/stretchr/testify/assert"
)

type UpdateTestSuite struct {
	suite.Suite
	dialect Dialect
	ctx     *CompilerContext
	users   TableElem
}

func (suite *UpdateTestSuite) SetupTest() {
	suite.dialect = NewDefaultDialect()
	suite.ctx = NewCompilerContext(suite.dialect)
	suite.users = Table(
		"users",
		Column("id", BigInt()).NotNull(),
		Column("email", Varchar()).NotNull().Unique(),
		PrimaryKey("email"),
	)
}

func (suite *UpdateTestSuite) TestUpdateSimple() {

	sql := Update(suite.users).
		Values(map[string]interface{}{
			"email": "robert@de.niro",
		}).Accept(suite.ctx)

	binds := suite.ctx.Binds

	assert.Contains(suite.T(), sql, "UPDATE users")
	assert.Contains(suite.T(), sql, "SET email = ?")
	assert.Equal(suite.T(), []interface{}{"robert@de.niro"}, binds)
}

func (suite *UpdateTestSuite) TestUpdateWhereReturning() {
	sql := Update(suite.users).
		Values(map[string]interface{}{"email": "robert@de.niro"}).
		Where(Eq(suite.users.C("email"), "al@pacino")).
		Returning(suite.users.C("id"), suite.users.C("email")).
		Accept(suite.ctx)
	binds := suite.ctx.Binds

	assert.Contains(suite.T(), sql, "UPDATE users")
	assert.Contains(suite.T(), sql, "SET email = ?")
	assert.Contains(suite.T(), sql, "WHERE email = ?")
	assert.Contains(suite.T(), sql, "RETURNING id, email")
	assert.Equal(suite.T(), []interface{}{
		"robert@de.niro",
		"al@pacino",
	}, binds)
}

func TestUpdateTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateTestSuite))
}
