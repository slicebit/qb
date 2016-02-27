package qb

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// NewEngine generates a new engine and returns it as an engine pointer
func NewEngine(driver string, dsn string) (*Engine, error) {

	conn, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	return &Engine{
		driver: driver,
		dsn:    dsn,
		db:     conn,
	}, err
}

// Engine is the generic struct for handling db connections
type Engine struct {
	driver string
	dsn    string
	db     *sql.DB
}

// Exec executes insert & update type queries and returns sql.Result and error
func (e *Engine) Exec(query *Query) (sql.Result, error) {

	stmt, err := e.db.Prepare(query.SQL())
	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(query.Bindings()...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// QueryRow wraps *sql.DB.QueryRow()
func (e *Engine) QueryRow(query *Query) *sql.Row {
	return e.db.QueryRow(query.SQL(), query.Bindings()...)
}

// Query wraps *sql.DB.Query()
func (e *Engine) Query(query *Query) (*sql.Rows, error) {
	return e.db.Query(query.SQL(), query.Bindings()...)
}

// DB returns sql.DB of wrapped engine connection
func (e *Engine) DB() *sql.DB {
	return e.db
}

// Ping pings the db using connection and returns error if connectivity is not present
func (e *Engine) Ping() error {
	return e.db.Ping()
}

// Driver returns the driver as string
func (e *Engine) Driver() string {
	return e.driver
}

// Dsn returns the connection dsn
func (e *Engine) Dsn() string {
	return e.dsn
}
