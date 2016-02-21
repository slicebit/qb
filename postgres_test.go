package qb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type PostgresUser struct {
	ID       int64  `qb:"constraints:primary_key"`
	Email    string `qb:"constraints:unique,notnull"`
	FullName string `qb:"constraints:notnull"`
	Password string `qb:"constraints:notnull"`
	Bio      string `qb:"type:text; constraints:null"`
}

type PostgresSession struct {
	SessionID int64     `qb:"constraints:primary_key"`
	UserID    int64     `qb:"constraints:ref(postgres_user.id)"`
	CreatedAt time.Time `qb:"constraints:notnull"`
	ExpiresAt time.Time `qb:"constraints:notnull"`
}

var engine *Engine
var metadata *MetaData

func TestSetup(t *testing.T) {
	var err error
	engine, err = NewEngine("postgres", "user=postgres dbname=qb_test sslmode=disable")
	assert.Nil(t, err)
	assert.NotNil(t, engine)

	metadata = NewMetaData(engine)
}

func TestCreateTables(t *testing.T) {

	metadata.Add(PostgresUser{})
	metadata.Add(PostgresSession{})
	err := metadata.CreateAll(engine)
	assert.Nil(t, err)
}

func TestInsertSampleData(t *testing.T) {

	jackNicholson, err := metadata.Table("postgres_user").Insert(map[string]interface{}{
		"id":        1,
		"email":     "jack@nicholson.com",
		"full_name": "Jack Nicholson",
		"password":  "jack-nicholson",
		"bio":       "Jack Nicholson, an American actor, producer, screen-writer and director, is a three-time Academy Award winner and twelve-time nominee.",
	})

	marlonBrando, err := metadata.Table("postgres_user").Insert(map[string]interface{}{
		"id":        2,
		"email":     "marlon@brando.com",
		"full_name": "Marlon Brando",
		"password":  "marlon-brando",
		"bio":       "Marlon Brando is widely considered the greatest movie actor of all time, rivaled only by the more theatrically oriented Laurence Olivier in terms of esteem.",
	})

	_, err = engine.Exec(jackNicholson)
	assert.Nil(t, err)

	fmt.Println(marlonBrando.SQL())
	fmt.Println(marlonBrando.Bindings())

	_, err = engine.Exec(marlonBrando)
	assert.Nil(t, err)
}

func TestInsertFail(t *testing.T) {

	aras, err := metadata.Table("postgres_user").Insert(map[string]interface{}{
		"invalid_column": "invalid_value",
	})

	_, err = engine.Exec(aras)
	assert.NotNil(t, err)
}

func TestDropTables(t *testing.T) {

	err := metadata.DropAll(engine)
	assert.Nil(t, err)
}
