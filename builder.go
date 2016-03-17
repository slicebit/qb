package qb

import (
	"fmt"
	"strings"
)

// NewBuilder generates a new dialect object
func NewBuilder(driver string) *Builder {
	return &Builder{
		query:   NewQuery(),
		dialect: NewDialect(driver),
	}
}

// Builder is a struct that holds an active query that it is used for building common sql queries
// it has all the common functions except multiple statements & table crudders
type Builder struct {
	query   *Query
	dialect Dialect
}

// Dialect returns the active dialect of builder
func (b *Builder) Dialect() Dialect {
	return b.dialect
}

// Reset clears query bindings and its errors
func (b *Builder) Reset() {
	b.query = NewQuery()
	b.dialect.Reset()
}

// Query returns the active query and resets the query.
// The query clauses and returns the sql and bindings
func (b *Builder) Query() *Query {
	query := b.query
	b.Reset()
	fmt.Printf("\n%s\n%s\n", query.SQL(), query.Bindings())
	return query
}

// Insert generates an "insert into %s(%s)" statement
func (b *Builder) Insert(table string) *Builder {
	clause := fmt.Sprintf("INSERT INTO %s", b.dialect.Escape(table))
	b.query.AddClause(clause)
	return b
}

// Values generates "values(%s)" statement and add bindings for each value
func (b *Builder) Values(m map[string]interface{}) *Builder {

	keys := []string{}
	values := []interface{}{}
	for k, v := range m {
		keys = append(keys, k)
		values = append(values, v)
		b.query.AddBinding(v)
	}

	b.query.AddClause(fmt.Sprintf("(%s)", strings.Join(keys, ", ")))

	placeholders := []string{}

	for _ = range values {
		placeholders = append(placeholders, b.dialect.Placeholder())
	}
	clause := fmt.Sprintf("VALUES (%s)", strings.Join(placeholders, ", "))
	b.query.AddClause(clause)
	return b
}

// Update generates "update %s" statement
func (b *Builder) Update(table string) *Builder {
	clause := fmt.Sprintf("UPDATE %s", b.dialect.Escape(table))
	b.query.AddClause(clause)
	return b
}

// Set generates "set a = placeholder" statement for each key a and add bindings for map value
func (b *Builder) Set(m map[string]interface{}) *Builder {
	updates := []string{}
	for k, v := range m {
		updates = append(updates, fmt.Sprintf("%s = %s", k, b.dialect.Placeholder()))
		b.query.AddBinding(v)
	}
	clause := fmt.Sprintf("SET %s", strings.Join(updates, ", "))
	b.query.AddClause(clause)
	return b
}

// Delete generates "delete" statement
func (b *Builder) Delete(table string) *Builder {
	b.query.AddClause(fmt.Sprintf("DELETE FROM %s", b.dialect.Escape(table)))
	return b
}

// Select generates "select %s" statement
func (b *Builder) Select(columns ...string) *Builder {
	clause := fmt.Sprintf("SELECT %s", strings.Join(columns, ", "))
	b.query.AddClause(clause)
	return b
}

// From generates "from %s" statement for each table name
func (b *Builder) From(tables ...string) *Builder {
	tbls := []string{}
	for _, v := range tables {
		tablePieces := strings.Split(v, " ")
		v = b.dialect.Escape(tablePieces[0])
		if len(tablePieces) > 1 {
			v = fmt.Sprintf("%s %s", v, tablePieces[1])
		}
		tbls = append(tbls, v)
	}
	b.query.AddClause(fmt.Sprintf("FROM %s", strings.Join(tbls, ", ")))
	return b
}

// InnerJoin generates "inner join %s on %s" statement for each expression
func (b *Builder) InnerJoin(table string, expressions ...string) *Builder {
	tablePieces := strings.Split(table, " ")

	v := b.dialect.Escape(tablePieces[0])
	if len(tablePieces) > 1 {
		v = fmt.Sprintf("%s %s", v, tablePieces[1])
	}
	b.query.AddClause(fmt.Sprintf("INNER JOIN %s ON %s", v, strings.Join(expressions, " ")))
	return b
}

