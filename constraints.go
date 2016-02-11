package qbit

import (
	"fmt"
	"strings"
)

// builtin constraints that should be used with tag
// e.g;
// type User struct {
//		Id int
//		PrimaryKey `qbit:"id"`
// }
type PrimaryKey Constraint
type ForeignKey Constraint
type CompositeUnique Constraint

type Constraint struct {
	Name string
}

// function generates generic null constraint
func Null() Constraint {
	return Constraint{"NULL"}
}

// function generates generic not null constraint
func NotNull() Constraint {
	return Constraint{"NOT NULL"}
}

// function generates generic default constraint
func Default(value interface{}) Constraint {
	return Constraint{fmt.Sprintf("DEFAULT `%v`", value)}
}

// function generates generic unique constraint
// if cols are givern, then composite unique constraint will be built
func Unique(cols ...string) Constraint {
	if len(cols) == 0 {
		return Constraint{"UNIQUE"}
	}
	return Constraint{fmt.Sprintf("UNIQUE(%s)", strings.Join(cols, ", "))}
}

// function generates generic primary key syntax
// if cols are given, then composite primary key will be built
func Primary(cols ...string) Constraint {
	if len(cols) == 0 {
		return Constraint{"PRIMARY KEY"}
	}
	constraint := Constraint{fmt.Sprintf("PRIMARY KEY(%s)", strings.Join(cols, ", "))}
	return constraint
}

// function generates generic foreign key syntax
func Foreign(cols string, reftable string, refcols string) Constraint {
	constraint := Constraint{
		fmt.Sprintf(
			"FOREIGN KEY (%s) REFERENCES %s(%s)",
			cols,
			reftable,
			refcols,
		),
	}
	return constraint
}
