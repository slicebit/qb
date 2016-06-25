package qb

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
)

// New generates a new Session given engine and returns session pointer
func New(driver string, dsn string) (*Session, error) {
	engine, err := NewEngine(driver, dsn)
	if err != nil {
		return nil, err
	}

	dialect := NewDialect(driver)

	engine.SetDialect(dialect)

	return &Session{
		statements: []*Stmt{},
		engine:     engine,
		dialect:    dialect,
		mapper:     Mapper(dialect),
		metadata:   MetaData(dialect),
		mutex:      &sync.Mutex{},
	}, nil
}

// Session is the composition of engine connection & orm mappings
type Session struct {
	builder    Builder
	filters    []Clause
	statements []*Stmt
	engine     *Engine
	mapper     MapperElem
	metadata   *MetaDataElem
	dialect    Dialect
	tx         *sql.Tx
	mutex      *sync.Mutex
}

func (s *Session) add(statement *Stmt) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var err error
	if s.tx == nil {
		s.tx, err = s.engine.DB().Begin()
		s.statements = []*Stmt{}
		if err != nil {
			panic(err)
		}
	}
	s.statements = append(s.statements, statement)
}

// Metadata Wrappers

// T returns the table given name as string
// It is for Query() function parameter generation
func (s *Session) T(name string) TableElem {
	return s.metadata.Table(name)
}

// AddTable adds a model to metadata that is mapped into table object
func (s *Session) AddTable(model interface{}) {
	s.metadata.Add(model)
}

// CreateAll creates all tables that are registered to metadata
func (s *Session) CreateAll() error {
	return s.metadata.CreateAll(s.engine)
}

// DropAll drops all tables that are registered to metadata
func (s *Session) DropAll() error {
	return s.metadata.DropAll(s.engine)
}

// Engine wrappers

// Engine returns the current sqlx wrapped engine
func (s *Session) Engine() *Engine {
	return s.engine
}

// Close closes engine db (sqlx) connection
func (s *Session) Close() {
	s.engine.DB().Close()
}

// AddStatement adds a statement given the query pointer retrieved from Build() function
func (s *Session) AddStatement(statement *Stmt) {
	s.add(statement)
}

// Dialect returns the current dialect of session
func (s *Session) Dialect() Dialect {
	return s.dialect
}

// Metadata returns the metadata of session
func (s *Session) Metadata() *MetaDataElem {
	return s.metadata
}

// Session Api

// Delete adds a single delete statement to the session
func (s *Session) Delete(model interface{}) {
	kv := s.mapper.ToMap(model, false)

	tableName := s.mapper.ModelName(model)

	d := Delete(s.metadata.Table(tableName))
	conditions := []Clause{}
	for k, v := range kv {
		conditions = append(conditions, Eq(s.metadata.Table(tableName).C(k), v))
	}

	stmt := d.Where(And(conditions...)).Build(s.dialect)
	s.add(stmt)
}

// Add adds a single model to the session. The query must be insert or update
func (s *Session) Add(model interface{}) {
	m := s.mapper.ToMap(model, false)
	tableName := s.mapper.ModelName(model)
	ups := Upsert(s.metadata.Table(tableName)).Values(m)
	statement := ups.Build(s.dialect)

	s.add(statement)
}

// AddAll adds multiple models an adds an insert statement to current queries
func (s *Session) AddAll(models ...interface{}) {
	for _, m := range models {
		s.Add(m)
	}
}

// Commit commits the current transaction with queries
func (s *Session) Commit() error {
	for _, statement := range s.statements {
		s.engine.log(statement)
		_, err := s.tx.Exec(statement.SQL(), statement.Bindings()...)
		if err != nil {
			s.tx = nil
			s.statements = []*Stmt{}
			return err
		}
	}

	err := s.tx.Commit()
	s.tx = nil
	s.statements = []*Stmt{}
	return err
}

// Rollback rollbacks the current transaction
func (s *Session) Rollback() error {
	if s.tx != nil {
		return s.tx.Rollback()
	}

	return errors.New("Current transaction is nil")
}

// Find returns a row given model properties
func (s *Session) Find(model interface{}) *Session {
	table := s.mapper.ModelName(model)
	modelMap := s.mapper.ToMap(model, true)

	cols := []Clause{}
	for k := range modelMap {
		cols = append(cols, s.T(table).C(k))
	}

	ands := []Clause{}

	for k := range modelMap {
		if modelMap[k] == nil {
			continue
		}
		ands = append(ands, Eq(s.metadata.Table(table).C(k), modelMap[k]))
	}

	s.builder = Select(cols...).From(s.T(table)).Where(And(ands...))
	return s
}

// Builder returns the active query builder of session
func (s *Session) Builder() Builder {
	return s.builder
}

// Statement builds the active query and returns it as a Stmt
func (s *Session) Statement() *Stmt {
	if s.isSelect() {
		if len(s.filters) > 0 {
			s.builder = (s.builder.(SelectStmt)).Where(And(s.filters...))
		}
	}

	statement := s.builder.Build(s.dialect)
	s.filters = []Clause{}
	s.builder = nil
	return statement
}

// Query starts a select statement given columns
func (s *Session) Query(clauses ...Clause) *Session {
	if len(clauses) == 0 {
		panic(fmt.Errorf("You must enter one or more column or aggregate paramater(s)"))
	} else {
		var table string
		for _, v := range clauses {
			if s.isCol(v) {
				table = (v.(ColumnElem)).Table
			}
		}
		s.builder = Select(clauses...)
		if table != "" {
			s.builder = (s.builder.(SelectStmt)).From(s.T(table))
		}
	}
	return s
}

