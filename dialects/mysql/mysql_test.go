package mysql

import (
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/slicebit/qb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var mysqlDsn = "root:@tcp(localhost:3306)/qb_test?charset=utf8"

type MysqlTestSuite struct {
	suite.Suite
	engine   *qb.Engine
	metadata *qb.MetaDataElem
}

func (suite *MysqlTestSuite) SetupTest() {
	var err error
	suite.engine, err = qb.New("mysql", mysqlDsn)

	assert.Nil(suite.T(), err)
	err = suite.engine.Ping()

	assert.Nil(suite.T(), err)
	suite.metadata = qb.MetaData()

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), suite.engine)

	suite.engine.DB().Exec("DROP TABLE IF EXISTS user")
	suite.engine.DB().Exec("DROP TABLE IF EXISTS session")
}

func (suite *MysqlTestSuite) TestUUID() {
	assert.Equal(suite.T(), "VARCHAR(36)", suite.engine.Dialect().CompileType(qb.UUID()))
}

func (suite *MysqlTestSuite) TestDialect() {
	dialect := qb.NewDialect("mysql")
	assert.Equal(suite.T(), true, dialect.SupportsUnsigned())
	assert.Equal(suite.T(), "test", dialect.Escape("test"))
	assert.Equal(suite.T(), false, dialect.Escaping())
	dialect.SetEscaping(true)
	assert.Equal(suite.T(), true, dialect.Escaping())
	assert.Equal(suite.T(), "`test`", dialect.Escape("test"))
	assert.Equal(suite.T(), []string{"`test`"}, dialect.EscapeAll([]string{"test"}))
	assert.Equal(suite.T(), "mysql", dialect.Driver())
}

func (suite *MysqlTestSuite) TestWrapError() {
	dialect := qb.NewDialect("mysql")
	err := errors.New("xxx")
	qbErr := dialect.WrapError(err)
	assert.Equal(suite.T(), err, qbErr.Orig)

	for _, tt := range []struct {
		mErr   uint16
		qbCode qb.ErrorCode
	}{
		{ER_SYNTAX_ERROR, qb.ErrProgramming},
		{ER_DATA_TOO_LONG, qb.ErrData},
		{ER_CANNOT_ADD_FOREIGN, qb.ErrIntegrity},
		{ER_FEATURE_DISABLED, qb.ErrNotSupported},
		{ER_CHECKREAD, qb.ErrOperational},
		{999, qb.ErrInternal},
	} {
		mErr := mysql.MySQLError{Number: tt.mErr}
		qbErr := dialect.WrapError(&mErr)
		assert.Equal(suite.T(), tt.qbCode, qbErr.Code)
	}
}

func (suite *MysqlTestSuite) TestMysql() {
	type User struct {
		ID       string         `db:"id"`
		Email    string         `db:"email"`
		FullName string         `db:"full_name"`
		Bio      sql.NullString `db:"bio"`
		Oscars   int            `db:"oscars"`
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
		qb.Column("id", qb.Varchar().Size(40)),
		qb.Column("email", qb.Varchar()).Unique().NotNull(),
		qb.Column("full_name", qb.Varchar()).NotNull(),
		qb.Column("bio", qb.Text()).Null(),
		qb.Column("oscars", qb.Int()).Default(0),
		qb.PrimaryKey("id"),
	)

	sessions := qb.Table(
		"session",
		qb.Column("id", qb.BigInt()).AutoIncrement(),
		qb.Column("user_id", qb.Varchar().Size(40)).NotNull(),
		qb.Column("auth_token", qb.Varchar().Size(40)).NotNull().Unique(),
		qb.Column("created_at", qb.Timestamp()).Null(),
		qb.Column("expires_at", qb.Timestamp()).Null(),
		qb.PrimaryKey("id"),
		qb.ForeignKey("user_id").References("user", "id"),
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

	res, err := suite.engine.Exec(ins)
	assert.Nil(suite.T(), err)

	lastInsertID, err := res.LastInsertId()
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), lastInsertID, int64(1))

	rowsAffected, err := res.RowsAffected()
	assert.Equal(suite.T(), rowsAffected, int64(1))

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
	sel = qb.Select(sessions.C("id"), sessions.C("user_id"), sessions.C("auth_token")).
		From(sessions).
		InnerJoin(users, sessions.C("user_id"), users.C("id")).
		Where(users.C("id").Eq("b6f8bfe3-a830-441a-a097-1777e6bfae95"))

	err = suite.engine.Select(sel, &sessionSlice)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), len(sessionSlice), 1)

	assert.Equal(suite.T(), sessionSlice[0].ID, int64(1))
	assert.Equal(suite.T(), sessionSlice[0].UserID, "b6f8bfe3-a830-441a-a097-1777e6bfae95")
	assert.Equal(suite.T(), sessionSlice[0].AuthToken, "e4968197-6137-47a4-ba79-690d8c552248")

	upd := qb.Update(users).
		Values(map[string]interface{}{
			"bio": sql.NullString{String: "nil", Valid: false},
		}).Where(users.C("id").Eq("b6f8bfe3-a830-441a-a097-1777e6bfae95"))

	_, err = suite.engine.Exec(upd)
	assert.Nil(suite.T(), err)

	sel = qb.Select(users.C("id"), users.C("email"), users.C("full_name"), users.C("bio")).
		From(users).
		Where(users.C("id").Eq("b6f8bfe3-a830-441a-a097-1777e6bfae95"))

	err = suite.engine.Get(sel, &user)
	assert.Equal(suite.T(), user.Bio, sql.NullString{String: "", Valid: false})

	del := qb.Delete(sessions).Where(sessions.C("auth_token").Eq("99e591f8-1025-41ef-a833-6904a0f89a38"))
	_, err = suite.engine.Exec(del)
	assert.Nil(suite.T(), err)

	// drop tables
	assert.Nil(suite.T(), suite.metadata.DropAll(suite.engine))
}

func (suite *MysqlTestSuite) TestUpsert() {
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

	ctx := qb.NewCompilerContext(NewDialect())
	sql := ups.Accept(ctx)
	binds := ctx.Binds()

	assert.Contains(suite.T(), sql, "INSERT INTO users")
	assert.Contains(suite.T(), sql, "id", "email", "created_at")
	assert.Contains(suite.T(), sql, "VALUES(?, ?, ?)")
	assert.Contains(suite.T(), sql, "ON DUPLICATE KEY UPDATE")
	assert.Contains(suite.T(), sql, "id = ?", "email = ?", "created_at = ?")
	assert.Contains(suite.T(), binds, "9883cf81-3b56-4151-ae4e-3903c5bc436d")
	assert.Contains(suite.T(), binds, "al@pacino.com")
	assert.Equal(suite.T(), 6, len(binds))
}

func TestMysqlTestSuite(t *testing.T) {
	suite.Run(t, new(MysqlTestSuite))
}

func init() {
	dsn := os.Getenv("QBTEST_MYSQL")
	if dsn != "" {
		mysqlDsn = dsn
	}
}
