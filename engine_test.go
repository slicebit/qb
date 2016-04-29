package qb

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

func TestInvalidEngine(t *testing.T) {
	engine, err := NewEngine("invalid", "")
	assert.NotEqual(t, err, nil)
	assert.Equal(t, engine, (*Engine)(nil))
}

func TestEngineExec(t *testing.T) {
	engine, err := NewEngine("postgres", "user=root dbname=pqtest")

	query := NewBuilder(engine.Driver()).
		Insert("user").
		Values(map[string]interface{}{
			"full_name": "Aras Can Akin",
		}).Query()
	assert.Equal(t, err, nil)

	res, err := engine.Exec(query)
	assert.Equal(t, res, nil)
	assert.NotNil(t, err)
}

func TestEngineFail(t *testing.T) {
	engine, err := NewEngine("sqlite3", "./qb_test.db")
	assert.Nil(t, err)

	query := NewBuilder(engine.Driver()).
		Insert("user").
		Values(map[string]interface{}{
			"full_name": "Aras Can Akin",
		}).Query()

	_, err = engine.Exec(query)
	assert.NotNil(t, err)
}
