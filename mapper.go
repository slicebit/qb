package qb

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/serenize/snaker"
	"strings"
)

const tagPrefix = "qb"

// NewMapper creates a new mapper struct and returns it as a mapper pointer
func NewMapper(driver string) *Mapper {
	return &Mapper{
		driver: driver,
	}
}

// Mapper is the generic struct for struct to table mapping
type Mapper struct {
	driver string
}

func (m *Mapper) extractValue(value string) string {

	hasParams := strings.Contains(value, "(") && strings.Contains(value, ")")

	if hasParams {
		startIndex := strings.Index(value, "(")
		endIndex := strings.Index(value, ")")
		return value[startIndex+1 : endIndex]
	}

	return ""
}

// ToRawMap converts a model struct to map without changing the field names.
func (m *Mapper) ToRawMap(model interface{}) map[string]interface{} {
	return structs.Map(model)
}

// ToMap converts a model struct to a map. Uninitialized fields are ignored.
// Fields are renamed using qb conventions
func (m *Mapper) ToMap(model interface{}) map[string]interface{} {

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
func (m *Mapper) ModelName(model interface{}) string {
	return snaker.CamelToSnake(structs.Name(model))
}

// ColName returns the sql column name of model
func (m *Mapper) ColName(col string) string {
	return snaker.CamelToSnake(col)
}

// ToType returns the type mapping of column.
// If tagType is, then colType would automatically be resolved.
// If tagType is not "", then automatic type resolving would be overridden by tagType
func (m *Mapper) ToType(colType string, tagType string) *Type {

	// convert tagType
	if tagType != "" {
		tagType = strings.ToUpper(tagType)
		return &Type{tagType}
	}

	// convert default type
	switch colType {
	case "string":
		return &Type{"VARCHAR(255)"}
	case "int":
		return &Type{"INT"}
	case "int64":
		return &Type{"BIGINT"}
	case "float32":
		return &Type{"FLOAT"}
	case "float64":
		return &Type{"FLOAT"}
	case "bool":
		return &Type{"BOOLEAN"}
	case "time.Time":
		return &Type{"TIMESTAMP"}
	case "*time.Time":
		return &Type{"TIMESTAMP"}
	default:
		return &Type{"VARCHAR"}
	}
}

// ToTable parses struct and converts it to a new table
func (m *Mapper) ToTable(model interface{}) (*Table, error) {

	modelName := m.ModelName(model)

	table := NewTable(m.driver, modelName, []Column{})
	adapter := NewAdapter(m.driver)

	//fmt.Printf("model name: %s\n\n", modelName)

	var col Column
	var rawTag string

	for _, f := range structs.Fields(model) {

		colName := m.ColName(f.Name())
		colType := fmt.Sprintf("%T", f.Value())

		rawTag = f.Tag(tagPrefix)

		constraints := []Constraint{}
		//fmt.Println()
		//fmt.Printf("field name: %s\n", colName)
		//fmt.Printf("field raw tag: %s\n", rawTag)
		//fmt.Printf("field type name: %T\n", f.Value())
		//fmt.Printf("field constraints: %v\n", constraints)
		//fmt.Println()

		// clean trailing spaces of tag
		rawTag = strings.Replace(f.Tag(tagPrefix), " ", "", -1)

		// parse tag
		tag, err := ParseTag(rawTag)
		if err != nil {
			return nil, err
		}

		if tag.Ignore {
			continue
		}

		// convert tag into constraints
		var constraint Constraint
		for _, v := range tag.Constraints {
			if v == "null" {
				constraint = Null()
			} else if v == "notnull" || v == "not_null" {
				constraint = NotNull()
			} else if v == "unique" {
				constraint = Constraint{
					Name: "UNIQUE",
				}
			} else if v == "auto_increment" || v == "autoincrement" {
				if m.driver == "mysql" {
					constraint = Constraint{
						Name: "AUTO_INCREMENT",
					}
				} else if m.driver == "sqlite3" {
					constraint = Constraint{
						Name: "AUTOINCREMENT",
					}
				} else {
					continue
				}
			} else if strings.Contains(v, "default") {
				constraint = Default(m.extractValue(v))
			} else if strings.Contains(v, "primary_key") {
				if adapter.SupportsInlinePrimaryKey() {
					constraint = Constraint{
						Name: "PRIMARY KEY",
					}
				} else {
					table.AddPrimary(colName)
					continue
				}
			} else if strings.Contains(v, "ref") && strings.Contains(v, "(") && strings.Contains(v, ")") {
				tc := strings.Split(m.extractValue(v), ".")
				table.AddRef(colName, tc[0], tc[1])
				continue
			} else if strings.Contains(v, "index") {
				if strings.Contains(f.Name(), "CompositeIndex") {
					is := strings.Split(v, ":")
					cols := strings.Split(is[1], ",")
					table.AddIndex(cols...)
				} else {
					table.AddIndex(colName)
				}
				continue
			} else {
				return nil, fmt.Errorf("Invalid constraint: %s", v)
			}
			constraints = append(constraints, constraint)
		}

		//fmt.Printf("field tag.Type: %s\n", tag.Type)
		//fmt.Printf("field tag.Constraints: %v\n", tag.Constraints)

		if strings.Contains(f.Name(), "CompositeIndex") {
			continue
		}

		col = Column{
			Name:        colName,
			Constraints: constraints,
			Type:        m.ToType(colType, tag.Type),
		}

		table.AddColumn(col)

		//fmt.Println()
	}

	return table, nil
}
