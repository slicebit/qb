package qb

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/serenize/snaker"
)

// New generates a new engine and returns it as an engine pointer
func New(driver string, dsn string) (*Engine, error) {
	conn, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	// set name mapper function
	conn.MapperFunc(func(name string) string {
		return snaker.CamelToSnake(name)
	})

	return &Engine{
		dialect: NewDialect(driver),
		dsn:     dsn,
		db:      conn,
		logger:  &DefaultLogger{LDefault, log.New(os.Stdout, "", -1)},
	}, err
}

// Engine is the generic struct for handling db connections
type Engine struct {
	dsn     string
	db      *sqlx.DB
	dialect Dialect
	logger  Logger
}

// Dialect returns the engine dialect
func (e Engine) Dialect() Dialect {
	return e.dialect
}

// SetDialect sets the current engine dialect
func (e Engine) SetDialect(dialect Dialect) {
	e.dialect = dialect
}

// TranslateError translates the native errors into qb.Error
func (e Engine) TranslateError(err error) error {
	if err != nil {
		return e.dialect.WrapError(err)
	}
	return nil
}

// Logger returns the active logger of engine
func (e *Engine) Logger() Logger {
	return e.logger
}

// SetLogger sets the logger of engine
func (e *Engine) SetLogger(logger Logger) {
	e.logger = logger
}

// SetLogFlags sets the log flags on the current logger
func (e *Engine) SetLogFlags(flags LogFlags) {
	e.logger.SetLogFlags(flags)
}

func (e *Engine) log(statement *Stmt) {
	logFlags := e.logger.LogFlags()
	if logFlags&LQuery != 0 {
		e.logger.Println("SQL:", statement.SQL())
	}
	if logFlags&LBindings != 0 {
		e.logger.Println("Bindings:", statement.Bindings())
	}
}

// Exec executes insert & update type queries and returns sql.Result and error
func (e *Engine) Exec(builder Builder) (sql.Result, error) {
	return e.ExecContext(context.Background(), builder)
}

// ExecContext executes insert & update type queries and returns sql.Result and error
func (e *Engine) ExecContext(ctx context.Context, builder Builder) (sql.Result, error) {
	statement := builder.Build(e.dialect)
	e.log(statement)
	res, err := e.db.ExecContext(ctx, statement.SQL(), statement.Bindings()...)
	return res, e.TranslateError(err)
}

// Row wraps a *sql.Row in order to translate errors
type Row struct {
	*sql.Row
	TranslateError func(error) error
}

// Scan wraps sql.Row.Scan()
func (r Row) Scan(dest ...interface{}) error {
	return r.TranslateError(r.Row.Scan(dest...))
}

// QueryRow wraps *sql.DB.QueryRow()
func (e *Engine) QueryRow(builder Builder) Row {
	return e.QueryRowContext(context.Background(), builder)
}

// QueryRowContext wraps *sql.DB.QueryRow()
func (e *Engine) QueryRowContext(ctx context.Context, builder Builder) Row {
	statement := builder.Build(e.dialect)
	e.log(statement)
	return Row{
		e.db.QueryRowContext(ctx, statement.SQL(), statement.Bindings()...),
		e.TranslateError,
	}
}

// Query wraps *sql.DB.Query()
func (e *Engine) Query(builder Builder) (*sql.Rows, error) {
	return e.QueryContext(context.Background(), builder)
}

// QueryContext wraps *sql.DB.QueryContext()
func (e *Engine) QueryContext(ctx context.Context, builder Builder) (*sql.Rows, error) {
	statement := builder.Build(e.dialect)
	e.log(statement)
	rows, err := e.db.QueryContext(ctx, statement.SQL(), statement.Bindings()...)
	return rows, e.TranslateError(err)
}

// Get maps the single row to a model
func (e *Engine) Get(builder Builder, model interface{}) error {
	return e.GetContext(context.Background(), builder, model)
}

// GetContext maps the single row to a model
func (e *Engine) GetContext(ctx context.Context, builder Builder, model interface{}) error {
	statement := builder.Build(e.dialect)
	e.log(statement)
	return e.TranslateError(
		e.db.GetContext(ctx, model, statement.SQL(), statement.Bindings()...))
}

// Select maps multiple rows to a model array
func (e *Engine) Select(builder Builder, model interface{}) error {
	return e.SelectContext(context.Background(), builder, model)
}

// SelectContext maps multiple rows to a model array
func (e *Engine) SelectContext(ctx context.Context, builder Builder, model interface{}) error {
	statement := builder.Build(e.dialect)
	e.log(statement)
	return e.TranslateError(
		e.db.SelectContext(ctx, model, statement.SQL(), statement.Bindings()...))
}

