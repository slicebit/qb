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
	Reset()
	SupportsInlinePrimaryKey() bool
}

// Default dialect is a type of dialect that can be used with unsupported sql drivers
type DefaultDialect struct{}

// Escape wraps the string with escape characters of the dialect
func (d *DefaultDialect) Escape(str string) string {
	return str
}

func (d *DefaultDialect) Placeholder() string {
	return "?"
}

func (d *DefaultDialect) Reset() {}

func (d *DefaultDialect) SupportsInlinePrimaryKey() bool { return false }

// Default dialect is a type of dialect that can be used with postgres sql driver
type PostgresDialect struct {
	bindingIndex int
}

func (d *PostgresDialect) Escape(str string) string {
	return fmt.Sprintf("\"%s\"", str)
}

func (d *PostgresDialect) Placeholder() string {
	d.bindingIndex++
	return fmt.Sprintf("$%d", d.bindingIndex)
}

func (d *PostgresDialect) Reset() { d.bindingIndex = 0 }

func (d *PostgresDialect) SupportsInlinePrimaryKey() bool { return true }

type MysqlDialect struct{}

func (d *MysqlDialect) Escape(str string) string {
	return fmt.Sprintf("`%s`", str)
}

func (d *MysqlDialect) Placeholder() string {
	return "?"
}

func (d *MysqlDialect) Reset() {}

func (d *MysqlDialect) SupportsInlinePrimaryKey() bool { return false }

type SqliteDialect struct{}

func (d *SqliteDialect) Escape(str string) string {
	return fmt.Sprintf("`%s`", str)
}

func (d *SqliteDialect) Placeholder() string {
	return "?"
}

func (d *SqliteDialect) Reset() {}

func (d *SqliteDialect) SupportsInlinePrimaryKey() bool { return true }
