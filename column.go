package qbit
import (
	"fmt"
	"strings"
)

// function generates a table column
func NewColumn(name string, t *Type, constraints []Constraint) Column {
	return Column{
		Name:        name,
		T:           t,
		Constraints: constraints,
	}
}

type Column struct {
	Name        string
	T           *Type
	Constraints []Constraint
}

func (c *Column) Sql() string {
	constraints := []string{}
	for _, v := range c.Constraints {
		constraints = append(constraints, v.Name)
	}
	return fmt.Sprintf("%s %s %s", c.Name, c.T.Sql(), strings.Join(constraints, " "))
}
