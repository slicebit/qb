package qb

import (
	"fmt"
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

func TestInvalidEngine(t *testing.T) {

	engine, err := NewEngine("invalid", "")
	assert.NotEqual(t, err, nil)
	assert.Equal(t, engine, (*Engine)(nil))
}

func TestEngineExec(t *testing.T) {

	engine, err := NewEngine("postgres", "user=root dbname=pqtest")

	query, err := NewBuilder(engine.Driver()).Insert("user", "full_name").Values("Aras Can Akin").Build()
	fmt.Println(query)
	assert.Equal(t, err, nil)

	res, err := engine.Exec(query)
	assert.Equal(t, res, nil)
	assert.NotEqual(t, err, nil)
}
