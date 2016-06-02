package qb

import (
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TypeTestSuite struct {
	suite.Suite
}

func (suite *TypeTestSuite) TestConstraints() {

	sizeType := Varchar().Size(255).Unique().NotNull().Default("hello")

	assert.Equal(suite.T(), sizeType.String(), "VARCHAR(255) UNIQUE NOT NULL DEFAULT 'hello'")

	precisionType := Type("FLOAT").Precision(2, 5).Null()

	assert.Equal(suite.T(), precisionType.String(), "FLOAT(2, 5) NULL")

	assert.Equal(suite.T(), Char().String(), "CHAR")
	assert.Equal(suite.T(), Varchar().String(), "VARCHAR")
	assert.Equal(suite.T(), Text().String(), "TEXT")
	assert.Equal(suite.T(), Int().String(), "INT")
	assert.Equal(suite.T(), SmallInt().String(), "SMALLINT")
	assert.Equal(suite.T(), BigInt().String(), "BIGINT")
	assert.Equal(suite.T(), Numeric().Precision(2, 5).String(), "NUMERIC(2, 5)")
	assert.Equal(suite.T(), Decimal().String(), "DECIMAL")
	assert.Equal(suite.T(), Float().String(), "FLOAT")
	assert.Equal(suite.T(), Boolean().String(), "BOOLEAN")
	assert.Equal(suite.T(), Timestamp().String(), "TIMESTAMP")
}

func TestTypeTestSuite(t *testing.T) {
	suite.Run(t, new(TypeTestSuite))
}
