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
		ID              string `qb:"constraints:primary_key"`
		SecondaryID     string `qb:"constraints:primary_key"`
		FacebookID      int64  `qb:"constraints:ref(facebook.id)"`
		FacebookProfile string `qb:"constraints:ref(facebook.profile_id)"`
		ProfileID       int64  `qb:"constraints:ref(profile.id)"`
		ProfileName     string `qb:"constraints:ref(profile.name)"`
		Email           string `qb:"type:varchar(255); constraints:unique,notnull"`
		FullName        string `qb:"constraints:notnull,default"`
		Password        string `qb:"type:text"`
		UserType        string `qb:"constraints:default(guest)"`
		Ignored         bool   `qb:"-"`
		Premium         bool
		CreatedAt       time.Time  `qb:"constraints:notnull"`
		DeletedAt       *time.Time `qb:"constraints:null;"`
		Level           int
		Money           float32
		Score           float64
		Unknown         UnknownType `qb:"index"`
		CompositeIndex  `qb:"index:full_name, password"`
	}

	dialect := NewDialect("mysql")
	mapper := Mapper(dialect)

	userTable, err := mapper.ToTable(User{})
	assert.Nil(t, err)

	ddl := userTable.Create(dialect)
	fmt.Println(ddl, "\n")

	assert.Contains(t, ddl, "CREATE TABLE user (")
	assert.Contains(t, ddl, "id VARCHAR(255)")
	assert.Contains(t, ddl, "secondary_id VARCHAR(255)")
	assert.Contains(t, ddl, "facebook_id BIGINT")
	assert.Contains(t, ddl, "facebook_profile VARCHAR(255)")
	assert.Contains(t, ddl, "profile_id BIGINT")
	assert.Contains(t, ddl, "profile_name VARCHAR(255)")
	assert.Contains(t, ddl, "email VARCHAR(255) UNIQUE NOT NULL")
	assert.Contains(t, ddl, "full_name VARCHAR(255) NOT NULL DEFAULT ''")
	assert.Contains(t, ddl, "password TEXT")
	assert.Contains(t, ddl, "user_type VARCHAR(255) DEFAULT 'guest'")
	assert.Contains(t, ddl, "premium BOOLEAN")
	assert.Contains(t, ddl, "created_at TIMESTAMP NOT NULL")
	assert.Contains(t, ddl, "deleted_at TIMESTAMP NULL")
	assert.Contains(t, ddl, "level INT")
	assert.Contains(t, ddl, "money FLOAT")
	assert.Contains(t, ddl, "score FLOAT")
	assert.Contains(t, ddl, "unknown VARCHAR(255)")
	assert.Contains(t, ddl, "PRIMARY KEY(id, secondary_id)")
	assert.Contains(t, ddl, "FOREIGN KEY(facebook_id, facebook_profile) REFERENCES facebook(id, profile_id)")
	assert.Contains(t, ddl, "FOREIGN KEY(profile_id, profile_name) REFERENCES profile(id, name)")
	assert.Contains(t, ddl, ")")
	assert.Contains(t, ddl, "CREATE INDEX i_unknown ON user(unknown)")
	assert.Contains(t, ddl, "CREATE INDEX i_full_name_password ON user(full_name, password);")
}

func TestMapperSqliteAutoIncrement(t *testing.T) {
	type User struct {
		ID int64 `qb:"constraints:auto_increment"`
	}

	dialect := NewDialect("sqlite3")
	mapper := Mapper(dialect)
	table, err := mapper.ToTable(User{})
	assert.Nil(t, err)
	ddl := table.Create(dialect)

	assert.Contains(t, ddl, "CREATE TABLE user (")
	assert.Contains(t, ddl, "id BIGINT")
	assert.Contains(t, ddl, ")")

	fmt.Println(ddl, "\n")
}

func TestMapperWithDBTag(t *testing.T) {
	type User struct {
		ID    string `db:"_id" qb:"type:varchar(36); constraints:primary_key"`
		Email string `qb:"constraints:unique, notnull"`
	}

	dialect := NewDialect("mysql")
	mapper := Mapper(dialect)
	table, err := mapper.ToTable(User{})
	assert.Nil(t, err)
	ddl := table.Create(dialect)

	assert.Contains(t, ddl, "CREATE TABLE user (")
	assert.Contains(t, ddl, "_id VARCHAR(36)")
	assert.Contains(t, ddl, "email VARCHAR(255) UNIQUE NOT NULL")
	assert.Contains(t, ddl, "PRIMARY KEY(_id)")

	m := mapper.ToMap(User{ID: "cba0667d-8c76-4999-9a55-84ffe572fb23", Email: "aras@slicebit.com"}, false)
	assert.Equal(t, m, map[string]interface{}{
		"_id":   "cba0667d-8c76-4999-9a55-84ffe572fb23",
		"email": "aras@slicebit.com",
	})

	fmt.Println(ddl, "\n")
}

func TestMapperPostgresAutoIncrement(t *testing.T) {
	type User struct {
		ID int64 `qb:"constraints:auto_increment"`
	}

	dialect := NewDialect("postgres")
	mapper := Mapper(dialect)
	table, err := mapper.ToTable(User{})
	assert.Nil(t, err)

	ddl := table.Create(dialect)
	assert.NotContains(t, ddl, "AUTOINCREMENT")
	assert.NotContains(t, ddl, "AUTO INCREMENT")
}

func TestMapperError(t *testing.T) {
	type UserErr struct {
		ID    string `qb:"type:varchar(255);tag_should_raise_err:val;"`
		Email string `qb:"wrongtag:"`
	}

	dialect := NewDialect("postgres")
	mapper := Mapper(dialect)

	userErrTable, err := mapper.ToTable(UserErr{})

	fmt.Println(err)
	assert.NotNil(t, err)
	assert.Zero(t, userErrTable)
}

