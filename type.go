package qb

import (
	"fmt"
	"strings"
)

// Char creates char type
func Char() TypeElem {
	return Type("CHAR")
}

// Varchar creates varchar type
func Varchar() TypeElem {
	return Type("VARCHAR").Size(255)
}

// Text creates text type
func Text() TypeElem {
	return Type("TEXT")
}

// Int creates int type
func Int() TypeElem {
	return Type("INT")
}

// SmallInt creates smallint type
func SmallInt() TypeElem {
	return Type("SMALLINT")
}

// BigInt creates bigint type
func BigInt() TypeElem {
	return Type("BIGINT")
}

// Numeric creates a numeric type
func Numeric() TypeElem {
	return Type("NUMERIC")
}

// Decimal creates a decimal type
func Decimal() TypeElem {
	return Type("DECIMAL")
}

// Float creates float type
func Float() TypeElem {
	return Type("FLOAT")
}

// Boolean creates boolean type
func Boolean() TypeElem {
	return Type("BOOLEAN")
}

// Timestamp creates timestamp type
func Timestamp() TypeElem {
	return Type("TIMESTAMP")
}

const defaultTypeSize = -1

// Type returns a new TypeElem while defining columns in table
func Type(name string) TypeElem {
	return TypeElem{
		Name:      name,
		size:      defaultTypeSize,
		precision: []int{},
	}
}

// TypeElem is the struct for defining column types
type TypeElem struct {
	Name        string
	constraints []ConstraintElem
	size        int
	precision   []int
	unique      bool
}

// String returns the clause as string
func (t TypeElem) String() string {
	constraintNames := []string{}
	for _, c := range t.constraints {
		constraintNames = append(constraintNames, c.String())
	}

	sizeSpecified := false
	if t.size != defaultTypeSize {
		sizeSpecified = true
	}

	if sizeSpecified {
		return strings.Trim(fmt.Sprintf("%s(%d) %s", t.Name, t.size, strings.Join(constraintNames, " ")), " ")
	}

	precisionSpecified := false
	if len(t.precision) > 0 {
		precisionSpecified = true
	}

	if precisionSpecified {
		precision := []string{}
		for _, p := range t.precision {
			precision = append(precision, fmt.Sprintf("%v", p))
		}
		return strings.Trim(fmt.Sprintf("%s(%s) %s", t.Name, strings.Join(precision, ", "), strings.Join(constraintNames, " ")), " ")
	}

	return strings.Trim(fmt.Sprintf("%s %s", t.Name, strings.Join(constraintNames, " ")), " ")
}

// Size adds size constraint to column type
func (t TypeElem) Size(size int) TypeElem {
	t.size = size
	return t
}

// Precision sets the precision of column type
// Note: Use it in Float, Decimal and Numeric types
func (t TypeElem) Precision(p int, s int) TypeElem {
	t.precision = []int{p, s}
	return t
}

// Default adds a default constraint to column type
func (t TypeElem) Default(def interface{}) TypeElem {
	t.constraints = append(t.constraints, Default(def))
	return t
}

// Null adds null constraint to column type
func (t TypeElem) Null() TypeElem {
	t.constraints = append(t.constraints, Null())
	return t
}

// NotNull adds not null constraint to column type
func (t TypeElem) NotNull() TypeElem {
	t.constraints = append(t.constraints, NotNull())
	return t
}

// Unique adds a unique constraint to column type
func (t TypeElem) Unique() TypeElem {
	t.constraints = append(t.constraints, Unique())
	t.unique = true
	return t
}

// Constraint adds a custom constraint to column type
func (t TypeElem) Constraint(name string) TypeElem {
	t.constraints = append(t.constraints, Constraint(name))
	return t
}
