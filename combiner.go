package qb

import (
	"strings"
	"fmt"
)

// buildCombiners generats and or statements and join them appropriately
func buildCombiners(dialect Dialect, combiner string, conditionals []Conditional) (string, []interface{}) {
	sqls := []string{}
	bindings := []interface{}{}
	for _, c := range conditionals {
		sql, values := c.Build(dialect)
		sqls = append(sqls, sql)
		bindings = append(bindings, values...)
	}

	return strings.Join(sqls, fmt.Sprintf(" %s ", combiner)), bindings
}

// And generates an AndClause given conditional clauses
func And(conditions ...Conditional) AndClause {
	return AndClause{conditions}
}

// AndClause is the base struct to keep and within the where clause
// It satisfies the Clause interface
type AndClause struct {
	conditions []Conditional
}

// Build compiles the and clause, joins the sql, returns sql and bindings
func (c AndClause) Build(dialect Dialect) (string, []interface{}) {
	return buildCombiners(dialect, "AND", c.conditions)
}

// Or generates an OrClause given conditional clauses
func Or(conditions ...Conditional) OrClause {
	return OrClause{conditions}
}

// OrClause is the base struct to keep or within the where clause
// It satisfies the Clause interface
type OrClause struct {
	conditions []Conditional
}

// Build compiles the or clause, joins the sql, returns sql and bindings
func (c OrClause) Build(dialect Dialect) (string, []interface{}) {
	return buildCombiners(dialect, "OR", c.conditions)
}
