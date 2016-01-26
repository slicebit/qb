package qbit

func Session(engine Engine) *session {
	return &session{
		engine: engine,
	}
}

type session struct {
	engine Engine
}
