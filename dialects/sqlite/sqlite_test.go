package qb

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/slicebit/qb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func asSQLBinds(clause qb.Clause, dialect qb.Dialect) (string, []interface{}) {
	ctx := qb.NewCompilerContext(dialect)
	return clause.Accept(ctx), ctx.Binds
}

type SqliteTestSuite struct {
	suite.Suite
	engine   *qb.Engine
	metadata *qb.MetaDataElem
}

func (suite *SqliteTestSuite) SetupTest() {
	var err error

	suite.engine, err = qb.New("sqlite3", "./qb_test.db")
	suite.engine.Dialect().SetEscaping(true)

	suite.metadata = qb.MetaData()

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), suite.engine)
}

func (suite *SqliteTestSuite) TestUUID() {
	assert.Equal(suite.T(), "VARCHAR(36)", suite.engine.Dialect().CompileType(qb.UUID()))
}

func (suite *SqliteTestSuite) TestDialect() {
	dialect := qb.NewDialect("sqlite")
	assert.Equal(suite.T(), false, dialect.SupportsUnsigned())
	assert.Equal(suite.T(), "test", dialect.Escape("test"))
	assert.Equal(suite.T(), false, dialect.Escaping())
	dialect.SetEscaping(true)
	assert.Equal(suite.T(), true, dialect.Escaping())
	assert.Equal(suite.T(), "\"test\"", dialect.Escape("test"))
	assert.Equal(suite.T(), []string{"\"test\""}, dialect.EscapeAll([]string{"test"}))
	assert.Equal(suite.T(), "sqlite3", dialect.Driver())
}

func (suite *SqliteTestSuite) TestWrapError() {
	dialect := qb.NewDialect("sqlite")
	err := errors.New("xxx")
	qbErr := dialect.WrapError(err)
	assert.Equal(suite.T(), err, qbErr.Orig)
}

