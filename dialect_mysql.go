package qb

import "fmt"

// MysqlDialect is a type of dialect that can be used with mysql driver
type MysqlDialect struct {
	escaping bool
}

// NewMysqlDialect returns a new MysqlDialect
func NewMysqlDialect() Dialect {
	return &MysqlDialect{false}
}

func init() {
	RegisterDialect("mysql", NewMysqlDialect)
}

// CompileType compiles a type into its DDL
func (d *MysqlDialect) CompileType(t TypeElem) string {
	if t.Name == "UUID" {
		return "VARCHAR(36)"
	}
	return DefaultCompileType(t, d.SupportsUnsigned())
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
	return escapeAll(d, strings[0:])
}

// SetEscaping sets the escaping parameter of dialect
func (d *MysqlDialect) SetEscaping(escaping bool) {
	d.escaping = escaping
}

// Escaping gets the escaping parameter of dialect
func (d *MysqlDialect) Escaping() bool {
	return d.escaping
}

// Placeholder returns the placeholder for bindings in the sql
func (d *MysqlDialect) Placeholder() string {
	return "?"
}

// Placeholders returns the placeholders for bindings in the sql
func (d *MysqlDialect) Placeholders(values ...interface{}) []string {
	return placeholders(d, values...)
}

// AutoIncrement generates auto increment sql of current dialect
func (d *MysqlDialect) AutoIncrement(column *ColumnElem) string {
	colSpec := d.CompileType(column.Type)
	if column.Options.PrimaryKey {
		colSpec += " PRIMARY KEY"
	}
	colSpec += " AUTO_INCREMENT"
	return colSpec
}

// Reset does nothing for the default driver
func (d *MysqlDialect) Reset() {}

// SupportsUnsigned returns whether driver supports unsigned type mappings or not
func (d *MysqlDialect) SupportsUnsigned() bool { return true }

// Driver returns the current driver of dialect
func (d *MysqlDialect) Driver() string {
	return "mysql"
}

func (d *MysqlDialect) GetCompiler() Compiler {
	return SQLCompiler{d}
}
