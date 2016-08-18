package qb

import "fmt"

// PostgresDialect is a type of dialect that can be used with postgres driver
type PostgresDialect struct {
	bindingIndex int
	escaping     bool
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
