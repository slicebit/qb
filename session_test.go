package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSession(t *testing.T) {
	session, err := New("postgres", postgresDsn)
	assert.NotNil(t, session.Engine())
	defer session.Close()
	assert.NotNil(t, session)
	assert.Nil(t, err)
}

func TestSessionCommitError(t *testing.T) {
	session, err := New("postgres", postgresDsn)
	defer session.Close()
	assert.Nil(t, err)
	users := Table(
		"user",
		Column("id", BigInt().NotNull()),
	)
	ins := Insert(users).Values(map[string]interface{}{}).Build(session.Dialect())
	session.AddStatement(ins)
	err = session.Rollback()
	assert.Nil(t, err)
	session.AddStatement(ins)
	err = session.Commit()
	assert.NotNil(t, err)
}

func TestSessionAddError(t *testing.T) {
	session, err := New("postgres", postgresDsn)
	session.Dialect().SetEscaping(true)
	assert.Nil(t, err)
	type User struct {
		ID string `qb:"constraints:primary_key"`
	}
	session.AddTable(User{})
	session.DropAll()
	//err = session.CreateAll()
	//assert.Nil(t, err)
	//

	session.Close()
	defer assert.Panics(t, func() {
		session.Add(&User{ID: "hello"})
	})
}

func TestSessionFail(t *testing.T) {
	session, err := New("unknown", "invalid")
	assert.Nil(t, session)
	assert.NotNil(t, err)
}

