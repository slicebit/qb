package qb

import (
	"fmt"
	"strings"
)

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
	return Constraint{fmt.Sprintf("DEFAULT '%v'", value)}
}

// Unique generates generic unique constraint
// if cols are given, then composite unique constraint will be built
func Unique(cols ...string) Constraint {
	if len(cols) == 0 {
		return Constraint{"UNIQUE"}
	}
	return Constraint{fmt.Sprintf("UNIQUE(%s)", strings.Join(cols, ", "))}
}
