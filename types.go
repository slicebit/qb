package qbit

import (
	"fmt"
	"strings"
)

type Type struct {
	Sql func() string
}

// generates common char type syntax
// mysql: 0 ⇐ n ⇐ 255, default is 1
// postgres: 1 ⇐ n ⇐ 1G, default is 1
func Char(size int) *Type {
	return &Type{
		Sql: func() string {
			return fmt.Sprintf("CHAR(%d)", size)
		},
	}
}

// generates common varchar type syntax
// mysql: 1 ⇐ n ⇐ 65535/charsize, n is mandatory
// postgres: 1 ⇐ n ⇐ 1G, default is 1G
func VarChar(size int) *Type {
	return &Type{
		Sql: func() string {
			return fmt.Sprintf("VARCHAR(%d)", size)
		},
	}
}

// generates common text type syntax
// mysql: 65555 bytes
// postgres: 1G
func Text() *Type {
	return &Type{
		Sql: func() string {
			return "TEXT"
		},
	}
}

// generates mysql mediumtext type syntax
// mysql: 16,777,215 (2^24 – 1) bytes
func MediumText() *Type {
	return &Type{
		Sql: func() string {
			return "MEDIUMTEXT"
		},
	}
}

// generates mysql longtext type syntax
// mysql: 4G (2^32 – 1) bytes
func LongText() *Type {
	return &Type{
		Sql: func() string {
			return "LONGTEXT"
		},
	}
}

// generates common bigint syntax
// mysql: Signed: -2^63 to 2^63-1, Unsigned: 0 to 2^64-1
// postgres: -2^63 to 2^63-1
func BigInt() *Type {
	return &Type{
		Sql: func() string {
			return "BIGINT"
		},
	}
}

// generates common int syntax
// mysql: Signed: -2^31 to 2^31-1 (2G), Unsigned: 0 to 2^32-1 (4G)
// postgres: -2^31 to 2^31-1 (2G)
func Int() *Type {
	return &Type{
		Sql: func() string {
			return "INT"
		},
	}
}

// generates common smallint syntax
// mysql: Signed: -32768 to 32767, Unsigned: 0 to 65535
func SmallInt() *Type {
	return &Type{
		Sql: func() string {
			return "SMALLINT"
		},
	}
}

// generates 32-bit auto-increment int for postgresql syntax
// postgres: -2^31 to 2^31-1 (2G)
func Serial() *Type {
	return &Type{
		Sql: func() string {
			return "SERIAL"
		},
	}
}

// generates 64-bit auto-increment int for postgresql syntax
// postgres: -2^63 to 2^63-1
func BigSerial() *Type {
	return &Type{
		Sql: func() string {
			return "BIGSERIAL"
		},
	}
}

// generates common numeric type syntax
// p: max number of all digits (both sides)
// s: max number of digits after the decimal point
func Numeric(p int, s int) *Type {
	return &Type{
		Sql: func() string {
			return fmt.Sprintf("NUMERIC(%d, %d)", p, s)
		},
	}
}

// generates float type syntax

// mysql: If only p is specified, p is the binary precision. if p and s are both specified, p is the maximum number of all digits (both sides of the decimal point),
// s is the maximum number of digits after the point. p and s are optional
// -3.402823466E+38 to -1.175494351E-38, 0, and 1.175494351E-38 to 3.402823466E+3

// postgres: p specifies the precision in binary digits
// 1E-307 to 1E+308, 1 ⇐ p ⇐ 53
func Float(p int, s ...int) *Type {
	return &Type{
		Sql: func() string {
			if len(s) == 0 {
				return fmt.Sprintf("FLOAT(%d)", p)
			} else {
				return fmt.Sprintf("FLOAT(%d, %d)", p, s[0])
			}
		},
	}
}

// generates double for mysql
// mysql: -1.797E+308 to -2.2254E-308, 0, and 2.225E-308 to 1.797E+308
// p is the maximum number of all digits (both sides of the decimal point), s is the maximum number of digits after the point, optional
func Double(p int, s int) *Type {
	return &Type{
		Sql: func() string {
			return fmt.Sprintf("DOUBLE(%d, %d)", p, s)
		},
	}
}

