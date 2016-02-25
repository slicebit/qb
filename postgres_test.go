package qb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type postgresUser struct {
	ID       int64  `qb:"type:bigserial; constraints:primary_key,auto_increment"`
	Email    string `qb:"constraints:unique,notnull"`
	FullName string `qb:"constraints:notnull"`
	Password string `qb:"constraints:notnull"`
	Bio      string `qb:"type:text; constraints:null"`
}

type postgresSession struct {
	SessionID int64     `qb:"constraints:primary_key"`
	UserID    int64     `qb:"constraints:ref(postgres_user.id)"`
	CreatedAt time.Time `qb:"constraints:notnull"`
	ExpiresAt time.Time `qb:"constraints:notnull"`
}

var postgresMetadata *MetaData

func TestPostgresSetup(t *testing.T) {
	postgresEngine, err := NewEngine("postgres", "user=postgres dbname=qb_test sslmode=disable")
	assert.Nil(t, err)
	assert.Nil(t, postgresEngine.Ping())
	assert.NotNil(t, postgresEngine)
	postgresMetadata = NewMetaData(postgresEngine)
}

func TestPostgresCreateTables(t *testing.T) {
	postgresMetadata.Add(postgresUser{})
	postgresMetadata.Add(postgresSession{})
	err := postgresMetadata.CreateAll()
	assert.Nil(t, err)
}

func TestPostgresInsertSampleData(t *testing.T) {

	jn := postgresMetadata.Table("postgres_user").Insert(map[string]interface{}{
		"email":     "jack@nicholson.com",
		"full_name": "Jack Nicholson",
		"password":  "jack-nicholson",
		"bio":       "Jack Nicholson, an American actor, producer, screen-writer and director, is a three-time Academy Award winner and twelve-time nominee.",
	})

	mb := postgresMetadata.Table("postgres_user").Insert(map[string]interface{}{
		"email":     "marlon@brando.com",
		"full_name": "Marlon Brando",
		"password":  "marlon-brando",
		"bio":       "Marlon Brando is widely considered the greatest movie actor of all time, rivaled only by the more theatrically oriented Laurence Olivier in terms of esteem.",
	})

	_, err := postgresMetadata.Engine().Exec(jn)
	assert.Nil(t, err)

	fmt.Println(mb.SQL())
	fmt.Println(mb.Bindings())

	_, err = postgresMetadata.Engine().Exec(mb)
	assert.Nil(t, err)
}

func TestPostgresInsertFail(t *testing.T) {

	ins := postgresMetadata.Table("postgres_user").Insert(map[string]interface{}{
		"invalid_column": "invalid_value",
	})

	_, err := postgresMetadata.Engine().Exec(ins)
	assert.NotNil(t, err)
}

func TestPostgresInsertTypeFail(t *testing.T) {

	ins := postgresMetadata.Table("postgres_user").Insert(map[string]interface{}{
		"email": 5,
	})

	_, err := postgresMetadata.Engine().Exec(ins)
	assert.NotNil(t, err)
}

func TestPostgresDropTables(t *testing.T) {
	defer postgresMetadata.Engine().DB().Close()
	err := postgresMetadata.DropAll()
	assert.Nil(t, err)
}
