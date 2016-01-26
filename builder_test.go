package qbit

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var builder *Builder

func TestBuilderInit(t *testing.T) {

	builder = NewBuilder()

	query, _, _ := builder.
		Select("id").
		From("user").
		Build()

	assert.Equal(t, query, "SELECT id\nFROM user;")
}

func TestBuilderSelectSimple(t *testing.T) {

	query, _, _ := builder.
		Select("id", "email", "name").
		From("user").
		Build()

	assert.Equal(t, query, "SELECT id, email, name\nFROM user;")
}

func TestBuilderSelectSingleCondition(t *testing.T) {

	query, bindings, _ := builder.
		Select("id", "email", "name").
		From("user").
		Where("id = $1", 5).
		Build()

	assert.Equal(t, query, "SELECT id, email, name\nFROM user\nWHERE id = $1;")
	assert.Equal(t, bindings, []interface{}{5})
}

func TestBuilderSelectOrderByMultiConditionWithAnd(t *testing.T) {

	query, bindings, _ := builder.
		Select("id", "email", "name").
		From("user").
		Where(builder.And("email = $1", "name = $2"), "a@b.c", "Aras Can Akin").
		OrderBy("email ASC, name DESC").
		Build()

	assert.Equal(t, query, "SELECT id, email, name\nFROM user\nWHERE (email = $1 AND name = $2)\nORDER BY email ASC, name DESC;")
	assert.Equal(t, bindings, []interface{}{"a@b.c", "Aras Can Akin"})

}

func TestBuilderSelectMultiConditionWithOr(t *testing.T) {

	query, bindings, _ := builder.
		Select("id", "email", "name").
		From("user").
		Where(builder.Or("email = $1", "name = $2"), "a@b.c", "Aras Can Akin").
		Limit(10, 15).
		Build()

	assert.Equal(t, query, "SELECT id, email, name\nFROM user\nWHERE email = $1 OR name = $2\nLIMIT 15 OFFSET 10;")
	assert.Equal(t, bindings, []interface{}{"a@b.c", "Aras Can Akin"})

}

func TestBuilderSelectAvgGroupByHaving(t *testing.T) {

	query, _, _ := builder.
		Select(builder.Avg("price")).
		From("products").
		GroupBy("category").
		Having(fmt.Sprintf("%s < 50", builder.Max("price"))).
		Build()

	assert.Equal(t, query, "SELECT AVG(price)\nFROM products\nGROUP BY category\nHAVING MAX(price) < 50;")
}

func TestBuilderSelectSumCount(t *testing.T) {

	query, _, _ := builder.
		Select(builder.Sum("price"), builder.Count("id")).
		From("products").
		Build()

	assert.Equal(t, query, "SELECT SUM(price), COUNT(id)\nFROM products;")
}

func TestBuilderSelectMinMax(t *testing.T) {

	query, _, _ := builder.
		Select(builder.Min("price"), builder.Max("price")).
		From("products").
		Build()

	assert.Equal(t, query, "SELECT MIN(price), MAX(price)\nFROM products;")
}

func TestBuilderSelectEqNeq(t *testing.T) {

	query, bindings, _ := builder.
		Select("id", "email", "name").
		From("user").
		Where(builder.And(
		builder.Eq("email", "a@b.c"),
		builder.NotEq("name", "Aras Can Akin"))).
		Build()

	assert.Equal(t, query, "SELECT id, email, name\nFROM user\nWHERE (email = ? AND name != ?);")
	assert.Equal(t, bindings, []interface{}{"a@b.c", "Aras Can Akin"})
}

func TestBuilderSelectInNotIn(t *testing.T) {

	query, bindings, _ := builder.
		Select("id", "email", "name").
		From("user").
		Where(builder.And(
		builder.In("name", "Aras Can Akin"),
		builder.NotIn("email", "a@b.c"),
	)).Build()

	assert.Equal(t, query, "SELECT id, email, name\nFROM user\nWHERE (name IN (?) AND email NOT IN (?));")
	assert.Equal(t, bindings, []interface{}{"Aras Can Akin", "a@b.c"})

}

func TestBuilderSelectGtGteStSte(t *testing.T) {

	query, bindings, _ := builder.
		Select("id", "age", "avg").
		From("goqb.user").
		Where(builder.And(
		builder.St("age", 35),
		builder.Gt("age", 18),
		builder.Ste("avg", 4.0),
		builder.Gte("avg", 2.8),
	)).Build()

	assert.Equal(t, query, "SELECT id, age, avg\nFROM goqb.user\nWHERE (age < ? AND age > ? AND avg <= ? AND avg >= ?);")
	assert.Equal(t, bindings, []interface{}{35, 18, 4.0, 2.8})
}