// generates double precision for postgres
// postgres: 1E-307 to 1E+308
func DoublePrecision() *Type {
	return &Type{
		Sql: func() string {
			return "DOUBLE PRECISION"
		},
	}
}

// generates Date type for common syntax
// mysql: January 01, 1000 to December 31, 9999
// postgres: January 01, 4713 BC to December 31, 5874897
func Date() *Type {
	return &Type{
		Sql: func() string {
			return "DATE"
		},
	}
}

// generates a date time type in common syntax
// mysql: -838:59:59 to 838:59:59
// postgres: 00:00:00.000000 to 23:59:59.999999
func Time() *Type {
	return &Type{
		Sql: func() string {
			return "TIME"
		},
	}
}

// generates a date time type in mysql syntax
// mysql: January 01, 1000 00:00:00 to December 31, 9999 23:59:59
func DateTime() *Type {
	return &Type{
		Sql: func() string {
			return "DATETIME"
		},
	}
}

// generates a timestamp type in common syntax
// mysql: January 01, 1970 00:00:01 to January 01, 2038 23:59:59
// postgres: January 01, 4713 BC 12:00:00.000000 AM to December 31, 5874897 11:59:59.999999 PM
func Timestamp() *Type {
	return &Type{
		Sql: func() string {
			return "TIMESTAMP"
		},
	}
}

// generates a year type in mysql syntax
// mysql: 4 digit format: 1901 to 2155 and 0; 2 digit format: 70 to 69 (1970 to 2069)
func Year() *Type {
	return &Type{
		Sql: func() string {
			return "YEAR"
		},
	}
}

// generates interval type in postgres syntax
// p is the number of digits in the fractional part of second, 0 ⇐ p ⇐ 6
// postgres: Interval units: microsecond, millisecond, second, minute, hour, day, week, month, year, decade, century and millennium
func Interval(p int) *Type {
	return &Type{
		Sql: func() string {
			return fmt.Sprintf("INTERVAL(%d)", p)
		},
	}
}

// generates bytea type in postgres syntax
// postgres: 231-1 bytes (2G)
func Bytea() *Type {
	return &Type{
		Sql: func() string {
			return "BYTEA"
		},
	}
}

// generates blob type in mysql syntax
// mysql: 65555 bytes
func Blob(n int) *Type {
	return &Type{
		Sql: func() string {
			return fmt.Sprintf("BLOB(%d)", n)
		},
	}
}

// generates blob type in mysql syntax
// mysql: 16,777,215 (2^24 – 1) bytes
func MediumBlob(n int) *Type {
	return &Type{
		Sql: func() string {
			return fmt.Sprintf("MEDIUMBLOB(%d)", n)
		},
	}
}

// generates blob type in mysql syntax
// mysql: 4G (2^32 – 1) bytes
func LongBlob(n int64) *Type {
	return &Type{
		Sql: func() string {
			return fmt.Sprintf("LONGBLOB(%d)", n)
		},
	}
}

// generates money type for postgresql syntax
// postgres: -92,233,720,368,547,758.08 to 92,233,720,368,547,758.07
func Money() *Type {
	return &Type{
		Sql: func() string {
			return "MONEY"
		},
	}
}

// generates common boolean type syntax
// TRUE, FALSE and NULL
func Boolean() *Type {
	return &Type{
		Sql: func() string {
			return "BOOLEAN"
		},
	}
}

// generates uuid type for postgresql syntax
// postgres: 16 bytes
func UUID() *Type {
	return &Type{
		Sql: func() string {
			return "UUID"
		},
	}
}

// generates enum type for mysql syntax
// mysql: Maximum 65,535 distinct values
func Enum(vals ...string) *Type {
	return &Type{
		Sql: func() string {
			for k, v := range vals {
				vals[k] = fmt.Sprintf("'%s'", v)
			}
			return fmt.Sprintf("ENUM(%s)", strings.Join(vals, ", "))
		},
	}
}
