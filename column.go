package qbit

// function generates a table column
func Column(name string, t Type, constraints []Constraint) *column {
	return &column{
		Name:        name,
		T:           t,
		Constraints: constraints,
	}
}

type column struct {
	Name        string
	T           Type
	Constraints []Constraint
}
