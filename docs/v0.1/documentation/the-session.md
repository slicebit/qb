---
title: "The Session"
excerpt: ""
---
In qb, the session object is the most useful & simplest object. It has all the building blocks of qb. A typical session object has the following dependencies;
[block:code]
{
  "codes": [
    {
      "code": "// Session is the composition of engine connection & orm mappings\ntype Session struct {\n\tqueries  []*Query  // queries in the current transaction\n\tmapper   *Mapper   // mapper object for shortcut funcs\n\tmetadata *MetaData // metadata object to keep table registry\n\ttx       *sql.Tx   // active transaction\n\tbuilder  *Builder  // query builder\n  mutex    *sync.Mutex // mutex for preventing race conditions while opening a transaction\n}",
      "language": "go"
    }
  ]
}
[/block]
The following object is the only object required to use every single function in qb. The other structs that can also build sql statements such as `Builder` is transactionless.

From this point, the following qb initialization is assumed.
[block:code]
{
  "codes": [
    {
      "code": "package main\n\nimport (\n\t\"github.com/aacanakin/qb\"\n)\n\nfunc main() {\n\n\tdb, err := qb.New(\"postgres\", \"user=postgres dbname=qb_test sslmode=disable\")\n\tif err != nil {\n\t\tpanic(err)\n\t}\n\tdefer db.Close()\n\n\tdb.Metadata().Add(User{})\n\tdb.Metadata().CreateAll()\n}",
      "language": "go"
    }
  ]
}
[/block]

[block:api-header]
{
  "type": "basic",
  "title": "Inserting"
}
[/block]
To insert rows on a table just add sample models and call `db.Add(model)`; 
[block:code]
{
  "codes": [
    {
      "code": "// generate an insert statement and add it to current transaction\ndb.Add(User{Name: \"Aras Can Akin\"})\n\n// commit\ndb.Commit()",
      "language": "go"
    }
  ]
}
[/block]
You can add multiple models by calling Add function sequentially. This would generate the following sql statements and bindings, add it to the current transaction.
[block:code]
{
  "codes": [
    {
      "code": "INSERT INTO user\n(name)\nVALUES ($1);\n[Aras Can Akin]",
      "language": "sql"
    }
  ]
}
[/block]

[block:api-header]
{
  "type": "basic",
  "title": "Updating"
}
[/block]
To update rows, use Update(table).Set(map[string]interface{}) chain as in the following example;
[block:code]
{
  "codes": [
    {
      "code": "query := db.Update(\"user\").\n  Set(map[string]interface{}{\n    \"name\": \"Aras Akin\",\n  }).\n\tWhere(db.Eq(\"name\", \"Aras Can Akin\")).\n  Query()\n\ndb.AddQuery(query)\ndb.Commit()",
      "language": "go"
    }
  ]
}
[/block]
As it can be easily noticed `Update` statement is a little more unique for the flexibility of update statements.

This type of syntax is also supported in `Select` statements.