// isCol returns if the clause is ColumnElem type
func (s *Session) isCol(clause Clause) bool {
	switch clause.(type) {
	case ColumnElem:
		return true
	default:
		return false
	}
}

// isSelect returns if the current builder is *Session
func (s *Session) isSelect() bool {
	switch s.builder.(type) {
	case SelectStmt:
		return true
	default:
		return false
	}
}

// Filter appends a filter to the current select statement
// NOTE: It currently only builds AndClause within the filters
// TODO: Add OR able filters
func (s *Session) Filter(conditional Clause) *Session {
	if !s.isSelect() {
		panic(fmt.Errorf("Please use Query(cols ...ColumnElem) before calling Filter()"))
	}

	s.filters = append(s.filters, conditional)
	return s
}

// From wraps select's From
// NOTE: You only need to set if Query() parameters are not columns
// No columns are in aggregate clauses
func (s *Session) From(table TableElem) *Session {
	if !s.isSelect() {
		panic(fmt.Errorf("Please use Query(cols ...ColumnElem) before calling From()"))
	}

	s.builder = (s.builder.(SelectStmt)).From(table)
	return s
}

// InnerJoin wraps select's InnerJoin
func (s *Session) InnerJoin(table TableElem, fromCol ColumnElem, col ColumnElem) *Session {
	if !s.isSelect() {
		panic(fmt.Errorf("Please use Query(cols ...ColumnElem) before calling InnerJoin()"))
	}

	s.builder = (s.builder.(SelectStmt)).InnerJoin(table, fromCol, col)
	return s
}

// CrossJoin wraps select's CrossJoin
func (s *Session) CrossJoin(table TableElem) *Session {
	if !s.isSelect() {
		panic(fmt.Errorf("Please use Query(cols ...ColumnElem) before calling CrossJoin()"))
	}

	s.builder = (s.builder.(SelectStmt)).CrossJoin(table)
	return s
}

// LeftJoin wraps select's LeftJoin
func (s *Session) LeftJoin(table TableElem, fromCol ColumnElem, col ColumnElem) *Session {
	if !s.isSelect() {
		panic(fmt.Errorf("Please use Query(cols ...ColumnElem) before calling LeftJoin()"))
	}

	s.builder = (s.builder.(SelectStmt)).LeftJoin(table, fromCol, col)
	return s
}

// RightJoin wraps select's RightJoin
func (s *Session) RightJoin(table TableElem, fromCol ColumnElem, col ColumnElem) *Session {
	if !s.isSelect() {
		panic(fmt.Errorf("Please use Query(cols ...ColumnElem) before calling RightJoin()"))
	}

	s.builder = (s.builder.(SelectStmt)).RightJoin(table, fromCol, col)
	return s
}

// GroupBy wraps the select's GroupBy
func (s *Session) GroupBy(cols ...ColumnElem) *Session {
	if !s.isSelect() {
		panic(fmt.Errorf("Please use Query(cols ...ColumnElem) before calling GroupBy()"))
	}

	s.builder = (s.builder.(SelectStmt)).GroupBy(cols...)
	return s
}

// Having wraps the select's Having
func (s *Session) Having(aggregate AggregateClause, op string, value interface{}) *Session {
	if !s.isSelect() {
		panic(fmt.Errorf("Please use Query(cols ...ColumnElem) before calling Having()"))
	}

	s.builder = (s.builder.(SelectStmt)).Having(aggregate, op, value)
	return s
}

// OrderBy wraps the select's OrderBy
func (s *Session) OrderBy(cols ...ColumnElem) *Session {
	if !s.isSelect() {
		panic(fmt.Errorf("Please use Query(cols ...ColumnElem) before calling OrderBy()"))
	}

	s.builder = (s.builder.(SelectStmt)).OrderBy(cols...).Asc()
	return s
}

// Asc wraps the select's Asc
func (s *Session) Asc() *Session {
	if !s.isSelect() {
		panic(fmt.Errorf("Please use Query(cols ...ColumnElem) & OrderBy() before calling Asc()"))
	}

	s.builder = (s.builder.(SelectStmt)).Asc()
	return s
}

// Desc wraps the select's Desc
func (s *Session) Desc() *Session {
	if !s.isSelect() {
		panic(fmt.Errorf("Please use Query(cols ...ColumnElem) & OrderBy() before calling Desc()"))
	}

	s.builder = (s.builder.(SelectStmt)).Desc()
	return s
}

// Limit wraps the select's Limit
func (s *Session) Limit(offset int, count int) *Session {
	if !s.isSelect() {
		panic(fmt.Errorf("Please use Query(cols ...ColumnElem) before calling Limit()"))
	}

	s.builder = (s.builder.(SelectStmt)).Limit(offset, count)
	return s
}

// Active query select & (insert/delete/update) ... returning ... finishers

// One returns the first record mapped as a model
// The interface should be struct pointer instead of struct
func (s *Session) One(model interface{}) error {
	return s.engine.Get(s.builder, model)
}

// All returns all the records mapped as a model slice
// The interface should be struct pointer instead of struct
func (s *Session) All(models interface{}) error {
	return s.engine.Select(s.builder, models)
}
