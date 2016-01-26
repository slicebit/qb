package qbit

import (
	"fmt"
	"strings"
)

func Query() *query {
	return &query{
		clauses:  []string{},
		bindings: []interface{}{},
	}
}

type query struct {
	clauses    []string
	bindings   []interface{}
}

func (q *query) AddClause(clause string) {
	q.clauses = append(q.clauses, clause)
}

func (q *query) AddBinding(bindings ...interface{}) {
	for _, v := range bindings {
		q.bindings = append(q.bindings, v)
	}
}

func (q *query) Clauses() []string {
	return q.clauses
}

func (q *query) Bindings() []interface{} {
	return q.bindings
}
