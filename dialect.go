package qb

// NewDialect returns a dialect pointer given driver
func NewDialect(driver string) Dialect {
	factory, ok := DialectRegistry[driver]
	if ok {
		return factory()
	}
	return &DefaultDialect{false}
}

// A DialectFactory is a Dialect Factory
type DialectFactory func() Dialect

// DialectRegistry is a global registry of dialects
var DialectRegistry = make(map[string]DialectFactory)

// RegisterDialect add a new dialect to the registry
func RegisterDialect(name string, factory DialectFactory) {
	DialectRegistry[name] = factory
}

// Dialect is the common interface for driver changes
// It is for fixing compatibility issues of different drivers
type Dialect interface {
	GetCompiler() Compiler
	CompileType(t TypeElem) string
	Escape(str string) string
	EscapeAll([]string) []string
	SetEscaping(escaping bool)
	Escaping() bool
	AutoIncrement(column *ColumnElem) string
	SupportsUnsigned() bool
	Driver() string
	ExtractError(err error, stmt *Stmt) Error
}

// common escape all
func escapeAll(dialect Dialect, strings []string) []string {
	for k, v := range strings {
		strings[k] = dialect.Escape(v)
	}

	return strings
}