// CrossJoin generates "cross join %s" statement for table
func (b *Builder) CrossJoin(table string) *Builder {
	tablePieces := strings.Split(table, " ")

	v := b.dialect.Escape(tablePieces[0])
	if len(tablePieces) > 1 {
		v = fmt.Sprintf("%s %s", v, tablePieces[1])
	}

	b.query.AddClause(fmt.Sprintf("CROSS JOIN %s", v))
	return b
}

// LeftOuterJoin generates "left outer join %s on %s" statement for each expression
func (b *Builder) LeftOuterJoin(table string, expressions ...string) *Builder {
	tablePieces := strings.Split(table, " ")

	v := b.dialect.Escape(tablePieces[0])
	if len(tablePieces) > 1 {
		v = fmt.Sprintf("%s %s", v, tablePieces[1])
	}

	b.query.AddClause(fmt.Sprintf("LEFT OUTER JOIN %s ON %s", v, strings.Join(expressions, " ")))
	return b
}

// RightOuterJoin generates "right outer join %s on %s" statement for each expression
func (b *Builder) RightOuterJoin(table string, expressions ...string) *Builder {
	tablePieces := strings.Split(table, " ")

	v := b.dialect.Escape(tablePieces[0])
	if len(tablePieces) > 1 {
		v = fmt.Sprintf("%s %s", v, tablePieces[1])
	}
	b.query.AddClause(fmt.Sprintf("RIGHT OUTER JOIN %s ON %s", v, strings.Join(expressions, " ")))
	return b
}

// FullOuterJoin generates "full outer join %s on %s" for each expression
func (b *Builder) FullOuterJoin(table string, expressions ...string) *Builder {
	tablePieces := strings.Split(table, " ")

	v := b.dialect.Escape(tablePieces[0])
	if len(tablePieces) > 1 {
		v = fmt.Sprintf("%s %s", v, tablePieces[1])
	}
	b.query.AddClause(fmt.Sprintf("FULL OUTER JOIN %s ON %s", v, strings.Join(expressions, " ")))
	return b
}

// Where generates "where %s" for the expression and adds bindings for each value
func (b *Builder) Where(expression string, bindings ...interface{}) *Builder {
	if expression == "" {
		return b
	}
	b.query.AddClause(fmt.Sprintf("WHERE %s", expression))
	b.query.AddBinding(bindings...)
	return b
}

// OrderBy generates "order by %s" for each expression
func (b *Builder) OrderBy(expressions ...string) *Builder {
	b.query.AddClause(fmt.Sprintf("ORDER BY %s", strings.Join(expressions, ", ")))
	return b
}

// GroupBy generates "group by %s" for each column
func (b *Builder) GroupBy(columns ...string) *Builder {
	b.query.AddClause(fmt.Sprintf("GROUP BY %s", strings.Join(columns, ", ")))
	return b
}

// Having generates "having %s" for each expression
func (b *Builder) Having(expressions ...string) *Builder {
	b.query.AddClause(fmt.Sprintf("HAVING %s", strings.Join(expressions, ", ")))
	return b
}

// Limit generates limit %d offset %d for offset and count
func (b *Builder) Limit(offset int, count int) *Builder {
	b.query.AddClause(fmt.Sprintf("LIMIT %d OFFSET %d", count, offset))
	return b
}

// aggregates

// Avg function generates "avg(%s)" statement for column
func (b *Builder) Avg(column string) string {
	return fmt.Sprintf("AVG(%s)", column)
}

// Count function generates "count(%s)" statement for column
func (b *Builder) Count(column string) string {
	return fmt.Sprintf("COUNT(%s)", column)
}

