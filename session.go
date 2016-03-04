package qb

import (
	"database/sql"
	"fmt"
)

// NewSession generates a new Session given engine and returns session pointer
func NewSession(metadata *MetaData) *Session {
	return &Session{
		queries: []*Query{},
		mapper:  NewMapper(metadata.Engine().Driver()),
		metadata: metadata,
	}
}

// Session is the composition of engine connection & orm mappings
type Session struct {
	queries []*Query
	mapper  *Mapper
	metadata *MetaData
	tx      *sql.Tx
}

func (s *Session) add(query *Query) {
	var err error
	if s.tx == nil {
		s.tx, err = s.metadata.Engine().DB().Begin()
		if err != nil {
			panic(err)
		}
	}
	s.queries = append(s.queries, query)
}

// Delete adds a single delete query to the session
func (s *Session) Delete(model interface{}) {

	kv := s.mapper.ConvertStructToMap(model)

	d := s.metadata.Table(s.mapper.ModelName(model)).Delete()

	ands := []string{}
	bindings := []interface{}{}
	for k, v := range kv {
		ands = append(ands, fmt.Sprintf("%s = %s", s.mapper.ColName(k), d.Placeholder()))
		bindings = append(bindings, v)
	}

	del := d.Where(d.And(ands...), bindings...).Query()
	s.add(del)
}

// Add adds a single query to the session. The query must be insert or update
func (s *Session) Add(model interface{}) {

	rawMap := s.mapper.ConvertStructToMap(model)

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

	return s.tx.Commit()
}

// Select makers
