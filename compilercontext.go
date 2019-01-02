package qb

// NewCompilerContext initialize a new compiler context
func NewCompilerContext(dialect Dialect) Context {
	return &CompilerContext{
		dialect:  dialect,
		compiler: dialect.GetCompiler(),
		binds:    []interface{}{},
	}
}

// CompilerContext is a data structure passed to all the Compiler visit
// functions. It contains the bindings, links to the Dialect and Compiler
// being used, and some contextual informations that can be used by the
// compiler functions to communicate during the compilation.
type CompilerContext struct {
	binds            []interface{}
	defaultTableName string
	inSubQuery       bool

	dialect  Dialect
	compiler Compiler
}

// Compiler returns the compiler in the context
func (ctx *CompilerContext) Compiler() Compiler {
	return ctx.compiler
}

// Dialect returns the dialect in the context
func (ctx *CompilerContext) Dialect() Dialect {
	return ctx.dialect
}

// Binds returns the binds registered to the context
func (ctx *CompilerContext) Binds() []interface{} {
	return ctx.binds
}

// ClearBinds clears all the binds that are registered to the context
func (ctx *CompilerContext) ClearBinds() {
	ctx.binds = []interface{}{}
}

// AddBinds registers new binds to the context
func (ctx *CompilerContext) AddBinds(binds ...interface{}) {
	for i := 0; i < len(binds); i++ {
		ctx.binds = append(ctx.binds, binds[i])
	}
}

// DefaultTableName returns the default table name registered to the context
func (ctx *CompilerContext) DefaultTableName() string {
	return ctx.defaultTableName
}

// SetDefaultTableName sets the default table name in context
func (ctx *CompilerContext) SetDefaultTableName(name string) {
	ctx.defaultTableName = name
}

// InSubQuery returns if InSubQuery is enabled
func (ctx *CompilerContext) InSubQuery() bool {
	return ctx.inSubQuery
}

// SetInSubQuery sets the inSubQuery in the context
func (ctx *CompilerContext) SetInSubQuery(inSubQuery bool) {
	ctx.inSubQuery = inSubQuery
}
