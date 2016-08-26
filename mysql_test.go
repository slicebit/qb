package qb

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
)

var mysqlDsn = "root:@tcp(localhost:3306)/qb_test?charset=utf8"

type MysqlTestSuite struct {
	suite.Suite
	db *Session
}

func (suite *MysqlTestSuite) SetupTest() {
	var err error
	suite.db, err = New("mysql", mysqlDsn)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), suite.db)
	suite.db.Engine().DB().Exec("DROP TABLE IF EXISTS user")
	suite.db.Engine().DB().Exec("DROP TABLE IF EXISTS session")
}

func (suite *MysqlTestSuite) TestUUID() {
	assert.Equal(suite.T(), "VARCHAR(36)", suite.db.Dialect().CompileType(UUID()))
}

func (suite *MysqlTestSuite) TestMysql() {
	type User struct {
		ID       string         `qb:"type:varchar(40); constraints:primary_key"`
		Email    string         `qb:"constraints:unique, notnull"`
		FullName string         `qb:"constraints:notnull"`
		Bio      sql.NullString `qb:"type:text; constraints:null"`
		Oscars   int            `qb:"constraints:default(0)"`
	}

	type Session struct {
		ID        int64     `qb:"type:bigint; constraints:primary_key, auto_increment"`
		UserID    string    `qb:"type:varchar(40); constraints:ref(user.id)"`
		AuthToken string    `qb:"type:varchar(40); constraints:notnull, unique"`
		CreatedAt time.Time `qb:"constraints:null"`
		ExpiresAt time.Time `qb:"constraints:null"`
	}

	var err error

	suite.db.AddTable(User{})
	suite.db.AddTable(Session{})

	err = suite.db.CreateAll()
	assert.Nil(suite.T(), err)

	// add sample user & session
	suite.db.AddAll(
		&User{
			ID:       "b6f8bfe3-a830-441a-a097-1777e6bfae95",
			Email:    "jack@nicholson.com",
			FullName: "Jack Nicholson",
			Bio:      sql.NullString{String: "Jack Nicholson, an American actor, producer, screen-writer and director, is a three-time Academy Award winner and twelve-time nominee.", Valid: true},
		}, &Session{
			UserID:    "b6f8bfe3-a830-441a-a097-1777e6bfae95",
			AuthToken: "e4968197-6137-47a4-ba79-690d8c552248",
			CreatedAt: time.Now(),
			ExpiresAt: time.Now().Add(24 * time.Hour),
		},
	)

	err = suite.db.Commit()
	assert.Nil(suite.T(), err)

	// find user
	var user User

	suite.db.Find(&User{ID: "b6f8bfe3-a830-441a-a097-1777e6bfae95"}).One(&user)

	assert.Equal(suite.T(), "jack@nicholson.com", user.Email)
	assert.Equal(suite.T(), "Jack Nicholson", user.FullName)
	assert.Equal(suite.T(), "Jack Nicholson, an American actor, producer, screen-writer and director, is a three-time Academy Award winner and twelve-time nominee.", user.Bio.String)

	// select using join
	sessions := []Session{}
	err = suite.db.Query(suite.db.T("session").C("user_id"), suite.db.T("session").C("id"), suite.db.T("session").C("auth_token")).
		InnerJoin(suite.db.T("user"), suite.db.T("session").C("user_id"), suite.db.T("user").C("id")).
		Filter(suite.db.T("user").C("id").Eq("b6f8bfe3-a830-441a-a097-1777e6bfae95")).
		All(&sessions)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, len(sessions))

	assert.Equal(suite.T(), int64(1), sessions[0].ID)
	assert.Equal(suite.T(), "b6f8bfe3-a830-441a-a097-1777e6bfae95", sessions[0].UserID)
	assert.Equal(suite.T(), "e4968197-6137-47a4-ba79-690d8c552248", sessions[0].AuthToken)

	user.Bio = sql.NullString{String: "nil", Valid: false}
	suite.db.Add(user)

	err = suite.db.Commit()
	assert.Nil(suite.T(), err)

	suite.db.Find(&User{ID: "b6f8bfe3-a830-441a-a097-1777e6bfae95"}).One(&user)
	assert.Equal(suite.T(), sql.NullString{String: "", Valid: false}, user.Bio)

	// delete session
	suite.db.Delete(&Session{AuthToken: "99e591f8-1025-41ef-a833-6904a0f89a38"})
	err = suite.db.Commit()
	assert.Nil(suite.T(), err)

	// drop tables
	assert.Nil(suite.T(), suite.db.DropAll())
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
