package qbit

import (
	"errors"
	"fmt"
	"strings"
)

func Builder() *builder {
	return &builder{
		prettyMode: true,
		query:      Query(),
		errors:     []error{},
	}
}

// builder is a subset of dialect could be used for common sql queries
// it has all the common functions except multiple statements & table crudders
type builder struct {
	prettyMode bool
	query      Query
	errors     []error
}

// function builds the query clauses and returns the sql and bindings
func (b *builder) Build() (string, []interface{}, error) {

	if b.HasError() {
		errs := []string{}
		for _, err := range b.errors {
			errs = append(errs, err.Error())
		}
		err := errors.New(strings.Join(errs, "\n"))
		return "", nil, err
	}

	var delimiter string

	if b.prettyMode {
		delimiter = "\n"
	} else {
		delimiter = " "
	}

	sql := fmt.Sprintf("%s;", strings.Join(b.query.Clauses(), delimiter))

	return sql, b.query.Bindings(), nil
}

func (b *builder) AddError(err error) {
	b.errors = append(b.errors, err)
}

func (b *builder) HasError() bool {
	return len(b.errors) > 0
}

func (b *builder) questionMarks(values ...interface{}) []string {
	questionMarks := make([]string, len(values))
	for k, _ := range values {
		questionMarks[k] = "?"
	}
	return questionMarks
}

// function generates an "insert into %s(%s)" statement
func (b *builder) Insert(table string, columns ...string) *builder {
	clause := fmt.Sprintf("INSERT INTO %s(%s)", table, strings.Join(columns, ", "))
	b.query.AddClause(clause)
	return b
}

// function generates "values(%s)" statement and add bindings for each value
func (b *builder) Values(values ...interface{}) *builder {
	b.query.AddBinding(values...)
	clause := fmt.Sprintf("VALUES (%s)", strings.Join(b.questionMarks(values...), ", "))
	b.query.AddClause(clause)
	return b
}

// function generates "update %s" statement
func (b *builder) Update(table string) *builder {
	clause := fmt.Sprintf("UPDATE %s", table)
	b.query.AddClause(clause)
	return b
}

// function generates "set a = ?" statement for each key a and add bindings for map value
func (b *builder) Set(m map[string]interface{}) *builder {
	updates := []string{}
	for k, v := range m {
		updates = append(updates, fmt.Sprintf("%s = ?", k))
		b.query.AddBinding(v)
	}
	clause := fmt.Sprintf("SET %s", strings.Join(updates, ", "))
	b.query.AddClause(clause)
	return b
}

// function generates "select %s" statement
func (b *builder) Select(columns ...string) *builder {
	clause := fmt.Sprintf("SELECT %s", strings.Join(columns, ", "))
	b.query.AddClause(clause)
	return b
}

// function generates "from %s" statement for each table name
func (b *builder) From(tables ...string) *builder {
	b.query.AddClause(fmt.Sprintf("FROM %s", strings.Join(tables, ", ")))
	return b
}

// function generates "inner join %s on %s" statement for each expression
func (b *builder) InnerJoin(table string, expressions ...string) *builder {
	b.query.AddClause(fmt.Sprintf("INNER JOIN %s ON %s", table, strings.Join(expressions, " ")))
	return b
}

// function generates "cross join %s" statement for table
func (b *builder) CrossJoin(table string) *builder {
	b.query.AddClause(fmt.Sprintf("CROSS JOIN %s", table))
	return b
}

// function generates "left outer join %s on %s" statement for each expression
func (b *builder) LeftOuterJoin(table string, expressions ...string) *builder {
	b.query.AddClause(fmt.Sprintf("LEFT OUTER JOIN %s ON %s", table, strings.Join(expressions, " ")))
	return b
}

// function generates "right outer join %s on %s" statement for each expression
func (b *builder) RightOuterJoin(table string, expressions ...string) *builder {
	b.query.AddClause(fmt.Sprintf("RIGHT OUTER JOIN %s ON %s", table, strings.Join(expressions, " ")))
	return b
}

// function generates "full outer join %s on %s" for each expression
func (b *builder) FullOuterJoin(table string, expressions ...string) *builder {
	b.query.AddClause(fmt.Sprintf("FULL OUTER JOIN %s ON %s", table, strings.Join(expressions, " ")))
	return b
}

// function generates "where %s" for the expression and adds bindings for each value
func (b *builder) Where(expression string, bindings ...interface{}) *builder {
	b.query.AddClause(fmt.Sprintf("WHERE %s", expression))
	b.query.AddBinding(bindings...)
	return b
}

// function generates "order by %s" for each expression
func (b *builder) OrderBy(expressions ...string) *builder {
	b.query.AddClause(fmt.Sprintf("ORDER BY %s", strings.Join(expressions, ", ")))
	return b
}

