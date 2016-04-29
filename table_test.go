package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTable(t *testing.T) {

	builder := NewBuilder("mysql")

	table := NewTable(
		builder,
		"user",
		[]Column{
			NewColumn(
				"id",
				NewType("BIGINT"),
				[]Constraint{},
			),
			NewColumn(
				"profile_id",
				NewType("BIGINT"),
				[]Constraint{},
			),
			NewColumn(
				"facebook_id",
				NewType("BIGINT"),
				[]Constraint{},
			),
			NewColumn(
				"email",
				NewType("VARCHAR(512)"),
				[]Constraint{
					Constraint{"UNIQUE"},
					NotNull(),
				},
			),
			NewColumn(
				"bio",
				NewType("TEXT"),
				[]Constraint{
					NotNull(),
				},
			),
			NewColumn(
				"gender",
				NewType("CHAR(16)"),
				[]Constraint{
					Default("female"),
				},
			),
			NewColumn(
				"birth_date",
				NewType("CHAR(16)"),
				[]Constraint{
					NotNull(),
				},
			),
		},
	)

	table.AddPrimary("id")
	table.AddRef("profile_id", "profile", "id")
	table.AddRef("facebook_id", "user_facebook", "id")

	q := "CREATE TABLE user(\n" +
		"\tid BIGINT,\n" +
		"\tprofile_id BIGINT,\n" +
		"\tfacebook_id BIGINT,\n" +
		"\temail VARCHAR(512) UNIQUE NOT NULL,\n" +
		"\tbio TEXT NOT NULL,\n" +
		"\tgender CHAR(16) DEFAULT 'female',\n" +
		"\tbirth_date CHAR(16) NOT NULL,\n" +
		"\tPRIMARY KEY (id),\n" +
		"\tFOREIGN KEY (profile_id) REFERENCES profile(id),\n" +
		"\tFOREIGN KEY (facebook_id) REFERENCES user_facebook(id)\n);"

	assert.Equal(t, table.SQL(), q)
}
