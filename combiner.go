package qb

// And generates an AndClause given conditional clauses
func And(clauses ...Clause) CombinerClause {
	return CombinerClause{"AND", clauses}
}

// Or generates an AndClause given conditional clauses
func Or(clauses ...Clause) CombinerClause {
	return CombinerClause{"OR", clauses}
}

// CombinerClause is for OR and AND clauses
type CombinerClause struct {
	operator string
	clauses  []Clause
}

// Accept calls the compiler VisitCombiner entry point
func (c CombinerClause) Accept(context Context) string {
	return context.Compiler().VisitCombiner(context, c)
}
