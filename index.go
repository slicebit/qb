package qb

import (
	"fmt"
	"strings"
)

type CompositeIndex string

func Index(table string, cols ...string) IndexElem {
	return IndexElem{
		Table:   table,
		Name:    fmt.Sprintf("i_%s", strings.Join(cols, "_")),
		Columns: cols,
	}
}

type IndexElem struct {
	Table   string
	Name    string
	Columns []string
}

func (i IndexElem) String(adapter Adapter) string {
	return fmt.Sprintf("CREATE INDEX %s ON %s(%s);", adapter.Escape(i.Name), adapter.Escape(i.Table), strings.Join(adapter.EscapeAll(i.Columns), ", "))
}
