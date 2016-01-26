package qbit

import (
	"errors"
	"fmt"
	"strings"
)

func NewBuilder() *Builder {
	return &Builder{
		prettyMode: true,
		query:      NewQuery(),
		errors:     []error{},
	}
}

const BUILDER_DELIMITER  = "\n"

// builder is a subset of dialect could be used for common sql queries
// it has all the common functions except multiple statements & table crudders
type Builder struct {
	prettyMode bool
	query      *Query
	errors     []error
}

// function resets query and its errors
func (b *Builder) Reset() {
	b.query = NewQuery()
	b.errors = []error{}
}

// function builds the query clauses and returns the sql and bindings
func (b *Builder) Build() (string, []interface{}, error) {

	if b.HasError() {
		errs := []string{}
		for _, err := range b.errors {
			errs = append(errs, err.Error())
		}
		err := errors.New(strings.Join(errs, "\n"))
		b.errors = []error{}
		return "", nil, err
	}

//	var delimiter string
//
//	if b.prettyMode {
//		delimiter = "\n"
//	} else {
//		delimiter = " "
//	}

	bindings := b.query.Bindings()
	sql := fmt.Sprintf("%s;", strings.Join(b.query.Clauses(), BUILDER_DELIMITER))

	b.query = NewQuery()

	return sql, bindings, nil
}

func (b *Builder) AddError(err error) {
	b.errors = append(b.errors, err)
}

func (b *Builder) HasError() bool {
	return len(b.errors) > 0
}

func (b *Builder) Errors() []error {
	return b.errors
}

func (b *Builder) questionMarks(values ...interface{}) []string {
	questionMarks := make([]string, len(values))
	for k, _ := range values {
		questionMarks[k] = "?"
	}
	return questionMarks
}

// function generates an "insert into %s(%s)" statement
func (b *Builder) Insert(table string, columns ...string) *Builder {
	clause := fmt.Sprintf("INSERT INTO %s(%s)", table, strings.Join(columns, ", "))
	b.query.AddClause(clause)
	return b
}

// function generates "values(%s)" statement and add bindings for each value
func (b *Builder) Values(values ...interface{}) *Builder {
	b.query.AddBinding(values...)
	clause := fmt.Sprintf("VALUES (%s)", strings.Join(b.questionMarks(values...), ", "))
	b.query.AddClause(clause)
	return b
}

// function generates "update %s" statement
func (b *Builder) Update(table string) *Builder {
	clause := fmt.Sprintf("UPDATE %s", table)
	b.query.AddClause(clause)
	return b
}

// function generates "set a = ?" statement for each key a and add bindings for map value
func (b *Builder) Set(m map[string]interface{}) *Builder {
	updates := []string{}
	for k, v := range m {
		updates = append(updates, fmt.Sprintf("%s = ?", k))
		b.query.AddBinding(v)
	}
	clause := fmt.Sprintf("SET %s", strings.Join(updates, ", "))
	b.query.AddClause(clause)
	return b
}

// function generates "delete" statement
func (b *Builder) Delete(table string) *Builder {
	b.query.AddClause(fmt.Sprintf("DELETE FROM %s", table))
	return b
}

// function generates "select %s" statement
func (b *Builder) Select(columns ...string) *Builder {
	clause := fmt.Sprintf("SELECT %s", strings.Join(columns, ", "))
	b.query.AddClause(clause)
	return b
}

// function generates "from %s" statement for each table name
func (b *Builder) From(tables ...string) *Builder {
	b.query.AddClause(fmt.Sprintf("FROM %s", strings.Join(tables, ", ")))
	return b
}

// function generates "inner join %s on %s" statement for each expression
func (b *Builder) InnerJoin(table string, expressions ...string) *Builder {
	b.query.AddClause(fmt.Sprintf("INNER JOIN %s ON %s", table, strings.Join(expressions, " ")))
	return b
}

// function generates "cross join %s" statement for table
func (b *Builder) CrossJoin(table string) *Builder {
	b.query.AddClause(fmt.Sprintf("CROSS JOIN %s", table))
	return b
}

// function generates "left outer join %s on %s" statement for each expression
func (b *Builder) LeftOuterJoin(table string, expressions ...string) *Builder {
	b.query.AddClause(fmt.Sprintf("LEFT OUTER JOIN %s ON %s", table, strings.Join(expressions, " ")))
	return b
}

// function generates "right outer join %s on %s" statement for each expression
func (b *Builder) RightOuterJoin(table string, expressions ...string) *Builder {
	b.query.AddClause(fmt.Sprintf("RIGHT OUTER JOIN %s ON %s", table, strings.Join(expressions, " ")))
	return b
}

// function generates "full outer join %s on %s" for each expression
func (b *Builder) FullOuterJoin(table string, expressions ...string) *Builder {
	b.query.AddClause(fmt.Sprintf("FULL OUTER JOIN %s ON %s", table, strings.Join(expressions, " ")))
	return b
}

// function generates "where %s" for the expression and adds bindings for each value
func (b *Builder) Where(expression string, bindings ...interface{}) *Builder {
	b.query.AddClause(fmt.Sprintf("WHERE %s", expression))
	b.query.AddBinding(bindings...)
	return b
}

