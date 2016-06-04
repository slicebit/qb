
# qb - the database toolkit for go

[![Join the chat at https://gitter.im/aacanakin/qb](https://badges.gitter.im/aacanakin/qb.svg)](https://gitter.im/aacanakin/qb?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Build Status](https://travis-ci.org/aacanakin/qb.svg?branch=master)](https://travis-ci.org/aacanakin/qb) [![Coverage Status](https://coveralls.io/repos/github/aacanakin/qb/badge.svg?branch=master)](https://coveralls.io/github/aacanakin/qb?branch=master) [![License (LGPL version 2.1)](https://img.shields.io/badge/license-GNU%20LGPL%20version%202.1-brightgreen.svg?style=flat)](http://opensource.org/licenses/LGPL-2.1) [![Go Report Card](https://goreportcard.com/badge/github.com/aacanakin/qb)](https://goreportcard.com/report/github.com/aacanakin/qb) [![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://godoc.org/github.com/aacanakin/qb)


**This project is currently pre 1.**

Currently, it's not feature complete. It can have potential bugs. There are no tests having concurrency race condition tests. It can crash especially in concurrency. 
Before 1.x releases, each major release could break backwards compatibility.

About qb
--------
qb is a database toolkit for easier db usage in go. It is inspired from python's most favorite orm sqlalchemy. qb is an orm as well as a query builder. It is quite modular in case of using just expression api and query building stuff.

[Documentation](https://qb.readme.io)
-------------
The documentation is hosted in [readme.io](https://qb.readme.io) which has great support for markdown docs. Currently, the docs is about 80% - 90% complete. The doc files will be added to this repo soon. Moreover, you can check the godoc from [here](https://godoc.org/github.com/aacanakin/qb). Contributions & Feedbacks in docs are welcome.

What's New (0.2)
----------------
The new table api design provides functionality for generating table objects like in the sqlalchemy's expression api.
Here's a full example to define a table;
```go
    db, err := qb.New("mysql", "root:@tcp(localhost:3306)/qb_test?charset=utf8")
	if err != nil {
		panic(err)
	}

	usersTable := qb.Table(
		"users",
		qb.Column("id", qb.Varchar().Size(40)),
		qb.Column("facebook_id", qb.BigInt()),
		qb.Column("email", qb.Varchar().Size(40).Unique()),
		qb.Column("device_id", qb.Varchar().Size(255).Unique()),
		qb.Column("session_id", qb.Varchar().Size(40)),
		qb.Column("auth_token", qb.Varchar().Size(40)),
		qb.Column("role_id", qb.Varchar().Size(40)),
		qb.PrimaryKey("id"),
		qb.ForeignKey().
			Ref("session_id", "sessions", "id").
			Ref("auth_token", "sessions", "auth_token").
			Ref("role_id", "roles", "id"),
		qb.UniqueKey("email", "device_id"),
	).Index("email", "device_id").Index("facebook_id")

	fmt.Println(usersTable.Create(db.Builder().Adapter()), "\n")

	db.Metadata().AddTable(usersTable)
	err = db.Metadata().CreateAll(db.Engine())
	if err != nil {
		panic(err)
	}
	
	// prints
	//CREATE TABLE users (
    //	session_id VARCHAR(40),
    //	auth_token VARCHAR(40),
    //	role_id VARCHAR(40),
    //	id VARCHAR(40),
    //	facebook_id BIGINT,
    //	email VARCHAR(40) UNIQUE,
    //	device_id VARCHAR(255) UNIQUE,
    //	PRIMARY KEY(id),
    //	FOREIGN KEY(session_id, auth_token) REFERENCES sessions(id, auth_token),
    //	FOREIGN KEY(role_id) REFERENCES roles(id),
    //	CONSTRAINT u_email_device_id UNIQUE(email, device_id)
    //);
    //CREATE INDEX i_email_device_id ON users(email, device_id);
    //CREATE INDEX i_facebook_id ON users(facebook_id);
	
```

Features
--------
- Support for postgres, mysql & sqlite3
- Simplistic query builder with no real magic
- Struct to table ddl mapper where initial table migrations can happen
- Expression builder which can be built almost any sql statements
- Transactional session api that auto map structs to queries
- Foreign Key definitions of structs using tags
- Single & composite column indices
- Relationships (soon..)

Installation
------------
Installation with glide;
```sh
glide get github.com/aacanakin/qb
```

Installation using go get;
```sh
go get -u github.com/aacanakin/qb
```
If you want to install test dependencies then;
```sh
go get -u -t github.com/aacanakin/qb
```

Quick Start
-----------
```go
package main

import (
	"fmt"
	"github.com/aacanakin/qb"
	"github.com/nu7hatch/gouuid"
)

type User struct {
	ID       string `qb:"type:uuid; constraints:primary_key"`
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
	db.Metadata().Add(User{})

	// create all tables registered to metadata
	db.Metadata().CreateAll()

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

	db.Metadata().DropAll() // drops all tables

}
```

Credits
-------
[Aras Can Akin](http://github.com/aacanakin)
[aodin](https://github.com/aodin)