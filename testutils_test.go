package qb

import (
	"testing"
)

func asSQL(clause Clause, dialect Dialect) (string, []interface{}) {
	ctx := NewCompilerContext(dialect)
	return clause.Accept(ctx), ctx.Binds
}

type TestingLogWriter struct {
	t *testing.T
}

func (w TestingLogWriter) Write(p []byte) (n int, err error) {
	w.t.Log(string(p))
	return len(p), nil
}
