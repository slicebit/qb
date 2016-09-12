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
	assert.Implements(suite.T(), (*Compiler)(nil), suite.def.GetCompiler())
	assert.Equal(suite.T(), false, suite.def.SupportsUnsigned())
	assert.Equal(suite.T(), "test", suite.def.Escape("test"))
	assert.Equal(suite.T(), false, suite.def.Escaping())
	suite.def.SetEscaping(true)
	assert.Equal(suite.T(), true, suite.def.Escaping())
	assert.Equal(suite.T(), "`test`", suite.def.Escape("test"))
	assert.Equal(suite.T(), []string{"`test`"}, suite.def.EscapeAll([]string{"test"}))
	assert.Equal(suite.T(), "", suite.def.Driver())

	autoincCol := Column("id", Int()).PrimaryKey().AutoIncrement()
	assert.Equal(suite.T(),
		"INT PRIMARY KEY AUTO INCREMENT",
		suite.def.AutoIncrement(&autoincCol))

}

func (suite *DialectTestSuite) TestMysqlDialect() {
	assert.Equal(suite.T(), true, suite.mysql.SupportsUnsigned())
	assert.Equal(suite.T(), "test", suite.mysql.Escape("test"))
	assert.Equal(suite.T(), false, suite.mysql.Escaping())
	suite.mysql.SetEscaping(true)
	assert.Equal(suite.T(), true, suite.mysql.Escaping())
	assert.Equal(suite.T(), "`test`", suite.mysql.Escape("test"))
	assert.Equal(suite.T(), []string{"`test`"}, suite.mysql.EscapeAll([]string{"test"}))
	assert.Equal(suite.T(), "mysql", suite.mysql.Driver())
}

func (suite *DialectTestSuite) TestPostgresDialect() {
	assert.Equal(suite.T(), false, suite.postgres.SupportsUnsigned())
	assert.Equal(suite.T(), "test", suite.postgres.Escape("test"))
	assert.Equal(suite.T(), false, suite.postgres.Escaping())
	suite.postgres.SetEscaping(true)
	assert.Equal(suite.T(), true, suite.postgres.Escaping())
	assert.Equal(suite.T(), "\"test\"", suite.postgres.Escape("test"))
	assert.Equal(suite.T(), []string{"\"test\""}, suite.postgres.EscapeAll([]string{"test"}))
	assert.Equal(suite.T(), "postgres", suite.postgres.Driver())

	col := Column("autoinc", Int()).AutoIncrement()
	assert.Equal(suite.T(), "SERIAL", suite.postgres.AutoIncrement(&col))

	col = Column("autoinc", BigInt()).AutoIncrement()
	assert.Equal(suite.T(), "BIGSERIAL", suite.postgres.AutoIncrement(&col))

	col = Column("autoinc", SmallInt()).AutoIncrement()
	assert.Equal(suite.T(), "SMALLSERIAL", suite.postgres.AutoIncrement(&col))
}

func (suite *DialectTestSuite) TestSqliteDialect() {
	assert.Equal(suite.T(), false, suite.sqlite.SupportsUnsigned())
	assert.Equal(suite.T(), "test", suite.sqlite.Escape("test"))
	assert.Equal(suite.T(), false, suite.sqlite.Escaping())
	suite.sqlite.SetEscaping(true)
	assert.Equal(suite.T(), true, suite.sqlite.Escaping())
	assert.Equal(suite.T(), "test", suite.sqlite.Escape("test"))
	assert.Equal(suite.T(), []string{"test"}, suite.sqlite.EscapeAll([]string{"test"}))
	assert.Equal(suite.T(), "sqlite3", suite.sqlite.Driver())
}

func TestDialectTestSuite(t *testing.T) {
	suite.Run(t, new(DialectTestSuite))
}
