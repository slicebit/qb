package qb

import "fmt"

// Column generates a ColumnElem given name and type
func Column(name string, t TypeElem) ColumnElem {
	return ColumnElem{name, t, ""}
}

// ColumnElem is the definition of any columns defined in a table
type ColumnElem struct {
	Name  string
	Type  TypeElem
	Table string // This field should be lazily set by Table() function
}

// String returns the column element as an sql clause
// It satisfies the TableClause interface
func (c ColumnElem) String(dialect Dialect) string {
	return fmt.Sprintf("%s %s", dialect.Escape(c.Name), c.Type.String())
}

// Build compiles the column element and returns sql, bindings
// It satisfies the Clause interface
func (c ColumnElem) Build(dialect Dialect) (string, []interface{}) {
	return fmt.Sprintf("%s", dialect.Escape(c.Name)), []interface{}{}
}

// conditional wrappers

// Like wraps the Like(col ColumnElem, pattern string)
func (c ColumnElem) Like(pattern string) Conditional {
	return Like(c, pattern)
}

// NotIn wraps the NotIn(col ColumnElem, values ...interface{})
func (c ColumnElem) NotIn(values ...interface{}) Conditional {
	return NotIn(c, values...)
}

// In wraps the In(col ColumnElem, values ...interface{})
func (c ColumnElem) In(values ...interface{}) Conditional {
	return In(c, values...)
}

// NotEq wraps the NotEq(col ColumnElem, value interface{})
func (c ColumnElem) NotEq(value interface{}) Conditional {
	return NotEq(c, value)
}

// Eq wraps the Eq(col ColumnElem, value interface{})
func (c ColumnElem) Eq(value interface{}) Conditional {
	return Eq(c, value)
}

// Gt wraps the Gt(col ColumnElem, value interface{})
func (c ColumnElem) Gt(value interface{}) Conditional {
	return Gt(c, value)
}

// St wraps the St(col ColumnElem, value interface{})
func (c ColumnElem) St(value interface{}) Conditional {
	return St(c, value)
}

// Gte wraps the Gte(col ColumnElem, value interface{})
func (c ColumnElem) Gte(value interface{}) Conditional {
	return Gte(c, value)
}

// Ste wraps the Ste(col ColumnElem, value interface{})
func (c ColumnElem) Ste(value interface{}) Conditional {
	return Ste(c, value)
}
