package qb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhere(t *testing.T) {
	ctx := NewCompilerContext(NewDefaultDialect())
	assert.Equal(t,
		"WHERE X",
		Where(SQLText("X")).Accept(ctx))
	assert.Equal(t,
		"WHERE (X AND Y)",
		Where(SQLText("X"), SQLText("Y")).Accept(ctx))
}

func TestWhereAnd(t *testing.T) {
	ctx := NewCompilerContext(NewDefaultDialect())
	assert.Equal(t,
		"WHERE (X AND Y)",
		Where(SQLText("X")).And(SQLText("Y")).Accept(ctx))
	assert.Equal(t,
		"WHERE (X AND Y AND Z)",
		Where(SQLText("X")).And(SQLText("Y"), SQLText("Z")).Accept(ctx))
}

func TestWhereOr(t *testing.T) {
	ctx := NewCompilerContext(NewDefaultDialect())
	assert.Equal(t,
		"WHERE (X OR Y)",
		Where(SQLText("X")).Or(SQLText("Y")).Accept(ctx))
	assert.Equal(t,
		"WHERE (X OR Y OR Z)",
		Where(SQLText("X")).Or(SQLText("Y"), SQLText("Z")).Accept(ctx))
}
