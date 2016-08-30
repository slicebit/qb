package qb

import (
	"fmt"
)

// Where generates a compilable where clause
func Where(clause SQLClause) WhereSQLClause {
	return WhereSQLClause{clause}
}

// WhereSQLClause is the base of any where clause when using expression api
type WhereSQLClause struct {
	clause SQLClause
}

// Build compiles the where clause, returns sql and bindings
func (c WhereSQLClause) Build(dialect Dialect) (string, []interface{}) {
	sql, bindings := c.clause.Build(dialect)
	return fmt.Sprintf("WHERE %s", sql), bindings
}
