package qb

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SelectTestSuite struct {
	suite.Suite
	sqlite   Dialect
	mysql    Dialect
	postgres Dialect
	users    TableElem
	sessions TableElem
}

func (suite *SelectTestSuite) SetupTest() {
	suite.sqlite = NewDialect("sqlite3")
	suite.sqlite.SetEscaping(true)
	suite.mysql = NewDialect("mysql")
	suite.mysql.SetEscaping(true)
	suite.postgres = NewDialect("postgres")
	suite.postgres.SetEscaping(true)

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

	var statement *Stmt
	statement = sel.Build(suite.sqlite)
	assert.Equal(suite.T(), "SELECT id\nFROM users;", statement.SQL())

	statement = sel.Build(suite.mysql)
	assert.Equal(suite.T(), "SELECT `id`\nFROM `users`;", statement.SQL())

	statement = sel.Build(suite.postgres)
	assert.Equal(suite.T(), "SELECT \"id\"\nFROM \"users\";", statement.SQL())
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

	var statement *Stmt

	statement = sel.Build(suite.sqlite)
	assert.Equal(suite.T(), "SELECT id\nFROM users\nWHERE (email = ? AND id != ?);", statement.SQL())
	assert.Equal(suite.T(), []interface{}{"al@pacino.com", 5}, statement.Bindings())

	statement = sel.Build(suite.mysql)
	assert.Equal(suite.T(), "SELECT `id`\nFROM `users`\nWHERE (`email` = ? AND `id` != ?);", statement.SQL())
	assert.Equal(suite.T(), []interface{}{"al@pacino.com", 5}, statement.Bindings())

	statement = sel.Build(suite.postgres)
	assert.Equal(suite.T(), "SELECT \"id\"\nFROM \"users\"\nWHERE (\"email\" = $1 AND \"id\" != $2);", statement.SQL())
	assert.Equal(suite.T(), []interface{}{"al@pacino.com", 5}, statement.Bindings())
}

func (suite *SelectTestSuite) TestOrderByLimit() {
	selOrderByDesc := Select(suite.sessions.C("id")).
		From(suite.sessions).
		Where(Eq(suite.sessions.C("user_id"), 5)).
		OrderBy(suite.sessions.C("id")).Desc().
		Limit(0, 20)

	var statement *Stmt
	statement = selOrderByDesc.Build(suite.sqlite)
	assert.Equal(suite.T(), "SELECT id\nFROM sessions\nWHERE user_id = ?\nORDER BY id DESC\nLIMIT 20 OFFSET 0;", statement.SQL())
	assert.Equal(suite.T(), []interface{}{5}, statement.Bindings())

	statement = selOrderByDesc.Build(suite.mysql)
	assert.Equal(suite.T(), "SELECT `id`\nFROM `sessions`\nWHERE `user_id` = ?\nORDER BY `id` DESC\nLIMIT 20 OFFSET 0;", statement.SQL())
	assert.Equal(suite.T(), []interface{}{5}, statement.Bindings())

	statement = selOrderByDesc.Build(suite.postgres)
	assert.Equal(suite.T(), "SELECT \"id\"\nFROM \"sessions\"\nWHERE \"user_id\" = $1\nORDER BY \"id\" DESC\nLIMIT 20 OFFSET 0;", statement.SQL())
	assert.Equal(suite.T(), []interface{}{5}, statement.Bindings())

	selWithoutOrder := Select(suite.sessions.C("id")).
		From(suite.sessions).
		Where(Eq(suite.sessions.C("user_id"), 5)).
		OrderBy(suite.sessions.C("id"))

	statement = selWithoutOrder.Build(suite.sqlite)
	assert.Equal(suite.T(), "SELECT id\nFROM sessions\nWHERE user_id = ?\nORDER BY id ASC;", statement.SQL())
	assert.Equal(suite.T(), []interface{}{5}, statement.Bindings())

	statement = selWithoutOrder.Build(suite.mysql)
	assert.Equal(suite.T(), "SELECT `id`\nFROM `sessions`\nWHERE `user_id` = ?\nORDER BY `id` ASC;", statement.SQL())
	assert.Equal(suite.T(), []interface{}{5}, statement.Bindings())

	statement = selWithoutOrder.Build(suite.postgres)
	assert.Equal(suite.T(), "SELECT \"id\"\nFROM \"sessions\"\nWHERE \"user_id\" = $1\nORDER BY \"id\" ASC;", statement.SQL())
	assert.Equal(suite.T(), []interface{}{5}, statement.Bindings())

	selOrderByAsc := Select(suite.sessions.C("id")).
		From(suite.sessions).
		Where(Eq(suite.sessions.C("user_id"), 5)).
		OrderBy(suite.sessions.C("id")).Asc()

	statement = selOrderByAsc.Build(suite.sqlite)
	assert.Equal(suite.T(), "SELECT id\nFROM sessions\nWHERE user_id = ?\nORDER BY id ASC;", statement.SQL())
	assert.Equal(suite.T(), []interface{}{5}, statement.Bindings())

	statement = selOrderByAsc.Build(suite.mysql)
	assert.Equal(suite.T(), "SELECT `id`\nFROM `sessions`\nWHERE `user_id` = ?\nORDER BY `id` ASC;", statement.SQL())
	assert.Equal(suite.T(), []interface{}{5}, statement.Bindings())

	statement = selOrderByAsc.Build(suite.postgres)
	assert.Equal(suite.T(), "SELECT \"id\"\nFROM \"sessions\"\nWHERE \"user_id\" = $1\nORDER BY \"id\" ASC;", statement.SQL())
	assert.Equal(suite.T(), []interface{}{5}, statement.Bindings())
}

