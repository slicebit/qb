package qb

import (
	"fmt"
	"strings"

	"github.com/slicebit/qb"
)

// MysqlDialect is a type of dialect that can be used with mysql driver
type MysqlDialect struct {
	escaping bool
}

// NewMysqlDialect returns a new MysqlDialect
func NewMysqlDialect() qb.Dialect {
	return &MysqlDialect{false}
}

func init() {
	qb.RegisterDialect("mysql", NewMysqlDialect)
}

// CompileType compiles a type into its DDL
func (d *MysqlDialect) CompileType(t qb.TypeElem) string {
	if t.Name == "UUID" {
		return "VARCHAR(36)"
	}
	return qb.DefaultCompileType(t, d.SupportsUnsigned())
}

// Escape wraps the string with escape characters of the dialect
func (d *MysqlDialect) Escape(str string) string {
	if d.escaping {
		return fmt.Sprintf("`%s`", str)
	}
	return str
}

// EscapeAll wraps all elements of string array
func (d *MysqlDialect) EscapeAll(strings []string) []string {
	return qb.EscapeAll(d, strings[0:])
}

// SetEscaping sets the escaping parameter of dialect
func (d *MysqlDialect) SetEscaping(escaping bool) {
	d.escaping = escaping
}

// Escaping gets the escaping parameter of dialect
func (d *MysqlDialect) Escaping() bool {
	return d.escaping
}

// AutoIncrement generates auto increment sql of current dialect
func (d *MysqlDialect) AutoIncrement(column *qb.ColumnElem) string {
	colSpec := d.CompileType(column.Type)
	if column.Options.InlinePrimaryKey {
		colSpec += " PRIMARY KEY"
	}
	colSpec += " AUTO_INCREMENT"
	return colSpec
}

// SupportsUnsigned returns whether driver supports unsigned type mappings or not
func (d *MysqlDialect) SupportsUnsigned() bool { return true }

// Driver returns the current driver of dialect
func (d *MysqlDialect) Driver() string {
	return "mysql"
}

// GetCompiler returns a MysqlCompiler
func (d *MysqlDialect) GetCompiler() qb.Compiler {
	return MysqlCompiler{qb.NewSQLCompiler(d)}
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
