package qb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var emptyBinds = []interface{}{}

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
		ForeignKey("main_group_id").References("group", "id"),
	)
)

var compileTests = []struct {
	clause Clause
	expect string
	binds  []interface{}
}{
	{SQLText("1"), "1", emptyBinds},
	{
		Join("LEFT JOIN", TTGroup, TTUser),
		"group\nLEFT JOIN user ON user.main_group_id = group.id",
		emptyBinds,
	},
	{
		Join("LEFT JOIN", TTGroup, TTUser, TTGroup.C("id").Eq(TTUser.C("id"))),
		"group\nLEFT JOIN user ON group.id = user.id",
		emptyBinds,
	},
	{
		Join("LEFT JOIN", TTGroup, TTUser, TTGroup.C("id"), TTUser.C("id")),
		"group\nLEFT JOIN user ON group.id = user.id",
		emptyBinds,
	},
	{
		Exists(Select(TTGroup.C("name")).From(TTGroup).Where(TTGroup.C("id").Eq(TTUser.C("main_group_id")))),
		"EXISTS(SELECT group.name\nFROM group\nWHERE group.id = user.main_group_id)",
		emptyBinds,
	},
	{
		NotExists(Select(TTGroup.C("name")).From(TTGroup).Where(TTGroup.C("id").Eq(TTUser.C("main_group_id")))),
		"NOT EXISTS(SELECT group.name\nFROM group\nWHERE group.id = user.main_group_id)",
		emptyBinds,
	},
	{
		Select(Exists(Select(SQLText("1")).From(TTGroup).Where(TTGroup.C("id").Eq(TTUser.C("main_group_id"))))),
		"SELECT EXISTS(SELECT 1\nFROM group\nWHERE group.id = user.main_group_id)",
		emptyBinds,
	},
	{
		Select(SQLText("1")).From(TTGroup).ForUpdate(),
		"SELECT 1\nFROM group\nFOR UPDATE",
		emptyBinds,
	},
	{
		Select(SQLText("1")).From(TTGroup).ForUpdate(TTUser, TTGroup),
		"SELECT 1\nFROM group\nFOR UPDATE OF user, group",
		emptyBinds,
	},
}

func TestCompile(t *testing.T) {
	compile := func(clause Clause) (string, []interface{}) {
		context := NewCompilerContext(NewDialect("default"))
		return clause.Accept(context), context.Binds()
	}

	for _, tt := range compileTests {
		actual, binds := compile(tt.clause)
		assert.Equal(t, tt.expect, actual)
		assert.Equal(t, tt.binds, binds)
	}
}
