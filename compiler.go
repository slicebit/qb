package qb

// Compiler is a visitor that produce SQL from various types of Clause
type Compiler interface {
	VisitAggregate(Context, AggregateClause) string
	VisitAlias(Context, AliasClause) string
	VisitBinary(Context, BinaryExpressionClause) string
	VisitBind(Context, BindClause) string
	VisitColumn(Context, ColumnElem) string
	VisitCombiner(Context, CombinerClause) string
	VisitDelete(Context, DeleteStmt) string
	VisitExists(Context, ExistsClause) string
	VisitForUpdate(Context, ForUpdateClause) string
	VisitHaving(Context, HavingClause) string
	VisitIn(Context, InClause) string
	VisitInsert(Context, InsertStmt) string
	VisitJoin(Context, JoinClause) string
	VisitLabel(Context, string) string
	VisitList(Context, ListClause) string
	VisitOrderBy(Context, OrderByClause) string
	VisitSelect(Context, SelectStmt) string
	VisitTable(Context, TableElem) string
	VisitText(Context, TextClause) string
	VisitUpdate(Context, UpdateStmt) string
	VisitUpsert(Context, UpsertStmt) string
	VisitWhere(Context, WhereClause) string
}
