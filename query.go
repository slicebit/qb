package qb

import (
	"fmt"
	"strings"
)

const defaultDelimiter = "\n"

// NewQuery creates a new query and returns its pointer
func NewQuery() *Query {
	return &Query{
		clauses:   []string{},
		bindings:  []interface{}{},
		delimiter: defaultDelimiter,
		bindingIndex: 0,
	}
}

// Query is the base abstraction for sql queries
type Query struct {
	clauses      []string
	bindings     []interface{}
	delimiter    string
	bindingIndex int
}

// SetDelimiter sets the delimiter of query
func (q *Query) SetDelimiter(delimiter string) {
	q.delimiter = delimiter
}

// AddClause appends a new clause to current query
func (q *Query) AddClause(clause string) {
	q.clauses = append(q.clauses, clause)
}

// AddBinding appends a new binding to current query
func (q *Query) AddBinding(bindings ...interface{}) {
	for _, v := range bindings {
		q.bindings = append(q.bindings, v)
	}
}

// Clauses returns all clauses of current query
func (q *Query) Clauses() []string {
	return q.clauses
}

// Bindings returns all bindings of current query
func (q *Query) Bindings() []interface{} {
	return q.bindings
}

// Placeholder generates a placeholder for binding.
// If the driver is postgres generates "$i" where i is incremental integer
// Otherwise it generates "?"
func (q *Query) placeholder(driver string) string {
	if driver == "postgres" {
		q.bindingIndex++
		return fmt.Sprintf("$%d", q.bindingIndex)
	}
	return "?"
}

// Placeholders generates multiple placeholders
// This is for builder to be able to temporarily put question marks
// The driver is unknown in builder
func (q *Query) QuestionMarks(values ...interface{}) []string {
	marks := []string{}
	for _ = range values {
		marks = append(marks, "?")
	}
	return marks
}

// SQL returns the query struct sql statement
func (q *Query) SQL(driver string) string {

	if len(q.clauses) > 0 {
		sql := fmt.Sprintf("%s;", strings.Join(q.clauses, q.delimiter))
		if driver == "postgres" {
			count := strings.Count(sql, "?")
			for i := 0; i < count; i++ {
				sql = strings.Replace(sql, "?", q.placeholder(driver), 1)
			}
		}
		return sql
	}

	return ""
}
