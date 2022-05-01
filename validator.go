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
	checks := NewChecks(v.schema)
	for _, f := range v.files {
		for _, check := range checks.List {
			err := check.Do(f)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
