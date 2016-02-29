package qb

import (
	"fmt"
	"strings"
)

// Dialect is the abstraction of set of any expression api functions
type Dialect interface {
	Query() *Query

	// clauses
	Insert(table string, columns ...string) *commonDialect
	Values(values ...interface{}) *commonDialect
	Update(table string) *commonDialect
	Set(params map[string]interface{}) *commonDialect
	Delete(table string) *commonDialect
	Select(columns ...string) *commonDialect
	From(tables ...string) *commonDialect
	InnerJoin(table string, expressions ...string) *commonDialect
	CrossJoin(table string) *commonDialect
	LeftOuterJoin(table string, expressions ...string) *commonDialect
	RightOuterJoin(table string, expressions ...string) *commonDialect
	FullOuterJoin(table string, expressions ...string) *commonDialect
	Where(expression string, bindings ...interface{}) *commonDialect
	OrderBy(expressions ...string) *commonDialect
	GroupBy(columns ...string) *commonDialect
	Having(expressions ...string) *commonDialect
	Limit(offset int, count int) *commonDialect

	// aggregates
	Avg(column string) string
	Count(column string) string
	Sum(column string) string
	Min(column string) string
	Max(column string) string

	// comparators
	In(key string, values ...interface{}) string
	NotIn(key string, values ...interface{}) string
	Eq(key string, value interface{}) string
	NotEq(key string, value interface{}) string
	Gt(key string, value interface{}) string
	Gte(key string, value interface{}) string
	St(key string, value interface{}) string
	Ste(key string, value interface{}) string
	And(expressions ...string) string
	Or(expressions ...string) string

	// table crudders
	CreateTable(table string, fields []string, constraints []string) *commonDialect
	AlterTable(table string) *commonDialect
	Add(colName string, colType string) *commonDialect
	Drop(colName string) *commonDialect
	DropTable(table string) *commonDialect
}

// NewBuilder generates a new builder object
func NewDialect(driver string) Dialect {
	switch driver {
	case "sqlite":
		return SqliteDialect{
			&commonDialect{query: NewQuery()},
		}

	case "mysql":
		return MysqlDialect{
			&commonDialect{query: NewQuery()},
		}
	case "postgres":
		return PostgresDialect{
			&commonDialect{query: NewQuery()},
			0,
		}
	default:
		panic(fmt.Errorf("Unsupported Driver: %s", driver))
	}
}

// Builder is a subset of dialect could be used for common sql queries
// it has all the common functions except multiple statements & table crudders
type commonDialect struct {
	query *Query
}

// Reset clears query bindings and its errors
func (cd *commonDialect) Reset() {
	cd.query = NewQuery()
}

// Build generates sql query and bindings from clauses and bindings.
// The query clauses and returns the sql and bindings
func (cd *commonDialect) Query() *Query {
	query := cd.query
	cd.Reset()
	return query
}

func (cd *commonDialect) placeholder() string {
	return "?"
}

func (cd *commonDialect) placeholders(values ...interface{}) []string {
	placeholders := make([]string, len(values))
	for k := range values {
		placeholders[k] = cd.placeholder()
	}
	return placeholders
}

// Insert generates an "insert into %s(%s)" statement
func (cd *commonDialect) Insert(table string, columns ...string) *commonDialect {
	clause := fmt.Sprintf("INSERT INTO %s(%s)", table, strings.Join(columns, ", "))
	cd.query.AddClause(clause)
	return cd
}

// Values generates "values(%s)" statement and add bindings for each value
func (cd *commonDialect) Values(values ...interface{}) *commonDialect {
	cd.query.AddBinding(values...)
	clause := fmt.Sprintf("VALUES (%s)", strings.Join(cd.placeholders(values...), ", "))
	cd.query.AddClause(clause)
	return cd
}

// Update generates "update %s" statement
func (cd *commonDialect) Update(table string) *commonDialect {
	clause := fmt.Sprintf("UPDATE %s", table)
	cd.query.AddClause(clause)
	return cd
}

// Set generates "set a = placeholder" statement for each key a and add bindings for map value
func (cd *commonDialect) Set(m map[string]interface{}) *commonDialect {
	updates := []string{}
	for k, v := range m {
		updates = append(updates, fmt.Sprintf("%s = %s", k, cd.placeholder()))
		cd.query.AddBinding(v)
	}
	clause := fmt.Sprintf("SET %s", strings.Join(updates, ", "))
	cd.query.AddClause(clause)
	return cd
}

