package qb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhere(t *testing.T) {
	assert.Equal(t,
		"WHERE X", asDefSQL(
			Where(SQLText("X"))))
	assert.Equal(t,
		"WHERE (X AND Y)", asDefSQL(
			Where(SQLText("X"), SQLText("Y"))))
}
