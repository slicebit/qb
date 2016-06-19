package qb

import (
	"fmt"
)

// Where generates a compilable where clause
func Where(clause Clause) WhereClause {
	return WhereClause{clause}
}

// WhereClause is the base of any where clause when using expression api
type WhereClause struct {
	clause Clause
}

// Build compiles the where clause, returns sql and bindings
func (c WhereClause) Build(dialect Dialect) (string, []interface{}) {
	sql, bindings := c.clause.Build(dialect)
	return fmt.Sprintf("WHERE %s", sql), bindings
}