func (suite *SelectTestSuite) TestJoin() {

	// inner join
	selInnerJoin := Select(suite.sessions.C("id"), suite.sessions.C("auth_token")).
		From(suite.sessions).
		InnerJoin(suite.users, suite.sessions.C("user_id"), suite.users.C("id")).
		Where(Eq(suite.sessions.C("user_id"), 5))

	assert.Equal(suite.T(), suite.sessions.C("user_id"), selInnerJoin.from.C("user_id"))
	assert.Panics(suite.T(), func() { selInnerJoin.from.C("invalid") })

	assert.Equal(suite.T(), len(suite.sessions.All())+len(suite.users.All()), len(selInnerJoin.from.All()))

	var statement *Stmt

	statement = selInnerJoin.Build(suite.sqlite)
	assert.Equal(suite.T(), "SELECT sessions.id, sessions.auth_token\nFROM sessions\nINNER JOIN users ON sessions.user_id = users.id\nWHERE sessions.user_id = ?;", statement.SQL())
	assert.Equal(suite.T(), []interface{}{5}, statement.Bindings())

	statement = selInnerJoin.Build(suite.mysql)
	assert.Equal(suite.T(), "SELECT `sessions`.`id`, `sessions`.`auth_token`\nFROM `sessions`\nINNER JOIN `users` ON `sessions`.`user_id` = `users`.`id`\nWHERE `sessions`.`user_id` = ?;", statement.SQL())
	assert.Equal(suite.T(), []interface{}{5}, statement.Bindings())

	statement = selInnerJoin.Build(suite.postgres)
	assert.Equal(suite.T(), "SELECT \"sessions\".\"id\", \"sessions\".\"auth_token\"\nFROM \"sessions\"\nINNER JOIN \"users\" ON \"sessions\".\"user_id\" = \"users\".\"id\"\nWHERE \"sessions\".\"user_id\" = $1;", statement.SQL())
	assert.Equal(suite.T(), []interface{}{5}, statement.Bindings())

	// left join
	selLeftJoin := Select(suite.sessions.C("id"), suite.sessions.C("auth_token")).
		From(suite.sessions).
		LeftJoin(suite.users, suite.sessions.C("user_id"), suite.users.C("id")).
		Where(Eq(suite.sessions.C("user_id"), 5))

	statement = selLeftJoin.Build(suite.sqlite)
	assert.Equal(suite.T(), "SELECT sessions.id, sessions.auth_token\nFROM sessions\nLEFT OUTER JOIN users ON sessions.user_id = users.id\nWHERE sessions.user_id = ?;", statement.SQL())
	assert.Equal(suite.T(), []interface{}{5}, statement.Bindings())

	statement = selLeftJoin.Build(suite.mysql)
	assert.Equal(suite.T(), "SELECT `sessions`.`id`, `sessions`.`auth_token`\nFROM `sessions`\nLEFT OUTER JOIN `users` ON `sessions`.`user_id` = `users`.`id`\nWHERE `sessions`.`user_id` = ?;", statement.SQL())
	assert.Equal(suite.T(), []interface{}{5}, statement.Bindings())

	statement = selLeftJoin.Build(suite.postgres)
	assert.Equal(suite.T(), "SELECT \"sessions\".\"id\", \"sessions\".\"auth_token\"\nFROM \"sessions\"\nLEFT OUTER JOIN \"users\" ON \"sessions\".\"user_id\" = \"users\".\"id\"\nWHERE \"sessions\".\"user_id\" = $1;", statement.SQL())
	assert.Equal(suite.T(), []interface{}{5}, statement.Bindings())

	// right join
	selRightJoin := Select(suite.sessions.C("id")).
		From(suite.sessions).
		RightJoin(suite.users, suite.sessions.C("user_id"), suite.users.C("id")).
		Where(Eq(suite.sessions.C("user_id"), 5))

	statement = selRightJoin.Build(suite.sqlite)
	assert.Equal(suite.T(), "SELECT sessions.id\nFROM sessions\nRIGHT OUTER JOIN users ON sessions.user_id = users.id\nWHERE sessions.user_id = ?;", statement.SQL())
	assert.Equal(suite.T(), []interface{}{5}, statement.Bindings())

	statement = selRightJoin.Build(suite.mysql)
	assert.Equal(suite.T(), "SELECT `sessions`.`id`\nFROM `sessions`\nRIGHT OUTER JOIN `users` ON `sessions`.`user_id` = `users`.`id`\nWHERE `sessions`.`user_id` = ?;", statement.SQL())
	assert.Equal(suite.T(), []interface{}{5}, statement.Bindings())

	statement = selRightJoin.Build(suite.postgres)
	assert.Equal(suite.T(), "SELECT \"sessions\".\"id\"\nFROM \"sessions\"\nRIGHT OUTER JOIN \"users\" ON \"sessions\".\"user_id\" = \"users\".\"id\"\nWHERE \"sessions\".\"user_id\" = $1;", statement.SQL())
	assert.Equal(suite.T(), []interface{}{5}, statement.Bindings())

	// cross join
	selCrossJoin := Select(suite.sessions.C("id")).
		From(suite.sessions).
		CrossJoin(suite.users).
		Where(Eq(suite.sessions.C("user_id"), 5))

	statement = selCrossJoin.Build(suite.sqlite)
	assert.Equal(suite.T(), "SELECT sessions.id\nFROM sessions\nCROSS JOIN users\nWHERE sessions.user_id = ?;", statement.SQL())
	assert.Equal(suite.T(), []interface{}{5}, statement.Bindings())

	statement = selCrossJoin.Build(suite.mysql)
	assert.Equal(suite.T(), "SELECT `sessions`.`id`\nFROM `sessions`\nCROSS JOIN `users`\nWHERE `sessions`.`user_id` = ?;", statement.SQL())
	assert.Equal(suite.T(), []interface{}{5}, statement.Bindings())

	statement = selCrossJoin.Build(suite.postgres)
	assert.Equal(suite.T(), "SELECT \"sessions\".\"id\"\nFROM \"sessions\"\nCROSS JOIN \"users\"\nWHERE \"sessions\".\"user_id\" = $1;", statement.SQL())
	assert.Equal(suite.T(), []interface{}{5}, statement.Bindings())
}

