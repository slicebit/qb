---
title: "Defining Tables"
excerpt: ""
---
In the [Defining Models](doc:defining-models) models] section, it is explained how to build tables from structs. However, there is a secondary option to create tables without using any structs. This method is a lightweight solution to build tables that can be used in manipulating data. This method also doesn't use any reflection.

The reason to have this in the expression api is simply because you may need extra flexibility to build table objects that can work with databases. Moreover, defining tables using this method would be more idiomatic to go. Here's an example user struct using orm api;
[block:code]
{
  "codes": [
    {
      "code": "type User struct {\n  ID string `db:\"_id\" qb:\"type:varchar(36); constraints:primary_key\"`\n  Email string `db:\"_email\" qb:\"type:varchar(64); constraints: unique, not_null, \"`\n}",
      "language": "go"
    }
  ]
}
[/block]
As you might notice, the tags are long enough to keep track of it. Here's the equivalent table using the expression api;
[block:code]
{
  "codes": [
    {
      "code": "usersTable := qb.Table(\n  \"user\",\n  qb.Column(\"_id\", qb.Varchar().Size(40)),\n  qb.Column(\"_email\", qb.Varchar().Size(64).Unique().NotNull()),\n  qb.PrimaryKey(\"_id\"),\n)\n\nfmt.Println(usersTable.Create(db.Builder().Adapter()))\n\n// prints\n/*\nCREATE TABLE user (\n\t_id VARCHAR(40),\n\t_email VARCHAR(64) UNIQUE NOT NULL,\n\tPRIMARY KEY(_id)\n);\n*/",
      "language": "go"
    }
  ]
}
[/block]
This method is more performant than the auto mapping method. Performance critical applications should use table api instead of orm.
[block:api-header]
{
  "type": "basic",
  "title": "Defining Column Types & Constraints"
}
[/block]
The `Column(name string, t TypeElem)` function generates a column and appends a column to the table. The following builtin types are available for columns;
[block:parameters]
{
  "data": {
    "0-0": "Char()",
    "h-0": "func",
    "0-1": "generates a char type",
    "1-0": "Varchar()",
    "1-1": "generates a varchar type with 255 as default size",
    "2-0": "Text()",
    "2-1": "generates a text type",
    "3-0": "SmallInt()",
    "3-1": "generates a small int type",
    "4-0": "Int()",
    "4-1": "generates an int type",
    "6-0": "Numeric()",
    "6-1": "generates a numeric type",
    "7-0": "Float()",
    "7-1": "generates a float type",
    "5-0": "BigInt()",
    "5-1": "generates an int64 type",
    "8-0": "Boolean()",
    "8-1": "generates a boolean type",
    "9-0": "Timestamp()",
    "9-1": "generates a timestamp type",
    "h-1": "desc"
  },
  "cols": 2,
  "rows": 10
}
[/block]
You can optionally define your own type using `Type(name string)` function as in the following;
[block:code]
{
  "codes": [
    {
      "code": "col := qb.Column(\"id\", qb.Type(\"UUID\"))",
      "language": "go"
    }
  ]
}
[/block]
After defining types, the constraints of that column can be defined using chaining as in the following exampe;
[block:code]
{
  "codes": [
    {
      "code": "col := qb.Column(\"email\", qb.Varchar().Size(40).NotNull().Unique())",
      "language": "go"
    }
  ]
}
[/block]
This means that we have an `email` column with `VARCHAR(40) NOT NULL UNIQUE` types & constraints.

Here is the builtin constraints you can use;
[block:parameters]
{
  "data": {
    "h-0": "func",
    "h-1": "Desc",
    "0-0": "Size(size int)",
    "0-1": "Adds a size constraint given size",
    "1-0": "Default(val interface{})",
    "1-1": "Adds a default constraint given value",
    "2-0": "Null()",
    "2-1": "Adds a nullable constraint",
    "3-0": "NotNull()",
    "3-1": "Adds a non nullable constraint",
    "4-0": "Unique()",
    "4-1": "Adds a unique constraint"
  },
  "cols": 2,
  "rows": 5
}
[/block]
You can optionally define your custom constraint using `Constraint(name string)` function as in the following;
[block:code]
{
  "codes": [
    {
      "code": "col := qb.Column(\"id\", qb.Int().Constraint(\"CHECK (id > 0)\"))",
      "language": "go"
    }
  ]
}
[/block]
This example adds email column a check constraint
[block:api-header]
{
  "type": "basic",
  "title": "Defining Table constraints"
}
[/block]
The table constraints (primary, foreign keys, unique keys) can be defined using builtin table constraints.
[block:code]
{
  "codes": [
    {
      "code": "usersTable := qb.Table(\n  \"user\",\n  qb.Column(\"id\", qb.Varchar().Size(36)),\n  qb.Column(\"address_id\", qb.Text()),\n  qb.PrimaryKey(\"id\"),\n  qb.ForeignKey().Ref(\"address_id\", \"address\", \"id\"),\n)",
      "language": "go"
    }
  ]
}
[/block]
This will tell qb to add a primary key table constraint on `id` column and a foreign key constraint on `address_id` referenced from `address.id`.

You can also add composite unique keys using `UniqueKey(cols ...string)` function. Note that it should be used `Unique()` in defining column type if a single unique constraint is required.
[block:code]
{
  "codes": [
    {
      "code": "usersTable := qb.Table(\n  \"user\",\n  qb.Column(\"id\", qb.Varchar().Size(36)),\n  qb.Column(\"address_id\", qb.Text()),\n  qb.UniqueKey(\"id\", \"address_id\"),\n)",
      "language": "go"
    }
  ]
}
[/block]

[block:api-header]
{
  "type": "basic",
  "title": "Defining table indices"
}
[/block]
Table indices can be defined by chaining `Table()` function of qb.
[block:code]
{
  "codes": [
    {
      "code": "usersTable := qb.Table(\n  \"user\",\n  qb.Column(\"id\", qb.Varchar().Size(36)),\n  qb.Column(\"address_id\", qb.Text()),\n).Index(\"id\").Index(\"id\", \"address_id\")",
      "language": "go"
    }
  ]
}
[/block]
As you might notice, `Index()` function can be used single indices as well as composite indices.