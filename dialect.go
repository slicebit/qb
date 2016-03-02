package qb

import (
	"fmt"
	"strings"
)

// NewBuilder generates a new builder object
func NewDialect(driver string) *Dialect {
	return &Dialect{
		driver:       driver,
		query:        NewQuery(),
		bindingIndex: 0,
	}
}

// Dialect is a subset of dialect could be used for common sql queries
// it has all the common functions except multiple statements & table crudders
type Dialect struct {
	driver       string
	query        *Query
	bindingIndex int
}

func (d *Dialect) placeholder() string {
	if d.driver == "postgres" {
		d.bindingIndex++
		return fmt.Sprintf("$%d", d.bindingIndex)
	}
	return "?"
}

func (d *Dialect) placeholders(values ...interface{}) []string {
	placeholders := make([]string, len(values))
	for k := range values {
		placeholders[k] = d.placeholder()
	}
	return placeholders
}

// Reset clears query bindings and its errors
func (d *Dialect) Reset() {
	d.query = NewQuery()
	d.bindingIndex = 0
}

// Build generates sql query and bindings from clauses and bindings.
// The query clauses and returns the sql and bindings
func (d *Dialect) Query() *Query {
	query := d.query
	d.Reset()
	return query
}

// Insert generates an "insert into %s(%s)" statement
func (d *Dialect) Insert(table string, columns ...string) *Dialect {
	clause := fmt.Sprintf("INSERT INTO %s(%s)", table, strings.Join(columns, ", "))
	d.query.AddClause(clause)
	return d
}

// Values generates "values(%s)" statement and add bindings for each value
func (d *Dialect) Values(values ...interface{}) *Dialect {
	d.query.AddBinding(values...)
	clause := fmt.Sprintf("VALUES (%s)", strings.Join(d.placeholders(values...), ", "))
	d.query.AddClause(clause)
	return d
}

// Update generates "update %s" statement
func (d *Dialect) Update(table string) *Dialect {
	clause := fmt.Sprintf("UPDATE %s", table)
	d.query.AddClause(clause)
	return d
}

// Set generates "set a = placeholder" statement for each key a and add bindings for map value
func (d *Dialect) Set(m map[string]interface{}) *Dialect {
	updates := []string{}
	for k, v := range m {
		updates = append(updates, fmt.Sprintf("%s = %s", k, d.placeholder()))
		d.query.AddBinding(v)
	}
	clause := fmt.Sprintf("SET %s", strings.Join(updates, ", "))
	d.query.AddClause(clause)
	return d
}

// Delete generates "delete" statement
func (d *Dialect) Delete(table string) *Dialect {
	d.query.AddClause(fmt.Sprintf("DELETE FROM %s", table))
	return d
}

// Select generates "select %s" statement
func (d *Dialect) Select(columns ...string) *Dialect {
	clause := fmt.Sprintf("SELECT %s", strings.Join(columns, ", "))
	d.query.AddClause(clause)
	return d
}

// From generates "from %s" statement for each table name
func (d *Dialect) From(tables ...string) *Dialect {
	d.query.AddClause(fmt.Sprintf("FROM %s", strings.Join(tables, ", ")))
	return d
}

// InnerJoin generates "inner join %s on %s" statement for each expression
func (d *Dialect) InnerJoin(table string, expressions ...string) *Dialect {
	d.query.AddClause(fmt.Sprintf("INNER JOIN %s ON %s", table, strings.Join(expressions, " ")))
	return d
}

// CrossJoin generates "cross join %s" statement for table
func (d *Dialect) CrossJoin(table string) *Dialect {
	d.query.AddClause(fmt.Sprintf("CROSS JOIN %s", table))
	return d
}

// LeftOuterJoin generates "left outer join %s on %s" statement for each expression
func (d *Dialect) LeftOuterJoin(table string, expressions ...string) *Dialect {
	d.query.AddClause(fmt.Sprintf("LEFT OUTER JOIN %s ON %s", table, strings.Join(expressions, " ")))
	return d
}

// RightOuterJoin generates "right outer join %s on %s" statement for each expression
func (d *Dialect) RightOuterJoin(table string, expressions ...string) *Dialect {
	d.query.AddClause(fmt.Sprintf("RIGHT OUTER JOIN %s ON %s", table, strings.Join(expressions, " ")))
	return d
}

// FullOuterJoin generates "full outer join %s on %s" for each expression
func (d *Dialect) FullOuterJoin(table string, expressions ...string) *Dialect {
	d.query.AddClause(fmt.Sprintf("FULL OUTER JOIN %s ON %s", table, strings.Join(expressions, " ")))
	return d
}

// Where generates "where %s" for the expression and adds bindings for each value
func (d *Dialect) Where(expression string, bindings ...interface{}) *Dialect {
	expression = strings.Replace(expression, "?", d.placeholder(), -1)
	d.query.AddClause(fmt.Sprintf("WHERE %s", expression))
	d.query.AddBinding(bindings...)
	return d
}

// OrderBy generates "order by %s" for each expression
func (d *Dialect) OrderBy(expressions ...string) *Dialect {
	d.query.AddClause(fmt.Sprintf("ORDER BY %s", strings.Join(expressions, ", ")))
	return d
}

