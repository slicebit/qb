package qb

import (
	"strings"
	"testing"
)

// func asDefSQL(clause Clause) string {
// 	return asSQL(clause, NewDialect("default"))
// }

// func asDefSQLBinds(clause Clause) (string, []interface{}) {
// 	return asSQLBinds(clause, NewDialect("default"))
// }

// func asSQL(clause Clause, dialect Dialect) string {
// 	sql, _ := asSQLBinds(clause, dialect)
// 	return sql
// }

// func asSQLBinds(clause Clause, dialect Dialect) (string, []interface{}) {
// 	ctx := NewCompilerContext(dialect)
// 	return clause.Accept(ctx), ctx.Binds
// }

type TestingLogWriter struct {
	t     *testing.T
	lines []string
}

func (w *TestingLogWriter) Write(p []byte) (n int, err error) {
	w.lines = append(w.lines, string(p))
	return len(p), nil
}

func (w *TestingLogWriter) Flush() {
	w.t.Log("Captured:\n" + strings.Join(w.lines, ""))
}
