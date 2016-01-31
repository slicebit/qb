package qbit

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypes(t *testing.T) {

	assert.Equal(t, SmallInt().Sql(), "SMALLINT")
	assert.Equal(t, Int().Sql(), "INT")
	assert.Equal(t, BigInt().Sql(), "BIGINT")
	assert.Equal(t, Numeric(6, 3).Sql(), "NUMERIC(6, 3)")
	assert.Equal(t, Char(40).Sql(), "CHAR(40)")
	assert.Equal(t, VarChar(40).Sql(), "VARCHAR(40)")
	assert.Equal(t, Text().Sql(), "TEXT")
	assert.Equal(t, Boolean().Sql(), "BOOLEAN")
	assert.Equal(t, Date().Sql(), "DATE")
	assert.Equal(t, DateTime().Sql(), "DATETIME")
	assert.Equal(t, Time().Sql(), "TIME")
}
