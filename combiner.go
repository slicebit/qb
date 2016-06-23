package qb

import (
	"fmt"
	"strings"
)

// buildCombiners generates and or statements and join them appropriately
func buildCombiners(dialect Dialect, combiner string, clauses []Clause) (string, []interface{}) {
	sqls := []string{}
	bindings := []interface{}{}
	for _, c := range clauses {
		sql, values := c.Build(dialect)
		sqls = append(sqls, sql)
		bindings = append(bindings, values...)
	}

	return fmt.Sprintf("(%s)", strings.Join(sqls, fmt.Sprintf(" %s ", combiner))), bindings
}

// And generates an AndClause given conditional clauses
func And(clauses ...Clause) AndClause {
	return AndClause{clauses}
}

// AndClause is the base struct to keep and within the where clause
// It satisfies the Clause interface
type AndClause struct {
	clauses []Clause
}

// Build compiles the and clause, joins the sql, returns sql and bindings
func (c AndClause) Build(dialect Dialect) (string, []interface{}) {
	return buildCombiners(dialect, "AND", c.clauses)
}

// Or generates an OrClause given conditional clauses
func Or(clauses ...Clause) OrClause {
	return OrClause{clauses}
}

// OrClause is the base struct to keep or within the where clause
// It satisfies the Clause interface
type OrClause struct {
	clauses []Clause
}

// Build compiles the or clause, joins the sql, returns sql and bindings
func (c OrClause) Build(dialect Dialect) (string, []interface{}) {
	return buildCombiners(dialect, "OR", c.clauses)
}
