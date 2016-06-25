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
	dialect := NewDialect("postgres")
	dialect.SetEscaping(true)
	engine.SetDialect(dialect)

	usersTable := Table(
		"users",
		Column("full_name", Varchar().NotNull()),
	)

	ins := Insert(usersTable).
		Values(map[string]interface{}{
			"full_name": "Al Pacino",
		})

	assert.Nil(t, err)

	res, err := engine.Exec(ins)
	assert.Equal(t, res, nil)
	assert.NotNil(t, err)
}

func TestEngineFail(t *testing.T) {
	engine, err := NewEngine("sqlite3", "./qb_test.db")
	engine.SetDialect(NewDialect("sqlite3"))
	assert.Nil(t, err)

	usersTable := Table(
		"users",
		Column("full_name", Varchar().NotNull()),
	)

	statement := Insert(usersTable).
		Values(map[string]interface{}{
			"full_name": "Robert De Niro",
		})

	_, err = engine.Exec(statement)
	assert.NotNil(t, err)
}
