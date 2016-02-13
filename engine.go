package qb

import "database/sql"

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
func (e *Engine) Exec(sql string, bindings []interface{}) (sql.Result, error) {

	stmt, err := e.db.Prepare(sql)
	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(bindings)
	if err != nil {
		return nil, err
	}

	return res, nil
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