func TestMapperInvalidConstraint(t *testing.T) {
	type InvalidConstraint struct {
		ID string `qb:"constraints:invalid_constraint"`
	}

	dialect := NewDialect("mysql")
	mapper := Mapper(dialect)

	invalidConstraintTable, err := mapper.ToTable(InvalidConstraint{})

	assert.Zero(t, invalidConstraintTable)
	assert.NotNil(t, err)
}

func TestNonZeroStruct(t *testing.T) {
	type User struct {
		ID int
	}

	mapper := Mapper(NewDialect("mysql"))
	m := mapper.ToMap(User{5}, false)
	assert.Equal(t, m, map[string]interface{}{"id": 5})
}

func TestMapperUtilFuncs(t *testing.T) {
	type UserErr struct {
		ID    string `qb:"type:varchar(255);tag_should_raise_err:val;"`
		Email string `qb:"wrongtag:"`
	}

	mapper := Mapper(NewDialect("mysql"))

	kv := mapper.ToMap(UserErr{}, false)
	assert.Equal(t, kv, map[string]interface{}{})
}

func TestMapperTypes(t *testing.T) {
	sqliteMapper := Mapper(NewDialect("sqlite3"))
	postgresMapper := Mapper(NewDialect("postgres"))
	mysqlMapper := Mapper(NewDialect("mysql"))

	assert.Equal(t, sqliteMapper.ToType("string", ""), Varchar().Size(255))
	assert.Equal(t, postgresMapper.ToType("string", ""), Varchar().Size(255))
	assert.Equal(t, mysqlMapper.ToType("string", ""), Varchar().Size(255))

	assert.Equal(t, sqliteMapper.ToType("int", ""), Int())
	assert.Equal(t, postgresMapper.ToType("int", ""), Int())
	assert.Equal(t, mysqlMapper.ToType("int", ""), Int())

	assert.Equal(t, sqliteMapper.ToType("int8", ""), SmallInt())
	assert.Equal(t, postgresMapper.ToType("int8", ""), SmallInt())
	assert.Equal(t, mysqlMapper.ToType("int8", ""), SmallInt())

	assert.Equal(t, sqliteMapper.ToType("int16", ""), SmallInt())
	assert.Equal(t, postgresMapper.ToType("int16", ""), SmallInt())
	assert.Equal(t, mysqlMapper.ToType("int16", ""), SmallInt())

	assert.Equal(t, sqliteMapper.ToType("int32", ""), Int())
	assert.Equal(t, postgresMapper.ToType("int32", ""), Int())
	assert.Equal(t, mysqlMapper.ToType("int32", ""), Int())

	assert.Equal(t, sqliteMapper.ToType("int64", ""), BigInt())
	assert.Equal(t, postgresMapper.ToType("int64", ""), BigInt())
	assert.Equal(t, mysqlMapper.ToType("int64", ""), BigInt())

	assert.Equal(t, sqliteMapper.ToType("uint", ""), BigInt())
	assert.Equal(t, postgresMapper.ToType("uint", ""), BigInt())
	assert.Equal(t, mysqlMapper.ToType("uint", ""), Type("INT UNSIGNED"))

	assert.Equal(t, sqliteMapper.ToType("uint8", ""), SmallInt())
	assert.Equal(t, postgresMapper.ToType("uint8", ""), SmallInt())
	assert.Equal(t, mysqlMapper.ToType("uint8", ""), Type("TINYINT UNSIGNED"))

	assert.Equal(t, sqliteMapper.ToType("uint16", ""), Int())
	assert.Equal(t, postgresMapper.ToType("uint16", ""), Int())
	assert.Equal(t, mysqlMapper.ToType("uint16", ""), Type("SMALLINT UNSIGNED"))

	assert.Equal(t, sqliteMapper.ToType("uint32", ""), BigInt())
	assert.Equal(t, postgresMapper.ToType("uint32", ""), BigInt())
	assert.Equal(t, mysqlMapper.ToType("uint32", ""), Type("INT UNSIGNED"))

	assert.Equal(t, sqliteMapper.ToType("uint64", ""), BigInt())
	assert.Equal(t, postgresMapper.ToType("uint64", ""), BigInt())
	assert.Equal(t, mysqlMapper.ToType("uint64", ""), Type("BIGINT UNSIGNED"))

	assert.Equal(t, sqliteMapper.ToType("float32", ""), Float())
	assert.Equal(t, postgresMapper.ToType("float32", ""), Float())
	assert.Equal(t, mysqlMapper.ToType("float32", ""), Float())

	assert.Equal(t, sqliteMapper.ToType("float64", ""), Float())
	assert.Equal(t, postgresMapper.ToType("float64", ""), Float())
	assert.Equal(t, mysqlMapper.ToType("float64", ""), Float())

	assert.Equal(t, sqliteMapper.ToType("bool", ""), Boolean())
	assert.Equal(t, postgresMapper.ToType("bool", ""), Boolean())
	assert.Equal(t, mysqlMapper.ToType("bool", ""), Boolean())

	assert.Equal(t, sqliteMapper.ToType("time.Time", ""), Timestamp())
	assert.Equal(t, postgresMapper.ToType("time.Time", ""), Timestamp())
	assert.Equal(t, mysqlMapper.ToType("time.Time", ""), Timestamp())

	assert.Equal(t, sqliteMapper.ToType("*time.Time", ""), Timestamp())
	assert.Equal(t, postgresMapper.ToType("*time.Time", ""), Timestamp())
	assert.Equal(t, mysqlMapper.ToType("*time.Time", ""), Timestamp())
}
