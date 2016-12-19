package qb

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TableTestSuite struct {
	suite.Suite
}

func (suite *TableTestSuite) TestTableSimpleCreateDrop() {
	dialect := NewDialect("default")
	usersTable := Table("users", Column("id", Varchar().Size(40)))
	assert.Equal(suite.T(), 1, len(usersTable.All()))

	ddl := usersTable.Create(dialect)
	assert.Equal(suite.T(), "CREATE TABLE users (\n\tid VARCHAR(40)\n);", ddl)

	statement := usersTable.Build(dialect)
	assert.Equal(suite.T(), "CREATE TABLE users (\n\tid VARCHAR(40)\n);", statement.SQL())
	assert.Equal(suite.T(), []interface{}{}, statement.Bindings())

	assert.Equal(suite.T(), "DROP TABLE users;", usersTable.Drop(dialect))
}

func (suite *TableTestSuite) TestTablePrimaryForeignKey() {
	usersTable := Table(
		"users",
		Column("id", Varchar().Size(40)),
		Column("session_id", Varchar().Size(40)),
		Column("auth_token", Varchar().Size(40)),
		Column("role_id", Varchar().Size(40)),
		PrimaryKey("id"),
		ForeignKey("session_id", "auth_token").
			References("sessions", "id", "auth_token"),
		ForeignKey("role_id").References("roles", "id"),
	)

	ddl := usersTable.Create(NewDialect("default"))
	assert.Contains(suite.T(), ddl, "CREATE TABLE users (")
	assert.Contains(suite.T(), ddl, "auth_token VARCHAR(40)")
	assert.Contains(suite.T(), ddl, "role_id VARCHAR(40)")
	assert.Contains(suite.T(), ddl, "id VARCHAR(40) PRIMARY KEY")
	assert.Contains(suite.T(), ddl, "session_id VARCHAR(40)")
	assert.Contains(suite.T(), ddl, "FOREIGN KEY(session_id, auth_token) REFERENCES sessions(id, auth_token)")
	assert.Contains(suite.T(), ddl, "FOREIGN KEY(role_id) REFERENCES roles(id)")
	assert.Contains(suite.T(), ddl, ");")
}

func (suite *TableTestSuite) TestTablePrimaryKey() {
	t := Table(
		"users",
		Column("id", Varchar().Size(40)).PrimaryKey(),
	)
	assert.Equal(suite.T(), []string{"id"}, t.PrimaryKeyConstraint.Columns)

	t = Table(
		"users",
		Column("fname", Varchar().Size(40)).PrimaryKey(),
		Column("lname", Varchar().Size(40)).PrimaryKey(),
	)

	assert.Equal(suite.T(), []string{"fname", "lname"}, t.PrimaryKeyConstraint.Columns)

	ddl := t.Create(NewDialect("default"))
	assert.Contains(suite.T(), ddl, "PRIMARY KEY(fname, lname)")

	assert.Panics(suite.T(), func() {
		Table(
			"users",
			Column("id", Varchar().Size(40)).PrimaryKey(),
			PrimaryKey("id"),
		)
	})
}

func (suite *TableTestSuite) TestTableUniqueCompositeUnique() {
	usersTable := Table(
		"users",
		Column("id", Varchar().Size(40)),
		Column("email", Varchar().Size(40)).Unique(),
		Column("device_id", Varchar().Size(255)).Unique(),
		UniqueKey("email", "device_id"),
	)

	ddl := usersTable.Create(NewDialect("default"))
	assert.Contains(suite.T(), ddl, "CREATE TABLE users (")
	assert.Contains(suite.T(), ddl, "id VARCHAR(40)")
	assert.Contains(suite.T(), ddl, "email VARCHAR(40) UNIQUE")
	assert.Contains(suite.T(), ddl, "device_id VARCHAR(255) UNIQUE")
	assert.Contains(suite.T(), ddl, "CONSTRAINT u_users_email_device_id UNIQUE(email, device_id)")
	assert.Contains(suite.T(), ddl, ");")
}

