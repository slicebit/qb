package qb

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
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

func (suite *DialectTestSuite) TestDefaultDialect() {
	assert.Equal(suite.T(), suite.def.SupportsUnsigned(), false)
	assert.Equal(suite.T(), suite.def.AutoIncrement(), "AUTO INCREMENT")
	assert.Equal(suite.T(), suite.def.Escape("test"), "test")
	assert.Equal(suite.T(), suite.def.Escaping(), false)
	suite.def.SetEscaping(true)
	assert.Equal(suite.T(), suite.def.Escaping(), true)
	assert.Equal(suite.T(), suite.def.Escape("test"), "`test`")
	assert.Equal(suite.T(), suite.def.EscapeAll([]string{"test"}), []string{"`test`"})
	assert.Equal(suite.T(), suite.def.Placeholder(), "?")
	assert.Equal(suite.T(), suite.def.Placeholders(5, 10), []string{"?", "?"})
	assert.Equal(suite.T(), suite.def.Driver(), "")
	suite.def.Reset() // does nothing
}

func (suite *DialectTestSuite) TestMysqlDialect() {
	assert.Equal(suite.T(), suite.mysql.SupportsUnsigned(), true)
	assert.Equal(suite.T(), suite.mysql.AutoIncrement(), "AUTO_INCREMENT")
	assert.Equal(suite.T(), suite.mysql.Escape("test"), "test")
	assert.Equal(suite.T(), suite.mysql.Escaping(), false)
	suite.mysql.SetEscaping(true)
	assert.Equal(suite.T(), suite.mysql.Escaping(), true)
	assert.Equal(suite.T(), suite.mysql.Escape("test"), "`test`")
	assert.Equal(suite.T(), suite.mysql.EscapeAll([]string{"test"}), []string{"`test`"})
	assert.Equal(suite.T(), suite.mysql.Placeholder(), "?")
	assert.Equal(suite.T(), suite.mysql.Placeholders(5, 10), []string{"?", "?"})
	assert.Equal(suite.T(), suite.mysql.Driver(), "mysql")
	suite.mysql.Reset() // does nothing
}

func (suite *DialectTestSuite) TestPostgresDialect() {
	assert.Equal(suite.T(), suite.postgres.SupportsUnsigned(), false)
	assert.Equal(suite.T(), suite.postgres.AutoIncrement(), "")
	assert.Equal(suite.T(), suite.postgres.Escape("test"), "test")
	assert.Equal(suite.T(), suite.postgres.Escaping(), false)
	suite.postgres.SetEscaping(true)
	assert.Equal(suite.T(), suite.postgres.Escaping(), true)
	assert.Equal(suite.T(), suite.postgres.Escape("test"), "\"test\"")
	assert.Equal(suite.T(), suite.postgres.EscapeAll([]string{"test"}), []string{"\"test\""})
	assert.Equal(suite.T(), suite.postgres.Placeholder(), "$1")
	assert.Equal(suite.T(), suite.postgres.Placeholder(), "$2")
	suite.postgres.Reset()
	assert.Equal(suite.T(), suite.postgres.Placeholder(), "$1")
	assert.Equal(suite.T(), suite.postgres.Placeholders(5, 10), []string{"$2", "$3"})
	assert.Equal(suite.T(), suite.postgres.Driver(), "postgres")
}

func (suite *DialectTestSuite) TestSqliteDialect() {
	assert.Equal(suite.T(), suite.sqlite.SupportsUnsigned(), false)
	assert.Equal(suite.T(), suite.sqlite.AutoIncrement(), "")
	//assert.Equal(suite.T(), suite.sqlite.AutoIncrement(), "AUTOINCREMENT")
	assert.Equal(suite.T(), suite.sqlite.Escape("test"), "test")
	assert.Equal(suite.T(), suite.sqlite.Escaping(), false)
	suite.sqlite.SetEscaping(true)
	assert.Equal(suite.T(), suite.sqlite.Escaping(), true)
	assert.Equal(suite.T(), suite.sqlite.Escape("test"), "test")
	assert.Equal(suite.T(), suite.sqlite.EscapeAll([]string{"test"}), []string{"test"})
	assert.Equal(suite.T(), suite.sqlite.Placeholder(), "?")
	assert.Equal(suite.T(), suite.sqlite.Placeholders(5, 10), []string{"?", "?"})
	assert.Equal(suite.T(), suite.sqlite.Driver(), "sqlite3")
	suite.sqlite.Reset() // does nothing
}

func TestAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(DialectTestSuite))
}
