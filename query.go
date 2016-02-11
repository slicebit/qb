package qbit

// NewQuery creates a new query and returns its pointer
func NewQuery() *Query {
	return &Query{
		clauses:  []string{},
		bindings: []interface{}{},
	}
}

// Query is the base abstraction for sql queries
type Query struct {
	clauses  []string
	bindings []interface{}
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
