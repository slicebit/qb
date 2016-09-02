package qb

import (
	"reflect"
)

// SQLText returns a raw SQL clause
func SQLText(text string) TextClause {
	return TextClause{Text: text}
}

// TextClause is a raw SQL clause
type TextClause struct {
	Text string
}

// Accept calls the compiler VisitText method
func (c TextClause) Accept(context *CompilerContext) string {
	return context.Compiler.VisitText(context, c)
}

// List returns a list-of-clauses clause
func List(clauses ...Clause) ListClause {
	return ListClause{
		Clauses: clauses,
	}
}

// ListClause is a list of clause elements (for IN operator for example)
type ListClause struct {
	Clauses []Clause
}

// Accept calls the compiler VisitList method
func (c ListClause) Accept(context *CompilerContext) string {
	return context.Compiler.VisitList(context, c)
}

// Bind a value
func Bind(value interface{}) BindClause {
	return BindClause{
		Value: value,
	}
}

// BindClause binds a value to a placeholder
type BindClause struct {
	Value interface{}
}

// Accept calls the compiler VisitBind method
func (c BindClause) Accept(context *CompilerContext) string {
	return context.Compiler.VisitBind(context, c)
}

// GetClauseFrom returns the value if already a Clause, or make one
// if it is a scalar value
func GetClauseFrom(value interface{}) Clause {
	if clause, ok := value.(Clause); ok {
		return clause
	}
	// For now we assume any non-clause is a Value:
	return Bind(value)
}

// GetListFrom returns a list clause from any list
//
// if 'any' is a Clause, it is wrapped into a ListClause unless it is
// a ListClause, in which case it is returned as-is
//
// if 'any' is not a list, it is wrapped as a clause, and returned a the
// only item of a new ListClause
//
// if 'any' is a list:
//   if len>1 each item is converted to a clause if not already one.
//   if len=1, the function is recursively called on it. This allow to use
//   this function on a variadic arg where a single list clause is passed.
//
func GetListFrom(any interface{}) Clause {
	if clause, ok := any.(Clause); ok {
		if _, isList := clause.(ListClause); isList {
			return clause
		}
		return List(clause)
	}

	v := reflect.ValueOf(any)
	switch kind := v.Kind(); kind {
	case reflect.Array, reflect.Slice:
		if v.Len() == 1 {
			return GetListFrom(v.Index(0).Interface())
		}
		var clauses []Clause
		for i := 0; i != v.Len(); i++ {
			clauses = append(clauses, GetClauseFrom(v.Index(i).Interface()))
		}
		return List(clauses...)
	default:
		return List(GetClauseFrom(any))
	}
}
