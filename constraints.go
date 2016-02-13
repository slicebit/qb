package qbit

import (
	"fmt"
	"strings"
)

// PrimaryKey is the builtin constraints that should be used with tag
// type User struct {
//		Id int
//		PrimaryKey `qbit:"id"`
// }
type PrimaryKey Constraint

// ForeignKey is the builtin constraint that should be used with tag
type ForeignKey Constraint

// CompositeUnique is the builtin multiple unique constraint that should be used with tag
type CompositeUnique Constraint

// Constraint is the generic struct for table level and column level constraints
type Constraint struct {
	Name string
}

// Null generates generic null constraint
func Null() Constraint {
	return Constraint{"NULL"}
}

// NotNull generates generic not null constraint
func NotNull() Constraint {
	return Constraint{"NOT NULL"}
}

// Default generates generic default constraint
func Default(value interface{}) Constraint {
	return Constraint{fmt.Sprintf("DEFAULT '%s'", value)}
}

// Unique generates generic unique constraint
// if cols are given, then composite unique constraint will be built
func Unique(cols ...string) Constraint {
	if len(cols) == 0 {
		return Constraint{"UNIQUE"}
	}
	return Constraint{fmt.Sprintf("UNIQUE(%s)", strings.Join(cols, ", "))}
}

// Primary generates generic primary key syntax
// if cols are given, then composite primary key will be built
func Primary(cols ...string) Constraint {
	if len(cols) == 0 {
		return Constraint{"PRIMARY KEY"}
	}
	constraint := Constraint{fmt.Sprintf("PRIMARY KEY(%s)", strings.Join(cols, ", "))}
	return constraint
}

// Foreign generates generic foreign key syntax
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
