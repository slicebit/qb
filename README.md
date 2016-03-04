# qb - the database toolkit for go
**This project is currently pre 1.**

Although the tests coverage are high, currently, it is not recommended to use it in production. It can currently crash especially in high concurrency.

About qb
--------
qb is a database toolkit for easier db usage in go. It is inspired from python's most favorite orm sqlalchemy. qb is an orm as well as a query builder. It is

Features
--------
- Support for postgres, mysql & sqlite
- Simplistic query builder with no real magic
- Struct to table ddl mapper where initial table migrations can happen
- Expression api which can be built almost any sql statements
- Transactional session api that auto map structs to queries
- Foreign Key definitions of structs using tags
- Relationships (soon..)

Quick Start - Session API
-------------------------
```go
import (
    "github.com/aacanakin/qb"
)

type User struct {
	ID       string `qb:"type:uuid; constraints:primary_key"`
	Email    string `qb:"constraints:unique, notnull"`
	FullName string `qb:"constraints:notnull"`
	Password string `qb:"constraints:notnull"`
	Bio      string `qb:"type:text; constraints:null"`
}

func main() {
    engine, err := qb.NewEngine("postgres", "user=postgres dbname=qb_test sslmode=disable")

    if err != nil {
        panic(err)
    }

    metadata := qb.NewMetadata(engine)
    session := qb.NewSession(metadata)

    session.Metadata().Add(&User{}, &Session{})
    err = session.Metadata().CreateAll()
    // Creates user table

    // insert user using session
    rdnId, _ := uuid.NewV4()
	rdn := &User{
		ID:       rdnId.String(),
		Email:    "robert@de-niro.com",
		FullName: "Robert De Niro",
		Password: "rdn",
	}

    session.Add(rdn)
    err = session.Commit()
    // inserts Robert De Niro to users table
}
```

Quick Start - Expression API
----------------------------
```go
// incoming...
```