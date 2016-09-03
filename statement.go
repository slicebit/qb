package qb

import (
	"fmt"
	"strings"
)

const defaultDelimiter = "\n"

// Statement creates a new query and returns its pointer
func Statement() *Stmt {
	return &Stmt{
		clauses:      []string{},
		bindings:     []interface{}{},
		delimiter:    defaultDelimiter,
		bindingIndex: 0,
	}
}

// Stmt is the base abstraction for all sql queries
type Stmt struct {
	clauses      []string
	bindings     []interface{}
	delimiter    string
	bindingIndex int
}

// Text is for executing raw sql
// It parses the sql and generates clauses from
func (s *Stmt) Text(sql string) {
	sql = strings.Replace(sql, ";", "", -1)
	sql = strings.Replace(sql, "\t", "", -1)
	sql = strings.Trim(sql, "\n")
	clauses := strings.Split(sql, "\n")
	for _, c := range clauses {
		s.clauses = append(s.clauses, c)
	}
}

// SetDelimiter sets the delimiter of query
func (s *Stmt) SetDelimiter(delimiter string) {
	s.delimiter = delimiter
}

// AddSQLClause appends a new clause to current query
func (s *Stmt) AddSQLClause(clause string) {
	s.clauses = append(s.clauses, clause)
}

// AddBinding appends a new binding to current query
func (s *Stmt) AddBinding(bindings ...interface{}) {
	for _, v := range bindings {
		s.bindings = append(s.bindings, v)
	}
}

// SQLClauses returns all clauses of current query
func (s *Stmt) SQLClauses() []string {
	return s.clauses
}

// Bindings returns all bindings of current query
func (s *Stmt) Bindings() []interface{} {
	return s.bindings
}

// SQL returns the query struct sql statement
func (s *Stmt) SQL() string {
	if len(s.clauses) > 0 {
		sql := fmt.Sprintf("%s;", strings.Join(s.clauses, s.delimiter))
		return sql
	}

	return ""
}
