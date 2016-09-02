package qb

import (
	"fmt"
	"strings"
)

func NewCompilerContext(dialect Dialect) *CompilerContext {
	return &CompilerContext{
		Dialect:  dialect,
		Compiler: dialect.GetCompiler(),
		Vars:     make(map[string]interface{}),
	}
}

type CompilerContext struct {
	Binds            []interface{}
	DefaultTableName string
	Vars             map[string]interface{}

	Dialect  Dialect
	Compiler Compiler
}

type Compiler interface {
	VisitAggregate(*CompilerContext, AggregateClause) string
	VisitColumn(*CompilerContext, ColumnElem) string
	VisitCombiner(*CompilerContext, CombinerClause) string
	VisitCondition(*CompilerContext, Conditional) string
	VisitDelete(*CompilerContext, DeleteStmt) string
	VisitHaving(*CompilerContext, HavingClause) string
	VisitJoin(*CompilerContext, JoinClause) string
	VisitLabel(*CompilerContext, string) string
	VisitOrderBy(*CompilerContext, OrderByClause) string
	VisitSelect(*CompilerContext, SelectStmt) string
	VisitWhere(*CompilerContext, WhereClause) string
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

func (c SQLCompiler) VisitCombiner(context *CompilerContext, combiner CombinerClause) string {
	sqls := []string{}
	for _, c := range combiner.clauses {
		sql := c.Accept(context)
		sqls = append(sqls, sql)
	}

	return fmt.Sprintf("(%s)", strings.Join(sqls, fmt.Sprintf(" %s ", combiner.operator)))
}

func (c SQLCompiler) VisitCondition(context *CompilerContext, condition Conditional) string {
	var sql string
	key := condition.Col.Accept(context)

	switch condition.Op {
	case "IN":
		sql = fmt.Sprintf("%s %s (%s)", key, condition.Op, strings.Join(context.Dialect.Placeholders(condition.Values...), ", "))
		context.Binds = append(context.Binds, condition.Values...)
	case "NOT IN":
		sql = fmt.Sprintf("%s %s (%s)", key, condition.Op, strings.Join(context.Dialect.Placeholders(condition.Values...), ", "))
		context.Binds = append(context.Binds, condition.Values...)
	case "LIKE":
		sql = fmt.Sprintf("%s %s '%s'", key, condition.Op, condition.Values[0])
	default:
		sql = fmt.Sprintf("%s %s %s", key, condition.Op, context.Dialect.Placeholder())
		context.Binds = append(context.Binds, condition.Values...)
	}
	return sql
}

func (c SQLCompiler) VisitDelete(context *CompilerContext, delete DeleteStmt) string {
	sql := "DELETE FROM "
	sql += context.Compiler.VisitLabel(context, delete.table.Name)

	if delete.where != nil {
		sql += "\n" + delete.where.Accept(context)
	}

	returning := []string{}
	for _, c := range delete.returning {
		returning = append(returning, context.Dialect.Escape(c.Name))
	}

	if len(returning) > 0 {
		sql += "\nRETURNING " + strings.Join(returning, ", ")
	}

	return sql
}

func (c SQLCompiler) VisitHaving(context *CompilerContext, having HavingClause) string {
	aggSQL := having.aggregate.Accept(context)
	context.Binds = append(context.Binds, having.value)
	return fmt.Sprintf("HAVING %s %s %s", aggSQL, having.op, context.Dialect.Placeholder())
}

func (c SQLCompiler) VisitJoin(context *CompilerContext, join JoinClause) string {
	sql := fmt.Sprintf(
		"%s %s",
		join.joinType,
		context.Compiler.VisitLabel(context, join.table.Name),
	)
	if (join.fromCol.Name != "") || (join.col.Name != "") {
		sql += " ON " + join.fromCol.Accept(context) + " = " + join.col.Accept(context)
	}

	return sql
}

func (c SQLCompiler) VisitLabel(context *CompilerContext, label string) string {
	return c.Dialect.Escape(label)
}

func (c SQLCompiler) VisitOrderBy(context *CompilerContext, orderBy OrderByClause) string {
	cols := []string{}
	for _, c := range orderBy.columns {
		cols = append(cols, c.Accept(context))
	}

	return fmt.Sprintf("ORDER BY %s %s", strings.Join(cols, ", "), orderBy.t)
}

func (c SQLCompiler) VisitSelect(context *CompilerContext, select_ SelectStmt) string {
	lines := []string{}
	addLine := func(s string) {
		lines = append(lines, s)
	}
	if len(select_.joins) == 0 {
		context.DefaultTableName = select_.from.Name
	}

	// select
	columns := []string{}
	for _, c := range select_.sel {
		sql := c.Accept(context)
		columns = append(columns, sql)
	}
	addLine(fmt.Sprintf("SELECT %s", strings.Join(columns, ", ")))

	// from
	addLine(fmt.Sprintf("FROM %s", context.Dialect.Escape(select_.from.Name)))

	// joins
	for _, j := range select_.joins {
		addLine(j.Accept(context))
	}

	// where
	if select_.where != nil {
		addLine(select_.where.Accept(context))
	}

	// group by
	groupByCols := []string{}
	for _, c := range select_.groupBy {
		groupByCols = append(groupByCols, context.Dialect.Escape(c.Name))
	}
	if len(groupByCols) > 0 {
		addLine(fmt.Sprintf("GROUP BY %s", strings.Join(groupByCols, ", ")))
	}

	// having
	for _, h := range select_.having {
		sql := h.Accept(context)
		addLine(sql)
	}

	// order by
	if select_.orderBy != nil {
		sql := select_.orderBy.Accept(context)
		addLine(sql)
	}

	if (select_.offset != nil) && (select_.count != nil) {
		addLine(fmt.Sprintf("LIMIT %d OFFSET %d", *select_.count, *select_.offset))
	}

	return strings.Join(lines, "\n")
}

func (c SQLCompiler) VisitWhere(context *CompilerContext, where WhereClause) string {
	return fmt.Sprintf("WHERE %s", where.clause.Accept(context))
}
