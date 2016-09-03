package qb

// conditional generators, comparator functions

// Like generates a like conditional sql clause
func Like(left Clause, right interface{}) BinaryExpressionClause {
	return BinaryExpression(left, "LIKE", GetClauseFrom(right))
}

// NotIn generates a not in conditional sql clause
func NotIn(left Clause, values ...interface{}) BinaryExpressionClause {
	return NotInClause(left, GetListFrom(values))
}

// NotInClause generates a NOT IN clause
func NotInClause(left Clause, right Clause) BinaryExpressionClause {
	return BinaryExpression(left, "NOT IN", right)
}

// In generates an in conditional sql clause
func In(left Clause, values ...interface{}) BinaryExpressionClause {
	return InClause(left, GetListFrom(values))
}

// InClause generates an in conditional sql clause
func InClause(left Clause, right Clause) BinaryExpressionClause {
	return BinaryExpression(left, "IN", right)
}

// NotEq generates a not equal conditional sql clause
func NotEq(left Clause, right interface{}) BinaryExpressionClause {
	return BinaryExpression(left, "!=", GetClauseFrom(right))
}

// Eq generates a equals conditional sql clause
func Eq(left Clause, right interface{}) BinaryExpressionClause {
	return BinaryExpression(left, "=", GetClauseFrom(right))
}

// Gt generates a greater than conditional sql clause
func Gt(left Clause, right interface{}) BinaryExpressionClause {
	return BinaryExpression(left, ">", GetClauseFrom(right))
}

// Lt generates a less than conditional sql clause
func Lt(left Clause, right interface{}) BinaryExpressionClause {
	return BinaryExpression(left, "<", GetClauseFrom(right))
}

// Gte generates a greater than or equal to conditional sql clause
func Gte(left Clause, right interface{}) BinaryExpressionClause {
	return BinaryExpression(left, ">=", GetClauseFrom(right))
}

// Lte generates a less than or equal to conditional sql clause
func Lte(left Clause, right interface{}) BinaryExpressionClause {
	return BinaryExpression(left, "<=", GetClauseFrom(right))
}

// BinaryExpression generates a condition object to use in update, delete & select statements
func BinaryExpression(left Clause, op string, right Clause) BinaryExpressionClause {
	return BinaryExpressionClause{
		Left:  left,
		Right: right,
		Op:    op,
	}
}

// BinaryExpressionClause is the base struct for any conditional statements in sql clauses
type BinaryExpressionClause struct {
	Left  Clause
	Right Clause
	Op    string
}

// Accept calls the compiler VisitBinary method
func (c BinaryExpressionClause) Accept(context *CompilerContext) string {
	return context.Compiler.VisitBinary(context, c)
}
