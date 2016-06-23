package qb

import (
	"fmt"
	"strings"
)

// Insert generates an insert statement and returns it
// Insert(usersTable).Values(map[string]interface{}{"id": 1})
func Insert(table TableElem) InsertStmt {
	return InsertStmt{
		table:     table,
		values:    map[string]interface{}{},
		returning: []ColumnElem{},
	}
}

// InsertStmt is the base struct for any insert statements
type InsertStmt struct {
	table     TableElem
	values    map[string]interface{}
	returning []ColumnElem
}

// Values accepts map[string]interface{} and forms the values map of insert statement
func (s InsertStmt) Values(values map[string]interface{}) InsertStmt {
	for k, v := range values {
		s.values[k] = v
	}
	return s
}

// Returning accepts the column names as strings and forms the returning array of insert statement
// NOTE: Please use it in only postgres dialect, otherwise it'll crash
func (s InsertStmt) Returning(cols ...ColumnElem) InsertStmt {
	for _, c := range cols {
		s.returning = append(s.returning, c)
	}
	return s
}

// Build generates a statement out of InsertStmt object
func (s InsertStmt) Build(dialect Dialect) *Stmt {
	statement := Statement()
	colNames := []string{}
	values := []string{}
	for k, v := range s.values {
		colNames = append(colNames, dialect.Escape(k))
		statement.AddBinding(v)
		values = append(values, dialect.Placeholder())
	}
	statement.AddClause(fmt.Sprintf("INSERT INTO %s(%s)", dialect.Escape(s.table.Name), strings.Join(colNames, ", ")))
	statement.AddClause(fmt.Sprintf("VALUES(%s)", strings.Join(values, ", ")))

	returning := []string{}
	for _, r := range s.returning {
		returning = append(returning, dialect.Escape(r.Name))
	}
	if len(s.returning) > 0 {
		statement.AddClause(fmt.Sprintf("RETURNING %s", strings.Join(returning, ", ")))
	}
	return statement
}