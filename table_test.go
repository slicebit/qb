package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTable(t *testing.T) {

	table := NewTable(
		"mysql",
		"user",
		[]Column{
			NewColumn(
				"id",
				BigInt(),
				[]Constraint{},
			),
			NewColumn(
				"profile_id",
				BigInt(),
				[]Constraint{},
			),
			NewColumn(
				"facebook_id",
				BigInt(),
				[]Constraint{},
			),
			NewColumn(
				"email",
				VarChar(512),
				[]Constraint{
					Constraint{"UNIQUE"},
					NotNull(),
				},
			),
			NewColumn(
				"bio",
				Text(),
				[]Constraint{
					NotNull(),
				},
			),
			NewColumn(
				"gender",
				Char(16),
				[]Constraint{
					Default("female"),
				},
			),
			NewColumn(
				"birth_date",
				Char(16),
				[]Constraint{
					NotNull(),
				},
			),
		},
		[]Constraint{
		//	Primary("id"),
		//	Foreign("profile_id", "profile", "id"),
		},
	)

	table.AddPrimary("id")
	table.AddRef("profile_id", "profile", "id")
	table.AddRef("facebook_id", "user_facebook", "id")

	q := "CREATE TABLE `user`(\n" +
		"\t`id` BIGINT,\n" +
		"\t`profile_id` BIGINT,\n" +
		"\t`facebook_id` BIGINT,\n" +
		"\t`email` VARCHAR(512) UNIQUE NOT NULL,\n" +
		"\t`bio` TEXT NOT NULL,\n" +
		"\t`gender` CHAR(16) DEFAULT 'female',\n" +
		"\t`birth_date` CHAR(16) NOT NULL,\n" +
		"\tPRIMARY KEY (`id`),\n" +
		"\tFOREIGN KEY (`profile_id`) REFERENCES `profile`(`id`),\n" +
		"\tFOREIGN KEY (`facebook_id`) REFERENCES `user_facebook`(`id`)\n);"

	assert.Equal(t, table.SQL(), q)
}

func TestTableInsert(t *testing.T) {

	table := NewTable(
		"mysql",
		"user",
		[]Column{
			NewColumn("id", BigInt(), []Constraint{}),
			NewColumn("full_name", VarChar(), []Constraint{Unique()}),
		},
		[]Constraint{})

	table.AddPrimary("id")

	kv := map[string]interface{}{
		"id":        1,
		"full_name": "Aras Can Akin",
	}

	query := table.Insert(kv).Query()

	assert.Contains(t, query.SQL(), "INSERT INTO `user`\n(")
	assert.Contains(t, query.SQL(), "id")
	assert.Contains(t, query.SQL(), "full_name")
	assert.Contains(t, query.SQL(), ")\nVALUES (?, ?);")
	assert.Contains(t, query.Bindings(), 1)
	assert.Contains(t, query.Bindings(), "Aras Can Akin")
}

func TestTableUpdate(t *testing.T) {

	table := NewTable(
		"mysql",
		"user",
		[]Column{
			NewColumn("id", BigInt(), []Constraint{}),
			NewColumn("full_name", VarChar(), []Constraint{Unique()}),
		},
		[]Constraint{})

	table.AddPrimary("id")

	query := table.
		Update(map[string]interface{}{"full_name": "Aras"}).
		Where("id = ?", 1).
		Query()

	assert.Equal(t, query.SQL(), "UPDATE `user`\nSET full_name = ?\nWHERE id = ?;")
	assert.Equal(t, query.Bindings(), []interface{}{"Aras", 1})
}

func TestTableDelete(t *testing.T) {

	table := NewTable(
		"mysql",
		"user",
		[]Column{
			NewColumn("id", BigInt(), []Constraint{}),
			NewColumn("full_name", VarChar(), []Constraint{Unique()}),
		},
		[]Constraint{})

	table.AddPrimary("id")

	query := table.
		Delete().
		Where("id = ?", 1).
		Query()

	assert.Equal(t, query.SQL(), "DELETE FROM `user`\nWHERE id = ?;")
	assert.Equal(t, query.Bindings(), []interface{}{1})
}
