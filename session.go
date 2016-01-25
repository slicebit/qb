package qbit

type Session struct {
	Engine Engine
}

func Session(engine Engine) *Session {
	return &Session{
		Engine: engine,
	}
}
