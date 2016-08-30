package qb

import (
	"fmt"
	"strings"
)

// Delete generates a delete statement and returns it for chaining
// qb.Delete(usersTable).Where(qb.Eq("id", 5))
func Delete(table TableElem) DeleteStmt {
	return DeleteStmt{
		table:     table,
		returning: []ColumnElem{},
	}
}

// DeleteStmt is the base struct for building delete queries
type DeleteStmt struct {
	table     TableElem
	where     *WhereSQLClause
	returning []ColumnElem
}

// Where adds a where clause to the current delete statement
func (s DeleteStmt) Where(clause SQLClause) DeleteStmt {
	s.where = &WhereSQLClause{clause}
	return s
}

// Returning accepts the column names as strings and forms the returning array of insert statement
// NOTE: Please use it in only postgres dialect, otherwise it'll crash
func (s DeleteStmt) Returning(cols ...ColumnElem) DeleteStmt {
	s.returning = append(s.returning, cols...)
	return s
}

// Build generates a statement out of DeleteStmt object
func (s DeleteStmt) Build(dialect Dialect) *Stmt {
	defer dialect.Reset()

	statement := Statement()
	statement.AddSQLClause(fmt.Sprintf("DELETE FROM %s", dialect.Escape(s.table.Name)))
	if s.where != nil {
		where, whereBindings := s.where.Build(dialect)
		statement.AddSQLClause(where)
		statement.AddBinding(whereBindings...)
	}

	returning := []string{}
	for _, c := range s.returning {
		returning = append(returning, dialect.Escape(c.Name))
	}

	if len(returning) > 0 {
		statement.AddSQLClause(fmt.Sprintf("RETURNING %s", strings.Join(returning, ", ")))
	}

	return statement
}
