package qbit

import (
//	"github.com/aacanakin/qbit/mysql"
//	"github.com/aacanakin/qbit/postgres"
//	"github.com/aacanakin/qbit/sqlite"
//	"log"
)

func NewMetadata(engine *Engine) *Metadata {

	//	var mapper Mapper
	//
	//	if engine.Driver() == "mysql" {
	//		mapper = mysql.NewMapper()
	//	} else if engine.Driver() == "sqlite" {
	//		mapper = sqlite.NewMapper()
	//	} else if engine.Driver() == "postgres" {
	//		mapper = postgres.NewMapper()
	//	} else {
	//		log.Fatalln("Invalid Driver: ", engine.Driver())
	//	}

	return &Metadata{
		tables: []Table{},
		engine: engine,
		//		mapper: mapper,
	}
}

type Metadata struct {
	tables []Table
	engine *Engine
}

func (m *Metadata) Add(model interface{}) {
	//	m.AddTable(m.mapper.Convert(model))
}

func (m *Metadata) AddTable(table Table) {
	m.tables = append(m.tables, table)
}

func (m *Metadata) Tables() []Table {
	return m.tables
}

func (m *Metadata) CreateAll(engine *Engine) error {
	return nil
}

func (m *Metadata) DropAll(engine *Engine) error {
	return nil
}
