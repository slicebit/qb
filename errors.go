package qb

// ErrorCode discriminates the types of errors that qb wraps, mainly the
// constraint errors
// The different kind of errors are based on the python dbapi errors
// (https://www.python.org/dev/peps/pep-0249/#exceptions)
type ErrorCode int

// Bit 8 and 9 are flags to separate interface errors from database errors
const (
	// ErrAny is for errors that could not be categorized by the dialect
	ErrAny ErrorCode = 0
	// ErrInterface is a bit mask for errors that are related to the database
	// interface rather than the database itself.
	ErrInterface ErrorCode = 1 << 8
	// ErrDatabase is a bit mask for errors that are related to the database.
	ErrDatabase ErrorCode = 1 << 9
)

// Database error codes are in bits 5 to 7, leaving bits 0 to 4 for detailed
// codes in a later version
const (
	// ErrData is for errors that are due to problems with the processed data
	// like division by zero, numeric value out of range, etc.
	ErrData ErrorCode = ErrDatabase | (iota + 1<<5)
	// ErrOperational is for errors that are related to the database's
	// operation and not necessarily under the control of the programmer, e.g.
	// an unexpected disconnect occurs, the data source name is not found, a
	// transaction could not be processed, a memory allocation error occurred
	// during processing, etc.
	ErrOperational
	// ErrIntegrity is when the relational integrity of the database is
	// affected, e.g. a foreign key check fails
	ErrIntegrity
	// ErrInternal is when the database encounters an internal error, e.g. the
	// cursor is not valid anymore, the transaction is out of sync, etc.
	ErrInternal
	// ErrProgramming is for programming errors, e.g. table not found or
	// already exists, syntax error in the SQL statement, wrong number of
	// parameters specified, etc.
	ErrProgramming
	// ErrNotSupported is in case a method or database API was used which
	// is not supported by the database, e.g. requesting a .rollback() on a
	// connection that does not support transaction or has transactions turned
	// off.
	ErrNotSupported
)

// IsInterfaceError returns true if the error is a Interface error
func (err ErrorCode) IsInterfaceError() bool {
	return err&ErrInterface != 0
}

// IsDatabaseError returns true if the error is a Database error
func (err ErrorCode) IsDatabaseError() bool {
	return err&ErrDatabase != 0
}

// Error wraps driver errors. It helps handling constraint error in
// a generic way, while still giving access to the original error
type Error struct {
	Code       ErrorCode
	Orig       error // The native error from the driver
	Table      string
	Column     string
	Constraint string
}

func (err Error) Unwrap() error {
	return err.Orig
}

func (err Error) Error() string {
	switch err.Code {
	case ErrAny:
		return "Uncategorized error: " + err.Orig.Error()
	case ErrInterface:
		return "Interface error: " + err.Orig.Error()
	case ErrDatabase:
		return "Database error: " + err.Orig.Error()
	case ErrData:
		return "Database data error: " + err.Orig.Error()
	case ErrOperational:
		return "Database operational error: " + err.Orig.Error()
	case ErrIntegrity:
		return "Database integrity error: " + err.Orig.Error()
	case ErrInternal:
		return "Database internal error: " + err.Orig.Error()
	case ErrProgramming:
		return "Database programming error: " + err.Orig.Error()
	default:
		return err.Orig.Error()
	}
}
