package qb

import (
	"fmt"
	"strings"
)

// Upsert generates an insert ... on (duplicate key/conflict) update statement
func Upsert(table TableElem) UpsertStmt {
	return UpsertStmt{
		table:     table,
		values:    map[string]interface{}{},
		returning: []ColumnElem{},
	}
}

// UpsertStmt is the base struct for any insert ... on conflict/duplicate key ... update ... statements
type UpsertStmt struct {
	table     TableElem
	values    map[string]interface{}
	returning []ColumnElem
}

// Values accepts map[string]interface{} and forms the values map of insert statement
func (s UpsertStmt) Values(vals map[string]interface{}) UpsertStmt {
	for k, v := range vals {
		s.values[k] = v
	}
	return s
}

// Returning accepts the column names as strings and forms the returning array of insert statement
// NOTE: Please use it in only postgres dialect, otherwise it'll crash
func (s UpsertStmt) Returning(cols ...string) UpsertStmt {
	for _, c := range cols {
		s.returning = append(s.returning, s.table.C(c))
	}
	return s
}

// Build generates a statement out of UpdateStmt object
// NOTE: It generates different statements for each driver
// For sqlite, it generates REPLACE INTO ... VALUES ...
// For mysql, it generates INSERT INTO ... VALUES ... ON DUPLICATE KEY UPDATE ...
// For postgres, it generates INSERT INTO ... VALUES ... ON CONFLICT(...) DO UPDATE SET ...
func (s UpsertStmt) Build(dialect Dialect) *Stmt {
	statement := Statement()

	colNames := []string{}
	values := []string{}
	for k, v := range s.values {
		colNames = append(colNames, dialect.Escape(k))
		statement.AddBinding(v)
		values = append(values, dialect.Placeholder())
	}

	updates := []string{}
	for k, v := range s.values {
		updates = append(updates, fmt.Sprintf("%s = %s", dialect.Escape(k), dialect.Placeholder()))
		statement.AddBinding(v)
	}

	switch dialect.Driver() {
	case "mysql":
		statement.AddClause(fmt.Sprintf("INSERT INTO %s(%s)", dialect.Escape(s.table.Name), strings.Join(colNames, ", ")))
		statement.AddClause(fmt.Sprintf("VALUES(%s)", strings.Join(values, ", ")))
		statement.AddClause(fmt.Sprintf("ON DUPLICATE KEY UPDATE %s", strings.Join(updates, ", ")))
		break
	case "postgres":
		statement.AddClause(fmt.Sprintf("INSERT INTO %s(%s)", dialect.Escape(s.table.Name), strings.Join(colNames, ", ")))
		statement.AddClause(fmt.Sprintf("VALUES(%s)", strings.Join(values, ", ")))
		uniqueCols := []string{}
		for _, c := range s.table.PrimaryCols() {
			uniqueCols = append(uniqueCols, dialect.Escape(c.Name))
		}
		statement.AddClause(fmt.Sprintf("ON CONFLICT (%s) DO UPDATE SET %s", strings.Join(uniqueCols, ", "), strings.Join(updates, ", ")))
		returning := []string{}
		for _, r := range s.returning {
			returning = append(returning, dialect.Escape(r.Name))
		}
		if len(s.returning) > 0 {
			statement.AddClause(fmt.Sprintf("RETURNING %s", strings.Join(returning, ", ")))
		}
		break
	case "sqlite3":
		statement.AddClause(fmt.Sprintf("REPLACE INTO %s(%s)", dialect.Escape(s.table.Name), strings.Join(colNames, ", ")))
		statement.AddClause(fmt.Sprintf("VALUES(%s)", strings.Join(values, ", ")))
		break
	}
	return statement
}
