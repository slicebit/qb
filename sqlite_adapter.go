package qb

import "fmt"

// SqliteAdapter is a type of adapter that can be used with sqlite driver
type SqliteAdapter struct{}

// Escape wraps the string with escape characters of the adapter
func (a *SqliteAdapter) Escape(str string) string {
	return fmt.Sprintf("`%s`", str)
}

// EscapeAll wraps all elements of string array
func (a *SqliteAdapter) EscapeAll(strings []string) []string {
	return escapeAll(a, strings[0:])
}

// Placeholder returns the placeholder for bindings in the sql
func (a *SqliteAdapter) Placeholder() string {
	return "?"
}

// Placeholders returns the placeholders for bindings in the sql
func (a *SqliteAdapter) Placeholders(values ...interface{}) []string {
	return placeholders(a, values...)
}

// Reset does nothing for the default driver
func (a *SqliteAdapter) Reset() {}

// SupportsInlinePrimaryKey returns whether the driver supports inline primary key definitions
func (a *SqliteAdapter) SupportsInlinePrimaryKey() bool { return true }

// Driver returns the current driver of adapter
func (a *SqliteAdapter) Driver() string {
	return "sqlite3"
}
