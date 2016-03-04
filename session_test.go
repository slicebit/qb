package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSession(t *testing.T) {

	engine, err := NewEngine("postgres", "user=root dbname=qb_test")

	assert.Equal(t, err, nil)
	assert.NotNil(t, engine)

	metadata := NewMetaData(engine)
	session := NewSession(metadata)

	assert.NotNil(t, session)
}
