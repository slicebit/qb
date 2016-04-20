package qb

// Type is the base abstraction for any sql column type
type Type struct {
	SQL string
}

// NewType creates a new type pointer given sql as string
func NewType(sql string) *Type {
	return &Type{sql}
}
