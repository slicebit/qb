package qbit

// NewTable generates a new table pointer given table name, column and table constraints
func NewTable(name string, columns []Column, constraints []Constraint) *Table {
	return &Table{
		name:        name,
		columns:     columns,
		constraints: constraints,
		builder:     NewBuilder(),
	}
}

// Table is the base abstraction for any sql table
type Table struct {
	name        string
	columns     []Column
	constraints []Constraint
	builder     *Builder
}

// SQL generates create table syntax of table
func (t *Table) SQL() string {

	cols := []string{}
	for _, v := range t.columns {
		cols = append(cols, v.SQL())
	}

	constraints := []string{}
	for _, v := range t.constraints {
		constraints = append(constraints, v.Name)
	}

	sql, _, _ := t.builder.CreateTable(t.name, cols, constraints).Build()
	return sql
}

// AddColumn appends a new column to current table
func (t *Table) AddColumn(column Column) {
	t.columns = append(t.columns, column)
}

// AddConstraint appends a new constraint to current table
func (t *Table) AddConstraint(c Constraint) {
	t.constraints = append(t.constraints, c)
}

// Constraints returns the constraint slice of current table
func (t *Table) Constraints() []Constraint {
	return t.constraints
}