// Sum function generates "sum(%s)" statement for column
func (b *Builder) Sum(column string) string {
	return fmt.Sprintf("SUM(%s)", column)
}

// Min function generates "min(%s)" statement for column
func (b *Builder) Min(column string) string {
	return fmt.Sprintf("MIN(%s)", column)
}

// Max function generates "max(%s)" statement for column
func (b *Builder) Max(column string) string {
	return fmt.Sprintf("MAX(%s)", column)
}

// expressions

// NotIn function generates "%s not in (%s)" for key and adds bindings for each value
func (b *Builder) NotIn(key string, values ...interface{}) string {
	b.query.AddBinding(values...)
	return fmt.Sprintf("%s NOT IN (%s)", key, strings.Join(b.query.QuestionMarks(values...), ","))
}

// In function generates "%s in (%s)" for key and adds bindings for each value
func (b *Builder) In(key string, values ...interface{}) string {
	b.query.AddBinding(values...)
	return fmt.Sprintf("%s IN (%s)", key, strings.Join(b.query.QuestionMarks(values...), ","))
}

// NotEq function generates "%s != placeholder" for key and adds binding for value
func (b *Builder) NotEq(key string, value interface{}) string {
	b.query.AddBinding(value)
	return fmt.Sprintf("%s != ?", key)
}

// Eq function generates "%s = placeholder" for key and adds binding for value
func (b *Builder) Eq(key string, value interface{}) string {
	b.query.AddBinding(value)
	return fmt.Sprintf("%s = ?", key)
}

// Gt function generates "%s > placeholder" for key and adds binding for value
func (b *Builder) Gt(key string, value interface{}) string {
	b.query.AddBinding(value)
	return fmt.Sprintf("%s > ?", key)
}

// Gte function generates "%s >= placeholder" for key and adds binding for value
func (b *Builder) Gte(key string, value interface{}) string {
	b.query.AddBinding(value)
	return fmt.Sprintf("%s >= ?", key)
}

// St function generates "%s < placeholder" for key and adds binding for value
func (b *Builder) St(key string, value interface{}) string {
	b.query.AddBinding(value)
	return fmt.Sprintf("%s < ?", key)
}

// Ste function generates "%s <= placeholder" for key and adds binding for value
func (b *Builder) Ste(key string, value interface{}) string {
	b.query.AddBinding(value)
	return fmt.Sprintf("%s <= ?", key)
}

// And function generates " AND " between any number of expressions
func (b *Builder) And(expressions ...string) string {
	if len(expressions) == 0 {
		return ""
	} else if len(expressions) == 1 {
		return expressions[0]
	}
	return fmt.Sprintf("(%s)", strings.Join(expressions, " AND "))
}

// Or function generates " OR " between any number of expressions
func (b *Builder) Or(expressions ...string) string {
	return strings.Join(expressions, " OR ")
}

// CreateTable generates generic CREATE TABLE statement
func (b *Builder) CreateTable(table string, fields []string, constraints []string) *Builder {

	b.query.AddClause(fmt.Sprintf("CREATE TABLE %s(", b.dialect.Escape(table)))

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

// AlterTable generates generic ALTER TABLE statement
func (b *Builder) AlterTable(table string) *Builder {
	b.query.AddClause(fmt.Sprintf("ALTER TABLE %s", table))
	return b
}

// DropTable generates generic DROP TABLE statement
func (b *Builder) DropTable(table string) *Builder {
	b.query.AddClause(fmt.Sprintf("DROP TABLE %s", b.dialect.Escape(table)))
	return b
}

// Add generates generic ADD COLUMN statement
func (b *Builder) Add(colName string, colType string) *Builder {
	b.query.AddClause(fmt.Sprintf("ADD %s %s", colName, colType))
	return b
}

// Drop generates generic DROP COLUMN statement
func (b *Builder) Drop(colName string) *Builder {
	b.query.AddClause(fmt.Sprintf("DROP %s", colName))
	return b
}
