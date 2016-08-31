package qb

// Where generates a compilable where clause
func Where(clause Clause) WhereClause {
	return WhereClause{clause}
}

// WhereClause is the base of any where clause when using expression api
type WhereClause struct {
	clause Clause
}

// Accept compiles the where clause, returns sql
func (c WhereClause) Accept(context *CompilerContext) string {
	return context.Compiler.VisitWhere(context, c)
}
