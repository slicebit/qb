package qbit

// NewSession generates a new Session given engine and returns session pointer
func NewSession(engine *Engine) *Session {
	return &Session{
		engine: engine,
	}
}

// Session is the composition of engine connection & orm mappings
type Session struct {
	engine *Engine
}

// Engine returns the current session engine
func (s *Session) Engine() *Engine {
	return s.engine
}
