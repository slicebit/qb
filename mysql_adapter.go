package qb

import "fmt"

// MysqlAdapter is a type of adapter that can be used with mysql driver
type MysqlAdapter struct {
	escaping bool
}

// Escape wraps the string with escape characters of the adapter
func (a *MysqlAdapter) Escape(str string) string {
	if a.escaping {
		return fmt.Sprintf("`%s`", str)
	}
	return str
}

// EscapeAll wraps all elements of string array
func (a *MysqlAdapter) EscapeAll(strings []string) []string {
	return escapeAll(a, strings[0:])
}

// SetEscaping sets the escaping parameter of adapter
func (a *MysqlAdapter) SetEscaping(escaping bool) {
	a.escaping = escaping
}

// GetEscaping gets the escaping parameter of adapter
func (a *MysqlAdapter) GetEscaping() bool {
	return a.escaping
}

// Placeholder returns the placeholder for bindings in the sql
func (a *MysqlAdapter) Placeholder() string {
	return "?"
}

// Placeholders returns the placeholders for bindings in the sql
func (a *MysqlAdapter) Placeholders(values ...interface{}) []string {
	return placeholders(a, values...)
}

// AutoIncrement generates auto increment sql of current adapter
func (a *MysqlAdapter) AutoIncrement() string {
	return "AUTO_INCREMENT"
}

// Reset does nothing for the default driver
func (a *MysqlAdapter) Reset() {}

// SupportsInlinePrimaryKey returns whether the driver supports inline primary key definitions
func (a *MysqlAdapter) SupportsInlinePrimaryKey() bool { return false }

// Driver returns the current driver of adapter
func (a *MysqlAdapter) Driver() string {
	return "mysql"
}
