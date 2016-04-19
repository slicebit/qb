
# qb - the database toolkit for go

[![Join the chat at https://gitter.im/aacanakin/qb](https://badges.gitter.im/aacanakin/qb.svg)](https://gitter.im/aacanakin/qb?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Build Status](https://travis-ci.org/aacanakin/qb.svg?branch=master)](https://travis-ci.org/aacanakin/qb) [![Coverage Status](https://coveralls.io/repos/github/aacanakin/qb/badge.svg?branch=master)](https://coveralls.io/github/aacanakin/qb?branch=master) [![License (LGPL version 2.1)](https://img.shields.io/badge/license-GNU%20LGPL%20version%202.1-brightgreen.svg?style=flat)](http://opensource.org/licenses/LGPL-2.1) [![Go Report Card](https://goreportcard.com/badge/github.com/aacanakin/qb)](https://goreportcard.com/report/github.com/aacanakin/qb) [![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://godoc.org/github.com/aacanakin/qb)


**This project is currently pre 1.**
More documentation will be coming soon.
Although the tests coverage are high, currently, it is not recommended to use it in production. It can currently crash especially in concurrency.

About qb
--------
qb is a database toolkit for easier db usage in go. It is inspired from python's most favorite orm sqlalchemy. qb is an orm as well as a query builder. It is quite modular in case of using just expression api and query building stuff.

Features
--------
- Support for postgres, mysql & sqlite
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
	fmt.Println(err)

	var user User
	db.Find(&User{ID: userID.String()}).One(&user)

	fmt.Println("id", user.ID)
	fmt.Println("email", user.Email)
	fmt.Println("full_name", user.FullName)

	db.Metadata().DropAll() // drops all tables

}
```
