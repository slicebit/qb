package qb

import (
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var postgresDsn = "user=postgres dbname=qb_test sslmode=disable"

type PostgresTestSuite struct {
	suite.Suite
	engine   *Engine
	metadata *MetaDataElem
}

func TestPostgresBlob(t *testing.T) {
	assert.Equal(t, "bytea", NewDialect("postgres").CompileType(Blob()))
}

func (suite *PostgresTestSuite) SetupTest() {
	var err error

	suite.engine, err = New("postgres", postgresDsn)
	suite.engine.Dialect().SetEscaping(true)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), suite.engine)

	suite.metadata = MetaData()
}

func (suite *PostgresTestSuite) TestUUID() {
	assert.Equal(suite.T(), "UUID", suite.engine.Dialect().CompileType(UUID()))
}

func (suite *PostgresTestSuite) TestPostgres() {
	type User struct {
		ID          string         `db:"id"`
		Email       string         `db:"email"`
		FullName    string         `db:"full_name"`
		Bio         sql.NullString `db:"bio"`
		Oscars      int            `db:"oscars"`
		IgnoreField string         `db:"-"`
	}

	type Session struct {
		ID        int64      `db:"id"`
		UserID    string     `db:"user_id"`
		AuthToken string     `db:"auth_token"`
		CreatedAt *time.Time `db:"created_at"`
		ExpiresAt *time.Time `db:"expires_at"`
	}

	users := Table(
		"user",
		Column("id", Type("UUID")),
		Column("email", Varchar()).Unique().NotNull(),
		Column("full_name", Varchar()).NotNull(),
		Column("bio", Text()).Null(),
		Column("oscars", Int()).Default(0),
		PrimaryKey("id"),
	)

	sessions := Table(
		"session",
		Column("id", Type("BIGSERIAL")),
		Column("user_id", Type("UUID")),
		Column("auth_token", Type("UUID")),
		Column("created_at", Timestamp()).NotNull(),
		Column("expires_at", Timestamp()).NotNull(),
		PrimaryKey("id"),
		ForeignKey("user_id").References("user", "id"),
	).Index("created_at", "expires_at")

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

	// duplicate constraint should return a proper ConstraintViolation
	_, err = suite.engine.Exec(ins)
	assert.Error(suite.T(), err)
	integErr, ok := err.(IntegrityError)
	if !ok {
		suite.T().Fatal()
	}
	assert.Equal(suite.T(), "user_pkey", integErr.Constraint)
	assert.Equal(suite.T(), "constraint error: user_pkey", integErr.Error())

	ins = Insert(sessions).Values(map[string]interface{}{
		"user_id":    "b6f8bfe3-a830-441a-a097-1777e6bfae95",
		"auth_token": "e4968197-6137-47a4-ba79-690d8c552248",
		"created_at": time.Now(),
		"expires_at": time.Now().Add(24 * time.Hour),
	}).Returning(sessions.C("id"))

	var id int64
	err = suite.engine.QueryRow(ins).Scan(&id)
	assert.Nil(suite.T(), err)

	statement := Insert(users).Values(map[string]interface{}{
		"id":        "b6f8bfe3-a830-441a-a097-1777e6bfae95",
		"email":     "jack@nicholson.com",
		"full_name": "Jack Nicholson",
		"bio":       sql.NullString{},
	})

	_, err = suite.engine.Exec(statement)
	assert.NotNil(suite.T(), err)

	statement = Insert(users).Values(map[string]interface{}{
		"id":        "cf28d117-a12d-4b75-acd8-73a7d3cbb15f",
		"email":     "jack@nicholson2.com",
		"full_name": "Jack Nicholson",
		"bio":       sql.NullString{},
	})

	_, err = suite.engine.Exec(statement)
	assert.Nil(suite.T(), err)

	// find user using QueryRow()
	sel := Select(users.C("id"), users.C("email"), users.C("full_name"), users.C("bio")).From(users).
		Where(users.C("id").Eq("cf28d117-a12d-4b75-acd8-73a7d3cbb15f"))

	row := suite.engine.QueryRow(sel)
	assert.NotNil(suite.T(), row)

	// find user using Query()
	rows, err := suite.engine.Query(sel)
	assert.Nil(suite.T(), err)
	if err != nil {
		suite.T().Fatal(err)
	}
	rowLength := 0
	for rows.Next() {
		rowLength++
	}
	assert.Equal(suite.T(), 1, rowLength)

	// find user using session api's Find()
	var user User

	sel = Select(users.C("id"), users.C("email"), users.C("full_name"), users.C("bio")).From(users).
		Where(users.C("id").Eq("b6f8bfe3-a830-441a-a097-1777e6bfae95"))

	err = suite.engine.Get(sel, &user)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), "jack@nicholson.com", user.Email)
	assert.Equal(suite.T(), "Jack Nicholson", user.FullName)
	assert.Equal(suite.T(), "Jack Nicholson, an American actor, producer, screen-writer and director, is a three-time Academy Award winner and twelve-time nominee.", user.Bio.String)

	// select using join
	sessionSlice := []Session{}

	sel = Select(sessions.C("id"), sessions.C("user_id"), sessions.C("auth_token"), sessions.C("created_at"), sessions.C("expires_at")).
		From(sessions).
		InnerJoin(users, sessions.C("user_id"), users.C("id")).
		Where(sessions.C("user_id").Eq("b6f8bfe3-a830-441a-a097-1777e6bfae95"))

	err = suite.engine.Select(sel, &sessionSlice)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, len(sessionSlice))

	assert.Equal(suite.T(), int64(1), sessionSlice[0].ID)
	assert.Equal(suite.T(), "b6f8bfe3-a830-441a-a097-1777e6bfae95", sessionSlice[0].UserID)
	assert.Equal(suite.T(), "e4968197-6137-47a4-ba79-690d8c552248", sessionSlice[0].AuthToken)

	// update user

	upd := Update(users).Values(map[string]interface{}{
		"bio": sql.NullString{Valid: false},
	})

	_, err = suite.engine.Exec(upd)

	assert.Nil(suite.T(), err)

	sel = Select(users.C("id"), users.C("email"), users.C("full_name"), users.C("bio")).From(users).
		Where(users.C("id").Eq("b6f8bfe3-a830-441a-a097-1777e6bfae95"))

	err = suite.engine.Get(sel, &user)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), user.Bio, sql.NullString{Valid: false})

	// delete session
	del := Delete(sessions).Where(sessions.C("auth_token").Eq("99e591f8-1025-41ef-a833-6904a0f89a38"))
	_, err = suite.engine.Exec(del)
	assert.Nil(suite.T(), err)

	// drop tables
	assert.Nil(suite.T(), suite.metadata.DropAll(suite.engine))
}

func TestPostgresTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}

func init() {
	dsn := os.Getenv("QBTEST_POSTGRES")
	if dsn != "" {
		postgresDsn = dsn
	}
}

func TestPostGresDialectExtractError(t *testing.T) {
	pg, err := New("postgres", postgresDsn)
	if err != nil {
		t.Fatal(err)
	}
	myErr := errors.New("some error")
	assert.Equal(t, NewQbError(myErr, nil), pg.dialect.ExtractError(myErr, nil))
}
