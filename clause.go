package qb

// Clause is the common interface of any SQL expression
type Clause interface {
	Accept(context *CompilerContext) string
}

// TableSQLClause is the common interface for ddl generators such as Column(), PrimaryKey(), ForeignKey().Ref(), etc.
type TableSQLClause interface {
	// String takes the dialect and returns the ddl as an sql string
	String(dialect Dialect) string
}

// Builder is the common interface for any statement builder in qb such as Insert(), Update(), Delete(), Select() query starters
type Builder interface {
	// Build takes a dialect and returns a stmt
	Build(dialect Dialect) *Stmt
}
