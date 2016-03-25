package qb

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMapper(t *testing.T) {

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

	mapper := NewMapper("mysql")

	userTable, err := mapper.ToTable(User{})

	assert.Nil(t, err)
	fmt.Println(userTable.SQL())
}

func TestMapperSqliteAutoIncrement(t *testing.T) {

	type User struct {
		ID int64 `qb:"constraints:auto_increment"`
	}

	mapper := NewMapper("sqlite3")
	table, err := mapper.ToTable(User{})

	assert.Nil(t, err)
	fmt.Println(table.SQL())
}

func TestMapperError(t *testing.T) {

	type UserErr struct {
		ID    string `qb:"type:varchar(255);tag_should_raise_err:val;"`
		Email string `qb:"wrongtag:"`
	}

	mapper := NewMapper("postgres")

	userErrTable, err := mapper.ToTable(UserErr{})

	assert.NotNil(t, err)
	assert.Empty(t, userErrTable)
}

type InvalidConstraint struct {
	ID string `qb:"constraints:invalid_constraint"`
}

func TestMapperInvalidConstraint(t *testing.T) {

	mapper := NewMapper("mysql")

	invalidConstraintTable, err := mapper.ToTable(InvalidConstraint{})

	assert.Nil(t, invalidConstraintTable)
	assert.NotNil(t, err)
}

func TestMapperUtilFuncs(t *testing.T) {

	type UserErr struct {
		ID    string `qb:"type:varchar(255);tag_should_raise_err:val;"`
		Email string `qb:"wrongtag:"`
	}

	mapper := NewMapper("mysql")

	assert.Equal(t, mapper.ColName("CreatedAt"), "created_at")

	kv := mapper.ToMap(UserErr{})
	assert.Equal(t, kv, map[string]interface{}{})
}
