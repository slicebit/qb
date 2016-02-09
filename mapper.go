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
const DEFAULT_FLOAT_P = 53

func NewMapper(driver string) *Mapper {

	var typeMapper TypeMapper
	switch driver {
	case "mysql":
		typeMapper = MysqlTypeMapper{}
		break
	case "postgres":
		typeMapper = PostgresTypeMapper{}
		break
	case "sqlite":
		typeMapper = SqliteTypeMapper{}
		break
	}

	return &Mapper{
		driver:     driver,
		typeMapper: typeMapper,
	}
}

type Mapper struct {
	driver     string
	typeMapper TypeMapper
}

func (m *Mapper) extractValue(value string) string {
	startIndex := strings.Index(value, "(")
	endIndex := strings.Index(value, ")")
	return value[startIndex+1 : endIndex]
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
			constraint = Unique()
		} else if v == "key" {
			constraint = Key()
		} else if v == "index" {
			constraint = Index()
		} else if strings.Contains(v, "default") {
			constraint = Default(m.extractValue(v))
		} else {
			return nil, errors.New(fmt.Sprintf("Invalid constraint: %s", v))
		}

		//else if v == "primary_key" {
		//			constraint = PrimaryKey()
		//}

		//else if strings.Contains(v, "foreignkey") {
		//			tableColumnPair := strings.Split(m.extractValue(v), ".")
		//			if len(tableColumnPair) != 2 {
		//				return nil, errors.New("Invalid foreign key tag. It should be foreign_key(table.column)")
		//			}
		//			// returns unformatted foreign key with parametric name
		//			constraint = ForeignKey("%s", )
		//}

		//		fmt.Println("Matched constraint: ", constraint.Name)

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

		// clean trailing spaces of tag
		rawTag = strings.Replace(f.Tag(TAG), " ", "", 1)

		constraints := []Constraint{}
		fmt.Printf("field name: %s\n", colName)
		fmt.Printf("field raw tag: %s\n", rawTag)
		fmt.Printf("field type name: %T\n", f.Value())
		fmt.Printf("field constraints: %v\n", constraints)

		if colType != "qbit.PrimaryKey" && colType != "qbit.ForeignKey" {

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

		} else if colType == "qbit.PrimaryKey" {

			table.AddConstraint(&Constraint{
				Name: fmt.Sprintf("(%s)", rawTag),
			})

		} else { // colType == "qbit.ForeignKey"

		}

		fmt.Println()

		col = Column{
			Name:        colName,
			Constraints: constraints,
			Type:        VarChar(255),
		}
		cols = append(cols, col)
		cols = append(cols, col)
	}

	//	cols, err := m.convertColumns(structs.Fields(model))
	//	if err != nil {
	//		return nil, err
	//	}

//	fmt.Println("cols: ", cols)

	return table, nil

	//	return NewTable(name, cols, []Constraint{}), nil
}
