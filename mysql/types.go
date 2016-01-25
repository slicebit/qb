package mysql

import (
	"fmt"
	"github.com/aacanakin/qbit"
	"strings"
)

// generates a float type in mysql syntax
// if only p is specified, p is the binary precision.
// if p and s are both specified, p is the maximum number of all digits (both sides of the decimal point),
// s is the maximum number of digits after the point. p and s are optional
func Float(p int, s int) *qbit.Type {
	return &qbit.Type{
		Sql: func() string {
			return fmt.Sprintf("FLOAT(%d, %d)", p, s)
		},
	}
}

// generates a date type in mysql syntax
func Date() *qbit.Type {
	return &qbit.Type{
		Sql: func() string {
			return "DATE"
		},
	}
}

// generates a date time type in mysql syntax
func DateTime() *qbit.Type {
	return &qbit.Type{
		Sql: func() string {
			return "DATETIME"
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

// generates a char type in mysql syntax
func Char(size int) *qbit.Type {
	return &qbit.Type{
		Sql: func() string {
			return fmt.Sprintf("CHAR(%d)", size)
		},
	}
}

// generates a varchar type in mysql syntax
func VarChar(size int) *qbit.Type {
	return &qbit.Type{
		Sql: func() string {
			return fmt.Sprintf("VARCHAR(%d)", size)
		},
	}
}

// generates a text type in mysql syntax
func Text() *qbit.Type {
	return &qbit.Type{
		Sql: func() string {
			return "TEXT"
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
