package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSession(t *testing.T) {

	session, err := New("postgres", "user=root dbname=qb_test")
	assert.NotNil(t, session)
	assert.Nil(t, err)
}
