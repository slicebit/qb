package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConstraints(t *testing.T) {
	assert.Equal(t, NotNull(), Constraint{"NOT NULL"})
	assert.Equal(t, Default(5), Constraint{"DEFAULT '5'"})
	assert.Equal(t, Default("-"), Constraint{"DEFAULT '-'"})
	assert.Equal(t, Unique(), Constraint{"UNIQUE"})
	assert.Equal(t, Unique("email", "name"), Constraint{"UNIQUE(email, name)"})
}
