package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypes(t *testing.T) {
	typ := NewType("FLOAT(4,6)")
	assert.Equal(t, typ.SQL, "FLOAT(4,6)")
}