func (suite *SelectTestSuite) TestGroupByHaving() {
	sel := Select(Count(suite.sessions.C("id"))).
		From(suite.sessions).
		GroupBy(suite.sessions.C("user_id")).
		Having(Sum(suite.sessions.C("id")), ">", 4)

	var statement *Stmt
	statement = sel.Build(suite.sqlite)
	assert.Equal(suite.T(), "SELECT COUNT(id)\nFROM sessions\nGROUP BY user_id\nHAVING SUM(id) > ?;", statement.SQL())
	assert.Equal(suite.T(), []interface{}{4}, statement.Bindings())

	statement = sel.Build(suite.mysql)
	assert.Equal(suite.T(), "SELECT COUNT(`id`)\nFROM `sessions`\nGROUP BY `user_id`\nHAVING SUM(`id`) > ?;", statement.SQL())
	assert.Equal(suite.T(), []interface{}{4}, statement.Bindings())

	statement = sel.Build(suite.postgres)
	assert.Equal(suite.T(), "SELECT COUNT(\"id\")\nFROM \"sessions\"\nGROUP BY \"user_id\"\nHAVING SUM(\"id\") > $1;", statement.SQL())
	assert.Equal(suite.T(), []interface{}{4}, statement.Bindings())
}

func (suite *SelectTestSuite) TestAlias() {
	sessionA := Alias("newname", suite.sessions)
	sel := Select(sessionA.C("id")).From(sessionA)
	st := sel.Build(suite.sqlite)
	assert.Equal(suite.T(), "SELECT id\nFROM sessions AS newname;", st.SQL())

	sel = Select(sessionA.All()...).From(sessionA)
	sql := sel.Build(suite.mysql).SQL()
	assert.Contains(suite.T(), sql, "`id`", st.SQL())
	assert.Contains(suite.T(), sql, "`user_id`", st.SQL())
	assert.Contains(suite.T(), sql, "`auth_token`", st.SQL())

	usersA := Alias("u", suite.users)
	sel = Select(usersA.C("email")).
		From(usersA).
		LeftJoin(sessionA, usersA.C("id"), sessionA.C("user_id")).
		Where(sessionA.C("auth_token").Eq("42"))
	st = sel.Build(suite.postgres)
	assert.Equal(suite.T(), `SELECT "u"."email"
FROM "users" AS "u"
LEFT OUTER JOIN "sessions" AS "newname" ON "u"."id" = "newname"."user_id"
WHERE "newname"."auth_token" = $1;`, st.SQL())
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
		GuessJoinOnClause(t1, t2)
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
