package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	TTGroup = Table(
		"group",
		Column("id", Int()).AutoIncrement().PrimaryKey(),
		Column("name", Text()).Unique(),
	)

	TTUser = Table(
		"user",
		Column("id", Int()).AutoIncrement().PrimaryKey(),
		Column("name", Text()).Unique(),
		Column("main_group_id", Int()),
		ForeignKey().Ref("main_group_id", "group", "id"),
	)
)

var compileTests = []struct {
	clause Clause
	expect string
	binds  []interface{}
}{
	{SQLText("1"), "1", nil},
	{
		Exists(Select(TTGroup.C("name")).From(TTGroup).Where(TTGroup.C("id").Eq(TTUser.C("main_group_id")))),
		"EXISTS(SELECT group.name\nFROM group\nWHERE group.id = user.main_group_id)",
		nil,
	},
	{
		NotExists(Select(TTGroup.C("name")).From(TTGroup).Where(TTGroup.C("id").Eq(TTUser.C("main_group_id")))),
		"NOT EXISTS(SELECT group.name\nFROM group\nWHERE group.id = user.main_group_id)",
		nil,
	},
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
