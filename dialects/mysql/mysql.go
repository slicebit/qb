package mysql

import (
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/slicebit/qb"
)

//go:generate go run ./tools/generrors.go
//go:generate gofmt -w errors.go

const ()

// Dialect is a type of dialect that can be used with mysql driver
type Dialect struct {
	escaping bool
}

// NewDialect returns a new MysqlDialect
func NewDialect() qb.Dialect {
	return &Dialect{false}
}

func init() {
	qb.RegisterDialect("mysql", NewDialect)
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
		return fmt.Sprintf("`%s`", str)
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
	colSpec := d.CompileType(column.Type)
	if column.Options.InlinePrimaryKey {
		colSpec += " PRIMARY KEY"
	}
	colSpec += " AUTO_INCREMENT"
	return colSpec
}

// SupportsUnsigned returns whether driver supports unsigned type mappings or not
func (d *Dialect) SupportsUnsigned() bool { return true }

// Driver returns the current driver of dialect
func (d *Dialect) Driver() string {
	return "mysql"
}

// GetCompiler returns a MysqlCompiler
func (d *Dialect) GetCompiler() qb.Compiler {
	return MysqlCompiler{qb.NewSQLCompiler(d)}
}

// WrapError wraps a native error in a qb Error
func (d *Dialect) WrapError(err error) qb.Error {
	qbErr := qb.Error{Orig: err}
	mErr, ok := err.(*mysql.MySQLError)
	if !ok {
		return qbErr
	}
	// Error mapping logic is copied from MySQL-python-1.2.5
	switch mErr.Number {
	case CR_COMMANDS_OUT_OF_SYNC,
		ER_DB_CREATE_EXISTS,
		ER_SYNTAX_ERROR,
		ER_PARSE_ERROR,
		ER_NO_SUCH_TABLE,
		ER_WRONG_DB_NAME,
		ER_WRONG_TABLE_NAME,
		ER_FIELD_SPECIFIED_TWICE,
		ER_INVALID_GROUP_FUNC_USE,
		ER_UNSUPPORTED_EXTENSION,
		ER_TABLE_MUST_HAVE_COLUMNS,
		ER_CANT_DO_THIS_DURING_AN_TRANSACTION:
		qbErr.Code = qb.ErrProgramming
	case WARN_DATA_TRUNCATED,
		ER_WARN_DATA_OUT_OF_RANGE,
		ER_NO_DEFAULT,
		ER_PRIMARY_CANT_HAVE_NULL,
		ER_DATA_TOO_LONG,
		ER_DATETIME_FUNCTION_OVERFLOW:
		qbErr.Code = qb.ErrData
	case ER_DUP_ENTRY,
		ER_DUP_UNIQUE,
		ER_NO_REFERENCED_ROW,
		ER_NO_REFERENCED_ROW_2,
		ER_ROW_IS_REFERENCED,
		ER_ROW_IS_REFERENCED_2,
		ER_CANNOT_ADD_FOREIGN:
		qbErr.Code = qb.ErrIntegrity
	case ER_WARNING_NOT_COMPLETE_ROLLBACK,
		ER_NOT_SUPPORTED_YET,
		ER_FEATURE_DISABLED,
		ER_UNKNOWN_STORAGE_ENGINE:
		qbErr.Code = qb.ErrNotSupported
	default:
		if mErr.Number < 1000 {
			qbErr.Code = qb.ErrInternal
		} else {
			qbErr.Code = qb.ErrOperational
		}
	}
	return qbErr
}

// MysqlCompiler is a SQLCompiler specialised for Mysql
type MysqlCompiler struct {
	qb.SQLCompiler
}

// VisitUpsert generates INSERT INTO ... VALUES ... ON DUPLICATE KEY UPDATE ...
func (MysqlCompiler) VisitUpsert(context *qb.CompilerContext, upsert qb.UpsertStmt) string {
	var (
		colNames []string
		values   []string
	)

	for k, v := range upsert.ValuesMap {
		colNames = append(colNames, context.Compiler.VisitLabel(context, k))
		context.Binds = append(context.Binds, v)
		values = append(values, "?")
	}

	updates := []string{}
	for k, v := range upsert.ValuesMap {
		updates = append(updates, fmt.Sprintf(
			"%s = %s",
			context.Dialect.Escape(k),
			"?",
		))
		context.Binds = append(context.Binds, v)
	}

	sql := fmt.Sprintf(
		"INSERT INTO %s(%s)\nVALUES(%s)\nON DUPLICATE KEY UPDATE %s",
		context.Dialect.Escape(upsert.Table.Name),
		strings.Join(colNames, ", "),
		strings.Join(values, ", "),
		strings.Join(updates, ", "),
	)

	return sql
}
