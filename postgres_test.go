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
	ID       string  `qb:"type:uuid; constraints:primary_key, auto_increment"`
	Email    string  `qb:"constraints:unique, notnull"`
	FullName string  `qb:"constraints:notnull"`
	Bio      *string `qb:"type:text; constraints:null"`
	Oscars   int     `qb:"constraints:default(0)"`
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

type PostgresExpressionTestSuite struct {
	suite.Suite
	metadata *MetaData
	builder  *Builder
	engine   *Engine
	session  *Session
}

func (suite *PostgresExpressionTestSuite) SetupTest() {
	engine, err := NewEngine("postgres", "user=postgres dbname=qb_test sslmode=disable")
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), engine)
	suite.engine = engine
	suite.builder = NewBuilder()
	suite.metadata = NewMetaData(engine)
	suite.session = NewSession(suite.metadata)
}

func (suite *PostgresExpressionTestSuite) TestPostgresExpression() {

	var err error

	// create tables
	suite.metadata.Add(pUser{})
	suite.metadata.Add(pSession{})

	err = suite.metadata.CreateAll()
	assert.Nil(suite.T(), err)

	fmt.Println()

	// insert user using dialect
	insUserJN := suite.builder.Insert("p_user").Values(
		map[string]interface{}{
			"id":        "b6f8bfe3-a830-441a-a097-1777e6bfae95",
			"email":     "jack@nicholson.com",
			"full_name": "Jack Nicholson",
			"bio":       "Jack Nicholson, an American actor, producer, screen-writer and director, is a three-time Academy Award winner and twelve-time nominee.",
		}).Query()

	fmt.Println(insUserJN.SQL(suite.metadata.Engine().Driver()))
	fmt.Println(insUserJN.Bindings())
	fmt.Println()

	_, err = suite.metadata.Engine().Exec(insUserJN)
	assert.Nil(suite.T(), err)

	// insert user using table
	ddlID, _ := uuid.NewV4()
	insUserDDL := suite.metadata.Table("p_user").Insert(
		map[string]interface{}{
			"id":        ddlID.String(),
			"email":     "daniel@day-lewis.com",
			"full_name": "Daniel Day-Lewis",
		}).Query()

	_, err = suite.metadata.Engine().Exec(insUserDDL)
	assert.Nil(suite.T(), err)

	// insert session using dialect
	insSession := suite.builder.Insert("p_session").Values(
		map[string]interface{}{
			"user_id":    "b6f8bfe3-a830-441a-a097-1777e6bfae95",
			"auth_token": "e4968197-6137-47a4-ba79-690d8c552248",
			"created_at": time.Now(),
			"expires_at": time.Now().Add(24 * time.Hour),
		}).Query()

	_, err = suite.metadata.Engine().Exec(insSession)
	assert.Nil(suite.T(), err)

	// select user using dialect
	selUser := suite.builder.
	Select("id", "email", "full_name", "bio").
	From("p_user").
	Where("p_user.id = ?", "b6f8bfe3-a830-441a-a097-1777e6bfae95").
	Query()

	var user pUser
	suite.metadata.Engine().QueryRow(selUser).Scan(&user.ID, &user.Email, &user.FullName, &user.Bio)

	assert.Equal(suite.T(), user.ID, "b6f8bfe3-a830-441a-a097-1777e6bfae95")
	assert.Equal(suite.T(), user.Email, "jack@nicholson.com")
	assert.Equal(suite.T(), user.FullName, "Jack Nicholson")

	// select sessions
	selSessions := suite.builder.
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
		assert.True(suite.T(), session.ID >= int64(1))
		assert.NotNil(suite.T(), session.CreatedAt)
		assert.NotNil(suite.T(), session.ExpiresAt)
		sessions = append(sessions, session)
	}

	assert.Equal(suite.T(), len(sessions), 1)

	// update session
	query := suite.builder.
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
	delSession := suite.builder.
	Delete("p_session").
	Where("auth_token = ?", "99e591f8-1025-41ef-a833-6904a0f89a38").
	Query()

	_, err = suite.metadata.Engine().Exec(delSession)
	assert.Nil(suite.T(), err)

	// insert failure
	insFail := suite.builder.
	Insert("p_user").
	Values(
		map[string]interface{}{
			"invalid_column": "invalid_value",
		}).
	Query()

	_, err = suite.metadata.Engine().Exec(insFail)
	assert.NotNil(suite.T(), err)

	// insert type failure
	insTypeFail := suite.builder.
	Insert("p_user").
	Values(map[string]interface{}{
		"email": 5,
	}).Query()

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

type PostgresSessionTestSuite struct {
	suite.Suite
	metadata *MetaData
	engine   *Engine
	session  *Session
}

func (suite *PostgresSessionTestSuite) SetupTest() {
	engine, err := NewEngine("postgres", "user=postgres dbname=qb_test sslmode=disable")
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), engine)
	suite.engine = engine
	suite.metadata = NewMetaData(engine)
	suite.session = NewSession(suite.metadata)
}

func (suite *PostgresSessionTestSuite) TestPostgresSession() {

	var err error

	// create tables
	suite.metadata.Add(pUser{})
	suite.metadata.Add(pSession{})

	err = suite.metadata.CreateAll()
	assert.Nil(suite.T(), err)

	// insert user using session
	rdnID, _ := uuid.NewV4()
	rdn := pUser{
		ID:       rdnID.String(),
		Email:    "robert@de-niro.com",
		FullName: "Robert De Niro",
		Oscars:   3,
	}

	apId, _ := uuid.NewV4()
	ap := pUser{
		ID:       apId.String(),
		Email:    "al@pacino.com",
		FullName: "Al Pacino",
		Oscars:   1,
	}

	suite.session.AddAll(rdn, ap)
	err = suite.session.Commit()
	assert.Nil(suite.T(), err)

	// find first user using session
	var usr pUser

	err = suite.session.Find(pUser{ID: rdnID.String()}).First(&usr)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), usr.Email, "robert@de-niro.com")
	assert.Equal(suite.T(), usr.FullName, "Robert De Niro")
	assert.Equal(suite.T(), usr.Oscars, 3)

	fmt.Println(usr)

	//// find filter by all using session
	oneOscarUsers := []pUser{}
	suite.session.Find(&pUser{Oscars:1}).All(&oneOscarUsers)

	fmt.Println("One oscar users;")
	fmt.Println(oneOscarUsers)

	// delete user using session api
	suite.session.Delete(rdn)
	err = suite.session.Commit()
	assert.Nil(suite.T(), err)
}

func TestPostgresTestSuite(t *testing.T) {
	suite.Run(t, new(PostgresExpressionTestSuite))
	suite.Run(t, new(PostgresSessionTestSuite))
}
