package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var compileTests = []struct {
	clause Clause
	expect string
	binds  []interface{}
}{
	{SQLText("1"), "1", nil},
}

func TestCompile(t *testing.T) {
	compile := func(clause Clause) (string, []interface{}) {
		context := NewCompilerContext(NewDialect("default"))
		return clause.Accept(context), context.Binds
	}

	for _, tt := range compileTests {
		actual, binds := compile(tt.clause)
		assert.Equal(t, tt.expect, actual)
		assert.Equal(t, tt.binds, binds)
	}
}
