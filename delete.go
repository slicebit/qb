package qb

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
	where     *WhereClause
	returning []ColumnElem
}

// Where adds a where clause to the current delete statement
func (s DeleteStmt) Where(clause Clause) DeleteStmt {
	s.where = &WhereClause{clause}
	return s
}

// Returning accepts the column names as strings and forms the returning array of insert statement
// NOTE: Please use it in only postgres dialect, otherwise it'll crash
func (s DeleteStmt) Returning(cols ...ColumnElem) DeleteStmt {
	s.returning = append(s.returning, cols...)
	return s
}

// Accept implements Clause.Accept
func (s DeleteStmt) Accept(context *CompilerContext) string {
	return context.Compiler.VisitDelete(context, s)
}

// Build generates a statement out of DeleteStmt object
func (s DeleteStmt) Build(dialect Dialect) *Stmt {
	defer dialect.Reset()

	context := NewCompilerContext(dialect)
	statement := Statement()
	statement.AddSQLClause(s.Accept(context))
	statement.AddBinding(context.Binds...)

	return statement
}
