package qb

import (
	"fmt"
	"strings"
)

// buildCombiners generates and or statements and join them appropriately
func buildCombiners(dialect Dialect, combiner string, clauses []SQLClause) (string, []interface{}) {
	sqls := []string{}
	bindings := []interface{}{}
	for _, c := range clauses {
		sql, values := c.Build(dialect)
		sqls = append(sqls, sql)
		bindings = append(bindings, values...)
	}

	return fmt.Sprintf("(%s)", strings.Join(sqls, fmt.Sprintf(" %s ", combiner))), bindings
}

// And generates an AndSQLClause given conditional clauses
func And(clauses ...SQLClause) AndSQLClause {
	return AndSQLClause{clauses}
}

// AndSQLClause is the base struct to keep and within the where clause
// It satisfies the SQLClause interface
type AndSQLClause struct {
	clauses []SQLClause
}

// Build compiles the and clause, joins the sql, returns sql and bindings
func (c AndSQLClause) Build(dialect Dialect) (string, []interface{}) {
	return buildCombiners(dialect, "AND", c.clauses)
}

// Or generates an OrSQLClause given conditional clauses
func Or(clauses ...SQLClause) OrSQLClause {
	return OrSQLClause{clauses}
}

// OrSQLClause is the base struct to keep or within the where clause
// It satisfies the SQLClause interface
type OrSQLClause struct {
	clauses []SQLClause
}

// Build compiles the or clause, joins the sql, returns sql and bindings
func (c OrSQLClause) Build(dialect Dialect) (string, []interface{}) {
	return buildCombiners(dialect, "OR", c.clauses)
}
