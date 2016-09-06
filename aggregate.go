package qb

// Avg function generates "avg(%s)" statement for clause
func Avg(clause Clause) AggregateClause {
	return Aggregate("AVG", clause)
}

// Count function generates "count(%s)" statement for clause
func Count(clause Clause) AggregateClause {
	return Aggregate("COUNT", clause)
}

// Sum function generates "sum(%s)" statement for clause
func Sum(clause Clause) AggregateClause {
	return Aggregate("SUM", clause)
}

// Min function generates "min(%s)" statement for clause
func Min(clause Clause) AggregateClause {
	return Aggregate("MIN", clause)
}

// Max function generates "max(%s)" statement for clause
func Max(clause Clause) AggregateClause {
	return Aggregate("MAX", clause)
}

// Aggregate generates a new aggregate clause given function & clause
func Aggregate(fn string, clause Clause) AggregateClause {
	return AggregateClause{fn, clause}
}

// AggregateClause is the base struct for building aggregate functions
type AggregateClause struct {
	fn     string
	clause Clause
}

// Accept calls the compiler VisitAggregate function
func (c AggregateClause) Accept(context *CompilerContext) string {
	return context.Compiler.VisitAggregate(context, c)
}
