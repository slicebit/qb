package qbit

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypes(t *testing.T) {

	assert.Equal(t, Char(40).Sql(), "CHAR(40)")
	assert.Equal(t, VarChar(40).Sql(), "VARCHAR(40)")
	assert.Equal(t, Text().Sql(), "TEXT")
	assert.Equal(t, MediumText().Sql(), "MEDIUMTEXT")
	assert.Equal(t, LongText().Sql(), "LONGTEXT")

	assert.Equal(t, BigInt().Sql(), "BIGINT")
	assert.Equal(t, Int().Sql(), "INT")
	assert.Equal(t, SmallInt().Sql(), "SMALLINT")

	assert.Equal(t, Serial().Sql(), "SERIAL")
	assert.Equal(t, BigSerial().Sql(), "BIGSERIAL")

	assert.Equal(t, Numeric(6, 3).Sql(), "NUMERIC(6, 3)")
	assert.Equal(t, Float(6, 3).Sql(), "FLOAT(6, 3)")
	assert.Equal(t, Float(6).Sql(), "FLOAT(6)")
	assert.Equal(t, Double(6, 3).Sql(), "DOUBLE(6, 3)")
	assert.Equal(t, DoublePrecision().Sql(), "DOUBLE PRECISION")

	assert.Equal(t, Date().Sql(), "DATE")
	assert.Equal(t, Time().Sql(), "TIME")
	assert.Equal(t, DateTime().Sql(), "DATETIME")
	assert.Equal(t, Timestamp().Sql(), "TIMESTAMP")
	assert.Equal(t, Year().Sql(), "YEAR")
	assert.Equal(t, Interval(4).Sql(), "INTERVAL(4)")

	assert.Equal(t, Bytea().Sql(), "BYTEA")
	assert.Equal(t, Blob(400).Sql(), "BLOB(400)")
	assert.Equal(t, MediumBlob(400).Sql(), "MEDIUMBLOB(400)")
	assert.Equal(t, LongBlob(400).Sql(), "LONGBLOB(400)")

	assert.Equal(t, Money().Sql(), "MONEY")

	assert.Equal(t, Boolean().Sql(), "BOOLEAN")

	assert.Equal(t, UUID().Sql(), "UUID")

	assert.Equal(t, Enum("premium", "trial", "free").Sql(), "ENUM('premium', 'trial', 'free')")
}
