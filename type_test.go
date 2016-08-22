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

	assert.Equal(suite.T(), "VARCHAR(255) UNIQUE NOT NULL DEFAULT 'hello'", sizeType.String(dialect))

	precisionType := Type("FLOAT").Precision(2, 5).Null()

	assert.Equal(suite.T(), "FLOAT(2, 5) NULL", precisionType.String(dialect))

	assert.Equal(suite.T(), "CHAR", Char().String(dialect))
	assert.Equal(suite.T(), "VARCHAR(255)", Varchar().String(dialect))
	assert.Equal(suite.T(), "TEXT", Text().String(dialect))
	assert.Equal(suite.T(), "INT", Int().String(dialect))
	assert.Equal(suite.T(), "SMALLINT", SmallInt().String(dialect))
	assert.Equal(suite.T(), "BIGINT", BigInt().String(dialect))
	assert.Equal(suite.T(), "NUMERIC(2, 5)", Numeric().Precision(2, 5).String(dialect))
	assert.Equal(suite.T(), "DECIMAL", Decimal().String(dialect))
	assert.Equal(suite.T(), "FLOAT", Float().String(dialect))
	assert.Equal(suite.T(), "BOOLEAN", Boolean().String(dialect))
	assert.Equal(suite.T(), "TIMESTAMP", Timestamp().String(dialect))

	assert.Equal(suite.T(), "INT TEST", Int().Constraint("TEST").String(dialect))
}

func TestTypeTestSuite(t *testing.T) {
	suite.Run(t, new(TypeTestSuite))
}

func (suite *TypeTestSuite) TestUnsigned() {
	dialect := NewDialect("mysql")
	assert.Equal(suite.T(), "BIGINT", BigInt().Signed().String(dialect))
	assert.Equal(suite.T(), "BIGINT UNSIGNED", BigInt().Unsigned().String(dialect))
	assert.Equal(suite.T(), "NUMERIC(2, 5) UNSIGNED", Numeric().Precision(2, 5).Unsigned().String(dialect))

	dialect = NewDialect("")
	assert.Equal(suite.T(), "INT", Int().Signed().String(dialect))
	assert.Equal(suite.T(), "SMALLINT", TinyInt().Unsigned().String(dialect))
	assert.Equal(suite.T(), "INT", SmallInt().Unsigned().String(dialect))
	assert.Equal(suite.T(), "BIGINT", Int().Unsigned().String(dialect))
	assert.Equal(suite.T(), "BIGINT", BigInt().Unsigned().String(dialect))
}
