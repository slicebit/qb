package qb

// DefaultAdapter is a type of adapter that can be used with unsupported sql drivers
type DefaultAdapter struct{}

// Escape wraps the string with escape characters of the adapter
func (a *DefaultAdapter) Escape(str string) string {
	return str
}

// EscapeAll wraps all elements of string array
func (a *DefaultAdapter) EscapeAll(strings []string) []string {
	return escapeAll(a, strings[0:])
}

// Placeholder returns the placeholder for bindings in the sql
func (a *DefaultAdapter) Placeholder() string {
	return "?"
}

// Placeholders returns the placeholders for bindings in the sql
func (a *DefaultAdapter) Placeholders(values ...interface{}) []string {
	return placeholders(a, values...)
}

// Reset does nothing for the default driver
func (a *DefaultAdapter) Reset() {}

// SupportsInlinePrimaryKey returns whether the driver supports inline primary key definitions
func (a *DefaultAdapter) SupportsInlinePrimaryKey() bool { return false }

// Driver returns the current driver of adapter
func (a *DefaultAdapter) Driver() string {
	return ""
}
