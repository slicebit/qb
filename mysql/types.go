package mysql

import (
	"fmt"
	"github.com/aacanakin/qbit"
	"strings"
)

// generates a date time type in mysql syntax
func DateTime() *qbit.Type {
	return &qbit.Type{
		Sql: func() string {
			return "DATETIME"
		},
	}
}

func Timestamp() *qbit.Type {
	return &qbit.Type{
		Sql: func() string {
			return "TIMESTAMP"
		},
	}
}

// generates a date time type in mysql syntax
func Time() *qbit.Type {
	return &qbit.Type{
		Sql: func() string {
			return "TIME"
		},
	}
}

// generates a tiny-text type in mysql syntax
func TinyText() *qbit.Type {
	return &qbit.Type{
		Sql: func() string {
			return "TINYTEXT"
		},
	}
}

// generates an enum type in mysql syntax
func Enum(vals []interface{}) *qbit.Type {

	for k, _ := range vals {
		vals[k] = fmt.Sprintf("'%s'", vals[k])
	}

	return &qbit.Type{
		Sql: func() string {
			return fmt.Sprintf("ENUM(%s)", strings.Join(vals, ", "))
		},
	}
}