func (suite *SqliteTestSuite) TestSqlite() {
	type User struct {
		ID       string         `db:"id"`
		Email    string         `db:"email"`
		FullName string         `db:"full_name"`
		Bio      sql.NullString `db:"bio"`
		Oscars   int            `db:"oscars"`
	}

	type Session struct {
		UserID    string    `db:"user_id"`
		AuthToken string    `db:"auth_token"`
		CreatedAt time.Time `db:"created_at"`
		ExpiresAt time.Time `db:"expires_at"`
	}

	users := qb.Table(
		"users",
		qb.Column("id", qb.Varchar().Size(40)),
		qb.Column("email", qb.Varchar()).NotNull().Unique(),
		qb.Column("full_name", qb.Varchar()).NotNull(),
		qb.Column("bio", qb.Text()).Null(),
		qb.Column("oscars", qb.Int()).NotNull().Default(0),
		qb.PrimaryKey("id"),
	)

	sessions := qb.Table(
		"sessions",
		qb.Column("user_id", qb.Varchar().Size(40)),
		qb.Column("auth_token", qb.Varchar().Size(40)).NotNull().Unique(),
		qb.Column("created_at", qb.Timestamp()).NotNull(),
		qb.Column("expires_at", qb.Timestamp()).NotNull(),
		qb.ForeignKey("user_id").References("users", "id"),
	)

	var err error

	suite.metadata.AddTable(users)
	suite.metadata.AddTable(sessions)

	err = suite.metadata.CreateAll(suite.engine)
	assert.Nil(suite.T(), err)

	ins := qb.Insert(users).Values(map[string]interface{}{
		"id":        "b6f8bfe3-a830-441a-a097-1777e6bfae95",
		"email":     "jack@nicholson.com",
		"full_name": "Jack Nicholson",
		"bio":       sql.NullString{String: "Jack Nicholson, an American actor, producer, screen-writer and director, is a three-time Academy Award winner and twelve-time nominee.", Valid: true},
	})

	_, err = suite.engine.Exec(ins)
	assert.Nil(suite.T(), err)

	ins = qb.Insert(sessions).Values(map[string]interface{}{
		"user_id":    "b6f8bfe3-a830-441a-a097-1777e6bfae95",
		"auth_token": "e4968197-6137-47a4-ba79-690d8c552248",
		"created_at": time.Now(),
		"expires_at": time.Now().Add(24 * time.Hour),
	})

	_, err = suite.engine.Exec(ins)
	assert.Nil(suite.T(), err)

	// find user
	var user User

	sel := qb.Select(users.C("id"), users.C("email"), users.C("full_name"), users.C("bio")).
		From(users).
		Where(users.C("id").Eq("b6f8bfe3-a830-441a-a097-1777e6bfae95"))

	err = suite.engine.Get(sel, &user)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), "jack@nicholson.com", user.Email)
	assert.Equal(suite.T(), "Jack Nicholson", user.FullName)
	assert.Equal(suite.T(), "Jack Nicholson, an American actor, producer, screen-writer and director, is a three-time Academy Award winner and twelve-time nominee.", user.Bio.String)

	// select using join
	sessionSlice := []Session{}

	sel = qb.Select(sessions.C("user_id"), sessions.C("auth_token"), sessions.C("created_at"), sessions.C("expires_at")).
		From(sessions).
		InnerJoin(users, sessions.C("user_id"), users.C("id")).
		Where(users.C("id").Eq("b6f8bfe3-a830-441a-a097-1777e6bfae95"))

	err = suite.engine.Select(sel, &sessionSlice)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, len(sessionSlice))

	assert.Equal(suite.T(), "b6f8bfe3-a830-441a-a097-1777e6bfae95", sessionSlice[0].UserID)
	assert.Equal(suite.T(), "e4968197-6137-47a4-ba79-690d8c552248", sessionSlice[0].AuthToken)

	upd := qb.Update(users).
		Values(map[string]interface{}{
			"bio": sql.NullString{String: "nil", Valid: false},
		}).
		Where(users.C("id").Eq("b6f8bfe3-a830-441a-a097-1777e6bfae95"))

	_, err = suite.engine.Exec(upd)
	assert.Nil(suite.T(), err)

	sel = qb.Select(users.C("id"), users.C("email"), users.C("full_name"), users.C("bio")).
		From(users).
		Where(users.C("id").Eq("b6f8bfe3-a830-441a-a097-1777e6bfae95"))

	err = suite.engine.Get(sel, &user)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), user.Bio, sql.NullString{String: "", Valid: false})
	assert.Equal(suite.T(), sql.NullString{String: "", Valid: false}, user.Bio)

	del := qb.Delete(sessions).Where(sessions.C("auth_token").Eq("99e591f8-1025-41ef-a833-6904a0f89a38"))

	// delete session
	_, err = suite.engine.Exec(del)
	assert.Nil(suite.T(), err)

	// drop tables
	assert.Nil(suite.T(), suite.metadata.DropAll(suite.engine))
}

func (suite *SqliteTestSuite) TestUpsert() {
	users := qb.Table(
		"users",
		qb.Column("id", qb.Varchar().Size(36)),
		qb.Column("email", qb.Varchar()).Unique(),
		qb.Column("created_at", qb.Timestamp()).NotNull(),
		qb.PrimaryKey("id"),
	)

	now := time.Now().UTC().String()

	ups := qb.Upsert(users).Values(map[string]interface{}{
		"id":         "9883cf81-3b56-4151-ae4e-3903c5bc436d",
		"email":      "al@pacino.com",
		"created_at": now,
	})

	sql, binds := asSQLBinds(ups, suite.engine.Dialect())
	assert.Contains(suite.T(), sql, `REPLACE INTO "users"`)
	assert.Contains(suite.T(), sql, "id", "email", "created_at")
	assert.Contains(suite.T(), sql, "VALUES(?, ?, ?)")
	assert.Contains(suite.T(), binds, "9883cf81-3b56-4151-ae4e-3903c5bc436d")
	assert.Contains(suite.T(), binds, "al@pacino.com")
	assert.Contains(suite.T(), binds, now)
	assert.Equal(suite.T(), 3, len(binds))
}

func (suite *SqliteTestSuite) TestSqliteAutoIncrement() {
	col := qb.Column("test", qb.Int()).AutoIncrement()
	assert.Panics(suite.T(), func() {
		col.String(suite.engine.Dialect())
	})

	col.Options.InlinePrimaryKey = true
	assert.Equal(suite.T(), "INTEGER PRIMARY KEY", suite.engine.Dialect().AutoIncrement(&col))
}

func TestSqliteTestSuite(t *testing.T) {
	suite.Run(t, new(SqliteTestSuite))
}