func TestSessionWrappings(t *testing.T) {
	qb, err := New("postgres", postgresDsn)
	assert.NotNil(t, qb)
	assert.Nil(t, err)

	users := Table(
		"users",
		Column("id", Varchar().Size(36)),
		Column("name", Varchar().NotNull()),
		Column("score", BigInt().Default(0)),
		PrimaryKey("id"),
	)

	sessions := Table(
		"sessions",
		Column("id", Varchar().Size(36)),
		Column("user_id", Varchar().Size(36)),
		Column("created_at", Timestamp().NotNull()),
		PrimaryKey("id"),
		ForeignKey().Ref("user_id", "users", "id"),
	)

	qb.Metadata().AddTable(users)
	qb.Metadata().AddTable(sessions)

	selInnerJoin := qb.
		Query(sessions.C("id"), sessions.C("created_at")).
		InnerJoin(users, sessions.C("user_id"), users.C("id")).
		Filter(sessions.C("id").Eq("9efbc9ab-7914-426c-8818-7d40b0427c8f")).
		Statement()

	assert.Equal(t, "SELECT sessions.id, sessions.created_at\nFROM sessions\nINNER JOIN users ON sessions.user_id = users.id\nWHERE (sessions.id = $1);", selInnerJoin.SQL())
	assert.Equal(t, []interface{}{"9efbc9ab-7914-426c-8818-7d40b0427c8f"}, selInnerJoin.Bindings())

	selLeftJoin := qb.Query(sessions.All()...).
		LeftJoin(users, sessions.C("user_id"), users.C("id")).
		OrderBy(sessions.C("created_at")).Desc().
		Limit(0, 20).
		Filter(sessions.C("user_id").Eq("9efbc9ab-7914-426c-8818-7d40b0427c8f")).
		Filter(sessions.C("user_id").NotEq("9efbc9ac-7914-426c-8818-7d40b0427c8f")).
		Filter(sessions.C("created_at").Lte("2016-06-10")).
		Filter(sessions.C("created_at").Lt("2016-06-10")).
		Filter(sessions.C("created_at").Gte("2016-06-09")).
		Filter(sessions.C("created_at").Gt("2016-06-09")).
		Statement()

	assert.Contains(t, selLeftJoin.SQL(), "id")
	assert.Contains(t, selLeftJoin.SQL(), "user_id")
	assert.Contains(t, selLeftJoin.SQL(), "created_at")
	assert.Contains(t, selLeftJoin.SQL(), "\nFROM sessions\nLEFT OUTER JOIN users ON sessions.user_id = users.id\nWHERE (sessions.user_id = $1 AND sessions.user_id != $2 AND sessions.created_at <= $3 AND sessions.created_at < $4 AND sessions.created_at >= $5 AND sessions.created_at > $6)\nORDER BY created_at DESC\nLIMIT 20 OFFSET 0;")
	assert.Equal(t, []interface{}{"9efbc9ab-7914-426c-8818-7d40b0427c8f", "9efbc9ac-7914-426c-8818-7d40b0427c8f", "2016-06-10", "2016-06-10", "2016-06-09", "2016-06-09"}, selLeftJoin.Bindings())

	selRightJoin := qb.Query(sessions.C("id"), sessions.C("user_id"), sessions.C("created_at")).
		RightJoin(users, sessions.C("user_id"), users.C("id")).
		OrderBy(sessions.C("created_at")).Desc().
		Limit(0, 20).
		Filter(sessions.C("user_id").In("9efbc9ab-7914-426c-8818-7d40b0427c8f")).
		Statement()

	assert.Equal(t, "SELECT sessions.id, sessions.user_id, sessions.created_at\nFROM sessions\nRIGHT OUTER JOIN users ON sessions.user_id = users.id\nWHERE (sessions.user_id IN ($1))\nORDER BY created_at DESC\nLIMIT 20 OFFSET 0;", selRightJoin.SQL())
	assert.Equal(t, []interface{}{"9efbc9ab-7914-426c-8818-7d40b0427c8f"}, selRightJoin.Bindings())

	selCrossJoin := qb.
		Query(sessions.C("id"), sessions.C("user_id"), sessions.C("created_at")).
		CrossJoin(users).
		OrderBy(sessions.C("created_at")).Asc().
		Limit(0, 20).
		Filter(sessions.C("user_id").NotIn("9efbc9ab-7914-426c-8818-7d40b0427c8f")).
		Statement()

	assert.Equal(t, "SELECT sessions.id, sessions.user_id, sessions.created_at\nFROM sessions\nCROSS JOIN users\nWHERE (sessions.user_id NOT IN ($1))\nORDER BY created_at ASC\nLIMIT 20 OFFSET 0;", selCrossJoin.SQL())
	assert.Equal(t, []interface{}{"9efbc9ab-7914-426c-8818-7d40b0427c8f"}, selCrossJoin.Bindings())

	selLike := qb.
		Query(users.C("id"), users.C("name")).
		Filter(users.C("name").Like("%Robert%")).
		Statement()

	assert.Equal(t, "SELECT id, name\nFROM users\nWHERE (users.name LIKE '%Robert%');", selLike.SQL())
	assert.Equal(t, []interface{}{}, selLike.Bindings())

	selAggCountMinMax := qb.
		Query(Count(users.C("id")), Max(users.C("name")), Min(users.C("name"))).
		From(qb.T("users")).
		GroupBy(users.C("name")).
		Having(Sum(users.C("score")), ">", 100).
		Statement()

	assert.Equal(t, "SELECT COUNT(id), MAX(name), MIN(name)\nFROM users\nGROUP BY name\nHAVING SUM(score) > $1;", selAggCountMinMax.SQL())
	assert.Equal(t, []interface{}{100}, selAggCountMinMax.Bindings())

	selAggAvgSum := qb.
		Query(Avg(users.C("score")), Sum(users.C("score"))).
		GroupBy(users.C("id")).
		Statement()

	assert.Equal(t, "SELECT AVG(score), SUM(score)\nFROM \nGROUP BY id;", selAggAvgSum.SQL())
	assert.Equal(t, []interface{}{}, selAggAvgSum.Bindings())

	assert.Panics(t, func() {
		qb.Query()
	})

	assert.Panics(t, func() {
		qb.From(users)
	})

	assert.Panics(t, func() {
		qb.Filter(users.C("id").Eq(""))
	})

	assert.Panics(t, func() {
		qb.InnerJoin(sessions, users.C("id"), sessions.C("user_id"))
	})

	assert.Panics(t, func() {
		qb.LeftJoin(sessions, users.C("id"), sessions.C("user_id"))
	})

	assert.Panics(t, func() {
		qb.RightJoin(sessions, users.C("id"), sessions.C("user_id"))
	})

	assert.Panics(t, func() {
		qb.CrossJoin(users)
	})

	assert.Panics(t, func() {
		qb.GroupBy(users.C("id"))
	})

	assert.Panics(t, func() {
		qb.Having(Sum(users.C("score")), ">", 100)
	})

	assert.Panics(t, func() {
		qb.OrderBy(sessions.C("created_at"))
	})

	assert.Panics(t, func() {
		qb.Asc()
	})

	assert.Panics(t, func() {
		qb.Desc()
	})

	assert.Panics(t, func() {
		qb.Limit(0, 20)
	})
}
