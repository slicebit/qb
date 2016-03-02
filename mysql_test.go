package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type mUser struct {
	ID       int64  `qb:"constraints:primary_key,auto_increment"`
	Email    string `qb:"constraints:unique,notnull"`
	FullName string `qb:"constraints:notnull"`
	Password string `qb:"constraints:notnull"`
	Bio      string `qb:"type:text; constraints:null"`
}

type mSession struct {
	SessionID int64     `qb:"constraints:primary_key"`
	UserID    int64     `qb:"constraints:ref(mysql_user.id)"`
	CreatedAt time.Time `qb:"constraints:notnull"`
	ExpiresAt time.Time `qb:"constraints:notnull"`
}

var mMetadata *MetaData

func TestMysqlSetup(t *testing.T) {
	mysqlEngine, err := NewEngine("mysql", "root:@tcp(localhost:3306)/qb_test?charset=utf8")
	assert.Nil(t, err)
	assert.Nil(t, mysqlEngine.Ping())
	assert.NotNil(t, mysqlEngine)
	mMetadata = NewMetaData(mysqlEngine)
}

func TestMysqlCreateTables(t *testing.T) {
	mMetadata.Add(mUser{})
	mMetadata.Add(mSession{})
	err := mMetadata.CreateAll()
	assert.Nil(t, err)
}