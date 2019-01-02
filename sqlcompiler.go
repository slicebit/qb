package qb

import (
	"fmt"
	"strings"
)

// NewSQLCompiler returns a new SQLCompiler
func NewSQLCompiler(dialect Dialect) SQLCompiler {
	return SQLCompiler{Dialect: dialect}
}

// SQLCompiler aims to provide a SQL ANSI-92 implementation of Compiler
type SQLCompiler struct {
	Dialect Dialect
}

// VisitAggregate compiles aggregate functions (COUNT, SUM...)
func (c SQLCompiler) VisitAggregate(context Context, aggregate AggregateClause) string {
	return fmt.Sprintf("%s(%s)", aggregate.fn, aggregate.clause.Accept(context))
}

// VisitAlias compiles a '<selectable> AS <aliasname>' SQL clause
func (SQLCompiler) VisitAlias(context Context, alias AliasClause) string {
	return fmt.Sprintf(
		"%s AS %s",
		alias.Selectable.Accept(context),
		context.Dialect().Escape(alias.Name),
	)
}

// VisitBinary compiles LEFT <op> RIGHT expressions
func (c SQLCompiler) VisitBinary(context Context, binary BinaryExpressionClause) string {
	return fmt.Sprintf(
		"%s %s %s",
		binary.Left.Accept(context),
		binary.Op,
		binary.Right.Accept(context),
	)
}

// VisitBind renders a bounded value
func (SQLCompiler) VisitBind(context Context, bind BindClause) string {
	context.AddBinds(bind.Value)
	return "?"
}

// VisitColumn returns a column name, optionnaly escaped depending on the dialect
// configuration
func (c SQLCompiler) VisitColumn(context Context, column ColumnElem) string {
	sql := ""
	if context.InSubQuery() || context.DefaultTableName() != column.Table {
		sql += c.Dialect.Escape(column.Table) + "."
	}
	sql += c.Dialect.Escape(column.Name)
	return sql
}

// VisitCombiner compiles AND and OR sql clauses
func (c SQLCompiler) VisitCombiner(context Context, combiner CombinerClause) string {
	sqls := []string{}
	for _, c := range combiner.clauses {
		sql := c.Accept(context)
		sqls = append(sqls, sql)
	}

	return fmt.Sprintf("(%s)", strings.Join(sqls, fmt.Sprintf(" %s ", combiner.operator)))
}

// VisitDelete compiles a DELETE statement
func (c SQLCompiler) VisitDelete(context Context, delete DeleteStmt) string {
	sql := "DELETE FROM " + delete.table.Accept(context)

	if delete.where != nil {
		sql += "\n" + delete.where.Accept(context)
	}

	returning := []string{}
	for _, c := range delete.returning {
		returning = append(returning, context.Dialect().Escape(c.Name))
	}

	if len(returning) > 0 {
		sql += "\nRETURNING " + strings.Join(returning, ", ")
	}

	return sql
}

// VisitExists compile a EXISTS clause
func (SQLCompiler) VisitExists(context Context, exists ExistsClause) string {
	var sql string
	if exists.Not {
		sql = "NOT "
	}
	sql += "EXISTS(%s)"
	context.SetInSubQuery(true)
	defer func() { context.SetInSubQuery(false) }()
	return fmt.Sprintf(sql, exists.Select.Accept(context))
}

// VisitForUpdate compiles a 'FOR UPDATE' clause
func (c SQLCompiler) VisitForUpdate(context Context, forUpdate ForUpdateClause) string {
	var sql = "FOR UPDATE"
	if len(forUpdate.Tables) != 0 {
		var tablenames []string
		for _, table := range forUpdate.Tables {
			tablenames = append(tablenames, table.Name)
		}
		sql += " OF " + strings.Join(tablenames, ", ")
	}
	return sql
}

// VisitHaving compiles a HAVING clause
func (c SQLCompiler) VisitHaving(context Context, having HavingClause) string {
	aggSQL := having.aggregate.Accept(context)
	return fmt.Sprintf("HAVING %s %s %s", aggSQL, having.op, Bind(having.value).Accept(context))
}

// VisitIn compiles a <left> (NOT) IN (<right>)
func (c SQLCompiler) VisitIn(context Context, in InClause) string {
	return fmt.Sprintf(
		"%s %s (%s)",
		in.Left.Accept(context),
		in.Op,
		in.Right.Accept(context),
	)
}

// VisitInsert compiles a INSERT statement
func (c SQLCompiler) VisitInsert(context Context, insert InsertStmt) string {
	context.SetDefaultTableName(insert.table.Name)
	defer func() { context.SetDefaultTableName("") }()

	cols := List()
	values := List()
	for k, v := range insert.values {
		cols.Clauses = append(cols.Clauses, insert.table.C(k))
		values.Clauses = append(values.Clauses, Bind(v))
	}

	sql := fmt.Sprintf(
		"INSERT INTO %s(%s)\nVALUES(%s)",
		insert.table.Accept(context),
		cols.Accept(context),
		values.Accept(context),
	)

	returning := []string{}
	for _, r := range insert.returning {
		returning = append(returning, r.Accept(context))
	}
	if len(insert.returning) > 0 {
		sql += fmt.Sprintf(
			"\nRETURNING %s",
			strings.Join(returning, ", "),
		)
	}

	return sql
}

