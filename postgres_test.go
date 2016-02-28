package qb

import (
	"github.com/stretchr/testify/assert"
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

var pMetadata *MetaData

func TestPostgresSetup(t *testing.T) {
	pEngine, err := NewEngine("postgres", "user=postgres dbname=qb_test sslmode=disable")
	assert.Nil(t, err)
	assert.Nil(t, pEngine.Ping())
	assert.NotNil(t, pEngine)
	pMetadata = NewMetaData(pEngine)
}

func TestPostgresCreateTables(t *testing.T) {
	pMetadata.Add(pUser{})
	pMetadata.Add(pSession{})
	err := pMetadata.CreateAll()
	assert.Nil(t, err)
}

func TestPostgresInsertSampleData(t *testing.T) {

	insUser := pMetadata.Table("p_user").Insert(map[string]interface{}{
		"id":        "b6f8bfe3-a830-441a-a097-1777e6bfae95",
		"email":     "jack@nicholson.com",
		"full_name": "Jack Nicholson",
		"password":  "jack-nicholson",
		"bio":       "Jack Nicholson, an American actor, producer, screen-writer and director, is a three-time Academy Award winner and twelve-time nominee.",
	})

	_, err := pMetadata.Engine().Exec(insUser)
	assert.Nil(t, err)

	insSession := pMetadata.Table("p_session").Insert(map[string]interface{}{
		"user_id":    "b6f8bfe3-a830-441a-a097-1777e6bfae95",
		"auth_token": "e4968197-6137-47a4-ba79-690d8c552248",
		"created_at": time.Now(),
		"expires_at": time.Now().Add(24 * time.Hour),
	})

	_, err = pMetadata.Engine().Exec(insSession)

	assert.Nil(t, err)
}

func TestPostgresSelectUser(t *testing.T) {

	query := NewBuilder(pMetadata.Engine().Driver()).Select("id", "email", "full_name", "bio").
		From("p_user").
		Where("p_user.id = ?", "b6f8bfe3-a830-441a-a097-1777e6bfae95").
		Build()

	var user pUser
	pMetadata.Engine().QueryRow(query).Scan(&user.ID, &user.Email, &user.FullName, &user.Bio)

	assert.Equal(t, user.ID, "b6f8bfe3-a830-441a-a097-1777e6bfae95")
	assert.Equal(t, user.Email, "jack@nicholson.com")
	assert.Equal(t, user.FullName, "Jack Nicholson")
	assert.Equal(t, user.Bio, "Jack Nicholson, an American actor, producer, screen-writer and director, is a three-time Academy Award winner and twelve-time nominee.")
}

func TestPostgresSelectSessions(t *testing.T) {

	query := NewBuilder(pMetadata.Engine().Driver()).
		Select("s.id", "s.auth_token", "s.created_at", "s.expires_at").
		From("p_user u").
		InnerJoin("p_session s", "u.id = s.user_id").
		Where("u.id = ?", "b6f8bfe3-a830-441a-a097-1777e6bfae95").
		Build()

	rows, err := pMetadata.Engine().Query(query)
	assert.Nil(t, err)
	if err != nil {
		defer rows.Close()
	}

	sessions := []pSession{}

	for rows.Next() {
		var session pSession
		rows.Scan(&session.ID, &session.AuthToken, &session.CreatedAt, &session.ExpiresAt)
		assert.Equal(t, session.ID, int64(1))
		assert.NotNil(t, session.CreatedAt)
		assert.NotNil(t, session.ExpiresAt)
		sessions = append(sessions, session)
	}

	assert.Equal(t, len(sessions), 1)
}

func TestPostgresUpdateSession(t *testing.T) {

	query := NewBuilder(pMetadata.Engine().Driver()).
		Update("p_session").
		Set(
		map[string]interface{}{
			"auth_token": "99e591f8-1025-41ef-a833-6904a0f89a38",
		}).
		Where("id = ?", 1).Build()

	_, err := pMetadata.Engine().Exec(query)
	assert.Nil(t, err)
}

func TestPostgresDeleteSession(t *testing.T) {
	query := NewBuilder(pMetadata.Engine().Driver()).
		Delete("p_session").
		Where("auth_token = ?", "99e591f8-1025-41ef-a833-6904a0f89a38").
		Build()

	_, err := pMetadata.Engine().Exec(query)
	assert.Nil(t, err)
}

func TestPostgresInsertFail(t *testing.T) {

	ins := pMetadata.Table("p_user").Insert(map[string]interface{}{
		"invalid_column": "invalid_value",
	})

	_, err := pMetadata.Engine().Exec(ins)
	assert.NotNil(t, err)
}

func TestPostgresInsertTypeFail(t *testing.T) {

	ins := pMetadata.Table("p_user").Insert(map[string]interface{}{
		"email": 5,
	})

	_, err := pMetadata.Engine().Exec(ins)
	assert.NotNil(t, err)
}

func TestPostgresDropTables(t *testing.T) {
	defer pMetadata.Engine().DB().Close()
	err := pMetadata.DropAll()
	assert.Nil(t, err)
}
