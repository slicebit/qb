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

	sizeType := Varchar().Size(255).Unique().NotNull().Default("hello")

	assert.Equal(suite.T(), sizeType.String(dialect), "VARCHAR(255) UNIQUE NOT NULL DEFAULT 'hello'")

	precisionType := Type("FLOAT").Precision(2, 5).Null()

	assert.Equal(suite.T(), precisionType.String(dialect), "FLOAT(2, 5) NULL")

	assert.Equal(suite.T(), Char().String(dialect), "CHAR")
	assert.Equal(suite.T(), Varchar().String(dialect), "VARCHAR(255)")
	assert.Equal(suite.T(), Text().String(dialect), "TEXT")
	assert.Equal(suite.T(), Int().String(dialect), "INT")
	assert.Equal(suite.T(), SmallInt().String(dialect), "SMALLINT")
	assert.Equal(suite.T(), BigInt().String(dialect), "BIGINT")
	assert.Equal(suite.T(), Numeric().Precision(2, 5).String(dialect), "NUMERIC(2, 5)")
	assert.Equal(suite.T(), Decimal().String(dialect), "DECIMAL")
	assert.Equal(suite.T(), Float().String(dialect), "FLOAT")
	assert.Equal(suite.T(), Boolean().String(dialect), "BOOLEAN")
	assert.Equal(suite.T(), Timestamp().String(dialect), "TIMESTAMP")
}

func TestTypeTestSuite(t *testing.T) {
	suite.Run(t, new(TypeTestSuite))
}

func (suite *TypeTestSuite) TestUnsigned() {
	dialect := NewDialect("mysql")
	assert.Equal(suite.T(), BigInt().Signed().String(dialect), "BIGINT")
	assert.Equal(suite.T(), BigInt().Unsigned().String(dialect), "BIGINT UNSIGNED")
	assert.Equal(suite.T(), Numeric().Precision(2, 5).Unsigned().String(dialect), "NUMERIC(2, 5) UNSIGNED")

	dialect = NewDialect("")
	assert.Equal(suite.T(), Int().Signed().String(dialect), "INT")
	assert.Equal(suite.T(), SmallInt().Unsigned().String(dialect), "INT")
	assert.Equal(suite.T(), Int().Unsigned().String(dialect), "BIGINT")
	assert.Equal(suite.T(), BigInt().Unsigned().String(dialect), "BIGINT")
}
