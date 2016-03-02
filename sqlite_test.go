package qb

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"time"
)

type sUser struct {
	ID       string `qb:"type:int64; constraints:primary_key, auto_increment"`
	Email    string `qb:"constraints:unique, notnull"`
	FullName string `qb:"constraints:notnull"`
	Password string `qb:"constraints:notnull"`
	Bio      string `qb:"type:text; constraints:null"`
}

type sSession struct {
	ID        int64     `qb:"type:bigint; constraints:primary_key"`
	UserID    string    `qb:"type:uuid; constraints:ref(p_user.id)"`
	AuthToken string    `qb:"type:uuid; constraints:notnull, unique"`
	CreatedAt time.Time `qb:"constraints:notnull"`
	ExpiresAt time.Time `qb:"constraints:notnull"`
}

type SqliteTestSuite struct {
	suite.Suite
	metadata *MetaData
	engine   *Engine
}

func (suite *SqliteTestSuite) SetupTest() {
	engine, err := NewEngine("sqlite3", "./qb_test.db")
	assert.Nil(suite.T(), err)
	suite.engine = engine
	suite.metadata = NewMetaData(suite.engine)
}

func (suite *SqliteTestSuite) TestCreateTables() {
	suite.metadata.Add(sUser{})
	suite.metadata.Add(sSession{})
	err := suite.metadata.CreateAll()
	assert.NotNil(suite.T(), err)
}
