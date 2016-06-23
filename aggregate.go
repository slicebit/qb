package qb

import (
	"fmt"
)

// Avg function generates "avg(%s)" statement for column
func Avg(column ColumnElem) AggregateClause {
	return Aggregate("AVG", column)
}

// Count function generates "count(%s)" statement for column
func Count(column ColumnElem) AggregateClause {
	return Aggregate("COUNT", column)
}

// Sum function generates "sum(%s)" statement for column
func Sum(column ColumnElem) AggregateClause {
	return Aggregate("SUM", column)
}

// Min function generates "min(%s)" statement for column
func Min(column ColumnElem) AggregateClause {
	return Aggregate("MIN", column)
}

// Max function generates "max(%s)" statement for column
func Max(column ColumnElem) AggregateClause {
	return Aggregate("MAX", column)
}

// Aggregate generates a new aggregate clause given function & column
func Aggregate(fn string, column ColumnElem) AggregateClause {
	return AggregateClause{fn, column}
}

// AggregateClause is the base struct for building aggregate functions
type AggregateClause struct {
	fn     string
	column ColumnElem
}

// Build compiles the aggregate clause and returns the sql and bindings
func (c AggregateClause) Build(dialect Dialect) (string, []interface{}) {
	bindings := []interface{}{}
	sql := fmt.Sprintf("%s(%s)", c.fn, dialect.Escape(c.column.Name))
	return sql, bindings
}
