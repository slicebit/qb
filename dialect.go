package qb

// NewDialect returns a dialect pointer given driver
func NewDialect(driver string) Dialect {
	dialect, ok := DialectRegistry[driver]
	if ok {
		return dialect
	}
	panic("No such dialect: " + driver)
}

// DialectRegistry is a global registry of dialects
var DialectRegistry = make(map[string]Dialect)

// RegisterDialect add a new dialect to the registry
func RegisterDialect(name string, dialect Dialect) {
	DialectRegistry[name] = dialect
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
	WrapError(err error) Error
}

// EscapeAll common escape all
func EscapeAll(dialect Dialect, strings []string) []string {
	for k, v := range strings {
		strings[k] = dialect.Escape(v)
	}

	return strings
}