func TestBuilderBasicInsert(t *testing.T) {

	query, bindings, _ := builder.
		Insert("user", "name", "email", "password").
		Values("Aras Can Akin", "a@b.c", "p4ssw0rd").
		Build()

	assert.Equal(t, query, "INSERT INTO user(name, email, password)\nVALUES (?, ?, ?);")
	assert.Equal(t, bindings, []interface{}{"Aras Can Akin", "a@b.c", "p4ssw0rd"})
}

//func TestBasicUpsert(t *testing.T) {
//
//	assert := assert.New(t)
//
//	query, bindings := builder.
//		Insert("user", "name", "email").
//		Values("Aras Can Akin", "aacanakin@gmail.com").
//		UpdateOnDuplicate(map[string]interface{}{
//		"count": 2,
//	}).Build()
//
//	assert.Equal(query, "INSERT INTO user(name, email) VALUES (?, ?) ON DUPLICATE KEY UPDATE count = ?;")
//	assert.Equal(bindings, []interface{}{"Aras Can Akin", "aacanakin@gmail.com", 2})
//
//}

func TestBuilderBasicUpdate(t *testing.T) {

	query, bindings, _ := builder.
		Update("user").
		Set(
		map[string]interface{}{
			"email": "a@b.c",
			"name":  "Aras",
		}).
		Where("id = ?", 5).
		Build()

	assert.Equal(t, query, "UPDATE user\nSET email = ?, name = ?\nWHERE id = ?;")
	assert.Equal(t, bindings, []interface{}{"a@b.c", "Aras", 5})
}

func TestBuilderDelete(t *testing.T) {

	query, bindings, _ := builder.
		Delete("user").
		Where("id = ?", 5).
		Build()

	assert.Equal(t, query, "DELETE FROM user\nWHERE id = ?;")
	assert.Equal(t, bindings, []interface{}{5})
}

func TestBuilderInnerJoin(t *testing.T) {

	query, bindings, _ := builder.
		Select("id", "name", "email").
		From("user").
		InnerJoin("email", "user.id = email.id").
		Where("id = ?", 5).
		Build()

	assert.Equal(t, query, "SELECT id, name, email\nFROM user\nINNER JOIN email ON user.id = email.id\nWHERE id = ?;")
	assert.Equal(t, bindings, []interface{}{5})
}

func TestBuilderLeftJoin(t *testing.T) {

	query, bindings, _ := builder.
		Select("id", "name").
		From("user").
		LeftOuterJoin("email", "user.id = email.id").
		Where("id = ?", 5).
		Build()

	assert.Equal(t, query, "SELECT id, name\nFROM user\nLEFT OUTER JOIN email ON user.id = email.id\nWHERE id = ?;")
	assert.Equal(t, bindings, []interface{}{5})
}

func TestBuilderRightJoin(t *testing.T) {

	query, bindings, _ := builder.
		Select("id", "email_address").
		From("user").
		RightOuterJoin("email", "user.id = email.id").
		Where("id = ?", 5).
		Build()

	assert.Equal(t, query, "SELECT id, email_address\nFROM user\nRIGHT OUTER JOIN email ON user.id = email.id\nWHERE id = ?;")
	assert.Equal(t, bindings, []interface{}{5})
}

func TestBuilderFullOuterJoin(t *testing.T) {

	query, bindings, _ := builder.
		Select("id", "name", "email").
		From("user").
		FullOuterJoin("email", "user.id = email.id").
		Where("id = ?", 5).
		Build()

	assert.Equal(t, query, "SELECT id, name, email\nFROM user\nFULL OUTER JOIN email ON user.id = email.id\nWHERE id = ?;")
	assert.Equal(t, bindings, []interface{}{5})

}

func TestBuilderCrossJoin(t *testing.T) {

	query, bindings, _ := builder.
		Select("id", "name", "email").
		From("user").
		CrossJoin("email").
		Where("id = ?", 5).
		Build()

	assert.Equal(t, query, "SELECT id, name, email\nFROM user\nCROSS JOIN email\nWHERE id = ?;")
	assert.Equal(t, bindings, []interface{}{5})
}
