package qb

import (
	"database/sql"
	"errors"
	"strings"
	"sync"
)

// New generates a new Session given engine and returns session pointer
func New(driver string, dsn string) (*Session, error) {
	engine, err := NewEngine(driver, dsn)
	if err != nil {
		return nil, err
	}

	builder := NewBuilder(engine.Driver())

	return &Session{
		queries:  []*QueryElem{},
		engine:   engine,
		mapper:   Mapper(builder.Adapter()),
		metadata: MetaData(builder),
		builder:  builder,
		mutex:    &sync.Mutex{},
	}, nil
}

// Session is the composition of engine connection & orm mappings
type Session struct {
	queries  []*QueryElem
	engine   *Engine
	mapper   MapperElem
	metadata *MetaDataElem
	tx       *sql.Tx
	builder  *Builder
	mutex    *sync.Mutex
}

func (s *Session) add(query *QueryElem) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	var err error
	if s.tx == nil {
		s.queries = []*QueryElem{}
		s.tx, err = s.engine.DB().Begin()
		if err != nil {
			panic(err)
		}
	}
	s.queries = append(s.queries, query)
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

// Engine returns the current sqlx wrapped engine
func (s *Session) Engine() *Engine {
	return s.engine
}

// Close closes engine db (sqlx) connection
func (s *Session) Close() {
	s.engine.DB().Close()
}

// Builder returns query builder
func (s *Session) Builder() *Builder {
	return s.builder
}

// AddQuery adds a query given the query pointer retrieved from Query() function
func (s *Session) AddQuery(query *QueryElem) {
	s.add(query)
}

// Query returns the active query built by session
func (s *Session) Query() *QueryElem {
	return s.builder.Query()
}

// Metadata returns the metadata of session
func (s *Session) Metadata() *MetaDataElem {
	return s.metadata
}

// Delete adds a single delete query to the session
func (s *Session) Delete(model interface{}) {
	kv := s.mapper.ToMap(model, false)

	tName := s.mapper.ModelName(model)

	d := s.builder.Delete(tName)
	ands := []string{}
	for k, v := range kv {
		ands = append(ands, s.Eq(k, v))
	}

	del := d.Where(d.And(ands...)).Query()
	s.add(del)
}

// Add adds a single model to the session. The query must be insert or update
func (s *Session) Add(model interface{}) {
	m := s.mapper.ToMap(model, false)
	q := s.builder.Insert(s.mapper.ModelName(model)).Values(m).Query()
	s.add(q)
}

// AddAll adds multiple models an adds an insert statement to current queries
func (s *Session) AddAll(models ...interface{}) {
	for _, m := range models {
		s.Add(m)
	}
}

// Commit commits the current transaction with queries
func (s *Session) Commit() error {
	for _, q := range s.queries {
		_, err := s.tx.Exec(q.SQL(), q.Bindings()...)
		if err != nil {
			s.tx = nil
			s.queries = []*QueryElem{}
			return err
		}
	}

	err := s.tx.Commit()
	s.tx = nil
	s.queries = []*QueryElem{}
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
	tName := s.mapper.ModelName(model)
	modelMap := s.mapper.ToMap(model, true)

	sqlColNames := []string{}
	for k := range modelMap {
		sqlColNames = append(sqlColNames, k)
	}

	ands := []string{}

	for k := range modelMap {
		if modelMap[k] == nil {
			continue
		}
		ands = append(ands, s.builder.Eq(k, modelMap[k]))
	}

	s.builder.
		Select(s.builder.Adapter().EscapeAll(sqlColNames)...).
		From(tName).
		Where(s.builder.And(ands...))

	return s
}

// One returns the first record mapped as a model
// The interface should be struct pointer instead of struct
func (s *Session) One(model interface{}) error {
	query := s.builder.Query()
	return s.engine.Get(query, model)
}

// All returns all the records mapped as a model slice
// The interface should be struct pointer instead of struct
func (s *Session) All(models interface{}) error {
	query := s.builder.Query()
	return s.engine.Select(query, models)
}

// builder overrides for session

// Update generates "update %s" statement
func (s *Session) Update(table string) *Session {
	s.builder.Update(table)
	return s
}

// Set generates "set a = placeholder" statement for each key a and add bindings for map value
func (s *Session) Set(m map[string]interface{}) *Session {
	s.builder.Set(m)
	return s
}

// Select generates "select %s" statement
func (s *Session) Select(columns ...string) *Session {
	s.builder.Select(columns...)
	return s
}

