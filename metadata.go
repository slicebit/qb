package qbit

import (
//	"github.com/aacanakin/qbit/mysql"
//	"github.com/aacanakin/qbit/postgres"
//	"github.com/aacanakin/qbit/sqlite"
//	"log"
)

func MetaData(engine *Engine) *metadata {

	return &metadata{
		tables: []Table{},
		engine: engine,
		mapper: NewMapper(engine.Driver()),
	}
}

type metadata struct {
	tables []Table
	engine *Engine
	mapper *Mapper
}

func (m *metadata) Add(model interface{}) {
	//	m.AddTable(m.mapper.Convert(model))
}

func (m *metadata) AddTable(table Table) {
	m.tables = append(m.tables, table)
}

func (m *metadata) Tables() []Table {
	return m.tables
}

func (m *metadata) CreateAll(engine *Engine) error {
	return nil
}

func (m *metadata) DropAll(engine *Engine) error {
	return nil
}
