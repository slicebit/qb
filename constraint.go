package qb

import (
	"fmt"
	"strings"
)

// Null generates generic null constraint
func Null() ConstraintElem {
	return ConstraintElem{"NULL"}
}

// NotNull generates generic not null constraint
func NotNull() ConstraintElem {
	return ConstraintElem{"NOT NULL"}
}

// Default generates generic default constraint
func Default(value interface{}) ConstraintElem {
	return ConstraintElem{fmt.Sprintf("DEFAULT '%v'", value)}
}

// Unique generates generic unique constraint
// if cols are given, then composite unique constraint will be built
func Unique() ConstraintElem {
	return ConstraintElem{"UNIQUE"}
}

// Constraint generates a custom constraint due to variation of adapters
func Constraint(name string) ConstraintElem {
	return ConstraintElem{name}
}

// ConstraintElem is the definition of column & table constraints
type ConstraintElem struct {
	Name string
}

// String returns the constraint as an sql clause
func (c ConstraintElem) String() string {
	return c.Name
}

// PrimaryKey generates a primary key constraint of any table
func PrimaryKey(cols ...string) PrimaryKeyConstraint {
	return PrimaryKeyConstraint{cols}
}

// PrimaryKeyConstraint is the definition of primary key constraints of any table
type PrimaryKeyConstraint struct {
	Columns []string
}

// String returns the primary key constraints as an sql clause
func (c PrimaryKeyConstraint) String(adapter Adapter) string {
	cols := []string{}
	for _, col := range c.Columns {
		cols = append(cols, adapter.Escape(col))
	}

	return fmt.Sprintf("PRIMARY KEY(%s)", strings.Join(cols, ", "))
}

// ForeignKey generates a foreign key for table constraint definitions
func ForeignKey() ForeignKeyConstraints {
	return ForeignKeyConstraints{[]Reference{}}
}

// ForeignKeyConstraint is the definition of foreign key in any table
type ForeignKeyConstraints struct {
	Refs []Reference
}

func (c ForeignKeyConstraints) String(adapter Adapter) string {
	clauses := []string{}
	for _, ref := range c.Refs {
		clauses = append(clauses, fmt.Sprintf(
			"\tFOREIGN KEY(%s) REFERENCES %s(%s)",
			strings.Join(adapter.EscapeAll(ref.Cols), ", "),
			adapter.Escape(ref.RefTable),
			strings.Join(adapter.EscapeAll(ref.RefCols), ", "),
		))
	}

	return strings.Join(clauses, ",\n")
}

// Ref generates a reference after the definition of foreign key by chaining
func (c ForeignKeyConstraints) Ref(col string, refTable string, refCol string) ForeignKeyConstraints {
	for k, v := range c.Refs {
		if refTable == v.RefTable {
			c.Refs[k].Cols = append(c.Refs[k].Cols, col)
			c.Refs[k].RefCols = append(c.Refs[k].RefCols, refCol)
			return c
		}
	}

	ref := Reference{[]string{}, refTable, []string{}}
	ref.Cols = append(ref.Cols, col)
	ref.RefCols = append(ref.RefCols, refCol)
	c.Refs = append(c.Refs, ref)
	return c
}

// Reference is the main struct for defining foreign key references
type Reference struct {
	Cols     []string
	RefTable string
	RefCols  []string
}

// UniqueKey generates UniqueKeyConstraint given columns as strings
func UniqueKey(cols ...string) UniqueKeyConstraint {
	return UniqueKeyConstraint{
		fmt.Sprintf("u_%s", strings.Join(cols, "_")),
		cols,
	}
}

// UniqueKeyConstraint is the base struct to define composite unique indexes of tables
type UniqueKeyConstraint struct {
	name string
	cols []string
}

// String generates composite unique indices as sql clause
func (c UniqueKeyConstraint) String(adapter Adapter) string {
	return fmt.Sprintf("CONSTRAINT %s UNIQUE(%s)", c.name, strings.Join(c.cols, ", "))
}
