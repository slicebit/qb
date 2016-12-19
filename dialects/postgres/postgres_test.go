package qb

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"github.com/slicebit/qb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
)

var postgresDsn = "user=postgres dbname=qb_test sslmode=disable"

func asSQLBinds(clause qb.Clause, dialect qb.Dialect) (string, []interface{}) {
	ctx := qb.NewCompilerContext(dialect)
	return clause.Accept(ctx), ctx.Binds
}

type PostgresTestSuite struct {
	suite.Suite
	engine   *qb.Engine
	metadata *qb.MetaDataElem
}

func TestPostgresBlob(t *testing.T) {
	assert.Equal(t, "bytea", qb.NewDialect("postgres").CompileType(qb.Blob()))
}

func (suite *PostgresTestSuite) SetupTest() {
	var err error

	suite.engine, err = qb.New("postgres", postgresDsn)
	suite.engine.Dialect().SetEscaping(true)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), suite.engine)

	suite.metadata = qb.MetaData()
}

func (suite *PostgresTestSuite) TestUUID() {
	assert.Equal(suite.T(), "UUID", suite.engine.Dialect().CompileType(qb.UUID()))
}

func (suite *PostgresTestSuite) TestDialect() {
	dialect := qb.NewDialect("postgres")
	assert.Equal(suite.T(), false, dialect.SupportsUnsigned())
	assert.Equal(suite.T(), "test", dialect.Escape("test"))
	assert.Equal(suite.T(), false, dialect.Escaping())
	dialect.SetEscaping(true)
	assert.Equal(suite.T(), true, dialect.Escaping())
	assert.Equal(suite.T(), "\"test\"", dialect.Escape("test"))
	assert.Equal(suite.T(), []string{"\"test\""}, dialect.EscapeAll([]string{"test"}))
	assert.Equal(suite.T(), "postgres", dialect.Driver())

	col := qb.Column("autoinc", qb.Int()).AutoIncrement()
	assert.Equal(suite.T(), "SERIAL", dialect.AutoIncrement(&col))

	col = qb.Column("autoinc", qb.BigInt()).AutoIncrement()
	assert.Equal(suite.T(), "BIGSERIAL", dialect.AutoIncrement(&col))

	col = qb.Column("autoinc", qb.SmallInt()).AutoIncrement()
	assert.Equal(suite.T(), "SMALLSERIAL", dialect.AutoIncrement(&col))
}

func (suite *PostgresTestSuite) TestWrapError() {
	dialect := qb.NewDialect("postgres")
	err := errors.New("xxx")
	qbErr := dialect.WrapError(err)
	assert.Equal(suite.T(), err, qbErr.Orig)
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

	users := qb.Table(
		"user",
		qb.Column("id", qb.Type("UUID")),
		qb.Column("email", qb.Varchar()).Unique().NotNull(),
		qb.Column("full_name", qb.Varchar()).NotNull(),
		qb.Column("bio", qb.Text()).Null(),
		qb.Column("oscars", qb.Int()).Default(0),
		qb.PrimaryKey("id"),
	)

	sessions := qb.Table(
		"session",
		qb.Column("id", qb.Type("BIGSERIAL")),
		qb.Column("user_id", qb.Type("UUID")),
		qb.Column("auth_token", qb.Type("UUID")),
		qb.Column("created_at", qb.Timestamp()).NotNull(),
		qb.Column("expires_at", qb.Timestamp()).NotNull(),
		qb.PrimaryKey("id"),
		qb.ForeignKey("user_id").References("user", "id"),
	).Index("created_at", "expires_at")

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

	ins = qb.Insert(sessions).Values(map[string]interface{}{
		"user_id":    "b6f8bfe3-a830-441a-a097-1777e6bfae95",
		"auth_token": "e4968197-6137-47a4-ba79-690d8c552248",
		"created_at": time.Now(),
		"expires_at": time.Now().Add(24 * time.Hour),
	}).Returning(sessions.C("id"))

	var id int64
	err = suite.engine.QueryRow(ins).Scan(&id)
	assert.Nil(suite.T(), err)

	statement := qb.Insert(users).Values(map[string]interface{}{
		"id":        "b6f8bfe3-a830-441a-a097-1777e6bfae95",
		"email":     "jack@nicholson.com",
		"full_name": "Jack Nicholson",
		"bio":       sql.NullString{},
	})

	_, err = suite.engine.Exec(statement)
	assert.NotNil(suite.T(), err)

	statement = qb.Insert(users).Values(map[string]interface{}{
		"id":        "cf28d117-a12d-4b75-acd8-73a7d3cbb15f",
		"email":     "jack@nicholson2.com",
		"full_name": "Jack Nicholson",
		"bio":       sql.NullString{},
	})

	_, err = suite.engine.Exec(statement)
	assert.Nil(suite.T(), err)

	// find user using QueryRow()
	sel := qb.Select(users.C("id"), users.C("email"), users.C("full_name"), users.C("bio")).From(users).
		Where(users.C("id").Eq("cf28d117-a12d-4b75-acd8-73a7d3cbb15f"))

	row := suite.engine.QueryRow(sel)
	assert.NotNil(suite.T(), row)

	// find user using Query()
	rows, err := suite.engine.Query(sel)
	assert.Nil(suite.T(), err)
	rowLength := 0
	for rows.Next() {
		rowLength++
	}
	assert.Equal(suite.T(), 1, rowLength)

	// find user using session api's Find()
	var user User

	sel = qb.Select(users.C("id"), users.C("email"), users.C("full_name"), users.C("bio")).From(users).
		Where(users.C("id").Eq("b6f8bfe3-a830-441a-a097-1777e6bfae95"))

	err = suite.engine.Get(sel, &user)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), "jack@nicholson.com", user.Email)
	assert.Equal(suite.T(), "Jack Nicholson", user.FullName)
	assert.Equal(suite.T(), "Jack Nicholson, an American actor, producer, screen-writer and director, is a three-time Academy Award winner and twelve-time nominee.", user.Bio.String)

	// select using join
	sessionSlice := []Session{}

	sel = qb.Select(sessions.C("id"), sessions.C("user_id"), sessions.C("auth_token"), sessions.C("created_at"), sessions.C("expires_at")).
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

	upd := qb.Update(users).Values(map[string]interface{}{
		"bio": sql.NullString{Valid: false},
	})

	_, err = suite.engine.Exec(upd)

	assert.Nil(suite.T(), err)

	sel = qb.Select(users.C("id"), users.C("email"), users.C("full_name"), users.C("bio")).From(users).
		Where(users.C("id").Eq("b6f8bfe3-a830-441a-a097-1777e6bfae95"))

	err = suite.engine.Get(sel, &user)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), user.Bio, sql.NullString{Valid: false})

	// delete session
	del := qb.Delete(sessions).Where(sessions.C("auth_token").Eq("99e591f8-1025-41ef-a833-6904a0f89a38"))
	_, err = suite.engine.Exec(del)
	assert.Nil(suite.T(), err)

	// drop tables
	assert.Nil(suite.T(), suite.metadata.DropAll(suite.engine))
}

