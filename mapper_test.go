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

	mapper := NewMapper(NewBuilder("mysql"))

	userTable, err := mapper.ToTable(User{})

	assert.Nil(t, err)
	fmt.Println(userTable.SQL())
}

func TestMapperSqliteAutoIncrement(t *testing.T) {

	type User struct {
		ID int64 `qb:"constraints:auto_increment"`
	}

	mapper := NewMapper(NewBuilder("sqlite3"))
	table, err := mapper.ToTable(User{})

	assert.Nil(t, err)
	fmt.Println(table.SQL())
}

func TestMapperError(t *testing.T) {

	type UserErr struct {
		ID    string `qb:"type:varchar(255);tag_should_raise_err:val;"`
		Email string `qb:"wrongtag:"`
	}

	mapper := NewMapper(NewBuilder("postgres"))

	userErrTable, err := mapper.ToTable(UserErr{})

	assert.NotNil(t, err)
	assert.Empty(t, userErrTable)
}

type InvalidConstraint struct {
	ID string `qb:"constraints:invalid_constraint"`
}

func TestMapperInvalidConstraint(t *testing.T) {

	mapper := NewMapper(NewBuilder("mysql"))

	invalidConstraintTable, err := mapper.ToTable(InvalidConstraint{})

	assert.Nil(t, invalidConstraintTable)
	assert.NotNil(t, err)
}

func TestMapperUtilFuncs(t *testing.T) {

	type UserErr struct {
		ID    string `qb:"type:varchar(255);tag_should_raise_err:val;"`
		Email string `qb:"wrongtag:"`
	}

	mapper := NewMapper(NewBuilder("mysql"))

	assert.Equal(t, mapper.ColName("CreatedAt"), "created_at")

	kv := mapper.ToMap(UserErr{})
	assert.Equal(t, kv, map[string]interface{}{})
}

func TestMapperTypes(t *testing.T) {
	sqliteMapper := NewMapper(NewBuilder("sqlite3"))
	postgresMapper := NewMapper(NewBuilder("postgres"))
	mysqlMapper := NewMapper(NewBuilder("mysql"))

	assert.Equal(t, sqliteMapper.ToType("string", ""), &Type{"VARCHAR(255)"})
	assert.Equal(t, postgresMapper.ToType("string", ""), &Type{"VARCHAR(255)"})
	assert.Equal(t, mysqlMapper.ToType("string", ""), &Type{"VARCHAR(255)"})

	assert.Equal(t, sqliteMapper.ToType("int", ""), &Type{"INT"})
	assert.Equal(t, postgresMapper.ToType("int", ""), &Type{"INT"})
	assert.Equal(t, mysqlMapper.ToType("int", ""), &Type{"INT"})

	assert.Equal(t, sqliteMapper.ToType("int8", ""), &Type{"SMALLINT"})
	assert.Equal(t, postgresMapper.ToType("int8", ""), &Type{"SMALLINT"})
	assert.Equal(t, mysqlMapper.ToType("int8", ""), &Type{"SMALLINT"})

	assert.Equal(t, sqliteMapper.ToType("int16", ""), &Type{"SMALLINT"})
	assert.Equal(t, postgresMapper.ToType("int16", ""), &Type{"SMALLINT"})
	assert.Equal(t, mysqlMapper.ToType("int16", ""), &Type{"SMALLINT"})

	assert.Equal(t, sqliteMapper.ToType("int32", ""), &Type{"INT"})
	assert.Equal(t, postgresMapper.ToType("int32", ""), &Type{"INT"})
	assert.Equal(t, mysqlMapper.ToType("int32", ""), &Type{"INT"})

	assert.Equal(t, sqliteMapper.ToType("int64", ""), &Type{"BIGINT"})
	assert.Equal(t, postgresMapper.ToType("int64", ""), &Type{"BIGINT"})
	assert.Equal(t, mysqlMapper.ToType("int64", ""), &Type{"BIGINT"})

	assert.Equal(t, sqliteMapper.ToType("uint", ""), &Type{"BIGINT"})
	assert.Equal(t, postgresMapper.ToType("uint", ""), &Type{"BIGINT"})
	assert.Equal(t, mysqlMapper.ToType("uint", ""), &Type{"INT UNSIGNED"})

	assert.Equal(t, sqliteMapper.ToType("uint8", ""), &Type{"SMALLINT"})
	assert.Equal(t, postgresMapper.ToType("uint8", ""), &Type{"SMALLINT"})
	assert.Equal(t, mysqlMapper.ToType("uint8", ""), &Type{"TINYINT UNSIGNED"})

	assert.Equal(t, sqliteMapper.ToType("uint16", ""), &Type{"INT"})
	assert.Equal(t, postgresMapper.ToType("uint16", ""), &Type{"INT"})
	assert.Equal(t, mysqlMapper.ToType("uint16", ""), &Type{"SMALLINT UNSIGNED"})

	assert.Equal(t, sqliteMapper.ToType("uint32", ""), &Type{"BIGINT"})
	assert.Equal(t, postgresMapper.ToType("uint32", ""), &Type{"BIGINT"})
	assert.Equal(t, mysqlMapper.ToType("uint32", ""), &Type{"INT UNSIGNED"})

	assert.Equal(t, sqliteMapper.ToType("uint64", ""), &Type{"BIGINT"})
	assert.Equal(t, postgresMapper.ToType("uint64", ""), &Type{"BIGINT"})
	assert.Equal(t, mysqlMapper.ToType("uint64", ""), &Type{"BIGINT UNSIGNED"})

	assert.Equal(t, sqliteMapper.ToType("float32", ""), &Type{"FLOAT"})
	assert.Equal(t, postgresMapper.ToType("float32", ""), &Type{"FLOAT"})
	assert.Equal(t, mysqlMapper.ToType("float32", ""), &Type{"FLOAT"})

	assert.Equal(t, sqliteMapper.ToType("float64", ""), &Type{"FLOAT"})
	assert.Equal(t, postgresMapper.ToType("float64", ""), &Type{"FLOAT"})
	assert.Equal(t, mysqlMapper.ToType("float64", ""), &Type{"FLOAT"})

	assert.Equal(t, sqliteMapper.ToType("bool", ""), &Type{"BOOLEAN"})
	assert.Equal(t, postgresMapper.ToType("bool", ""), &Type{"BOOLEAN"})
	assert.Equal(t, mysqlMapper.ToType("bool", ""), &Type{"BOOLEAN"})

	assert.Equal(t, sqliteMapper.ToType("time.Time", ""), &Type{"TIMESTAMP"})
	assert.Equal(t, postgresMapper.ToType("time.Time", ""), &Type{"TIMESTAMP"})
	assert.Equal(t, mysqlMapper.ToType("time.Time", ""), &Type{"TIMESTAMP"})

	assert.Equal(t, sqliteMapper.ToType("*time.Time", ""), &Type{"TIMESTAMP"})
	assert.Equal(t, postgresMapper.ToType("*time.Time", ""), &Type{"TIMESTAMP"})
	assert.Equal(t, mysqlMapper.ToType("*time.Time", ""), &Type{"TIMESTAMP"})
}
