package qb

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

// Engine returns the currently bound engine of metadata
func (m *MetaData) Engine() *Engine {
	return m.engine
}

// Add retrieves the struct and converts it using mapper and appends to tables slice
func (m *MetaData) Add(model interface{}) {

	table, err := m.mapper.ToTable(model)
	if err != nil {
		panic(err)
	}

	m.AddTable(table)
}

// AddTable appends table to tables slice
func (m *MetaData) AddTable(table *Table) {
	m.tables = append(m.tables, table)
}

// Table returns the metadata registered table object. It returns nil if table is not found
func (m *MetaData) Table(name string) *Table {

	for _, t := range m.tables {
		if t.name == name {
			return t
		}
	}

	return nil
}

// Tables returns the current tables slice
func (m *MetaData) Tables() []*Table {
	return m.tables
}

// CreateAll creates all the tables added to metadata
func (m *MetaData) CreateAll() error {

	tx, err := m.engine.DB().Begin()
	if err != nil {
		return err
	}

	for _, t := range m.tables {
		_, err = tx.Exec(t.SQL())
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	return err
}

// DropAll drops all the tables which is added to metadata
func (m *MetaData) DropAll() error {

	b := NewBuilder(m.engine.Driver())

	tx, err := m.engine.DB().Begin()
	if err != nil {
		return err
	}

	for i := len(m.tables) - 1; i >= 0; i-- {
		drop := b.DropTable(m.tables[i].Name()).Query()
		_, err = tx.Exec(drop.SQL())
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	return err
}
