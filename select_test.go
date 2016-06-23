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
		Column("email", Varchar().NotNull().Unique()),
		Column("password", Varchar().NotNull()),
		PrimaryKey("id"),
	)

	suite.sessions = Table(
		"sessions",
		Column("id", BigInt()),
		Column("user_id", BigInt()),
		Column("auth_token", Varchar().Size(36).Unique().NotNull()),
		PrimaryKey("id"),
		ForeignKey().Ref("user_id", "users", "id"),
	)
}

func (suite *SelectTestSuite) TestSimpleSelect() {
	sel := Select(suite.users.C("id")).From(suite.users)

	var statement *Stmt
	statement = sel.Build(suite.sqlite)
	assert.Equal(suite.T(), statement.SQL(), "SELECT id\nFROM users;")

	statement = sel.Build(suite.mysql)
	assert.Equal(suite.T(), statement.SQL(), "SELECT `id`\nFROM `users`;")

	statement = sel.Build(suite.postgres)
	assert.Equal(suite.T(), statement.SQL(), "SELECT \"id\"\nFROM \"users\";")
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
	assert.Equal(suite.T(), statement.SQL(), "SELECT id\nFROM users\nWHERE (users.email = ?) AND (users.id != ?);")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{"al@pacino.com", 5})

	statement = sel.Build(suite.mysql)
	assert.Equal(suite.T(), statement.SQL(), "SELECT `id`\nFROM `users`\nWHERE (`users`.`email` = ?) AND (`users`.`id` != ?);")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{"al@pacino.com", 5})

	statement = sel.Build(suite.postgres)
	assert.Equal(suite.T(), statement.SQL(), "SELECT \"id\"\nFROM \"users\"\nWHERE (\"users\".\"email\" = $1) AND (\"users\".\"id\" != $2);")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{"al@pacino.com", 5})
}

func (suite *SelectTestSuite) TestOrderByLimit() {
	selOrderByDesc := Select(suite.sessions.C("id")).
		From(suite.sessions).
		Where(Eq(suite.sessions.C("user_id"), 5)).
		OrderBy(suite.sessions.C("id")).Desc().
		Limit(0, 20)

	var statement *Stmt
	statement = selOrderByDesc.Build(suite.sqlite)
	assert.Equal(suite.T(), statement.SQL(), "SELECT id\nFROM sessions\nWHERE (sessions.user_id = ?)\nORDER BY id DESC\nLIMIT 20 OFFSET 0;")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{5})

	statement = selOrderByDesc.Build(suite.mysql)
	assert.Equal(suite.T(), statement.SQL(), "SELECT `id`\nFROM `sessions`\nWHERE (`sessions`.`user_id` = ?)\nORDER BY `id` DESC\nLIMIT 20 OFFSET 0;")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{5})

	statement = selOrderByDesc.Build(suite.postgres)
	assert.Equal(suite.T(), statement.SQL(), "SELECT \"id\"\nFROM \"sessions\"\nWHERE (\"sessions\".\"user_id\" = $1)\nORDER BY \"id\" DESC\nLIMIT 20 OFFSET 0;")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{5})

	selWithoutOrder := Select(suite.sessions.C("id")).
		From(suite.sessions).
		Where(Eq(suite.sessions.C("user_id"), 5)).
		OrderBy(suite.sessions.C("id"))

	statement = selWithoutOrder.Build(suite.sqlite)
	assert.Equal(suite.T(), statement.SQL(), "SELECT id\nFROM sessions\nWHERE (sessions.user_id = ?)\nORDER BY id ASC;")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{5})

	statement = selWithoutOrder.Build(suite.mysql)
	assert.Equal(suite.T(), statement.SQL(), "SELECT `id`\nFROM `sessions`\nWHERE (`sessions`.`user_id` = ?)\nORDER BY `id` ASC;")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{5})

	statement = selWithoutOrder.Build(suite.postgres)
	assert.Equal(suite.T(), statement.SQL(), "SELECT \"id\"\nFROM \"sessions\"\nWHERE (\"sessions\".\"user_id\" = $1)\nORDER BY \"id\" ASC;")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{5})

	selOrderByAsc := Select(suite.sessions.C("id")).
		From(suite.sessions).
		Where(Eq(suite.sessions.C("user_id"), 5)).
		OrderBy(suite.sessions.C("id")).Asc()

	statement = selOrderByAsc.Build(suite.sqlite)
	assert.Equal(suite.T(), statement.SQL(), "SELECT id\nFROM sessions\nWHERE (sessions.user_id = ?)\nORDER BY id ASC;")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{5})

	statement = selOrderByAsc.Build(suite.mysql)
	assert.Equal(suite.T(), statement.SQL(), "SELECT `id`\nFROM `sessions`\nWHERE (`sessions`.`user_id` = ?)\nORDER BY `id` ASC;")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{5})

	statement = selOrderByAsc.Build(suite.postgres)
	assert.Equal(suite.T(), statement.SQL(), "SELECT \"id\"\nFROM \"sessions\"\nWHERE (\"sessions\".\"user_id\" = $1)\nORDER BY \"id\" ASC;")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{5})
}

