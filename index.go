package qb

import (
	"fmt"
	"strings"
)

// CompositeIndex is the struct definition when building composite indices for any struct that will be mapped into a table
type CompositeIndex string

// Index generates an index clause given table and columns as params
func Index(table string, cols ...string) IndexElem {
	return IndexElem{
		Table:   table,
		Name:    fmt.Sprintf("i_%s", strings.Join(cols, "_")),
		Columns: cols,
	}
}

// IndexElem is the definition of any index elements for a table
type IndexElem struct {
	Table   string
	Name    string
	Columns []string
}

// String returns the index element as an sql clause
func (i IndexElem) String(dialect Dialect) string {
	return fmt.Sprintf("CREATE INDEX %s ON %s(%s);", dialect.Escape(i.Name), dialect.Escape(i.Table), strings.Join(dialect.EscapeAll(i.Columns), ", "))
}
