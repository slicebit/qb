package qb

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/serenize/snaker"
	"strings"
)

const tagPrefix = "qb"

// Mapper creates a new mapper struct and returns it as a mapper pointer
func Mapper() MapperElem {
	return MapperElem{}
}

// MapperElem is the generic struct for struct to table mapping
type MapperElem struct{}

func (m MapperElem) extractValue(value string) string {
	hasParams := strings.Contains(value, "(") && strings.Contains(value, ")")

	if hasParams {
		startIndex := strings.Index(value, "(")
		endIndex := strings.Index(value, ")")
		return value[startIndex+1 : endIndex]
	}

	return ""
}

// ToMap converts a model struct to a map. Uninitialized fields are optionally ignored if includeZeroFields is true.
// Fields are renamed using qb conventions. db tag overrides column names
func (m MapperElem) ToMap(model interface{}, includeZeroFields bool) map[string]interface{} {
	fields := structs.Fields(model)
	kv := map[string]interface{}{}
	for _, f := range fields {
		if strings.Contains(f.Tag("qb"), "-") {
			continue
		}

		if f.Name() == "qb.CompositeIndex" || f.Name() == "CompositeIndex" {
			continue
		}

		isZero := f.IsZero()
		if isZero && !includeZeroFields {
			continue
		}

		var val interface{}
		if f.IsZero() {
			val = nil
		} else {
			val = f.Value()
		}

		dbTag := ParseDBTag(f.Tag("db"))
		if dbTag != "" {
			kv[dbTag] = val
		} else {
			kv[snaker.CamelToSnake(f.Name())] = val
		}
	}
	return kv
}

// ModelName returns the sql table name of model
func (m MapperElem) ModelName(model interface{}) string {
	return snaker.CamelToSnake(structs.Name(model))
}

// ToType returns the type mapping of column.
// If tagType is, then colType would automatically be resolved.
// If tagType is not "", then automatic type resolving would be overridden by tagType
func (m *MapperElem) ToType(colType string, tagType string) TypeElem {
	// convert tagType
	if tagType != "" {
		tagType = strings.ToUpper(tagType)
		return Type(tagType)
	}
	// convert default type
	switch colType {
	case "string":
		return Varchar().Size(255)
	case "int":
		return Int()
	case "int8":
		return TinyInt()
	case "int16":
		return SmallInt()
	case "int32":
		return Int()
	case "int64":
		return BigInt()
	case "uint":
		return Int().Unsigned()
	case "uint8":
		return TinyInt().Unsigned()
	case "uint16":
		return SmallInt().Unsigned()
	case "uint32":
		return Int().Unsigned()
	case "uint64":
		return BigInt().Unsigned()
	case "float32":
		return Float() // TODO: Not sure
	case "float64":
		return Float() // TODO: Not sure
	case "bool":
		return Boolean()
	case "time.Time":
		return Timestamp()
	case "*time.Time":
		return Timestamp()
	default:
		return Varchar().Size(255)
	}
}

// ToTable parses struct and converts it to a new table
func (m *MapperElem) ToTable(model interface{}) (TableElem, error) {
	modelName := m.ModelName(model)

	colClauses := []TableClause{}
	constraintClauses := []TableClause{}

	for _, f := range structs.Fields(model) {

		tag, err := ParseTag(f.Tag(tagPrefix))
		if err != nil {
			return TableElem{}, err
		}

		if tag.Ignore {
			continue
		}

		dbTag := ParseDBTag(f.Tag("db"))
		var colName string
		if dbTag != "" {
			colName = dbTag
		} else {
			colName = snaker.CamelToSnake(f.Name())
		}

		colType := m.ToType(fmt.Sprintf("%T", f.Value()), tag.Type)

		// convert tag into constraints
		for _, v := range tag.Constraints {
			if v == "null" {
				colType = colType.Null()
			} else if v == "notnull" || v == "not_null" {
				colType = colType.NotNull()
			} else if v == "unique" {
				colType = colType.Unique()
			} else if v == "auto_increment" || v == "autoincrement" {
				colType = colType.AutoIncrement()
			} else if strings.Contains(v, "default") {
				colType = colType.Default(m.extractValue(v))
			} else if strings.Contains(v, "primary_key") {
				// TODO: Possible performance issue, fix this when possible, maybe table.AddPrimary option should be thought
				pkDefined := false
				for i, tc := range constraintClauses {
					switch tc.(type) {
					case PrimaryKeyConstraint:
						pk := tc.(PrimaryKeyConstraint)
						pk.Columns = append(pk.Columns, colName)
						constraintClauses[i] = pk
						pkDefined = true
						break
					}
				}
				if pkDefined {
					continue
				}
				constraintClauses = append(constraintClauses, PrimaryKey(colName))
			} else if strings.Contains(v, "ref") && strings.Contains(v, "(") && strings.Contains(v, ")") {
				// TODO: Possible performance issue, fix this when possible, maybe table.AddRef option should be thought
				fkp := strings.Split(m.extractValue(v), ".")
				fkDefined := false
				for i, tc := range constraintClauses {
					switch tc.(type) {
					case ForeignKeyConstraints:
						fk := tc.(ForeignKeyConstraints)
						fk = fk.Ref(colName, fkp[0], fkp[1])
						constraintClauses[i] = fk
						fkDefined = true
						break
					}
				}

				if fkDefined {
					continue
				}

				constraintClauses = append(constraintClauses, ForeignKey().Ref(colName, fkp[0], fkp[1]))
			} else if strings.Contains(v, "index") {
				if strings.Contains(f.Name(), "CompositeIndex") {
					is := strings.Split(v, ":")
					cols := strings.Split(is[1], ",")
					constraintClauses = append(constraintClauses, Index(modelName, cols...))
				} else {
					constraintClauses = append(constraintClauses, Index(modelName, colName))
				}
			} else {
				return TableElem{}, fmt.Errorf("Invalid constraint: %s", v)
			}
		}

		if strings.Contains(f.Name(), "CompositeIndex") {
			continue
		}

		colClauses = append(colClauses, Column(colName, colType))
	}

	return Table(modelName, append(colClauses, constraintClauses...)...), nil
}
