package qbit

import (
	//	"errors"
	"fmt"
	"github.com/fatih/structs"
	"github.com/serenize/snaker"
	//	"reflect"
	"strings"
	//	"errors"
	"errors"
)

const TAG = "qbit"

func NewMapper(driver string) *Mapper {
	return &Mapper{
		driver: driver,
	}
}

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

func (m *Mapper) ConvertType(colType string, tagType string) *Type {

	// convert tagType
	if tagType != "" {
		tagType = strings.ToUpper(tagType)
		return &Type{func() string { return tagType }}
	}

	// convert default type
	switch colType {
	case "string":
		return VarChar()
	case "int":
		return Int()
	case "int64":
		return BigInt()
	case "float32":
		return Float()
	case "float64":
		return Float()
	case "bool":
		return Boolean()
	case "uuid.UUID":
		if m.driver == "postgres" {
			return UUID()
		}
		return VarChar(36)
	case "time.Time":
		return Timestamp()
	case "*time.Time":
		return Timestamp()
	default:
		return VarChar()
	}
}

func (m *Mapper) convertConstraints(rawConstraints []string) ([]Constraint, error) {

	constraints := []Constraint{}

	var constraint Constraint
	for _, v := range rawConstraints {

		if v == "null" {
			constraint = Null()
		} else if v == "notnull" {
			constraint = NotNull()
		} else if v == "unique" {
			constraint = Constraint{
				Name: "UNIQUE",
			}
		} else if v == "index" {
			constraint = Constraint{
				Name: "INDEX",
			}
		} else if strings.Contains(v, "default") {
			constraint = Default(m.extractValue(v))
		} else {
			return nil, errors.New(fmt.Sprintf("Invalid constraint: %s", v))
		}

		constraints = append(constraints, constraint)
	}

	return constraints, nil
}

func (m *Mapper) Convert(model interface{}) (*Table, error) {

	modelName := snaker.CamelToSnake(structs.Name(model))

	table := &Table{
		name:        modelName,
		columns:     []Column{},
		constraints: []Constraint{},
		builder:     NewBuilder(),
	}

	fmt.Printf("model name: %s\n\n", modelName)

	var col Column
	var rawTag string

	for _, f := range structs.Fields(model) {

		colName := snaker.CamelToSnake(f.Name())
		colType := fmt.Sprintf("%T", f.Value())

		rawTag = f.Tag(TAG)

		constraints := []Constraint{}
		fmt.Printf("field name: %s\n", colName)
		fmt.Printf("field raw tag: %s\n", rawTag)
		fmt.Printf("field type name: %T\n", f.Value())
		fmt.Printf("field constraints: %v\n", constraints)

		if colType == "qbit.PrimaryKey" {
			table.AddConstraint(Constraint{
				Name: fmt.Sprintf("PRIMARY KEY (%s)", rawTag),
			})
		} else if colType == "qbit.ForeignKey" {
			rawTag = strings.Replace(rawTag, "references", "REFERENCES", 1)
			table.AddConstraint(Constraint{
				Name: fmt.Sprintf("FOREIGN KEY %s", rawTag),
			})
		} else if colType == "qbit.Index" {
			rawTagPieces := strings.Split(rawTag, ";")
			for _, v := range rawTagPieces {
				table.AddConstraint(Constraint{
					Name: fmt.Sprintf("INDEX(%s)", v),
				})
			}
		} else if colType == "qbit.CompositeUnique" {
			rawTagPieces := strings.Split(rawTag, ";")
			for _, v := range rawTagPieces {
				table.AddConstraint(Constraint{
					Name: fmt.Sprintf("UNIQUE(%s)", v),
				})
			}
		} else {

			// clean trailing spaces of tag
			rawTag = strings.Replace(f.Tag(TAG), " ", "", 1)

			// parse tag
			tag, err := ParseTag(rawTag)
			if err != nil {
				return nil, err
			}

			// convert tag into constraints
			constraints, err = m.convertConstraints(tag.Constraints)
			if err != nil {
				return nil, err
			}

			fmt.Printf("field tag.Type: %s\n", tag.Type)
			fmt.Printf("field tag.Constraints: %v\n", tag.Constraints)

			col = Column{
				Name:        colName,
				Constraints: constraints,
				Type:        m.ConvertType(colType, tag.Type), // TODO: map type
			}

			table.AddColumn(col)
		}

		fmt.Println()
	}

	return table, nil
}
