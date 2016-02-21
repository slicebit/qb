package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSession(t *testing.T) {

	engine, err := NewEngine("postgres", "user=root dbname=qb_test")

	assert.Equal(t, err, nil)

	session := NewSession(engine)

	assert.Equal(t, session.Engine(), engine)
}