func (suite *PostgresTestSuite) TestAutoIncrement() {
	col := qb.Column("id", qb.BigInt()).AutoIncrement()
	assert.Equal(suite.T(),
		"BIGSERIAL",
		suite.engine.Dialect().AutoIncrement(&col))

	col = qb.Column("id", qb.SmallInt()).AutoIncrement()
	assert.Equal(suite.T(),
		"SMALLSERIAL",
		suite.engine.Dialect().AutoIncrement(&col))

	col = qb.Column("id", qb.Int()).AutoIncrement()
	assert.Equal(suite.T(),
		"SERIAL",
		suite.engine.Dialect().AutoIncrement(&col))

	col = qb.Column("id", qb.Int()).AutoIncrement()
	col.Options.InlinePrimaryKey = true
	assert.Equal(suite.T(),
		"SERIAL PRIMARY KEY",
		suite.engine.Dialect().AutoIncrement(&col))
}

func (suite *PostgresTestSuite) TestUpsert() {
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

	assert.Contains(suite.T(), sql, "INSERT INTO \"users\"")
	assert.Contains(suite.T(), sql, "\"id\"", "\"email\"")
	assert.Contains(suite.T(), sql, "VALUES($1, $2, $3)")
	assert.Contains(suite.T(), sql, "ON CONFLICT", "DO UPDATE SET")
	assert.Contains(suite.T(), binds, "9883cf81-3b56-4151-ae4e-3903c5bc436d")
	assert.Contains(suite.T(), binds, "al@pacino.com")
	assert.Equal(suite.T(), 6, len(binds))

	ups = qb.Upsert(users).
		Values(map[string]interface{}{
			"id":    "9883cf81-3b56-4151-ae4e-3903c5bc436d",
			"email": "al@pacino.com",
		}).
		Returning(users.C("id"), users.C("email"))

	sql, binds = asSQLBinds(ups, suite.engine.Dialect())
	assert.Contains(suite.T(), sql, "INSERT INTO \"users\"")
	assert.Contains(suite.T(), sql, "\"id\"", "\"email\"")
	assert.Contains(suite.T(), sql, "ON CONFLICT")
	assert.Contains(suite.T(), sql, "DO UPDATE SET")
	assert.Contains(suite.T(), sql, "VALUES($1, $2)")
	assert.Contains(suite.T(), sql, "RETURNING \"id\", \"email\"")
	assert.Contains(suite.T(), binds, "9883cf81-3b56-4151-ae4e-3903c5bc436d")
	assert.Contains(suite.T(), binds, "al@pacino.com")
	assert.Equal(suite.T(), 4, len(binds))
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
