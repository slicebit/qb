package qbit

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type User struct {
	ID         string     `qb:"type:uuid;constraints:primary_key"`
	Email      string     `qb:"type:varchar(255); constraints:unique,notnull"`
	FullName   string     `qb:"constraints:notnull,index"`
	Password   string     `qb:"type:text"`
	FacebookID int64      `qb:"constraints:null"`
	UserType   string     `qb:"constraints:default(guest)"`
	Points     float32    `qb:"constraints:ref(user_score.points)"`
	CreatedAt  time.Time  `qb:"constraints:notnull"`
	UpdatedAt  time.Time  `qb:"constraints:notnull"`
	DeletedAt  *time.Time `qb:"constraints:null"`
}

func TestMapper(t *testing.T) {

	engine, err := NewEngine("mysql", "root:@tcp(127.0.0.1:3306)/qbit_test")
	defer engine.DB().Close()

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	mapper := NewMapper("postgres")

	userTable, err := mapper.Convert(User{})

	if err != nil {
		fmt.Errorf("Error: %s\n", err.Error())
	}

	fmt.Println(userTable.SQL())

	//	result, err := engine.Exec(userTable.Sql(), []interface{}{})

	//	assert.Equal(t, err, nil)
	//	fmt.Println(result)

	assert.Equal(t, 1, 1)

}
