package qb

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultDialect(t *testing.T) {
	dialect := NewDialect("default")
	assert.Implements(t, (*Compiler)(nil), dialect.GetCompiler())
	assert.Equal(t, false, dialect.SupportsUnsigned())
	assert.Equal(t, "test", dialect.Escape("test"))
	assert.Equal(t, false, dialect.Escaping())
	dialect.SetEscaping(true)
	assert.Equal(t, true, dialect.Escaping())
	assert.Equal(t, "`test`", dialect.Escape("test"))
	assert.Equal(t, []string{"`test`"}, dialect.EscapeAll([]string{"test"}))
	assert.Equal(t, "", dialect.Driver())

	autoincCol := Column("id", Int()).PrimaryKey().AutoIncrement()
	assert.Equal(t,
		"INT PRIMARY KEY AUTO INCREMENT",
		dialect.AutoIncrement(&autoincCol))

	err := errors.New("xxx")
	qbErr := dialect.WrapError(err)
	assert.Equal(t, err, qbErr.Orig)
}
