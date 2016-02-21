package qb

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var b *Builder

func TestBuilderInit(t *testing.T) {

	b = NewBuilder("mysql")

	query, _ := b.
		Select("id").
		From("user").
		Build()

	assert.Equal(t, query.SQL(), "SELECT id\nFROM user;")
	b.Reset()
}

func TestBuilderError(t *testing.T) {

	b.AddError(errors.New("Syntax Error"))
	assert.Equal(t, b.Errors(), []error{errors.New("Syntax Error")})
	assert.Equal(t, b.HasError(), true)

	query, err := b.
		Select("id").
		From("user").
		Build()

	assert.Equal(t, query.SQL(), "")
	assert.Equal(t, err, errors.New(strings.Join([]string{"Syntax Error"}, "\n")))
	b.Reset()
}

func TestBuilderSelectSimple(t *testing.T) {

	query, _ := b.
		Select("id", "email", "name").
		From("user").
		Build()

	assert.Equal(t, query.SQL(), "SELECT id, email, name\nFROM user;")
}

func TestBuilderSelectSingleCondition(t *testing.T) {

	query, _ := b.
		Select("id", "email", "name").
		From("user").
		Where("id = $1", 5).
		Build()

	assert.Equal(t, query.SQL(), "SELECT id, email, name\nFROM user\nWHERE id = $1;")
	assert.Equal(t, query.Bindings(), []interface{}{5})
}

func TestBuilderSelectOrderByMultiConditionWithAnd(t *testing.T) {

	query, _ := b.
		Select("id", "email", "name").
		From("user").
		Where(b.And("email = $1", "name = $2"), "a@b.c", "Aras Can Akin").
		OrderBy("email ASC, name DESC").
		Build()

	assert.Equal(t, query.SQL(), "SELECT id, email, name\nFROM user\nWHERE (email = $1 AND name = $2)\nORDER BY email ASC, name DESC;")
	assert.Equal(t, query.Bindings(), []interface{}{"a@b.c", "Aras Can Akin"})

}

func TestBuilderSelectMultiConditionWithOr(t *testing.T) {

	query, _ := b.
		Select("id", "email", "name").
		From("user").
		Where(b.Or("email = $1", "name = $2"), "a@b.c", "Aras Can Akin").
		Limit(10, 15).
		Build()

	assert.Equal(t, query.SQL(), "SELECT id, email, name\nFROM user\nWHERE email = $1 OR name = $2\nLIMIT 15 OFFSET 10;")
	assert.Equal(t, query.Bindings(), []interface{}{"a@b.c", "Aras Can Akin"})

}

func TestBuilderSelectAvgGroupByHaving(t *testing.T) {

	query, _ := b.
		Select(b.Avg("price")).
		From("products").
		GroupBy("category").
		Having(fmt.Sprintf("%s < 50", b.Max("price"))).
		Build()

	assert.Equal(t, query.SQL(), "SELECT AVG(price)\nFROM products\nGROUP BY category\nHAVING MAX(price) < 50;")
}

func TestBuilderSelectSumCount(t *testing.T) {

	query, _ := b.
		Select(b.Sum("price"), b.Count("id")).
		From("products").
		Build()

	assert.Equal(t, query.SQL(), "SELECT SUM(price), COUNT(id)\nFROM products;")
}

func TestBuilderSelectMinMax(t *testing.T) {

	query, _ := b.
		Select(b.Min("price"), b.Max("price")).
		From("products").
		Build()

	assert.Equal(t, query.SQL(), "SELECT MIN(price), MAX(price)\nFROM products;")
}

func TestBuilderSelectEqNeq(t *testing.T) {

	query, _ := b.
		Select("id", "email", "name").
		From("user").
		Where(b.And(
		b.Eq("email", "a@b.c"),
		b.NotEq("name", "Aras Can Akin"))).
		Build()

	assert.Equal(t, query.SQL(), "SELECT id, email, name\nFROM user\nWHERE (email = ? AND name != ?);")
	assert.Equal(t, query.Bindings(), []interface{}{"a@b.c", "Aras Can Akin"})
}

func TestBuilderSelectInNotIn(t *testing.T) {

	query, _ := b.
		Select("id", "email", "name").
		From("user").
		Where(b.And(
		b.In("name", "Aras Can Akin"),
		b.NotIn("email", "a@b.c"),
	)).Build()

	assert.Equal(t, query.SQL(), "SELECT id, email, name\nFROM user\nWHERE (name IN (?) AND email NOT IN (?));")
	assert.Equal(t, query.Bindings(), []interface{}{"Aras Can Akin", "a@b.c"})

}

func TestBuilderSelectGtGteStSte(t *testing.T) {

	query, _ := b.
		Select("id", "age", "avg").
		From("goqb.user").
		Where(b.And(
		b.St("age", 35),
		b.Gt("age", 18),
		b.Ste("avg", 4.0),
		b.Gte("avg", 2.8),
	)).Build()

	assert.Equal(t, query.SQL(), "SELECT id, age, avg\nFROM goqb.user\nWHERE (age < ? AND age > ? AND avg <= ? AND avg >= ?);")
	assert.Equal(t, query.Bindings(), []interface{}{35, 18, 4.0, 2.8})
}

