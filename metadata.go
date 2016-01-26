package qbit

func NewMetadata() *Metadata {
	return &Metadata{
		tables: []Table{},
	}
}

type Metadata struct {
	tables []Table
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