func (suite *TableTestSuite) TestTableIndex() {
	usersTable := Table(
		"users",
		Column("id", Varchar().Size(40)),
		Column("email", Varchar().Size(40)).Unique(),
		Index("users", "id"),
		Index("users", "email"),
		Index("users", "id", "email"),
	)
	ddl := usersTable.Create(NewDialect("default"))
	assert.Contains(suite.T(), ddl, "CREATE TABLE users (")
	assert.Contains(suite.T(), ddl, "id VARCHAR(40)")
	assert.Contains(suite.T(), ddl, "email VARCHAR(40) UNIQUE")
	assert.Contains(suite.T(), ddl, ")")
	assert.Contains(suite.T(), ddl, "CREATE INDEX i_id ON users(id)")
	assert.Contains(suite.T(), ddl, "CREATE INDEX i_email ON users(email)")
	assert.Contains(suite.T(), ddl, "CREATE INDEX i_id_email ON users(id, email);")

	assert.Equal(suite.T(), ColumnElem{Name: "id", Type: Varchar().Size(40), Table: "users"}, usersTable.C("id"))
	assert.Zero(suite.T(), usersTable.C("nonExisting"))
}

func (suite *TableTestSuite) TestTableIndexChain() {
	usersTable := Table("users", Column("id", Varchar().Size(40))).Index("id")
	ddl := usersTable.Create(NewDialect("default"))
	assert.Equal(suite.T(), "CREATE TABLE users (\n\tid VARCHAR(40)\n);\nCREATE INDEX i_id ON users(id);", ddl)
}

func (suite *TableTestSuite) TestTableStarters() {
	users := Table(
		"users",
		Column("id", Varchar().Size(40)),
		Column("email", Varchar().Size(40)).Unique(),
		PrimaryKey("id"),
	)

	dialect := NewDialect("default")

	ins := users.
		Insert().
		Values(map[string]interface{}{
			"id":    "5a73ef89-cf0a-4c51-ab8c-cc273ebb3a55",
			"email": "al@pacino.com",
		}).
		Build(dialect)

	assert.Contains(suite.T(), ins.SQL(), "INSERT INTO users")
	assert.Contains(suite.T(), ins.SQL(), "id")
	assert.Contains(suite.T(), ins.SQL(), "email")
	assert.Contains(suite.T(), ins.SQL(), "VALUES(?, ?)")
	assert.Contains(suite.T(), ins.Bindings(), "5a73ef89-cf0a-4c51-ab8c-cc273ebb3a55")
	assert.Contains(suite.T(), ins.Bindings(), "al@pacino.com")

	ups := users.Upsert()
	assert.Equal(suite.T(), users, ups.Table)

	upd := users.
		Update().
		Values(map[string]interface{}{
			"email": "al@pacino.com",
		}).
		Where(users.C("id").Eq("5a73ef89-cf0a-4c51-ab8c-cc273ebb3a55")).
		Build(dialect)

	assert.Equal(suite.T(), "UPDATE users\nSET email = ?\nWHERE id = ?;", upd.SQL())
	assert.Equal(suite.T(), []interface{}{"al@pacino.com", "5a73ef89-cf0a-4c51-ab8c-cc273ebb3a55"}, upd.Bindings())

	del := users.
		Delete().
		Where(users.C("id").Eq("5a73ef89-cf0a-4c51-ab8c-cc273ebb3a55")).
		Build(dialect)

	assert.Equal(suite.T(), "DELETE FROM users\nWHERE users.id = ?;", del.SQL())
	assert.Equal(suite.T(), []interface{}{"5a73ef89-cf0a-4c51-ab8c-cc273ebb3a55"}, del.Bindings())

	sel := users.
		Select(users.C("id"), users.C("email")).
		Where(users.C("id").Eq("5a73ef89-cf0a-4c51-ab8c-cc273ebb3a55")).
		Build(dialect)

	assert.Equal(suite.T(), "SELECT id, email\nFROM users\nWHERE id = ?;", sel.SQL())
	assert.Equal(suite.T(), []interface{}{"5a73ef89-cf0a-4c51-ab8c-cc273ebb3a55"}, sel.Bindings())
}

func TestTableTestSuite(t *testing.T) {
	suite.Run(t, new(TableTestSuite))
}
