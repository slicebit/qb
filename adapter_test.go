package qb

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AdapterTestSuite struct {
	suite.Suite
	postgres Adapter
	mysql    Adapter
	sqlite   Adapter
	def      Adapter
}

func (suite *AdapterTestSuite) SetupTest() {
	suite.postgres = NewAdapter("postgres")
	suite.mysql = NewAdapter("mysql")
	suite.sqlite = NewAdapter("sqlite3")
	suite.def = NewAdapter("default")
}

func (suite *AdapterTestSuite) TestDefaultAdapter() {
	assert.Equal(suite.T(), suite.def.Escape("test"), "test")
	assert.Equal(suite.T(), suite.def.EscapeAll([]string{"test"}), []string{"test"})
	assert.Equal(suite.T(), suite.def.Placeholder(), "?")
	assert.Equal(suite.T(), suite.def.Placeholders(5, 10), []string{"?", "?"})
	assert.Equal(suite.T(), suite.def.SupportsInlinePrimaryKey(), false)
	assert.Equal(suite.T(), suite.def.Driver(), "")
	suite.def.Reset() // does nothing
}

func (suite *AdapterTestSuite) TestMysqlAdapter() {
	assert.Equal(suite.T(), suite.mysql.Escape("test"), "`test`")
	assert.Equal(suite.T(), suite.mysql.EscapeAll([]string{"test"}), []string{"`test`"})
	assert.Equal(suite.T(), suite.mysql.Placeholder(), "?")
	assert.Equal(suite.T(), suite.mysql.Placeholders(5, 10), []string{"?", "?"})
	assert.Equal(suite.T(), suite.mysql.SupportsInlinePrimaryKey(), false)
	assert.Equal(suite.T(), suite.mysql.Driver(), "mysql")
	suite.mysql.Reset() // does nothing
}

func (suite *AdapterTestSuite) TestPostgresAdapter() {
	assert.Equal(suite.T(), suite.postgres.Escape("test"), "\"test\"")
	assert.Equal(suite.T(), suite.postgres.EscapeAll([]string{"test"}), []string{"\"test\""})
	assert.Equal(suite.T(), suite.postgres.Placeholder(), "$1")
	assert.Equal(suite.T(), suite.postgres.Placeholder(), "$2")
	suite.postgres.Reset()
	assert.Equal(suite.T(), suite.postgres.Placeholder(), "$1")
	assert.Equal(suite.T(), suite.postgres.Placeholders(5, 10), []string{"$2", "$3"})
	assert.Equal(suite.T(), suite.postgres.SupportsInlinePrimaryKey(), true)
	assert.Equal(suite.T(), suite.postgres.Driver(), "postgres")
}

func (suite *AdapterTestSuite) TestSqliteAdapter() {
	assert.Equal(suite.T(), suite.sqlite.Escape("test"), "`test`")
	assert.Equal(suite.T(), suite.sqlite.EscapeAll([]string{"test"}), []string{"`test`"})
	assert.Equal(suite.T(), suite.sqlite.Placeholder(), "?")
	assert.Equal(suite.T(), suite.sqlite.Placeholders(5, 10), []string{"?", "?"})
	assert.Equal(suite.T(), suite.sqlite.SupportsInlinePrimaryKey(), true)
	assert.Equal(suite.T(), suite.sqlite.Driver(), "sqlite3")
	suite.sqlite.Reset() // does nothing
}

func TestAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(AdapterTestSuite))
}
