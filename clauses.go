package qb

// SQLText returns a raw SQL clause
func SQLText(text string) TextClause {
	return TextClause{Text: text}
}

// TextClause is a raw SQL clause
type TextClause struct {
	Text string
}

// Accept calls the compiler VisitText method
func (c TextClause) Accept(context *CompilerContext) string {
	return context.Compiler.VisitText(context, c)
}

// List returns a list-of-clauses clause
func List(clauses ...Clause) ListClause {
	return ListClause{
		Clauses: clauses,
	}
}

// ListClause is a list of clause elements (for IN operator for example)
type ListClause struct {
	Clauses []Clause
}

// Accept calls the compiler VisitList method
func (c ListClause) Accept(context *CompilerContext) string {
	return context.Compiler.VisitList(context, c)
}

// Bind a value
func Bind(value interface{}) BindClause {
	return BindClause{
		Value: value,
	}
}

// BindClause binds a value to a placeholder
type BindClause struct {
	Value interface{}
}

// Accept calls the compiler VisitBind method
func (c BindClause) Accept(context *CompilerContext) string {
	return context.Compiler.VisitBind(context, c)
}

// GetClauseFrom returns the value if already a Clause, or make one
// if it is a scalar value
func GetClauseFrom(value interface{}) Clause {
	if clause, ok := value.(Clause); ok {
		return clause
	}
	// For now we assume any non-clause is a Value:
	return Bind(value)
}

// GetListFrom returns a list clause from any list
//
// If only one value is passed and is a ListClause, it is returned
// as-is.
// In any other case, a ListClause is built with each value wrapped
// by a Bind() if not already a Clause
func GetListFrom(values ...interface{}) Clause {
	if len(values) == 1 {
		if clause, ok := values[0].(ListClause); ok {
			return clause
		}
	}

	var clauses []Clause
	for _, value := range values {
		clauses = append(clauses, GetClauseFrom(value))
	}
	return List(clauses...)
}
