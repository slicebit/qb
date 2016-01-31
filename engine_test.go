package qbit

import (
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEngine(t *testing.T) {

	engine, err := NewEngine("postgres", "user=root dbname=pqtest")

	assert.Equal(t, err, nil)
	assert.Equal(t, engine.Driver(), "postgres")
	assert.Equal(t, engine.Ping(), engine.DB().Ping())
	assert.Equal(t, engine.Dsn(), "user=root dbname=pqtest")
}
