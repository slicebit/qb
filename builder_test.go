package qb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type BuilderTestSuite struct {
	suite.Suite
	builder *Builder
}

func (suite *BuilderTestSuite) SetupTest() {
	suite.builder = NewBuilder("mysql")
}

func (suite *BuilderTestSuite) TestBuilderInit() {

	query := suite.builder.
		Select("id").
		From("user").
		Query()

	assert.Equal(suite.T(), query.SQL(), "SELECT id\nFROM user;")
}

func (suite *BuilderTestSuite) TestBuilderSelectSimple() {

	query := suite.builder.
		Select("id", "email", "name").
		From("user").
		Where("").
		Query()

	assert.Equal(suite.T(), query.SQL(), "SELECT id, email, name\nFROM user;")
}

func (suite *BuilderTestSuite) TestBuilderEmptyAnd() {
	assert.Equal(suite.T(), suite.builder.And(), "")
}

func (suite *BuilderTestSuite) TestBuilderSelectSingleCondition() {

	query := suite.builder.
		Select("id", "email", "name").
		From("user").
		Where("id = ?", 5).
		Query()

	assert.Equal(suite.T(), query.SQL(), "SELECT id, email, name\nFROM user\nWHERE id = ?;")
	assert.Equal(suite.T(), query.Bindings(), []interface{}{5})
}

func (suite *BuilderTestSuite) TestBuilderSelectOrderByMultiConditionWithAnd() {

	query := suite.builder.
		Select("id", "email", "name").
		From("user").
		Where(suite.builder.And("email = ?", "name = ?"), "a@b.c", "Aras Can Akin").
		OrderBy("email ASC, name DESC").
		Query()

	assert.Equal(suite.T(), query.SQL(), "SELECT id, email, name\nFROM user\nWHERE (email = ? AND name = ?)\nORDER BY email ASC, name DESC;")
	assert.Equal(suite.T(), query.Bindings(), []interface{}{"a@b.c", "Aras Can Akin"})

}

func (suite *BuilderTestSuite) TestBuilderSelectMultiConditionWithOr() {

	query := suite.builder.
		Select("id", "email", "name").
		From("user").
		Where(suite.builder.Or("email = $1", "name = $2"), "a@b.c", "Aras Can Akin").
		Limit(10, 15).
		Query()

	assert.Equal(suite.T(), query.SQL(), "SELECT id, email, name\nFROM user\nWHERE email = $1 OR name = $2\nLIMIT 15 OFFSET 10;")
	assert.Equal(suite.T(), query.Bindings(), []interface{}{"a@b.c", "Aras Can Akin"})
}

func (suite *BuilderTestSuite) TestBuilderSelectAvgGroupByHaving() {

	query := suite.builder.
		Select(suite.builder.Avg("price")).
		From("products").
		GroupBy("category").
		Having(fmt.Sprintf("%s < 50", suite.builder.Max("price"))).
		Query()

	assert.Equal(suite.T(), query.SQL(), "SELECT AVG(price)\nFROM products\nGROUP BY category\nHAVING MAX(price) < 50;")
}

func (suite *BuilderTestSuite) TestBuilderSelectSumCount() {

	query := suite.builder.
		Select(suite.builder.Sum("price"), suite.builder.Count("id")).
		From("products").
		Query()

	assert.Equal(suite.T(), query.SQL(), "SELECT SUM(price), COUNT(id)\nFROM products;")
}

func (suite *BuilderTestSuite) TestBuilderSelectMinMax() {

	query := suite.builder.
		Select(suite.builder.Min("price"), suite.builder.Max("price")).
		From("products").
		Query()

	assert.Equal(suite.T(), query.SQL(), "SELECT MIN(price), MAX(price)\nFROM products;")
}

