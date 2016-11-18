package qb

import "log"

// These are the log flags qb can use
const (
	// LDefault is the default flag that logs nothing
	LDefault LogFlags = 0
	// LQuery Flag to log queries
	LQuery LogFlags = 1 << iota
	// LBindings Flag to log bindings
	LBindings
)

// LogFlags is the type we use for flags that can be passed
// to the logger
type LogFlags uint

// Logger is the std logger interface of the qb engine
type Logger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})

	LogFlags() LogFlags
	SetLogFlags(LogFlags)
}

// DefaultLogger is the default logger of qb engine unless engine.SetLogger() is not called
type DefaultLogger struct {
	logFlags LogFlags
	*log.Logger
}

// SetLogFlags sets the logflags
// It is for changing engine log flags
func (l *DefaultLogger) SetLogFlags(logFlags LogFlags) {
	l.logFlags = logFlags
}

// LogFlags gets the logflags as an int
func (l *DefaultLogger) LogFlags() LogFlags {
	return l.logFlags
}
