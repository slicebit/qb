package qb

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"
)

// NewSession generates a new Session given engine and returns session pointer
func NewSession(metadata *MetaData) *Session {
	return &Session{
		queries:  []*Query{},
		mapper:   NewMapper(metadata.Engine().Driver()),
		metadata: metadata,
		dialect:  NewDialect(metadata.Engine().Driver()),
	}
}

// Session is the composition of engine connection & orm mappings
type Session struct {
	queries  []*Query
	mapper   *Mapper
	metadata *MetaData
	tx       *sql.Tx
	dialect  *Dialect
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

// Delete adds a single delete query to the session
func (s *Session) Delete(model interface{}) {

	kv := s.mapper.ToMap(model)

	tName := s.mapper.ModelName(model)

	d := s.metadata.Table(tName).Delete()
	ands := []string{}
	bindings := []interface{}{}

	pcols := s.metadata.Table(tName).PrimaryKey()

	// if table has primary key
	if len(pcols) > 0 {

		for _, pk := range pcols {
			b := kv[pk]
			ands = append(ands, fmt.Sprintf("%s = %s", pk, d.Placeholder()))
			bindings = append(bindings, b)
		}

	} else {
		for k, v := range kv {
			ands = append(ands, fmt.Sprintf("%s = %s", s.mapper.ColName(k), d.Placeholder()))
			bindings = append(bindings, v)
		}
	}

	del := d.Where(d.And(ands...), bindings...).Query()
	s.add(del)
}

// Add adds a single query to the session. The query must be insert or update
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

// Select makers

// Find returns a row given model properties.
func (s *Session) Find(model interface{}) *Session {

	tName := s.mapper.ModelName(model)

	rModelMap := s.mapper.ToRawMap(model)
	modelMap := s.mapper.ToMap(model)

	sqlColNames := []string{}
	for k, _ := range rModelMap {
		sqlColNames = append(sqlColNames, s.mapper.ColName(k))
	}

	sort.Strings(sqlColNames)

	s.dialect = NewDialect(s.metadata.Engine().Driver())

	s.dialect.Select(sqlColNames...).From(tName)

	ands := []string{}
	bindings := []interface{}{}

	for k, v := range modelMap {
		ands = append(ands, fmt.Sprintf("%s = %s", s.mapper.ColName(k), s.dialect.Placeholder()))
		bindings = append(bindings, v)
	}

	s.dialect.Where(s.dialect.And(ands...), bindings...)
	return s
}

// TODO: Finish these implementations
// Query starts a select dialect given the model properties
func (s *Session) Query(model interface{}) *Session {
	return s
}

// FilterBy builds where statements given the conditions as map[string]interface{}
func (s *Session) FilterBy(m map[string]interface{}) *Session {
	return s
}

// Filter build complex filter statements such as gt, gte, st, ste, in, avg, count, etc.
func (s *Session) Filter() *Session {
	return s
}

// Join performs a join with another struct given model an optionally given explicit conditions
func (s *Session) Join(model interface{}, exConditions ...interface{}) *Session {
	return s
}

// First returns the first record mapped as a model
// The interface should be struct pointer instead of struct
func (s *Session) First(model interface{}) error {

	colNames := []string{}
	modelMap := s.mapper.ToRawMap(model)

	for k, _ := range modelMap {
		colNames = append(colNames, k)
	}

	sort.Strings(colNames)

	query := s.dialect.Query()
	rows, err := s.metadata.Engine().Query(query)
	if err != nil {
		return err
	}

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	values := make([]interface{}, len(cols))
	valuePtrs := make([]interface{}, len(cols))

	defer rows.Close()
	for rows.Next() {

		for i, _ := range cols {
			valuePtrs[i] = &values[i]
		}

		err = rows.Scan(valuePtrs...)
		if err != nil {
			return err
		}

		for i, _ := range cols {

			var v interface{}

			val := values[i]

			b, ok := val.([]byte)

			if ok {
				v = string(b)
			} else {
				v = val
			}

			modelMap[colNames[i]] = v
		}

		s.mapper.ToStruct(modelMap, model)
		return nil
	}

	return errors.New("Record not found")
}

// All returns all the records mapped as a model slice
// The interface should be struct pointer instead of struct
func (s *Session) All(models []interface{}) error {
	return errors.New("Record not found")
}
