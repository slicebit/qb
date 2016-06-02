package qb

// TODO: Metadata should not use builder, it should only use adapter
// MetaData creates a new MetaData object and returns
func MetaData(engine *Engine, builder *Builder) *MetaDataElem {
	return &MetaDataElem{
		tables:  []TableElem{},
		engine:  engine,
		mapper:  Mapper(builder.Adapter()),
		builder: builder,
	}
}

// MetaDataElem is the container for database structs and tables
type MetaDataElem struct {
	tables  []TableElem
	engine  *Engine
	mapper  MapperElem
	builder *Builder
}

// Engine returns the currently bound engine of metadata
func (m *MetaDataElem) Engine() *Engine {
	return m.engine
}

// Add retrieves the struct and converts it using mapper and appends to tables slice
func (m *MetaDataElem) Add(model interface{}) {
	table, err := m.mapper.ToTable(model)
	if err != nil {
		panic(err)
	}

	m.AddTable(table)
}

// AddTable appends table to tables slice
func (m *MetaDataElem) AddTable(table TableElem) {
	m.tables = append(m.tables, table)
}

// Table returns the metadata registered table object. It returns nil if table is not found
func (m *MetaDataElem) Table(name string) *TableElem {
	for _, t := range m.tables {
		if t.Name == name {
			return &t
		}
	}

	return nil
}

// Tables returns the current tables slice
func (m *MetaDataElem) Tables() []TableElem {
	return m.tables
}

// CreateAll creates all the tables added to metadata
func (m *MetaDataElem) CreateAll() error {
	tx, err := m.engine.DB().Begin()
	if err != nil {
		return err
	}

	for _, t := range m.tables {
		_, err = tx.Exec(t.Create(m.builder.Adapter()))
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	return err
}

// DropAll drops all the tables which is added to metadata
func (m *MetaDataElem) DropAll() error {
	tx, err := m.engine.DB().Begin()
	if err != nil {
		return err
	}

	for i := len(m.tables) - 1; i >= 0; i-- {
		drop := m.builder.DropTable(m.tables[i].Name).Query()
		_, err = tx.Exec(drop.SQL())
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	return err
}
