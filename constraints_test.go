package qbit

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConstraints(t *testing.T) {

	assert.Equal(t, NotNull(), Constraint{"NOT NULL"})
	assert.Equal(t, Default(5), Constraint{"DEFAULT `5`"})
	assert.Equal(t, Default("-"), Constraint{"DEFAULT `-`"})
}
