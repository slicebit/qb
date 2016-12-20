package qb

type Error interface {
	error
	Orig() error
	Stmt() *Stmt
}

func NewQbError(err error, stmt *Stmt) QbError {
	return QbError{err, stmt}
}

type QbError struct {
	orig      error
	Statement *Stmt
}

func (e QbError) Error() string {
	return e.orig.Error()
}

func (e QbError) Stmt() *Stmt {
	return e.Statement
}

func (e QbError) Orig() error {
	return e.orig
}

// IntegrityError is an error that can be returned by a qb.dialect
// if a constraint is violated
type IntegrityError struct {
	QbError
	Constraint string
}

// Error implements the error interface
func (e IntegrityError) Error() string {
	return "constraint error: " + e.Constraint
}
