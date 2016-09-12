package qb

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
func (s UpsertStmt) Values(values map[string]interface{}) UpsertStmt {
	for k, v := range values {
		s.values[k] = v
	}
	return s
}

// Returning accepts the column names as strings and forms the returning array of insert statement
// NOTE: Please use it in only postgres dialect, otherwise it'll crash
func (s UpsertStmt) Returning(cols ...ColumnElem) UpsertStmt {
	for _, c := range cols {
		s.returning = append(s.returning, c)
	}
	return s
}

// Accept calls the compiler VisitUpsert function
func (s UpsertStmt) Accept(context *CompilerContext) string {
	return context.Compiler.VisitUpsert(context, s)
}

// Build generates a statement out of UpdateStmt object
func (s UpsertStmt) Build(dialect Dialect) *Stmt {
	context := NewCompilerContext(dialect)
	statement := Statement()
	statement.AddSQLClause(s.Accept(context))
	statement.AddBinding(context.Binds...)

	return statement
}
