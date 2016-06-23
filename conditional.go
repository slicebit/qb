package qb

import (
	"fmt"
	"strings"
)

// conditional generators, comparator functions

// Like generates a like conditional sql clause
func Like(col ColumnElem, pattern string) Conditional {
	return Condition(col, "LIKE", pattern)
}

// NotIn generates a not in conditional sql clause
func NotIn(col ColumnElem, values ...interface{}) Conditional {
	return Condition(col, "NOT IN", values...)
}

// In generates an in conditional sql clause
func In(col ColumnElem, values ...interface{}) Conditional {
	return Condition(col, "IN", values...)
}

// NotEq generates a not equal conditional sql clause
func NotEq(col ColumnElem, value interface{}) Conditional {
	return Condition(col, "!=", value)
}

// Eq generates a equals conditional sql clause
func Eq(col ColumnElem, value interface{}) Conditional {
	return Condition(col, "=", value)
}

// Gt generates a greater than conditional sql clause
func Gt(col ColumnElem, value interface{}) Conditional {
	return Condition(col, ">", value)
}

// St generates a smaller than conditional sql clause
func St(col ColumnElem, value interface{}) Conditional {
	return Condition(col, "<", value)
}

// Gte generates a greater than or equal to conditional sql clause
func Gte(col ColumnElem, value interface{}) Conditional {
	return Condition(col, ">=", value)
}

// Ste generates a smaller than or equal to conditional sql clause
func Ste(col ColumnElem, value interface{}) Conditional {
	return Condition(col, "<=", value)
}

// Condition generates a condition object to use in update, delete & select statements
func Condition(col ColumnElem, op string, values ...interface{}) Conditional {
	return Conditional{col, values, op}
}

// Conditional is the base struct for any conditional statements in sql clauses
type Conditional struct {
	Col    ColumnElem
	Values []interface{}
	Op     string
}

// Build compiles the conditional element
func (c Conditional) Build(dialect Dialect) (string, []interface{}) {
	var sql string
	key := dialect.Escape(c.Col.Name)
	if c.Col.Table != "" {
		key = fmt.Sprintf("%s.%s", dialect.Escape(c.Col.Table), key)
	}

	switch c.Op {
	case "IN":
		sql = fmt.Sprintf("%s %s (%s)", key, c.Op, strings.Join(dialect.Placeholders(c.Values...), ", "))
		return sql, c.Values
	case "NOT IN":
		sql = fmt.Sprintf("%s %s (%s)", key, c.Op, strings.Join(dialect.Placeholders(c.Values...), ", "))
		return sql, c.Values
	case "LIKE":
		sql = fmt.Sprintf("%s %s '%s'", key, c.Op, c.Values[0])
		return sql, []interface{}{}
	default:
		sql = fmt.Sprintf("%s %s %s", key, c.Op, dialect.Placeholder())
		return sql, c.Values
	}
}