// function generates "group by %s" for each column
func (b *builder) GroupBy(columns ...string) *builder {
	b.query.AddClause(fmt.Sprintf("GROUP BY %s", strings.Join(columns, ", ")))
	return b
}

// function generates "having %s" for each expression
func (b *builder) Having(expressions ...string) *builder {
	b.query.AddClause(fmt.Sprintf("HAVING %s", strings.Join(expressions, ", ")))
	return b
}

// function generates limit %d offset %d for offset and count
func (b *builder) Limit(offset int, count int) *builder {
	b.query.AddClause(fmt.Sprintf("LIMIT %d OFFSET %d", count, offset))
	return b
}

// aggregates

// function generates "avg(%s)" statement for column
func (b *builder) Avg(column string) string {
	return fmt.Sprintf("AVG(%s)", column)
}

// function generates "count(%s)" statement for column
func (b *builder) Count(column string) string {
	return fmt.Sprintf("COUNT(%s)", column)
}

// function generates "sum(%s)" statement for column
func (b *builder) Sum(column string) string {
	return fmt.Sprintf("SUM(%s)", column)
}

// function generates "min(%s)" statement for column
func (b *builder) Min(column string) string {
	return fmt.Sprintf("MIN(%s)", column)
}

// function generates "max(%s)" statement for column
func (b *builder) Max(column string) string {
	return fmt.Sprintf("MAX(%s)", column)
}

// expressions

// function generates "%s not in (%s)" for key and adds bindings for each value
func (b *builder) NotIn(key string, values ...interface{}) string {
	b.query.AddBinding(values...)
	return fmt.Sprintf("%s NOT IN (%s)", key, strings.Join(b.questionMarks(values...), ","))
}

// function generates "%s in (%s)" for key and adds bindings for each value
func (b *builder) In(key string, values ...interface{}) string {
	b.query.AddBinding(values...)
	return fmt.Sprintf("%s IN (%s)", key, strings.Join(b.questionMarks(values...), ","))
}

// function generates "%s != ?" for key and adds binding for value
func (b *builder) NotEq(key string, value interface{}) string {
	b.query.AddBinding(value)
	return fmt.Sprintf("%s != ?", key)
}

// function generates "%s = ?" for key and adds binding for value
func (b *builder) Eq(key string, value interface{}) string {
	b.query.AddBinding(value)
	return fmt.Sprintf("%s = ?", key)
}

// function generates "%s > ?" for key and adds binding for value
func (b *builder) Gt(key string, value interface{}) string {
	b.query.AddBinding(value)
	return fmt.Sprintf("%s > ?", key)
}

// function generates "%s >= ?" for key and adds binding for value
func (b *builder) Gte(key string, value interface{}) string {
	b.query.AddBinding(value)
	return fmt.Sprintf("%s >= ?", key)
}

// function generates "%s < ?" for key and adds binding for value
func (b *builder) St(key string, value interface{}) string {
	b.query.AddBinding(value)
	return fmt.Sprintf("%s < ?", key)
}

// function generates "%s <= ?" for key and adds binding for value
func (b *builder) Ste(key string, value interface{}) string {
	b.query.AddBinding(value)
	return fmt.Sprintf("%s <= ?", key)
}

// function generates " AND " between any number of expressions
func (b *builder) And(expressions ...string) string {
	return fmt.Sprintf("(%s)", strings.Join(expressions, " AND "))
}

// function generates " OR " between any number of expressions
func (b *builder) Or(expressions ...string) string {
	return strings.Join(expressions, " OR ")
}

// function generates generic CREATE TABLE statement
func (b *builder) CreateTable(table string, fields []string, constraints ...string) *builder {

	b.query.AddClause(fmt.Sprintf("CREATE TABLE %s("))

	for _, f := range fields {
		b.query.AddClause(fmt.Sprintf("%s,", f))
	}

	for _, c := range constraints {
		b.query.AddClause(fmt.Sprintf("%s, ", c))
	}

	b.query.AddClause(")")
	return b
}

// function generates generic ALTER TABLE statement
func (b *builder) AlterTable(table string) *builder {

	b.query.AddClause(fmt.Sprintf("ALTER TABLE %s", table))
	return b
}

// function generates generic DROP TABLE statement
func (b *builder) DropTable(table string) *builder {

	b.query.AddClause(fmt.Sprintf("DROP TABLE %s", table))
	return b
}

// function generates generic ADD COLUMN statement
func (b *builder) Add(colName string, colType string) *builder {

	b.query.AddClause(fmt.Sprintf("ADD %s %s", colName, colType))
	return b
}

// function generates generic DROP COLUMN statement
func (b *builder) Drop(colName string) *builder {

	b.query.AddClause(fmt.Sprintf("DROP %s", colName))
	return b
}
