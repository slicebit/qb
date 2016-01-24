package qbit

type metadata struct {
	tables []table
}

func MetaData() *metadata {
	return &metadata{
		tables: new([]table),
	}
}

func (m *metadata) AddTable(table table) {
	m.tables = append(m.tables, table)
}

func (m *metadata) Tables() {
	return m.tables
}

func (m *metadata) CreateAll(engine *Engine) error {
	return nil
}

func (m *metadata) DropAll(engine *Engine) error {
	return nil
}
