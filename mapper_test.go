package qb

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type User struct {
	ID         string     `qb:"constraints:primary_key"`
	Email      string     `qb:"type:varchar(255); constraints:unique,notnull"`
	FullName   string     `qb:"constraints:notnull"`
	Password   string     `qb:"type:text"`
	FacebookID int64      `qb:"constraints:null"`
	UserType   string     `qb:"constraints:default(guest)"`
	CreatedAt  time.Time  `qb:"constraints:notnull"`
	UpdatedAt  time.Time  `qb:"constraints:notnull"`
	DeletedAt  *time.Time `qb:"constraints:null"`
}

type UserScore struct {
	UserID string `qb:"constraints:ref(user.id),primary_key"`
	Score  int64
}

type UserErr struct {
	ID    string `qb:"type:varchar(255);tag_should_raise_err:val;"`
	Email string `qb:"wrongtag:"`
}

func TestMapper(t *testing.T) {

	mapper := NewMapper("mysql")

	userTable, err := mapper.Convert(User{})
	userScoreTable, err := mapper.Convert(UserScore{})

	assert.Nil(t, err)
	fmt.Println(userTable.SQL())
	fmt.Println(userScoreTable.SQL())
}

func TestMapperError(t *testing.T) {

	mapper := NewMapper("postgres")

	userErrTable, err := mapper.Convert(UserErr{})

	fmt.Println(err)
	assert.NotNil(t, err)
	assert.Empty(t, userErrTable)
}
