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
