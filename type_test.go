package qb

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TypeTestSuite struct {
	suite.Suite
}

func (suite *TypeTestSuite) TestTypes() {
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
	assert.Equal(suite.T(), "BLOB", dialect.CompileType(Blob()))
}

func (suite *TypeTestSuite) TestUnsigned() {
	assert.Equal(suite.T(), "BIGINT", DefaultCompileType(BigInt().Signed(), true))
	assert.Equal(suite.T(), "BIGINT UNSIGNED", DefaultCompileType(BigInt().Unsigned(), true))
	assert.Equal(suite.T(), "NUMERIC(2, 5) UNSIGNED", DefaultCompileType(Numeric().Precision(2, 5).Unsigned(), true))

	assert.Equal(suite.T(), "INT", DefaultCompileType(Int().Signed(), false))
	assert.Equal(suite.T(), "SMALLINT", DefaultCompileType(TinyInt().Unsigned(), false))
	assert.Equal(suite.T(), "INT", DefaultCompileType(SmallInt().Unsigned(), false))
	assert.Equal(suite.T(), "BIGINT", DefaultCompileType(Int().Unsigned(), false))
	assert.Equal(suite.T(), "BIGINT", DefaultCompileType(BigInt().Unsigned(), false))
}

func TestTypeTestSuite(t *testing.T) {
	suite.Run(t, new(TypeTestSuite))
}