func TestBuilderBasicInsert(t *testing.T) {

	query, _ := b.
		Insert("user", "name", "email", "password").
		Values("Aras Can Akin", "a@b.c", "p4ssw0rd").
		Build()

	assert.Equal(t, query.SQL(), "INSERT INTO user(name, email, password)\nVALUES (?, ?, ?);")
	assert.Equal(t, query.Bindings(), []interface{}{"Aras Can Akin", "a@b.c", "p4ssw0rd"})
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

	query, _ := b.
		Update("user").
		Set(
		map[string]interface{}{
			"email": "a@b.c",
			"name":  "Aras",
		}).
		Where("id = ?", 5).
		Build()

	assert.Equal(t, query.SQL(), "UPDATE user\nSET email = ?, name = ?\nWHERE id = ?;")
	assert.Equal(t, query.Bindings(), []interface{}{"a@b.c", "Aras", 5})
}

func TestBuilderDelete(t *testing.T) {

	query, _ := b.
		Delete("user").
		Where("id = ?", 5).
		Build()

	assert.Equal(t, query.SQL(), "DELETE FROM user\nWHERE id = ?;")
	assert.Equal(t, query.Bindings(), []interface{}{5})
}

func TestBuilderInnerJoin(t *testing.T) {

	query, _ := b.
		Select("id", "name", "email").
		From("user").
		InnerJoin("email", "user.id = email.id").
		Where("id = ?", 5).
		Build()

	assert.Equal(t, query.SQL(), "SELECT id, name, email\nFROM user\nINNER JOIN email ON user.id = email.id\nWHERE id = ?;")
	assert.Equal(t, query.Bindings(), []interface{}{5})
}

func TestBuilderLeftJoin(t *testing.T) {

	query, _ := b.
		Select("id", "name").
		From("user").
		LeftOuterJoin("email", "user.id = email.id").
		Where("id = ?", 5).
		Build()

	assert.Equal(t, query.SQL(), "SELECT id, name\nFROM user\nLEFT OUTER JOIN email ON user.id = email.id\nWHERE id = ?;")
	assert.Equal(t, query.Bindings(), []interface{}{5})
}

func TestBuilderRightJoin(t *testing.T) {

	query, _ := b.
		Select("id", "email_address").
		From("user").
		RightOuterJoin("email", "user.id = email.id").
		Where("id = ?", 5).
		Build()

	assert.Equal(t, query.SQL(), "SELECT id, email_address\nFROM user\nRIGHT OUTER JOIN email ON user.id = email.id\nWHERE id = ?;")
	assert.Equal(t, query.Bindings(), []interface{}{5})
}

func TestBuilderFullOuterJoin(t *testing.T) {

	query, _ := b.
		Select("id", "name", "email").
		From("user").
		FullOuterJoin("email", "user.id = email.id").
		Where("id = ?", 5).
		Build()

	assert.Equal(t, query.SQL(), "SELECT id, name, email\nFROM user\nFULL OUTER JOIN email ON user.id = email.id\nWHERE id = ?;")
	assert.Equal(t, query.Bindings(), []interface{}{5})

}

func TestBuilderCrossJoin(t *testing.T) {

	query, _ := b.
		Select("id", "name", "email").
		From("user").
		CrossJoin("email").
		Where("id = ?", 5).
		Build()

	assert.Equal(t, query.SQL(), "SELECT id, name, email\nFROM user\nCROSS JOIN email\nWHERE id = ?;")
	assert.Equal(t, query.Bindings(), []interface{}{5})
}

func TestBuilderCreateTable(t *testing.T) {

	query, _ := b.
		CreateTable("user",
		[]string{
			"id UUID PRIMARY KEY",
			"email CHAR(255) NOT NULL",
			"name VARCHAR(255) NOT NULL",
			"username VARCHAR(255) NOT NULL",
		},
		[]string{
			Constraint{"UNIQUE(email, name)"}.Name,
			Constraint{"UNIQUE(username)"}.Name,
		},
	).Build()

	qct := `CREATE TABLE user(
	id UUID PRIMARY KEY,
	email CHAR(255) NOT NULL,
	name VARCHAR(255) NOT NULL,
	username VARCHAR(255) NOT NULL,
	UNIQUE(email, name),
	UNIQUE(username)
);`
	assert.Equal(t, query.SQL(), qct)
}

func TestBuilderAlterTableAddColumn(t *testing.T) {

	query, _ := b.
		AlterTable("user").
		Add("name", "TEXT").
		Build()

	assert.Equal(t, query.SQL(), "ALTER TABLE user\nADD name TEXT;")
}

func TestBuilderAlterTableDropColumn(t *testing.T) {

	query, _ := b.
		AlterTable("user").
		Drop("name").
		Build()

	assert.Equal(t, query.SQL(), "ALTER TABLE user\nDROP name;")
}

func TestBuilderDropTable(t *testing.T) {

	query, _ := b.
		DropTable("user").
		Build()

	assert.Equal(t, query.SQL(), "DROP TABLE user;")
}
