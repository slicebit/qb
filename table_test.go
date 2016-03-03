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
		"\tid BIGINT,\n" +
		"\tprofile_id BIGINT,\n" +
		"\tfacebook_id BIGINT,\n" +
		"\temail VARCHAR(512) UNIQUE NOT NULL,\n" +
		"\tbio TEXT NOT NULL,\n" +
		"\tgender CHAR(16) DEFAULT 'female',\n" +
		"\tbirth_date CHAR(16) NOT NULL,\n" +
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

	query := table.Insert(kv).Query()

	assert.Equal(t, query.SQL(), "INSERT INTO user(id, full_name)\nVALUES (?, ?);")
	assert.Equal(t, query.Bindings(), []interface{}{1, "Aras Can Akin"})
}

func TestTableUpdate(t *testing.T) {

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

	query := table.
		Update(map[string]interface{}{"full_name": "Aras"}).
		Where("id = ?", 1).
		Query()

	assert.Equal(t, query.SQL(), "UPDATE user\nSET full_name = ?\nWHERE id = ?;")
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
		[]Constraint{
			Primary("id"),
		})

	query := table.
		Delete().
		Where("id = ?", 1).
		Query()

	assert.Equal(t, query.SQL(), "DELETE FROM user\nWHERE id = ?;")
	assert.Equal(t, query.Bindings(), []interface{}{1})
}
