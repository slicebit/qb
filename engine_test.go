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
	dialect := NewDialect(engine.Driver())

	usersTable := Table(
		"users",
		Column("full_name", Varchar().NotNull()),
	)

	statement := Insert(usersTable).
		Values(map[string]interface{}{
			"full_name": "Al Pacino",
		}).Build(dialect)

	assert.Equal(t, err, nil)

	res, err := engine.Exec(statement)
	assert.Equal(t, res, nil)
	assert.NotNil(t, err)
}

func TestEngineFail(t *testing.T) {
	engine, err := NewEngine("sqlite3", "./qb_test.db")
	dialect := NewDialect(engine.Driver())
	assert.Nil(t, err)

	usersTable := Table(
		"users",
		Column("full_name", Varchar().NotNull()),
	)

	statement := Insert(usersTable).
		Values(map[string]interface{}{
			"full_name": "Robert De Niro",
		}).Build(dialect)

	_, err = engine.Exec(statement)
	assert.NotNil(t, err)
}