func (suite *SelectTestSuite) TestJoin() {

	// inner join
	selInnerJoin := Select(suite.sessions.C("id"), suite.sessions.C("auth_token")).
		From(suite.sessions).
		InnerJoin(suite.users, suite.sessions.C("user_id"), suite.users.C("id")).
		Where(Eq(suite.sessions.C("user_id"), 5))

	var statement *Stmt

	statement = selInnerJoin.Build(suite.sqlite)
	assert.Equal(suite.T(), statement.SQL(), "SELECT sessions.id, sessions.auth_token\nFROM sessions\nINNER JOIN users ON sessions.user_id = users.id\nWHERE (sessions.user_id = ?);")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{5})

	statement = selInnerJoin.Build(suite.mysql)
	assert.Equal(suite.T(), statement.SQL(), "SELECT `sessions`.`id`, `sessions`.`auth_token`\nFROM `sessions`\nINNER JOIN `users` ON `sessions`.`user_id` = `users`.`id`\nWHERE (`sessions`.`user_id` = ?);")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{5})

	statement = selInnerJoin.Build(suite.postgres)
	assert.Equal(suite.T(), statement.SQL(), "SELECT \"sessions\".\"id\", \"sessions\".\"auth_token\"\nFROM \"sessions\"\nINNER JOIN \"users\" ON \"sessions\".\"user_id\" = \"users\".\"id\"\nWHERE (\"sessions\".\"user_id\" = $1);")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{5})

	// left join
	selLeftJoin := Select(suite.sessions.C("id"), suite.sessions.C("auth_token")).
		From(suite.sessions).
		LeftJoin(suite.users, suite.sessions.C("user_id"), suite.users.C("id")).
		Where(Eq(suite.sessions.C("user_id"), 5))

	statement = selLeftJoin.Build(suite.sqlite)
	assert.Equal(suite.T(), statement.SQL(), "SELECT sessions.id, sessions.auth_token\nFROM sessions\nLEFT OUTER JOIN users ON sessions.user_id = users.id\nWHERE (sessions.user_id = ?);")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{5})

	statement = selLeftJoin.Build(suite.mysql)
	assert.Equal(suite.T(), statement.SQL(), "SELECT `sessions`.`id`, `sessions`.`auth_token`\nFROM `sessions`\nLEFT OUTER JOIN `users` ON `sessions`.`user_id` = `users`.`id`\nWHERE (`sessions`.`user_id` = ?);")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{5})

	statement = selLeftJoin.Build(suite.postgres)
	assert.Equal(suite.T(), statement.SQL(), "SELECT \"sessions\".\"id\", \"sessions\".\"auth_token\"\nFROM \"sessions\"\nLEFT OUTER JOIN \"users\" ON \"sessions\".\"user_id\" = \"users\".\"id\"\nWHERE (\"sessions\".\"user_id\" = $1);")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{5})

	// right join
	selRightJoin := Select(suite.sessions.C("id")).
		From(suite.sessions).
		RightJoin(suite.users, suite.sessions.C("user_id"), suite.users.C("id")).
		Where(Eq(suite.sessions.C("user_id"), 5))

	statement = selRightJoin.Build(suite.sqlite)
	assert.Equal(suite.T(), statement.SQL(), "SELECT sessions.id\nFROM sessions\nRIGHT OUTER JOIN users ON sessions.user_id = users.id\nWHERE (sessions.user_id = ?);")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{5})

	statement = selRightJoin.Build(suite.mysql)
	assert.Equal(suite.T(), statement.SQL(), "SELECT `sessions`.`id`\nFROM `sessions`\nRIGHT OUTER JOIN `users` ON `sessions`.`user_id` = `users`.`id`\nWHERE (`sessions`.`user_id` = ?);")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{5})

	statement = selRightJoin.Build(suite.postgres)
	assert.Equal(suite.T(), statement.SQL(), "SELECT \"sessions\".\"id\"\nFROM \"sessions\"\nRIGHT OUTER JOIN \"users\" ON \"sessions\".\"user_id\" = \"users\".\"id\"\nWHERE (\"sessions\".\"user_id\" = $1);")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{5})

	// cross join
	selCrossJoin := Select(suite.sessions.C("id")).
		From(suite.sessions).
		CrossJoin(suite.users).
		Where(Eq(suite.sessions.C("user_id"), 5))

	statement = selCrossJoin.Build(suite.sqlite)
	assert.Equal(suite.T(), statement.SQL(), "SELECT sessions.id\nFROM sessions\nCROSS JOIN users\nWHERE (sessions.user_id = ?);")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{5})

	statement = selCrossJoin.Build(suite.mysql)
	assert.Equal(suite.T(), statement.SQL(), "SELECT `sessions`.`id`\nFROM `sessions`\nCROSS JOIN `users`\nWHERE (`sessions`.`user_id` = ?);")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{5})

	statement = selCrossJoin.Build(suite.postgres)
	assert.Equal(suite.T(), statement.SQL(), "SELECT \"sessions\".\"id\"\nFROM \"sessions\"\nCROSS JOIN \"users\"\nWHERE (\"sessions\".\"user_id\" = $1);")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{5})
}

func (suite *SelectTestSuite) TestGroupByHaving() {
	sel := Select(Count(suite.sessions.C("id"))).
		From(suite.sessions).
		GroupBy(suite.sessions.C("user_id")).
		Having(Sum(suite.sessions.C("id")), ">", 4)

	var statement *Stmt
	statement = sel.Build(suite.sqlite)
	assert.Equal(suite.T(), statement.SQL(), "SELECT COUNT(id)\nFROM sessions\nGROUP BY user_id\nHAVING SUM(id) > ?;")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{4})

	statement = sel.Build(suite.mysql)
	assert.Equal(suite.T(), statement.SQL(), "SELECT COUNT(`id`)\nFROM `sessions`\nGROUP BY `user_id`\nHAVING SUM(`id`) > ?;")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{4})

	statement = sel.Build(suite.postgres)
	assert.Equal(suite.T(), statement.SQL(), "SELECT COUNT(\"id\")\nFROM \"sessions\"\nGROUP BY \"user_id\"\nHAVING SUM(\"id\") > $1;")
	assert.Equal(suite.T(), statement.Bindings(), []interface{}{4})
}

func TestSelectTestSuite(t *testing.T) {
	suite.Run(t, new(SelectTestSuite))
}