// DB returns sql.DB of wrapped engine connection
func (e *Engine) DB() *sqlx.DB {
	return e.db
}

// Ping pings the db using connection and returns error if connectivity is not present
func (e *Engine) Ping() error {
	return e.db.Ping()
}

// PingContext pings the db using connection and returns error if connectivity is not present
func (e *Engine) PingContext(ctx context.Context) error {
	return e.db.PingContext(ctx)
}

// Close closes the sqlx db connection
func (e *Engine) Close() error {
	return e.db.Close()
}

// Driver returns the driver as string
func (e *Engine) Driver() string {
	return e.dialect.Driver()
}

// Dsn returns the connection dsn
func (e *Engine) Dsn() string {
	return e.dsn
}

// Begin begins a transaction and return a *qb.Tx
func (e *Engine) Begin() (*Tx, error) {
	return e.BeginContext(context.Background())
}

// BeginContext begins a transaction and return a *qb.Tx
func (e *Engine) BeginContext(ctx context.Context) (*Tx, error) {
	tx, err := e.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, e.dialect.WrapError(err)
	}
	return &Tx{e, tx}, nil
}

// Tx is an in-progress database transaction
type Tx struct {
	engine *Engine
	tx     *sqlx.Tx
}

// Tx returns the underlying *sqlx.Tx
func (tx *Tx) Tx() *sqlx.Tx {
	return tx.tx
}

// Commit commits the transaction
func (tx *Tx) Commit() error {
	return tx.tx.Commit()
}

// Rollback aborts the transaction
func (tx *Tx) Rollback() error {
	return tx.tx.Rollback()
}

// Exec executes insert & update type queries and returns sql.Result and error
func (tx *Tx) Exec(builder Builder) (sql.Result, error) {
	return tx.ExecContext(context.Background(), builder)
}

// ExecContext executes insert & update type queries and returns sql.Result and error
func (tx *Tx) ExecContext(ctx context.Context, builder Builder) (sql.Result, error) {
	statement := builder.Build(tx.engine.dialect)
	tx.engine.log(statement)
	res, err := tx.tx.ExecContext(ctx, statement.SQL(), statement.Bindings()...)
	return res, tx.engine.TranslateError(err)
}

// QueryRow wraps *sql.DB.QueryRow()
func (tx *Tx) QueryRow(builder Builder) Row {
	return tx.QueryRowContext(context.Background(), builder)
}

// QueryRowContext wraps *sql.DB.QueryRow()
func (tx *Tx) QueryRowContext(ctx context.Context, builder Builder) Row {
	statement := builder.Build(tx.engine.dialect)
	tx.engine.log(statement)
	return Row{
		tx.tx.QueryRowContext(ctx, statement.SQL(), statement.Bindings()...),
		tx.engine.TranslateError,
	}
}

// Query wraps *sql.DB.Query()
func (tx *Tx) Query(builder Builder) (*sql.Rows, error) {
	return tx.QueryContext(context.Background(), builder)
}

// QueryContext wraps *sql.DB.QueryContext()
func (tx *Tx) QueryContext(ctx context.Context, builder Builder) (*sql.Rows, error) {
	statement := builder.Build(tx.engine.dialect)
	tx.engine.log(statement)
	rows, err := tx.tx.QueryContext(ctx, statement.SQL(), statement.Bindings()...)
	return rows, tx.engine.TranslateError(err)
}

// Get maps the single row to a model
func (tx *Tx) Get(builder Builder, model interface{}) error {
	return tx.GetContext(context.Background(), builder, model)
}

// GetContext maps the single row to a model
func (tx *Tx) GetContext(ctx context.Context, builder Builder, model interface{}) error {
	statement := builder.Build(tx.engine.dialect)
	tx.engine.log(statement)
	return tx.engine.TranslateError(
		tx.tx.GetContext(ctx, model, statement.SQL(), statement.Bindings()...))
}

// Select maps multiple rows to a model array
func (tx *Tx) Select(builder Builder, model interface{}) error {
	return tx.SelectContext(context.Background(), builder, model)
}

// SelectContext maps multiple rows to a model array
func (tx *Tx) SelectContext(ctx context.Context, builder Builder, model interface{}) error {
	statement := builder.Build(tx.engine.dialect)
	tx.engine.log(statement)
	return tx.engine.TranslateError(
		tx.tx.SelectContext(ctx, model, statement.SQL(), statement.Bindings()...))
}
