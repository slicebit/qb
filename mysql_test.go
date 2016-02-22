package qb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type mysqlUser struct {
	ID       int64  `qb:"constraints:primary_key,auto_increment"`
	Email    string `qb:"constraints:unique,notnull"`
	FullName string `qb:"constraints:notnull"`
	Password string `qb:"constraints:notnull"`
	Bio      string `qb:"type:text; constraints:null"`
}

type mysqlSession struct {
	SessionID int64     `qb:"constraints:primary_key"`
	UserID    int64     `qb:"constraints:ref(mysql_user.id)"`
	CreatedAt time.Time `qb:"constraints:notnull"`
	ExpiresAt time.Time `qb:"constraints:notnull"`
}

var mysqlMetadata *MetaData

func TestMysqlSetup(t *testing.T) {
	mysqlEngine, err := NewEngine("mysql", "root:@tcp(localhost:3306)/qb_test?charset=utf8")
	assert.Nil(t, err)
	assert.Nil(t, mysqlEngine.Ping())
	assert.NotNil(t, mysqlEngine)
	mysqlMetadata = NewMetaData(mysqlEngine)
}

func TestMysqlCreateTables(t *testing.T) {
	mysqlMetadata.Add(mysqlUser{})
	mysqlMetadata.Add(mysqlSession{})
	err := mysqlMetadata.CreateAll()
	assert.Nil(t, err)
}

func TestMysqlInsertSampleData(t *testing.T) {

	jn := mysqlMetadata.Table("mysql_user").Insert(map[string]interface{}{
		"email":     "jack@nicholson.com",
		"full_name": "Jack Nicholson",
		"password":  "jack-nicholson",
		"bio":       "Jack Nicholson, an American actor, producer, screen-writer and director, is a three-time Academy Award winner and twelve-time nominee.",
	})

	mb := mysqlMetadata.Table("mysql_user").Insert(map[string]interface{}{
		"email":     "marlon@brando.com",
		"full_name": "Marlon Brando",
		"password":  "marlon-brando",
		"bio":       "Marlon Brando is widely considered the greatest movie actor of all time, rivaled only by the more theatrically oriented Laurence Olivier in terms of esteem.",
	})

	_, err := mysqlMetadata.Engine().Exec(jn)
	assert.Nil(t, err)

	fmt.Println(mb.SQL())
	fmt.Println(mb.Bindings())

	_, err = mysqlMetadata.Engine().Exec(mb)
	assert.Nil(t, err)
}

func TestMysqlInsertFail(t *testing.T) {

	ins := mysqlMetadata.Table("mysql_user").Insert(map[string]interface{}{
		"invalid_column": "invalid_value",
	})

	_, err := mysqlMetadata.Engine().Exec(ins)
	assert.NotNil(t, err)
}

func TestMysqlDropTables(t *testing.T) {
	defer mysqlMetadata.Engine().DB().Close()
	err := mysqlMetadata.DropAll()
	assert.Nil(t, err)
}
