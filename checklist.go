package csvalidator

type Checker interface {
	Do(f File) error
}

type Check struct {
	List []Checker
}

func NewCheck(schema schema) Check {
	var list []Checker
	list = append(list, NewColumnName(schema))
	list = append(list, NewMissingColumn(schema))
	return Check{List: list}
}