// Delete generates "delete" statement
func (cd *commonDialect) Delete(table string) *commonDialect {
	cd.query.AddClause(fmt.Sprintf("DELETE FROM %s", table))
	return cd
}

// Select generates "select %s" statement
func (cd *commonDialect) Select(columns ...string) *commonDialect {
	clause := fmt.Sprintf("SELECT %s", strings.Join(columns, ", "))
	cd.query.AddClause(clause)
	return cd
}

// From generates "from %s" statement for each table name
func (cd *commonDialect) From(tables ...string) *commonDialect {
	cd.query.AddClause(fmt.Sprintf("FROM %s", strings.Join(tables, ", ")))
	return cd
}

// InnerJoin generates "inner join %s on %s" statement for each expression
func (cd *commonDialect) InnerJoin(table string, expressions ...string) *commonDialect {
	cd.query.AddClause(fmt.Sprintf("INNER JOIN %s ON %s", table, strings.Join(expressions, " ")))
	return cd
}

// CrossJoin generates "cross join %s" statement for table
func (cd *commonDialect) CrossJoin(table string) *commonDialect {
	cd.query.AddClause(fmt.Sprintf("CROSS JOIN %s", table))
	return cd
}

// LeftOuterJoin generates "left outer join %s on %s" statement for each expression
func (cd *commonDialect) LeftOuterJoin(table string, expressions ...string) *commonDialect {
	cd.query.AddClause(fmt.Sprintf("LEFT OUTER JOIN %s ON %s", table, strings.Join(expressions, " ")))
	return cd
}

// RightOuterJoin generates "right outer join %s on %s" statement for each expression
func (cd *commonDialect) RightOuterJoin(table string, expressions ...string) *commonDialect {
	cd.query.AddClause(fmt.Sprintf("RIGHT OUTER JOIN %s ON %s", table, strings.Join(expressions, " ")))
	return cd
}

// FullOuterJoin generates "full outer join %s on %s" for each expression
func (cd *commonDialect) FullOuterJoin(table string, expressions ...string) *commonDialect {
	cd.query.AddClause(fmt.Sprintf("FULL OUTER JOIN %s ON %s", table, strings.Join(expressions, " ")))
	return cd
}

// Where generates "where %s" for the expression and adds bindings for each value
func (cd *commonDialect) Where(expression string, bindings ...interface{}) *commonDialect {
	expression = strings.Replace(expression, "?", cd.placeholder(), -1)
	cd.query.AddClause(fmt.Sprintf("WHERE %s", expression))
	cd.query.AddBinding(bindings...)
	return cd
}

// OrderBy generates "order by %s" for each expression
func (cd *commonDialect) OrderBy(expressions ...string) *commonDialect {
	cd.query.AddClause(fmt.Sprintf("ORDER BY %s", strings.Join(expressions, ", ")))
	return cd
}

// GroupBy generates "group by %s" for each column
func (cd *commonDialect) GroupBy(columns ...string) *commonDialect {
	cd.query.AddClause(fmt.Sprintf("GROUP BY %s", strings.Join(columns, ", ")))
	return cd
}

// Having generates "having %s" for each expression
func (cd *commonDialect) Having(expressions ...string) *commonDialect {
	cd.query.AddClause(fmt.Sprintf("HAVING %s", strings.Join(expressions, ", ")))
	return cd
}

// Limit generates limit %d offset %d for offset and count
func (cd *commonDialect) Limit(offset int, count int) *commonDialect {
	cd.query.AddClause(fmt.Sprintf("LIMIT %d OFFSET %d", count, offset))
	return cd
}

// aggregates

// Avg function generates "avg(%s)" statement for column
func (cd *commonDialect) Avg(column string) string {
	return fmt.Sprintf("AVG(%s)", column)
}

// Count function generates "count(%s)" statement for column
func (cd *commonDialect) Count(column string) string {
	return fmt.Sprintf("COUNT(%s)", column)
}

// Sum function generates "sum(%s)" statement for column
func (cd *commonDialect) Sum(column string) string {
	return fmt.Sprintf("SUM(%s)", column)
}

// Min function generates "min(%s)" statement for column
func (cd *commonDialect) Min(column string) string {
	return fmt.Sprintf("MIN(%s)", column)
}

// Max function generates "max(%s)" statement for column
func (cd *commonDialect) Max(column string) string {
	return fmt.Sprintf("MAX(%s)", column)
}