func (suite *BuilderTestSuite) TestBuilderSelectEqNeq() {

	query := suite.builder.
		Select("id", "email", "name").
		From("user").
		Where(suite.builder.And(
			suite.builder.Eq("email", "a@b.c"),
			suite.builder.NotEq("name", "Aras Can Akin"))).
		Query()

	assert.Equal(suite.T(), query.SQL(), "SELECT id, email, name\nFROM user\nWHERE (email = ? AND name != ?);")
	assert.Equal(suite.T(), query.Bindings(), []interface{}{"a@b.c", "Aras Can Akin"})
}

func (suite *BuilderTestSuite) TestBuilderSelectInNotIn() {

	query := suite.builder.
		Select("id", "email", "name").
		From("user").
		Where(suite.builder.And(
			suite.builder.In("name", "Aras Can Akin"),
			suite.builder.NotIn("email", "a@b.c"),
		)).Query()

	assert.Equal(suite.T(), query.SQL(), "SELECT id, email, name\nFROM user\nWHERE (name IN (?) AND email NOT IN (?));")
	assert.Equal(suite.T(), query.Bindings(), []interface{}{"Aras Can Akin", "a@b.c"})

}

func (suite *BuilderTestSuite) TestBuilderSelectGtGteStSte() {

	query := suite.builder.
		Select("id", "age", "avg").
		From("goqb.user").
		Where(suite.builder.And(
			suite.builder.St("age", 35),
			suite.builder.Gt("age", 18),
			suite.builder.Ste("avg", 4.0),
			suite.builder.Gte("avg", 2.8),
		)).Query()

	assert.Equal(suite.T(), query.SQL(), "SELECT id, age, avg\nFROM goqb.user\nWHERE (age < ? AND age > ? AND avg <= ? AND avg >= ?);")
	assert.Equal(suite.T(), query.Bindings(), []interface{}{35, 18, 4.0, 2.8})
}

func (suite *BuilderTestSuite) TestBuilderBasicInsert() {

	fields := map[string]interface{}{
		"name":     "Aras Can Akin",
		"email":    "a@b.c",
		"password": "p4ssw0rd",
	}

	query := suite.builder.
		Insert("user").
		Values(fields).
		Returning("email").
		Query()

	assert.Contains(suite.T(), query.SQL(), "INSERT INTO user\n(")
	assert.Contains(suite.T(), query.SQL(), "name")
	assert.Contains(suite.T(), query.SQL(), "email")
	assert.Contains(suite.T(), query.SQL(), "password")
	assert.Contains(suite.T(), query.SQL(), "\nVALUES (?, ?, ?)")
	assert.Contains(suite.T(), query.SQL(), "RETURNING email;")
	assert.Contains(suite.T(), query.Bindings(), "Aras Can Akin")
	assert.Contains(suite.T(), query.Bindings(), "a@b.c")
	assert.Contains(suite.T(), query.Bindings(), "p4ssw0rd")
}

func (suite *BuilderTestSuite) TestBuilderBasicUpdate() {

	query := suite.builder.
		Update("user").
		Set(
			map[string]interface{}{
				"email": "a@b.c",
				"name":  "Aras",
			}).
		Where("id = ?", 5).
		Query()

	assert.Contains(suite.T(), query.SQL(), "UPDATE user\nSET")
	assert.Contains(suite.T(), query.SQL(), "email = ?")
	assert.Contains(suite.T(), query.SQL(), "name = ?")
	assert.Contains(suite.T(), query.SQL(), "WHERE id = ?;")
	assert.Contains(suite.T(), query.Bindings(), "a@b.c")
	assert.Contains(suite.T(), query.Bindings(), "Aras")
	assert.Contains(suite.T(), query.Bindings(), 5)
}

func (suite *BuilderTestSuite) TestBuilderDelete() {

	query := suite.builder.
		Delete("user").
		Where("id = ?", 5).
		Query()

	assert.Equal(suite.T(), query.SQL(), "DELETE FROM user\nWHERE id = ?;")
	assert.Equal(suite.T(), query.Bindings(), []interface{}{5})
}

