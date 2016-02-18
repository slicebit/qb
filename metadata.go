package qb

import (
//	"github.com/aacanakin/qbit/mysql"
//	"github.com/aacanakin/qbit/postgres"
//	"github.com/aacanakin/qbit/sqlite"
//	"log"
)

// NewMetaData creates a new MetaData object and returns
func NewMetaData(engine *Engine) *MetaData {

	return &MetaData{
		tables: []*Table{},
		engine: engine,
		mapper: NewMapper(engine.Driver()),
	}
}

// MetaData is the container for database structs and tables
type MetaData struct {
	tables []*Table
	engine *Engine
	mapper *Mapper
}

// Add retrieves the struct and converts it using mapper and appends to tables slice
func (m *MetaData) Add(model interface{}) {

	table, err := m.mapper.Convert(model)
	if err != nil {
		panic(err)
	}

	m.AddTable(table)
}

// AddTable appends table to tables slice
func (m *MetaData) AddTable(table *Table) {
	m.tables = append(m.tables, table)
}

// Tables returns the current tables slice
func (m *MetaData) Tables() []*Table {
	return m.tables
}

// CreateAll creates all the tables added to metadata
func (m *MetaData) CreateAll(engine *Engine) error {
	return nil
}

// DropAll drops all the tables which is added to metadata
func (m *MetaData) DropAll(engine *Engine) error {
	return nil
}
