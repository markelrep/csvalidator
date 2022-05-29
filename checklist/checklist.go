package checklist

import (
	"github.com/markelrep/csvalidator/files"
	"github.com/markelrep/csvalidator/schema"
)

type Checker interface {
	Do(f files.File) error
}

type Checklist struct {
	List []Checker
}

func NewChecklist(schema schema.Schema) Checklist {
	var list []Checker
	list = append(list, NewColumnName(schema))
	list = append(list, NewMissingColumn(schema))
	return Checklist{List: list}
}
