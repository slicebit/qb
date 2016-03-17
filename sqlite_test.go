package qb

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type sUser struct {
	ID       string `qb:"constraints:primary_key"`
	Email    string `qb:"constraints:unique, notnull"`
	FullName string `qb:"constraints:notnull"`
	Password string `qb:"constraints:notnull"`
	Bio      string `qb:"type:text; constraints:null"`
}

type sSession struct {
	ID        int64     `qb:"type:bigint; constraints:primary_key, auto_increment"`
	UserID    string    `qb:"type:uuid; constraints:ref(s_user.id)"`
	AuthToken string    `qb:"type:uuid; constraints:notnull, unique"`
	CreatedAt time.Time `qb:"constraints:notnull"`
	ExpiresAt time.Time `qb:"constraints:notnull"`
}

type SqliteTestSuite struct {
	suite.Suite
	session *Session
}

func (suite *SqliteTestSuite) SetupTest() {
	var err error

	engine, err := NewEngine("sqlite3", ":memory:")
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), engine)
	engine.DB().Exec("CREATE DATABASE qb_test;")

	suite.session, err = New("sqlite3", ":memory:")
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), suite.session)
}

func (suite *SqliteTestSuite) TestSqlite() {

	type User struct {
		ID       string         `qb:"type:varchar(40); constraints:primary_key"`
		Email    string         `qb:"constraints:unique, notnull"`
		FullName string         `qb:"constraints:notnull"`
		Bio      sql.NullString `qb:"type:text; constraints:null"`
		Oscars   int            `qb:"constraints:default(0)"`
	}

	type Session struct {
		ID        int       `qb:"type:int; constraints:primary_key, auto_increment"`
		UserID    string    `qb:"type:varchar(40); constraints:ref(user.id)"`
		AuthToken string    `qb:"type:varchar(40); constraints:notnull, unique"`
		CreatedAt time.Time `qb:"constraints:notnull"`
		ExpiresAt time.Time `qb:"constraints:notnull"`
	}

	var err error

	suite.session.Metadata().Add(User{})
	suite.session.Metadata().Add(Session{})

	err = suite.session.Metadata().CreateAll()
	assert.Nil(suite.T(), err)

	// add sample user
	suite.session.Add(&User{
		ID:       "b6f8bfe3-a830-441a-a097-1777e6bfae95",
		Email:    "jack@nicholson.com",
		FullName: "Jack Nicholson",
		Bio:      sql.NullString{String: "Jack Nicholson, an American actor, producer, screen-writer and director, is a three-time Academy Award winner and twelve-time nominee.", Valid: true},
	})

	// add sample session
	suite.session.Add(&Session{
		UserID:    "b6f8bfe3-a830-441a-a097-1777e6bfae95",
		AuthToken: "e4968197-6137-47a4-ba79-690d8c552248",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})

	err = suite.session.Commit()
	assert.Nil(suite.T(), err)

	// find user
	var user User

	suite.session.Find(&User{ID: "b6f8bfe3-a830-441a-a097-1777e6bfae95"}).First(&user)

	assert.Equal(suite.T(), user.Email, "jack@nicholson.com")
	assert.Equal(suite.T(), user.FullName, "Jack Nicholson")
	assert.Equal(suite.T(), user.Bio.String, "Jack Nicholson, an American actor, producer, screen-writer and director, is a three-time Academy Award winner and twelve-time nominee.")

	// select using join
	sessions := []Session{}
	err = suite.session.Select("s.id", "s.user_id", "s.auth_token").
		From("session s").
		InnerJoin("user u", "u.id = s.user_id").
		Where("s.user_id = ?", "b6f8bfe3-a830-441a-a097-1777e6bfae95").
		All(&sessions)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), len(sessions), 1)

	assert.Equal(suite.T(), sessions[0].ID, 1)
	assert.Equal(suite.T(), sessions[0].UserID, "b6f8bfe3-a830-441a-a097-1777e6bfae95")
	assert.Equal(suite.T(), sessions[0].AuthToken, "e4968197-6137-47a4-ba79-690d8c552248")

	// update user
	update := suite.session.
		Update("user").
		Set(map[string]interface{}{
			"bio": nil,
		}).
		Where("id = ?", "b6f8bfe3-a830-441a-a097-1777e6bfae95").
		Query()

	suite.session.AddQuery(update)
	err = suite.session.Commit()
	assert.Nil(suite.T(), err)

	suite.session.Find(&User{ID: "b6f8bfe3-a830-441a-a097-1777e6bfae95"}).First(&user)
	assert.Equal(suite.T(), user.Bio, sql.NullString{String: "", Valid: false})

	// delete session
	suite.session.Delete(&Session{AuthToken: "99e591f8-1025-41ef-a833-6904a0f89a38"})
	err = suite.session.Commit()
	assert.Nil(suite.T(), err)

	// drop tables
	assert.Nil(suite.T(), suite.session.Metadata().DropAll())

	// fail model
	type FailModel struct {
		ID int64 `qb:"type:notype"`
	}

	assert.Panics(suite.T(), func() {
		suite.session.Add(FailModel{})
	})
}

func TestSqliteTestSuite(t *testing.T) {
	suite.Run(t, new(SqliteTestSuite))
}
