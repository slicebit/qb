---
title: "The Builder"
excerpt: ""
---
The builder api is a way to build complex queries without performance losses from functions using reflection api. The builder uses none of reflection api calls.

This project was started as the only query builder. As the table definitions needed, it is decided to build all the needed toolkit for db based app development.

In the examples, the following initializations are assumed;
[block:code]
{
  "codes": [
    {
      "code": "package main\n\nimport (\n\t\"fmt\"\n\t\"github.com/aacanakin/qb\"\n)\n\nfunc main() {\n\n\tengine, err := qb.NewEngine(\n\t\t\"mysql\",\n\t\t\"root@tcp(localhost:3306)/qb_test?charset=utf8\"\n  )\n\n\tif err != nil {\n\t\tpanic(err)\n\t}\n  \n  // create a builder instance\n  b := qb.NewBuilder(engine.Driver())\n}",
      "language": "go"
    }
  ]
}
[/block]
In qb, `Engine` is one of the core objects that handle db connections & query execution.
[block:api-header]
{
  "type": "basic",
  "title": "Optional escaping"
}
[/block]
qb also provides query escaping for building query which may contain keywords of the database driver. However, due to community feedback, escaping idea was not found as any good. So, there is optional escaping implemented in builder. You can set the escaping method by the following;
[block:code]
{
  "codes": [
    {
      "code": "// create a builder instance\nb := qb.NewBuilder(engine.Driver())\nb.SetEscaping(true)",
      "language": "go"
    }
  ]
}
[/block]

[block:api-header]
{
  "type": "basic",
  "title": "Log flags"
}
[/block]
The builder provides a simple logging mechanism to log queries and bindings. The logging happens when `Query()` function is called. The log flags can be changed by the following;
[block:code]
{
  "codes": [
    {
      "code": "// create a builder instance\nb := qb.NewBuilder(engine.Driver())\nb.SetLogFlags(qb.LQuery|qb.LBindings)",
      "language": "go"
    }
  ]
}
[/block]

[block:parameters]
{
  "data": {
    "h-0": "Log flag",
    "0-0": "LDefault",
    "h-1": "Description",
    "0-1": "The default log flag of qb that means no logging at all.",
    "1-0": "LQuery",
    "1-1": "The `Query()` call would only log the sql statement that is generated",
    "2-0": "LBindings",
    "2-1": "The `Query()` call would only log the bindings array that is generated"
  },
  "cols": 2,
  "rows": 3
}
[/block]

[block:api-header]
{
  "type": "basic",
  "title": "Inserting"
}
[/block]
A simple insert statement generation would be the following;
[block:code]
{
  "codes": [
    {
      "code": "query := b.\n  Insert(\"user\").\n  Values(map[string]interface{}{\n    \"name\":  \"Aras Can\",\n    \"email\": \"aras@slicebit.com\",\n  }).Query()\n\nfmt.Println(query.SQL())\nfmt.Println(query.Bindings())\n\n// outputs\n\n// sql\n// INSERT INTO user\n// (name, email)\n// VALUES (?, ?);\n\n// bindings\n// [Aras Can aras@slicebit.com]",
      "language": "go"
    }
  ]
}
[/block]
The `Query()` call produces a Query struct that has sql string and bindings array. Moreover, `Values()` call understands the column names and escapes them correctly.
[block:api-header]
{
  "type": "basic",
  "title": "Updating"
}
[/block]
A simple update statement generation would be the following;
[block:code]
{
  "codes": [
    {
      "code": "query := b.\n\t\tUpdate(\"user\").\n\t\tSet(map[string]interface{}{\n\t\t\t\"name\": \"Aras Can Akin\",\n\t\t}).\n\t\tWhere(b.Eq(\"name\", \"Aras Can\")).\n\t\tQuery()\n\nfmt.Println(query.SQL())\nfmt.Println(query.Bindings())\n\n// outputs\n// UPDATE user\n// SET name = ?\n// WHERE name = ?;\n// [Aras Can Akin, Aras Can]",
      "language": "go"
    }
  ]
}
[/block]

[block:api-header]
{
  "type": "basic",
  "title": "Deleting"
}
[/block]
A simple delete example is the following;
[block:code]
{
  "codes": [
    {
      "code": "query := b.\n  Delete(\"user\").\n  Where(b.Eq(\"name\", \"Aras Can\")).\n  Query()\n\nfmt.Println(query.SQL())\nfmt.Println(query.Bindings())\n\n// outputs\n// DELETE FROM user\n// WHERE name = ?;\n// [Aras Can]",
      "language": "go"
    }
  ]
}
[/block]

