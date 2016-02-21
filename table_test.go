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
			Primary("id"),
			Foreign("profile_id", "profile", "id"),
		},
	)

	table.AddConstraint(Foreign("facebook_id", "user_facebook", "id"))

	q := "CREATE TABLE user(\n" +
		"\t`id` BIGINT,\n" +
		"\t`profile_id` BIGINT,\n" +
		"\t`facebook_id` BIGINT,\n" +
		"\t`email` VARCHAR(512) UNIQUE NOT NULL,\n" +
		"\t`bio` TEXT NOT NULL,\n" +
		"\t`gender` CHAR(16) DEFAULT 'female',\n" +
		"\t`birth_date` CHAR(16) NOT NULL,\n" +
		"\tPRIMARY KEY(id),\n" +
		"\tFOREIGN KEY (profile_id) REFERENCES profile(id),\n" +
		"\tFOREIGN KEY (facebook_id) REFERENCES user_facebook(id)\n);"

	assert.Equal(t, table.SQL(), q)
	assert.Equal(t, table.Constraints(), []Constraint{
		Primary("id"),
		Foreign("profile_id", "profile", "id"),
		Foreign("facebook_id", "user_facebook", "id"),
	})
}

func TestTableInsert(t *testing.T) {

	table := NewTable(
		"mysql",
		"user",
		[]Column{
			NewColumn("id", BigInt(), []Constraint{}),
			NewColumn("full_name", VarChar(), []Constraint{Unique()}),
		},
		[]Constraint{
			Primary("id"),
		})

	kv := map[string]interface{}{
		"id":        1,
		"full_name": "Aras Can Akin",
	}
	query, _ := table.Insert(kv)

	assert.Equal(t, query.SQL(), "INSERT INTO user(id, full_name)\nVALUES (?, ?);")
	assert.Equal(t, query.Bindings(), []interface{}{1, "Aras Can Akin"})
}
