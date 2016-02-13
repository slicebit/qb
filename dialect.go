package qb

// Dialect is the abstraction of set of any expression api functions
type Dialect interface {
	Build() (string, []interface{})

	// clauses
	Insert(table string, columns ...string) *Dialect
	Values(values ...interface{}) *Dialect
	Update(table string) *Dialect
	Set(params map[string]interface{}) *Dialect
	Delete(table string) *Dialect
	Select(columns ...string) *Dialect
	From(tables ...string) *Dialect
	InnerJoin(table string, expressions ...string) *Dialect
	CrossJoin(table string) *Dialect
	LeftOuterJoin(table string, expressions ...string) *Dialect
	RightOuterJoin(table string, expressions ...string) *Dialect
	FullOuterJoin(table string, expressions ...string) *Dialect
	Where(expression string, bindings ...interface{}) *Dialect
	OrderBy(expressions ...string) *Dialect
	GroupBy(columns ...string) *Dialect
	Having(expressions ...string) *Dialect
	Limit(offset int, count int) *Dialect

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

	// multiple statement generators
	Upsert(table string, params map[string]interface{}) *Dialect

	// table crudders
	//	CreateDb(name string) *Dialect
	CreateTable(table string, fields []string, constraints []string) *Dialect
	AlterTable(table string) *Dialect
	Add(colName string, colType string, constraints []string) *Dialect
	Drop(colName string) *Dialect
	DropTable(table string) *Dialect
}
