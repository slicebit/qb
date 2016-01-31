package mysql

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTypes(t *testing.T) {

	assert.Equal(t, Date().Sql(), "DATE")
	assert.Equal(t, DateTime().Sql(), "DATETIME")
	assert.Equal(t, Timestamp().Sql(), "TIMESTAMP")
	assert.Equal(t, Time().Sql(), "TIME")
	assert.Equal(t, Time().Sql(), "TIME")
	assert.Equal(t, TinyText().Sql(), "TINYTEXT")
	assert.Equal(t, Enum([]string{"PREMIUM", "FREETRIAL"}).Sql(), "ENUM('PREMIUM', 'FREETRIAL')")

}