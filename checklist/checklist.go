package checklist

import (
	"github.com/markelrep/csvalidator/files"
	"github.com/markelrep/csvalidator/schema"
)

// Checker is common interface which should implement each check in checklist
type Checker interface {
	// Do is doing check of file
	Do(f *files.File) error
	Enqueue(row files.Row)
	Done()
}

// Checklist is list of checks which should be applied to the file
type Checklist struct {
	List []Checker
}

// NewChecklist creates a new Checklist
func NewChecklist(schema schema.Schema) Checklist {
	var list []Checker
	if len(schema.Columns) == 0 {
		return Checklist{}
	}
	list = append(list, NewColumnName(schema))
	list = append(list, NewMissingColumn(schema))
	list = append(list, NewColumnRegexpMatch(schema))
	list = append(list, NewColumnExactContain(schema))
	return Checklist{List: list}
}
