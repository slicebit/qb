package qb

import (
	"fmt"
)

func NewCompilerContext(compiler Compiler) *CompilerContext {
	return &CompilerContext{
		Compiler: compiler,
		Vars:     make(map[string]interface{}),
	}
}

type CompilerContext struct {
	Binds            []interface{}
	DefaultTableName string
	Vars             map[string]interface{}

	Compiler Compiler
}

type Compiler interface {
	VisitAggregate(*CompilerContext, AggregateClause) string
	VisitColumn(*CompilerContext, ColumnElem) string
}

type SQLCompiler struct {
	Dialect Dialect
}

func (c SQLCompiler) VisitAggregate(context *CompilerContext, aggregate AggregateClause) string {
	return fmt.Sprintf("%s(%s)", aggregate.fn, aggregate.column.Accept(context))
}

func (c SQLCompiler) VisitColumn(context *CompilerContext, column ColumnElem) string {
	sql := ""
	if context.DefaultTableName != column.Table {
		sql += c.Dialect.Escape(column.Table) + "."
	}
	sql += c.Dialect.Escape(column.Name)
	return sql
}
