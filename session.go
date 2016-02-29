package qb

// NewSession generates a new Session given engine and returns session pointer
func NewSession(engine *Engine) *Session {
	return &Session{
		engine:  engine,
		queries: []*Query{},
	}
}

// Session is the composition of engine connection & orm mappings
type Session struct {
	engine  *Engine
	queries []*Query
}

// Engine returns the current session engine
func (s *Session) Engine() *Engine {
	return s.engine
}

func (s *Session) add(query *Query) {
	s.queries = append(s.queries, query)
}

// Delete adds a single delete query to the session
func (s *Session) Delete(query *Query) {
	s.add(query)
}

// Add adds a single query to the session. The query must be insert or update
func (s *Session) Add(query *Query) {
	s.add(query)
}

// AddAll adds multiple queries to the session. The query must be insert or update
func (s *Session) AddAll(queries ...*Query) {
	s.queries = append(s.queries, queries...)
}
