package qbit

func Table(name string, columns []Column, constraints []Constraint) *table {
	return &table{
		name: name,
		columns: columns,
		constraints: constraints,
	}
}

type table struct {
	name string
	columns []Column
	constraints []Constraint
}

func (t *table) String(engine string) string {
	return ""
}