// GroupBy generates "group by %s" for each column
func (d *Dialect) GroupBy(columns ...string) *Dialect {
	d.query.AddClause(fmt.Sprintf("GROUP BY %s", strings.Join(columns, ", ")))
	return d
}

// Having generates "having %s" for each expression
func (d *Dialect) Having(expressions ...string) *Dialect {
	d.query.AddClause(fmt.Sprintf("HAVING %s", strings.Join(expressions, ", ")))
	return d
}

// Limit generates limit %d offset %d for offset and count
func (d *Dialect) Limit(offset int, count int) *Dialect {
	d.query.AddClause(fmt.Sprintf("LIMIT %d OFFSET %d", count, offset))
	return d
}

// aggregates

// Avg function generates "avg(%s)" statement for column
func (d *Dialect) Avg(column string) string {
	return fmt.Sprintf("AVG(%s)", column)
}

// Count function generates "count(%s)" statement for column
func (d *Dialect) Count(column string) string {
	return fmt.Sprintf("COUNT(%s)", column)
}

// Sum function generates "sum(%s)" statement for column
func (d *Dialect) Sum(column string) string {
	return fmt.Sprintf("SUM(%s)", column)
}

// Min function generates "min(%s)" statement for column
func (d *Dialect) Min(column string) string {
	return fmt.Sprintf("MIN(%s)", column)
}

// Max function generates "max(%s)" statement for column
func (d *Dialect) Max(column string) string {
	return fmt.Sprintf("MAX(%s)", column)
}

// expressions

// NotIn function generates "%s not in (%s)" for key and adds bindings for each value
func (d *Dialect) NotIn(key string, values ...interface{}) string {
	d.query.AddBinding(values...)
	return fmt.Sprintf("%s NOT IN (%s)", key, strings.Join(d.placeholders(values...), ","))
}

// In function generates "%s in (%s)" for key and adds bindings for each value
func (d *Dialect) In(key string, values ...interface{}) string {
	d.query.AddBinding(values...)
	return fmt.Sprintf("%s IN (%s)", key, strings.Join(d.placeholders(values...), ","))
}

// NotEq function generates "%s != placeholder" for key and adds binding for value
func (d *Dialect) NotEq(key string, value interface{}) string {
	d.query.AddBinding(value)
	return fmt.Sprintf("%s != %s", key, d.placeholder())
}

// Eq function generates "%s = placeholder" for key and adds binding for value
func (d *Dialect) Eq(key string, value interface{}) string {
	d.query.AddBinding(value)
	return fmt.Sprintf("%s = %s", key, d.placeholder())
}

// Gt function generates "%s > placeholder" for key and adds binding for value
func (d *Dialect) Gt(key string, value interface{}) string {
	d.query.AddBinding(value)
	return fmt.Sprintf("%s > %s", key, d.placeholder())
}

// Gte function generates "%s >= placeholder" for key and adds binding for value
func (d *Dialect) Gte(key string, value interface{}) string {
	d.query.AddBinding(value)
	return fmt.Sprintf("%s >= %s", key, d.placeholder())
}

// St function generates "%s < placeholder" for key and adds binding for value
func (d *Dialect) St(key string, value interface{}) string {
	d.query.AddBinding(value)
	return fmt.Sprintf("%s < %s", key, d.placeholder())
}

// Ste function generates "%s <= placeholder" for key and adds binding for value
func (d *Dialect) Ste(key string, value interface{}) string {
	d.query.AddBinding(value)
	return fmt.Sprintf("%s <= %s", key, d.placeholder())
}

// And function generates " AND " between any number of expressions
func (d *Dialect) And(expressions ...string) string {
	return fmt.Sprintf("(%s)", strings.Join(expressions, " AND "))
}

// Or function generates " OR " between any number of expressions
func (d *Dialect) Or(expressions ...string) string {
	return strings.Join(expressions, " OR ")
}

// CreateTable generates generic CREATE TABLE statement
func (d *Dialect) CreateTable(table string, fields []string, constraints []string) *Dialect {

	d.query.AddClause(fmt.Sprintf("CREATE TABLE %s(", table))

	for k, f := range fields {
		clause := fmt.Sprintf("\t%s", f)
		if len(fields)-1 > k || len(constraints) > 0 {
			clause += ","
		}
		d.query.AddClause(clause)
	}

	for k, c := range constraints {
		constraint := fmt.Sprintf("\t%s", c)
		if len(constraints)-1 > k {
			constraint += ","
		}
		d.query.AddClause(fmt.Sprintf("%s", constraint))
	}

	d.query.AddClause(")")
	return d
}

// AlterTable generates generic ALTER TABLE statement
func (d *Dialect) AlterTable(table string) *Dialect {
	d.query.AddClause(fmt.Sprintf("ALTER TABLE %s", table))
	return d
}

// DropTable generates generic DROP TABLE statement
func (d *Dialect) DropTable(table string) *Dialect {
	d.query.AddClause(fmt.Sprintf("DROP TABLE %s", table))
	return d
}

// Add generates generic ADD COLUMN statement
func (d *Dialect) Add(colName string, colType string) *Dialect {
	d.query.AddClause(fmt.Sprintf("ADD %s %s", colName, colType))
	return d
}

// Drop generates generic DROP COLUMN statement
func (d *Dialect) Drop(colName string) *Dialect {
	d.query.AddClause(fmt.Sprintf("DROP %s", colName))
	return d
}
