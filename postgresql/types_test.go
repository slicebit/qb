package postgresql

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypes(t *testing.T) {

	assert.Equal(t, Serial().Sql(), "SERIAL")
	assert.Equal(t, BigSerial().Sql(), "BIGSERIAL")
	assert.Equal(t, Real().Sql(), "REAL")
	assert.Equal(t, Time(1).Sql(), "TIME(1)")
	assert.Equal(t, Timestamp(1, true).Sql(), "TIMESTAMP(1) WITH TIMEZONE")
	assert.Equal(t, Timestamp(1, false).Sql(), "TIMESTAMP(1) WITHOUT TIMEZONE")
	assert.Equal(t, Interval(1).Sql(), "INTERVAL(1)")
	assert.Equal(t, Bytea().Sql(), "BYTEA")
	assert.Equal(t, Money().Sql(), "MONEY")
	assert.Equal(t, UUID().Sql(), "UUID")
}
