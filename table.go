package qb

import (
	"fmt"
	"strings"
)

// Table generates table struct given name and clauses
func Table(name string, clauses ...TableSQLClause) TableElem {
	table := TableElem{
		Name:                  name,
		Columns:               map[string]ColumnElem{},
		ForeignKeyConstraints: ForeignKeyConstraints{},
		Indices:               []IndexElem{},
	}

	var pkeyCols []ColumnElem

	for _, clause := range clauses {
		switch clause.(type) {
		case ColumnElem:
			col := clause.(ColumnElem)
			if col.Options.PrimaryKey {
				pkeyCols = append(pkeyCols, col)
			}
			col.Table = name
			table.Columns[col.Name] = col
			break
		case PrimaryKeyConstraint:
			table.PrimaryKeyConstraint = clause.(PrimaryKeyConstraint)
			break
		//case ForeignKeyConstraints:
		//	table.ForeignKeyConstraints.FKeys = append(
		//		table.ForeignKeyConstraints.FKeys,
		//		clause.(ForeignKeyConstraints).FKeys...)
		case ForeignKeyConstraint:
			table.ForeignKeyConstraints.FKeys = append(
				table.ForeignKeyConstraints.FKeys,
				clause.(ForeignKeyConstraint),
			)
			break
		case UniqueKeyConstraint:
			table.UniqueKeyConstraint = clause.(UniqueKeyConstraint)
			break
		case IndexElem:
			table.Indices = append(table.Indices, clause.(IndexElem))
			break
		}
	}

	if len(pkeyCols) > 0 && table.PrimaryKeyConstraint.Columns != nil {
		panic(fmt.Sprintf("Table %s has both 'PrimaryKey()' columns (%#v) and a PrimaryKeyConstraint. Only only should be set", name, pkeyCols))
	}
	if len(pkeyCols) > 1 {
		var pkeyNames []string
		for _, col := range pkeyCols {
			pkeyNames = append(pkeyNames, col.Name)
		}
		table.PrimaryKeyConstraint = PrimaryKey(pkeyNames...)
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

// DefaultName returns the name of the table
func (t TableElem) DefaultName() string {
	return t.Name
}

// All returns all columns of table as a column slice
func (t TableElem) All() []Clause {
	cols := []Clause{}
	for _, v := range t.Columns {
		cols = append(cols, v)
	}
	return cols
}

// ColumnList columns of the table
func (t TableElem) ColumnList() []ColumnElem {
	cols := []ColumnElem{}
	for _, v := range t.Columns {
		cols = append(cols, v)
	}
	return cols
}

// Index appends an IndexElem to current table without giving table name
func (t TableElem) Index(cols ...string) TableElem {
	t.Indices = append(t.Indices, Index(t.Name, cols...))
	return t
}

// Create generates create table syntax and returns it as a query struct
func (t TableElem) Create(dialect Dialect) string {
	statement := Statement()
	statement.AddSQLClause(fmt.Sprintf("CREATE TABLE %s (", dialect.Escape(t.Name)))

	colClauses := []string{}
	for _, col := range t.Columns {
		colClauses = append(colClauses, fmt.Sprintf("\t%s", col.String(dialect)))
	}

	if len(t.PrimaryKeyConstraint.Columns) > 0 {
		colClauses = append(colClauses, fmt.Sprintf("\t%s", t.PrimaryKeyConstraint.String(dialect)))
	}

	if len(t.ForeignKeyConstraints.FKeys) > 0 {
		colClauses = append(colClauses, t.ForeignKeyConstraints.String(dialect))
	}

	if t.UniqueKeyConstraint.name != "" {
		colClauses = append(colClauses, fmt.Sprintf("\t%s", t.UniqueKeyConstraint.String(dialect)))
	}

	statement.AddSQLClause(strings.Join(colClauses, ",\n"))

	statement.AddSQLClause(")")

	ddl := statement.SQL()

	indexSqls := []string{}
	for _, index := range t.Indices {
		iSQLClause := index.String(dialect)
		indexSqls = append(indexSqls, iSQLClause)
	}

	sqls := []string{ddl}
	sqls = append(sqls, indexSqls...)

	return strings.Join(sqls, "\n")
}

// Build generates a Statement object out of table ddl
func (t TableElem) Build(dialect Dialect) *Stmt {
	sql := t.Create(dialect)
	statement := Statement()
	statement.AddSQLClause(strings.Trim(sql, ";")) // TODO: Remove this ugly hack
	return statement
}

// PrimaryCols returns the columns that are primary key to the table
func (t TableElem) PrimaryCols() []ColumnElem {
	primaryCols := []ColumnElem{}
	pkCols := t.PrimaryKeyConstraint.Columns
	for _, pkCol := range pkCols {
		primaryCols = append(primaryCols, t.C(pkCol))
	}
	return primaryCols
}

// Drop generates drop table syntax and returns it as a query struct
func (t TableElem) Drop(dialect Dialect) string {
	stmt := Statement()
	stmt.AddSQLClause(fmt.Sprintf("DROP TABLE %s", dialect.Escape(t.Name)))
	return stmt.SQL()
}

// C returns the column name given col
func (t TableElem) C(name string) ColumnElem {
	return t.Columns[name]
}

// query starters

// Insert starts an insert statement by setting the table parameter
func (t TableElem) Insert() InsertStmt {
	return Insert(t)
}

// Update starts an update statement by setting the table parameter
func (t TableElem) Update() UpdateStmt {
	return Update(t)
}

// Delete starts a delete statement by setting the table parameter
func (t TableElem) Delete() DeleteStmt {
	return Delete(t)
}

// Upsert starts an upsert statement by setting the table parameter
func (t TableElem) Upsert() UpsertStmt {
	return Upsert(t)
}

// Select starts a select statement by setting from table
func (t TableElem) Select(clauses ...Clause) SelectStmt {
	return Select(clauses...).From(t)
}

// Accept implements Clause.Accept
func (t TableElem) Accept(context *CompilerContext) string {
	return context.Compiler.VisitTable(context, t)
}
