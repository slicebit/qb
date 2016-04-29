package qb

import (
	"fmt"
	"strings"
)

// NewTable generates a new table pointer given table name, column and table constraints
func NewTable(builder *Builder, name string, columns []Column) *Table {
	return &Table{
		name:        name,
		columns:     columns,
		primaryCols: []string{},
		refs:        []ref{},
		builder:     builder,
		indices:     []*Index{},
	}
}

// Table is the base abstraction for any sql table
type Table struct {
	name        string
	columns     []Column
	primaryCols []string
	refs        []ref
	builder     *Builder
	indices     []*Index
}

// Column returns the table column given column name
func (t *Table) Column(name string) (Column, error) {
	for _, c := range t.columns {
		if c.Name == name {
			return c, nil
		}
	}

	return Column{}, fmt.Errorf("Invalid column %s", name)
}

// Name returns the table name
func (t *Table) Name() string {
	return t.name
}

// SQL generates create table syntax of table
func (t *Table) SQL() string {
	cols := []string{}
	for _, v := range t.columns {
		cols = append(cols, v.SQL(t.builder.Adapter()))
	}

	constraints := []string{}

	// build primary key constraints using primaryCols
	if len(t.primaryCols) > 0 {
		for k, col := range t.primaryCols {
			t.primaryCols[k] = t.builder.Adapter().Escape(col)
		}
		constraints = append(constraints, fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(t.primaryCols, ", ")))
	}

	// build foreign key constraints using refCols
	for _, ref := range t.refs {
		for k := range ref.cols {
			ref.cols[k] = t.builder.Adapter().Escape(ref.cols[k])
			ref.refCols[k] = t.builder.Adapter().Escape(ref.refCols[k])
		}
		constraints = append(constraints, fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(%s)", strings.Join(ref.cols, ", "), t.builder.Adapter().Escape(ref.refTable), strings.Join(ref.refCols, ", ")))
	}

	tableSQL := t.builder.CreateTable(t.name, cols, constraints).Query().SQL()

	indexSqls := []string{}
	for _, index := range t.indices {
		q := t.builder.CreateIndex(index.Name(), index.Table(), index.Columns()...).Query()
		indexSqls = append(indexSqls, q.SQL())
	}

	sqls := []string{tableSQL}
	sqls = append(sqls, indexSqls...)

	return strings.Join(sqls, "\n")
}

// AddColumn appends a new column to current table
func (t *Table) AddColumn(column Column) {
	t.columns = append(t.columns, column)
}

// AddPrimary appends a primary column that will be lazily built as a primary key constraint
func (t *Table) AddPrimary(col string) {
	t.primaryCols = append(t.primaryCols, col)
}

type ref struct {
	cols     []string
	refTable string
	refCols  []string
}

// AddRef appends a new reference struct that will be lazily built as a foreign key constraint
func (t *Table) AddRef(col string, refTable string, refCol string) {
	if len(t.refs) > 0 {
		for k, ref := range t.refs {
			if refTable == ref.refTable {
				t.refs[k].cols = append(t.refs[k].cols, col)
				t.refs[k].refCols = append(t.refs[k].refCols, refCol)
				return
			}
		}
	}

	r := ref{[]string{}, refTable, []string{}}
	r.cols = append(r.cols, col)
	r.refCols = append(r.refCols, refCol)
	t.refs = append(t.refs, r)
}

// AddIndex appends a new index that will be lazily created in SQL() function
func (t *Table) AddIndex(columns ...string) {
	indexName := fmt.Sprintf("index_%s", strings.Join(columns, "_"))
	t.indices = append(t.indices, NewIndex(t.name, indexName, columns...))
}