The following sql statement & bindings will be produced within a transaction;
[block:code]
{
  "codes": [
    {
      "code": "UPDATE user\nSET name = $1\nWHERE name = $2;\n[Aras Akin, Aras Can Akin]",
      "language": "sql"
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
Deletes are done by the following session call;
[block:code]
{
  "codes": [
    {
      "code": "// insert the row\ndb.Add(User{Name: \"Aras Can Akin\"})\ndb.Commit()\n\n// delete it\ndb.Delete(User{Name: \"Aras Can Akin\"})\ndb.Commit()",
      "language": "go"
    }
  ]
}
[/block]
The statement would produce the following sql statements & bindings within a transaction;
[block:code]
{
  "codes": [
    {
      "code": "DELETE FROM user\nWHERE name = $1;\n[Aras Can Akin]",
      "language": "sql"
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
There are more than one way to build select queries. In `qb.Session`, there exists a shortcut namely `Find()` that finds a model that matches the struct values as in the following;
[block:api-header]
{
  "type": "basic",
  "title": "Find.One()"
}
[/block]

[block:code]
{
  "codes": [
    {
      "code": "db.Add(User{Name: \"Aras Can Akin\"})\ndb.Commit()\n\nvar user User\nerr = db.Find(User{Name: \"Aras Can Akin\"}).One(&user)\nif err != nil {\n  fmt.Println(err)\n}\n\nfmt.Printf(\"id=%d, name=%s\\n\", user.ID, user.Name)",
      "language": "go"
    }
  ]
}
[/block]
The `Find()` call would produce the following sql statement with bindings;
[block:code]
{
  "codes": [
    {
      "code": "SELECT name, id\nFROM user\nWHERE name = $1;\n[Aras Can Akin]",
      "language": "sql"
    }
  ]
}
[/block]

[block:api-header]
{
  "type": "basic",
  "title": "Find().All()"
}
[/block]
`Find(model interface{}).All(models interface{})` returns all rows that is matched by struct example as in the following;
[block:code]
{
  "codes": [
    {
      "code": "db.Add(User{Name: \"Aras Can Akin\", Email: \"aras@gmail.com\"})\ndb.Add(User{Name: \"Aras Can Akin\", Email: \"aras@slicebit.com\"})\ndb.Commit()\n\nvar users []User\nerr = db.Find(User{Name: \"Aras Can Akin\"}).All(&users)\nif err != nil {\n  fmt.Println(err)\n}\n\nfor _, u := range users {\n  fmt.Printf(\"<User id=%d name=%s email=%s>\\n\", u.ID, u.Name, u.Email)\n}",
      "language": "go"
    }
  ]
}
[/block]
`Find().All()` call would produce the following sql statements and bindings;
[block:code]
{
  "codes": [
    {
      "code": "SELECT id, name, email\nFROM user\nWHERE name = $1;\n[Aras Can Akin]",
      "language": "sql"
    }
  ]
}
[/block]
[The Builder](doc:the-builder)] explained how sql statements are built by function chaining. The more complex select statements can be built by both builder and the session. Query building by func chaining can be also achieved using `qb.Session`. Lets make a complex selective query using joins;
[block:code]
{
  "codes": [
    {
      "code": "type User struct {\n  ID    int64  `qb:\"type:bigserial; constraints:primary_key\"`\n  Name  string `qb:\"constraints:not_null\"`\n  Email string `qb:\"constraints:not_null, unique\"`\n}\n\ntype Session struct {\n  UserID    int64  `qb:\"constraints:ref(user.id)\"`\n  AuthToken string `qb:\"type:uuid\"`\n  Agent     string `qb:\"constraints:not_null\"`\n}\n\ndb, err := qb.New(\"postgres\", \"user=postgres dbname=qb_test sslmode=disable\")\nif err != nil {\n  panic(err)\n}\ndefer db.Close()\n\ndb.Metadata().Add(User{})\ndb.Metadata().Add(Session{})\ndb.Metadata().CreateAll()\n\ndb.Add(User{Name: \"Aras Can Akin\", Email: \"aras@gmail.com\"})\ndb.Add(User{Name: \"Aras Can Akin\", Email: \"aras@slicebit.com\"})\n\ndb.Add(Session{\n  UserID:    1,\n  AuthToken: \"f1d2a8af-b048-479a-99c8-3725805299cf\",\n  Agent:     \"android\"})\n\ndb.Add(Session{\n  UserID:    1,\n  AuthToken: \"9bb95918-1cc1-4ab1-b3b1-29bf385d00bb\",\n  Agent:     \"chrome\"})\n\ndb.Add(Session{\n  UserID:    2,\n  AuthToken: \"0470f8ae-3e36-4a83-a4d1-7562173c48c6\",\n  Agent:     \"ios\"})\n\ndb.Commit()\n\nvar sessions []Session\nerr = db.\n  Select(\"s.user_id, s.auth_token, s.agent\").\n  From(\"session s\").\n  InnerJoin(\"user u\", \"s.user_id = u.id\").\n  Where(\"u.name = ?\", \"Aras Can Akin\").\n  OrderBy(\"s.agent\").\n  All(&sessions)\n\nif err != nil {\n  fmt.Println(err)\n}\n\nfor _, s := range sessions {\n  fmt.Printf(\"<Session user_id=%d auth_token=%s agent=%s>\\n\",\n             s.UserID,\n             s.AuthToken,\n             s.Agent)\n}\n\n// outputs\n// <Session user_id=1 auth_token=f1d2a8af-b048-479a-99c8-3725805299cf agent=android>\n// <Session user_id=1 auth_token=9bb95918-1cc1-4ab1-b3b1-29bf385d00bb agent=chrome>\n// <Session user_id=2 auth_token=0470f8ae-3e36-4a83-a4d1-7562173c48c6 agent=ios>",
      "language": "go"
    }
  ]
}
[/block]
Let's take a look at the last complex Select statement with use of `qb.Session`. The chain starting with `Select()` is just as the same as the builder way. However, builder doesn't have any `All()` functions which parses the struct and maps the values for each iteration in the result set.

**So, the very first question comes to mind is;**
Why do I have to bother with the builder api?

The `One()` and `All()` function calls uses sqlx's mapper functions which uses a lot of reflection. Therefore, it may slow down your query calls when dealing with large & complex structs. Therefore, performance critical apis should use builder api.