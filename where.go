package qb

// Where generates a compilable where clause
func Where(clauses ...Clause) WhereClause {
	var clause Clause
	if len(clauses) == 1 {
		clause = clauses[0]
	} else {
		clause = And(clauses...)
	}
	return WhereClause{clause}
}

// WhereClause is the base of any where clause when using expression api
type WhereClause struct {
	clause Clause
}

// Accept compiles the where clause, returns sql
func (c WhereClause) Accept(context Context) string {
	return context.Compiler().VisitWhere(context, c)
}

// And combine the current clause and the new ones with a And()
func (c WhereClause) And(clauses ...Clause) WhereClause {
	clauses = append([]Clause{c.clause}, clauses...)
	c.clause = And(clauses...)
	return c
}

// Or combine the current clause and the new ones with a Or()
func (c WhereClause) Or(clauses ...Clause) WhereClause {
	clauses = append([]Clause{c.clause}, clauses...)
	c.clause = Or(clauses...)
	return c
}
