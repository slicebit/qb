package qb

// NewDialect returns a dialect pointer given driver
func NewDialect(driver string) Dialect {
	switch driver {
	case "postgres":
		return &PostgresDialect{escaping: false, bindingIndex: 0}
	case "mysql":
		return &MysqlDialect{false}
	case "sqlite3":
		return &SqliteDialect{false}
	default:
		return &DefaultDialect{false}
	}
}

// Dialect is the common adapter for driver changes
// It is for fixing compatibility issues of different drivers
type Dialect interface {
	Escape(str string) string
	EscapeAll([]string) []string
	SetEscaping(escaping bool)
	Escaping() bool
	Placeholder() string
	Placeholders(values ...interface{}) []string
	AutoIncrement() string
	Reset()
	SupportsUnsigned() bool
	Driver() string
}

// common escape all
func escapeAll(dialect Dialect, strings []string) []string {
	for k, v := range strings {
		strings[k] = dialect.Escape(v)
	}

	return strings
}

// common placeholders
func placeholders(dialect Dialect, values ...interface{}) []string {
	placeholders := []string{}
	for range values {
		placeholders = append(placeholders, dialect.Placeholder())
	}
	return placeholders
}
