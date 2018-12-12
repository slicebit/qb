package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aacanakin/qb"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var postgresDsn = "user=postgres dbname=qb_test sslmode=disable"

type PostgresTestSuite struct {
	suite.Suite
	engine   *qb.Engine
	metadata *qb.MetaDataElem
	ctx      *qb.CompilerContext
}

func (suite *PostgresTestSuite) SetupTest() {
	var err error

	suite.engine, err = qb.New("postgres", postgresDsn)
	suite.ctx = qb.NewCompilerContext(suite.engine.Dialect())

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), suite.engine)

	suite.metadata = qb.MetaData()
}

func (suite *PostgresTestSuite) TestPostgresBlob() {
	dialect := NewDialect()
	assert.Equal(suite.T(), "bytea", dialect.CompileType(qb.Blob()))
}

func (suite *PostgresTestSuite) TestUUID() {
	dialect := NewDialect()
	assert.Equal(suite.T(), "UUID", dialect.CompileType(qb.UUID()))
}

func (suite *PostgresTestSuite) TestDialectSimple() {
	dialect := NewDialect()
	assert.Equal(suite.T(), false, dialect.SupportsUnsigned())
	assert.Equal(suite.T(), "test", dialect.Escape("test"))
	assert.Equal(suite.T(), false, dialect.Escaping())
	assert.Equal(suite.T(), "postgres", dialect.Driver())
}

func (suite *PostgresTestSuite) TestDialectEscaping() {
	dialect := NewDialect()
	dialect.SetEscaping(true)
	assert.Equal(suite.T(), true, dialect.Escaping())
	assert.Equal(suite.T(), "\"test\"", dialect.Escape("test"))
	assert.Equal(suite.T(), []string{"\"test\""}, dialect.EscapeAll([]string{"test"}))
}

func (suite *PostgresTestSuite) TestDialectIntAutoIncrement() {
	dialect := NewDialect()
	col := qb.Column("autoinc", qb.Int()).AutoIncrement()
	assert.Equal(suite.T(), "SERIAL", dialect.AutoIncrement(&col))
}

func (suite *PostgresTestSuite) TestDialectBigIntAutoIncrement() {
	dialect := NewDialect()
	col := qb.Column("autoinc", qb.BigInt()).AutoIncrement()
	assert.Equal(suite.T(), "BIGSERIAL", dialect.AutoIncrement(&col))
}

func (suite *PostgresTestSuite) TestDialectSmallIntAutoIncrement() {
	dialect := NewDialect()
	col := qb.Column("autoinc", qb.SmallInt()).AutoIncrement()
	assert.Equal(suite.T(), "SMALLSERIAL", dialect.AutoIncrement(&col))
}

func (suite *PostgresTestSuite) TestWrapError() {
	err := errors.New("xxx")
	dialect := NewDialect()
	qbErr := dialect.WrapError(err)
	assert.Equal(suite.T(), err, qbErr.Orig)

	for _, tt := range []struct {
		pgCode string
		qbCode qb.ErrorCode
	}{
		{"0A000", qb.ErrNotSupported},
		{"20000", qb.ErrProgramming},
		{"21000", qb.ErrProgramming},
		{"22000", qb.ErrData},
		{"23000", qb.ErrIntegrity},
		{"24000", qb.ErrInternal},
		{"27000", qb.ErrOperational},
		{"2D000", qb.ErrInternal},
		{"34000", qb.ErrOperational},
		{"39000", qb.ErrInternal},
		{"3D000", qb.ErrProgramming},
		{"40000", qb.ErrOperational},
		{"42000", qb.ErrProgramming},
		{"54000", qb.ErrOperational},
		{"F0000", qb.ErrInternal},
		{"HV000", qb.ErrOperational},
		{"P0000", qb.ErrInternal},
		{"ZZ000", qb.ErrDatabase},
	} {
		pgErr := pq.Error{Code: pq.ErrorCode(tt.pgCode)}
		qbErr := suite.engine.Dialect().WrapError(&pgErr)
		assert.Equal(suite.T(), tt.qbCode, qbErr.Code)
	}
}

