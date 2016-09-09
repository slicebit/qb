package qb

import "fmt"

// Selectable is any clause from which we can select columns and is suitable
// as a FROM clause element
type Selectable interface {
	Clause
	All() []Clause
	ColumnList() []ColumnElem
	C(column string) ColumnElem
	DefaultName() string
}

// Select generates a select statement and returns it
func Select(clauses ...Clause) SelectStmt {
	return SelectStmt{
		sel:     clauses,
		groupBy: []ColumnElem{},
		having:  []HavingClause{},
	}
}

// SelectStmt is the base struct for building select statements
type SelectStmt struct {
	sel         []Clause
	from        Selectable
	groupBy     []ColumnElem
	orderBy     *OrderByClause
	having      []HavingClause
	WhereClause *WhereClause
	offset      *int
	count       *int
}

// Select sets the selected columns
func (s SelectStmt) Select(clauses ...Clause) SelectStmt {
	s.sel = clauses
	return s
}

// From sets the from selectable of select statement
func (s SelectStmt) From(selectable Selectable) SelectStmt {
	s.from = selectable
	return s
}

// Where sets the where clause of select statement
func (s SelectStmt) Where(clauses ...Clause) SelectStmt {
	where := Where(clauses...)
	s.WhereClause = &where
	return s
}

// InnerJoin appends an inner join clause to the select statement
func (s SelectStmt) InnerJoin(right Selectable, onClause ...Clause) SelectStmt {
	return s.From(Join("INNER JOIN", s.from, right, onClause...))
}

// CrossJoin appends an cross join clause to the select statement
func (s SelectStmt) CrossJoin(right Selectable) SelectStmt {
	return s.From(Join("CROSS JOIN", s.from, right, nil))
}

// LeftJoin appends an left outer join clause to the select statement
func (s SelectStmt) LeftJoin(right Selectable, onClause ...Clause) SelectStmt {
	return s.From(Join("LEFT OUTER JOIN", s.from, right, onClause...))
}

