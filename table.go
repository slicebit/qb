package qbit

func NewTable(name string, columns []Column, constraints []Constraint) *Table {
	return &Table{
		name:        name,
		columns:     columns,
		constraints: constraints,
		builder:     NewBuilder(),
	}
}

type Table struct {
	name        string
	columns     []Column
	constraints []Constraint
	builder     *Builder
}

func (t *Table) Sql() string {

	cols := []string{}
	for _, v := range t.columns {
		cols = append(cols, v.Sql())
	}

	constraints := []string{}
	for _, v := range t.constraints {
		constraints = append(constraints, v.Name)
	}

	sql, _, _ := t.builder.CreateTable(t.name, cols, constraints).Build()
	return sql
}

func (t *Table) AddConstraint(c Constraint) {
	t.constraints = append(t.constraints, c)
}

func (t *Table) Constraints() []Constraint {
	return t.constraints
}
