package qb

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSession(t *testing.T) {

	session, err := New("postgres", "user=root dbname=qb_test")
	defer session.Close()
	assert.Equal(t, session.Metadata().Engine(), session.Engine())
	assert.NotNil(t, session)
	assert.Nil(t, err)
}

func TestSessionFail(t *testing.T) {
	session, err := New("unknown", "invalid")
	assert.Nil(t, session)
	assert.NotNil(t, err)
}

func TestSessionWrappings(t *testing.T) {
	session, err := New("postgres", "user=postgres dbname=qb_test sslmode=disable")
	assert.NotNil(t, session)
	assert.Nil(t, err)

	q := session.
		Select("c1", "c2").
		From("t1").
		InnerJoin("t2", "t1.c1 = t2.tc2").
		CrossJoin("t3").
		LeftOuterJoin("t4", "t1.c1 = t4.tc4").
		RightOuterJoin("t5", "t1.c1 = t5.tc5").
		FullOuterJoin("t6", "t1.c1 = t6.tc6").
		OrderBy("t1.c1").
		GroupBy("t1.c1").
		Having("t1.c1 > 5").
		Limit(0, 1).
		Query()

	assert.Equal(t, q.SQL(), "SELECT c1, c2\nFROM \"t1\"\nINNER JOIN \"t2\" ON t1.c1 = t2.tc2\nCROSS JOIN \"t3\"\nLEFT OUTER JOIN \"t4\" ON t1.c1 = t4.tc4\nRIGHT OUTER JOIN \"t5\" ON t1.c1 = t5.tc5\nFULL OUTER JOIN \"t6\" ON t1.c1 = t6.tc6\nORDER BY t1.c1\nGROUP BY t1.c1\nHAVING t1.c1 > 5\nLIMIT 1 OFFSET 0;")

	assert.Equal(t, session.Avg("money"), "AVG(money)")
	assert.Equal(t, session.Count("money"), "COUNT(money)")
	assert.Equal(t, session.Sum("money"), "SUM(money)")
	assert.Equal(t, session.Min("money"), "MIN(money)")
	assert.Equal(t, session.Max("money"), "MAX(money)")

	assert.Equal(t, session.NotIn("email", "gmail"), "email NOT IN ($1)")
	assert.Equal(t, session.In("email", "gmail"), "email IN ($2)")
	assert.Equal(t, session.NotEq("email", "gmail"), "email != $3")
	assert.Equal(t, session.Eq("email", "gmail"), "email = $4")

	assert.Equal(t, session.Gt("money", 5), "money > $5")
	assert.Equal(t, session.Gte("money", 5), "money >= $6")
	assert.Equal(t, session.St("money", 5), "money < $7")
	assert.Equal(t, session.Ste("money", 5), "money <= $8")

	assert.Equal(t, session.And(), "")
	assert.Equal(t, session.Or(), "")

	assert.NotNil(t, session.Builder())
}
