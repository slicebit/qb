package qb

import (
	"strings"
	"testing"
)

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
