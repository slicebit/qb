package qbit

import "database/sql"

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

type Engine struct {
	driver string
	dsn    string
	db     *sql.DB
}

func (e *Engine) DB() *sql.DB {
	return e.db
}

func (e *Engine) Ping() error {
	return e.db.Ping()
}

func (e *Engine) Driver() string {
	return e.driver
}

func (e *Engine) Dsn() string {
	return e.dsn
}
