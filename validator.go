package csvalidator

import (
	"fmt"
)

type Validator struct {
	schema schema
	files  []File
}

func NewValidator(path, schemaPath string, firstHeader bool) (Validator, error) {
	s, err := parseSchema(schemaPath)
	if err != nil {
		return Validator{}, fmt.Errorf("failed create validator: %w", err)
	}

	files, err := NewFiles(path, firstHeader)
	if err != nil {
		return Validator{}, err
	}

	return Validator{
		schema: s,
		files:  files,
	}, nil
}

func (v Validator) Validate() error {
	for _, f := range v.files {
		missing := make([]string, 0, len(f.headers))
		for _, header := range v.schema.Columns {
			if _, ok := f.headers[header.Name]; !ok && header.Required {
				missing = append(missing, header.Name)
			}
		}
		if len(missing) > 0 {
			return fmt.Errorf("required headers are missing: %v", missing)
		}

		for i, row := range f.records {
			if i == 0 && f.firstIsHeader {
				for i := range row {
					expected := v.schema.Columns[i].Name
					got := row[i]
					if expected != got {
						return fmt.Errorf("validaton failed, column name is wrong, expected: %v, got: %v", expected, got)
					}
				}
			}
		}
	}
	return nil
}
