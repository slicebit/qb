package qb

import (
	"fmt"
	"strings"
)

// PostgresDialect is a type of dialect that can be used with postgres driver
type PostgresDialect struct {
	bindingIndex int
	escaping     bool
}

// NewPostgresDialect returns a new PostgresDialect
func NewPostgresDialect() Dialect {
	return &PostgresDialect{escaping: false, bindingIndex: 0}
}

func init() {
	RegisterDialect("postgres", NewPostgresDialect)
}

// CompileType compiles a type into its DDL
func (d *PostgresDialect) CompileType(t TypeElem) string {
	return DefaultCompileType(t, d.SupportsUnsigned())
}

// Escape wraps the string with escape characters of the dialect
func (d *PostgresDialect) Escape(str string) string {
	if d.escaping {
		return fmt.Sprintf("\"%s\"", str)
	}
	return str
}

// EscapeAll wraps all elements of string array
func (d *PostgresDialect) EscapeAll(strings []string) []string {
	return escapeAll(d, strings[0:])
}

// SetEscaping sets the escaping parameter of dialect
func (d *PostgresDialect) SetEscaping(escaping bool) {
	d.escaping = escaping
}

// Escaping gets the escaping parameter of dialect
func (d *PostgresDialect) Escaping() bool {
	return d.escaping
}

// Placeholder returns the placeholder for bindings in the sql
func (d *PostgresDialect) Placeholder() string {
	d.bindingIndex++
	return fmt.Sprintf("$%d", d.bindingIndex)
}

// Placeholders returns the placeholders for bindings in the sql
func (d *PostgresDialect) Placeholders(values ...interface{}) []string {
	return placeholders(d, values...)
}

// AutoIncrement generates auto increment sql of current dialect
func (d *PostgresDialect) AutoIncrement(column *ColumnElem) string {
	var colSpec string
	if column.Type.Name == "BIGINT" {
		colSpec = "BIGSERIAL"
	} else if column.Type.Name == "SMALLINT" {
		colSpec = "SMALLSERIAL"
	} else {
		colSpec = "SERIAL"
	}
	if column.Options.PrimaryKey {
		colSpec += " PRIMARY KEY"
	}
	return colSpec
}

// Reset clears the binding index for postgres driver
func (d *PostgresDialect) Reset() { d.bindingIndex = 0 }

// SupportsUnsigned returns whether driver supports unsigned type mappings or not
func (d *PostgresDialect) SupportsUnsigned() bool { return false }

// Driver returns the current driver of dialect
func (d *PostgresDialect) Driver() string {
	return "postgres"
}

func (d *PostgresDialect) GetCompiler() Compiler {
	return PostgresCompiler{SQLCompiler{d}}
}

type PostgresCompiler struct {
	SQLCompiler
}

// VisitUpsert generates INSERT INTO ... VALUES ... ON CONFLICT(...) DO UPDATE SET ...
func (PostgresCompiler) VisitUpsert(context *CompilerContext, upsert UpsertStmt) string {
	var (
		colNames []string
		values   []string
	)
	for k, v := range upsert.values {
		colNames = append(colNames, context.Compiler.VisitLabel(context, k))
		context.Binds = append(context.Binds, v)
		values = append(values, context.Dialect.Placeholder())
	}

	var updates []string
	for k, v := range upsert.values {
		updates = append(updates, fmt.Sprintf(
			"%s = %s",
			context.Dialect.Escape(k),
			context.Dialect.Placeholder()))
		context.Binds = append(context.Binds, v)
	}

	var uniqueCols []string
	for _, c := range upsert.table.PrimaryCols() {
		uniqueCols = append(uniqueCols, context.Compiler.VisitLabel(context, c.Name))
	}

	sql := fmt.Sprintf(
		"INSERT INTO %s(%s)\nVALUES(%s)\nON CONFLICT (%s) DO UPDATE SET %s",
		context.Compiler.VisitLabel(context, upsert.table.Name),
		strings.Join(colNames, ", "),
		strings.Join(values, ", "),
		strings.Join(uniqueCols, ", "),
		strings.Join(updates, ", "))

	var returning []string
	for _, r := range upsert.returning {
		returning = append(returning, context.Compiler.VisitLabel(context, r.Name))
	}
	if len(upsert.returning) > 0 {
		sql += fmt.Sprintf(
			"RETURNING %s",
			strings.Join(returning, ", "),
		)
	}
	return sql
}
