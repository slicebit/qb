package qb_test

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/slicebit/qb"
	_ "github.com/slicebit/qb/dialects/sqlite"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEngine(t *testing.T) {
	engine, err := qb.New("sqlite3", ":memory:")

	assert.Equal(t, nil, err)
	assert.Equal(t, "sqlite3", engine.Driver())
	assert.Equal(t, engine.DB().Ping(), engine.Ping())
	assert.Equal(t, ":memory:", engine.Dsn())
}

func TestInvalidEngine(t *testing.T) {
	engine, err := qb.New("invalid", "")
	assert.NotEqual(t, nil, err)
	assert.Equal(t, (*qb.Engine)(nil), engine)
}

func TestEngineExec(t *testing.T) {
	engine, err := qb.New("sqlite3", ":memory:")
	dialect := qb.NewDialect("sqlite")
	dialect.SetEscaping(true)
	engine.SetDialect(dialect)

	usersTable := qb.Table(
		"users",
		qb.Column("full_name", qb.Varchar()).NotNull(),
	)

	ins := qb.Insert(usersTable).
		Values(map[string]interface{}{
			"full_name": "Al Pacino",
		})

	assert.Nil(t, err)

	res, err := engine.Exec(ins)
	assert.Equal(t, nil, res)
	assert.NotNil(t, err)
}

func TestEngineFail(t *testing.T) {
	engine, err := qb.New("sqlite3", ":memory:")
	defer engine.Close()
	engine.SetDialect(qb.NewDialect("sqlite3"))
	assert.Nil(t, err)

	usersTable := qb.Table(
		"users",
		qb.Column("full_name", qb.Varchar()).NotNull(),
	)

	statement := qb.Insert(usersTable).
		Values(map[string]interface{}{
			"full_name": "Robert De Niro",
		})

	_, err = engine.Exec(statement)
	assert.NotNil(t, err)
}

func TestTx(t *testing.T) {
	engine, err := qb.New("sqlite3", ":memory:")
	assert.Nil(t, err)
	defer engine.Close()

	engine.SetDialect(qb.NewDialect("sqlite3"))

	usersTable := qb.Table(
		"users",
		qb.Column("full_name", qb.Varchar()).NotNull(),
	)

	_, err = engine.DB().Exec(usersTable.Create(engine.Dialect()))
	assert.Nil(t, err)

	countStmt := qb.Select(qb.Count(usersTable.C("full_name"))).From(usersTable)

	tx, err := engine.Begin()
	assert.Equal(t, nil, err)

	assert.NotNil(t, tx.Tx())

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
	engine, err := qb.New("sqlite3", "file:///dev/null?_txlock=exclusive")
	assert.Nil(t, err)
	_, err = engine.Begin()
	assert.NotNil(t, err)
}

func TestEngineQuery(t *testing.T) {
	engine, err := qb.New("sqlite3", ":memory:")
	assert.Nil(t, err)
	rows, err := engine.Query(qb.Select(qb.SQLText("1")))
	assert.Nil(t, err)
	assert.True(t, rows.Next())
	var value int
	assert.Nil(t, rows.Scan(&value))
	assert.Equal(t, 1, value)
	assert.False(t, rows.Next())
}

func TestEngineGet(t *testing.T) {
	var s struct {
		Value int `db:"value"`
	}
	engine, err := qb.New("sqlite3", ":memory:")
	assert.Nil(t, err)
	assert.Nil(t, engine.Get(qb.Select(qb.SQLText("1 AS value")), &s))
	assert.Equal(t, 1, s.Value)
}

func TestEngineSelect(t *testing.T) {
	var s []struct {
		Value int `db:"value"`
	}
	engine, err := qb.New("sqlite3", ":memory:")
	assert.Nil(t, err)
	assert.Nil(t, engine.Select(qb.Select(qb.SQLText("1 AS value")), &s))
	assert.Equal(t, 1, len(s))
	assert.Equal(t, 1, s[0].Value)
}
