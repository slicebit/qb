package qb

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEngine(t *testing.T) {
	engine, err := New("sqlite3", ":memory:")

	assert.Equal(t, nil, err)
	assert.Equal(t, "sqlite3", engine.Driver())
	assert.Equal(t, engine.DB().Ping(), engine.Ping())
	assert.Equal(t, ":memory:", engine.Dsn())
}

func TestInvalidEngine(t *testing.T) {
	engine, err := New("invalid", "")
	assert.NotEqual(t, nil, err)
	assert.Equal(t, (*Engine)(nil), engine)
}

func TestEngineExec(t *testing.T) {
	engine, err := New("sqlite3", ":memory:")
	dialect := NewDialect("sqlite")
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
	defer engine.Close()
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

func TestTx(t *testing.T) {
	engine, err := New("sqlite3", ":memory:")
	assert.Nil(t, err)
	defer engine.Close()

	engine.SetDialect(NewDialect("sqlite3"))

	usersTable := Table(
		"users",
		Column("full_name", Varchar()).NotNull(),
	)

	_, err = engine.DB().Exec(usersTable.Create(engine.Dialect()))
	assert.Nil(t, err)

	countStmt := Select(Count(usersTable.C("full_name"))).From(usersTable)

	tx, err := engine.Begin()
	assert.Equal(t, nil, err)

	assert.Equal(t, tx.tx, tx.Tx())

	_, err = tx.Exec(
		usersTable.Insert().
			Values(map[string]interface{}{
				"full_name": "Robert De Niro",
			}),
	)
	assert.Equal(t, nil, err)
	var count int
	row := tx.QueryRow(countStmt)
	assert.Nil(t, row.Scan(&count))
	assert.Equal(t, 1, count)

	tx.Commit()

	row = engine.QueryRow(countStmt)
	assert.Equal(t, nil, row.Scan(&count))
	assert.Equal(t, 1, count)

	tx, err = engine.Begin()
	assert.Equal(t, nil, err)

	_, err = tx.Exec(usersTable.Insert().
		Values(map[string]interface{}{
			"full_name": "Al Pacino",
		}),
	)
	assert.Equal(t, nil, err)

	rows, err := tx.Query(countStmt)
	assert.Equal(t, nil, err)
	assert.True(t, rows.Next())
	assert.Equal(t, nil, rows.Scan(&count))
	assert.Equal(t, 2, count)

	tx.Rollback()

	row = engine.QueryRow(countStmt)
	assert.Equal(t, nil, row.Scan(&count))
	assert.Equal(t, 1, count)

	tx, _ = engine.Begin()

	assert.Nil(t, nil)

	var user struct{ FullName string }
	var users []struct{ FullName string }

	assert.Nil(t,
		tx.Get(usersTable.Select(usersTable.C("full_name")), &user),
	)
	assert.Equal(t, "Robert De Niro", user.FullName)

	assert.Nil(t,
		tx.Select(usersTable.Select(usersTable.C("full_name")), &users),
	)
	assert.Equal(t, "Robert De Niro", users[0].FullName)

}

func TestTxBeginError(t *testing.T) {
	engine, err := New("sqlite3", "file:///dev/null?_txlock=exclusive")
	assert.Nil(t, err)
	_, err = engine.Begin()
	assert.NotNil(t, err)
}
