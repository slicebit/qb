package qbit

type TagMapper struct {

	// contains default, null, notnull, unique, primary_key, foreign_key(table.column), check(condition > 0)
	Constraints []string

	// contains type(size) or type parameters
	Type string
}

func (tm *TagMapper) ConvertTag(col string) *TagMapper {

	tm := &TagMapper{
		Constraints: []string{},
		Type:        "",
	}

	return tm
}
