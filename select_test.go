package qb

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SelectTestSuite struct {
	suite.Suite
	users    TableElem
	sessions TableElem
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
}

func (suite *SelectTestSuite) TestSimpleSelect() {
	sel := Select(suite.users.C("id")).From(suite.users)
	assert.Equal(suite.T(), "SELECT id\nFROM users", asDefSQL(sel))

	sel = sel.Select(Count(suite.users.C("id")))
	assert.Equal(suite.T(), "SELECT COUNT(id)\nFROM users", asDefSQL(sel))
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

	sql, binds := asDefSQLBinds(sel)
	assert.Equal(suite.T(), "SELECT id\nFROM users\nWHERE (email = ? AND id != ?)", sql)
	assert.Equal(suite.T(), []interface{}{"al@pacino.com", 5}, binds)
}

func (suite *SelectTestSuite) TestOrderByLimit() {
	selOrderByDesc := Select(suite.sessions.C("id")).
		From(suite.sessions).
		Where(Eq(suite.sessions.C("user_id"), 5)).
		OrderBy(suite.sessions.C("id")).Desc().
		Limit(20)

	sql, binds := asDefSQLBinds(selOrderByDesc)
	assert.Equal(suite.T(), "SELECT id\nFROM sessions\nWHERE user_id = ?\nORDER BY id DESC\nLIMIT 20", sql)
	assert.Equal(suite.T(), []interface{}{5}, binds)

	selWithoutOrder := Select(suite.sessions.C("id")).
		From(suite.sessions).
		Where(Eq(suite.sessions.C("user_id"), 5)).
		OrderBy(suite.sessions.C("id")).
		Offset(12)

	sql, binds = asDefSQLBinds(selWithoutOrder)
	assert.Equal(suite.T(), "SELECT id\nFROM sessions\nWHERE user_id = ?\nORDER BY id ASC\nOFFSET 12", sql)
	assert.Equal(suite.T(), []interface{}{5}, binds)

	selOrderByAsc := Select(suite.sessions.C("id")).
		From(suite.sessions).
		Where(Eq(suite.sessions.C("user_id"), 5)).
		OrderBy(suite.sessions.C("id")).Asc().
		LimitOffset(20, 12)

	sql, binds = asDefSQLBinds(selOrderByAsc)
	assert.Equal(suite.T(), "SELECT id\nFROM sessions\nWHERE user_id = ?\nORDER BY id ASC\nLIMIT 20 OFFSET 12", sql)
	assert.Equal(suite.T(), []interface{}{5}, binds)
}

func (suite *SelectTestSuite) TestJoin() {

	// inner join
	selInnerJoin := Select(suite.sessions.C("id"), suite.sessions.C("auth_token")).
		From(suite.sessions).
		InnerJoin(suite.users, suite.sessions.C("user_id"), suite.users.C("id")).
		Where(Eq(suite.sessions.C("user_id"), 5))

	assert.Equal(suite.T(), suite.sessions.C("user_id"), selInnerJoin.FromClause.C("user_id"))
	assert.Panics(suite.T(), func() { selInnerJoin.FromClause.C("invalid") })

	assert.Equal(suite.T(), len(suite.sessions.All())+len(suite.users.All()), len(selInnerJoin.FromClause.All()))

	sql, binds := asDefSQLBinds(selInnerJoin)

	assert.Equal(suite.T(), "SELECT sessions.id, sessions.auth_token\nFROM sessions\nINNER JOIN users ON sessions.user_id = users.id\nWHERE sessions.user_id = ?", sql)
	assert.Equal(suite.T(), []interface{}{5}, binds)

	// left join
	selLeftJoin := Select(suite.sessions.C("id"), suite.sessions.C("auth_token")).
		From(suite.sessions).
		LeftJoin(suite.users, suite.sessions.C("user_id"), suite.users.C("id")).
		Where(Eq(suite.sessions.C("user_id"), 5))

	sql, binds = asDefSQLBinds(selLeftJoin)

	assert.Equal(suite.T(), "SELECT sessions.id, sessions.auth_token\nFROM sessions\nLEFT OUTER JOIN users ON sessions.user_id = users.id\nWHERE sessions.user_id = ?", sql)
	assert.Equal(suite.T(), []interface{}{5}, binds)

	// right join
	selRightJoin := Select(suite.sessions.C("id")).
		From(suite.sessions).
		RightJoin(suite.users, suite.sessions.C("user_id"), suite.users.C("id")).
		Where(Eq(suite.sessions.C("user_id"), 5))

	sql, binds = asDefSQLBinds(selRightJoin)
	assert.Equal(suite.T(), "SELECT sessions.id\nFROM sessions\nRIGHT OUTER JOIN users ON sessions.user_id = users.id\nWHERE sessions.user_id = ?", sql)
	assert.Equal(suite.T(), []interface{}{5}, binds)

	// cross join
	selCrossJoin := Select(suite.sessions.C("id")).
		From(suite.sessions).
		CrossJoin(suite.users).
		Where(Eq(suite.sessions.C("user_id"), 5))

	sql, binds = asDefSQLBinds(selCrossJoin)
	assert.Equal(suite.T(), "SELECT sessions.id\nFROM sessions\nCROSS JOIN users\nWHERE sessions.user_id = ?", sql)
	assert.Equal(suite.T(), []interface{}{5}, binds)
}

