package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypes(t *testing.T) {

	assert.Equal(t, Char(40).SQL(), "CHAR(40)")
	assert.Equal(t, VarChar(40).SQL(), "VARCHAR(40)")
	assert.Equal(t, Text().SQL(), "TEXT")
	assert.Equal(t, MediumText().SQL(), "MEDIUMTEXT")
	assert.Equal(t, LongText().SQL(), "LONGTEXT")

	assert.Equal(t, BigInt().SQL(), "BIGINT")
	assert.Equal(t, Int().SQL(), "INT")
	assert.Equal(t, SmallInt().SQL(), "SMALLINT")

	assert.Equal(t, Serial().SQL(), "SERIAL")
	assert.Equal(t, BigSerial().SQL(), "BIGSERIAL")

	assert.Equal(t, Numeric().SQL(), "NUMERIC(6, 2)")
	assert.Equal(t, Numeric(7).SQL(), "NUMERIC(7, 2)")
	assert.Equal(t, Numeric(6, 3).SQL(), "NUMERIC(6, 3)")
	assert.Equal(t, Float().SQL(), "FLOAT")
	assert.Equal(t, Float(6).SQL(), "FLOAT(6)")
	assert.Equal(t, Double(6, 3).SQL(), "DOUBLE(6, 3)")
	assert.Equal(t, DoublePrecision().SQL(), "DOUBLE PRECISION")

	assert.Equal(t, Date().SQL(), "DATE")
	assert.Equal(t, Time().SQL(), "TIME")
	assert.Equal(t, DateTime().SQL(), "DATETIME")
	assert.Equal(t, Timestamp().SQL(), "TIMESTAMP")
	assert.Equal(t, Year().SQL(), "YEAR")
	assert.Equal(t, Interval(4).SQL(), "INTERVAL(4)")

	assert.Equal(t, Bytea().SQL(), "BYTEA")
	assert.Equal(t, Blob(400).SQL(), "BLOB(400)")
	assert.Equal(t, MediumBlob(400).SQL(), "MEDIUMBLOB(400)")
	assert.Equal(t, LongBlob().SQL(), "LONGBLOB")

	assert.Equal(t, Money().SQL(), "MONEY")

	assert.Equal(t, Boolean().SQL(), "BOOLEAN")

	assert.Equal(t, UUID().SQL(), "UUID")

	assert.Equal(t, Enum("premium", "trial", "free").SQL(), "ENUM('premium', 'trial', 'free')")
}
