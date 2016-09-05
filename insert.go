package qb

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

// Accept implements Clause.Accept
func (s InsertStmt) Accept(context *CompilerContext) string {
	return context.Compiler.VisitInsert(context, s)
}

// Build generates a statement out of InsertStmt object
func (s InsertStmt) Build(dialect Dialect) *Stmt {
	defer dialect.Reset()

	statement := Statement()
	context := NewCompilerContext(dialect)
	statement.AddSQLClause(s.Accept(context))
	statement.AddBinding(context.Binds...)

	return statement
}
