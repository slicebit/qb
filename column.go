package qbit

import (
	"fmt"
	"strings"
)

// function generates a table column
func NewColumn(name string, t *Type, constraints []Constraint) Column {
	return Column{
		Name:        name,
		Type:        t,
		Constraints: constraints,
	}
}

type Column struct {
	Name        string
	Type        *Type
	Constraints []Constraint
}

func (c *Column) Sql() string {

	constraints := []string{}
	for _, v := range c.Constraints {
		constraints = append(constraints, v.Name)
	}

	colPieces := []string{
		c.Name,
		c.Type.Sql(),
	}

	if len(constraints) > 0 {
		colPieces = append(colPieces, strings.Join(constraints, " "))
	}

	return fmt.Sprintf("%s", strings.Join(colPieces, " "))
}
