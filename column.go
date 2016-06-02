package qb

import "fmt"

// Column generates a ColumnElem given name and type
func Column(name string, t TypeElem) ColumnElem {
	return ColumnElem{name, t}
}

// ColumnElem is the definition of any columns defined in a table
type ColumnElem struct {
	Name string
	Type TypeElem
}

// String returns the column element as an sql clause
func (c ColumnElem) String(adapter Adapter) string {
	return fmt.Sprintf("%s %s", adapter.Escape(c.Name), c.Type.String())
}
