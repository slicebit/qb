package qb

import "fmt"

// DefaultDialect is a type of dialect that can be used with unsupported sql drivers
type DefaultDialect struct {
	escaping bool
}

// CompileType compiles a type into its DDL
func (d *DefaultDialect) CompileType(t TypeElem) string {
	return DefaultCompileType(t, d.SupportsUnsigned())
}

// Escape wraps the string with escape characters of the dialect
func (d *DefaultDialect) Escape(str string) string {
	if d.escaping {
		return fmt.Sprintf("`%s`", str)
	}
	return str
}

// EscapeAll wraps all elements of string array
func (d *DefaultDialect) EscapeAll(strings []string) []string {
	return escapeAll(d, strings[0:])
}

// SetEscaping sets the escaping parameter of dialect
func (d *DefaultDialect) SetEscaping(escaping bool) {
	d.escaping = escaping
}

// Escaping gets the escaping parameter of dialect
func (d *DefaultDialect) Escaping() bool {
	return d.escaping
}

// Placeholder returns the placeholder for bindings in the sql
func (d *DefaultDialect) Placeholder() string {
	return "?"
}

// Placeholders returns the placeholders for bindings in the sql
func (d *DefaultDialect) Placeholders(values ...interface{}) []string {
	return placeholders(d, values...)
}

// AutoIncrement generates auto increment sql of current dialect
func (d *DefaultDialect) AutoIncrement(column *ColumnElem) string {
	colSpec := d.CompileType(column.Type)
	if column.Options.PrimaryKey {
		colSpec += " PRIMARY KEY"
	}
	colSpec += " AUTO INCREMENT"
	return colSpec
}

// Reset does nothing for the default driver
func (d *DefaultDialect) Reset() {}

// SupportsUnsigned returns whether driver supports unsigned type mappings or not
func (d *DefaultDialect) SupportsUnsigned() bool { return false }

// Driver returns the current driver of dialect
func (d *DefaultDialect) Driver() string {
	return ""
}

func (d *DefaultDialect) GetCompiler() Compiler {
	return SQLCompiler{d}
}
