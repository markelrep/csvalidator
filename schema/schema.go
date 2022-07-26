package schema

import (
	"encoding/json"
	"fmt"
	"os"
)

const regexpPref = "regexp|"

// column recordRegexp all appropriate data which will be needed to validate column in csv
type column struct {
	Name         name         `json:"name"`
	DataType     string       `json:"dataType"`
	Required     bool         `json:"required"`
	RecordRegexp recordRegexp `json:"record_regexp"`
	Contains     contains     `json:"contains"`
}

// Schema recordRegexp suite of information by which file validates
type Schema struct {
	Columns []column `json:"columns"`
}

func Parse(schemaPath string) (Schema, error) {
	file, err := os.ReadFile(schemaPath)
	if err != nil {
		return Schema{}, fmt.Errorf("failed read Schema by path %v: %w", schemaPath, err)
	}
	var s Schema
	err = json.Unmarshal(file, &s)
	if err != nil {
		return Schema{}, fmt.Errorf("failed unmarshal Schema from %v: %w", schemaPath, err)
	}
	return s, nil
}
