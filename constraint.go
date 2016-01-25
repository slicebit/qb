package qbit

import (
	"fmt"
	"strings"
)

type Constraint string

func NotNull() *Constraint {
	return "NOT NULL"
}

func Default(value interface{}) *Constraint {
	return fmt.Sprintf("DEFAULT `%s`", value)
}

func Unique(cols ...string) *Constraint {
	if len(cols) > 0 {
		return fmt.Sprintf("UNIQUE")
	}
	return fmt.Sprintf("UNIQUE(%s)", strings.Join(cols, ", "))
}

func Key() *Constraint {
	return "KEY"
}

func PrimaryKey(cols ...string) *Constraint {
	if len(cols) > 0 {
		return fmt.Sprintf("PRIMARY KEY")
	}
	return fmt.Sprintf("PRIMARY KEY(%s)", strings.Join(cols, ", "))
}

func ForeignKey(cols string, table string, refcols string) *Constraint {
	return fmt.Sprintf(
		"FOREIGN KEY (%s) REFERENCES %s ($s)",
		cols,
		table,
		refcols,
	)
}
