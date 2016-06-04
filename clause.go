package qb

// Clause is the common interface for any table clause in qb.
// NOTE: Do not mix this with builder's internal clauses array.
type Clause interface {
	String(adapter Adapter) string
}
