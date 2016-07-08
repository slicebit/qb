![alt text](https://github.com/aacanakin/qb/raw/master/qb_logo_128.png "qb: the database toolkit for go")

# qb - the database toolkit for go

[![Join the chat at https://gitter.im/aacanakin/qb](https://badges.gitter.im/aacanakin/qb.svg)](https://gitter.im/aacanakin/qb?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Build Status](https://travis-ci.org/aacanakin/qb.svg?branch=master)](https://travis-ci.org/aacanakin/qb) [![Coverage Status](https://coveralls.io/repos/github/aacanakin/qb/badge.svg?branch=master)](https://coveralls.io/github/aacanakin/qb?branch=master) [![License (LGPL version 2.1)](https://img.shields.io/badge/license-GNU%20LGPL%20version%202.1-brightgreen.svg?style=flat)](http://opensource.org/licenses/LGPL-2.1) [![Go Report Card](https://goreportcard.com/badge/github.com/aacanakin/qb)](https://goreportcard.com/report/github.com/aacanakin/qb) [![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://godoc.org/github.com/aacanakin/qb)

**This project is currently pre 1.**

Currently, it's not feature complete. It can have potential bugs. There are no tests covering concurrency race conditions. It can crash especially in concurrency. 
Before 1.x releases, each major release could break backwards compatibility.

About qb
--------
qb is a database toolkit for easier db usage in go. It is inspired from python's most favorite orm sqlalchemy. qb is an orm as well as a query builder. It is quite modular in case of using just expression api and query building stuff.

[Documentation](https://qb.readme.io)
-------------
The documentation is hosted in [readme.io](https://qb.readme.io) which has great support for markdown docs. Currently, the docs are about 80% - 90% complete. The doc files will be added to this repo soon. Moreover, you can check the godoc from [here](https://godoc.org/github.com/aacanakin/qb). Contributions & Feedbacks in docs are welcome.

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
glide get github.com/aacanakin/qb
```

0.1 installation with glide;
```sh
glide get github.com/aacanakin/qb#0.1
```

Installation using go get;
```sh
go get -u github.com/aacanakin/qb
```
If you want to install test dependencies then;
```sh
go get -u -t github.com/aacanakin/qb
```

Quick Start - ORM
-----------------
```go
package main

import (
	"fmt"
	"github.com/aacanakin/qb"
	"github.com/nu7hatch/gouuid"
)

type User struct {
	ID       string `db:"_id" qb:"type:uuid; constraints:primary_key"`
	Email    string `qb:"constraints:unique, notnull"`
	FullName string `qb:"constraints:notnull"`
	Bio      string `qb:"type:text; constraints:null"`
}

func main() {

	db, err := qb.New("postgres", "user=postgres dbname=qb_test sslmode=disable")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	// add table to metadata
	db.AddTable(User{})

	// create all tables registered to metadata
	db.CreateAll()

	userID, _ := uuid.NewV4()
	db.Add(&User{
		ID:       userID.String(),
		Email:    "robert@de-niro.com",
		FullName: "Robert De Niro",
	})

	err = db.Commit() // insert user
	if err != nil {
	    fmt.Println(err)
	    return
	}

	var user User
	db.Find(&User{ID: userID.String()}).One(&user)

	fmt.Println("id", user.ID)
	fmt.Println("email", user.Email)
	fmt.Println("full_name", user.FullName)

	db.DropAll() // drops all tables

}
```

QuickStart - Expression API
---------------------------
```go
package main

import (
	"fmt"
	"github.com/aacanakin/qb"
)

func main() {
	db, _ := qb.New("sqlite3", ":memory:")
	defer db.Close()

	db.Dialect().SetEscaping(true)

	actors := qb.Table(
		"actor",
		qb.Column("id", qb.Varchar().Size(36)),
		qb.Column("name", qb.Varchar().NotNull()),
		qb.PrimaryKey("id"),
	)

	db.Metadata().AddTable(actors)
	err := db.CreateAll()
	if err != nil {
		panic(err)
	}

	ins := actors.Insert().Values(map[string]interface{}{
		"id":   "3af82cdc-4d21-473b-a175-cbc3f9119eda",
		"name": "Robert De Niro",
	})

	_, err = db.Engine().Exec(ins)
	if err != nil {
		panic(err)
	}

	sel := actors.
		Select(actors.C("name"), actors.C("id")).
		Where(actors.C("name").Eq("Robert De Niro"))

	var name string
	var id string

	db.Engine().QueryRow(sel).Scan(&name, &id)
	fmt.Printf("<User name=%s id=%s/>\n", name, id)

	// outputs
	// <User name=Robert De Niro id=3af82cdc-4d21-473b-a175-cbc3f9119eda/>
}
```

Credits
-------
- [Aras Can Akın](https://github.com/aacanakin)
- [Onur Şentüre](https://github.com/onursenture)
- [Aaron O. Ellis](https://github.com/aodin)
- [Shawn Smith](https://github.com/shawnps)
