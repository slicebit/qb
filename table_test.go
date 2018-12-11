package qb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TableTestSuite struct {
	suite.Suite
	dialect Dialect
}

func (suite *TableTestSuite) SetupTest() {
	suite.dialect = NewDefaultDialect()
}

func (suite *TableTestSuite) TestTableSimpleCreate() {
	usersTable := Table("users", Column("id", Varchar().Size(40)))
	assert.Equal(suite.T(), 1, len(usersTable.All()))

	ddl := usersTable.Create(suite.dialect)
	assert.Contains(suite.T(), ddl, "CREATE TABLE users (")
	assert.Contains(suite.T(), ddl, "id VARCHAR(40)")
	assert.Contains(suite.T(), ddl, ");")

	statement := usersTable.Build(suite.dialect)
	sql := statement.SQL()
	assert.Contains(suite.T(), sql, "CREATE TABLE users (")
	assert.Contains(suite.T(), sql, "id VARCHAR(40)")
	assert.Contains(suite.T(), sql, ");")
	assert.Equal(suite.T(), []interface{}{}, statement.Bindings())
}

func (suite *TableTestSuite) TestTableSimpleDrop() {
	usersTable := Table("users", Column("id", Varchar().Size(40)))

	assert.Equal(suite.T(), "DROP TABLE users;", usersTable.Drop(suite.dialect))
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

	ddl := usersTable.Create(suite.dialect)
	assert.Contains(suite.T(), ddl, "CREATE TABLE users (")
	assert.Contains(suite.T(), ddl, "auth_token VARCHAR(40)")
	assert.Contains(suite.T(), ddl, "role_id VARCHAR(40)")
	assert.Contains(suite.T(), ddl, "id VARCHAR(40) PRIMARY KEY")
	assert.Contains(suite.T(), ddl, "session_id VARCHAR(40)")
	assert.Contains(suite.T(), ddl, "FOREIGN KEY(session_id, auth_token) REFERENCES sessions(id, auth_token)")
	assert.Contains(suite.T(), ddl, "FOREIGN KEY(role_id) REFERENCES roles(id)")
	assert.Contains(suite.T(), ddl, ");")
}

func (suite *TableTestSuite) TestTableSimplePrimaryKey() {
	users := Table(
		"users",
		Column("id", Varchar().Size(40)).PrimaryKey(),
	)
	assert.Equal(suite.T(), []string{"id"}, users.PrimaryKeyConstraint.Columns)
}

func (suite *TableTestSuite) TestTableCompositePrimaryKey() {

	users := Table(
		"users",
		Column("fname", Varchar().Size(40)).PrimaryKey(),
		Column("lname", Varchar().Size(40)).PrimaryKey(),
	)

	assert.Equal(suite.T(), []string{"fname", "lname"}, users.PrimaryKeyConstraint.Columns)
	cols := users.PrimaryCols()
	assert.Equal(suite.T(), 2, len(cols))
	assert.Equal(suite.T(), "fname", cols[0].Name)
	assert.Equal(suite.T(), "lname", cols[1].Name)

	ddl := users.Create(suite.dialect)
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

	ddl := usersTable.Create(suite.dialect)
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
	ddl := usersTable.Create(suite.dialect)
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
	ddl := usersTable.Create(suite.dialect)
	assert.Contains(suite.T(), ddl, "CREATE TABLE users (")
	assert.Contains(suite.T(), ddl, "id VARCHAR(40)")
	assert.Contains(suite.T(), ddl, ");")
	assert.Contains(suite.T(), ddl, "CREATE INDEX i_id ON users(id);")
}

func (suite *TableTestSuite) TestTableStarters() {
	users := Table(
		"users",
		Column("id", Varchar().Size(40)),
		Column("email", Varchar().Size(40)).Unique(),
		PrimaryKey("id"),
	)

	ins := users.
		Insert().
		Values(map[string]interface{}{
			"id":    "5a73ef89-cf0a-4c51-ab8c-cc273ebb3a55",
			"email": "al@pacino.com",
		}).
		Build(suite.dialect)

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
		Build(suite.dialect)

	updSQL := upd.SQL()
	assert.Contains(suite.T(), updSQL, "UPDATE users")
	assert.Contains(suite.T(), updSQL, "SET email = ?")
	assert.Contains(suite.T(), updSQL, "WHERE id = ?;")

	assert.Equal(suite.T(), []interface{}{
		"al@pacino.com",
		"5a73ef89-cf0a-4c51-ab8c-cc273ebb3a55",
	}, upd.Bindings())

	del := users.
		Delete().
		Where(users.C("id").Eq("5a73ef89-cf0a-4c51-ab8c-cc273ebb3a55")).
		Build(suite.dialect)

	delSQL := del.SQL()

	assert.Contains(suite.T(), delSQL, "DELETE FROM users")
	assert.Contains(suite.T(), delSQL, "WHERE users.id = ?;")
	assert.Equal(suite.T(), []interface{}{"5a73ef89-cf0a-4c51-ab8c-cc273ebb3a55"}, del.Bindings())

	sel := users.
		Select(users.C("id"), users.C("email")).
		Where(users.C("id").Eq("5a73ef89-cf0a-4c51-ab8c-cc273ebb3a55")).
		Build(suite.dialect)

	selSQL := sel.SQL()

	assert.Contains(suite.T(), selSQL, "SELECT id, email")
	assert.Contains(suite.T(), selSQL, "FROM users")
	assert.Contains(suite.T(), selSQL, "WHERE id = ?;")
	assert.Equal(suite.T(), []interface{}{"5a73ef89-cf0a-4c51-ab8c-cc273ebb3a55"}, sel.Bindings())
}

func TestTableTestSuite(t *testing.T) {
	suite.Run(t, new(TableTestSuite))
}
