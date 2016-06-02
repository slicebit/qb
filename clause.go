package qb

type Clause interface {
	String(adapter Adapter) string
}