func (suite *SelectTestSuite) TestGroupByHaving() {
	sel := Select(Count(suite.sessions.C("id"))).
		From(suite.sessions).
		GroupBy(suite.sessions.C("user_id")).
		Having(Sum(suite.sessions.C("id")), ">", 4)

	sql, binds := asDefSQLBinds(sel)
	assert.Equal(suite.T(), "SELECT COUNT(id)\nFROM sessions\nGROUP BY user_id\nHAVING SUM(id) > ?", sql)
	assert.Equal(suite.T(), []interface{}{4}, binds)
}

func (suite *SelectTestSuite) TestAlias() {
	sessionA := Alias("newname", suite.sessions)
	sel := Select(sessionA.C("id")).From(sessionA)
	assert.Equal(suite.T(), "SELECT id\nFROM sessions AS newname", asDefSQL(sel))

	sel = Select(sessionA.All()...).From(sessionA)
	sql := asDefSQL(sel)
	assert.Contains(suite.T(), sql, "id", sql)
	assert.Contains(suite.T(), sql, "user_id", sql)
	assert.Contains(suite.T(), sql, "auth_token", sql)

	usersA := Alias("u", suite.users)
	sel = Select(usersA.C("email")).
		From(usersA).
		LeftJoin(sessionA, usersA.C("id"), sessionA.C("user_id")).
		Where(sessionA.C("auth_token").Eq("42"))
	sql = asDefSQL(sel)
	assert.Equal(suite.T(), `SELECT u.email
FROM users AS u
LEFT OUTER JOIN sessions AS newname ON u.id = newname.user_id
WHERE newname.auth_token = ?`, sql)
}

func (suite *SelectTestSuite) TestGuessJoinOnClause() {
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

	assert.Equal(suite.T(), "t3.c1 = t1.c1", asDefSQL(GuessJoinOnClause(t3, t1)))
	assert.Equal(suite.T(), "t3.c1 = t1.c1", asDefSQL(GuessJoinOnClause(t1, t3)))
	assert.Equal(suite.T(), "(t4.c1 = t1.c1 AND t4.c2 = t1.c2)", asDefSQL(GuessJoinOnClause(t4, t1)))

	assert.Panics(suite.T(), func() {
		GuessJoinOnClause(t2, t3)
	})
}

func (suite *SelectTestSuite) TestMakeJoinOnClause() {
	assert.Panics(suite.T(), func() {
		MakeJoinOnClause(TableElem{}, TableElem{}, And(), And(), And())
	})
}

func TestSelectTestSuite(t *testing.T) {
	suite.Run(t, new(SelectTestSuite))
}
