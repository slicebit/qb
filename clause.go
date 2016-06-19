package qb

// Compiles is the standard interface for any compilable sql clause
// Compiling means to post process any sql clauses if needed such as escaping, putting placeholders, etc.
type Compiles interface {
	// Build is the key function of any sql clause
	// It returns sql as string and bindings as []interface{}
	Build(dialect Dialect) (string, []interface{})
}

// Clause is the key interface for any sql clause
type Clause interface {
	Compiles
	// String returns the dialect agnostic sql clause and bindings.
	// It returns :varname as placeholders instead of $n or ?.
	//String() (string, []interface{})
}

// TableClause is the common interface for ddl generators such as Column(), PrimaryKey(), ForeignKey().Ref(), etc.
type TableClause interface {
	// String takes the dialect and returns the ddl as an sql string
	String(dialect Dialect) string
}

// Builder is the common interface for any statement builder in qb such as Insert(), Update(), Delete(), Select() query starters
type Builder interface {
	// Build takes a dialect and returns a stmt
	Build(dialect Dialect) *Stmt
}