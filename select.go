package qb

import (
	"fmt"
	"strings"
)

// Select generates a select statement and returns it
func Select(clauses ...SQLClause) SelectStmt {
	return SelectStmt{
		sel:     clauses,
		joins:   []JoinSQLClause{},
		groupBy: []ColumnElem{},
		having:  []HavingSQLClause{},
	}
}

// SelectStmt is the base struct for building select statements
type SelectStmt struct {
	sel     []SQLClause
	from    TableElem
	joins   []JoinSQLClause
	groupBy []ColumnElem
	orderBy *OrderBySQLClause
	having  []HavingSQLClause
	where   *WhereSQLClause
	offset  *int
	count   *int
}

// From sets the from table of select statement
func (s SelectStmt) From(table TableElem) SelectStmt {
	s.from = table
	return s
}

// Where sets the where clause of select statement
func (s SelectStmt) Where(clause SQLClause) SelectStmt {
	where := Where(clause)
	s.where = &where
	return s
}

// InnerJoin appends an inner join clause to the select statement
func (s SelectStmt) InnerJoin(table TableElem, fromCol ColumnElem, col ColumnElem) SelectStmt {
	join := join("INNER JOIN", s.from, table, fromCol, col)
	s.joins = append(s.joins, join)
	return s
}

// CrossJoin appends an cross join clause to the select statement
func (s SelectStmt) CrossJoin(table TableElem) SelectStmt {
	join := join("CROSS JOIN", s.from, table, ColumnElem{}, ColumnElem{})
	s.joins = append(s.joins, join)
	return s
}

// LeftJoin appends an left outer join clause to the select statement
func (s SelectStmt) LeftJoin(table TableElem, fromCol ColumnElem, col ColumnElem) SelectStmt {
	join := join("LEFT OUTER JOIN", s.from, table, fromCol, col)
	s.joins = append(s.joins, join)
	return s
}

// RightJoin appends a right outer join clause to select statement
func (s SelectStmt) RightJoin(table TableElem, fromCol ColumnElem, col ColumnElem) SelectStmt {
	join := join("RIGHT OUTER JOIN", s.from, table, fromCol, col)
	s.joins = append(s.joins, join)
	return s
}

// OrderBy generates an OrderBySQLClause and sets select statement's orderbyclause
// OrderBy(usersTable.C("id")).Asc()
// OrderBy(usersTable.C("email")).Desc()
func (s SelectStmt) OrderBy(columns ...ColumnElem) SelectStmt {
	s.orderBy = &OrderBySQLClause{columns, "ASC"}
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
func (s SelectStmt) Having(aggregate AggregateSQLClause, op string, value interface{}) SelectStmt {
	s.having = append(s.having, HavingSQLClause{aggregate, op, value})
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

	statement := Statement()

	// select
	columns := []string{}
	for _, c := range s.sel {
		sql, _ := c.Build(dialect)
		if len(s.joins) > 0 {
			columns = append(columns, fmt.Sprintf("%s.%s", dialect.Escape(s.from.Name), sql))
		} else {
			columns = append(columns, sql)
		}

	}
	statement.AddSQLClause(fmt.Sprintf("SELECT %s", strings.Join(columns, ", ")))

	// from
	statement.AddSQLClause(fmt.Sprintf("FROM %s", dialect.Escape(s.from.Name)))

	// joins
	for _, j := range s.joins {
		sql, _ := j.Build(dialect)
		statement.AddSQLClause(sql)
	}

	// where
	if s.where != nil {
		where, bindings := s.where.Build(dialect)
		statement.AddSQLClause(where)
		statement.AddBinding(bindings...)
	}

	// group by
	groupByCols := []string{}
	for _, c := range s.groupBy {
		groupByCols = append(groupByCols, dialect.Escape(c.Name))
	}
	if len(groupByCols) > 0 {
		statement.AddSQLClause(fmt.Sprintf("GROUP BY %s", strings.Join(groupByCols, ", ")))
	}

	// having
	for _, h := range s.having {
		sql, bindings := h.Build(dialect)
		statement.AddSQLClause(sql)
		statement.AddBinding(bindings...)
	}

	// order by
	if s.orderBy != nil {
		sql, _ := s.orderBy.Build(dialect)
		statement.AddSQLClause(sql)
	}

	if (s.offset != nil) && (s.count != nil) {
		statement.AddSQLClause(fmt.Sprintf("LIMIT %d OFFSET %d", *s.count, *s.offset))
	}

	return statement
}

func join(joinType string, fromTable TableElem, table TableElem, fromCol ColumnElem, col ColumnElem) JoinSQLClause {
	return JoinSQLClause{
		joinType,
		fromTable,
		table,
		fromCol,
		col,
	}
}

// JoinSQLClause is the base struct for generating join clauses when using select
// It satisfies SQLClause interface
type JoinSQLClause struct {
	joinType  string
	fromTable TableElem
	table     TableElem
	fromCol   ColumnElem
	col       ColumnElem
}

// Build generates join sql & bindings out of JoinSQLClause struct
func (c JoinSQLClause) Build(dialect Dialect) (string, []interface{}) {

	if (c.fromCol.Name == "") && (c.col.Name == "") {
		return fmt.Sprintf(
			"%s %s",
			c.joinType,
			dialect.Escape(c.table.Name),
		), []interface{}{}
	}

	return fmt.Sprintf(
		"%s %s ON %s.%s = %s.%s",
		c.joinType,
		dialect.Escape(c.table.Name),
		dialect.Escape(c.fromTable.Name),
		dialect.Escape(c.fromCol.Name),
		dialect.Escape(c.table.Name),
		dialect.Escape(c.col.Name),
	), []interface{}{}
}

// OrderBySQLClause is the base struct for generating order by clauses when using select
// It satisfies SQLClause interface
type OrderBySQLClause struct {
	columns []ColumnElem
	t       string
}

// Build generates an order by clause
func (c OrderBySQLClause) Build(dialect Dialect) (string, []interface{}) {
	cols := []string{}
	for _, c := range c.columns {
		cols = append(cols, dialect.Escape(c.Name))
	}

	return fmt.Sprintf("ORDER BY %s %s", strings.Join(cols, ", "), c.t), []interface{}{}
}

// HavingSQLClause is the base struct for generating having clauses when using select
// It satisfies SQLClause interface
type HavingSQLClause struct {
	aggregate AggregateSQLClause
	op        string
	value     interface{}
}

// Build generates having sql & bindings out of HavingSQLClause struct
func (c HavingSQLClause) Build(dialect Dialect) (string, []interface{}) {
	aggSQL, bindings := c.aggregate.Build(dialect)
	bindings = append(bindings, c.value)
	return fmt.Sprintf("HAVING %s %s %s", aggSQL, c.op, dialect.Placeholder()), bindings
}