// expressions

// NotIn function generates "%s not in (%s)" for key and adds bindings for each value
func (cd *commonDialect) NotIn(key string, values ...interface{}) string {
	cd.query.AddBinding(values...)
	return fmt.Sprintf("%s NOT IN (%s)", key, strings.Join(cd.placeholders(values...), ","))
}

// In function generates "%s in (%s)" for key and adds bindings for each value
func (cd *commonDialect) In(key string, values ...interface{}) string {
	cd.query.AddBinding(values...)
	return fmt.Sprintf("%s IN (%s)", key, strings.Join(cd.placeholders(values...), ","))
}

// NotEq function generates "%s != placeholder" for key and adds binding for value
func (cd *commonDialect) NotEq(key string, value interface{}) string {
	cd.query.AddBinding(value)
	return fmt.Sprintf("%s != %s", key, cd.placeholder())
}

// Eq function generates "%s = placeholder" for key and adds binding for value
func (cd *commonDialect) Eq(key string, value interface{}) string {
	cd.query.AddBinding(value)
	return fmt.Sprintf("%s = %s", key, cd.placeholder())
}

// Gt function generates "%s > placeholder" for key and adds binding for value
func (cd *commonDialect) Gt(key string, value interface{}) string {
	cd.query.AddBinding(value)
	return fmt.Sprintf("%s > %s", key, cd.placeholder())
}

// Gte function generates "%s >= placeholder" for key and adds binding for value
func (cd *commonDialect) Gte(key string, value interface{}) string {
	cd.query.AddBinding(value)
	return fmt.Sprintf("%s >= %s", key, cd.placeholder())
}

// St function generates "%s < placeholder" for key and adds binding for value
func (cd *commonDialect) St(key string, value interface{}) string {
	cd.query.AddBinding(value)
	return fmt.Sprintf("%s < %s", key, cd.placeholder())
}

// Ste function generates "%s <= placeholder" for key and adds binding for value
func (cd *commonDialect) Ste(key string, value interface{}) string {
	cd.query.AddBinding(value)
	return fmt.Sprintf("%s <= %s", key, cd.placeholder())
}

// And function generates " AND " between any number of expressions
func (cd *commonDialect) And(expressions ...string) string {
	return fmt.Sprintf("(%s)", strings.Join(expressions, " AND "))
}

// Or function generates " OR " between any number of expressions
func (cd *commonDialect) Or(expressions ...string) string {
	return strings.Join(expressions, " OR ")
}

// CreateTable generates generic CREATE TABLE statement
func (cd *commonDialect) CreateTable(table string, fields []string, constraints []string) *commonDialect {

	cd.query.AddClause(fmt.Sprintf("CREATE TABLE %s(", table))

	for k, f := range fields {
		clause := fmt.Sprintf("\t%s", f)
		if len(fields)-1 > k || len(constraints) > 0 {
			clause += ","
		}
		cd.query.AddClause(clause)
	}

	for k, c := range constraints {
		constraint := fmt.Sprintf("\t%s", c)
		if len(constraints)-1 > k {
			constraint += ","
		}
		cd.query.AddClause(fmt.Sprintf("%s", constraint))
	}

	cd.query.AddClause(")")
	return cd
}

// AlterTable generates generic ALTER TABLE statement
func (cd *commonDialect) AlterTable(table string) *commonDialect {

	cd.query.AddClause(fmt.Sprintf("ALTER TABLE %s", table))
	return cd
}

// DropTable generates generic DROP TABLE statement
func (cd *commonDialect) DropTable(table string) *commonDialect {

	cd.query.AddClause(fmt.Sprintf("DROP TABLE %s", table))
	return cd
}

// Add generates generic ADD COLUMN statement
func (cd *commonDialect) Add(colName string, colType string) *commonDialect {

	cd.query.AddClause(fmt.Sprintf("ADD %s %s", colName, colType))
	return cd
}

// Drop generates generic DROP COLUMN statement
func (cd *commonDialect) Drop(colName string) *commonDialect {
	cd.query.AddClause(fmt.Sprintf("DROP %s", colName))
	return cd
}

type SqliteDialect struct {
	*commonDialect
}

type MysqlDialect struct {
	*commonDialect
}

type PostgresDialect struct {
	*commonDialect
	bindingIndex int
}

func (pd *PostgresDialect) placeholder() string {
	pd.bindingIndex++
	return fmt.Sprintf("$%d", pd.bindingIndex)
}
