package qb

import "fmt"

// PostgresAdapter is a type of adapter that can be used with postgres driver
type PostgresAdapter struct {
	bindingIndex int
	escaping     bool
}

// Escape wraps the string with escape characters of the adapter
func (a *PostgresAdapter) Escape(str string) string {
	if a.escaping {
		return fmt.Sprintf("\"%s\"", str)
	}
	return str
}

// EscapeAll wraps all elements of string array
func (a *PostgresAdapter) EscapeAll(strings []string) []string {
	return escapeAll(a, strings[0:])
}

// SetEscaping sets the escaping parameter of adapter
func (a *PostgresAdapter) SetEscaping(escaping bool) {
	a.escaping = escaping
}

// GetEscaping gets the escaping parameter of adapter
func (a *PostgresAdapter) GetEscaping() bool {
	return a.escaping
}

// Placeholder returns the placeholder for bindings in the sql
func (a *PostgresAdapter) Placeholder() string {
	a.bindingIndex++
	return fmt.Sprintf("$%d", a.bindingIndex)
}

// Placeholders returns the placeholders for bindings in the sql
func (a *PostgresAdapter) Placeholders(values ...interface{}) []string {
	return placeholders(a, values...)
}

// AutoIncrement generates auto increment sql of current adapter
func (a *PostgresAdapter) AutoIncrement() string {
	return ""
}

// Reset clears the binding index for postgres driver
func (a *PostgresAdapter) Reset() { a.bindingIndex = 0 }

// SupportsInlinePrimaryKey returns whether the driver supports inline primary key definitions
func (a *PostgresAdapter) SupportsInlinePrimaryKey() bool { return true }

// Driver returns the current driver of adapter
func (a *PostgresAdapter) Driver() string {
	return "postgres"
}
