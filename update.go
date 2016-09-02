package qb

import (
	"fmt"
	"strings"
)

// Update generates an update statement and returns it
// qb.Update(usersTable).
// Values(map[string]interface{}{"id": 1}).
// Where(qb.Eq("id", 5))
func Update(table TableElem) UpdateStmt {
	return UpdateStmt{
		table:     table,
		values:    map[string]interface{}{},
		returning: []ColumnElem{},
	}
}

// UpdateStmt is the base struct for any update statements
type UpdateStmt struct {
	table     TableElem
	values    map[string]interface{}
	returning []ColumnElem
	where     *WhereSQLClause
}

// Build generates a statement out of UpdateStmt object
func (s UpdateStmt) Build(dialect Dialect) *Stmt {
	defer dialect.Reset()

	statement := Statement()
	statement.AddSQLClause(fmt.Sprintf("UPDATE %s", dialect.Escape(s.table.Name)))
	sets := []string{}
	bindings := []interface{}{}
	for k, v := range s.values {
		sets = append(sets, fmt.Sprintf("%s = %s", dialect.Escape(k), dialect.Placeholder()))
		bindings = append(bindings, v)
	}

	if len(sets) > 0 {
		statement.AddSQLClause(fmt.Sprintf("SET %s", strings.Join(sets, ", ")))
	}

	if s.where != nil {
		where, whereBindings := s.where.Build(dialect)
		bindings = append(bindings, whereBindings...)
		statement.AddSQLClause(where)
	}

	returning := []string{}
	for _, c := range s.returning {
		returning = append(returning, dialect.Escape(c.Name))
	}

	if len(returning) > 0 {
		statement.AddSQLClause(fmt.Sprintf("RETURNING %s", strings.Join(returning, ", ")))
	}

	statement.AddBinding(bindings...)

	return statement
}

// Values accepts map[string]interface{} and forms the values map of insert statement
func (s UpdateStmt) Values(values map[string]interface{}) UpdateStmt {
	for k, v := range values {
		s.values[s.table.C(k).Name] = v
	}
	return s
}

// Returning accepts the column names as strings and forms the returning array of insert statement
// NOTE: Please use it in only postgres dialect, otherwise it'll crash
func (s UpdateStmt) Returning(cols ...ColumnElem) UpdateStmt {
	for _, c := range cols {
		s.returning = append(s.returning, c)
	}
	return s
}

// Where adds a where clause to update statement and returns the update statement
func (s UpdateStmt) Where(clause SQLClause) UpdateStmt {
	s.where = &WhereSQLClause{clause}
	return s
}
