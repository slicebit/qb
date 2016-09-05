package qb

// conditional generators, comparator functions

// Like generates a like conditional sql clause
func Like(col ColumnElem, pattern string) Conditional {
	return Condition(col, "LIKE", pattern)
}

// NotIn generates a not in conditional sql clause
func NotIn(col ColumnElem, values ...interface{}) Conditional {
	return Condition(col, "NOT IN", values...)
}

// In generates an in conditional sql clause
func In(col ColumnElem, values ...interface{}) Conditional {
	return Condition(col, "IN", values...)
}

// NotEq generates a not equal conditional sql clause
func NotEq(col ColumnElem, value interface{}) Conditional {
	return Condition(col, "!=", value)
}

// Eq generates a equals conditional sql clause
func Eq(col ColumnElem, value interface{}) Conditional {
	return Condition(col, "=", value)
}

// Gt generates a greater than conditional sql clause
func Gt(col ColumnElem, value interface{}) Conditional {
	return Condition(col, ">", value)
}

// Lt generates a less than conditional sql clause
func Lt(col ColumnElem, value interface{}) Conditional {
	return Condition(col, "<", value)
}

// Gte generates a greater than or equal to conditional sql clause
func Gte(col ColumnElem, value interface{}) Conditional {
	return Condition(col, ">=", value)
}

// Lte generates a less than or equal to conditional sql clause
func Lte(col ColumnElem, value interface{}) Conditional {
	return Condition(col, "<=", value)
}

// Condition generates a condition object to use in update, delete & select statements
func Condition(col ColumnElem, op string, values ...interface{}) Conditional {
	return Conditional{col, values, op}
}

// Conditional is the base struct for any conditional statements in sql clauses
type Conditional struct {
	Col    ColumnElem
	Values []interface{}
	Op     string
}

// Accept compiles the conditional element
func (c Conditional) Accept(context *CompilerContext) string {
	return context.Compiler.VisitCondition(context, c)
}
