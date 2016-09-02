package qb

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TypeTestSuite struct {
	suite.Suite
}

func (suite *TypeTestSuite) TestConstraints() {
	dialect := NewDialect("")

	precisionType := Type("FLOAT").Precision(2, 5)

	assert.Equal(suite.T(), "FLOAT(2, 5)", dialect.CompileType(precisionType))

	assert.Equal(suite.T(), "CHAR", dialect.CompileType(Char()))
	assert.Equal(suite.T(), "VARCHAR(255)", dialect.CompileType(Varchar()))
	assert.Equal(suite.T(), "TEXT", dialect.CompileType(Text()))
	assert.Equal(suite.T(), "INT", dialect.CompileType(Int()))
	assert.Equal(suite.T(), "SMALLINT", dialect.CompileType(SmallInt()))
	assert.Equal(suite.T(), "BIGINT", dialect.CompileType(BigInt()))
	assert.Equal(suite.T(), "NUMERIC(2, 5)", dialect.CompileType(Numeric().Precision(2, 5)))
	assert.Equal(suite.T(), "DECIMAL", dialect.CompileType(Decimal()))
	assert.Equal(suite.T(), "FLOAT", dialect.CompileType(Float()))
	assert.Equal(suite.T(), "BOOLEAN", dialect.CompileType(Boolean()))
	assert.Equal(suite.T(), "TIMESTAMP", dialect.CompileType(Timestamp()))
}

func TestTypeTestSuite(t *testing.T) {
	suite.Run(t, new(TypeTestSuite))
}

func (suite *TypeTestSuite) TestUnsigned() {
	dialect := NewDialect("mysql")
	assert.Equal(suite.T(), "BIGINT", dialect.CompileType(BigInt().Signed()))
	assert.Equal(suite.T(), "BIGINT UNSIGNED", dialect.CompileType(BigInt().Unsigned()))
	assert.Equal(suite.T(), "NUMERIC(2, 5) UNSIGNED", dialect.CompileType(Numeric().Precision(2, 5).Unsigned()))

	dialect = NewDialect("")
	assert.Equal(suite.T(), "INT", dialect.CompileType(Int().Signed()))
	assert.Equal(suite.T(), "SMALLINT", dialect.CompileType(TinyInt().Unsigned()))
	assert.Equal(suite.T(), "INT", dialect.CompileType(SmallInt().Unsigned()))
	assert.Equal(suite.T(), "BIGINT", dialect.CompileType(Int().Unsigned()))
	assert.Equal(suite.T(), "BIGINT", dialect.CompileType(BigInt().Unsigned()))
}
