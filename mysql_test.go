package qb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type mUser struct {
	ID       string `qb:"constraints:primary_key"`
	Email    string `qb:"constraints:unique, notnull"`
	FullName string `qb:"constraints:notnull"`
	Password string `qb:"constraints:notnull"`
	Bio      string `qb:"type:text; constraints:null"`
}

type mSession struct {
	ID        int64     `qb:"type:bigint; constraints:primary_key, auto_increment"`
	UserID    string    `qb:"constraints:ref(m_user.id)"`
	AuthToken string    `qb:"constraints:notnull, unique"`
	CreatedAt time.Time `qb:"constraints:notnull"`
	ExpiresAt time.Time `qb:"constraints:notnull"`
}

type MysqlTestSuite struct {
	suite.Suite
	metadata *MetaData
	dialect  *Dialect
}

func (suite *MysqlTestSuite) SetupTest() {
	engine, err := NewEngine("mysql", "root:@tcp(localhost:3306)/qb_test?charset=utf8")
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), engine)
	suite.dialect = NewDialect(engine.Driver())
	suite.metadata = NewMetaData(engine)
}

func (suite *MysqlTestSuite) TestMysql() {

	var err error

	// create tables
	suite.metadata.Add(mUser{})
	suite.metadata.Add(mSession{})
	err = suite.metadata.CreateAll()
	assert.Nil(suite.T(), err)

	// insert user
	insUser := suite.dialect.
		Insert("m_user", "id", "email", "full_name", "password", "bio").
		Values("b6f8bfe3-a830-441a-a097-1777e6bfae95", "jack@nicholson.com", "Jack Nicholson", "jack-nicholson", "Jack Nicholson, an American actor, producer, screen-writer and director, is a three-time Academy Award winner and twelve-time nominee.").
		Query()

	fmt.Println(insUser.SQL())
	fmt.Println(insUser.Bindings())
	_, err = suite.metadata.Engine().Exec(insUser)
	assert.Nil(suite.T(), err)

	// insert session
	insSession := suite.dialect.
		Insert("m_session", "user_id", "auth_token", "created_at", "expires_at").
		Values("b6f8bfe3-a830-441a-a097-1777e6bfae95", "e4968197-6137-47a4-ba79-690d8c552248", time.Now(), time.Now().Add(24*time.Hour)).
		Query()

	_, err = suite.metadata.Engine().Exec(insSession)

	fmt.Println(insSession.SQL())
	fmt.Println(insSession.Bindings())
	assert.Nil(suite.T(), err)

	// select user
	selUser := suite.dialect.
		Select("id", "email", "full_name", "bio").
		From("m_user").
		Where("m_user.id = ?", "b6f8bfe3-a830-441a-a097-1777e6bfae95").
		Query()

	var user mUser
	suite.metadata.Engine().QueryRow(selUser).Scan(&user.ID, &user.Email, &user.FullName, &user.Bio)

	assert.Equal(suite.T(), user.ID, "b6f8bfe3-a830-441a-a097-1777e6bfae95")
	assert.Equal(suite.T(), user.Email, "jack@nicholson.com")
	assert.Equal(suite.T(), user.FullName, "Jack Nicholson")
	assert.Equal(suite.T(), user.Bio, "Jack Nicholson, an American actor, producer, screen-writer and director, is a three-time Academy Award winner and twelve-time nominee.")

	// select sessions
	selSessions := suite.dialect.
		Select("s.id", "s.auth_token", "s.created_at", "s.expires_at").
		From("m_user u").
		InnerJoin("m_session s", "u.id = s.user_id").
		Where("u.id = ?", "b6f8bfe3-a830-441a-a097-1777e6bfae95").
		Query()

	rows, err := suite.metadata.Engine().Query(selSessions)
	assert.Nil(suite.T(), err)
	if err != nil {
		defer rows.Close()
	}

	sessions := []mSession{}

	for rows.Next() {
		var session mSession
		rows.Scan(&session.ID, &session.AuthToken, &session.CreatedAt, &session.ExpiresAt)
		assert.Equal(suite.T(), session.ID, int64(1))
		assert.NotNil(suite.T(), session.CreatedAt)
		assert.NotNil(suite.T(), session.ExpiresAt)
		sessions = append(sessions, session)
	}

	assert.Equal(suite.T(), len(sessions), 1)

	// update session
	query := suite.dialect.
		Update("m_session").
		Set(
			map[string]interface{}{
				"auth_token": "99e591f8-1025-41ef-a833-6904a0f89a38",
			},
		).
		Where("id = ?", 1).Query()

	_, err = suite.metadata.Engine().Exec(query)
	assert.Nil(suite.T(), err)

	// delete session
	delSession := suite.dialect.
		Delete("m_session").
		Where("auth_token = ?", "99e591f8-1025-41ef-a833-6904a0f89a38").
		Query()

	_, err = suite.metadata.Engine().Exec(delSession)
	assert.Nil(suite.T(), err)

	// insert failure
	insFail := suite.dialect.
		Insert("m_user", "invalid_column").
		Values("invalid_value").
		Query()

	_, err = suite.metadata.Engine().Exec(insFail)
	assert.NotNil(suite.T(), err)

	// insert type failure
	//insTypeFail := suite.dialect.
	//	Insert("m_user", "email").
	//	Values(5).
	//	Query()
	//
	//_, err = suite.metadata.Engine().Exec(insTypeFail)
	//assert.NotNil(suite.T(), err)

	// drop tables
	err = suite.metadata.DropAll()
}

func TestMysqlsTestSuite(t *testing.T) {
	suite.Run(t, new(MysqlTestSuite))
}
