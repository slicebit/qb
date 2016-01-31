package qbit

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestConstraints(t *testing.T) {

	assert.Equal(t, NotNull(), Constraint{"NOT NULL"})
	assert.Equal(t, Default(5), Constraint{"DEFAULT `5`"})
	assert.Equal(t, Default("-"), Constraint{"DEFAULT `-`"})
	assert.Equal(t, Unique(), Constraint{"UNIQUE"})
	assert.Equal(t, Unique("email", "name"), Constraint{"UNIQUE(email, name)"})
	assert.Equal(t, Key(), Constraint{"KEY"})
	assert.Equal(t, PrimaryKey(), Constraint{"PRIMARY KEY"})
	assert.Equal(t, PrimaryKey("email", "password"), Constraint{"PRIMARY KEY(email, password)"})
	assert.Equal(t, ForeignKey("user_id", "profile", "user_id"), Constraint{"FOREIGN KEY (user_id) REFERENCES profile(user_id)"})
}
