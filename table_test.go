package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTable(t *testing.T) {

	table := NewTable(
		"user",
		[]Column{
			NewColumn(
				"id",
				BigInt(),
				[]Constraint{

				},
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
		},
	)

	q := "CREATE TABLE user(\n" +
		"\t`id` BIGINT,\n" +
		"\t`email` VARCHAR(512) UNIQUE NOT NULL,\n" +
		"\t`bio` TEXT NOT NULL,\n" +
		"\t`gender` CHAR(16) DEFAULT 'female',\n" +
		"\t`birth_date` CHAR(16) NOT NULL,\n" +
		"\tPRIMARY KEY(id)\n);"

	assert.Equal(t, table.SQL(), q)
}
