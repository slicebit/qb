package qb

import (
	"fmt"
)

// NewAdapter returns a adapter pointer given driver
func NewAdapter(driver string) Adapter {
	switch driver {
	case "postgres":
		return &PostgresAdapter{bindingIndex: 0}
	case "mysql":
		return &MysqlAdapter{}
	case "sqlite3":
		return &SqliteAdapter{}
	default:
		return &DefaultAdapter{}
	}
}

// Adapter is the common adapter for driver changes
// It is for fixing compatibility issues of different drivers
type Adapter interface {
	Escape(str string) string
	EscapeAll([]string) []string
	Placeholder() string
	Placeholders(values ...interface{}) []string
	Reset()
	SupportsInlinePrimaryKey() bool
	Driver() string
}

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

// PostgresAdapter is a type of adapter that can be used with postgres driver
type PostgresAdapter struct {
	bindingIndex int
}

// Escape wraps the string with escape characters of the adapter
func (a *PostgresAdapter) Escape(str string) string {
	return fmt.Sprintf("\"%s\"", str)
}

// EscapeAll wraps all elements of string array
func (a *PostgresAdapter) EscapeAll(strings []string) []string {
	return escapeAll(a, strings[0:])
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

// Reset clears the binding index for postgres driver
func (a *PostgresAdapter) Reset() { a.bindingIndex = 0 }

// SupportsInlinePrimaryKey returns whether the driver supports inline primary key definitions
func (a *PostgresAdapter) SupportsInlinePrimaryKey() bool { return true }

// Driver returns the current driver of adapter
func (a *PostgresAdapter) Driver() string {
	return "postgres"
}

// MysqlAdapter is a type of adapter that can be used with mysql driver
type MysqlAdapter struct{}

// Escape wraps the string with escape characters of the adapter
func (a *MysqlAdapter) Escape(str string) string {
	return fmt.Sprintf("`%s`", str)
}

// EscapeAll wraps all elements of string array
func (a *MysqlAdapter) EscapeAll(strings []string) []string {
	return escapeAll(a, strings[0:])
}

// Placeholder returns the placeholder for bindings in the sql
func (a *MysqlAdapter) Placeholder() string {
	return "?"
}

// Placeholders returns the placeholders for bindings in the sql
func (a *MysqlAdapter) Placeholders(values ...interface{}) []string {
	return placeholders(a, values...)
}

// Reset does nothing for the default driver
func (a *MysqlAdapter) Reset() {}

// SupportsInlinePrimaryKey returns whether the driver supports inline primary key definitions
func (a *MysqlAdapter) SupportsInlinePrimaryKey() bool { return false }

// Driver returns the current driver of adapter
func (a *MysqlAdapter) Driver() string {
	return "mysql"
}

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

// common escape all
func escapeAll(adapter Adapter, strings []string) []string {
	for k, v := range strings {
		strings[k] = adapter.Escape(v)
	}

	return strings
}

// common placeholders
func placeholders(adapter Adapter, values ...interface{}) []string {
	placeholders := []string{}
	for _ = range values {
		placeholders = append(placeholders, adapter.Placeholder())
	}
	return placeholders
}
