# qb - the database toolkit for go
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
- Relationships (soon..)

Quick Start
-----------
```go
import (
    "github.com/aacanakin/qb"
    "github.com/nu7hatch/gouuid"
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

    session.Metadata().Add(&User{})
    err = session.Metadata().CreateAll() // Creates user table

    // insert user using session
    userID, _ := uuid.NewV4()
	user := &User{
		ID:       userID.String(),
		Email:    "robert@de-niro.com",
		FullName: "Robert De Niro",
		Password: "rdn",
	}

    session.Add(user)
    err = session.Commit()
    // user is inserted into db
}
```