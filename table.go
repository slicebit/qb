package qb

import (
	"fmt"
	"strings"
)

// Table generates table struct given name and clauses
func Table(name string, clauses ...Clause) TableElem {
	table := TableElem{
		Name:                  name,
		Columns:               map[string]ColumnElem{},
		ForeignKeyConstraints: ForeignKeyConstraints{},
		Indices:               []IndexElem{},
	}

	for _, clause := range clauses {
		switch clause.(type) {
		case ColumnElem:
			col := clause.(ColumnElem)
			table.Columns[col.Name] = col
			break
		case PrimaryKeyConstraint:
			table.PrimaryKeyConstraint = clause.(PrimaryKeyConstraint)
			break
		case ForeignKeyConstraints:
			table.ForeignKeyConstraints = clause.(ForeignKeyConstraints)
			break
		case UniqueKeyConstraint:
			table.UniqueKeyConstraint = clause.(UniqueKeyConstraint)
			break
		case IndexElem:
			table.Indices = append(table.Indices, clause.(IndexElem))
			break
		}
	}

	return table
}

// TableElem is the definition of any sql table
type TableElem struct {
	Name                  string
	Columns               map[string]ColumnElem
	PrimaryKeyConstraint  PrimaryKeyConstraint
	ForeignKeyConstraints ForeignKeyConstraints
	UniqueKeyConstraint   UniqueKeyConstraint
	Indices               []IndexElem
}

// Index appends an IndexElem to current table without giving table name
func (t TableElem) Index(cols ...string) TableElem {
	t.Indices = append(t.Indices, Index(t.Name, cols...))
	return t
}

// Create generates create table syntax and returns it as a query struct
func (t TableElem) Create(adapter Adapter) string {
	query := Query()
	query.AddClause(fmt.Sprintf("CREATE TABLE %s (", adapter.Escape(t.Name)))

	colClauses := []string{}
	for _, col := range t.Columns {
		colClauses = append(colClauses, fmt.Sprintf("\t%s", col.String(adapter)))
	}

	if len(t.PrimaryKeyConstraint.Columns) > 0 {
		colClauses = append(colClauses, fmt.Sprintf("\t%s", t.PrimaryKeyConstraint.String(adapter)))
	}

	if len(t.ForeignKeyConstraints.Refs) > 0 {
		colClauses = append(colClauses, t.ForeignKeyConstraints.String(adapter))
	}

	if t.UniqueKeyConstraint.name != "" {
		colClauses = append(colClauses, fmt.Sprintf("\t%s", t.UniqueKeyConstraint.String(adapter)))
	}

	query.AddClause(strings.Join(colClauses, ",\n"))

	query.AddClause(")")

	ddl := query.SQL()

	indexSqls := []string{}
	for _, index := range t.Indices {
		iClause := index.String(adapter)
		indexSqls = append(indexSqls, iClause)
	}

	sqls := []string{ddl}
	sqls = append(sqls, indexSqls...)

	return strings.Join(sqls, "\n")
}

// Drop generates drop table syntax and returns it as a query struct
func (t TableElem) Drop(adapter Adapter) string {
	query := Query()
	query.AddClause(fmt.Sprintf("DROP TABLE %s", adapter.Escape(t.Name)))
	return query.SQL()
}

// C returns the column name given col
func (t TableElem) C(name string) ColumnElem {
	return t.Columns[name]
}
