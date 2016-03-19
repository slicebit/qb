package qb

import (
	"fmt"
)

// NewDialect returns a dialect pointer given driver
func NewDialect(driver string) Dialect {
	switch driver {
	case "postgres":
		return &PostgresDialect{bindingIndex: 0}
	case "mysql":
		return &MysqlDialect{}
	case "sqlite3":
		return &SqliteDialect{}
	default:
		return &DefaultDialect{}
	}
}

// Dialect is the common adapter for driver changes
// It is for fixing compatibility issues of different drivers
type Dialect interface {
	Escape(str string) string
	Placeholder() string
	Placeholders(values ...interface{}) []string
	Reset()
	SupportsInlinePrimaryKey() bool
}

// DefaultDialect is a type of dialect that can be used with unsupported sql drivers
type DefaultDialect struct{}

// Escape wraps the string with escape characters of the dialect
func (d *DefaultDialect) Escape(str string) string {
	return str
}

// Placeholder returns the placeholder for bindings in the sql
func (d *DefaultDialect) Placeholder() string {
	return "?"
}

// Placeholders returns the placeholders for bindings in the sql
func (d *DefaultDialect) Placeholders(values ...interface{}) []string {
	placeholders := []string{}
	for _ = range values {
		placeholders = append(placeholders, d.Placeholder())
	}
	return placeholders
}

// Reset does nothing for the default driver
func (d *DefaultDialect) Reset() {}

// SupportsInlinePrimaryKey returns whether the driver supports inline primary key definitions
func (d *DefaultDialect) SupportsInlinePrimaryKey() bool { return false }

// PostgresDialect is a type of dialect that can be used with postgres driver
type PostgresDialect struct {
	bindingIndex int
}

// Escape wraps the string with escape characters of the dialect
func (d *PostgresDialect) Escape(str string) string {
	return fmt.Sprintf("\"%s\"", str)
}

// Placeholder returns the placeholder for bindings in the sql
func (d *PostgresDialect) Placeholder() string {
	d.bindingIndex++
	return fmt.Sprintf("$%d", d.bindingIndex)
}

// Placeholders returns the placeholders for bindings in the sql
func (d *PostgresDialect) Placeholders(values ...interface{}) []string {
	placeholders := []string{}
	for _ = range values {
		placeholders = append(placeholders, d.Placeholder())
	}
	return placeholders
}

// Reset clears the binding index for postgres driver
func (d *PostgresDialect) Reset() { d.bindingIndex = 0 }

// SupportsInlinePrimaryKey returns whether the driver supports inline primary key definitions
func (d *PostgresDialect) SupportsInlinePrimaryKey() bool { return true }

// MysqlDialect is a type of dialect that can be used with mysql driver
type MysqlDialect struct{}

// Escape wraps the string with escape characters of the dialect
func (d *MysqlDialect) Escape(str string) string {
	return fmt.Sprintf("`%s`", str)
}

// Placeholder returns the placeholder for bindings in the sql
func (d *MysqlDialect) Placeholder() string {
	return "?"
}

// Placeholders returns the placeholders for bindings in the sql
func (d *MysqlDialect) Placeholders(values ...interface{}) []string {
	placeholders := []string{}
	for _ = range values {
		placeholders = append(placeholders, d.Placeholder())
	}
	return placeholders
}

// Reset does nothing for the default driver
func (d *MysqlDialect) Reset() {}

// SupportsInlinePrimaryKey returns whether the driver supports inline primary key definitions
func (d *MysqlDialect) SupportsInlinePrimaryKey() bool { return false }

// SqliteDialect is a type of dialect that can be used with sqlite driver
type SqliteDialect struct{}

// Escape wraps the string with escape characters of the dialect
func (d *SqliteDialect) Escape(str string) string {
	return fmt.Sprintf("`%s`", str)
}

// Placeholder returns the placeholder for bindings in the sql
func (d *SqliteDialect) Placeholder() string {
	return "?"
}

// Placeholders returns the placeholders for bindings in the sql
func (d *SqliteDialect) Placeholders(values ...interface{}) []string {
	placeholders := []string{}
	for _ = range values {
		placeholders = append(placeholders, d.Placeholder())
	}
	return placeholders
}

// Reset does nothing for the default driver
func (d *SqliteDialect) Reset() {}

// SupportsInlinePrimaryKey returns whether the driver supports inline primary key definitions
func (d *SqliteDialect) SupportsInlinePrimaryKey() bool { return true }
