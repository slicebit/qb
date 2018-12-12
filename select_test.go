package qb

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SelectTestSuite struct {
	suite.Suite
	users    TableElem
	sessions TableElem
	ctx      *CompilerContext
	dialect  Dialect
}

func (suite *SelectTestSuite) SetupTest() {
	suite.users = Table(
		"users",
		Column("id", BigInt()),
		Column("email", Varchar()).NotNull().Unique(),
		Column("password", Varchar()).NotNull(),
		PrimaryKey("id"),
	)

	suite.sessions = Table(
		"sessions",
		Column("id", BigInt()),
		Column("user_id", BigInt()),
		Column("auth_token", Varchar().Size(36)).Unique().NotNull(),
		PrimaryKey("id"),
		ForeignKey("user_id").References("users", "id"),
	)

	suite.dialect = NewDefaultDialect()
	suite.ctx = NewCompilerContext(suite.dialect)
}

func (suite *SelectTestSuite) TestSelectSimple() {
	sel := Select(suite.users.C("id")).From(suite.users)
	assert.Equal(suite.T(), "SELECT id\nFROM users", sel.Accept(suite.ctx))
}

func (suite *SelectTestSuite) TestSelectAggregate() {
	sel := Select(suite.users.C("id")).From(suite.users)
	selCount := sel.Select(Count(suite.users.C("id")))
	assert.Equal(suite.T(), "SELECT COUNT(id)\nFROM users", selCount.Accept(suite.ctx))
}

func (suite *SelectTestSuite) TestSelectWhere() {
	sel := Select(suite.users.C("id")).
		From(suite.users).
		Where(
			And(
				Eq(suite.users.C("email"), "al@pacino.com"),
				NotEq(suite.users.C("id"), 5),
			),
		)

	sql := sel.Accept(suite.ctx)
	binds := suite.ctx.Binds

	assert.Equal(suite.T(), "SELECT id\nFROM users\nWHERE (email = ? AND id != ?)", sql)
	assert.Equal(suite.T(), []interface{}{"al@pacino.com", 5}, binds)
}

func (suite *SelectTestSuite) TestSelectOrderByLimit() {
	selOrderByDesc := Select(suite.sessions.C("id")).
		From(suite.sessions).
		Where(Eq(suite.sessions.C("user_id"), 5)).
		OrderBy(suite.sessions.C("id")).Desc().
		Limit(20)

	sql := selOrderByDesc.Accept(suite.ctx)
	binds := suite.ctx.Binds

	assert.Equal(suite.T(), "SELECT id\nFROM sessions\nWHERE user_id = ?\nORDER BY id DESC\nLIMIT 20", sql)
	assert.Equal(suite.T(), []interface{}{5}, binds)
}

func (suite *SelectTestSuite) TestSelectWithoutOrder() {
	selWithoutOrder := Select(suite.sessions.C("id")).
		From(suite.sessions).
		Where(Eq(suite.sessions.C("user_id"), 5)).
		OrderBy(suite.sessions.C("id")).
		Offset(12)

	sql := selWithoutOrder.Accept(suite.ctx)
	binds := suite.ctx.Binds

	assert.Equal(suite.T(), "SELECT id\nFROM sessions\nWHERE user_id = ?\nORDER BY id ASC\nOFFSET 12", sql)
	assert.Equal(suite.T(), []interface{}{5}, binds)
}

func (suite *SelectTestSuite) TestSelectOrderByAsc() {
	selOrderByAsc := Select(suite.sessions.C("id")).
		From(suite.sessions).
		Where(Eq(suite.sessions.C("user_id"), 5)).
		OrderBy(suite.sessions.C("id")).Asc().
		LimitOffset(20, 12)

	sql := selOrderByAsc.Accept(suite.ctx)
	binds := suite.ctx.Binds

	assert.Equal(suite.T(), "SELECT id\nFROM sessions\nWHERE user_id = ?\nORDER BY id ASC\nLIMIT 20 OFFSET 12", sql)
	assert.Equal(suite.T(), []interface{}{5}, binds)
}

func (suite *SelectTestSuite) TestSelectInnerJoin() {
	selInnerJoin := Select(suite.sessions.C("id"), suite.sessions.C("auth_token")).
		From(suite.sessions).
		InnerJoin(suite.users, suite.sessions.C("user_id"), suite.users.C("id")).
		Where(Eq(suite.sessions.C("user_id"), 5))

	sql := selInnerJoin.Accept(suite.ctx)
	binds := suite.ctx.Binds

	assert.Equal(suite.T(), suite.sessions.C("user_id"), selInnerJoin.FromClause.C("user_id"))
	assert.Panics(suite.T(), func() { selInnerJoin.FromClause.C("invalid") })
	assert.Equal(suite.T(), len(suite.sessions.All())+len(suite.users.All()), len(selInnerJoin.FromClause.All()))

	assert.Equal(suite.T(), "SELECT sessions.id, sessions.auth_token\nFROM sessions\nINNER JOIN users ON sessions.user_id = users.id\nWHERE sessions.user_id = ?", sql)
	assert.Equal(suite.T(), []interface{}{5}, binds)
}

func (suite *SelectTestSuite) TestSelectLeftJoin() {
	selLeftJoin := Select(suite.sessions.C("id"), suite.sessions.C("auth_token")).
		From(suite.sessions).
		LeftJoin(suite.users, suite.sessions.C("user_id"), suite.users.C("id")).
		Where(Eq(suite.sessions.C("user_id"), 5))

	sql := selLeftJoin.Accept(suite.ctx)
	binds := suite.ctx.Binds

	assert.Equal(suite.T(), "SELECT sessions.id, sessions.auth_token\nFROM sessions\nLEFT OUTER JOIN users ON sessions.user_id = users.id\nWHERE sessions.user_id = ?", sql)
	assert.Equal(suite.T(), []interface{}{5}, binds)
}

