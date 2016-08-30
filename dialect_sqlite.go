package qb

// SqliteDialect is a type of dialect that can be used with sqlite driver
type SqliteDialect struct {
	escaping bool
}

// NewSqliteDialect instanciate a SqliteDialect
func NewSqliteDialect() Dialect {
	return &SqliteDialect{false}
}

func init() {
	RegisterDialect("sqlite3", NewSqliteDialect)
	RegisterDialect("sqlite", NewSqliteDialect)
}

// CompileType compiles a type into its DDL
func (d *SqliteDialect) CompileType(t TypeElem) string {
	if t.Name == "UUID" {
		return "VARCHAR(36)"
	}
	return DefaultCompileType(t, d.SupportsUnsigned())
}

// Escape wraps the string with escape characters of the dialect
func (d *SqliteDialect) Escape(str string) string {
	return str
}

// EscapeAll wraps all elements of string array
func (d *SqliteDialect) EscapeAll(strings []string) []string {
	return escapeAll(d, strings[0:])
}

// SetEscaping sets the escaping parameter of dialect
func (d *SqliteDialect) SetEscaping(escaping bool) {
	d.escaping = escaping
}

// Escaping gets the escaping parameter of dialect
func (d *SqliteDialect) Escaping() bool {
	return d.escaping
}

// Placeholder returns the placeholder for bindings in the sql
func (d *SqliteDialect) Placeholder() string {
	return "?"
}

// Placeholders returns the placeholders for bindings in the sql
func (d *SqliteDialect) Placeholders(values ...interface{}) []string {
	return placeholders(d, values...)
}

// AutoIncrement generates auto increment sql of current dialect
func (d *SqliteDialect) AutoIncrement(colum *ColumnElem) string {
	return "INTEGER PRIMARY KEY"
}

// Reset does nothing for the default driver
func (d *SqliteDialect) Reset() {}

// SupportsUnsigned returns whether driver supports unsigned type mappings or not
func (d *SqliteDialect) SupportsUnsigned() bool { return false }

// Driver returns the current driver of dialect
func (d *SqliteDialect) Driver() string {
	return "sqlite3"
}

func (d *SqliteDialect) GetCompiler() Compiler {
	return SQLCompiler{d}
}
