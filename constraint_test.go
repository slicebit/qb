package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConstraints(t *testing.T) {

	assert.Equal(t, Constraint("NULL"), Null())
	assert.Equal(t, Constraint("NOT NULL"), NotNull())
	assert.Equal(t, Constraint("DEFAULT '5'"), Default(5))
	assert.Equal(t, Constraint("UNIQUE"), Unique())
	assert.Equal(t, ConstraintElem{"CHECK id > 5"}, Constraint("CHECK id > 5"))
	assert.Equal(t, "NOT NULL", NotNull().String())

	sqlite := NewDialect("sqlite3")
	sqlite.SetEscaping(true)

	mysql := NewDialect("mysql")
	mysql.SetEscaping(true)

	postgres := NewDialect("postgres")
	postgres.SetEscaping(true)

	assert.Equal(t, "PRIMARY KEY(id)", PrimaryKey("id").String(sqlite))
	assert.Equal(t, "PRIMARY KEY(`id`, `email`)", PrimaryKey("id", "email").String(mysql))
	assert.Equal(t, "PRIMARY KEY(\"id\", \"email\")", PrimaryKey("id", "email").String(postgres))

	assert.Contains(t, ForeignKey("user_id").References("users", "id").String(sqlite), "FOREIGN KEY(user_id) REFERENCES users(id)")
	assert.Contains(t, ForeignKey("user_id").References("users", "id").String(mysql), "FOREIGN KEY(`user_id`) REFERENCES `users`(`id`)")
	assert.Contains(t, ForeignKey("user_id").References("users", "id").String(postgres), "FOREIGN KEY(\"user_id\") REFERENCES \"users\"(\"id\")")
	assert.Contains(t, ForeignKey("user_id", "user_email").References("users", "id", "email").String(sqlite), "FOREIGN KEY(user_id, user_email) REFERENCES users(id, email)")

	assert.Panics(t, func() {
		ForeignKey().OnUpdate("invalid")
	})
	assert.Panics(t, func() {
		ForeignKey().OnDelete("invalid")
	})
	assert.Equal(t,
		"\tFOREIGN KEY(user_id) REFERENCES users(id) ON DELETE SET NULL",
		ForeignKey("user_id").References("users", "id").OnDelete("SET NULL").String(sqlite),
	)
	assert.Equal(t,
		"\tFOREIGN KEY(user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE",
		ForeignKey("user_id").References("users", "id").OnUpdate("CASCADE").OnDelete("CASCADE").String(sqlite),
	)

	assert.Equal(t, "CONSTRAINT u_id_email UNIQUE(id, email)", UniqueKey("id", "email").String(sqlite))
	assert.Equal(t, "CONSTRAINT u_id_email UNIQUE(`id`, `email`)", UniqueKey("id", "email").String(mysql))
	assert.Equal(t, "CONSTRAINT u_id_email UNIQUE(\"id\", \"email\")", UniqueKey("id", "email").String(postgres))
}
