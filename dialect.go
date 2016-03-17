package qb

import "fmt"

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

type Dialect interface {
	Escape(str string) string
	Placeholder() string
	Reset()
}

type DefaultDialect struct{}

func (d *DefaultDialect) Escape(str string) string {
	return str
}

func (d *DefaultDialect) Placeholder() string {
	return "?"
}

func (d *DefaultDialect) Reset() {}

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

type MysqlDialect struct{}

func (d *MysqlDialect) Escape(str string) string {
	return fmt.Sprintf("`%s`", str)
}

func (d *MysqlDialect) Placeholder() string {
	return "?"
}

func (d *MysqlDialect) Reset() {}

type SqliteDialect struct{}

func (d *SqliteDialect) Escape(str string) string {
	return fmt.Sprintf("`%s`", str)
}

func (d *SqliteDialect) Placeholder() string {
	return "?"
}

func (d *SqliteDialect) Reset() {}