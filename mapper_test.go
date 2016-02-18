package qb

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type UnknownType struct{}

type User struct {
	ID          string `qb:"constraints:primary_key"`
	FacebookID  int64  `qb:"constraints:ref(facebook.id)"`
	ProfileID   int64  `qb:"constraints:ref(profile.id)"`
	ProfileName string `qb:"constraints:ref(profile.name)"`
	Email       string `qb:"type:varchar(255); constraints:unique,notnull"`
	FullName    string `qb:"constraints:notnull,default"`
	Password    string `qb:"type:text"`
	UserType    string `qb:"constraints:default(guest)"`
	Premium     bool
	CreatedAt   time.Time  `qb:"constraints:notnull"`
	DeletedAt   *time.Time `qb:"constraints:null"`
	Level       int
	Money       float32
	Score       float64
	Unknown     UnknownType
}

func TestMapper(t *testing.T) {

	mapper := NewMapper("mysql")

	userTable, err := mapper.Convert(User{})

	assert.Nil(t, err)
	fmt.Println(userTable.SQL())
	//	fmt.Println(userScoreTable.SQL())
}

type UserErr struct {
	ID    string `qb:"type:varchar(255);tag_should_raise_err:val;"`
	Email string `qb:"wrongtag:"`
}

func TestMapperError(t *testing.T) {

	mapper := NewMapper("postgres")

	userErrTable, err := mapper.Convert(UserErr{})

	assert.NotNil(t, err)
	assert.Empty(t, userErrTable)
}

type InvalidConstraint struct {
	ID string `qb:"constraints:invalid_constraint"`
}

func TestMapperInvalidConstraint(t *testing.T) {

	mapper := NewMapper("mysql")

	invalidConstraintTable, err := mapper.Convert(InvalidConstraint{})

	assert.Nil(t, invalidConstraintTable)
	assert.NotNil(t, err)
}
