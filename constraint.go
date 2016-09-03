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

// Constraint generates a custom constraint due to variation of dialects
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
func (c PrimaryKeyConstraint) String(dialect Dialect) string {
	cols := []string{}
	for _, col := range c.Columns {
		cols = append(cols, dialect.Escape(col))
	}

	return fmt.Sprintf("PRIMARY KEY(%s)", strings.Join(cols, ", "))
}

// ForeignKey generates a foreign key for table constraint definitions
func ForeignKey(cols ...string) ForeignKeyConstraint {
	return ForeignKeyConstraint{Cols: cols}
}

// ForeignKeyConstraints is the definition of foreign keys in any table
type ForeignKeyConstraints struct {
	FKeys []ForeignKeyConstraint
}

func (c ForeignKeyConstraints) String(dialect Dialect) string {
	clauses := []string{}
	for _, fkey := range c.FKeys {
		clauses = append(clauses, fkey.String(dialect))
	}

	return strings.Join(clauses, ",\n")
}

// ForeignKeyConstraint is the main struct for defining foreign key references
type ForeignKeyConstraint struct {
	Cols     []string
	RefTable string
	RefCols  []string
}

func (fkey ForeignKeyConstraint) String(dialect Dialect) string {
	return fmt.Sprintf(
		"\tFOREIGN KEY(%s) REFERENCES %s(%s)",
		strings.Join(dialect.EscapeAll(fkey.Cols), ", "),
		dialect.Escape(fkey.RefTable),
		strings.Join(dialect.EscapeAll(fkey.RefCols), ", "),
	)
}

// References set the reference part of the foreign key
func (fkey ForeignKeyConstraint) References(refTable string, refCols ...string) ForeignKeyConstraint {
	fkey.RefTable = refTable
	fkey.RefCols = refCols
	return fkey
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
func (c UniqueKeyConstraint) String(dialect Dialect) string {
	return fmt.Sprintf("CONSTRAINT %s UNIQUE(%s)", c.name, strings.Join(dialect.EscapeAll(c.cols), ", "))
}
