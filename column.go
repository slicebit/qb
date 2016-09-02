package qb

import (
	"fmt"
	"strings"
)

// Column generates a ColumnElem given name and type
func Column(name string, t TypeElem) ColumnElem {
	return ColumnElem{
		Name: name,
		Type: t,
	}
}

// ColumnOptions holds options for a column
type ColumnOptions struct {
	AutoIncrement bool
	PrimaryKey    bool
	Unique        bool
}

// ColumnElem is the definition of any columns defined in a table
type ColumnElem struct {
	Name        string
	Type        TypeElem
	Table       string // This field should be lazily set by Table() function
	Constraints []ConstraintElem
	Options     ColumnOptions
}

// AutoIncrement set up “auto increment” semantics for an integer column.
// Depending on the dialect, the column may be required to be a PrimaryKey too.
func (c ColumnElem) AutoIncrement() ColumnElem {
	c.Options.AutoIncrement = true
	return c
}

// PrimaryKey add the column to the primary key
func (c ColumnElem) PrimaryKey() ColumnElem {
	c.Options.PrimaryKey = true
	return c
}

// String returns the column element as an sql clause
// It satisfies the TableSQLClause interface
func (c ColumnElem) String(dialect Dialect) string {
	colSpec := ""
	if c.Options.AutoIncrement {
		colSpec = dialect.AutoIncrement(&c)
	}
	if colSpec == "" {
		colSpec = dialect.CompileType(c.Type)
		constraintNames := []string{}
		for _, constraint := range c.Constraints {
			constraintNames = append(constraintNames, constraint.String())
		}
		if len(constraintNames) != 0 {
			colSpec = fmt.Sprintf("%s %s", colSpec, strings.Join(constraintNames, " "))
		}
	}
	res := fmt.Sprintf("%s %s", dialect.Escape(c.Name), colSpec)
	return res
}

// Build compiles the column element and returns sql, bindings
// It satisfies the SQLClause interface
func (c ColumnElem) Build(dialect Dialect) (string, []interface{}) {
	return dialect.Escape(c.Name), []interface{}{}
}

// constraints setters

// Default adds a default constraint to column type
func (c ColumnElem) Default(def interface{}) ColumnElem {
	c.Constraints = append(c.Constraints, Default(def))
	return c
}

// Null adds null constraint to column type
func (c ColumnElem) Null() ColumnElem {
	c.Constraints = append(c.Constraints, Null())
	return c
}

// NotNull adds not null constraint to column type
func (c ColumnElem) NotNull() ColumnElem {
	c.Constraints = append(c.Constraints, NotNull())
	return c
}

// Unique adds a unique constraint to column type
func (c ColumnElem) Unique() ColumnElem {
	c.Constraints = append(c.Constraints, Unique())
	c.Options.Unique = true
	return c
}

// Constraint adds a custom constraint to column type
func (c ColumnElem) Constraint(name string) ColumnElem {
	c.Constraints = append(c.Constraints, Constraint(name))
	return c
}

// conditional wrappers

// Like wraps the Like(col ColumnElem, pattern string)
func (c ColumnElem) Like(pattern string) SQLClause {
	return Like(c, pattern)
}

// NotIn wraps the NotIn(col ColumnElem, values ...interface{})
func (c ColumnElem) NotIn(values ...interface{}) SQLClause {
	return NotIn(c, values...)
}

// In wraps the In(col ColumnElem, values ...interface{})
func (c ColumnElem) In(values ...interface{}) SQLClause {
	return In(c, values...)
}

// NotEq wraps the NotEq(col ColumnElem, value interface{})
func (c ColumnElem) NotEq(value interface{}) SQLClause {
	return NotEq(c, value)
}

// Eq wraps the Eq(col ColumnElem, value interface{})
func (c ColumnElem) Eq(value interface{}) SQLClause {
	return Eq(c, value)
}

// Gt wraps the Gt(col ColumnElem, value interface{})
func (c ColumnElem) Gt(value interface{}) SQLClause {
	return Gt(c, value)
}

// Lt wraps the Lt(col ColumnElem, value interface{})
func (c ColumnElem) Lt(value interface{}) SQLClause {
	return Lt(c, value)
}

// Gte wraps the Gte(col ColumnElem, value interface{})
func (c ColumnElem) Gte(value interface{}) SQLClause {
	return Gte(c, value)
}

// Lte wraps the Lte(col ColumnElem, value interface{})
func (c ColumnElem) Lte(value interface{}) SQLClause {
	return Lte(c, value)
}