func (suite *PostgresTestSuite) TestPostgres() {
	type Actor struct {
		ID          string         `db:"id"`
		Email       string         `db:"email"`
		FullName    string         `db:"full_name"`
		Bio         sql.NullString `db:"bio"`
		Oscars      int            `db:"oscars"`
		IgnoreField string         `db:"-"`
	}

	type Session struct {
		ID        int64      `db:"id"`
		ActorID   string     `db:"actor_id"`
		AuthToken string     `db:"auth_token"`
		CreatedAt *time.Time `db:"created_at"`
		ExpiresAt *time.Time `db:"expires_at"`
	}

	actorsTable := qb.Table(
		"actors",
		qb.Column("id", qb.Type("UUID")),
		qb.Column("email", qb.Varchar()).Unique().NotNull(),
		qb.Column("full_name", qb.Varchar()).NotNull(),
		qb.Column("bio", qb.Text()).Null(),
		qb.Column("oscars", qb.Int()).Default(0),
		qb.PrimaryKey("id"),
	)

	sessionsTable := qb.Table(
		"sessions",
		qb.Column("id", qb.Type("BIGSERIAL")),
		qb.Column("actor_id", qb.Type("UUID")),
		qb.Column("auth_token", qb.Type("UUID")),
		qb.Column("created_at", qb.Timestamp()).NotNull(),
		qb.Column("expires_at", qb.Timestamp()).NotNull(),
		qb.PrimaryKey("id"),
		qb.ForeignKey("actor_id").References("actors", "id"),
	).Index("created_at", "expires_at")

	var err error

	suite.metadata.AddTable(actorsTable)
	suite.metadata.AddTable(sessionsTable)

	err = suite.metadata.CreateAll(suite.engine)
	fmt.Println("Metadata create all", err)
	assert.Nil(suite.T(), err)

	ins := qb.Insert(actorsTable).Values(map[string]interface{}{
		"id":        "b6f8bfe3-a830-441a-a097-1777e6bfae95",
		"email":     "jack@nicholson.com",
		"full_name": "Jack Nicholson",
		"bio":       sql.NullString{String: "Jack Nicholson, an American actor, producer, screen-writer and director, is a three-time Academy Award winner and twelve-time nominee.", Valid: true},
	})

	_, err = suite.engine.Exec(ins)

	ins = qb.Insert(sessionsTable).Values(map[string]interface{}{
		"actor_id":   "b6f8bfe3-a830-441a-a097-1777e6bfae95",
		"auth_token": "e4968197-6137-47a4-ba79-690d8c552248",
		"created_at": time.Now(),
		"expires_at": time.Now().Add(24 * time.Hour),
	}).Returning(sessionsTable.C("id"))

	var id int64
	err = suite.engine.QueryRow(ins).Scan(&id)
	assert.Nil(suite.T(), err)

	statement := qb.Insert(actorsTable).Values(map[string]interface{}{
		"id":        "b6f8bfe3-a830-441a-a097-1777e6bfae95",
		"email":     "jack@nicholson.com",
		"full_name": "Jack Nicholson",
		"bio":       sql.NullString{},
	})

	_, err = suite.engine.Exec(statement)
	assert.NotNil(suite.T(), err)

	statement = qb.Insert(actorsTable).Values(map[string]interface{}{
		"id":        "cf28d117-a12d-4b75-acd8-73a7d3cbb15f",
		"email":     "jack@nicholson2.com",
		"full_name": "Jack Nicholson",
		"bio":       sql.NullString{},
	})

	_, err = suite.engine.Exec(statement)
	assert.Nil(suite.T(), err)

	// find user using QueryRow()
	sel := qb.Select(
		actorsTable.C("id"),
		actorsTable.C("email"),
		actorsTable.C("full_name"),
		actorsTable.C("bio")).
		From(actorsTable).
		Where(actorsTable.C("id").Eq("cf28d117-a12d-4b75-acd8-73a7d3cbb15f"))

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
	var actor Actor

	sel = qb.Select(
		actorsTable.C("id"),
		actorsTable.C("email"),
		actorsTable.C("full_name"),
		actorsTable.C("bio")).
		From(actorsTable).
		Where(actorsTable.C("id").Eq("b6f8bfe3-a830-441a-a097-1777e6bfae95"))

	err = suite.engine.Get(sel, &actor)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), "jack@nicholson.com", actor.Email)
	assert.Equal(suite.T(), "Jack Nicholson", actor.FullName)
	assert.Equal(suite.T(), "Jack Nicholson, an American actor, producer, screen-writer and director, is a three-time Academy Award winner and twelve-time nominee.", actor.Bio.String)

	// select using join
	sessionSlice := []Session{}

	sel = qb.Select(
		sessionsTable.C("id"),
		sessionsTable.C("actor_id"),
		sessionsTable.C("auth_token"),
		sessionsTable.C("created_at"),
		sessionsTable.C("expires_at")).
		From(sessionsTable).
		InnerJoin(actorsTable, sessionsTable.C("actor_id"), actorsTable.C("id")).
		Where(sessionsTable.C("actor_id").Eq("b6f8bfe3-a830-441a-a097-1777e6bfae95"))

	err = suite.engine.Select(sel, &sessionSlice)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, len(sessionSlice))

	assert.Equal(suite.T(), int64(1), sessionSlice[0].ID)
	assert.Equal(suite.T(), "b6f8bfe3-a830-441a-a097-1777e6bfae95", sessionSlice[0].ActorID)
	assert.Equal(suite.T(), "e4968197-6137-47a4-ba79-690d8c552248", sessionSlice[0].AuthToken)

	// update user

	upd := qb.Update(actorsTable).Values(map[string]interface{}{
		"bio": sql.NullString{Valid: false},
	})

	_, err = suite.engine.Exec(upd)

	assert.Nil(suite.T(), err)

	sel = qb.Select(
		actorsTable.C("id"),
		actorsTable.C("email"),
		actorsTable.C("full_name"),
		actorsTable.C("bio")).
		From(actorsTable).
		Where(actorsTable.C("id").Eq("b6f8bfe3-a830-441a-a097-1777e6bfae95"))

	err = suite.engine.Get(sel, &actor)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), actor.Bio, sql.NullString{Valid: false})

	// delete session
	del := qb.Delete(sessionsTable).Where(
		sessionsTable.C("auth_token").Eq("99e591f8-1025-41ef-a833-6904a0f89a38"),
	)
	_, err = suite.engine.Exec(del)
	assert.Nil(suite.T(), err)

	// drop tables
	assert.Nil(suite.T(), suite.metadata.DropAll(suite.engine))
}

