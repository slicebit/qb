package qb

// NewAdapter returns a adapter pointer given driver
func NewAdapter(driver string) Adapter {
	switch driver {
	case "postgres":
		return &PostgresAdapter{escaping: false, bindingIndex: 0}
	case "mysql":
		return &MysqlAdapter{false}
	case "sqlite3":
		return &SqliteAdapter{false}
	default:
		return &DefaultAdapter{false}
	}
}

// Adapter is the common adapter for driver changes
// It is for fixing compatibility issues of different drivers
type Adapter interface {
	Escape(str string) string
	EscapeAll([]string) []string
	SetEscaping(escaping bool)
	GetEscaping() bool
	Placeholder() string
	Placeholders(values ...interface{}) []string
	AutoIncrement() string
	Reset()
	SupportsInlinePrimaryKey() bool
	Driver() string
}

// common escape all
func escapeAll(adapter Adapter, strings []string) []string {
	for k, v := range strings {
		strings[k] = adapter.Escape(v)
	}

	return strings
}

// common placeholders
func placeholders(adapter Adapter, values ...interface{}) []string {
	placeholders := []string{}
	for _ = range values {
		placeholders = append(placeholders, adapter.Placeholder())
	}
	return placeholders
}
