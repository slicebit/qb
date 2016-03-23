package qb

import (
	"database/sql"
	"fmt"
	"strings"
)

// New generates a new Session given engine and returns session pointer
func New(driver string, dsn string) (*Session, error) {

	engine, err := NewEngine(driver, dsn)
	if err != nil {
		return nil, err
	}

	return &Session{
		queries:  []*Query{},
		mapper:   NewMapper(engine.Driver()),
		metadata: NewMetaData(engine),
		builder:  NewBuilder(engine.Driver()),
	}, nil
}

// Session is the composition of engine connection & orm mappings
type Session struct {
	queries  []*Query
	mapper   *Mapper
	metadata *MetaData
	tx       *sql.Tx
	builder  *Builder
}

func (s *Session) add(query *Query) {
	var err error
	if s.tx == nil {
		s.queries = []*Query{}
		s.tx, err = s.metadata.Engine().DB().Begin()
		if err != nil {
			panic(err)
		}
	}
	s.queries = append(s.queries, query)
}

// Close closes engine db (sqlx) connection
func (s *Session) Close() {
	s.metadata.Engine().DB().Close()
}

// Builder returns query builder
func (s *Session) Builder() *Builder{
	return s.builder
}

// AddQuery adds a query given the query pointer retrieved from Query() function
func (s *Session) AddQuery(query *Query) {
	s.add(query)
}

// Query returns the active query built by session
func (s *Session) Query() *Query {
	return s.builder.Query()
}

// Metadata returns the metadata of session
func (s *Session) Metadata() *MetaData {
	return s.metadata
}

// Delete adds a single delete query to the session
func (s *Session) Delete(model interface{}) {

	kv := s.mapper.ToMap(model)

	tName := s.mapper.ModelName(model)

	d := s.metadata.Table(tName).Delete()
	ands := []string{}
	bindings := []interface{}{}
	for k, v := range kv {
		ands = append(ands, fmt.Sprintf("%s = %s", s.mapper.ColName(k), s.builder.Dialect().Placeholder()))
		bindings = append(bindings, v)
	}

	del := d.Where(d.And(ands...), bindings...).Query()
	s.add(del)
}

// Add adds a single model to the session. The query must be insert or update
func (s *Session) Add(model interface{}) {

	rawMap := s.mapper.ToMap(model)

	kv := map[string]interface{}{}

	for k, v := range rawMap {
		kv[s.mapper.ColName(k)] = v
	}

	q := s.metadata.Table(s.mapper.ModelName(model)).Insert(kv).Query()
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
			return err
		}
	}

	err := s.tx.Commit()
	s.tx = nil
	return err
}

// Find returns a row given model properties
func (s *Session) Find(model interface{}) *Session {

	tName := s.mapper.ModelName(model)
	rModelMap := s.mapper.ToRawMap(model)

	sqlColNames := []string{}
	for k := range rModelMap {
		sqlColNames = append(sqlColNames, s.mapper.ColName(k))
	}

	s.builder = NewBuilder(s.metadata.Engine().Driver())
	s.builder.Select(s.builder.Dialect().EscapeAll(sqlColNames)...).From(tName)

	modelMap := s.mapper.ToMap(model)

	ands := []string{}
	bindings := []interface{}{}

	for k, v := range modelMap {
		ands = append(ands, fmt.Sprintf("%s = %s", s.mapper.ColName(k), s.builder.Dialect().Placeholder()))
		bindings = append(bindings, v)
	}

	s.builder.Where(s.builder.And(ands...), bindings...)
	return s
}

// First returns the first record mapped as a model
// The interface should be struct pointer instead of struct
func (s *Session) First(model interface{}) error {
	query := s.builder.Query()
	return s.metadata.Engine().Get(query, model)
}

// All returns all the records mapped as a model slice
// The interface should be struct pointer instead of struct
func (s *Session) All(models interface{}) error {
	query := s.builder.Query()
	return s.metadata.Engine().Select(query, models)
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
	expression = strings.Replace(expression, "?", s.builder.Dialect().Placeholder(), -1)
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
