---
title: "Defining Models"
excerpt: ""
---
To define models, qb uses structs with tagging options. This tutorial will continue without defining qb again and again. Therefore it is assumed here that qb is initialized by the following code;
[block:callout]
{
  "type": "info",
  "title": "Examples",
  "body": "Although the library works in mysql and sqlite, the following examples would be in postgres"
}
[/block]

[block:code]
{
  "codes": [
    {
      "code": "import (\n\t\"github.com/aacanakin/qb\" \n)\n\ndb, err := qb.New(\"postgres\", \"user=postgres dbname=qb_test sslmode=disable\")\nif err != nil {\n  panic(err)\n}\ndefer db.Close()",
      "language": "go"
    }
  ]
}
[/block]
The simplest model definition would be the following;
[block:code]
{
  "codes": [
    {
      "code": "type User struct {\n  ID int64\n  Name string\n}\n\n// register struct to metadata to be mapped\ndb.Metadata().Add(User{})\n\n// create all tables registered to metadata within a transaction\ndb.Metadata().CreateAll()",
      "language": "go"
    }
  ]
}
[/block]
The code will generate a transaction and creates the following sql statement;
[block:code]
{
  "codes": [
    {
      "code": "CREATE TABLE user(\n\tid BIGINT,\n\tname VARCHAR(255)\n);",
      "language": "sql"
    }
  ]
}
[/block]
It is noticeable here that the most simple table would not solve any of our real world problems. Therefore, let's add some constraints.
[block:api-header]
{
  "type": "basic",
  "title": "Ignoring struct fields"
}
[/block]
It may be required to have ignoring for a pre or post processed field in a struct that won't be in a database field. Use `-` character to ignore a struct field as in the following; 
[block:code]
{
  "codes": [
    {
      "code": "type User struct {\n\t\tID        int64      `qb:\"constraints:primary_key\"`\n\t\tSessions  []Session  `qb:\"-\"`\n}",
      "language": "go"
    }
  ]
}
[/block]
In this struct definition, the `Sessions` field would not be added in the database.
[block:api-header]
{
  "type": "basic",
  "title": "Overriding colum names"
}
[/block]
You can also override column names using `db` tag. The following example shows how to override column names in the db tables;
[block:code]
{
  "codes": [
    {
      "code": "type User struct {\n  ID int64 `db:\"_id\" qb:\"constraints:primary_key\"` \n}",
      "language": "go"
    }
  ]
}
[/block]
This will tell qb to map `ID` field into `_id` in the db table.
[block:api-header]
{
  "type": "basic",
  "title": "Constraints"
}
[/block]
In qb there are several constraint definitions. The most common is to use struct tags. Here's the following;
[block:code]
{
  "codes": [
    {
      "code": "type User struct {\n\tID   int64 `qb:\"constraints:primary_key\"`\n  Name string\n}\n\n// register struct to metadata to be mapped\ndb.Metadata().Add(User{})\n\n// create all tables registered to metadata within a transaction\ndb.Metadata().CreateAll()",
      "language": "go"
    }
  ]
}
[/block]
This would generate the following sql statement in a transaction;
[block:code]
{
  "codes": [
    {
      "code": "CREATE TABLE user(\n\tid BIGINT PRIMARY KEY,\n\tname VARCHAR(255)\n);",
      "language": "sql"
    }
  ]
}
[/block]
Let's improve our model and use not null & unique constraints
[block:code]
{
  "codes": [
    {
      "code": "type User struct {\n\t\tID        int64      `qb:\"constraints:primary_key\"`\n\t\tName      string     `qb:\"constraints:not_null\"`\n\t\tEmail     string     `qb:\"constraints:unique, not_null\"`\n\t\tCreatedAt time.Time  `qb:\"constraints:not_null\"`\n\t\tDeletedAt *time.Time `qb:\"constraints:null\"`\n}\n\ndb.Metadata().Add(User{})",
      "language": "go"
    }
  ]
}
[/block]
Notice here that `DeletedAt` field is a time.Time pointer instead of time.Time. It is because to make `DeletedAt` field nullable. You can also have pointer types with `not null` constraints.

This would generate the following sql statement in a transaction;
[block:code]
{
  "codes": [
    {
      "code": "CREATE TABLE user(\n\tid BIGINT PRIMARY KEY,\n\tname VARCHAR(255) NOT NULL,\n\temail VARCHAR(255) UNIQUE NOT NULL,\n\tcreated_at TIMESTAMP NOT NULL,\n\tdeleted_at TIMESTAMP NULL\n);",
      "language": "sql"
    }
  ]
}
[/block]

[block:api-header]
{
  "type": "basic",
  "title": "Foreign Keys"
}
[/block]
Foreign keys lets you define relationships between tables. In qb, foreign keys are done in the following example;
[block:code]
{
  "codes": [
    {
      "code": "type User struct {\n\t\tID        int64      `qb:\"constraints:primary_key\"`\n\t\tName      string     `qb:\"constraints:not_null\"`\n\t\tEmail     string     `qb:\"constraints:unique, not_null\"`\n\t\tCreatedAt time.Time  `qb:\"constraints:not_null\"`\n\t\tDeletedAt *time.Time `qb:\"constraints:null\"`\n\t}\n\ntype Session struct {\n\t\tID        string `qb:\"type:uuid; constraints:primary_key\"`\n\t\tUserID    int64  `qb:\"constraints:ref(user.id)\"`\n\t\tAuthToken string `qb:\"type:uuid\", constraints:not_null`\n}\n\ndb.Metadata().Add(User{})\ndb.Metadata().Add(Session{})\ndb.Metadata().CreateAll()",
      "language": "go"
    }
  ]
}
[/block]
As you might notice, qb supports definitions of multiple constraints like in the `User.Email` field.

