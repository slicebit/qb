package qb

import (
	"fmt"
	"strings"
)

// NewTable generates a new table pointer given table name, column and table constraints
func NewTable(name string, columns []Column, constraints []Constraint) *Table {
	return &Table{
		name:        name,
		columns:     columns,
		constraints: constraints,
		builder:     NewBuilder(),
		primaryCols: []string{},
		refs:        []ref{},
	}
}

// Table is the base abstraction for any sql table
type Table struct {
	name        string
	columns     []Column
	constraints []Constraint
	builder     *Builder
	primaryCols []string
	refs        []ref
}

// SQL generates create table syntax of table
func (t *Table) SQL() string {

	cols := []string{}
	for _, v := range t.columns {
		cols = append(cols, v.SQL())
	}

	constraints := []string{}

	// build primary key constraints using primaryCols
	if len(t.primaryCols) > 0 {
		constraints = append(constraints, fmt.Sprintf("PRIMARY KEY(%s)", strings.Join(t.primaryCols, ", ")))
	}

	// build foreign key constraints using refCols
	for _, ref := range t.refs {
		constraints = append(constraints, fmt.Sprintf("FOREIGN KEY(%s) REFERENCES %s(%s)", strings.Join(ref.cols, ", "), ref.refTable, strings.Join(ref.refCols, ", ")))
	}

	for _, v := range t.constraints {
		constraints = append(constraints, v.Name)
	}

	query, err := t.builder.CreateTable(t.name, cols, constraints).Build()
	if err != nil {
		return ""
	}

	return query.SQL()
}

// AddColumn appends a new column to current table
func (t *Table) AddColumn(column Column) {
	t.columns = append(t.columns, column)
}

// AddConstraint appends a new constraint to current table
func (t *Table) AddConstraint(c Constraint) {
	t.constraints = append(t.constraints, c)
}

// AddPrimary appends a primary column that will be lazily built as a primary key constraint
func (t *Table) AddPrimary(col string) {
	t.primaryCols = append(t.primaryCols, col)
}

type ref struct {
	cols     []string
	refTable string
	refCols  []string
}

// AddRef appends a new reference struct that will be lazily built as a foreign key constraint
func (t *Table) AddRef(col string, refTable string, refCol string) {

	if len(t.refs) > 0 {
		for k, ref := range t.refs {
			if refTable == ref.refTable {
				t.refs[k].cols = append(t.refs[k].cols, fmt.Sprintf("`%s`", col))
				t.refs[k].refCols = append(t.refs[k].refCols, fmt.Sprintf("`%s`", refCol))
				return
			}
		}
	}

	r := ref{[]string{}, refTable, []string{}}
	r.cols = append(r.cols, col)
	r.refCols = append(r.refCols, refCol)
	t.refs = append(t.refs, r)
}

// Constraints returns the constraint slice of current table
func (t *Table) Constraints() []Constraint {
	return t.constraints
}

func (t *Table) Insert(kv map[string]interface{}) (*Query, error) {

	keys := make([]string, len(kv))
	vals := make([]interface{}, len(kv))

	for k, v := range kv {
		keys = append(keys, k)
		vals = append(vals, v)
	}

	return t.builder.Insert(t.name, keys...).Values(vals...).Build()
}