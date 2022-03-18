package csvalidator

import "fmt"

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

	for _, row := range v.records {
		for i, _ := range row {
			expected := v.schema.Columns[i].Name
			got := row[i]
			if expected != got {
				return fmt.Errorf("validaton failed, column name is wrong, expected: %v, got: %v", expected, got)
			}
		}
	}
	return nil
}
