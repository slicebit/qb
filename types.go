package qbit

import "fmt"

type Type struct {
	Sql func() string
}

// generates common smallint syntax
func SmallInt() *Type {
	return &Type{
		Sql: func() string {
			return "SMALLINT"
		},
	}
}

// generates common int syntax
func Int() *Type {
	return &Type{
		Sql: func() string {
			return "INT"
		},
	}
}

// generates common bigint syntax
func BigInt() *Type {
	return &Type{
		Sql: func() string {
			return "BIGINT"
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

// generates common char type syntax
func Char(size int) *Type {
	return &Type{
		Sql: func() string {
			return fmt.Sprintf("CHAR(%d)", size)
		},
	}
}

// generates common varchar type syntax
func VarChar(size int) *Type {
	return &Type{
		Sql: func() string {
			return fmt.Sprintf("VARCHAR(%d)", size)
		},
	}
}

// generates common text type syntax
func Text() *Type {
	return &Type{
		Sql: func() string {
			return "TEXT"
		},
	}
}

// generates common boolean type syntax
func Boolean() *Type {
	return &Type{
		Sql: func() string {
			return "BOOLEAN"
		},
	}
}