func (suite *PostgresTestSuite) TestAutoIncrement() {
	dialect := NewDialect()
	col := qb.Column("id", qb.BigInt()).AutoIncrement()
	assert.Equal(suite.T(),
		"BIGSERIAL",
		dialect.AutoIncrement(&col))

	col = qb.Column("id", qb.SmallInt()).AutoIncrement()
	assert.Equal(suite.T(),
		"SMALLSERIAL",
		dialect.AutoIncrement(&col))

	col = qb.Column("id", qb.Int()).AutoIncrement()
	assert.Equal(suite.T(),
		"SERIAL",
		dialect.AutoIncrement(&col))

	col = qb.Column("id", qb.Int()).AutoIncrement()
	col.Options.InlinePrimaryKey = true
	assert.Equal(suite.T(),
		"SERIAL PRIMARY KEY",
		dialect.AutoIncrement(&col))
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

	sql := ups.Accept(suite.ctx)
	binds := suite.ctx.Binds

	assert.Contains(suite.T(), sql, "INSERT INTO users")
	assert.Contains(suite.T(), sql, "id", "email")
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

	ctx := qb.NewCompilerContext(NewDialect())
	sql = ups.Accept(ctx)
	binds = ctx.Binds

	assert.Contains(suite.T(), sql, "INSERT INTO users")
	assert.Contains(suite.T(), sql, "id", "email")
	assert.Contains(suite.T(), sql, "ON CONFLICT")
	assert.Contains(suite.T(), sql, "DO UPDATE SET")
	assert.Contains(suite.T(), sql, "VALUES($1, $2)")
	assert.Contains(suite.T(), sql, "RETURNING id, email")
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
