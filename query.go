package qbit

import (
	"fmt"
	"strings"
)

func Query() *Query {
	return &Query{
		clauses:  []string{},
		bindings: []interface{}{},
	}
}

type Query struct {
	clauses    []string
	bindings   []interface{}
}

func (q *Query) AddClause(clause string) {
	q.clauses = append(q.clauses, clause)
}

func (q *Query) AddBinding(bindings ...interface{}) {
	for _, v := range bindings {
		q.bindings = append(q.bindings, v)
	}
}

func (q *Query) Clauses() []string {
	return q.clauses
}

func (q *Query) Bindings() []interface{} {
	return q.bindings
}
