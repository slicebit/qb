package qbit

func NewSession(engine *Engine) *session {
	return &session{
		engine: engine,
	}
}

type session struct {
	engine *Engine
}

func (s *session) Engine() *Engine {
	return s.engine
}
