package csvalidator

import "fmt"

type Checker interface {
	Do(f File) error
}

type Checks struct {
	List []Checker
}

func NewChecks(schema schema) Checks {
	var list []Checker
	list = append(list, NewColumnName(schema))
	list = append(list, NewMissingColumn(schema))
	return Checks{List: list}
}

type MissingColumn struct {
	schema schema
}

func NewMissingColumn(schema schema) MissingColumn {
	return MissingColumn{schema: schema}
}

func (mc MissingColumn) Do(f File) error {
	missing := make([]string, 0, len(f.headers))
	for _, header := range mc.schema.Columns {
		if _, ok := f.headers[header.Name]; !ok && header.Required {
			missing = append(missing, header.Name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("required headers are missing: %v", missing)
	}
	return nil
}

type ColumnName struct {
	schema schema
}

func NewColumnName(schema schema) ColumnName {
	return ColumnName{schema: schema}
}

func (cn ColumnName) Do(f File) error {
	for i, row := range f.records {
		if i == 0 && f.firstIsHeader {
			for i := range row {
				expected := cn.schema.Columns[i].Name
				got := row[i]
				if expected != got {
					return fmt.Errorf("validaton failed, column name is wrong, expected: %v, got: %v", expected, got)
				}
			}
		}
	}
	return nil
}
