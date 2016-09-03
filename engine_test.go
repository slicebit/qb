package qb

import (
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEngine(t *testing.T) {
	engine, err := New("postgres", "user=root dbname=pqtest")

	assert.Equal(t, nil, err)
	assert.Equal(t, "postgres", engine.Driver())
	assert.Equal(t, engine.DB().Ping(), engine.Ping())
	assert.Equal(t, "user=root dbname=pqtest", engine.Dsn())
}

func TestInvalidEngine(t *testing.T) {
	engine, err := New("invalid", "")
	assert.NotEqual(t, nil, err)
	assert.Equal(t, (*Engine)(nil), engine)
}

func TestEngineExec(t *testing.T) {
	engine, err := New("postgres", "user=root dbname=pqtest")
	dialect := NewDialect("postgres")
	dialect.SetEscaping(true)
	engine.SetDialect(dialect)

	usersTable := Table(
		"users",
		Column("full_name", Varchar()).NotNull(),
	)

	ins := Insert(usersTable).
		Values(map[string]interface{}{
			"full_name": "Al Pacino",
		})

	assert.Nil(t, err)

	res, err := engine.Exec(ins)
	assert.Equal(t, nil, res)
	assert.NotNil(t, err)
}

func TestEngineFail(t *testing.T) {
	engine, err := New("sqlite3", "./qb_test.db")
	engine.SetDialect(NewDialect("sqlite3"))
	assert.Nil(t, err)

	usersTable := Table(
		"users",
		Column("full_name", Varchar()).NotNull(),
	)

	statement := Insert(usersTable).
		Values(map[string]interface{}{
			"full_name": "Robert De Niro",
		})

	_, err = engine.Exec(statement)
	assert.NotNil(t, err)
}