// VisitJoin compiles a JOIN (ON) clause
func (c SQLCompiler) VisitJoin(context Context, join JoinClause) string {
	sql := fmt.Sprintf(
		"%s\n%s %s",
		join.Left.Accept(context),
		join.JoinType,
		join.Right.Accept(context),
	)
	if join.OnClause != nil {
		sql += " ON " + join.OnClause.Accept(context)
	}

	return sql
}

// VisitLabel returns a single label, optionally escaped
func (c SQLCompiler) VisitLabel(context Context, label string) string {
	return c.Dialect.Escape(label)
}

// VisitList compiles a list of values
func (c SQLCompiler) VisitList(context Context, list ListClause) string {
	var clauses []string
	for _, clause := range list.Clauses {
		clauses = append(clauses, clause.Accept(context))
	}
	return strings.Join(clauses, ", ")
}

// VisitOrderBy compiles a ORDER BY sql clause
func (c SQLCompiler) VisitOrderBy(context Context, OrderByClause OrderByClause) string {
	cols := []string{}
	for _, c := range OrderByClause.columns {
		cols = append(cols, c.Accept(context))
	}

	return fmt.Sprintf("ORDER BY %s %s", strings.Join(cols, ", "), OrderByClause.t)
}

// VisitSelect compiles a SELECT statement
func (c SQLCompiler) VisitSelect(context Context, selectStmt SelectStmt) string {
	lines := []string{}
	addLine := func(s string) {
		lines = append(lines, s)
	}
	if !context.InSubQuery() && selectStmt.FromClause != nil {
		// context.DefaultTableName = selectStmt.FromClause.DefaultName()
		context.SetDefaultTableName(selectStmt.FromClause.DefaultName())
	}

	// select
	columns := []string{}
	for _, c := range selectStmt.SelectList {
		sql := c.Accept(context)
		columns = append(columns, sql)
	}
	addLine(fmt.Sprintf("SELECT %s", strings.Join(columns, ", ")))

	// from
	if selectStmt.FromClause != nil {
		addLine(fmt.Sprintf("FROM %s", selectStmt.FromClause.Accept(context)))
	}

	// where
	if selectStmt.WhereClause != nil {
		addLine(selectStmt.WhereClause.Accept(context))
	}

	// group by
	groupByCols := []string{}
	for _, c := range selectStmt.GroupByClause {
		groupByCols = append(groupByCols, context.Dialect().Escape(c.Name))
	}
	if len(groupByCols) > 0 {
		addLine(fmt.Sprintf("GROUP BY %s", strings.Join(groupByCols, ", ")))
	}

	// having
	for _, h := range selectStmt.HavingClause {
		sql := h.Accept(context)
		addLine(sql)
	}

	// order by
	if selectStmt.OrderByClause != nil {
		sql := selectStmt.OrderByClause.Accept(context)
		addLine(sql)
	}

	if (selectStmt.OffsetValue != nil) || (selectStmt.LimitValue != nil) {
		var tokens []string
		if selectStmt.LimitValue != nil {
			tokens = append(tokens, fmt.Sprintf("LIMIT %d", *selectStmt.LimitValue))
		}
		if selectStmt.OffsetValue != nil {
			tokens = append(tokens, fmt.Sprintf("OFFSET %d", *selectStmt.OffsetValue))
		}
		addLine(strings.Join(tokens, " "))
	}

	if selectStmt.ForUpdateClause != nil {
		addLine(selectStmt.ForUpdateClause.Accept(context))
	}

	return strings.Join(lines, "\n")
}

// VisitTable returns a table name, optionally escaped
func (SQLCompiler) VisitTable(context Context, table TableElem) string {
	return context.Compiler().VisitLabel(context, table.Name)
}

// VisitText return a raw SQL clause as is
func (SQLCompiler) VisitText(context Context, text TextClause) string {
	return text.Text
}

// VisitUpdate compiles a UPDATE statement
func (c SQLCompiler) VisitUpdate(context Context, update UpdateStmt) string {
	context.SetDefaultTableName(update.table.Name)
	defer func() { context.SetDefaultTableName("") }()

	sql := "UPDATE " + update.table.Accept(context)

	sets := List()

	for k, v := range update.values {
		sets.Clauses = append(sets.Clauses,
			Eq(update.table.C(k), Bind(v)))
	}

	if len(sets.Clauses) > 0 {
		sql += "\nSET " + sets.Accept(context)
	}

	if update.where != nil {
		sql += "\n" + update.where.Accept(context)
	}

	returning := []string{}
	for _, c := range update.returning {
		returning = append(returning, context.Dialect().Escape(c.Name))
	}

	if len(returning) > 0 {
		sql += "\nRETURNING " + strings.Join(returning, ", ")
	}

	return sql
}

// VisitUpsert is not implemented and will panic.
// It should be implemented in each dialect
func (c SQLCompiler) VisitUpsert(context Context, upsert UpsertStmt) string {
	panic("Upsert is not Implemented in this compiler")
}

// VisitWhere compiles a WHERE clause
func (c SQLCompiler) VisitWhere(context Context, where WhereClause) string {
	return fmt.Sprintf("WHERE %s", where.clause.Accept(context))
}