// function generates "order by %s" for each expression
func (b *Builder) OrderBy(expressions ...string) *Builder {
	b.query.AddClause(fmt.Sprintf("ORDER BY %s", strings.Join(expressions, ", ")))
	return b
}

// function generates "group by %s" for each column
func (b *Builder) GroupBy(columns ...string) *Builder {
	b.query.AddClause(fmt.Sprintf("GROUP BY %s", strings.Join(columns, ", ")))
	return b
}

// function generates "having %s" for each expression
func (b *Builder) Having(expressions ...string) *Builder {
	b.query.AddClause(fmt.Sprintf("HAVING %s", strings.Join(expressions, ", ")))
	return b
}

// function generates limit %d offset %d for offset and count
func (b *Builder) Limit(offset int, count int) *Builder {
	b.query.AddClause(fmt.Sprintf("LIMIT %d OFFSET %d", count, offset))
	return b
}

// aggregates

// function generates "avg(%s)" statement for column
func (b *Builder) Avg(column string) string {
	return fmt.Sprintf("AVG(%s)", column)
}

// function generates "count(%s)" statement for column
func (b *Builder) Count(column string) string {
	return fmt.Sprintf("COUNT(%s)", column)
}

// function generates "sum(%s)" statement for column
func (b *Builder) Sum(column string) string {
	return fmt.Sprintf("SUM(%s)", column)
}

// function generates "min(%s)" statement for column
func (b *Builder) Min(column string) string {
	return fmt.Sprintf("MIN(%s)", column)
}

// function generates "max(%s)" statement for column
func (b *Builder) Max(column string) string {
	return fmt.Sprintf("MAX(%s)", column)
}

// expressions

// function generates "%s not in (%s)" for key and adds bindings for each value
func (b *Builder) NotIn(key string, values ...interface{}) string {
	b.query.AddBinding(values...)
	return fmt.Sprintf("%s NOT IN (%s)", key, strings.Join(b.questionMarks(values...), ","))
}

// function generates "%s in (%s)" for key and adds bindings for each value
func (b *Builder) In(key string, values ...interface{}) string {
	b.query.AddBinding(values...)
	return fmt.Sprintf("%s IN (%s)", key, strings.Join(b.questionMarks(values...), ","))
}

// function generates "%s != ?" for key and adds binding for value
func (b *Builder) NotEq(key string, value interface{}) string {
	b.query.AddBinding(value)
	return fmt.Sprintf("%s != ?", key)
}

// function generates "%s = ?" for key and adds binding for value
func (b *Builder) Eq(key string, value interface{}) string {
	b.query.AddBinding(value)
	return fmt.Sprintf("%s = ?", key)
}

// function generates "%s > ?" for key and adds binding for value
func (b *Builder) Gt(key string, value interface{}) string {
	b.query.AddBinding(value)
	return fmt.Sprintf("%s > ?", key)
}

// function generates "%s >= ?" for key and adds binding for value
func (b *Builder) Gte(key string, value interface{}) string {
	b.query.AddBinding(value)
	return fmt.Sprintf("%s >= ?", key)
}

// function generates "%s < ?" for key and adds binding for value
func (b *Builder) St(key string, value interface{}) string {
	b.query.AddBinding(value)
	return fmt.Sprintf("%s < ?", key)
}

// function generates "%s <= ?" for key and adds binding for value
func (b *Builder) Ste(key string, value interface{}) string {
	b.query.AddBinding(value)
	return fmt.Sprintf("%s <= ?", key)
}

// function generates " AND " between any number of expressions
func (b *Builder) And(expressions ...string) string {
	return fmt.Sprintf("(%s)", strings.Join(expressions, " AND "))
}

// function generates " OR " between any number of expressions
func (b *Builder) Or(expressions ...string) string {
	return strings.Join(expressions, " OR ")
}

// function generates generic CREATE TABLE statement
func (b *Builder) CreateTable(table string, fields []string, constraints []string) *Builder {

	b.query.AddClause(fmt.Sprintf("CREATE TABLE %s(", table))

	for k, f := range fields {
		clause := fmt.Sprintf("\t%s", f)
		if len(fields)-1 > k || len(constraints) > 0 {
			clause += ","
		}
		b.query.AddClause(clause)
	}

	for k, c := range constraints {
		constraint := fmt.Sprintf("\t%s", c)
		if len(constraints)-1 > k {
			constraint += ","
		}
		b.query.AddClause(fmt.Sprintf("%s", constraint))
	}

	b.query.AddClause(")")
	return b
}

// function generates generic ALTER TABLE statement
func (b *Builder) AlterTable(table string) *Builder {

	b.query.AddClause(fmt.Sprintf("ALTER TABLE %s", table))
	return b
}

// function generates generic DROP TABLE statement
func (b *Builder) DropTable(table string) *Builder {

	b.query.AddClause(fmt.Sprintf("DROP TABLE %s", table))
	return b
}

// function generates generic ADD COLUMN statement
func (b *Builder) Add(colName string, colType string) *Builder {

	b.query.AddClause(fmt.Sprintf("ADD %s %s", colName, colType))
	return b
}

// function generates generic DROP COLUMN statement
func (b *Builder) Drop(colName string) *Builder {

	b.query.AddClause(fmt.Sprintf("DROP %s", colName))
	return b
}