[block:api-header]
{
  "type": "basic",
  "title": "Selecting"
}
[/block]
The builder api provides functionality for building complex select statements. All the function supports can be seen at [here](https://godoc.org/github.com/aacanakin/qb#Builder).

A basic select statement would be the following;
[block:code]
{
  "codes": [
    {
      "code": "query := b.\n  Select(\"id\", \"name\").\n\tFrom(\"user\").\n  Where(b.Eq(\"name\", \"Aras Can\")).\n  Query()\n\nfmt.Println(query.SQL())\nfmt.Println(query.Bindings())\n\n// outputs\n// SELECT id, name\n// FROM user\n// WHERE name = ?;\n// [Aras Can]",
      "language": "go"
    }
  ]
}
[/block]
You can achieve multiple where conditioning using `b.And(queries ...string)` and `b.Or(queries ...string)` as follows;
[block:code]
{
  "codes": [
    {
      "code": "query := b.\n\t\tSelect(\"id\", \"name\").\n\t\tFrom(\"user\").\n\t\tWhere(\n\t\t\tb.And(\n\t\t\t\tb.Eq(\"name\", \"Aras Can\"),\n\t\t\t\tb.NotEq(\"id\", 1),\n\t\t\t),\n\t\t).\n    OrderBy(\"name ASC\").\n\t\tQuery()\n\nfmt.Println(query.SQL())\nfmt.Println(query.Bindings())\n\n// outputs\n// SELECT id, name\n// FROM user\n// WHERE (name = ? AND id != ?);\n// ORDER BY name ASC;\n// [Aras Can 1]",
      "language": "text"
    }
  ]
}
[/block]
The or function is used just as the same with `And()`. The only difference is `Or()` function generates "OR" between conditions.
[block:api-header]
{
  "type": "basic",
  "title": "Comparators"
}
[/block]
The following table shows the helper comparator functions in builder;
[block:code]
{
  "codes": [
    {
      "code": "// NotIn function generates \"%s not in (%s)\" for key and adds bindings for each value\nNotIn(key string, values ...interface{}) string\n\n// In function generates \"%s in (%s)\" for key and adds bindings for each value\nIn(key string, values ...interface{}) string\n\n// NotEq function generates \"%s != placeholder\" for key and adds binding for value\nNotEq(key string, value interface{}) string\n\n// Eq function generates \"%s = placeholder\" for key and adds binding for value\nEq(key string, value interface{}) string\n\n// Gt function generates \"%s > placeholder\" for key and adds binding for value\nGt(key string, value interface{}) string\n\n// Gte function generates \"%s >= placeholder\" for key and adds binding for value\nGte(key string, value interface{}) string\n\n// St function generates \"%s < placeholder\" for key and adds binding for value\nSt(key string, value interface{}) string\n\n// Ste function generates \"%s <= placeholder\" for key and adds binding for value\nSte(key string, value interface{}) string",
      "language": "go"
    }
  ]
}
[/block]
Note that you can still build where clauses using plain text as in the following;
[block:code]
{
  "codes": [
    {
      "code": "query := b.\n\t\tSelect(\"id\", \"name\").\n\t\tFrom(\"user\").\n\t\tWhere(\"name = ?\", \"Aras Can\").\n\t\tQuery()\n\n// outputs\n// SELECT id, name\n// FROM user\n// WHERE name = ?;\n// [Aras Can]",
      "language": "go"
    }
  ]
}
[/block]
As you can see, where also accepts plain conditionals. However, sql statement would not use escape characters. Therefore, you need to manually write escape characters or call `b.Adapter.Escape("name")` which escapes the string.
[block:callout]
{
  "type": "danger",
  "title": "Where",
  "body": "Plain text where doesn't escape column conditionals. Be careful!"
}
[/block]

[block:api-header]
{
  "type": "basic",
  "title": "Aggregates"
}
[/block]
The builder api also provides support for aggregate queries like max, min, count. The following list shows all the aggregate queries supported by builder api;
[block:code]
{
  "codes": [
    {
      "code": "// GroupBy generates \"group by %s\" for each column\nGroupBy(columns ...string) *Builder\n\n// Having generates \"having %s\" for each expression\nHaving(expressions ...string) *Builder\n\n// Avg function generates \"avg(%s)\" statement for column\nAvg(column string) string\n\n// Count function generates \"count(%s)\" statement for column\nCount(column string) string\n\n// Sum function generates \"sum(%s)\" statement for column\nSum(column string) string\n\n// Min function generates \"min(%s)\" statement for column\nMin(column string) string\n\n// Max function generates \"max(%s)\" statement for column\nMax(column string) string",
      "language": "go"
    }
  ]
}
[/block]
Note here that only `GroupBy()` & `Having()` functions can be used by chaining. The other functions are required to be called by using builder instance.
[block:api-header]
{
  "type": "basic",
  "title": "Executing"
}
[/block]
Unless the session api executing queries requires to have engine object initialized. The engine initialization is shown at the top of [The Builder](doc:the-builder) documentation. After building a query, the engine helper functions should execute the `Query` object that builder has built.

The following example shows how to execute `Insert` & `Update` based statements;
[block:code]
{
  "codes": [
    {
      "code": "query := b.\n\t\tInsert(\"user\").\n\t\tValues(map[string]interface{}{\n\t\t\t\"name\": \"Aras\",\n\t\t\t\"id\":   1,\n\t\t}).\n\t\tQuery()\n\nresult, err := engine.Exec(query)\nif err != nil {\n  fmt.Println(err)\n  return\n}\n\nlid, err := result.LastInsertId()\nra, err := result.RowsAffected()\n\nfmt.Println(lid)\nfmt.Println(ra)\n\n// outputs\n// 1\n// 1",
      "language": "go"
    }
  ]
}
[/block]
As it can be seen, the `Exec` function returns `sql.Result` in the `database/sql` package.

Selective statements use `Query()` & `QueryRow()` functions (which doesn't use reflection api) as well as sqlx package functions like `Get()` & `Select()`. The following example would show a simple execution of select statements.
[block:code]
{
  "codes": [
    {
      "code": "var id int\nvar name string\n\nquery := b.\n  Select(\"id\", \"name\").\n  From(\"user\").\n  Where(b.Eq(\"name\", \"Aras\")).\n  Limit(0, 1).\n  Query()\n\nengine.QueryRow(query).Scan(&id, &name)\n\nfmt.Printf(\"<User id=%d name=%s>\\n\", id, name)\n\n// outputs\n// <User id=1 name=Aras>",
      "language": "go"
    }
  ]
}
[/block]
Selecting multiple rows can be achieved by `Query()` as in the following;
[block:code]
{
  "codes": [
    {
      "code": "var id int\nvar name string\n\nquery := b.\n  Select(\"id\", \"name\").\n  From(\"user\").\n  Query()\n\nrows, err := engine.Query(query)\n\nif err != nil {\n  fmt.Println(err)\n  return\n}\n\ndefer rows.Close()\nfor rows.Next() {\n  err := rows.Scan(&id, &name)\n  if err != nil {\n    fmt.Println(err)\n    return\n  }\n\n  fmt.Printf(\"<User id=%d name=%s>\\n\", id, name)\n}\n\n// outputs\n// <User id=1 name=Aras>\n// <User id=2 name=Can>",
      "language": "go"
    }
  ]
}
[/block]
The reflection heavy `Get()` & `Select()` functions are also available for qb.
The `Get()` example is the following;
[block:code]
{
  "codes": [
    {
      "code": "type User struct {\n  ID   int\n  Name string\n}\n\nvar user User\n\nquery := b.\n  Select(\"id\", \"name\").\n  From(\"user\").\n  Where(b.Eq(\"name\", \"Aras\")).\n  Limit(0, 1).\n  Query()\n\nerr = engine.Get(query, &user)\n\nif err != nil {\n  fmt.Println(err)\n  return\n}\n\nfmt.Printf(\"<User id=%d name=%s>\\n\", user.ID, user.Name)\n\n// outputs\n// <User id=1 name=Aras>",
      "language": "go"
    }
  ]
}
[/block]
The `Select()` example is the following;
[block:code]
{
  "codes": [
    {
      "code": "type User struct {\n  ID   int    `db:\"id\"`\n  Name string `db:\"name\"`\n}\n\nusers := []User{}\n\nquery := b.\n  Select(\"name\", \"id\").\n  From(\"user\").\n  Query()\n\nfmt.Println(query.SQL())\nfmt.Println(query.Bindings())\n\nerr = engine.Select(query, &users)\n\nif err != nil {\n  fmt.Println(err)\n  return\n}\n\nfor _, u := range users {\n  fmt.Printf(\"<User id=%d name=%s>\\n\", u.ID, u.Name)\n}\n\n// outputs\n// <User id=1 name=Aras>\n// <User id=2 name=Can>",
      "language": "go"
    }
  ]
}
[/block]