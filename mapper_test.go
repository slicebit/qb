package qb

import (
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
	mapper := Mapper()

	userTable, err := mapper.ToTable(User{})
	assert.Nil(t, err)

	ddl := userTable.Create(dialect)

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
	mapper := Mapper()
	table, err := mapper.ToTable(User{})
	assert.Nil(t, err)
	ddl := table.Create(dialect)

	assert.Contains(t, ddl, "CREATE TABLE user (")
	assert.Contains(t, ddl, "id INTEGER PRIMARY KEY")
	assert.Contains(t, ddl, ")")
}

func TestMapperWithDBTag(t *testing.T) {
	type User struct {
		ID    string `db:"_id" qb:"type:varchar(36); constraints:primary_key"`
		Email string `qb:"constraints:unique, notnull"`
	}

	dialect := NewDialect("mysql")
	mapper := Mapper()
	table, err := mapper.ToTable(User{})
	assert.Nil(t, err)
	ddl := table.Create(dialect)

	assert.Contains(t, ddl, "CREATE TABLE user (")
	assert.Contains(t, ddl, "_id VARCHAR(36)")
	assert.Contains(t, ddl, "email VARCHAR(255) UNIQUE NOT NULL")
	assert.Contains(t, ddl, "PRIMARY KEY(_id)")

	m := mapper.ToMap(User{ID: "cba0667d-8c76-4999-9a55-84ffe572fb23", Email: "aras@slicebit.com"}, false)
	assert.Equal(t, map[string]interface{}{
		"_id":   "cba0667d-8c76-4999-9a55-84ffe572fb23",
		"email": "aras@slicebit.com",
	}, m)
}

func TestMapperPostgresAutoIncrement(t *testing.T) {
	type User struct {
		ID int64 `qb:"constraints:auto_increment"`
	}

	dialect := NewDialect("postgres")
	mapper := Mapper()
	table, err := mapper.ToTable(User{})
	assert.Nil(t, err)

	ddl := table.Create(dialect)
	assert.Contains(t, ddl, "id BIGSERIAL")
}

func TestMapperError(t *testing.T) {
	type UserErr struct {
		ID    string `qb:"type:varchar(255);tag_should_raise_err:val;"`
		Email string `qb:"wrongtag:"`
	}

	mapper := Mapper()

	userErrTable, err := mapper.ToTable(UserErr{})

	assert.NotNil(t, err)
	assert.Zero(t, userErrTable)
}

func TestMapperInvalidConstraint(t *testing.T) {
	type InvalidConstraint struct {
		ID string `qb:"constraints:invalid_constraint"`
	}

	mapper := Mapper()

	invalidConstraintTable, err := mapper.ToTable(InvalidConstraint{})

	assert.Zero(t, invalidConstraintTable)
	assert.NotNil(t, err)
}

func TestNonZeroStruct(t *testing.T) {
	type User struct {
		ID int
	}

	mapper := Mapper()
	m := mapper.ToMap(User{5}, false)
	assert.Equal(t, map[string]interface{}{"id": 5}, m)
}

func TestMapperUtilFuncs(t *testing.T) {
	type UserErr struct {
		ID    string `qb:"type:varchar(255);tag_should_raise_err:val;"`
		Email string `qb:"wrongtag:"`
	}

	mapper := Mapper()

	kv := mapper.ToMap(UserErr{}, false)
	assert.Equal(t, map[string]interface{}{}, kv)
}

func TestMapperTypes(t *testing.T) {
	mapper := Mapper()

	assert.Equal(t, Varchar().Size(255), mapper.ToType("string", ""))

	assert.Equal(t, Int(), mapper.ToType("int", ""))

	assert.Equal(t, TinyInt(), mapper.ToType("int8", ""))

	assert.Equal(t, SmallInt(), mapper.ToType("int16", ""))

	assert.Equal(t, Int(), mapper.ToType("int32", ""))

	assert.Equal(t, BigInt(), mapper.ToType("int64", ""))

	assert.Equal(t, Int().Unsigned(), mapper.ToType("uint", ""))

	assert.Equal(t, TinyInt().Unsigned(), mapper.ToType("uint8", ""))

	assert.Equal(t, SmallInt().Unsigned(), mapper.ToType("uint16", ""))

	assert.Equal(t, Int().Unsigned(), mapper.ToType("uint32", ""))

	assert.Equal(t, BigInt().Unsigned(), mapper.ToType("uint64", ""))

	assert.Equal(t, Float(), mapper.ToType("float32", ""))

	assert.Equal(t, Float(), mapper.ToType("float64", ""))

	assert.Equal(t, Boolean(), mapper.ToType("bool", ""))

	assert.Equal(t, Timestamp(), mapper.ToType("time.Time", ""))

	assert.Equal(t, Timestamp(), mapper.ToType("*time.Time", ""))
}
