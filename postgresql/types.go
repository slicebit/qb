package postgresql

import (
	"fmt"
	"github.com/aacanakin/qbit"
)

// generates 32-bit auto-increment int for postgresql syntax
func Serial() *qbit.Type {
	return &qbit.Type{
		Sql: func() string {
			return "SERIAL"
		},
	}
}

// generates 64-bit auto-increment int for postgresql syntax
func BigSerial() *qbit.Type {
	return &qbit.Type{
		Sql: func() string {
			return "BIGSERIAL"
		},
	}
}

// generates Real type for postgresql syntax
func Real() *qbit.Type {
	return &qbit.Type{
		Sql: func() string {
			return "REAL"
		},
	}
}

// generates Time type for postgresql syntax
// p: number of digits in the fractional part of second
func Time(p int) *qbit.Type {
	return &qbit.Type{
		Sql: func() string {
			return fmt.Sprintf("TIME(%d)", p)
		},
	}
}

// generates Time type for postgresql syntax
// p: number of digits in the fractional part of second
func Timestamp(p int, withTimezone bool) *qbit.Type {

	var tz string
	if withTimezone {
		tz = "WITH TIMEZONE"
	} else {
		tz = "WITHOUT TIMEZONE"
	}

	return &qbit.Type{
		Sql: func() string {
			return fmt.Sprintf("TIMESTAMP(%d) %s", p, tz)
		},
	}
}

// generates interval type for postgresql syntax
// p: number of digits in the fractional part of second
func Interval(p int) *qbit.Type {
	return &qbit.Type{
		Sql: func() string {
			return fmt.Sprintf("INTERVAL(%d)", p)
		},
	}
}

// generates bytea type for postgresql syntax
func Bytea() *qbit.Type {
	return &qbit.Type{
		Sql: func() string {
			return "BYTEA"
		},
	}
}

// generates money type for postgresql syntax
func Money() *qbit.Type {
	return &qbit.Type{
		Sql: func() string {
			return "MONEY"
		},
	}
}

// generates uuid type for postgresql syntax
func UUID() *qbit.Type {
	return &qbit.Type{
		Sql: func() string {
			return "UUID"
		},
	}
}
