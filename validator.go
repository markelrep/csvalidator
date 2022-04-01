package csvalidator

import (
	"fmt"
)

type Validator struct {
	records     [][]string
	schema      schema
	firstHeader bool
	headers     map[string]struct{}
}

func NewValidator(filePath, schemaPath string, firstHeader bool) (Validator, error) {
	s, err := parseSchema(schemaPath)
	if err != nil {
		return Validator{}, fmt.Errorf("failed create validator: %w", err)
	}
	records, err := readCSV(filePath)
	if err != nil {
		return Validator{}, fmt.Errorf("failed create validator: %w", err)
	}

	req := make(map[string]struct{})
	for i, row := range records {
		if i == 0 && firstHeader {
			for _, r := range row {
				req[r] = struct{}{}
			}
		}
		break
	}

	return Validator{
		records:     records,
		schema:      s,
		firstHeader: firstHeader,
		headers:     req,
	}, nil
}

func (v Validator) Validate() error {
	missing := make([]string, 0, len(v.headers))
	for _, header := range v.schema.Columns {
		if _, ok := v.headers[header.Name]; !ok && header.Required {
			missing = append(missing, header.Name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("required headers are missing: %v", missing)
	}

	// TODO Separate header validation from records valudation
	records := v.records
	if v.firstHeader {
		header := v.records[0]
		records = v.records[1:]
		for column := range header {
			expected := v.schema.Columns[column].Name
			got := header[column]
			if expected != got {
				return fmt.Errorf("validaton failed, column name is wrong, expected: %v, got: %v", expected, got)
			}
		}
	}

	for row := range records {
		for column := range records[row] {
			logPrefix := fmt.Sprintf("[Row %v | Column %v]", row+1, column+1)
			if err := v.schema.Columns[column].Validate(records[row][column]); err != nil {
				return fmt.Errorf("%v %v", logPrefix, err)
			}
		}
	}
	return nil
}
