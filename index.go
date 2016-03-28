package qb

// CompositeIndex is the container for multiple column indices
type CompositeIndex struct{}

// NewIndex generates a new Index struct given column array
func NewIndex(table string, name string, columns ...string) *Index {
	return &Index{
		table:   table,
		columns: columns,
		name:    name,
	}
}

// Index is the struct for generating table indices
type Index struct {
	table   string
	columns []string
	name    string
}

// Table returns the table property of index
func (i *Index) Table() string {
	return i.table
}

// Columns returns the columns property of index
func (i *Index) Columns() []string {
	return i.columns
}

// Name returns the name property of index
func (i *Index) Name() string {
	return i.name
}