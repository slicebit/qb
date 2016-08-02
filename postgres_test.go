package qb

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
)

var postgresDsn = "user=postgres dbname=qb_test sslmode=disable"

type PostgresTestSuite struct {
	suite.Suite
	db *Session
}

func (suite *PostgresTestSuite) SetupTest() {

	var err error

	suite.db, err = New("postgres", postgresDsn)
	suite.db.Dialect().SetEscaping(true)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), suite.db)
}

func (suite *PostgresTestSuite) TestPostgres() {
	type User struct {
		ID          string         `db:"_id" qb:"type:uuid; constraints:primary_key"`
		Email       string         `qb:"constraints:unique, notnull"`
		FullName    string         `qb:"constraints:notnull"`
		Bio         sql.NullString `qb:"type:text; constraints:null"`
		Oscars      int            `qb:"constraints:default(0)"`
		IgnoreField string         `qb:"-"`
	}

	type Session struct {
		ID             int64     `qb:"type:bigserial; constraints:primary_key"`
		UserID         string    `qb:"type:uuid; constraints:ref(user._id)"`
		AuthToken      string    `qb:"type:uuid; constraints:notnull, unique; index"`
		CreatedAt      time.Time `qb:"constraints:notnull"`
		ExpiresAt      time.Time `qb:"constraints:notnull"`
		CompositeIndex `qb:"index:created_at, expires_at"`
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

	statement := Insert(suite.db.T("user")).Values(map[string]interface{}{
		"_id":       "b6f8bfe3-a830-441a-a097-1777e6bfae95",
		"email":     "jack@nicholson.com",
		"full_name": "Jack Nicholson",
		"bio":       sql.NullString{},
	})

	_, err = suite.db.Engine().Exec(statement)
	assert.NotNil(suite.T(), err)
	fmt.Println("Duplicate error; ", err)

	statement = Insert(suite.db.T("user")).
		Values(map[string]interface{}{
			"_id":       "cf28d117-a12d-4b75-acd8-73a7d3cbb15f",
			"email":     "jack@nicholson2.com",
			"full_name": "Jack Nicholson",
			"bio":       sql.NullString{},
		})

	_, err = suite.db.Engine().Exec(statement)
	assert.Nil(suite.T(), err)

	err = suite.db.Rollback()
	assert.NotNil(suite.T(), err)

	// find user using QueryRow()
	sel := suite.db.Find(&User{ID: "cf28d117-a12d-4b75-acd8-73a7d3cbb15f"}).Builder()
	row := suite.db.Engine().QueryRow(sel)
	assert.NotNil(suite.T(), row)

	// find user using Query()
	sel = suite.db.Find(&User{ID: "cf28d117-a12d-4b75-acd8-73a7d3cbb15f"}).Builder()
	rows, err := suite.db.Engine().Query(sel)
	assert.Nil(suite.T(), err)
	rowLength := 0
	for rows.Next() {
		rowLength++
	}
	assert.Equal(suite.T(), rowLength, 1)

	// find user using session api's Find()
	var user User

	suite.db.Find(&User{ID: "b6f8bfe3-a830-441a-a097-1777e6bfae95"}).One(&user)

	assert.Equal(suite.T(), user.Email, "jack@nicholson.com")
	assert.Equal(suite.T(), user.FullName, "Jack Nicholson")
	assert.Equal(suite.T(), user.Bio.String, "Jack Nicholson, an American actor, producer, screen-writer and director, is a three-time Academy Award winner and twelve-time nominee.")

	// select using join
	sessions := []Session{}

	err = suite.db.Query(
		suite.db.T("session").C("user_id"),
		suite.db.T("session").C("id"),
		suite.db.T("session").C("auth_token"),
		suite.db.T("session").C("created_at"),
		suite.db.T("session").C("expires_at")).
		InnerJoin(suite.db.T("user"), suite.db.T("session").C("user_id"), suite.db.T("user").C("_id")).
		Filter(suite.db.T("session").C("user_id").Eq("b6f8bfe3-a830-441a-a097-1777e6bfae95")).
		All(&sessions)

	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), len(sessions), 1)

	assert.Equal(suite.T(), sessions[0].ID, int64(1))
	assert.Equal(suite.T(), sessions[0].UserID, "b6f8bfe3-a830-441a-a097-1777e6bfae95")
	assert.Equal(suite.T(), sessions[0].AuthToken, "e4968197-6137-47a4-ba79-690d8c552248")

	// update user
	user.Bio = sql.NullString{String: "nil", Valid: false}
	suite.db.Add(user)

	err = suite.db.Commit()
	assert.Nil(suite.T(), err)

	suite.db.Find(&User{ID: "b6f8bfe3-a830-441a-a097-1777e6bfae95"}).One(&user)
	assert.Equal(suite.T(), user.Bio, sql.NullString{String: "", Valid: false})

	// delete session
	suite.db.Delete(&Session{AuthToken: "99e591f8-1025-41ef-a833-6904a0f89a38"})
	err = suite.db.Commit()
	assert.Nil(suite.T(), err)

	// drop tables
	assert.Nil(suite.T(), suite.db.DropAll())
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