// RightJoin appends a right outer join clause to select statement
func (s SelectStmt) RightJoin(right Selectable, onClause ...Clause) SelectStmt {
	return s.From(Join("RIGHT OUTER JOIN", s.from, right, onClause...))
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

// Accept calls the compiler VisitSelect method
func (s SelectStmt) Accept(context *CompilerContext) string {
	return context.Compiler.VisitSelect(context, s)
}

// Build compiles the select statement and returns the Stmt
func (s SelectStmt) Build(dialect Dialect) *Stmt {
	context := NewCompilerContext(dialect)
	statement := Statement()
	statement.AddSQLClause(s.Accept(context))
	statement.AddBinding(context.Binds...)

	return statement
}

type joinOnClauseCandidate struct {
	source TableElem
	fkey   ForeignKeyConstraint
	target TableElem
}

func getTable(sel Selectable) (TableElem, bool) {
	switch t := sel.(type) {
	case TableElem:
		return t, true
	case *TableElem:
		return *t, true
	default:
		return TableElem{}, false
	}
}

// GuessJoinOnClause finds a join 'ON' clause between two tables
func GuessJoinOnClause(left Selectable, right Selectable) Clause {
	leftTable, ok := getTable(left)
	if !ok {
		panic("left Selectable is not a Table: Cannot guess join onClause")
	}
	rightTable, ok := getTable(right)
	if !ok {
		panic("right Selectable is not a Table: Cannot guess join onClause")
	}

	var candidates []joinOnClauseCandidate

	for _, fkey := range leftTable.ForeignKeyConstraints.FKeys {
		if fkey.RefTable != rightTable.Name {
			continue
		}
		candidates = append(
			candidates,
			joinOnClauseCandidate{leftTable, fkey, rightTable})
	}

	for _, fkey := range rightTable.ForeignKeyConstraints.FKeys {
		if fkey.RefTable != leftTable.Name {
			continue
		}
		candidates = append(
			candidates,
			joinOnClauseCandidate{rightTable, fkey, leftTable})
	}
	switch len(candidates) {
	case 0:
		panic(fmt.Sprintf(
			"No foreign keys found between %s and %s",
			leftTable.Name, rightTable.Name))
	case 1:
		candidate := candidates[0]
		var clauses []Clause
		for i, col := range candidate.fkey.Cols {
			refCol := candidate.fkey.RefCols[i]
			clauses = append(
				clauses,
				Eq(candidate.source.C(col), candidate.target.C(refCol)),
			)
		}
		if len(clauses) == 1 {
			return clauses[0]
		}
		return And(clauses...)
	default:
		panic(fmt.Sprintf(
			"Found %d foreign keys between %s and %s",
			len(candidates), leftTable.Name, rightTable.Name))
	}
}

// MakeJoinOnClause assemble a 'ON' clause for a join from either:
// 0 clause: attempt to guess the join clause (only if left & right are tables),
//           otherwise panics
// 1 clause: returns it
// 2 clauses: returns a Eq() of both
// otherwise if panics
func MakeJoinOnClause(left Selectable, right Selectable, onClause ...Clause) Clause {
	switch len(onClause) {
	case 0:
		return GuessJoinOnClause(left, right)
	case 1:
		return onClause[0]
	case 2:
		return Eq(onClause[0], onClause[1])
	default:
		panic("Cannot make a join condition with more than 2 clauses")
	}
}

// Join returns a new JoinClause
// onClause can be one of:
// - 0 clause: attempt to guess the join clause (only if left & right are tables),
//             otherwise panics
// - 1 clause: use it directly
// - 2 clauses: use a Eq() of both
func Join(joinType string, left Selectable, right Selectable, onClause ...Clause) JoinClause {
	return JoinClause{
		JoinType: joinType,
		Left:     left,
		Right:    right,
		OnClause: MakeJoinOnClause(left, right, onClause...),
	}
}

// JoinClause is the base struct for generating join clauses when using select
// It satisfies Clause interface
type JoinClause struct {
	JoinType string
	Left     Selectable
	Right    Selectable
	OnClause Clause
}

// Accept calls the compiler VisitJoin method
func (c JoinClause) Accept(context *CompilerContext) string {
	return context.Compiler.VisitJoin(context, c)
}

// All returns the columns from both sides of the join
func (c JoinClause) All() []Clause {
	return append(c.Left.All(), c.Right.All()...)
}

// ColumnList returns the columns from both sides of the join
func (c JoinClause) ColumnList() []ColumnElem {
	return append(c.Left.ColumnList(), c.Right.ColumnList()...)
}

// C returns the first column with the given name
// If columns from both sides of the join match the name,
// the one from the left side will be returned.
func (c JoinClause) C(name string) ColumnElem {
	for _, c := range c.ColumnList() {
		if c.Name == name {
			return c
		}
	}
	panic(fmt.Sprintf("No such column '%s' in join %v", name, c))
}

// DefaultName returns an empty string because Joins have no name by default
func (c JoinClause) DefaultName() string {
	return ""
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

// Accept calls the compiler VisitHaving function
func (c HavingClause) Accept(context *CompilerContext) string {
	return context.Compiler.VisitHaving(context, c)
}

// Alias returns a new AliasClause
func Alias(name string, selectable Selectable) AliasClause {
	return AliasClause{
		Name:       name,
		Selectable: selectable,
	}
}

// AliasClause is a ALIAS sql clause
type AliasClause struct {
	Name       string
	Selectable Selectable
}

// Accept calls the compiler VisitAlias function
func (c AliasClause) Accept(context *CompilerContext) string {
	return context.Compiler.VisitAlias(context, c)
}

// C returns the aliased selectable column with the given name.
// Before returning it, the 'Table' field is updated with alias
// name so that they can be used in Select()
func (c AliasClause) C(name string) ColumnElem {
	col := c.Selectable.C(name)
	col.Table = c.Name
	return col
}

// All returns the aliased selectable columns with their "Table"
// field updated with the alias name
func (c AliasClause) All() []Clause {
	var clauses []Clause
	for _, col := range c.ColumnList() {
		clauses = append(clauses, col)
	}
	return clauses
}

// ColumnList returns the aliased selectable columns with their "Table"
// field updated with the alias name
func (c AliasClause) ColumnList() []ColumnElem {
	var cols []ColumnElem
	for _, col := range c.Selectable.ColumnList() {
		col.Table = c.Name
		cols = append(cols, col)
	}
	return cols
}

// DefaultName returns the alias name
func (c AliasClause) DefaultName() string {
	return c.Name
}
