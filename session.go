package qb

import (
	"database/sql"
	"errors"
	"fmt"
)

// NewSession generates a new Session given engine and returns session pointer
func NewSession(metadata *MetaData) *Session {
	return &Session{
		queries:  []*Query{},
		mapper:   NewMapper(metadata.Engine().Driver()),
		metadata: metadata,
	}
}

// Session is the composition of engine connection & orm mappings
type Session struct {
	queries  []*Query
	mapper   *Mapper
	metadata *MetaData
	tx       *sql.Tx
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
func (s *Session) Find(model interface{}) error {

	tName := s.mapper.ModelName(model)

	rawModelMap := s.mapper.ToRawMap(model)
	modelMap := s.mapper.ToMap(model)

	colNames := []string{}
	sqlColNames := []string{}
	for k, _ := range rawModelMap {
		colNames = append(colNames, k)
		sqlColNames = append(sqlColNames, s.mapper.ColName(k))
	}

	d := NewDialect(s.metadata.Engine().Driver()).Select(sqlColNames...).From(tName)
	ands := []string{}
	bindings := []interface{}{}

	for k, v := range modelMap {
		ands = append(ands, fmt.Sprintf("%s = %s", s.mapper.ColName(k), d.Placeholder()))
		bindings = append(bindings, v)
	}

	sel := d.Where(d.And(ands...), bindings...).Query()
	rows, err := s.metadata.Engine().Query(sel)
	if err != nil {
		return err
	}

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	rawResult := make([][]byte, len(cols))
	result := make([]interface{}, len(cols))

	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i, _ := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			return err
		}

		for i, raw := range rawResult {
			if raw == nil {
				result[i] = nil
			} else {
				result[i] = string(raw)
			}
		}

		for i := 0; i < len(colNames); i++ {
			rawModelMap[colNames[i]] = result[i]
		}

		s.mapper.ToStruct(rawModelMap, model)
		return nil
	}

	return errors.New("Record not found")
}
