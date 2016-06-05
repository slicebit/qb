---
title: "Introduction"
excerpt: ""
---
qb is a database toolkit for easier implementation for db heavy apps in go. The aim is to build a lightweight architecture on top of [sqlx](https://github.com/jmoiron/sqlx) library for painless db implementation in go. It is inspired from python's most favorite sql package called [sqlalchemy](http://www.sqlalchemy.org/)
[block:api-header]
{
  "type": "basic",
  "title": "Features"
}
[/block]
- Support for postgres, mysql, mariadb & sqlite
- A simple query builder api
- A transactional session api implemented on top of query builder that have commonly used db queries such as Find(), All(), etc.
- A metadata api where tables can be generated and created from structs
- A struct to table mapper with use of tagging
- A query to struct mapper with use of sqlx library
- Foreign key reference definitions with use of tagging
- Single & composite column indices with use of tagging
- Relationships (coming soon..)
[block:api-header]
{
  "type": "basic",
  "title": "Raison d'être"
}
[/block]
The reason this package is being developed is mainly because of my personal curiosity. It is currently a hobby project of mine. However, there are more several reasons. I have played with most of the ormish libraries in go. At the time, neither of them were complete and there is a quite good [post](http://www.hydrogen18.com/blog/golang-orms-and-why-im-still-not-using-one.html) about that which resonates the point.

Moreover, there is this tweet I had posted about the go orm world;

[block:image]
{
  "images": [
    {
      "image": [
        "https://www.filepicker.io/api/file/oO46wdbToSY2NVZKYRnQ",
        "Screen Shot 2016-03-08 at 12.05.18 AM.png",
        "1286",
        "710",
        "#333237",
        ""
      ],
      "sizing": "smart",
      "border": true
    }
  ]
}
[/block]
From my perspective, I think qb solves most of the problems I'm suffering when using an orm in go and hopefully, it would be useful to anyone that has the similar problems with me.
[block:api-header]
{
  "type": "basic",
  "title": "Installation"
}
[/block]
To install qb, a simple go get would do the trick;
[block:code]
{
  "codes": [
    {
      "code": "go get -v -u github.com/aacanakin/qb",
      "language": "shell"
    }
  ]
}
[/block]
To get the test dependencies add a -t flag;
[block:code]
{
  "codes": [
    {
      "code": "go get -v -u -t github.com/aacanakin/qb",
      "language": "shell"
    }
  ]
}
[/block]
Moreover, [glide](https://glide.sh/) is also supported;
[block:code]
{
  "codes": [
    {
      "code": "glide get github.com/aacanakin/qb",
      "language": "shell"
    }
  ]
}
[/block]

[block:api-header]
{
  "type": "basic",
  "title": "Quick Start"
}
[/block]

[block:code]
{
  "codes": [
    {
      "code": "package main\n\nimport (\n    \"fmt\"\n    \"github.com/aacanakin/qb\"\n    \"github.com/nu7hatch/gouuid\"\n)\n\ntype User struct {\n    ID       string `qb:\"type:uuid; constraints:primary_key\"`\n    Email    string `qb:\"constraints:unique, notnull\"`\n    FullName string `qb:\"constraints:notnull\"`\n    Bio      string `qb:\"type:text; constraints:null\"`\n}\n\nfunc main() {\n\n    db, err := qb.New(\"postgres\", \"user=postgres dbname=qb_test sslmode=disable\")\n    if err != nil {\n        panic(err)\n    }\n\n    defer db.Close()\n\n    // add table to metadata\n    db.Metadata().Add(User{})\n\n    // create all tables registered to metadata\n    db.Metadata().CreateAll()\n\n    userID, _ := uuid.NewV4()\n    db.Add(&User{\n        ID:       userID.String(),\n        Email:    \"robert@de-niro.com\",\n        FullName: \"Robert De Niro\",\n    })\n\n    err = db.Commit() // insert user\n    fmt.Println(err)\n\n    var user User\n    db.Find(&User{ID: userID.String()}).First(&user)\n\n    fmt.Println(\"id\", user.ID)\n    fmt.Println(\"email\", user.Email)\n    fmt.Println(\"full_name\", user.FullName)\n\n    db.Metadata().DropAll() // drops all tables\n\n}",
      "language": "go"
    }
  ]
}
[/block]