package qbit

import "database/sql"

func Engine(driver string, dsn string) (*engine, error) {

	conn, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	return &engine{
		driver: driver,
		dsn:    dsn,
		db:     conn,
	}, err
}

type engine struct {
	driver string
	dsn    string
	db     *sql.DB
}

func (e *engine) DB() *sql.DB {
	return e.db
}

func (e *engine) Ping() error {
	return nil
}

func (e *engine) Driver() string {
	return e.driver
}

func (e *engine) Dsn() string {
	return e.dsn
}
