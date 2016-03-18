package qb

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"github.com/stretchr/testify/assert"
)

type DialectTestSuite struct {
	suite.Suite
	postgres Dialect
	mysql    Dialect
	sqlite   Dialect
	def      Dialect
}

func (suite *DialectTestSuite) SetupTest() {
	suite.postgres = NewDialect("postgres")
	suite.mysql = NewDialect("mysql")
	suite.sqlite = NewDialect("sqlite3")
	suite.def = NewDialect("default")
}

func (suite *DialectTestSuite) TestMysqlDialect() {
	assert.Equal(suite.T(), suite.mysql.Escape("test"), "`test`")
	assert.Equal(suite.T(), suite.mysql.Placeholder(), "?")
	assert.Equal(suite.T(), suite.mysql.SupportsInlinePrimaryKey(), false)
}

func (suite *DialectTestSuite) TestPostgresDialect() {
	assert.Equal(suite.T(), suite.postgres.Escape("test"), "\"test\"")
	assert.Equal(suite.T(), suite.postgres.Placeholder(), "$1")
	assert.Equal(suite.T(), suite.postgres.Placeholder(), "$2")
	suite.postgres.Reset()
	assert.Equal(suite.T(), suite.postgres.Placeholder(), "$1")
	assert.Equal(suite.T(), suite.postgres.SupportsInlinePrimaryKey(), true)
}

func (suite *DialectTestSuite) TestSqliteDialect() {
	assert.Equal(suite.T(), suite.sqlite.Escape("test"), "`test`")
	assert.Equal(suite.T(), suite.sqlite.Placeholder(), "?")
	assert.Equal(suite.T(), suite.sqlite.SupportsInlinePrimaryKey(), true)
}

func TestDialectTestSuite(t *testing.T) {
	suite.Run(t, new(DialectTestSuite))
}