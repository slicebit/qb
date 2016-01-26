package qbit

import (
	"fmt"
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
					PrimaryKey(),
				},
			),
			NewColumn(
				"email",
				VarChar(512),
				[]Constraint{
					Unique(),
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
		[]Constraint{},
	)

	q := fmt.Sprintf(`CREATE TABLE user(
	id BIGINT PRIMARY KEY,
	email VARCHAR(512) UNIQUE NOT NULL,
	bio TEXT NOT NULL,
	gender CHAR(16) DEFAULT %s,
	birth_date CHAR(16) NOT NULL
);`, "`female`")

	assert.Equal(t, table.Sql(), q)
}