func (suite *SelectTestSuite) TestSelectRightJoin() {
	selRightJoin := Select(suite.sessions.C("id")).
		From(suite.sessions).
		RightJoin(suite.users, suite.sessions.C("user_id"), suite.users.C("id")).
		Where(Eq(suite.sessions.C("user_id"), 5))

	sql := selRightJoin.Accept(suite.ctx)
	binds := suite.ctx.Binds

	assert.Equal(suite.T(), "SELECT sessions.id\nFROM sessions\nRIGHT OUTER JOIN users ON sessions.user_id = users.id\nWHERE sessions.user_id = ?", sql)
	assert.Equal(suite.T(), []interface{}{5}, binds)
}

func (suite *SelectTestSuite) TestSelectCrossJoin() {
	selCrossJoin := Select(suite.sessions.C("id")).
		From(suite.sessions).
		CrossJoin(suite.users).
		Where(Eq(suite.sessions.C("user_id"), 5))

	sql := selCrossJoin.Accept(suite.ctx)
	binds := suite.ctx.Binds

	assert.Equal(suite.T(), "SELECT sessions.id\nFROM sessions\nCROSS JOIN users\nWHERE sessions.user_id = ?", sql)
	assert.Equal(suite.T(), []interface{}{5}, binds)
}

func (suite *SelectTestSuite) TestSelectGroupByHaving() {
	sel := Select(Count(suite.sessions.C("id"))).
		From(suite.sessions).
		GroupBy(suite.sessions.C("user_id")).
		Having(Sum(suite.sessions.C("id")), ">", 4)

	sql := sel.Accept(suite.ctx)
	binds := suite.ctx.Binds

	assert.Equal(suite.T(), "SELECT COUNT(id)\nFROM sessions\nGROUP BY user_id\nHAVING SUM(id) > ?", sql)
	assert.Equal(suite.T(), []interface{}{4}, binds)
}

func (suite *SelectTestSuite) TestSelectAliasFrom() {
	sessionAlias := Alias("newname", suite.sessions)

	sel := Select(sessionAlias.C("id")).From(sessionAlias)
	assert.Equal(suite.T(), "SELECT id\nFROM sessions AS newname", sel.Accept(suite.ctx))
}

func (suite *SelectTestSuite) TestSelectAliasAll() {
	sessionAlias := Alias("newname", suite.sessions)

	sel := Select(sessionAlias.C("id")).From(sessionAlias)
	sql := sel.Accept(suite.ctx)

	assert.Contains(suite.T(), sql, "id", sql)
	assert.Contains(suite.T(), sql, "sessions AS newname", sql)
}

func (suite *SelectTestSuite) TestSelectAliasWhereMultipleTable() {
	sessionAlias := Alias("newname", suite.sessions)
	usersAlias := Alias("u", suite.users)
	sel := Select(usersAlias.C("email")).
		From(usersAlias).
		LeftJoin(sessionAlias, usersAlias.C("id"), sessionAlias.C("user_id")).
		Where(sessionAlias.C("auth_token").Eq("42"))

	sql := sel.Accept(suite.ctx)
	expected := strings.Join([]string{
		"SELECT u.email",
		"FROM users AS u",
		"LEFT OUTER JOIN sessions AS newname ON u.id = newname.user_id",
		"WHERE newname.auth_token = ?",
	}, "\n")
	assert.Equal(suite.T(), expected, sql)
}

func (suite *SelectTestSuite) TestSelectGuessJoinOnClause() {
	t1 := Table(
		"t1",
		Column("c1", Int()),
		Column("c2", Int()),
	)
	t2 := Table(
		"t2",
		Column("c1", Int()),
		Column("c2", Int()),
	)
	t3 := Table(
		"t3",
		Column("c1", Int()),
		Column("c2", Int()),
		ForeignKey("c1").References("t1", "c1"),
		ForeignKey("c1").References("t2", "c1"),
		ForeignKey("c2").References("t2", "c2"),
	)
	t4 := Table(
		"t4",
		Column("c1", Int()),
		Column("c2", Int()),
		ForeignKey("c1", "c2").References("t1", "c1", "c2"),
	)

	assert.Panics(suite.T(), func() {
		GuessJoinOnClause(t1, Alias("tt", t3))
	})

	assert.Panics(suite.T(), func() {
		GuessJoinOnClause(Alias("tt", t3), t2)
	})

	assert.Panics(suite.T(), func() {
		GuessJoinOnClause(t1, &t2)
	})

	assert.Equal(suite.T(), "t3.c1 = t1.c1", GuessJoinOnClause(t3, t1).Accept(suite.ctx))
	assert.Equal(suite.T(), "t3.c1 = t1.c1", GuessJoinOnClause(t1, t3).Accept(suite.ctx))
	assert.Equal(suite.T(), "(t4.c1 = t1.c1 AND t4.c2 = t1.c2)", GuessJoinOnClause(t4, t1).Accept(suite.ctx))

	assert.Panics(suite.T(), func() {
		GuessJoinOnClause(t2, t3)
	})
}

func (suite *SelectTestSuite) TestSelectMakeJoinOnClause() {
	assert.Panics(suite.T(), func() {
		MakeJoinOnClause(TableElem{}, TableElem{}, And(), And(), And())
	})
}

func TestSelectTestSuite(t *testing.T) {
	suite.Run(t, new(SelectTestSuite))
}
