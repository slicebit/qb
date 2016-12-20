package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUpsert(t *testing.T) {
	def := NewDialect("default")

	users := Table(
		"users",
		Column("id", Varchar().Size(36)),
		Column("email", Varchar()).Unique(),
		Column("created_at", Timestamp()).NotNull(),
		PrimaryKey("id"),
	)

	now := time.Now().UTC().String()

	ups := Upsert(users).Values(map[string]interface{}{
		"id":         "9883cf81-3b56-4151-ae4e-3903c5bc436d",
		"email":      "al@pacino.com",
		"created_at": now,
	})

	assert.Panics(t, func() {
		ups.Build(def)
	})
}
