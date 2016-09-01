package qb

// Selectable is any clause from which we can select columns and is suitable
// as a FROM clause element
type Selectable interface {
	Clause
	C(column string) ColumnElem
	DefaultName() string
}

// Select generates a select statement and returns it
func Select(clauses ...Clause) SelectStmt {
	return SelectStmt{
		sel:     clauses,
		joins:   []JoinClause{},
		groupBy: []ColumnElem{},
		having:  []HavingClause{},
	}
}

// SelectStmt is the base struct for building select statements
type SelectStmt struct {
	sel     []Clause
	from    Selectable
	joins   []JoinClause
	groupBy []ColumnElem
	orderBy *OrderByClause
	having  []HavingClause
	where   *WhereClause
	offset  *int
	count   *int
}

// From sets the from selectable of select statement
func (s SelectStmt) From(selectable Selectable) SelectStmt {
	s.from = selectable
	return s
}

// Where sets the where clause of select statement
func (s SelectStmt) Where(clause Clause) SelectStmt {
	where := Where(clause)
	s.where = &where
	return s
}

// InnerJoin appends an inner join clause to the select statement
func (s SelectStmt) InnerJoin(right Selectable, leftCol ColumnElem, rightCol ColumnElem) SelectStmt {
	join := join("INNER JOIN", s.from, right, leftCol, rightCol)
	s.joins = append(s.joins, join)
	return s
}

// CrossJoin appends an cross join clause to the select statement
func (s SelectStmt) CrossJoin(right Selectable) SelectStmt {
	join := join("CROSS JOIN", s.from, right, ColumnElem{}, ColumnElem{})
	s.joins = append(s.joins, join)
	return s
}

// LeftJoin appends an left outer join clause to the select statement
func (s SelectStmt) LeftJoin(right Selectable, leftCol ColumnElem, rightCol ColumnElem) SelectStmt {
	join := join("LEFT OUTER JOIN", s.from, right, leftCol, rightCol)
	s.joins = append(s.joins, join)
	return s
}

// RightJoin appends a right outer join clause to select statement
func (s SelectStmt) RightJoin(right Selectable, leftCol ColumnElem, rightCol ColumnElem) SelectStmt {
	join := join("RIGHT OUTER JOIN", s.from, right, leftCol, rightCol)
	s.joins = append(s.joins, join)
	return s
}

// OrderBy generates an OrderByClause and sets select statement's orderbyclause
// OrderBy(usersTable.C("id")).Asc()
// OrderBy(usersTable.C("email")).Desc()
func (s SelectStmt) OrderBy(columns ...ColumnElem) SelectStmt {
	s.orderBy = &OrderByClause{columns, "ASC"}
	return s
}

// Asc sets the t type of current order by clause
// NOTE: Please use it after calling OrderBy()
func (s SelectStmt) Asc() SelectStmt {
	s.orderBy.t = "ASC"
	return s
}

// Desc sets the t type of current order by clause
// NOTE: Please use it after calling OrderBy()
func (s SelectStmt) Desc() SelectStmt {
	s.orderBy.t = "DESC"
	return s
}

// GroupBy appends columns to group by clause of the select statement
func (s SelectStmt) GroupBy(cols ...ColumnElem) SelectStmt {
	s.groupBy = append(s.groupBy, cols...)
	return s
}

// Having appends a having clause to select statement
func (s SelectStmt) Having(aggregate AggregateClause, op string, value interface{}) SelectStmt {
	s.having = append(s.having, HavingClause{aggregate, op, value})
	return s
}

// Limit sets the offset & count values of the select statement
func (s SelectStmt) Limit(offset int, count int) SelectStmt {
	s.offset = &offset
	s.count = &count
	return s
}

// Build compiles the select statement and returns the Stmt
func (s SelectStmt) Build(dialect Dialect) *Stmt {
	defer dialect.Reset()

	context := NewCompilerContext(dialect)
	statement := Statement()
	statement.AddSQLClause(context.Compiler.VisitSelect(context, s))
	statement.AddBinding(context.Binds...)

	return statement
}

func join(joinType string, left Selectable, right Selectable, leftCol ColumnElem, rightCol ColumnElem) JoinClause {
	return JoinClause{
		joinType,
		left,
		right,
		leftCol,
		rightCol,
	}
}

// JoinClause is the base struct for generating join clauses when using select
// It satisfies Clause interface
type JoinClause struct {
	joinType string
	left     Selectable
	right    Selectable
	leftCol  ColumnElem
	rightCol ColumnElem
}

func (c JoinClause) Accept(context *CompilerContext) string {
	return context.Compiler.VisitJoin(context, c)
}

// OrderByClause is the base struct for generating order by clauses when using select
// It satisfies SQLClause interface
type OrderByClause struct {
	columns []ColumnElem
	t       string
}

// Accept generates an order by clause
func (c OrderByClause) Accept(context *CompilerContext) string {
	return context.Compiler.VisitOrderBy(context, c)
}

// HavingClause is the base struct for generating having clauses when using select
// It satisfies SQLClause interface
type HavingClause struct {
	aggregate AggregateClause
	op        string
	value     interface{}
}

// Accept generates having sql & bindings out of HavingClause struct
func (c HavingClause) Accept(context *CompilerContext) string {
	return context.Compiler.VisitHaving(context, c)
}
