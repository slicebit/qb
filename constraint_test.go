package qb

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestConstraints(t *testing.T) {

	assert.Equal(t, Null(), Constraint("NULL"))
	assert.Equal(t, NotNull(), Constraint("NOT NULL"))
	assert.Equal(t, Default(5), Constraint("DEFAULT '5'"))
	assert.Equal(t, Unique(), Constraint("UNIQUE"))
	assert.Equal(t, Constraint("CHECK id > 5"), ConstraintElem{"CHECK id > 5"})
	assert.Equal(t, NotNull().String(), "NOT NULL")

	sqlite := NewDialect("sqlite3")
	sqlite.SetEscaping(true)

	mysql := NewDialect("mysql")
	mysql.SetEscaping(true)

	postgres := NewDialect("postgres")
	postgres.SetEscaping(true)

	assert.Equal(t, PrimaryKey("id").String(sqlite), "PRIMARY KEY(id)")
	assert.Equal(t, PrimaryKey("id", "email").String(mysql), "PRIMARY KEY(`id`, `email`)")
	assert.Equal(t, PrimaryKey("id", "email").String(postgres), "PRIMARY KEY(\"id\", \"email\")")

	assert.Contains(t, ForeignKey().Ref("user_id", "users", "id").String(sqlite), "FOREIGN KEY(user_id) REFERENCES users(id)")
	assert.Contains(t, ForeignKey().Ref("user_id", "users", "id").String(mysql), "FOREIGN KEY(`user_id`) REFERENCES `users`(`id`)")
	assert.Contains(t, ForeignKey().Ref("user_id", "users", "id").String(postgres), "FOREIGN KEY(\"user_id\") REFERENCES \"users\"(\"id\")")
	assert.Contains(t, ForeignKey().Ref("user_id", "users", "id").Ref("user_email", "users", "email").String(sqlite), "FOREIGN KEY(user_id, user_email) REFERENCES users(id, email)")

	assert.Equal(t, UniqueKey("id", "email").String(sqlite), "CONSTRAINT u_id_email UNIQUE(id, email)")
	assert.Equal(t, UniqueKey("id", "email").String(mysql), "CONSTRAINT u_id_email UNIQUE(`id`, `email`)")
	assert.Equal(t, UniqueKey("id", "email").String(postgres), "CONSTRAINT u_id_email UNIQUE(\"id\", \"email\")")
}
