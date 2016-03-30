package qb

import (
	"fmt"
	"strings"
)

// NewColumn creates and returns a table column
func NewColumn(name string, t *Type, constraints []Constraint) Column {
	return Column{
		Name:        name,
		Type:        t,
		Constraints: constraints,
	}
}

// Column is the base abstraction for table struct columns
type Column struct {
	Name        string
	Type        *Type
	Constraints []Constraint
}

// SQL returns column as an sql statement
func (c *Column) SQL(driver string) string {

	adapter := NewAdapter(driver)

	constraints := []string{}
	for _, v := range c.Constraints {
		constraints = append(constraints, v.Name)
	}

	colPieces := []string{
		fmt.Sprintf("%s", adapter.Escape(c.Name)),
		c.Type.SQL,
	}

	if len(constraints) > 0 {
		colPieces = append(colPieces, strings.Join(constraints, " "))
	}

	return fmt.Sprintf("%s", strings.Join(colPieces, " "))
}
