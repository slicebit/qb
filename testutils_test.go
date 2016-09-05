package qb

func asSQL(clause Clause, dialect Dialect) (string, []interface{}) {
	ctx := NewCompilerContext(dialect)
	return clause.Accept(ctx), ctx.Binds
}
