package qbit

func NewSession(engine *Engine) *Session {
	return &Session{
		engine: engine,
	}
}

type Session struct {
	engine *Engine
}

func (s *Session) Engine() *Engine {
	return s.engine
}
