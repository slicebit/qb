package qb

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
	"errors"
)

type SqliteTestSuite struct {
	suite.Suite
	engine   *Engine
	metadata *MetaDataElem
}

func (suite *SqliteTestSuite) SetupTest() {
	var err error

	suite.engine, err = New("sqlite3", "./qb_test.db")
	suite.engine.Logger().SetLogFlags(LQuery | LBindings)
	suite.engine.Dialect().SetEscaping(true)

	suite.metadata = MetaData()

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), suite.engine)
}

func (suite *SqliteTestSuite) TestUUID() {
	assert.Equal(suite.T(), "VARCHAR(36)", suite.engine.Dialect().CompileType(UUID()))
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

	users := Table(
		"users",
		Column("id", Varchar().Size(40)),
		Column("email", Varchar()).NotNull().Unique(),
		Column("full_name", Varchar()).NotNull(),
		Column("bio", Text()).Null(),
		Column("oscars", Int()).NotNull().Default(0),
		PrimaryKey("id"),
	)

	sessions := Table(
		"sessions",
		Column("user_id", Varchar().Size(40)),
		Column("auth_token", Varchar().Size(40)).NotNull().Unique(),
		Column("created_at", Timestamp()).NotNull(),
		Column("expires_at", Timestamp()).NotNull(),
		ForeignKey("user_id").References("users", "id"),
	)

	var err error

	suite.metadata.AddTable(users)
	suite.metadata.AddTable(sessions)

	err = suite.metadata.CreateAll(suite.engine)
	assert.Nil(suite.T(), err)

	ins := Insert(users).Values(map[string]interface{}{
		"id":        "b6f8bfe3-a830-441a-a097-1777e6bfae95",
		"email":     "jack@nicholson.com",
		"full_name": "Jack Nicholson",
		"bio":       sql.NullString{String: "Jack Nicholson, an American actor, producer, screen-writer and director, is a three-time Academy Award winner and twelve-time nominee.", Valid: true},
	})

	_, err = suite.engine.Exec(ins)
	assert.Nil(suite.T(), err)

	ins = Insert(sessions).Values(map[string]interface{}{
		"user_id":    "b6f8bfe3-a830-441a-a097-1777e6bfae95",
		"auth_token": "e4968197-6137-47a4-ba79-690d8c552248",
		"created_at": time.Now(),
		"expires_at": time.Now().Add(24 * time.Hour),
	})

	_, err = suite.engine.Exec(ins)
	assert.Nil(suite.T(), err)

	// find user
	var user User

	sel := Select(users.C("id"), users.C("email"), users.C("full_name"), users.C("bio")).
		From(users).
		Where(users.C("id").Eq("b6f8bfe3-a830-441a-a097-1777e6bfae95"))

	err = suite.engine.Get(sel, &user)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), "jack@nicholson.com", user.Email)
	assert.Equal(suite.T(), "Jack Nicholson", user.FullName)
	assert.Equal(suite.T(), "Jack Nicholson, an American actor, producer, screen-writer and director, is a three-time Academy Award winner and twelve-time nominee.", user.Bio.String)

	// select using join
	sessionSlice := []Session{}

	sel = Select(sessions.C("user_id"), sessions.C("auth_token"), sessions.C("created_at"), sessions.C("expires_at")).
		From(sessions).
		InnerJoin(users, sessions.C("user_id"), users.C("id")).
		Where(users.C("id").Eq("b6f8bfe3-a830-441a-a097-1777e6bfae95"))

	err = suite.engine.Select(sel, &sessionSlice)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, len(sessionSlice))

	assert.Equal(suite.T(), "b6f8bfe3-a830-441a-a097-1777e6bfae95", sessionSlice[0].UserID)
	assert.Equal(suite.T(), "e4968197-6137-47a4-ba79-690d8c552248", sessionSlice[0].AuthToken)

	upd := Update(users).
		Values(map[string]interface{}{
			"bio": sql.NullString{String: "nil", Valid: false},
		}).
		Where(users.C("id").Eq("b6f8bfe3-a830-441a-a097-1777e6bfae95"))

	_, err = suite.engine.Exec(upd)
	assert.Nil(suite.T(), err)

	sel = Select(users.C("id"), users.C("email"), users.C("full_name"), users.C("bio")).
		From(users).
		Where(users.C("id").Eq("b6f8bfe3-a830-441a-a097-1777e6bfae95"))

	err = suite.engine.Get(sel, &user)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), user.Bio, sql.NullString{String: "", Valid: false})
	assert.Equal(suite.T(), sql.NullString{String: "", Valid: false}, user.Bio)

	del := Delete(sessions).Where(sessions.C("auth_token").Eq("99e591f8-1025-41ef-a833-6904a0f89a38"))

	// delete session
	_, err = suite.engine.Exec(del)
	assert.Nil(suite.T(), err)

	// drop tables
	assert.Nil(suite.T(), suite.metadata.DropAll(suite.engine))
}

func TestSQLiteExtractError(t *testing.T) {
	engine, err := New("sqlite3", "./qb_test.db")
	if err != nil {
		t.Fatal(err)
	}
	myErr := errors.New("some error")
	assert.Equal(t, NewQbError(myErr, nil), engine.dialect.ExtractError(myErr, nil))
}

func (suite *SqliteTestSuite) TestSqliteAutoIncrement() {
	col := Column("test", Int()).AutoIncrement()
	assert.Panics(suite.T(), func() {
		col.String(suite.engine.Dialect())
	})
}

func TestSqliteTestSuite(t *testing.T) {
	suite.Run(t, new(SqliteTestSuite))
}