func (suite *BuilderTestSuite) TestBuilderInnerJoin() {

	query := suite.builder.
		Select("id", "name", "email").
		From("user").
		InnerJoin("email", "user.id = email.id").
		Where("id = ?", 5).
		Query()

	assert.Equal(suite.T(), query.SQL(), "SELECT id, name, email\nFROM user\nINNER JOIN email ON user.id = email.id\nWHERE id = ?;")
	assert.Equal(suite.T(), query.Bindings(), []interface{}{5})
}

func (suite *BuilderTestSuite) TestBuilderLeftJoin() {

	query := suite.builder.
		Select("id", "name").
		From("user").
		LeftOuterJoin("email e", "user.id = e.id").
		Where("id = ?", 5).
		Query()

	assert.Equal(suite.T(), query.SQL(), "SELECT id, name\nFROM user\nLEFT OUTER JOIN email e ON user.id = e.id\nWHERE id = ?;")
	assert.Equal(suite.T(), query.Bindings(), []interface{}{5})
}

func (suite *BuilderTestSuite) TestBuilderRightJoin() {

	query := suite.builder.
		Select("id", "email_address").
		From("user").
		RightOuterJoin("email e", "user.id = e.id").
		Where("id = ?", 5).
		Query()

	assert.Equal(suite.T(), query.SQL(), "SELECT id, email_address\nFROM user\nRIGHT OUTER JOIN email e ON user.id = e.id\nWHERE id = ?;")
	assert.Equal(suite.T(), query.Bindings(), []interface{}{5})
}

func (suite *BuilderTestSuite) TestBuilderFullOuterJoin() {

	query := suite.builder.
		Select("id", "name", "email").
		From("user").
		FullOuterJoin("email e", "user.id = e.id").
		Where("id = ?", 5).
		Query()

	assert.Equal(suite.T(), query.SQL(), "SELECT id, name, email\nFROM user\nFULL OUTER JOIN email e ON user.id = e.id\nWHERE id = ?;")
	assert.Equal(suite.T(), query.Bindings(), []interface{}{5})

}

func (suite *BuilderTestSuite) TestBuilderCrossJoin() {

	query := suite.builder.
		Select("id", "name", "email").
		From("user").
		CrossJoin("email e").
		Where("id = ?", 5).
		Query()

	assert.Equal(suite.T(), query.SQL(), "SELECT id, name, email\nFROM user\nCROSS JOIN email e\nWHERE id = ?;")
	assert.Equal(suite.T(), query.Bindings(), []interface{}{5})
}

func (suite *BuilderTestSuite) TestBuilderCreateTable() {

	query := suite.builder.
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
		).Query()

	qct := fmt.Sprintf(`CREATE TABLE %s(
	id UUID PRIMARY KEY,
	email CHAR(255) NOT NULL,
	name VARCHAR(255) NOT NULL,
	username VARCHAR(255) NOT NULL,
	UNIQUE(email, name),
	UNIQUE(username)
);`, "user")
	assert.Equal(suite.T(), query.SQL(), qct)
}

func (suite *BuilderTestSuite) TestBuilderAlterTableAddColumn() {

	query := suite.builder.
		AlterTable("user").
		Add("name", "TEXT").
		Query()

	assert.Equal(suite.T(), query.SQL(), "ALTER TABLE user\nADD name TEXT;")
}

func (suite *BuilderTestSuite) TestBuilderAlterTableDropColumn() {

	query := suite.builder.
		AlterTable("user").
		Drop("name").
		Query()

	assert.Equal(suite.T(), query.SQL(), "ALTER TABLE user\nDROP name;")
}

func (suite *BuilderTestSuite) TestBuilderDropTable() {

	query := suite.builder.
		DropTable("user").
		Query()

	assert.Equal(suite.T(), query.SQL(), "DROP TABLE user;")
}

func TestBuilderSuite(t *testing.T) {
	suite.Run(t, new(BuilderTestSuite))
}
