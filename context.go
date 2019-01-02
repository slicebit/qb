package qb

// Context is the base definition of a compiler context
type Context interface {
	Compiler() Compiler
	Dialect() Dialect
	Binds() []interface{}
	ClearBinds()
	AddBinds(...interface{})
	DefaultTableName() string
	SetDefaultTableName(name string)
	InSubQuery() bool
	SetInSubQuery(inSubQuery bool)
}
