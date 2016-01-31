package qbit

import (
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEngine(t *testing.T) {

	engine, err := NewEngine("postgres", "user:password@tcp(127.0.0.1:3306)/hello")

	assert.Equal(t, err, nil)
	assert.Equal(t, engine.Driver(), "postgres")
	assert.Equal(t, engine.Ping(), engine.DB().Ping())
	assert.Equal(t, engine.Dsn(), "user:password@tcp(127.0.0.1:3306)/hello")
}
