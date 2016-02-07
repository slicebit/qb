package qbit

import (
	"fmt"
	"strings"
)

type Constraint struct {
	Name string
}

func Null() Constraint {
	return Constraint{"NULL"}
}

func NotNull() Constraint {
	return Constraint{"NOT NULL"}
}

func Default(value interface{}) Constraint {
	return Constraint{fmt.Sprintf("DEFAULT `%v`", value)}
}

func Unique(cols ...string) Constraint {
	if len(cols) == 0 {
		return Constraint{"UNIQUE"}
	}
	return Constraint{fmt.Sprintf("UNIQUE(%s)", strings.Join(cols, ", "))}
}

func Key() Constraint {
	return Constraint{"KEY"}
}

func PrimaryKey(cols ...string) Constraint {
	if len(cols) == 0 {
		return Constraint{"PRIMARY KEY"}
	}
	return Constraint{fmt.Sprintf("PRIMARY KEY(%s)", strings.Join(cols, ", "))}
}

func ForeignKey(cols string, table string, refcols string) Constraint {
	return Constraint{fmt.Sprintf(
		"FOREIGN KEY (%s) REFERENCES %s(%s)",
		cols,
		table,
		refcols,
	)}
}
