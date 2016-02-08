package qbit

import (
	"errors"
	"fmt"
	"github.com/fatih/structs"
	"github.com/serenize/snaker"
	"reflect"
	"strings"
)

const TAG = "qbit"
const DEFAULT_FLOAT_P = 53

func NewMapper() *Mapper {
	return &Mapper{}
}

type Mapper struct {
}

func (m *Mapper) extractValue(value string) string {
	startIndex := strings.Index(value, "(")
	endIndex := strings.Index(value, ")")
	return value[startIndex+1 : endIndex]
}

func (m *Mapper) convertConstraints(c string) []Constraint {

	constraints := []Constraint{}
	rawConstraints := strings.Split(c, ",")

	var constraint Constraint
	for _, v := range rawConstraints {

		if v == "null" {
			constraint = Null()
		} else if v == "notnull" {
			constraint = NotNull()
		} else if v == "primary_key" {
			constraint = PrimaryKey()
		} else if v == "unique" {
			constraint = Unique()
		} else if v == "key" {
			constraint = Key()
		} else if v == "index" {
			constraint = Index()
		} else if strings.Contains(v, "default") {
			constraint = Default(m.extractValue(v))
		} else {
			constraint = Constraint{}
		}

		fmt.Println("Matched constraint: ", constraint.Name)

		constraints = append(constraints, constraint)
	}

	return constraints
}

func (m *Mapper) convertType(t string) *Type {
	return &Type{}
}

func (m *Mapper) convertDefaultType(kind reflect.Kind) *Type {

	if kind.String() == "string" {
		return VarChar(255)
	} else if kind.String() == "int" {
		return Int()
	} else if kind.String() == "int64" {
		return BigInt()
	} else if kind.String() == "float32" || kind.String() == "float64" {
		return Float(DEFAULT_FLOAT_P)
	} else if kind.String() == "bool" {
		return Boolean()
	}

	// default is string
	return VarChar(255)
}

func (m *Mapper) convertColumns(fields []*structs.Field) ([]Column, error) {

	cols := []Column{}

	var col Column
	var tag string

	for _, f := range fields {
		col = Column{}
		col.Name = snaker.CamelToSnake(f.Name())

		fmt.Printf("field name: %s\n", snaker.CamelToSnake(f.Name()))
		fmt.Printf("field tag: %s\n", f.Tag(TAG))
		fmt.Printf("field type: %s\n", f.Kind())

		// clean trailing spaces of tag
		tag = strings.Replace(f.Tag(TAG), " ", "", 1)

		fmt.Println("new tag: ", tag)
		for _, t := range strings.Split(tag, ";") {

			fmt.Println(t)
			col.T = m.convertDefaultType(f.Kind())

			if t != "" {
				tTypes := strings.Split(t, ":")
				fmt.Println(tTypes)

				if len(tTypes) != 2 {
					return nil, errors.New("Invalid tag types. Please make sure that tag type values doesn't contain ':' character")
				}

				if tTypes[0] == "constraints" || tTypes[0] == "constraint" {
					col.Constraints = m.convertConstraints(tTypes[1])
				} else if tTypes[0] == "type" {
					col.T = m.convertType(tTypes[1])
				}
			}
		}

		fmt.Println()
		cols = append(cols, col)
	}

	return []Column{}, nil
}

func (m *Mapper) Convert(model interface{}) (*Table, error) {

	//	name := strings.ToLower(structs.Name(model))

	//	cols, err := m.convertColumns(structs.Fields(model))
	//	if err != nil {
	//		return nil, err
	//	}

	fmt.Printf("model name: %T", model)

	//	val := reflect.ValueOf(model)
	//	for i := 0; i < val.NumField(); i++ {
	//
	//		valueField := val.Field(i)
	//		typeField := val.Field(i).Type()
	//		tag := typeField.Field(i).Tag
	//
	//		fmt.Printf("Field Name: %s, Type: %s, Tag(qbit): %s\n",
	//			valueField.String(),
	//			typeField.String(),
	//			tag.Get("qbit"),
	//		)
	//
	//		fmt.Println()
	//	}
	//
	return nil, nil

	//	return NewTable(name, cols, []Constraint{}), nil
}
