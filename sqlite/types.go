package sqlite

import "github.com/aacanakin/qbit"

// generates a date type in sqlite syntax
func Date() *qbit.Type {
	return &qbit.Type{
		Sql: func() string {
			return "DATE"
		},
	}
}

// generates a date time type in sqlite syntax
func DateTime() *qbit.Type {
	return &qbit.Type{
		Sql: func() string {
			return "DATETIME"
		},
	}
}

// generates a date time type in sqlite syntax
func Time() *qbit.Type {
	return &qbit.Type{
		Sql: func() string {
			return "TIME"
		},
	}
}
