package qb

import "log"

// There are the log flags qb can use
// The default log flag is not at all logging
const (
	LDefault = iota
	// log query flag
	LQuery
	// log bindings flag
	LBindings
)

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

	LogFlags() int
	SetLogFlags(int)
}

// DefaultLogger is the default logger of qb engine unless engine.SetLogger() is not called
type DefaultLogger struct {
	logFlags int
	*log.Logger
}

// SetLogFlags sets the logflags
// It is for changing engine log flags
func (l DefaultLogger) SetLogFlags(logFlags int) {
	l.logFlags = logFlags
}

// LogFlags gets the logflags as an int
func (l DefaultLogger) LogFlags() int {
	return l.logFlags
}