// From generates "from %s" statement for each table name
func (s *Session) From(tables ...string) *Session {
	s.builder.From(tables...)
	return s
}

// InnerJoin generates "inner join %s on %s" statement for each expression
func (s *Session) InnerJoin(table string, expressions ...string) *Session {
	s.builder.InnerJoin(table, expressions...)
	return s
}

// CrossJoin generates "cross join %s" statement for table
func (s *Session) CrossJoin(table string) *Session {
	s.builder.CrossJoin(table)
	return s
}

// LeftOuterJoin generates "left outer join %s on %s" statement for each expression
func (s *Session) LeftOuterJoin(table string, expressions ...string) *Session {
	s.builder.LeftOuterJoin(table, expressions...)
	return s
}

// RightOuterJoin generates "right outer join %s on %s" statement for each expression
func (s *Session) RightOuterJoin(table string, expressions ...string) *Session {
	s.builder.RightOuterJoin(table, expressions...)
	return s
}

// FullOuterJoin generates "full outer join %s on %s" for each expression
func (s *Session) FullOuterJoin(table string, expressions ...string) *Session {
	s.builder.FullOuterJoin(table, expressions...)
	return s
}

// Where generates "where %s" for the expression and adds bindings for each value
func (s *Session) Where(expression string, bindings ...interface{}) *Session {
	expression = strings.Replace(expression, "?", s.builder.Adapter().Placeholder(), -1)
	s.builder.Where(expression, bindings...)
	return s
}

// OrderBy generates "order by %s" for each expression
func (s *Session) OrderBy(expressions ...string) *Session {
	s.builder.OrderBy(expressions...)
	return s
}

// GroupBy generates "group by %s" for each column
func (s *Session) GroupBy(columns ...string) *Session {
	s.builder.GroupBy(columns...)
	return s
}

// Having generates "having %s" for each expression
func (s *Session) Having(expressions ...string) *Session {
	s.builder.Having(expressions...)
	return s
}

// Limit generates limit %d offset %d for offset and count
func (s *Session) Limit(offset int, count int) *Session {
	s.builder.Limit(offset, count)
	return s
}

// aggregates

// Avg function generates "avg(%s)" statement for column
func (s *Session) Avg(column string) string {
	return s.builder.Avg(column)
}

// Count function generates "count(%s)" statement for column
func (s *Session) Count(column string) string {
	return s.builder.Count(column)
}

// Sum function generates "sum(%s)" statement for column
func (s *Session) Sum(column string) string {
	return s.builder.Sum(column)
}

// Min function generates "min(%s)" statement for column
func (s *Session) Min(column string) string {
	return s.builder.Min(column)
}

// Max function generates "max(%s)" statement for column
func (s *Session) Max(column string) string {
	return s.builder.Max(column)
}

// expressions

// NotIn function generates "%s not in (%s)" for key and adds bindings for each value
func (s *Session) NotIn(key string, values ...interface{}) string {
	return s.builder.NotIn(key, values...)
}

// In function generates "%s in (%s)" for key and adds bindings for each value
func (s *Session) In(key string, values ...interface{}) string {
	return s.builder.In(key, values...)
}

// NotEq function generates "%s != placeholder" for key and adds binding for value
func (s *Session) NotEq(key string, value interface{}) string {
	return s.builder.NotEq(key, value)
}

// Eq function generates "%s = placeholder" for key and adds binding for value
func (s *Session) Eq(key string, value interface{}) string {
	return s.builder.Eq(key, value)
}

// Gt function generates "%s > placeholder" for key and adds binding for value
func (s *Session) Gt(key string, value interface{}) string {
	return s.builder.Gt(key, value)
}

// Gte function generates "%s >= placeholder" for key and adds binding for value
func (s *Session) Gte(key string, value interface{}) string {
	return s.builder.Gte(key, value)
}

// St function generates "%s < placeholder" for key and adds binding for value
func (s *Session) St(key string, value interface{}) string {
	return s.builder.St(key, value)
}

// Ste function generates "%s <= placeholder" for key and adds binding for value
func (s *Session) Ste(key string, value interface{}) string {
	return s.builder.Ste(key, value)
}

// And function generates " AND " between any number of expressions
func (s *Session) And(expressions ...string) string {
	return s.builder.And(expressions...)
}

// Or function generates " OR " between any number of expressions
func (s *Session) Or(expressions ...string) string {
	return s.builder.Or(expressions...)
}