This would generate the following sql statements within a transaction;
[block:code]
{
  "codes": [
    {
      "code": "CREATE TABLE user(\n\tid BIGINT PRIMARY KEY,\n\tname VARCHAR(255) NOT NULL,\n\temail VARCHAR(255) UNIQUE NOT NULL,\n\tcreated_at TIMESTAMP NOT NULL,\n\tdeleted_at TIMESTAMP NULL\n);\n\nCREATE TABLE session(\n\tid UUID PRIMARY KEY,\n\tuser_id BIGINT,\n\tauth_token UUID,\n\tFOREIGN KEY (user_id) REFERENCES user(id)\n);",
      "language": "sql"
    }
  ]
}
[/block]
As it might be noticed, it is really really easy to build simple relationships between tables.
[block:api-header]
{
  "type": "basic",
  "title": "Types"
}
[/block]
The following type mappings are applied in qb types;
[block:parameters]
{
  "data": {
    "h-0": "go type",
    "h-1": "mysql",
    "h-2": "postgres",
    "h-3": "sqlite",
    "0-0": "string",
    "0-1": "VARCHAR(255)",
    "0-2": "VARCHAR(255)",
    "0-3": "VARCHAR(255)",
    "1-0": "int",
    "1-1": "INT",
    "1-2": "INT",
    "1-3": "INT",
    "2-0": "int8",
    "2-1": "SMALLINT",
    "2-2": "SMALLINT",
    "2-3": "SMALLINT",
    "3-0": "int16",
    "3-1": "SMALLINT",
    "3-2": "SMALLINT",
    "3-3": "SMALLINT",
    "4-0": "int32",
    "4-1": "INT",
    "4-2": "INT",
    "4-3": "INT",
    "5-0": "int64",
    "5-1": "BIGINT",
    "5-2": "BIGINT",
    "5-3": "BIGINT",
    "6-0": "uint",
    "6-1": "INT UNSIGNED",
    "6-2": "BIGINT",
    "6-3": "BIGINT",
    "7-0": "uint8",
    "7-1": "TINYINT UNSIGNED",
    "7-2": "SMALLINT",
    "7-3": "SMALLINT",
    "8-0": "uint16",
    "8-1": "SMALLINT UNSIGNED",
    "8-2": "INT",
    "8-3": "INT",
    "9-0": "uint32",
    "9-1": "INT UNSIGNED",
    "9-2": "BIGINT",
    "9-3": "BIGINT",
    "10-0": "uint64",
    "10-1": "BIGINT UNSIGNED",
    "10-2": "BIGINT",
    "10-3": "BIGINT",
    "11-0": "float32",
    "11-1": "FLOAT",
    "11-2": "FLOAT",
    "11-3": "FLOAT",
    "12-0": "float64",
    "12-1": "FLOAT",
    "12-2": "FLOAT",
    "12-3": "FLOAT",
    "13-0": "bool",
    "13-1": "BOOLEAN",
    "13-2": "BOOLEAN",
    "13-3": "BOOLEAN",
    "14-0": "time.Time or *time.Time",
    "14-1": "TIMESTAMP",
    "14-2": "TIMESTAMP",
    "14-3": "TIMESTAMP",
    "15-0": "other types",
    "15-1": "VARCHAR",
    "15-2": "VARCHAR",
    "15-3": "VARCHAR"
  },
  "cols": 4,
  "rows": 16
}
[/block]

[block:callout]
{
  "type": "danger",
  "title": "Need feedbacks in float mappings",
  "body": "It is currently not clear to map floating point types. Feedbacks and contributions are welcome!"
}
[/block]

[block:api-header]
{
  "type": "basic",
  "title": "Enforcing types"
}
[/block]
Type enforcements can be achieved using the `type` tag as in the following example;
[block:code]
{
  "codes": [
    {
      "code": "type Session struct {\n\t\tID string `qb:\"type:uuid; constraints:primary_key\"`\n}",
      "language": "go"
    }
  ]
}
[/block]
As you can see although the type defined is string, type tag enforces the type to be uuid.

Enforcing types is a great feature when you need to define database specific types such as `uuid`, `datetime`, `decimal`, etc. 
[block:code]
{
  "codes": [
    {
      "code": "CREATE TABLE session(\n\tid UUID PRIMARY KEY\n);",
      "language": "sql"
    }
  ]
}
[/block]

[block:api-header]
{
  "type": "basic",
  "title": "Indexing"
}
[/block]
Indexing is also supported in qb. The following example shows how to do indexing with using struct tag `index`;
[block:code]
{
  "codes": [
    {
      "code": "type Session struct {\n\t\tID                string `qb:\"type:uuid; constraints:primary_key\"`\n\t\tUserID            int64  `qb:\"constraints:ref(user.id)\"`\n\t\tAuthToken         string `qb:\"type:uuid\", constraints:not_null; index`\n\t\tqb.CompositeIndex `qb:\"index:id, user_id\"`\n}",
      "language": "go"
    }
  ]
}
[/block]
This definition would create the following sql statements within a transaction;
[block:code]
{
  "codes": [
    {
      "code": "CREATE TABLE session(\n\tid UUID PRIMARY KEY,\n\tuser_id BIGINT,\n\tauth_token UUID,\n\tFOREIGN KEY (user_id) REFERENCES user(id)\n);\n\nCREATE INDEX index_user_id ON session (user_id);\nCREATE INDEX index_id_user_id ON session (id, user_id);",
      "language": "sql"
    }
  ]
}
[/block]
As it might be noticed qb supports single indices as well as composite indices. This could be useful when there are selective queries with both querying a single column and multiple columns.