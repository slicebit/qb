package qb

import "fmt"

// DefaultAdapter is a type of adapter that can be used with unsupported sql drivers
type DefaultAdapter struct {
	escaping bool
}

// Escape wraps the string with escape characters of the adapter
func (a *DefaultAdapter) Escape(str string) string {
	if a.escaping {
		return fmt.Sprintf("`%s`", str)
	}
	return str
}

// EscapeAll wraps all elements of string array
func (a *DefaultAdapter) EscapeAll(strings []string) []string {
	return escapeAll(a, strings[0:])
}

// SetEscaping sets the escaping parameter of adapter
func (a *DefaultAdapter) SetEscaping(escaping bool) {
	a.escaping = escaping
}

// GetEscaping gets the escaping parameter of adapter
func (a *DefaultAdapter) GetEscaping() bool {
	return a.escaping
}

// Placeholder returns the placeholder for bindings in the sql
func (a *DefaultAdapter) Placeholder() string {
	return "?"
}

// Placeholders returns the placeholders for bindings in the sql
func (a *DefaultAdapter) Placeholders(values ...interface{}) []string {
	return placeholders(a, values...)
}

// AutoIncrement generates auto increment sql of current adapter
func (a *DefaultAdapter) AutoIncrement() string {
	return "AUTO INCREMENT"
}

// Reset does nothing for the default driver
func (a *DefaultAdapter) Reset() {}

// SupportsInlinePrimaryKey returns whether the driver supports inline primary key definitions
func (a *DefaultAdapter) SupportsInlinePrimaryKey() bool { return false }

// SupportsUnsigned returns whether driver supports unsigned type mappings or not
func (a *DefaultAdapter) SupportsUnsigned() bool { return false }

// Driver returns the current driver of adapter
func (a *DefaultAdapter) Driver() string {
	return ""
}
