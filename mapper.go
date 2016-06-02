package qb

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/serenize/snaker"
	"strings"
)

const tagPrefix = "qb"

// NewMapper creates a new mapper struct and returns it as a mapper pointer
func Mapper(adapter Adapter) MapperElem {
	return MapperElem{adapter: adapter}
}

// Mapper is the generic struct for struct to table mapping
type MapperElem struct {
	//builder *Builder
	adapter Adapter
}

func (m MapperElem) extractValue(value string) string {
	hasParams := strings.Contains(value, "(") && strings.Contains(value, ")")

	if hasParams {
		startIndex := strings.Index(value, "(")
		endIndex := strings.Index(value, ")")
		return value[startIndex+1 : endIndex]
	}

	return ""
}

// ToRawMap converts a model struct to map without changing the field names.
func (m MapperElem) ToRawMap(model interface{}) map[string]interface{} {
	return structs.Map(model)
}

// ToMap converts a model struct to a map. Uninitialized fields are ignored.
// Fields are renamed using qb conventions
func (m MapperElem) ToMap(model interface{}) map[string]interface{} {
	fields := structs.Fields(model)
	kv := map[string]interface{}{}
	for _, f := range fields {
		if f.IsZero() {
			continue
		}

		kv[m.ColName(f.Name())] = f.Value()
	}
	return kv
}

// ModelName returns the sql table name of model
func (m MapperElem) ModelName(model interface{}) string {
	return snaker.CamelToSnake(structs.Name(model))
}

// ColName returns the sql column name of model
func (m MapperElem) ColName(col string) string {
	return snaker.CamelToSnake(col)
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
		return SmallInt()
	case "int16":
		return SmallInt()
	case "int32":
		return Int()
	case "int64":
		return BigInt()
	case "uint":
		if m.adapter.SupportsUnsigned() {
			return Type("INT UNSIGNED")
		}
		return BigInt()
	case "uint8":
		if m.adapter.SupportsUnsigned() {
			return Type("TINYINT UNSIGNED")
		}
		return SmallInt()
	case "uint16":
		if m.adapter.SupportsUnsigned() {
			return Type("SMALLINT UNSIGNED")
		}
		return Int()
	case "uint32":
		if m.adapter.SupportsUnsigned() {
			return Type("INT UNSIGNED")
		}
		return BigInt()
	case "uint64":
		if m.adapter.SupportsUnsigned() {
			return Type("BIGINT UNSIGNED")
		}
		return BigInt()
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

	colClauses := []Clause{}
	tableClauses := []Clause{}

	for _, f := range structs.Fields(model) {

		tag, err := ParseTag(f.Tag(tagPrefix))
		if err != nil {
			return TableElem{}, err
		}

		if tag.Ignore {
			continue
		}

		colName := m.ColName(f.Name())
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
				c := m.adapter.AutoIncrement()

				// it doesn't support auto increment
				if c == "" {
					continue
				} else {
					colType = colType.Constraint(c)
				}
			} else if strings.Contains(v, "default") {
				colType = colType.Default(m.extractValue(v))
			} else if strings.Contains(v, "primary_key") {
				// TODO: Possible performance issue, fix this when possible, maybe table.AddPrimary option should be thought
				pkDefined := false
				for i, tc := range tableClauses {
					switch tc.(type) {
					case PrimaryKeyConstraint:
						pk := tc.(PrimaryKeyConstraint)
						pk.Columns = append(pk.Columns, colName)
						tableClauses[i] = pk
						pkDefined = true
						break
					}
				}
				if pkDefined {
					continue
				}
				tableClauses = append(tableClauses, PrimaryKey(colName))
			} else if strings.Contains(v, "ref") && strings.Contains(v, "(") && strings.Contains(v, ")") {
				// TODO: Possible performance issue, fix this when possible, maybe table.AddRef option should be thought
				fkp := strings.Split(m.extractValue(v), ".")
				fkDefined := false
				for i, tc := range tableClauses {
					switch tc.(type) {
					case ForeignKeyConstraints:
						fk := tc.(ForeignKeyConstraints)
						fk = fk.Ref(colName, fkp[0], fkp[1])
						tableClauses[i] = fk
						fkDefined = true
						break
					}
				}

				if fkDefined {
					continue
				}

				tableClauses = append(tableClauses, ForeignKey().Ref(colName, fkp[0], fkp[1]))
			} else if strings.Contains(v, "index") {
				if strings.Contains(f.Name(), "CompositeIndex") {
					is := strings.Split(v, ":")
					cols := strings.Split(is[1], ",")
					tableClauses = append(tableClauses, Index(modelName, cols...))
				} else {
					tableClauses = append(tableClauses, Index(modelName, colName))
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

	return Table(modelName, append(colClauses, tableClauses...)...), nil
}
