package qb

import (
	"fmt"
	"github.com/nu7hatch/gouuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type pUser struct {
	ID       string `qb:"type:uuid; constraints:primary_key, auto_increment"`
	Email    string `qb:"constraints:unique, notnull"`
	FullName string `qb:"constraints:notnull"`
	Password string `qb:"constraints:notnull"`
	Bio      string `qb:"type:text; constraints:null"`
}

type pSession struct {
	ID        int64     `qb:"type:bigserial; constraints:primary_key"`
	UserID    string    `qb:"type:uuid; constraints:ref(p_user.id)"`
	AuthToken string    `qb:"type:uuid; constraints:notnull, unique"`
	CreatedAt time.Time `qb:"constraints:notnull"`
	ExpiresAt time.Time `qb:"constraints:notnull"`
}

type pFailModel struct {
	ID int64 `qb:"type:notype"`
}

type PostgresTestSuite struct {
	suite.Suite
	metadata *MetaData
	dialect  *Dialect
	engine   *Engine
	session  *Session
}

func (suite *PostgresTestSuite) SetupTest() {
	engine, err := NewEngine("postgres", "user=postgres dbname=qb_test sslmode=disable")
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), engine)
	suite.engine = engine
	suite.dialect = NewDialect(engine.Driver())
	suite.metadata = NewMetaData(engine)
	suite.session = NewSession(engine)
}

func (suite *PostgresTestSuite) TestPostgres() {

	var err error

	// create tables
	suite.metadata.Add(pUser{})
	suite.metadata.Add(pSession{})

	err = suite.metadata.CreateAll()
	assert.Nil(suite.T(), err)

	// insert user using dialect
	insUser := suite.dialect.
		Insert("p_user", "id", "email", "full_name", "password", "bio").
		Values("b6f8bfe3-a830-441a-a097-1777e6bfae95", "jack@nicholson.com", "Jack Nicholson", "jack-nicholson", "Jack Nicholson, an American actor, producer, screen-writer and director, is a three-time Academy Award winner and twelve-time nominee.").
		Query()

	fmt.Println(insUser.SQL())
	fmt.Println(insUser.Bindings())
	_, err = suite.metadata.Engine().Exec(insUser)
	assert.Nil(suite.T(), err)

	// insert user using session
	ddlId, _ := uuid.NewV4()
	ddl := pUser{
		ID:       ddlId.String(),
		Email:    "daniel@day-lewis.com",
		FullName: "Daniel Day-Lewis",
		Password: "ddl",
		Bio:      "Born in London, England, Daniel Michael Blake Day-Lewis is the second child of Cecil Day-Lewis (A.K.A. Nicholas Blake) (Poet Laureate of England) and his second wife, Jill Balcon. His maternal grandfather was Sir Michael Balcon, an important figure in the history of British cinema and head of the famous Ealing Studios.",
	}

	suite.session.AddAll(ddl)
	err = suite.session.Commit()
	assert.Nil(suite.T(), err)

	// insert session using dialect
	insSession := suite.dialect.
		Insert("p_session", "user_id", "auth_token", "created_at", "expires_at").
		Values("b6f8bfe3-a830-441a-a097-1777e6bfae95", "e4968197-6137-47a4-ba79-690d8c552248", time.Now(), time.Now().Add(24*time.Hour)).
		Query()

	_, err = suite.metadata.Engine().Exec(insSession)

	fmt.Println(insSession.SQL())
	fmt.Println(insSession.Bindings())
	assert.Nil(suite.T(), err)

	// select user
	selUser := suite.dialect.
		Select("id", "email", "full_name", "bio").
		From("p_user").
		Where("p_user.id = ?", "b6f8bfe3-a830-441a-a097-1777e6bfae95").
		Query()

	var user pUser
	suite.metadata.Engine().QueryRow(selUser).Scan(&user.ID, &user.Email, &user.FullName, &user.Bio)

	assert.Equal(suite.T(), user.ID, "b6f8bfe3-a830-441a-a097-1777e6bfae95")
	assert.Equal(suite.T(), user.Email, "jack@nicholson.com")
	assert.Equal(suite.T(), user.FullName, "Jack Nicholson")
	assert.Equal(suite.T(), user.Bio, "Jack Nicholson, an American actor, producer, screen-writer and director, is a three-time Academy Award winner and twelve-time nominee.")

	// select sessions
	selSessions := suite.dialect.
		Select("s.id", "s.auth_token", "s.created_at", "s.expires_at").
		From("p_user u").
		InnerJoin("p_session s", "u.id = s.user_id").
		Where("u.id = ?", "b6f8bfe3-a830-441a-a097-1777e6bfae95").
		Query()

	rows, err := suite.metadata.Engine().Query(selSessions)
	assert.Nil(suite.T(), err)
	if err != nil {
		defer rows.Close()
	}

	sessions := []pSession{}

	for rows.Next() {
		var session pSession
		rows.Scan(&session.ID, &session.AuthToken, &session.CreatedAt, &session.ExpiresAt)
		assert.Equal(suite.T(), session.ID, int64(1))
		assert.NotNil(suite.T(), session.CreatedAt)
		assert.NotNil(suite.T(), session.ExpiresAt)
		sessions = append(sessions, session)
	}

	assert.Equal(suite.T(), len(sessions), 1)

	// update session
	query := suite.dialect.
		Update("p_session").
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
		Delete("p_session").
		Where("auth_token = ?", "99e591f8-1025-41ef-a833-6904a0f89a38").
		Query()

	_, err = suite.metadata.Engine().Exec(delSession)
	assert.Nil(suite.T(), err)

	// insert failure
	insFail := suite.dialect.
		Insert("p_user", "invalid_column").
		Values("invalid_value").
		Query()

	_, err = suite.metadata.Engine().Exec(insFail)
	assert.NotNil(suite.T(), err)

	// insert type failure
	insTypeFail := suite.dialect.
		Insert("p_user", "email").
		Values(5).
		Query()

	_, err = suite.metadata.Engine().Exec(insTypeFail)
	assert.NotNil(suite.T(), err)

	// drop tables
	err = suite.metadata.DropAll()
	assert.Nil(suite.T(), err)

	// metadata create all fail
	metadata := NewMetaData(suite.engine)
	metadata.Add(pFailModel{})

	assert.NotNil(suite.T(), metadata.CreateAll())
	assert.NotNil(suite.T(), metadata.DropAll())
}

func TestPostgresTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}
