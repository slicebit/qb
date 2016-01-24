package qbit

type T interface {
	String(driver string) string
}

type Type struct {
	size   int
	String func(driver string) string
}

func Integer() *T {
	return &Type{
		String: func(driver string) {
			if driver == "mysql" {

			} else if driver == "postgresql" {

			} else if driver == "sqlite" {

			}
		},
	}
}
