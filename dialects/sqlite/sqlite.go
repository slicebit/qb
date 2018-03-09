package sqlite

import (
	"fmt"
	"strings"

	"github.com/aacanakin/qb"
	"github.com/mattn/go-sqlite3"
)

// Dialect is a type of dialect that can be used with sqlite driver
type Dialect struct {
	escaping bool
}

// NewDialect creates a new sqlite3 dialect
func NewDialect() qb.Dialect {
	return &Dialect{false}
}

func init() {
	qb.RegisterDialect("sqlite3", NewDialect())
	qb.RegisterDialect("sqlite", NewDialect())
}

// CompileType compiles a type into its DDL
func (d *Dialect) CompileType(t qb.TypeElem) string {
	if t.Name == "UUID" {
		return "VARCHAR(36)"
	}
	return qb.DefaultCompileType(t, d.SupportsUnsigned())
}

// Escape wraps the string with escape characters of the dialect
func (d *Dialect) Escape(str string) string {
	if d.escaping {
		return fmt.Sprintf(`"%s"`, str)
	}
	return str
}

// EscapeAll wraps all elements of string array
func (d *Dialect) EscapeAll(strings []string) []string {
	return qb.EscapeAll(d, strings[0:])
}

// SetEscaping sets the escaping parameter of dialect
func (d *Dialect) SetEscaping(escaping bool) {
	d.escaping = escaping
}

// Escaping gets the escaping parameter of dialect
func (d *Dialect) Escaping() bool {
	return d.escaping
}

// AutoIncrement generates auto increment sql of current dialect
func (d *Dialect) AutoIncrement(column *qb.ColumnElem) string {
	if !column.Options.InlinePrimaryKey {
		panic("Sqlite does not support non-primarykey autoincrement columns")
	}
	return "INTEGER PRIMARY KEY"
}

// SupportsUnsigned returns whether driver supports unsigned type mappings or not
func (d *Dialect) SupportsUnsigned() bool { return false }

// Driver returns the current driver of dialect
func (d *Dialect) Driver() string {
	return "sqlite3"
}

// GetCompiler returns a SqliteCompiler
func (d *Dialect) GetCompiler() qb.Compiler {
	return SqliteCompiler{qb.NewSQLCompiler(d)}
}

// WrapError wraps a native error in a qb Error
func (d *Dialect) WrapError(err error) qb.Error {
	qbErr := qb.Error{Orig: err}
	sErr, ok := err.(sqlite3.Error)
	if !ok {
		return qbErr
	}
	switch sErr.Code {
	case sqlite3.ErrInternal,
		sqlite3.ErrNotFound,
		sqlite3.ErrNomem:
		qbErr.Code = qb.ErrInternal
	case sqlite3.ErrError,
		sqlite3.ErrPerm,
		sqlite3.ErrAbort,
		sqlite3.ErrBusy,
		sqlite3.ErrLocked,
		sqlite3.ErrReadonly,
		sqlite3.ErrInterrupt,
		sqlite3.ErrIoErr,
		sqlite3.ErrFull,
		sqlite3.ErrCantOpen,
		sqlite3.ErrProtocol,
		sqlite3.ErrEmpty,
		sqlite3.ErrSchema:
		qbErr.Code = qb.ErrOperational
	case sqlite3.ErrCorrupt:
		qbErr.Code = qb.ErrDatabase
	case sqlite3.ErrTooBig:
		qbErr.Code = qb.ErrData
	case sqlite3.ErrConstraint,
		sqlite3.ErrMismatch:
		qbErr.Code = qb.ErrIntegrity
	case sqlite3.ErrMisuse:
		qbErr.Code = qb.ErrProgramming
	default:
		qbErr.Code = qb.ErrDatabase
	}
	return qbErr
}

// SqliteCompiler is a SQLCompiler specialised for Sqlite
type SqliteCompiler struct {
	qb.SQLCompiler
}

// VisitUpsert generates the following sql: REPLACE INTO ... VALUES ...
func (SqliteCompiler) VisitUpsert(context *qb.CompilerContext, upsert qb.UpsertStmt) string {
	var (
		colNames []string
		values   []string
	)
	for k, v := range upsert.ValuesMap {
		colNames = append(colNames, context.Compiler.VisitLabel(context, k))
		context.Binds = append(context.Binds, v)
		values = append(values, "?")
	}

	sql := fmt.Sprintf(
		"REPLACE INTO %s(%s)\nVALUES(%s)",
		context.Compiler.VisitLabel(context, upsert.Table.Name),
		strings.Join(colNames, ", "),
		strings.Join(values, ", "),
	)

	return sql
}
