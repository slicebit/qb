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

// TinyInt creates tinyint type
func TinyInt() TypeElem {
	return Type("TINYINT")
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

// UUID creates a UUID type
func UUID() TypeElem {
	return Type("UUID")
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
	Name      string
	size      int
	precision []int
	unsigned  bool
}

// DefaultCompileType is a default implementation for Dialect.CompileType
func DefaultCompileType(t TypeElem, supportsUnsigned bool) string {
	name := t.Name

	if t.unsigned && !supportsUnsigned {
		// use a bigger int type so the unsigned values can fit in
		switch name {
		case "TINYINT":
			name = "SMALLINT"
		case "SMALLINT":
			name = "INT"
		case "INT":
			name = "BIGINT"
		}
	}

	sizeSpecified := false
	precisionSpecified := false
	if t.size != defaultTypeSize {
		sizeSpecified = true
	} else if len(t.precision) > 0 {
		precisionSpecified = true
	}

	if sizeSpecified {
		name = fmt.Sprintf("%s(%d)", name, t.size)
	} else if precisionSpecified {
		precision := []string{}
		for _, p := range t.precision {
			precision = append(precision, fmt.Sprintf("%v", p))
		}
		name = fmt.Sprintf("%s(%s)", name, strings.Join(precision, ", "))
	}

	if t.unsigned && supportsUnsigned {
		name = fmt.Sprintf("%s UNSIGNED", name)
	}
	return name
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

// Unsigned change the column type to 'unsigned'
// Note: Use it in Float, Decimal and Numeric types
func (t TypeElem) Unsigned() TypeElem {
	t.unsigned = true
	return t
}

// Signed change the column type to 'signed'
// Note: Use it in Float, Decimal and Numeric types
func (t TypeElem) Signed() TypeElem {
	t.unsigned = false
	return t
}
