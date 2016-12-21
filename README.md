![alt text](https://github.com/slicebit/qb/raw/master/qb_logo_128.png "qb: the database toolkit for go")

# qb - the database toolkit for go

[![Join the chat at https://gitter.im/aacanakin/qb](https://badges.gitter.im/aacanakin/qb.svg)](https://gitter.im/aacanakin/qb?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Build Status](https://travis-ci.org/slicebit/qb.svg?branch=master)](https://travis-ci.org/slicebit/qb)
[![Coverage Status](https://coveralls.io/repos/github/slicebit/qb/badge.svg?branch=master)](https://coveralls.io/github/slicebit/qb?branch=master)
[![License (LGPL version 2.1)](https://img.shields.io/badge/license-GNU%20LGPL%20version%202.1-brightgreen.svg?style=flat)](http://opensource.org/licenses/LGPL-2.1)
[![Go Report Card](https://goreportcard.com/badge/github.com/slicebit/qb)](https://goreportcard.com/report/github.com/slicebit/qb)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://godoc.org/github.com/slicebit/qb)

**This project is currently pre 1.**

Currently, it's not feature complete. It can have potential bugs. There are no tests covering concurrency race conditions. It can crash especially in concurrency.
Before 1.x releases, each major release could break backwards compatibility.

About qb
--------
qb is a database toolkit for easier db queries in go. It is inspired from python's best orm, namely sqlalchemy. qb is an orm(sqlx) as well as a query builder. It is quite modular in case of using just expression api and query building stuff.

[Documentation](https://qb.readme.io)
-------------
The documentation is hosted in [readme.io](https://qb.readme.io) which has great support for markdown docs. Currently, the docs are about 80% - 90% complete. The doc files will be added to this repo soon. Moreover, you can check the godoc from [here](https://godoc.org/github.com/slicebit/qb). Contributions & Feedbacks in docs are welcome.

Features
--------
- Support for postgres(9.5.+), mysql & sqlite3
- Powerful expression API for building queries & table ddls
- Struct to table ddl mapper where initial table migrations can happen
- Transactional session api that auto map structs to queries
- Foreign key definitions
- Single & Composite column indices
- Relationships (soon.. probably in 0.3 milestone)

Installation
------------
Installation with glide;
```sh
glide get github.com/slicebit/qb
```

0.2 installation with glide;
```sh
glide get github.com/slicebit/qb#0.2
```

Installation using go get;
```sh
go get -u github.com/slicebit/qb
```
If you want to install test dependencies then;
```sh
go get -u -t github.com/slicebit/qb
```

Quick Start
-----------
```go
package main

import (
	"fmt"
	"github.com/slicebit/qb"
	_ "github.com/mattn/go-sqlite3"
    _ "github.com/slicebit/qb/dialects/sqlite"
)

type User struct {
	ID       string `db:"id"`
	Email    string `db:"email"`
	FullName string `db:"full_name"`
	Oscars   int    `db:"oscars"`
}

func main() {

	users := qb.Table(
		"users",
		qb.Column("id", qb.Varchar().Size(40)),
		qb.Column("email", qb.Varchar()).NotNull().Unique(),
		qb.Column("full_name", qb.Varchar()).NotNull(),
		qb.Column("oscars", qb.Int()).NotNull().Default(0),
		qb.PrimaryKey("id"),
	)

	db, err := qb.New("sqlite3", "./qb_test.db")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	metadata := qb.MetaData()

	// add table to metadata
	metadata.AddTable(users)

	// create all tables registered to metadata
	metadata.CreateAll(db)
	defer metadata.DropAll(db) // drops all tables

	ins := qb.Insert(users).Values(map[string]interface{}{
		"id":        "b6f8bfe3-a830-441a-a097-1777e6bfae95",
		"email":     "jack@nicholson.com",
		"full_name": "Jack Nicholson",
	})

	_, err = db.Exec(ins)
	if err != nil {
		panic(err)
	}

	// find user
	var user User

	sel := qb.Select(users.C("id"), users.C("email"), users.C("full_name")).
		From(users).
		Where(users.C("id").Eq("b6f8bfe3-a830-441a-a097-1777e6bfae95"))

	err = db.Get(sel, &user)
	fmt.Printf("%+v\n", user)
}
```

Credits
-------
- [Aras Can Akın](https://github.com/aacanakin)
- [Christophe de Vienne](https://github.com/cdevienne)
- [Onur Şentüre](https://github.com/onursenture)
- [Aaron O. Ellis](https://github.com/aodin)
- [Shawn Smith](https://github.com/shawnps)
