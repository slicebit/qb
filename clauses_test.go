package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSQLText(t *testing.T) {
	text := SQLText("1")
	assert.Equal(t, "1", text.Text)
}
