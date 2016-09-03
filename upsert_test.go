package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUpsert(t *testing.T) {
	sqlite := NewDialect("sqlite3")
	sqlite.SetEscaping(true)

	mysql := NewDialect("mysql")
	mysql.SetEscaping(true)

	postgres := NewDialect("postgres")
	postgres.SetEscaping(true)

	def := NewDialect("default")

	users := Table(
		"users",
		Column("id", Varchar().Size(36)),
		Column("email", Varchar()).Unique(),
		Column("created_at", Timestamp()).NotNull(),
		PrimaryKey("id"),
	)

	var statement *Stmt

	now := time.Now().UTC().String()

	ups := Upsert(users).Values(map[string]interface{}{
		"id":         "9883cf81-3b56-4151-ae4e-3903c5bc436d",
		"email":      "al@pacino.com",
		"created_at": now,
	})

	assert.Panics(t, func() {
		ups.Build(def)
	})

	statement = ups.Build(sqlite)
	assert.Contains(t, statement.SQL(), "REPLACE INTO users")
	assert.Contains(t, statement.SQL(), "id", "email", "created_at")
	assert.Contains(t, statement.SQL(), "VALUES(?, ?, ?)")
	assert.Contains(t, statement.Bindings(), "9883cf81-3b56-4151-ae4e-3903c5bc436d")
	assert.Contains(t, statement.Bindings(), "al@pacino.com")
	assert.Contains(t, statement.Bindings(), now)
	assert.Equal(t, 3, len(statement.Bindings()))

	statement = ups.Build(mysql)
	assert.Contains(t, statement.SQL(), "INSERT INTO `users`")
	assert.Contains(t, statement.SQL(), "`id`", "`email`", "`created_at`")
	assert.Contains(t, statement.SQL(), "VALUES(?, ?, ?)")
	assert.Contains(t, statement.SQL(), "ON DUPLICATE KEY UPDATE")
	assert.Contains(t, statement.SQL(), "`id` = ?", "`email` = ?", "`created_at` = ?")
	assert.Contains(t, statement.Bindings(), "9883cf81-3b56-4151-ae4e-3903c5bc436d")
	assert.Contains(t, statement.Bindings(), "al@pacino.com")
	assert.Equal(t, 6, len(statement.Bindings()))

	statement = ups.Build(postgres)
	assert.Contains(t, statement.SQL(), "INSERT INTO \"users\"")
	assert.Contains(t, statement.SQL(), "\"id\"", "\"email\"")
	assert.Contains(t, statement.SQL(), "VALUES($1, $2, $3)")
	assert.Contains(t, statement.SQL(), "ON CONFLICT", "DO UPDATE SET")
	assert.Contains(t, statement.Bindings(), "9883cf81-3b56-4151-ae4e-3903c5bc436d")
	assert.Contains(t, statement.Bindings(), "al@pacino.com")
	assert.Equal(t, 6, len(statement.Bindings()))

	statement = Upsert(users).
		Values(map[string]interface{}{
			"id":    "9883cf81-3b56-4151-ae4e-3903c5bc436d",
			"email": "al@pacino.com",
		}).
		Returning(users.C("id"), users.C("email")).
		Build(postgres)

	assert.Contains(t, statement.SQL(), "INSERT INTO \"users\"")
	assert.Contains(t, statement.SQL(), "\"id\"", "\"email\"")
	assert.Contains(t, statement.SQL(), "ON CONFLICT")
	assert.Contains(t, statement.SQL(), "DO UPDATE SET")
	assert.Contains(t, statement.SQL(), "VALUES($1, $2)")
	assert.Contains(t, statement.SQL(), "RETURNING \"id\", \"email\";")
	assert.Contains(t, statement.Bindings(), "9883cf81-3b56-4151-ae4e-3903c5bc436d")
	assert.Contains(t, statement.Bindings(), "al@pacino.com")
	assert.Equal(t, 4, len(statement.Bindings()))
}
