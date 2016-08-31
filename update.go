package qb

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
	where     *WhereClause
}

// Accept implements Clause.Accept
func (s UpdateStmt) Accept(context *CompilerContext) string {
	return context.Compiler.VisitUpdate(context, s)
}

// Build generates a statement out of UpdateStmt object
func (s UpdateStmt) Build(dialect Dialect) *Stmt {
	defer dialect.Reset()

	context := NewCompilerContext(dialect)
	statement := Statement()
	statement.AddSQLClause(s.Accept(context))
	statement.AddBinding(context.Binds...)

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
func (s UpdateStmt) Where(clause Clause) UpdateStmt {
	s.where = &WhereClause{clause}
	return s
}
